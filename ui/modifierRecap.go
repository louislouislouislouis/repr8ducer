package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type modifierRecap struct {
	urls       []string
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func (m modifierRecap) Init() tea.Cmd {
	return nil
}

func (m modifierRecap) Update(msg tea.Msg) (modifierRecap, tea.Cmd) {
	switch msg := msg.(type) {
	case urlDetectionMsg:
		m.urls = msg.val
		for _, url := range m.urls {
			input := textinput.New()
			input.Placeholder = url
			m.inputs = append(m.inputs, input)
		}
	}
	return m, nil
}

func (m modifierRecap) View() string {
	text := fmt.Sprintf("We have detected the following %d urls:\n", len(m.urls))

	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	m.inputs[0].SetSuggestions([]string{"a", "b", "c"})
	return text + b.String()
}

func updateModifierRecap(urls []string) tea.Cmd {
	return func() tea.Msg {
		return urlDetectionMsg{
			val: urls,
		}
	}
}
