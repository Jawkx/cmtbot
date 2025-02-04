package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var optionBlockStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("4")).
	Foreground(lipgloss.Color("0")).Margin(0, 1)

func splitByLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func getStyledFilenames(fileString string) string {

	var titleStyle = lipgloss.NewStyle().MarginLeft(3)
	var baseLabelStyle = lipgloss.NewStyle().Width(6).MarginLeft(6)

	var newFileStyle = baseLabelStyle.Foreground(lipgloss.Color("2"))
	var editFileStyle = baseLabelStyle.Foreground(lipgloss.Color("3"))
	var deleteFileStyle = baseLabelStyle.Foreground(lipgloss.Color("1"))

	if fileString == "" {
		return splitByLines(
			titleStyle.Render("No staged files"),
			"",
			optionBlockStyle.Render("[q]uit"),
		)
	}

	fileStrings := strings.Split(fileString, "\n")

	var styledStagedFiles string

	for _, line := range fileStrings {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		status, fileName := parts[0], parts[1]

		var styled string
		switch status {
		case "A":
			styled = newFileStyle.Render("Add") + fileName
		case "M":
			styled = editFileStyle.Render("Edit") + fileName
		case "D":
			styled = deleteFileStyle.Render("Delete") + fileName
		default:
			styled = line
		}

		styledStagedFiles += styled + "\n"
	}

	return splitByLines(titleStyle.Render("Staged files:"), styledStagedFiles,
		optionBlockStyle.Render("[c]ontinue")+optionBlockStyle.Render("[q]uit"),
	)
}
