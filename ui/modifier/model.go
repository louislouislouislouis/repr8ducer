package modifier

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/louislouislouislouis/repr8ducer/ui/commands"
	"github.com/louislouislouislouis/repr8ducer/ui/message"
	"github.com/louislouislouislouis/repr8ducer/ui/styles"
)

type ModifierModel struct {
	basePath        string
	IsActive        bool
	inputs          []textinput.Model
	inputIdxFocused int
	override        bool
	mode            modifierMode
}

type modifierMode int

const (
	override modifierMode = iota
	inputs
)

func (m ModifierModel) Init() tea.Cmd {
	return nil
}

func (m ModifierModel) getUrlsMapping() map[string]string {
	urlMapping := make(map[string]string, len(m.inputs))
	for _, input := range m.inputs {
		urlMapping[input.Placeholder] = input.Value()
	}
	return urlMapping
}

func (m ModifierModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.mode == override {
			switch msg.String() {
			case "left":
				if m.override {
					return m.handleLeftKey()
				}
				m.override = true
				return m, nil

			case "right":
				m.override = false
				return m, nil

			case "enter":
				m.mode = inputs
				var cmd tea.Cmd
				if !m.override {
					cmd = tea.Quit
				} else {
					if (len(m.inputs)) > 0 {
						m.inputs[0].Focus()
						m.inputIdxFocused = 0
						m.inputs[0].PromptStyle = styles.TextInputFocusedStyle
						m.inputs[0].TextStyle = styles.TextInputFocusedStyle
					}
				}
				return m, cmd

			case "Y", "y":
				m.override = true
				m.mode = inputs
				return m, nil
			case "N", "n":
				m.override = false
				return m, tea.Quit
			}
		} else {
			switch msg.String() {
			// Set focus to next input
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				if s == "enter" && m.inputIdxFocused == len(m.inputs)-1 {
					return m, commands.SetUrlsFromFolder(m.basePath, m.getUrlsMapping())
				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.inputIdxFocused--
				} else {
					m.inputIdxFocused++
				}

				if m.inputIdxFocused >= len(m.inputs) {
					m.inputIdxFocused = len(m.inputs) - 1
				} else if m.inputIdxFocused < 0 {
					m.inputIdxFocused = 0
				}

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.inputIdxFocused {
						// Set focused state
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = styles.TextInputFocusedStyle
						m.inputs[i].TextStyle = styles.TextInputFocusedStyle
						continue
					}
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = lipgloss.NewStyle()
					m.inputs[i].TextStyle = lipgloss.NewStyle()
				}

				return m, tea.Batch(cmds...)
			}
		}

	case message.UrlReplacementMsg:
		return m, tea.Quit
	case message.UrlDetectionMsg:
		m.basePath = msg.Val.BasePath
		m.inputs = make([]textinput.Model, len(msg.Val.Urls))
		for i, url := range msg.Val.Urls {
			input := textinput.New()
			input.Placeholder = url
			m.inputs[i] = input
		}
	}
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m ModifierModel) View() string {
	text := fmt.Sprintf("We have detected the following %d urls:\n", len(m.inputs))

	s := strings.Builder{}

	s.WriteString("Wanna override those ? ")
	yesStyle := lipgloss.NewStyle()
	noStyle := styles.ChooseButtonStyle
	if m.override {
		yesStyle = styles.ChooseButtonStyle
		noStyle = lipgloss.NewStyle()
	}
	s.WriteString(yesStyle.Render("[Yes](Y)"))
	s.WriteString("      ")
	s.WriteString(noStyle.Render("[No](N)"))
	s.WriteString("\n")
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	return styles.FocusStyle.Render(text + s.String() + b.String())
}
