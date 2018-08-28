# Slacklog

very simple logrus hook for slack. Use slack attachment for better logging information.

## Getting Started

```go
hook := slacklog.CreateSlackHook(
  "https://your-slackhook-url",
	"useful-label", // label and env will send to slack as footer 
	"#your-log-channel",
	os.Getenv("ENV"),
)

log.AddHook(hook)

log.Fatal("fatal error") // send to #your-log-channel

log.WithFileds(log.Fields{
  "user": user, // format it to string by fmt.Sprint
  "time": time, // format by time.RFC1123 if type is time.Time
}).Info("some useful information")

```