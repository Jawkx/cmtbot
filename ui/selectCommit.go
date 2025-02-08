package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func SelectCommit(commits []string, cursor int) string {
	var (
		listWidth = 45
		boxWidth  = 55
		listStyle = lipgloss.NewStyle().Width(listWidth).Margin(0, 2)
		boxStyle  = lipgloss.NewStyle().
				Padding(0, 1).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("63"))
		cursorChar   = "> "
		legendsStyle = lipgloss.NewStyle().Faint(true).MarginLeft(4)
	)

	listItems := make([]string, len(commits)+3)
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
	boxHeight := calculateMaxHeight(commits, boxWidth)

	var boxContent string

	if len(commits) > 0 {
		boxContent = commits[cursor]
	} else {
		boxContent = "No commits."
	}

	box := boxStyle.Width(boxWidth).Height(boxHeight).Render(boxContent)

	header := legendsStyle.Render("(↑/↓: arrows, j/k: move, enter: select)")

	content := lipgloss.JoinHorizontal(lipgloss.Top, listStyle.Render(list), box)

	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
	)

	return ui
}

func calculateMaxHeight(commits []string, maxWidth int) int {
	var maxHeight int

	for _, commit := range commits {
		var commitHeight int
		lines := strings.Split(commit, "\n")

		for _, line := range lines {
			lineLen := len(line)

			if lineLen != 0 {
				commitHeight = commitHeight + int(math.Ceil(float64(maxWidth)/float64(lineLen)))
			}
			commitHeight = commitHeight + 1
		}

		if commitHeight > maxHeight {
			maxHeight = commitHeight
		}

	}

	return maxHeight
}
