package main

import (
    "os"

    "github.com/jsirianni/relay/common/alert"

    "github.com/pkg/errors"
)

const (
    typeSlack = "slack"
    envDestType = "RELAY_DEST_TYPE"
    envSlackHookURL = "RELAY_SLACK_HOOK_URL"
    envSlackChannel = "RELAY_SLACK_CHANNEL"
)

func initDest() (alert.Alert, error) {
    destType := os.Getenv(envDestType)
    if destType == "" {
        return nil, errors.New("destination type is not set")
    }

    if destType == typeSlack {
        hookURL := os.Getenv(envSlackHookURL)
        channel := os.Getenv(envSlackChannel)
        return alert.NewSlack(hookURL, channel)
    }

    return nil, errors.New(destType + " is not supported")
}
