package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/mailer"
	"github.com/sikozonpc/social/internal/models"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*models.User
	Token string `json:"token"`
}

type Handlers struct {
	usersRepository *repositories.UsersRepository
	mailer          mailer.Client
	logger          *zap.SugaredLogger
	config          config.Config
	authenticator   Authenticator
}

func NewHandlers(
	usersRepository *repositories.UsersRepository,
	mailer mailer.Client,
	logger *zap.SugaredLogger,
	config config.Config,
	authenticator Authenticator,
) Handlers {
	return Handlers{
		usersRepository: usersRepository,
		mailer:          mailer,
		logger:          logger,
		config:          config,
		authenticator:   authenticator,
	}
}

// RegisterUserHandler godoc
//
//	@Summary		Registers a User
//	@Description	Registers a User
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/User [post]
func (h *Handlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	// hash the User password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	user := &models.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(passwordHash),
		IsActive: false,
		RoleID:   1,
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err = h.usersRepository.CreateAndInvite(ctx, user, hashToken, h.config.Mail.Exp)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			utils.BadRequestResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", h.config.FrontendURL, plainToken)

	// send mail
	status, err := h.mailer.SendActivationEmail(mailer.UserWelcomeTemplate, user.Username, user.Email, activationURL)
	if err != nil {
		h.logger.Errorw("error sending welcome email", "error", err)

		// rollback User creation if email fails (SAGA pattern)
		if err := h.usersRepository.DeleteByID(ctx, user.ID); err != nil {
			h.logger.Errorw("error deleting User", "error", err)
		}

		utils.InternalServerError(w, r, err)
		return
	}

	h.logger.Infow("Email sent", "status code", status)

	utils.WriteJSONMessage(w, http.StatusCreated,
		"User created successfully. Please check your email to confirm your account.")
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// LoginHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		LoginPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/login [post]
func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.usersRepository.GetOneByEmail(r.Context(), payload.Email)
	slog.Debug("GetOneByEmail result: ", "user", user, "err", err)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.UnauthorizedErrorResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	if !user.IsActive {
		utils.UnauthorizedErrorResponse(w, r, fmt.Errorf("user is not active"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utils.UnauthorizedErrorResponse(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(h.config.Auth.Token.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.config.Auth.Token.Iss,
		"aud": h.config.Auth.Token.Iss,
	}

	token, err := h.authenticator.GenerateToken(claims)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, TokenResponse{
		Token: token,
	})
}

// ConfirmHandler godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/confirm/{token} [put]
func (h *Handlers) ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := h.usersRepository.Activate(r.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.NotFoundResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	utils.WriteJSONMessage(w, http.StatusOK, "User activated successfully")
}
