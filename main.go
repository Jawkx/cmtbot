package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jawkx/cmtbot/llm"
	"github.com/Jawkx/cmtbot/ui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pelletier/go-toml/v2"
)

const (
	SHOW_DIFF_STATE         = "show_diff"
	GENERATING_COMMIT_STATE = "generating_commit_state"
	SELECT_COMMIT_STATE     = "select_commit_state"
)

type Config struct {
	ApiBase   string `toml:"api_base"`
	ApiKeyEnv string `toml:"api_key_env"`
	ModelName string `toml:"model_name"`
	NumOfMsg  int    `toml:"num_of_msg"`
}

type model struct {
	state     string
	diffFiles string
	diff      string
	err       error

	messages []string
	cursor   int
	// services
	llmService *llm.LlmService
	spinner    spinner.Model
}

func LoadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("error getting home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "cmtbot", "cmtbot.toml")

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return Config{
			ApiBase:   "https://openrouter.ai/api/v1/chat/completions",
			ApiKeyEnv: "OPENROUTER_API_KEY",
			ModelName: "google/gemini-flash-1.5",
			NumOfMsg:  5,
		}, nil
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = toml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return config, nil
}

func initialModel() model {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		// You might want to exit here or provide a fallback.
	}

	diffFiles, _ := getStagedFiles()
	diff, _ := getStagedDiff()
	llmService := llm.NewLlmService(cfg.ApiBase, cfg.ApiKeyEnv, cfg.ModelName)
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		state:      SHOW_DIFF_STATE,
		diffFiles:  diffFiles,
		diff:       diff,
		spinner:    s,
		err:        err, // Store the config loading error
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
				cfg, _ := LoadConfig()
				return m, generateCommitMessagesCmd(m.llmService, m.diff, cfg.NumOfMsg)
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
			content = ui.SelectCommit(m.messages, m.cursor)
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
