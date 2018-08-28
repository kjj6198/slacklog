package slacklog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kjj6198/requests"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Text        string        `json:"text"`
	Username    string        `json:"username"`
	Channel     string        `json:"channel"`
	Attachments []*Attachment `json:"attachments"`
}

type Attachment struct {
	Fallback   string   `json:"fallback"`
	Text       string   `json:"text"`
	Pretext    string   `json:"pretext"`
	Color      string   `json:"color"`
	Fields     []*Field `json:"fields"`
	TimeStamp  int64    `json:"ts"`
	Footer     string   `json:"footer"`
	FooterIcon string   `json:"footer_icon"`
}

type SlackError struct {
	Code   int
	Reason string
}

func (e *SlackError) Error() string {
	return fmt.Sprintf("[Slack]: %d %s", e.Code, e.Reason)
}

type Field struct {
	Title    string `json:"title"`
	TitleURL string `json:"title_url"`
	Value    string `json:"value"`
	Short    bool   `json:"short"`
}

type SlackClient struct {
	URL string
}

func mapLevelToColor(level logrus.Level) string {
	switch level {
	case logrus.ErrorLevel:
	case logrus.FatalLevel:
		return "danger"
	case logrus.WarnLevel:
		return "warning"
	case logrus.InfoLevel:
		return "good"
	default:
		return "danger"
	}

	return "good"
}

func (client *SlackClient) SendMessage(msg ...*Message) error {
	config := requests.Config{
		URL:     client.URL,
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	ctx := context.Background()
	for _, message := range msg {
		data, _ := json.Marshal(message)
		result := make(map[string]interface{})
		json.Unmarshal(data, &result)
		config.Body = result
		resp, body, err := requests.Request(ctx, config)

		if resp.StatusCode >= 400 {
			return &SlackError{
				Code:   resp.StatusCode,
				Reason: body,
			}
		}

		return errors.Wrap(err, "can not send message to slack")
	}

	return nil
}
