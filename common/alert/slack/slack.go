package slack

import (
	"os"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jsirianni/relay/util/logger"

	"github.com/pkg/errors"
)

const envSlackDebug = "SLACK_DEBUG"
var slackDebug = false

type Slack struct {
	HookURL string
	Channel string
	Log     logger.Logger
}

type payload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (slack Slack) Message(message string) error {
	// set debug, ignore parse errors
	x := os.Getenv(envSlackDebug)
	slackDebug, _ = strconv.ParseBool(x)

	if err := slack.validateArgs(message); err != nil {
		return errors.Wrap(err, "slack configuration failed validation")
	}

	return slack.sendPayload(message)
}

func (slack Slack) sendPayload(m string) error {
	p, err := json.Marshal(payload{Channel:slack.Channel,Text:m,})
	if err != nil {
		return nil
	}

	if slackDebug {
		slack.Log.Trace("slack payload: " + string(p))
	}

	req, err := http.NewRequest("POST", slack.HookURL, bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
    	CheckRedirect: func(req *http.Request, via []*http.Request) error {
        	return http.ErrUseLastResponse
    	},
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		if b == nil {
			slack.Log.Trace("slack response body is empty")
			b = []byte("")
		}
		return errors.New("slack returned status: " + strconv.Itoa(resp.StatusCode) + " " + string(b))
	}
	return nil
}
