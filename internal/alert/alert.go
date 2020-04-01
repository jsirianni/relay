package alert

import (
	"github.com/pkg/errors"

	"github.com/jsirianni/relay/internal/util/logger"
	"github.com/jsirianni/relay/internal/alert/slack"
	"github.com/jsirianni/relay/internal/alert/terminal"
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
