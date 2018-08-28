package slacklog

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/kjj6198/requests"
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
		return "danger"
	case logrus.FatalLevel:
		return "#fe6565"
	case logrus.WarnLevel:
		return "warning"
	case logrus.InfoLevel:
		return "good"
	default:
		return "danger"
	}
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

type SlackLogrusHook struct {
	Client *SlackClient
	Label  string
}

func (hook *SlackLogrusHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

func (hook *SlackLogrusHook) Fire(entry *logrus.Entry) error {
	var fields []*Field
	for key, val := range map[string]interface{}(entry.Data) {
		if timeVal, ok := val.(time.Time); ok {
			val = timeVal.Format(time.RFC1123)
		}

		fields = append(fields, &Field{
			Title: key,
			Value: fmt.Sprint(val),
			Short: len(fmt.Sprint(val)) < 25,
		})
	}

	message := &Message{
		Text: fmt.Sprintf("*[%s] %s*", strings.ToUpper(entry.Level.String()), entry.Message),
		Attachments: []*Attachment{
			&Attachment{
				Text:      entry.Message,
				Color:     mapLevelToColor(entry.Level),
				Fields:    fields,
				TimeStamp: entry.Time.Unix(),
				Footer:    hook.Label,
			},
		},
	}

	return errors.Wrap(hook.Client.SendMessage(message), "slack logging error")
}

func CreateSlackHook(webhookURL string, label string) logrus.Hook {
	return &SlackLogrusHook{
		Client: &SlackClient{URL: webhookURL},
		Label:  label,
	}
}
