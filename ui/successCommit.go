package ui

import "github.com/charmbracelet/lipgloss"

func SuccessCommit() string {
	var (
		containerStyle = lipgloss.NewStyle().MarginLeft(2)
		messageStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		legendsStyle   = lipgloss.NewStyle().Faint(true).MarginTop(2)
	)

	return containerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			messageStyle.Render("Successfully commited"),
			legendsStyle.Render("q: quit"),
		),
	)
}
