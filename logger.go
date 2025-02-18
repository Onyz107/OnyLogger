package onylogger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

// Spinner represents a loading spinner.
type Spinner struct {
	stopChan chan struct{}
	message  string
}

// NewSpinner creates a new Spinner instance with a specified message.
func NewSpinner(message string) *Spinner {
	return &Spinner{
		stopChan: make(chan struct{}),
		message:  message,
	}
}

// Start begins the spinner animation.
func (s *Spinner) Start() {
	go func() {
		frames := []string{"|", "/", "-", "\\"}
		for {
			for _, frame := range frames {
				select {
				case <-s.stopChan:
					return
				default:
					fmt.Print(slogMessage("LOADING", fmt.Sprintf("%s %s\r", frame, s.message)))
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
}

// Stop stops the spinner animation.
func (s *Spinner) Stop() {
	close(s.stopChan)
	fmt.Print("\r\n") // Clear the spinner line
}

// DebugEnabled toggles debug logging.
var DebugEnabled bool = false

// ColorsEnabled toggles colored output.
var ColorsEnabled bool = true

// SuppressOutput when true disables all log output.
var SuppressOutput bool = false

// Info logs an info message with any number of arguments.
func Info(args ...any) {
	logMessage("INFO", fmt.Sprint(args...))
}

// Sinfo returns an info message with any number of arguments.
func Sinfo(args ...any) string {
	return slogMessage("INFO", fmt.Sprint(args...))
}

// Infof logs a formatted info message.
func Infof(format string, a ...any) {
	Info(fmt.Sprintf(format, a...))
}

// Sinfof returns a formatted info message.
func Sinfof(format string, a ...any) string {
	return Sinfo(fmt.Sprintf(format, a...))
}

// Warning logs a warning message with any number of arguments.
func Warning(args ...any) {
	logMessage("WARNING", fmt.Sprint(args...))
}

// Swarning returns a warning message with any number of arguments.
func Swarning(args ...any) string {
	return slogMessage("WARNING", fmt.Sprint(args...))
}

// Warningf logs a formatted warning message.
func Warningf(format string, a ...any) {
	Warning(fmt.Sprintf(format, a...))
}

// Swarningf returns a formatted warning message.
func Swarningf(format string, a ...any) string {
	return Swarning(fmt.Sprintf(format, a...))
}

// Error logs an error message with any number of arguments.
func Error(args ...any) {
	logMessage("ERROR", fmt.Sprint(args...))
}

// Serror returns an error message with any number of arguments.
func Serror(args ...any) string {
	return slogMessage("ERROR", fmt.Sprint(args...))
}

// Errorf logs a formatted error message.
func Errorf(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
}

// Serrorf returns a formatted error message.
func Serrorf(format string, a ...any) string {
	return Serror(fmt.Sprintf(format, a...))
}

// Debug logs a debug message with any number of arguments if debugging is enabled.
func Debug(args ...any) {
	if DebugEnabled {
		logMessage("DEBUG", fmt.Sprint(args...))
	}
}

// Sdebug returns a debug message with any number of arguments if debugging is enabled.
func Sdebug(args ...any) string {
	if DebugEnabled {
		return slogMessage("DEBUG", fmt.Sprint(args...))
	}
	return ""
}

// Debugf logs a formatted debug message if debugging is enabled.
func Debugf(format string, a ...any) {
	if DebugEnabled {
		Debug(fmt.Sprintf(format, a...))
	}
}

// Sdebugf returns a formatted debug message if debugging is enabled.
func Sdebugf(format string, a ...any) string {
	if DebugEnabled {
		return Sdebug(fmt.Sprintf(format, a...))
	}
	return ""
}

// Success logs a success message with any number of arguments.
func Success(args ...any) {
	logMessage("SUCCESS", fmt.Sprint(args...))
}

// Ssuccess returns a success message with any number of arguments.
func Ssuccess(args ...any) string {
	return slogMessage("SUCCESS", fmt.Sprint(args...))
}

// Successf logs a formatted success message.
func Successf(format string, a ...any) {
	Success(fmt.Sprintf(format, a...))
}

// Ssuccessf returns a formatted success message.
func Ssuccessf(format string, a ...any) string {
	return Ssuccess(fmt.Sprintf(format, a...))
}

// UserInput logs a prompt and reads user input until a new line.
func UserInput(promptMessage string, userInput *string) {
	logMessage("USERINPUT", promptMessage)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		*userInput = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading input:", err)
	}
}

// MUserInput logs a prompt and reads user input until an empty line is read.
func MUserInput(promptMessage string, userInput *string) {
	logMessage("USERINPUT", promptMessage)
	scanner := bufio.NewScanner(os.Stdin)

	var input string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		input += line + "\n"
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading input:", err)
		return
	}

	*userInput = input
}

// InitLogger configures log output.
func InitLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0) // Remove default timestamp
}

const (
	colorReset   = "\033[0m"
	colorMagenta = "\033[35m"
	colorYellow  = "\033[33m"
	colorRed     = "\033[31m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
)

// Emojis for a thematic touch
const (
	emojiInfo    = "📜"
	emojiWarning = "⚠️"
	emojiError   = "❌"
	emojiDebug   = "🐛"
	emojiSuccess = "✅"
	emojiLoading = "⏳"
	emojiInput   = "📝"
)

// logMessage logs a message with a specific level and formats it.
func logMessage(level, message string) {
	if SuppressOutput {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var color, emoji string

	switch level {
	case "INFO":
		color = colorMagenta
		emoji = emojiInfo
	case "WARNING":
		color = colorYellow
		emoji = emojiWarning
	case "ERROR":
		color = colorRed
		emoji = emojiError
	case "DEBUG":
		color = colorCyan
		emoji = emojiDebug
	case "SUCCESS":
		color = colorGreen
		emoji = emojiSuccess
	case "USERINPUT":
		color = colorReset
		emoji = emojiInput
	case "LOADING":
		color = colorCyan
		emoji = emojiLoading
	default:
		color = colorReset
	}

	if !ColorsEnabled {
		color = ""
	}

	formattedMessage := fmt.Sprintf("%s[%s]%s [%s] %s", color, timestamp, colorReset, emoji, message)
	if level != "USERINPUT" {
		log.Println(formattedMessage)
	} else {
		fmt.Print(formattedMessage)
	}
}

// slogMessage logs a message with a specific level and formats it.
func slogMessage(level, message string) string {
	if SuppressOutput {
		return ""
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var emoji string

	switch level {
	case "INFO":
		emoji = emojiInfo
	case "WARNING":
		emoji = emojiWarning
	case "ERROR":
		emoji = emojiError
	case "DEBUG":
		emoji = emojiDebug
	case "SUCCESS":
		emoji = emojiSuccess
	case "USERINPUT":
		emoji = emojiInput
	case "LOADING":
		emoji = emojiLoading
	default:
	}

	formattedMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, emoji, message)
	return formattedMessage
}
