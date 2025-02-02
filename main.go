package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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

	messages, err := generateCommitMessages(diff)
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

func generateCommitMessages(diff string) ([]string, error) {
	prompt := fmt.Sprintf(`Generate 5 clear commit messages for these code changes.
Follow conventional commit format. Example format:
feat: add new search functionality
fix: resolve database connection issues

Code diff:
%s`, diff)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	resp, err := makeAPIRequest(reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	var messages []string
	for _, choice := range result.Choices {
		messages = append(messages, choice.Message.Content)
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
