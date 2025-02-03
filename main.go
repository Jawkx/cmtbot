package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"os"
	"time"
)

const (
	apiBase   = "https://openrouter.ai/api/v1/chat/completions"
	apiKeyEnv = "OPENROUTER_API_KEY"
	modelName = "deepseek/deepseek-r1-distill-qwen-32b"
)

func main() {
	diff, err := getStagedDiff()
	if err != nil {
		fmt.Printf("Error getting diff: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	spinnerDone := make(chan bool)
	go func() {
		s := spinner.New().Title("Generating commit messages...").Context(ctx)
		s.Run()
		spinnerDone <- true
	}()

	messages, err := generateCommitMessages(ctx, diff, 5)
	cancel()
	<-spinnerDone

	if err != nil {
		fmt.Printf("Error generating messages: %v\n", err)
		os.Exit(1)
	}

	selected := selectMessage(messages)
	if selected == "" {
		fmt.Println("Commit cancelled")
		return
	}

	if err := commitChanges(selected); err != nil {
		fmt.Printf("Error committing: %v\n", err)
		os.Exit(1)
	}
}

func selectMessage(messages []string) string {
	fmt.Println("\nSelect a commit message:")
	for i, msg := range messages {
		fmt.Printf("%d. %s\n", i+1, msg)
	}

	fmt.Print("\nEnter number (0 to cancel): ")
	var selection int
	fmt.Scanln(&selection)

	fmt.Printf("SELECTION: %d", selection)
	if selection < 1 || selection > len(messages) {
		return ""
	}
	return messages[selection-1]
}
