package ui

import (
	"strings"

	"github.com/Jawkx/cmtbot/git"
	"github.com/Jawkx/cmtbot/model"
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

	baseLabelStyle      = lipgloss.NewStyle().Width(11).Margin(0, 2)
	newFileStyle        = baseLabelStyle.Foreground(lipgloss.Color("2"))
	editFileStyle       = baseLabelStyle.Foreground(lipgloss.Color("3"))
	deleteFileStyle     = baseLabelStyle.Foreground(lipgloss.Color("1"))
	renamedFileStyle    = baseLabelStyle.Foreground(lipgloss.Color("4"))
	copiedFileStyle     = baseLabelStyle.Foreground(lipgloss.Color("5"))
	unmodifiedFileStyle = baseLabelStyle.Foreground(lipgloss.Color("6"))
	unknownFileStyle    = baseLabelStyle.Foreground(lipgloss.Color("9"))
	legendsStyle        = lipgloss.NewStyle().Faint(true).PaddingRight(1)
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
		case model.StatusAdded:
			styled = newFileStyle.Render(
				model.StatusAdded.Name(),
			) + file.Path + " " + file.Status.String()
		case model.StatusModified:
			styled = editFileStyle.Render(
				model.StatusModified.Name(),
			) + file.Path + " " + file.Status.String()
		case model.StatusDeleted:
			styled = deleteFileStyle.Render(
				model.StatusDeleted.Name(),
			) + file.Path
		case model.StatusRenamed:
			styled = renamedFileStyle.Render(
				model.StatusRenamed.Name(),
			) + file.Path
		case model.StatusCopied:
			styled = copiedFileStyle.Render(
				model.StatusCopied.Name(),
			) + file.Path
		case model.StatusUnmodified:
			styled = unmodifiedFileStyle.Render(
				model.StatusUnmodified.Name(),
			) + file.Path
		case model.StatusUnknown:
			styled = unknownFileStyle.Render(
				model.StatusUnknown.Name(),
			) + file.Path
		default:
			styled = unknownFileStyle.Render("?") + file.Path
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
