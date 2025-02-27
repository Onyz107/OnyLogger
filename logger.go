package onylogger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//////////////////////////// OPTIONS SELECTOR ////////////////////////////

var (
	primaryColor   = lipgloss.Color("#D53F8C")
	secondaryColor = lipgloss.Color("#9F7AEA")
	accentColor    = lipgloss.Color("#6B46C1")
	bgColor        = lipgloss.Color("#1A0B2E")
	textColor      = lipgloss.Color("#F8F9FA")

	titleStyle        lipgloss.Style
	selectedItemStyle lipgloss.Style
	itemStyle         lipgloss.Style
	footerStyle       lipgloss.Style
	optionsWrapper    lipgloss.Style
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		Quit:  key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc", "quit")),
	}
}

type option struct {
	title string
	desc  string
}

type model struct {
	options     []option
	cursor      int
	selected    int
	title       string
	keys        keyMap
	choice      string
	quitting    bool
	width       int
	height      int
	centerItems bool
}

func updateStyles(width int, centerItems bool) {
	optionWidth := min(width-8, 60)
	contentWidth := optionWidth - 2 // Account for border

	titleStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Background(bgColor).
		Bold(true).
		Padding(1, 2).
		MarginBottom(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor)

	if centerItems {
		titleStyle = titleStyle.Width(optionWidth).Align(lipgloss.Center)
	}

	baseItemStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Margin(0, 1) // Add horizontal margin for border spacing

	selectedItemStyle = baseItemStyle.Copy().
		Foreground(bgColor).
		Background(primaryColor).
		Bold(true)

	itemStyle = baseItemStyle.Copy().
		Foreground(textColor)

	if centerItems {
		selectedItemStyle = selectedItemStyle.Width(contentWidth).Align(lipgloss.Center)
		itemStyle = itemStyle.Width(contentWidth).Align(lipgloss.Center)
	}

	optionsWrapper = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 0)

	if centerItems {
		optionsWrapper = optionsWrapper.Width(optionWidth)
	}

	footerStyle = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Padding(1, 2)

	if centerItems {
		footerStyle = footerStyle.Width(optionWidth).Align(lipgloss.Center)
	}
}

func initialModel(title string, optionsDescs map[string]string, centerItems bool) model {
	var opts []option
	for key, value := range optionsDescs {
		opts = append(opts, option{title: key, desc: value})
	}

	return model{
		options:     opts,
		title:       title,
		keys:        newKeyMap(),
		cursor:      0,
		selected:    -1,
		centerItems: centerItems,
	}
}

// Not for public use, this is to implement an interface
func (m model) Init() tea.Cmd {
	return nil
}

// Not for public use, this is to implement an interface
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		updateStyles(m.width, m.centerItems)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Enter):
			m.selected = m.cursor
			m.choice = m.options[m.cursor].title
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// Not for public use, this is to implement an interface
func (m model) View() string {
	if m.quitting && m.selected != -1 {
		return fmt.Sprintf("\n✨ You selected: %s ✨\n\n", m.choice)
	}

	titleContent := titleStyle.Render(m.title)
	optionsContent := ""

	for i, opt := range m.options {
		var style lipgloss.Style
		var descStyle lipgloss.Style
		text := fmt.Sprintf("  %s", opt.title)

		if m.cursor == i {
			text = fmt.Sprintf("» %s", opt.title)
			style = selectedItemStyle
			descStyle = style.Copy().
				Foreground(textColor).
				Background(accentColor).
				Italic(true)
		} else {
			style = itemStyle
			descStyle = style.Copy().
				Foreground(secondaryColor).
				Italic(true)
		}

		optionsContent += style.Render(text) + "\n"
		if opt.desc != "" {
			optionsContent += descStyle.Render(opt.desc) + "\n\n"
		}
	}

	optionsSection := optionsWrapper.Render(optionsContent)
	separator := strings.Repeat("─", max(m.width-8, 60)-4)

	helpText := footerStyle.Render("↑/↓: navigate • enter: select • q/esc: quit")
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleContent,
		optionsSection,
		"\n"+separator+"\n",
		helpText,
	)

	if m.centerItems {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().MaxWidth(min(m.width-4, 62)).Render(content),
		)
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(content)
}

// SelectOption displays a selection menu with the given title and options,
// allowing the user to navigate and select an option. The selected option's
// title is returned as a string. If no option is selected, an error is returned.
//
// Parameters:
// - title: The title of the selection menu.
// - optionsDescs: A map where keys are option titles and values are option descriptions.
// - centerItems: A boolean indicating whether to center the items in the menu.
//
// Returns:
// - string: The choice selected.
// - error: An error if no option was selected or if there was an issue running the program.
func SelectOption(title string, optionsDescs map[string]string, centerItems bool) (string, error) {
	m := initialModel(title, optionsDescs, centerItems)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	if m, ok := finalModel.(model); ok && m.selected != -1 {
		return m.choice, nil
	}
	return "", fmt.Errorf("no option was selected")
}

//////////////////////////// END OPTIONS SELECTOR ////////////////////////////

///////////////////////// SPINNER //////////////////////////////////////

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

/////////////////////////// END SPINNER ///////////////////////////

//////////////////////////// FLAGS ////////////////////////////

// DebugEnabled toggles debug logging.
var DebugEnabled bool = false

// ColorsEnabled toggles colored output.
var ColorsEnabled bool = true

// SuppressOutput when true disables all log output.
var SuppressOutput bool = false

//////////////////////////// END FLAGS ////////////////////////////

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
