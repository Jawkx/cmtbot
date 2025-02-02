package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	apiBase   = "https://openrouter.ai/api/v1/chat/completions"
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

	var messages []string
	spinnerDone := make(chan bool)
	go func() {
		s := spinner.New().Title("Generating commit messages...").Context(ctx)
		s.Run()
		spinnerDone <- true
	}()

	messages, err = generateCommitMessages(ctx, diff)
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

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCommitMessages(ctx context.Context, diff string) ([]string, error) {
	prompt := fmt.Sprintf(
		`You are an expert at following the Conventional Commit specification. Given the git diff listed below, please generate a commit message for me, only reply with the raw generated commit message and don't include any explaination of context

Code diff:
%s`,
		diff,
	)

	resultChan := make(chan string, 5)
	errChan := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func() {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				resultChan <- ""
				return
			default:
				reqBody, _ := json.Marshal(map[string]interface{}{
					"model": modelName,
					"messages": []map[string]string{
						{"role": "user", "content": prompt},
					},
				})

				resp, err := makeAPIRequest(reqBody)
				if err != nil {
					errChan <- err
					resultChan <- ""
					return
				}
				defer resp.Body.Close()

				var result struct {
					Choices []struct {
						Message struct {
							Content string `json:"content"`
						} `json:"message"`
					} `json:"choices"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					errChan <- err
					resultChan <- ""
					return
				}

				if len(result.Choices) > 0 {
					resultChan <- result.Choices[0].Message.Content
				} else {
					resultChan <- ""
				}
				errChan <- nil
			}
		}()
	}

	var messages []string
	var firstError error

	for i := 0; i < 5; i++ {
		msg := <-resultChan
		err := <-errChan
		if err != nil && firstError == nil {
			firstError = err
		}
		if msg != "" {
			messages = append(messages, msg)
		}
	}

	if len(messages) == 0 && firstError != nil {
		return nil, firstError
	}

	return messages, nil
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

func commitChanges(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func makeAPIRequest(body []byte) (*http.Response, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")

	req, err := http.NewRequest("POST", apiBase, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return resp, nil
}
