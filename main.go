package main

import (
	"fmt"

	"os"

	"github.com/Jawkx/cmtbot/llm"
	"github.com/Jawkx/cmtbot/ui"
	"github.com/charmbracelet/bubbles/spinner"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	apiBase   = "https://openrouter.ai/api/v1/chat/completions"
	apiKeyEnv = "OPENROUTER_API_KEY"
	modelName = "google/gemini-flash-1.5"
	numOfMsg  = 5
)

const (
	SHOW_DIFF_STATE         = "show_diff"
	GENERATING_COMMIT_STATE = "generating_commit_state"
	SELECT_COMMIT_STATE     = "select_commit_state"
)

type model struct {
	state     string
	diffFiles string
	diff      string
	err       error

	messages []string
	// services
	llmService *llm.LlmService
	spinner    spinner.Model
}

func initialModel() model {
	diffFiles, _ := getStagedFiles()
	diff, _ := getStagedDiff()
	llmService := llm.NewLlmService(apiBase, apiKeyEnv, modelName)
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		state:     SHOW_DIFF_STATE,
		diffFiles: diffFiles,
		diff:      diff,
		spinner:   s,

		llmService: llmService,
	}
}

type commitMsgsResultMsg struct {
	messages []string
	err      error
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case commitMsgsResultMsg:
		m.state = SELECT_COMMIT_STATE
		m.messages = msg.messages
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "c":
			if m.state == SHOW_DIFF_STATE {
				m.state = GENERATING_COMMIT_STATE
				return m, generateCommitMessagesCmd(m.llmService, m.diff, numOfMsg)
			}
		}

		return m, nil
	}

	return m, nil
}

func generateCommitMessagesCmd(llmService *llm.LlmService, diff string, numOfMsg int) tea.Cmd {
	return func() tea.Msg {
		messages, err := llmService.GenerateCommitMessages(diff, numOfMsg)
		return commitMsgsResultMsg{messages: messages, err: err}
	}
}

func (m model) View() string {
	var content string

	if m.state == SHOW_DIFF_STATE {
		content = ui.StagedFiles(m.diffFiles)
	}

	if m.state == GENERATING_COMMIT_STATE {
		content = "  " + m.spinner.View() + "Generating commit message..."
	}

	if m.state == SELECT_COMMIT_STATE {
		if m.err != nil {
			content = fmt.Sprintf("Error: %v", m.err)
		} else {
			content = "SELECT COMMIT MESSAGE"
		}
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
