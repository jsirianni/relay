package alert

import (
	"github.com/pkg/errors"

	"github.com/jsirianni/relay/common/alert/slack"
)

type Alert interface {
	// Message takes a message as a string and sends it
	// to the configured destination
	Message(message string) error
}

func NewSlack(hookURL, channel string) (Alert, error) {
	if hookURL == "" {
		return nil, errors.New("slack hook url is not set")
	}

	if channel == "" {
		return nil, errors.New("slack channel is not set")
	}

	var a Alert = slack.Slack{hookURL,channel}
	return a, nil
}
