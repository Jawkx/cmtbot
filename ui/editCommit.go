package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

func EditCommit(textArea textarea.Model) string {
	var (
		containerStyle = lipgloss.NewStyle().MarginLeft(2)
		legendsStyle   = lipgloss.NewStyle().Faint(true).MarginTop(1)
	)

	return containerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			textArea.View(),
			legendsStyle.Render("ctrl+s: commit, esc: back"),
		),
	)
}
