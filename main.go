package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jawkx/cmtbot/llm"
	"github.com/Jawkx/cmtbot/ui"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
)

var version string

type State int

const (
	SHOW_DIFF_STATE State = iota
	GENERATING_COMMIT_STATE
	SELECT_COMMIT_STATE
	EDIT_COMMIT_STATE
	COMMITING_RESULT_STATE
	COMMITED_CHANGES_STATE
)

type Config struct {
	ApiBase   string `toml:"api_base"`
	ApiKeyEnv string `toml:"api_key_env"`
	ModelName string `toml:"model_name"`
	NumOfMsg  int    `toml:"num_of_msg"`
	Prompt    string `toml:"prompt"`
}

type model struct {
	state     State
	diffFiles string
	diff      string
	err       error

	width  int
	height int

	messages          []string
	cursor            int
	succeedCommitHash string

	// services
	llmService *llm.LlmService

	// components
	textArea textarea.Model
	spinner  spinner.Model
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

	llmService := llm.NewLlmService(cfg.ApiBase, cfg.ApiKeyEnv, cfg.ModelName, cfg.Prompt)

	s := spinner.New(spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))))
	s.Spinner = spinner.Dot

	ti := textarea.New()
	ti.ShowLineNumbers = false

	return model{
		state:      SHOW_DIFF_STATE,
		diffFiles:  diffFiles,
		diff:       diff,
		spinner:    s,
		err:        err, // Store the config loading error
		llmService: llmService,
		textArea:   ti,
	}
}

type commitMsgsResultMsg struct {
	messages []string
	err      error
}

type commitChangesResultMsg struct {
	hash string
	err  error
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textArea.SetWidth(m.width)
		return m, nil

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case commitMsgsResultMsg:
		m.state = SELECT_COMMIT_STATE
		m.messages = msg.messages
		m.err = msg.err
		return m, nil

	case commitChangesResultMsg:
		m.state = COMMITED_CHANGES_STATE
		m.err = msg.err
		m.succeedCommitHash = msg.hash
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch m.state {
		case SHOW_DIFF_STATE:
			return m.handleShowDiffState(msg)
		case SELECT_COMMIT_STATE:
			return m.handleSelectCommitState(msg)
		case EDIT_COMMIT_STATE:
			return m.handleEditCommitState(msg)
		}
	}

	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func generateCommitMessagesCmd(llmService *llm.LlmService, diff string, numOfMsg int) tea.Cmd {
	return func() tea.Msg {
		messages, err := llmService.GenerateCommitMessages(diff, numOfMsg)
		return commitMsgsResultMsg{messages: messages, err: err}
	}
}

func commitChangesCmd(message string) tea.Cmd {
	return func() tea.Msg {
		hash, err := commitChanges(message)
		return commitChangesResultMsg{err: err, hash: hash}
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
			content = ui.SelectCommit(m.messages, m.cursor, m.width)
		}
	}

	if m.state == EDIT_COMMIT_STATE {
		content = ui.EditCommit(m.textArea)
	}

	if m.state == COMMITING_RESULT_STATE {
		content = "  " + m.spinner.View() + "Commiting changes..."
	}

	if m.state == COMMITED_CHANGES_STATE {
		if m.err != nil {
			content = lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Render(fmt.Sprintf("Error: %v", m.err))
		} else {
			content = ui.SuccessCommit(m.succeedCommitHash)
		}
	}

	return "\n" + content
}

func main() {

	versionFlag := flag.Bool("v", false, "Print version info")
	flag.Parse()

	if *versionFlag {
		fmt.Println("cmtbot version:", version)
		os.Exit(0)
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
