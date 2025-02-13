package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jawkx/cmtbot/git"
	"github.com/Jawkx/cmtbot/llm"
	"github.com/Jawkx/cmtbot/ui"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type model struct {
	state       State
	stagedFiles []git.StagedFile
	diff        string
	err         error

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

func initialModel() model {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		// You might want to exit here or provide a fallback.
	}

	stagedFiles, _ := git.GetStagedFiles()
	diff, _ := git.GetStagedDiff()

	llmService := llm.NewLlmService(cfg.ApiBase, cfg.ApiKeyEnv, cfg.ModelName, cfg.Prompt)

	s := spinner.New(spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))))
	s.Spinner = spinner.Dot

	ti := textarea.New()
	ti.ShowLineNumbers = false

	return model{
		state:       SHOW_DIFF_STATE,
		stagedFiles: stagedFiles,
		diff:        diff,
		spinner:     s,
		err:         err, // Store the config loading error
		llmService:  llmService,
		textArea:    ti,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textArea.SetWidth(m.width)

		break

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		break

	case commitMsgsResultMsg:
		m.state = SELECT_COMMIT_STATE
		m.messages = msg.messages
		m.err = msg.err
		break

	case commitChangesResultMsg:
		m.state = COMMITED_CHANGES_STATE
		m.err = msg.err
		m.succeedCommitHash = msg.hash
		break

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

		var mNew tea.Model
		switch m.state {
		case SHOW_DIFF_STATE:
			mNew, cmd = m.handleShowDiffState(msg)
		case SELECT_COMMIT_STATE:
			mNew, cmd = m.handleSelectCommitState(msg)
		case EDIT_COMMIT_STATE:
			mNew, cmd = m.handleEditCommitState(msg)
		}
		m = mNew.(model)
		cmds = append(cmds, cmd)
	}

	m.textArea, cmd = m.textArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var content string

	if m.state == SHOW_DIFF_STATE {
		content = ui.StagedFiles(m.stagedFiles)
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
