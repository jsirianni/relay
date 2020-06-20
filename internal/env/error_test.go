package env

import (
    "testing"

    "github.com/pkg/errors"
)

func TestIsEnvNotSetErrorTrue(t *testing.T) {
    err := errors.New("some error")
    err = errors.Wrap(err, notSetERR)

    if !IsEnvNotSetError(err) {
        t.Errorf("expected IsEnvNotSetError() to return true when error includes: " + notSetERR)
    }
}

func TestIsEnvNotSetErrorFalse(t *testing.T) {
    err := errors.New("some error")

    if IsEnvNotSetError(err) {
        t.Errorf("expected IsEnvNotSetError() to return false when error does not include : " + notSetERR)
    }
}
