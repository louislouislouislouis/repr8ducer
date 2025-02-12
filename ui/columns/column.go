package columns

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/louislouislouislouis/repr8ducer/ui/common"
	"github.com/louislouislouislouis/repr8ducer/ui/message"
	"github.com/louislouislouislouis/repr8ducer/ui/styles"
)

type Column struct {
	width         int
	height        int
	isFocused     bool
	isInitialized bool
	list          list.Model
	current       string
	spinner       spinner.Model
}

func NewColumn(
	width, height int,
	isFocus, isInitialized bool,
	list list.Model,
	current string,
	spinner spinner.Model,
) Column {
	return Column{
		width:         width,
		height:        height,
		isFocused:     isFocus,
		isInitialized: isInitialized,
		current:       current,
		spinner:       spinner,
		list:          list,
	}
}

func (c *Column) SetFocus(val bool) {
	c.isFocused = val
}

func (c Column) Init() tea.Cmd {
	return nil
}

func (c Column) View() string {
	var text string
	if c.isInitialized {
		text = c.list.View()
	} else {
		text = c.spinner.View()
	}

	if c.isFocused {
		return styles.FocusStyle.Width(c.width).Render(text)
	}
	return styles.UnfocusStyle.Width(c.width).Render(text)
}

func (c Column) Update(msg tea.Msg) (Column, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		c.width = (msg.Width / 3) - 2
		c.height = msg.Height - 8
		c.list.SetSize(c.width, c.height)

	case message.ListUpdateMsg:
		c.isInitialized = true
		cmds := tea.Batch(
			c.list.SetItems(msg.Val),
			c.list.NewStatusMessage(msg.StatusTxt),
		)
		selectedIdx := 0
		for idx, item := range msg.Val {
			if item.(common.DisplayedItem).Title() == msg.PreSelectedValue {
				selectedIdx = idx
				break
			}
		}
		c.list.Select(selectedIdx)
		return c, cmds

	case spinner.TickMsg:
		var cmd tea.Cmd
		c.spinner, cmd = c.spinner.Update(msg)
		return c, cmd

	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}
