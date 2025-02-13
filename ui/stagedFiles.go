package ui

import (
	"strings"

	"github.com/Jawkx/cmtbot/git"
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
	deleteFileStyle = baseLabelStyle.Foreground(lipgloss.Color("1"))
	legendsStyle    = lipgloss.NewStyle().Faint(true).PaddingRight(1)
)

func StagedFiles(stagedFiles []git.StagedFile) string {
	if len(stagedFiles) == 0 {
		return containerStyle.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("No staged files"),
			legendsStyle.Render("q: quit"),
		))
	}

	var builder strings.Builder

	for _, file := range stagedFiles {
		var styled string
		switch file.Status {
		case StatusAdd:
			styled = newFileStyle.Render("Add") + file.Path
		case StatusModify:
			styled = editFileStyle.Render("Modify") + file.Path
		case StatusDelete:
			styled = deleteFileStyle.Render("Delete") + file.Path
		default:
			styled = file.Path
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
