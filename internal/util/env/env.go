package env

import (
    "os"

    "github.com/pkg/errors"
)

const (
    envProjectID = "RELAY_PROJECT_ID"
    envLogLevel  = "RELAY_LOG_LEVEL"
)

func LogLevel() (string, error) {
    return optional(envLogLevel)
}

func ProjectID() (string, error) {
    return required(envProjectID)
}

func required(e string) (string, error) {
    x := os.Getenv(e)
    if x == "" {
        return "", errors.New(e + " is not set in the environment")
    }
    return x, nil
}

func optional(e string) (string, error) {
    return os.Getenv(e), nil
}
