package onylogger

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type OnyLogger struct {
	*logrus.Logger
}

type emojiFormatter struct {
	levelEmojis map[logrus.Level]string
}

const (
	colorReset   = "\033[0m"
	colorMagenta = "\033[35m"
	colorYellow  = "\033[33m"
	colorRed     = "\033[31m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
)

func (f *emojiFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Use custom emoji if provided; otherwise use the default for the log level.
	emoji, ok := entry.Data["emoji"].(string)
	if !ok {
		emoji = f.levelEmojis[entry.Level]
	}

	var colorCode string
	switch entry.Level {
	case logrus.InfoLevel:
		colorCode = colorMagenta // Magenta for Info

		if logType, exists := entry.Data["log_type"].(string); exists && logType == "input" {
			colorCode = colorReset // No Color for Input
		}
	case logrus.WarnLevel:
		colorCode = colorYellow // Yellow
	case logrus.ErrorLevel:
		colorCode = colorRed // Red
	case logrus.DebugLevel:
		colorCode = colorCyan // Cyan
	default:
		colorCode = colorReset // Default (no color)
	}

	// Apply color to the timestamp
	timestamp := colorCode + entry.Time.Format("2006-01-02 15:04:05") + "\033[0m"

	var logMsg strings.Builder
	logMsg.WriteString("[")
	logMsg.WriteString(timestamp)
	logMsg.WriteString("] ")
	logMsg.WriteString(emoji)
	logMsg.WriteString(entry.Message)

	// Only add a newline if "no_newline" is not set to true.
	if noNewline, ok := entry.Data["no_newline"].(bool); !ok || !noNewline {
		logMsg.WriteString("\n")
	}

	return []byte(logMsg.String()), nil
}

func New() *OnyLogger {
	log := logrus.New()
	log.SetFormatter(&emojiFormatter{
		levelEmojis: map[logrus.Level]string{
			logrus.InfoLevel:  "[📜] ",
			logrus.WarnLevel:  "[⚠️ ] ",
			logrus.ErrorLevel: "[❌] ",
			logrus.DebugLevel: "[🐛] ",
		},
	})
	return &OnyLogger{Logger: log}
}

// LogAndAssignInput logs the provided message with the "📝" emoji without a newline,
// then reads user input and assigns it to the provided pointer.
func (l *OnyLogger) Input(message string, userInput *string) {
	// Chain the WithField calls so both custom fields are set.
	l.WithField("log_type", "input").
		WithField("emoji", "[📝] ").
		WithField("no_newline", true).
		Info(message)
	fmt.Print(" ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	*userInput = scanner.Text()
}
