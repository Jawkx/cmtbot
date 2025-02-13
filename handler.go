package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) handleShowDiffState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "c":
		if len(m.stagedFiles) != 0 {
			m.state = GENERATING_COMMIT_STATE
			cfg, _ := LoadConfig()
			return m, generateCommitMessagesCmd(m.llmService, m.diff, m.stagedFiles, cfg.NumOfMsg)
		}
	}
	return m, nil
}

func (m model) handleSelectCommitState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.cursor++
		if m.cursor >= len(m.messages) {
			m.cursor = 0
		}
	case "k", "up":
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.messages) - 1
		}
	case "enter":
		m.state = COMMITING_RESULT_STATE
		return m, generateCommitChangesCmd(m.messages[m.cursor])
	case "e":
		m.state = EDIT_COMMIT_STATE
		m.textArea.SetValue(m.messages[m.cursor])
		return m, m.textArea.Focus()
	}
	return m, nil
}

func (m model) handleEditCommitState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+s":
		m.state = COMMITING_RESULT_STATE
		return m, generateCommitChangesCmd(m.textArea.Value())
	}
	return m, nil
}
