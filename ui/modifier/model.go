package modifier

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/louislouislouislouis/repr8ducer/ui/message"
	"github.com/louislouislouislouis/repr8ducer/ui/styles"
)

type ModifierModel struct {
	IsActive        bool
	urlDetectionMsg []string
	inputs          []textinput.Model
	cursorMode      cursor.Mode
	isFocused       bool
	urls            []string
}

func (m ModifierModel) Init() tea.Cmd {
	return nil
}

func (m ModifierModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case message.UrlDetectionMsg:
		m.urls = msg.Val
		for _, url := range m.urls {
			input := textinput.New()
			input.Placeholder = url
			m.inputs = append(m.inputs, input)
		}
	}
	return m, nil
}

func (m ModifierModel) View() string {
	text := fmt.Sprintf("We have detected the following %d urls:\n", len(m.urls))

	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	if m.isFocused {
		return styles.FocusStyle.Render(text + b.String())
	} else {
		return styles.UnfocusStyle.Render(text + b.String())
	}
}

func UpdateModifierRecap(urls []string) tea.Cmd {
	return func() tea.Msg {
		return message.UrlDetectionMsg{
			Val: urls,
		}
	}
}
