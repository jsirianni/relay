package env

import (
    "os"

    "github.com/jsirianni/relay/util/logger"
    "github.com/pkg/errors"
)

const (
    envProjectID = "RELAY_PROJECT_ID"
    envLogLevel  = "RELAY_LOG_LEVEL"
)

func ENVLogLevel() (string, error) {
    logLevel := os.Getenv(envLogLevel)
    if logLevel == "" {
        return logger.InfoLVL, nil
    }
    return logLevel, nil
}

func ENVProjectID() (string, error) {
    projectID := os.Getenv(envProjectID)
    if projectID == "" {
        return "", errors.New(envProjectID + " is not set in the environment")
    }
    return projectID, nil
}
