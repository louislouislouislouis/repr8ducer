package ui

import "github.com/charmbracelet/lipgloss"

var (
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")). // Couleur blanche
			Background(lipgloss.Color("#F0F")).    // Fond bleu
			Bold(true).
			Padding(1, 2)
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0")).Bold(true)
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
)
