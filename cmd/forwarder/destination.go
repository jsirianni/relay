package main

import (
    "os"
    "strings"

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

// configure destination types and set the default
func initDests() error {
    f.Alert = make(map[string]alert.Alert)

    d := os.Getenv(envDestType)
    if d == "" {
        return errors.New("default destination type is not set")
    }

    destTypes := strings.Split(d, ",")
    if len(destTypes) == 0 {
        return errors.New("no destination types specified")
    }

    // set default if only only destination is specified
    if len(destTypes) == 1 {
        f.DefaultAlert = destTypes[0]
    }
    // only check destType global if multiple alertt types are requested
    if len(destTypes) > 1 {
        if destType == "" {
            return errors.New("--dest-type must be set if multiple destination types are specified, otherwise a default cannot be determined")
        }
    }

    for _, d := range destTypes {
        a, err := initDest(d)
        if err != nil {
            return err
        }
        f.Alert[d] = a
    }

    return nil
}

func initDest(d string) (alert.Alert, error) {
    if d == typeSlack {
        hookURL := os.Getenv(envSlackHookURL)
        channel := os.Getenv(envSlackChannel)
        return alert.NewSlack(hookURL, channel, f.Log)
    }

    if d == typeTerm {
        return alert.NewTerminal(f.Log)
    }

    if d == typeSendGrid {
        from := os.Getenv(envSendGridFromEmail)
        to   := os.Getenv(envSendGridToEmail)
        apiKey := os.Getenv(envSendGridAPIKey)
        return alert.NewSendGrid(from, to, apiKey, f.Log)
    }

    return nil, errors.New(d + " is not supported")
}
