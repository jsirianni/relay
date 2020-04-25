package env

import (
    "os"

    "github.com/pkg/errors"
)

const (
    // Google queue
    envProjectID = "RELAY_GOOGLE_PROJECT_ID"

    // Global
    envLogLevel  = "RELAY_LOG_LEVEL"
)

func LogLevel() (string, error) {
    return optional(envLogLevel)
}

func GoogleProjectID() (string, error) {
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
