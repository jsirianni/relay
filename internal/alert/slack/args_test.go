package slack

import (
    "testing"
)

func newTestSlack() Slack {
    return Slack{
        HookURL: "http://slack",
        Channel: "channel",
    }
}

func TestValidateArgs(t *testing.T) {
    s := newTestSlack()

    if err := s.validateArgs("test"); err != nil {
        t.Errorf("expected validateArgs to pass")
    }

    if s.validateArgs("") == nil {
        t.Errorf("expcted an error when validateArgs was not given a message")
    }

    s.HookURL = ""
    if s.validateArgs("test") == nil {
        t.Errorf("expected an error when hook url is not set")
    }

    s = newTestSlack()

    s.Channel = ""
    if s.validateArgs("test") == nil {
        t.Errorf("expected an error when channel is not set")
    }
}
