package terminal

import (
    "encoding/json"

    "github.com/jsirianni/relay/internal/logger"
)

type Terminal struct {
    Log logger.Logger
}

func (t Terminal) Message(message string) error {
    log := logger.Logger{}
    if err := log.Configure(logger.InfoLVL); err != nil {
        return err
    }
    log.Info(message)
    return nil
}

func (t Terminal) Type() string {
    return "terminal"
}

func (t Terminal) Config() ([]byte, error) {
    return json.Marshal(t)
}
