package terminal

import (
    "testing"
)

func TestStdout(t *testing.T) {
    a := Terminal{}
    if err := a.Message("hello"); err != nil {
        t.Errorf("expected terminal.Message to return a nil error, got: " + err.Error())
    }
}

func TestType(t *testing.T) {
    a := Terminal{}
    if a.Type() != "terminal" {
        t.Errorf("expected type to be terminal")
    }
}

func TestConfig(t *testing.T) {
    // test if conig returns a nil error when log is not set
    a := Terminal{}
    if _, err := a.Config(); err != nil {
        t.Errorf(err.Error())
    }
}
