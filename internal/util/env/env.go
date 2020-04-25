package env

import (
    "os"

    "github.com/pkg/errors"
)

// global constants for frontend and forwarder services
const (
    envLogLevel = "RELAY_LOG_LEVEL"
    envTopic    = "RELAY_TOPIC"
    envSub      = "RELAY_SUBSCRIPTION"
)

// frontend
const (
    envFrontPort = "RELAY_FRONTEND_PORT"
)

// forwarder
const (
    envQueueType = "RELAY_QUEUE_TYPE"
)

// google services
const (
    // Google queue
    envGoogleProjectID = "RELAY_GOOGLE_PROJECT_ID"
)

func FrontendPort() (string, error) {
    return read(envFrontPort)
}

func GoogleProjectID() (string, error) {
    return read(envGoogleProjectID)
}

func LogLevel() (string, error) {
    return read(envLogLevel)
}

func QueueType() (string, error) {
    return read(envQueueType)
}

func Subscription() (string, error) {
    return read(envSub)
}

func Topic() (string, error) {
    return read(envTopic)
}

func read(e string) (string, error) {
    x := os.Getenv(e)
    if x == "" {
        return "", errors.New(e+": "+notSetERR)
    }
    return x, nil
}
