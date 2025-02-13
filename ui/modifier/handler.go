package modifier

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/louislouislouislouis/repr8ducer/ui/commands"
)

func (m ModifierModel) handleLeftKey() (ModifierModel, tea.Cmd) {
	return m, commands.ChangeMainModelFocus(0)
}

func (m ModifierModel) handleRightKey() (ModifierModel, tea.Cmd) {
	return m, nil
}
