package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func SelectCommit(commits []string, cursor int) string {
	var (
		listWidth = 50
		boxWidth  = 50
		boxHeight = calculateMaxHeight(commits, boxWidth-2)
		listStyle = lipgloss.NewStyle().
				Width(listWidth).
				Height(boxHeight).
				PaddingTop(1)
		boxStyle = lipgloss.NewStyle().
				Padding(0, 1).
				MarginLeft(1).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("63"))
		cursorChar   = "> "
		legendsStyle = lipgloss.NewStyle().Faint(true).PaddingRight(1)
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

	legends := legendsStyle.Render("(↑/↓: arrows, j/k: move, enter: select)")

	list := listStyle.Render(lipgloss.JoinVertical(lipgloss.Left, listItems...))
	listWithLegends := lipgloss.JoinVertical(lipgloss.Right, list, legends)

	var boxContent string

	if len(commits) > 0 {
		boxContent = commits[cursor]
	} else {
		boxContent = "No commits."
	}

	box := boxStyle.Width(boxWidth).Height(boxHeight).Render(boxContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, listWithLegends, box)

	ui := lipgloss.JoinVertical(
		lipgloss.Left,
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
				commitHeight = commitHeight + int(math.Ceil(float64(lineLen)/float64(maxWidth)))
			} else {
				commitHeight = commitHeight + 1
			}

		}

		if commitHeight > maxHeight {
			maxHeight = commitHeight
		}

	}

	return maxHeight
}
