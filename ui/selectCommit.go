package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func SelectCommit(commits []string, cursor int) string {
	var (
		listWidth = 35
		boxWidth  = 45
		listStyle = lipgloss.NewStyle().Width(listWidth).Margin(0, 2)
		boxStyle  = lipgloss.NewStyle().
				Width(boxWidth).
				Padding(0, 1).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63"))
		cursorChar = "> "
	)

	listItems := make([]string, len(commits))
	for i, commit := range commits {
		lines := strings.SplitN(commit, "\n", 2)
		firstLine := lines[0]

		if len(firstLine) > listWidth-len(cursorChar) {
			firstLine = firstLine[:listWidth-3-len(cursorChar)] + "..."
		}

		if i == cursor {
			listItems[i] = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render(cursorChar + firstLine)
		} else {
			listItems[i] = "  " + firstLine
		}
	}

	list := lipgloss.JoinVertical(lipgloss.Left, listItems...)
	var boxContent string
	if len(commits) > 0 {
		boxContent = commits[cursor]
	} else {
		boxContent = "No commits."
	}
	box := boxStyle.Render(boxContent)

	ui := lipgloss.JoinHorizontal(lipgloss.Top, listStyle.Render(list), box)

	return ui
}
