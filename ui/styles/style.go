package styles

import "github.com/charmbracelet/lipgloss"

var (
	docStyle = lipgloss.
			NewStyle()

	BigTitleStyle = lipgloss.
			NewStyle().
			Align(lipgloss.Center)

	FocusStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	UnfocusStyle = FocusStyle.
			BorderForeground(lipgloss.Color("#bababa"))

	TextInputFocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ChooseButtonStyle = lipgloss.NewStyle().Background(lipgloss.Color("205"))
)
