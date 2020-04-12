package sendgrid

import (
    "testing"
    "encoding/json"
)

func TestType(t *testing.T) {
    s := SendGrid{}
    if s.Type() != "sendgrid" {
        t.Errorf("expected type to be sendgrid")
    }
}

func TestConfig(t *testing.T) {
    s := SendGrid{
        FromEmail: "test",
        ToEmail: "test",
        APIKey: "test",
    }

    b, err := s.Config()
    if err != nil {
        t.Errorf(err.Error())
    }

    if err := json.Unmarshal(b, &s); err != nil {
        t.Errorf(err.Error())
        return
    }

    if s.FromEmail != "test" {
        t.Errorf("expected from email to be test")
    }
    if s.ToEmail != "test" {
        t.Errorf("expected to email to be test")
    }
    if s.APIKey != "test" {
        t.Errorf("expected api key to be test")
    }
}
