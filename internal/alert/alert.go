package alert

import (
	"github.com/pkg/errors"

	"github.com/jsirianni/relay/internal/util/logger"
	"github.com/jsirianni/relay/internal/alert/slack"
	"github.com/jsirianni/relay/internal/alert/terminal"
	"github.com/jsirianni/relay/internal/alert/sendgrid"
)

type Alert interface {
	// Message takes a message as a string and sends it
	// to the configured destination
	Message(message string) error

	// The alert type (Slack, Terminal, etc)
	Type() string

	// Config returns a json []byte value
	Config() ([]byte, error)
}

func NewSlack(hookURL, channel string, l logger.Logger) (Alert, error) {
	if hookURL == "" {
		return nil, errors.New("slack hook url is not set")
	}

	if channel == "" {
		return nil, errors.New("slack channel is not set")
	}

	return slack.Slack{hookURL,channel,l}, nil
}

func NewTerminal(l logger.Logger) (Alert, error) {
	return terminal.Terminal{l}, nil
}

func NewSendGrid(fromEmail, toEmail, apiKey string, l logger.Logger) (Alert, error) {
	if fromEmail == "" {
		return nil, errors.New("sendgrid from email is not set")
	}

	if toEmail == "" {
		return nil, errors.New("sendgrid to email is not set")
	}

	if apiKey == "" {
		return nil, errors.New("sendgrid api key not set")
	}

	if l.Configured() == false {
		return nil, errors.New("sendgrid logger config is nil")
	}

	return sendgrid.SendGrid{
		FromEmail: fromEmail,
		ToEmail:   toEmail,
		APIKey:    apiKey,
		Log:       l,
	}, nil
}
