package flockhook

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    timeLayout "store-api/pkg/data/constant"
    "store-api/pkg/errs"
    "store-api/pkg/httpclient"

    "github.com/sirupsen/logrus"
)

// FlockHook to send logs error to Flock
type Hook struct {
    Channel    string
    HTTPClient httpclient.Client
    levels     []logrus.Level
}

// NewPapertrailHook creates a UDP hook to be added to an instance of logger.
func NewFlockHook(hook *Hook) (*Hook, error) {
    var err error
    return hook, err
}

// Fire is called when a log event is fired.
func (hook *Hook) Fire(entry *logrus.Entry) error {
    date := time.Now().Format(timeLayout.DefaultDatetimeLayout)
    data, _ := json.Marshal(entry.Data)

    if isSend, ok := entry.Data["flock"].(bool); ok {
        //fmt.Println(isSend)
        if !isSend {
            return nil
        }
    }

    d := string(data)
    if d == "{}" && entry.Level != logrus.DebugLevel {
        // Get more clue
        stack := errs.GetFileAndFuncForLogrus()
        d = fmt.Sprintf("{file: %s}", stack)
        fmt.Printf("Logrus: %v\n", d)
    }

    msg := entry.Message
    if len(entry.Message) > 301 {
        msg = entry.Message[:300]
    }
    payload := fmt.Sprintf("Log%s [%s]: %s - %s (%s) --- ",
        hook.LevelToIcon(entry.Level), date, msg, d, os.Getenv("APP_NAME"))

    hook.SendMsgFlock(payload)
    return nil
}

func (hook *Hook) LevelToIcon(level logrus.Level) string {
    switch level {
    case logrus.PanicLevel:
        return "Panic :exclamation:"
    case logrus.FatalLevel:
        return "FatalError :red_circle: :red_circle: :red_circle:"
    case logrus.ErrorLevel:
        return "Error :red_circle:"
    case logrus.WarnLevel:
        return "Warn :warning:"
    case logrus.InfoLevel:
        return "Info :information_source:"
    case logrus.DebugLevel:
        return "Debug :footprints:"
    default:
        return ""
    }
}

// SetLevels specify nessesary levels for this hook.
func (hook *Hook) SetLevels(lvs []logrus.Level) {
    hook.levels = lvs
}

// Levels returns the available logging levels.
func (hook *Hook) Levels() []logrus.Level {

    if hook.levels == nil {
        return []logrus.Level{
            logrus.PanicLevel,
            logrus.FatalLevel,
            logrus.ErrorLevel,
            logrus.WarnLevel,
            //logrus.InfoLevel, // Except INFO
            logrus.DebugLevel,
        }
    }

    return hook.levels
}

type FlockLogrusResponse struct {
    UID string `json:"uid"`
}

//SendMsgFlock send msg to default Flock
func (h *Hook) SendMsgFlock(message string) int {

    data := map[string]string{"text": message}
    headers := map[string]string{"Content-Type": "application/json"}

    code, _ := h.HTTPClient.Post(h.Channel, data, headers, &FlockLogrusResponse{})
    return code
}
