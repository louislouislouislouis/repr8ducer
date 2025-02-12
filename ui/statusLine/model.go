package statusline

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/louislouislouislouis/repr8ducer/ui/message"
)

type StatusLine struct {
	Text string
}

func (m StatusLine) Init() tea.Cmd {
	return nil
}

func (m StatusLine) Update(msg tea.Msg) (StatusLine, tea.Cmd) {
	switch msg := msg.(type) {
	case message.InfoMsg:
		m.Text = msg.Val
	}
	return m, nil
}

func (m StatusLine) View() string {
	return m.Text
}

func UpdateStatusLine(msg string) tea.Cmd {
	return func() tea.Msg {
		return message.InfoMsg{
			Val: msg,
		}
	}
}
