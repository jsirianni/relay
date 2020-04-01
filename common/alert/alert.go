package alert

import (
	"github.com/pkg/errors"

	"github.com/jsirianni/relay/util/logger"
	"github.com/jsirianni/relay/common/alert/slack"
	"github.com/jsirianni/relay/common/alert/terminal"
)

type Alert interface {
	// Message takes a message as a string and sends it
	// to the configured destination
	Message(message string) error
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

func NewTerminal() (Alert, error) {
	return terminal.Terminal{}, nil
}
