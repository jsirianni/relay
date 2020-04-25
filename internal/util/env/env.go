package env

import (
    "os"

    "github.com/pkg/errors"
)

// global constants for frontend and forwarder services
const (
    ENVLogLevel = "RELAY_LOG_LEVEL"
    ENVTopic    = "RELAY_TOPIC"
    ENVSub      = "RELAY_SUBSCRIPTION"
)

// frontend
const (
    ENVFrontPort = "RELAY_FRONTEND_PORT"
)

// forwarder
const (
    ENVQueueType = "RELAY_QUEUE_TYPE"
)

// google services
const (
    // Google queue
    ENVGoogleProjectID = "RELAY_GOOGLE_PROJECT_ID"
)

func FrontendPort() (string, error) {
    return read(ENVFrontPort)
}

func GoogleProjectID() (string, error) {
    return read(ENVGoogleProjectID)
}

func LogLevel() (string, error) {
    return read(ENVLogLevel)
}

func QueueType() (string, error) {
    return read(ENVQueueType)
}

func Subscription() (string, error) {
    return read(ENVSub)
}

func Topic() (string, error) {
    return read(ENVTopic)
}

func read(e string) (string, error) {
    x := os.Getenv(e)
    if x == "" {
        return "", errors.New(e+": "+notSetERR)
    }
    return x, nil
}
