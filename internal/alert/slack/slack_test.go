package slack

import (
    "testing"
    "encoding/json"
)

func TestType(t *testing.T) {
    s := Slack{}
    if s.Type() != "slack" {
        t.Errorf("expected type to be 'slack'")
    }
}

func TestConfig(t *testing.T) {
    s := Slack{
        HookURL: "test",
        Channel: "test",
    }

    b, err := s.Config()
    if err != nil {
        t.Errorf(err.Error())
    }

    if err := json.Unmarshal(b, &s); err != nil {
        t.Errorf(err.Error())
        return
    }

    if s.HookURL != "test" {
        t.Errorf("expected hookurl to be test")
    }
    if s.Channel != "test" {
        t.Errorf("expected channel to be test")
    }
}
