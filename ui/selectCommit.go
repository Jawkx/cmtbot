package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func SelectCommit(commits []string, cursor int, viewportWidth int) string {

	var (
		listWidth = viewportWidth/2 - 2
		boxWidth  = viewportWidth/2 - 2
	)

	var isVertical bool

	if listWidth < 35 {
		isVertical = true
		listWidth = viewportWidth - 4
		boxWidth = viewportWidth - 4
	}

	var (
		boxHeight = calculateMaxHeight(commits, boxWidth-2)
		listStyle = lipgloss.NewStyle().
				Width(listWidth).
				PaddingTop(1)
		boxStyle = lipgloss.NewStyle().
				Padding(0, 1).
				MarginLeft(1).
				Width(boxWidth).Height(boxHeight).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("63"))
		cursorChar     = "> "
		truncateString = "..."
		legendsStyle   = lipgloss.NewStyle().Faint(true).PaddingLeft(len(cursorChar))
	)

	legends := legendsStyle.Render("(↑/↓;j/k: move, enter: select, q: quit)")
	listItems := make([]string, len(commits))
	for i, commit := range commits {
		lines := strings.SplitN(commit, "\n", 2)
		firstLine := lines[0]

		if len(firstLine) > listWidth-len(cursorChar) {
			firstLine = firstLine[:listWidth-len(truncateString)-len(cursorChar)] + truncateString
		}

		if i == cursor {
			listItems[i] = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render(cursorChar + firstLine)
		} else {
			listItems[i] = "  " + firstLine
		}
	}

	list := listStyle.Render(lipgloss.JoinVertical(lipgloss.Left, listItems...))
	listWithLegends := lipgloss.JoinVertical(lipgloss.Left, list, legends)

	var boxContent string

	if len(commits) > 0 {
		boxContent = commits[cursor]
	} else {
		boxContent = "No commits."
	}

	box := boxStyle.Render(boxContent)

	var ui string

	if isVertical {
		ui = lipgloss.JoinVertical(lipgloss.Center, box, listWithLegends)
	} else {
		ui = lipgloss.JoinHorizontal(lipgloss.Top, listWithLegends, box)
	}

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
