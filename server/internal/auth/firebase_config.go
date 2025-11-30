package auth

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type Client interface {
	VerifyIDToken(ctx context.Context, idToken string) (*Token, error)
}

type Token struct {
	AuthTime int64                  `json:"auth_time"`
	Issuer   string                 `json:"iss"`
	Audience string                 `json:"aud"`
	Expires  int64                  `json:"exp"`
	IssuedAt int64                  `json:"iat"`
	Subject  string                 `json:"sub,omitempty"`
	UID      string                 `json:"uid,omitempty"`
	Claims   map[string]interface{} `json:"-"`
}

type client struct {
	firebaseClient *auth.Client
}

func (c *client) VerifyIDToken(ctx context.Context, idToken string) (*Token, error) {
	token, err := c.firebaseClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return &Token{
		AuthTime: token.AuthTime,
		Issuer:   token.Issuer,
		Audience: token.Audience,
		Expires:  token.Expires,
		IssuedAt: token.IssuedAt,
		Subject:  token.Subject,
		UID:      token.UID,
		Claims:   token.Claims,
	}, nil
}

func Setup(ctx context.Context) Client {
	emulatorHost := os.Getenv("FIREBASE_AUTH_EMULATOR_HOST")
	if emulatorHost != "" {
		fmt.Printf("--- Firebase Auth Emulator Detected: Connecting to %s ---\n", emulatorHost)
		// Use the project ID the emulator is configured with.
		// This should match the 'aud' claim in the JWT token.
		config := &firebase.Config{ProjectID: "demo-no-project"}
		app, err := firebase.NewApp(ctx, config)
		if err != nil {
			panic(fmt.Sprintf("error initializing app for emulator: %v\n", err))
		}

		firebaseClient, err := app.Auth(ctx)
		if err != nil {
			panic(fmt.Sprintf("error getting Auth client for emulator: %v\n", err))
		}
		return &client{
			firebaseClient: firebaseClient,
		}
	}

	fmt.Println("--- No Firebase Auth Emulator Detected: Connecting to Production ---")
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("error initializing app: %v\n", err))
	}

	firebaseClient, err := app.Auth(ctx)
	if err != nil {
		panic(fmt.Sprintf("error getting Auth client: %v\n", err))
	}

	return &client{
		firebaseClient: firebaseClient,
	}
}
