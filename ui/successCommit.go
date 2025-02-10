package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func SuccessCommit(commitHash string) string {
	var (
		containerStyle = lipgloss.NewStyle().MarginLeft(2)
		titleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		messageStyle   = lipgloss.NewStyle()
		legendsStyle   = lipgloss.NewStyle().Faint(true).MarginTop(1)
	)

	return containerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Successfully commited diff"),
			messageStyle.Render(commitHash),
			legendsStyle.Render("q: quit"),
		),
	)
}
