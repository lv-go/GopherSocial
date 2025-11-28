package mailer

import "log/slog"

type mockMailerClient struct {
}

func NewMockMailerClient() Client {
	return mockMailerClient{}
}

func (m mockMailerClient) SendActivationEmail(templateFile, username, email, activationUrl string) (int, error) {
	slog.Info("Sending mock email",
		"templateFile", templateFile,
		"username", username,
		"email", email,
		"activationUrl", activationUrl)
	return 200, nil
}
