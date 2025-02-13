package main

import (
	"github.com/Jawkx/cmtbot/git"
	"github.com/Jawkx/cmtbot/llm"
	tea "github.com/charmbracelet/bubbletea"
)

type commitMsgsResultMsg struct {
	messages []string
	err      error
}

func generateCommitMessagesCmd(
	llmService *llm.LlmService,
	diff string,
	diffFiles []git.StagedFile,
	numOfMsg int,
) tea.Cmd {
	return func() tea.Msg {
		messages, err := llmService.GenerateCommitMessages(diff, diffFiles, numOfMsg)
		return commitMsgsResultMsg{messages: messages, err: err}
	}
}

type commitChangesResultMsg struct {
	hash string
	err  error
}

func generateCommitChangesCmd(message string) tea.Cmd {
	return func() tea.Msg {
		hash, err := git.CommitChanges(message)
		return commitChangesResultMsg{err: err, hash: hash}
	}
}
