package alert

import (
    "testing"

    "github.com/jsirianni/relay/internal/util/logger"
)

func newLogger() logger.Logger {
    l := logger.Logger{}
    if err := l.Configure("trace"); err != nil {
        panic(err)
    }
    return l
}

func TestNewSlack(t *testing.T) {
    if _, err := NewSlack("", "test", newLogger()); err == nil {
        t.Errorf("expected an error when slack hook url is empty, got nil")
    }

    if _, err := NewSlack("test", "", newLogger()); err == nil {
        t.Errorf("expected an error when slack channel is empty, got nil")
    }

    l := logger.Logger{}
    if _, err := NewSlack("test", "test", l); err == nil {
        t.Errorf("expected an error when logger is not configured, got nil")
    }

    a, err := NewSlack("test", "test", newLogger())
    if err != nil {
        t.Errorf(err.Error())
    }

    if a.Type() != "slack" {
        t.Errorf("expected alert interface Type() method to return 'slack', got: " + a.Type())
    }
}

func TestNewTerminal(t *testing.T) {
    l := logger.Logger{}
    if _, err := NewTerminal(l); err == nil {
        t.Errorf("expected an error when logger is not configured, got nil")
    }

    a, err := NewTerminal(newLogger())
    if err != nil {
        t.Errorf(err.Error())
    }

    if a.Type() != "terminal" {
        t.Errorf("expected alert interface Type() method to return 'terminal', got: " + a.Type())
    }
}

func TestNewSendGrid(t *testing.T) {
    if _, err := NewSendGrid("test", "test", "", newLogger()); err == nil {
        t.Errorf("expected an error when sendgrid apiKey is empty, got nil")
    }

    if _, err := NewSendGrid("test", "", "test", newLogger()); err == nil {
        t.Errorf("expected an error when sendgrid to email is empty, got nil")
    }

    if _, err := NewSendGrid("", "test", "test", newLogger()); err == nil {
        t.Errorf("expected an error when sendgrid from email is empty, got nil")
    }

    l := logger.Logger{}
    if _, err := NewSendGrid("test", "test", "test", l); err == nil {
        t.Errorf("expected an error when logger is not configured, got nil")
    }

    a, err := NewSendGrid("test", "test", "test", newLogger())
    if err != nil {
        t.Errorf(err.Error())
    }

    if a.Type() != "sendgrid" {
        t.Errorf("expected alert interface Type() method to return 'sendgrid', got: " + a.Type())
    }
}
