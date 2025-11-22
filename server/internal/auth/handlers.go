package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sikozonpc/social/internal/config"
	"github.com/sikozonpc/social/internal/mailer"
	"github.com/sikozonpc/social/internal/models"
	"github.com/sikozonpc/social/internal/repositories"
	"github.com/sikozonpc/social/internal/store"
	"github.com/sikozonpc/social/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
		Role: models.Role{
			Name: "User",
		},
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err = h.usersRepository.CreateAndInvite(ctx, user, hashToken, h.config.Mail.Exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			utils.BadRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			utils.BadRequestResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	activationURL := fmt.Sprintf("%s/confirm/%s", h.config.FrontendURL, plainToken)

	isProdEnv := h.config.Env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// send mail
	status, err := h.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
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

	utils.WriteJSON(w, http.StatusCreated, userWithToken)
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type LoginResponse struct {
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
	if err != nil {
		switch err {
		case store.ErrNotFound:
			utils.UnauthorizedErrorResponse(w, r, err)
		default:
			utils.InternalServerError(w, r, err)
		}
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

	utils.WriteJSON(w, http.StatusOK, LoginResponse{
		Token: token,
	})
}
