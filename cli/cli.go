package cli

import (
	"dodo-lang/evaluator"
	"dodo-lang/lexer"
	"dodo-lang/object"
	"dodo-lang/parser"
	"fmt"
	"os/user"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Cli struct {
	env         *object.Environment
	verboseMode bool
}

func New(verbose bool) *Cli {
	env := object.NewEnvironment()
	return &Cli{env: env, verboseMode: verbose}
}

func (c *Cli) Init() *tea.Program {
	p := tea.NewProgram(initialModel(c.env, c.verboseMode))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error initialzing CLI: %v", err)
	}

	return p
}

type (
	errMsg error
)

type model struct {
	input        string
	output       string
	prompt       string
	history      []string
	historyPos   int
	cursorPos    int
	textInput    textinput.Model
	err          error
	parserErrors []string
	verboseMode  bool
	env          *object.Environment
}

var suggestions = []string{"let", "if", "else", "true", "false", "fn", "return", "()"}

func initialModel(env *object.Environment, verbose bool) model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 1024
	ti.Width = 128
	ti.ShowSuggestions = true
	ti.SetSuggestions(suggestions)

	return model{
		output:      "",
		history:     []string{},
		historyPos:  0,
		cursorPos:   0,
		textInput:   ti,
		err:         nil,
		verboseMode: verbose,
		env:         env,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.history = append(m.history, m.textInput.Value())
			evaluated, err := evalInput(m.textInput.Value(), m.env)

			// Reset previous
			m.textInput.Reset()
			m.output = ""

			if err != nil {
				m.parserErrors = err
				return m, nil
			} else {
				// Clear errors
				m.parserErrors = make([]string, 0)
			}

			m.output = evaluated
			return m, nil
		case tea.KeyUp:
			if len(m.history) == 0 {
				return m, nil
			}

			m.historyPos = (m.historyPos - 1 + len(m.history)) % len(m.history)
			m.textInput.SetValue(m.history[m.historyPos])

			return m, nil
		case tea.KeyDown:
			if m.historyPos >= len(m.history)-1 {
				m.textInput.Reset()
				return m, nil
			}

			m.historyPos = m.historyPos + 1
			m.textInput.SetValue(m.history[m.historyPos])

			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyRunes:
			if len(msg.Runes) > 0 {
				cursorPos := m.textInput.Position()
				currentValue := m.textInput.Value()
				var newValue = ""
				switch msg.Runes[0] {
				case '(', '[', '{':
					tk := msg.Runes[0]
					closingTags := map[rune]string{'(': ")", '[': "]", '{': "}"}
					if closingTag, ok := closingTags[tk]; ok {
						newValue = currentValue[:cursorPos] + string(tk) + closingTag + currentValue[cursorPos:]
						m.textInput.SetValue(newValue)
						m.textInput.SetCursor(cursorPos + 1)
						return m, cmd
					}
				}

			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	greeting := fmt.Sprintf("\nHello %s! You are currently running the Dodo programming language.\n", user.Username)
	styledOutput := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(m.output)

	styledParserErrors := ""

	if m.verboseMode {
		styledParserErrors = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3022")).Render(strings.Join(m.parserErrors, "\n"))
	}

	return fmt.Sprintf(
		"%s%s%s\n%s\n%s\n",
		greeting,
		"Now, type some commands below :)\n",
		styledParserErrors,
		m.textInput.View(),
		styledOutput,
	)
}

func evalInput(input string, env *object.Environment) (string, []string) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return "", p.GetParserErrors()
	}

	switch evaluated := evaluator.Eval(program, env).(type) {
	case *object.Error:
		if evaluated != nil {
			return evaluated.Inspect(), nil
		}
	default:
		if evaluated != nil {
			return evaluated.Inspect(), nil
		}
	}

	return "", nil
}
