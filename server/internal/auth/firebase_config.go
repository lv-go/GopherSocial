package auth

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type Client = auth.Client

func Setup(ctx context.Context) *Client {
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

		client, err := app.Auth(ctx)
		if err != nil {
			panic(fmt.Sprintf("error getting Auth client for emulator: %v\n", err))
		}
		return client
	}

	fmt.Println("--- No Firebase Auth Emulator Detected: Connecting to Production ---")
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("error initializing app: %v\n", err))
	}

	client, err := app.Auth(ctx)
	if err != nil {
		panic(fmt.Sprintf("error getting Auth client: %v\n", err))
	}

	return client
}
