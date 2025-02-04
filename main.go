package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	apiBase   = "https://openrouter.ai/api/v1/chat/completions"
	apiKeyEnv = "OPENROUTER_API_KEY"
	modelName = "deepseek/deepseek-r1-distill-qwen-32b"
)

const (
	SHOW_DIFF_STATE      = "show_diff"
	SHOW_FULL_DIFF_STATE = "show_full_diff"
)

type model struct {
	state     string
	diffFiles string
}

func initialModel() model {
	diffFiles, _ := getStagedFiles()
	return model{
		state:     SHOW_DIFF_STATE,
		diffFiles: diffFiles,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var content string

	if m.state == SHOW_DIFF_STATE {
		content = getStyledFilenames(m.diffFiles)
	}

	return "\n" + content
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// func main() {
// 	diff, err := getStagedDiff()
// 	if err != nil {
// 		fmt.Printf("Error getting diff: %v\n", err)
// 		os.Exit(1)
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()
//
// 	messages, err := generateCommitMessages(ctx, diff, 5)
//
// 	if err != nil {
// 		fmt.Printf("Error generating messages: %v\n", err)
// 		os.Exit(1)
// 	}
//
// 	selected := selectMessage(messages)
// 	if selected == "" {
// 		fmt.Println("Commit cancelled")
// 		return
// 	}
//
// 	if err := commitChanges(selected); err != nil {
// 		fmt.Printf("Error committing: %v\n", err)
// 		os.Exit(1)
// 	}
// }
