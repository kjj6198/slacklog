package slacklog

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SlackLogrusHook struct {
	Client  *SlackClient
	Label   string
	Channel string
	Env     string
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
		Channel: hook.Channel,
		Text:    fmt.Sprintf("*[%s] %s*", strings.ToUpper(entry.Level.String()), entry.Message),
		Attachments: []*Attachment{
			&Attachment{
				Text:      entry.Message,
				Color:     mapLevelToColor(entry.Level),
				Fields:    fields,
				TimeStamp: entry.Time.Unix(),
				Footer:    fmt.Sprintf("%s - %s", hook.Env, hook.Label),
			},
		},
	}

	return errors.Wrap(hook.Client.SendMessage(message), "slack logging error")
}

func CreateSlackHook(
	webhookURL string,
	label string,
	channel string,
	env string,
) logrus.Hook {
	return &SlackLogrusHook{
		Client:  &SlackClient{URL: webhookURL},
		Label:   label,
		Channel: channel,
		Env:     env,
	}
}
