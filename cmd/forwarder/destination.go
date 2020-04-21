package main

import (
    "os"

    "github.com/jsirianni/relay/internal/alert"

    "github.com/pkg/errors"
)

const (
    typeSlack    = "slack"
    typeTerm     = "terminal"
    typeSendGrid = "sendgrid"

    envDestType = "RELAY_DEST_TYPE"
    envSlackHookURL = "RELAY_SLACK_HOOK_URL"
    envSlackChannel = "RELAY_SLACK_CHANNEL"
    envSendGridFromEmail = "RELAY_SENDGRID_FROM_EMAIL"
    envSendGridToEmail   = "RELAY_SENDGRID_TO_EMAIL"
    envSendGridAPIKey    = "RELAY_SENDGRID_API_KEY"
)

func initDest() (alert.Alert, error) {
    destType := os.Getenv(envDestType)
    if destType == "" {
        return nil, errors.New("destination type is not set")
    }

    if destType == typeSlack {
        hookURL := os.Getenv(envSlackHookURL)
        channel := os.Getenv(envSlackChannel)
        return alert.NewSlack(hookURL, channel, f.Log)
    }

    if destType == typeTerm {
        return alert.NewTerminal(f.Log)
    }

    if destType == typeSendGrid {
        from := os.Getenv(envSendGridFromEmail)
        to   := os.Getenv(envSendGridToEmail)
        apiKey := os.Getenv(envSendGridAPIKey)
        return alert.NewSendGrid(from, to, apiKey, f.Log)
    }

    return nil, errors.New(destType + " is not supported")
}
