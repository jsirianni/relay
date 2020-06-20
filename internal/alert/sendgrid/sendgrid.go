package sendgrid

import (
    "encoding/json"

    "github.com/jsirianni/relay/internal/logger"

    realsendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
    FromEmail string
    ToEmail   string
    APIKey    string
    Log       logger.Logger
}

func (s SendGrid) Message(message string) error {
    from := mail.NewEmail("relay", s.FromEmail)
	subject := "relay alert"
	to := mail.NewEmail("relay", s.ToEmail)
	plainTextContent := message
	htmlContent := ""

	m := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := realsendgrid.NewSendClient(s.APIKey)
	r, err := client.Send(m)
	if err != nil {
		return err
	}
    if s.Log.Level() == "trace" {
        b, err := json.Marshal(r)
        if err != nil {
            s.Log.Error(err)
            return nil
        }

        s.Log.Trace(string(b))
    }
    return nil
}

func (s SendGrid) Type() string {
    return "sendgrid"
}

func (s SendGrid) Config() ([]byte, error) {
    return json.Marshal(s)
}
