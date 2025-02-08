package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	StatusAdd    = "A"
	StatusModify = "M"
	StatusDelete = "D"
)

var (
	optionBlockStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("4")).
				Foreground(lipgloss.Color("0")).Margin(0, 1)

	titleStyle     = lipgloss.NewStyle().MarginBottom(1).Foreground(lipgloss.Color("205"))
	containerStyle = lipgloss.NewStyle().MarginLeft(3)

	baseLabelStyle  = lipgloss.NewStyle().Width(6).Margin(0, 2)
	newFileStyle    = baseLabelStyle.Foreground(lipgloss.Color("2"))
	editFileStyle   = baseLabelStyle.Foreground(lipgloss.Color("3"))
	deleteFileStyle = baseLabelStyle.Foreground(lipgloss.Color("201"))
	legendsStyle    = lipgloss.NewStyle().Faint(true).PaddingRight(1)
)

func StagedFiles(fileString string) string {
	if fileString == "" {
		return containerStyle.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("No staged files"),
			legendsStyle.Render("q: quit"),
		))
	}

	fileStrings := strings.Split(fileString, "\n")
	var builder strings.Builder

	for _, line := range fileStrings {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		status, fileName := parts[0], parts[1]

		var styled string
		switch status {
		case StatusAdd:
			styled = newFileStyle.Render("Add") + fileName
		case StatusModify:
			styled = editFileStyle.Render("Modify") + fileName
		case StatusDelete:
			styled = deleteFileStyle.Render("Delete") + fileName
		default:
			styled = line
		}

		builder.WriteString(styled + "\n")
	}

	return containerStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Staged files:"),
		builder.String(),
		legendsStyle.Render("c: continue, q/ctrl+c: quit"),
	))
}
