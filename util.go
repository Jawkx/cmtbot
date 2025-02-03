package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func generateCommitMessages(ctx context.Context, diff string, numOfMessages int) ([]string, error) {
	prompt := fmt.Sprintf(
		`You are an expert at following the Conventional Commit specification. Given the git diff listed below, please generate a commit message for me following the rules below strictly

        Rules: 
        1. Only reply with the raw generated commit message
        2. Don't wrap the message in code tags

Code diff:
%s`,
		diff,
	)

	resultChan := make(chan string, numOfMessages)
	errChan := make(chan error, numOfMessages)

	for i := 0; i < numOfMessages; i++ {
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

	for i := 0; i < numOfMessages; i++ {
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

func makeAPIRequest(body []byte) (*http.Response, error) {
	apiKey := os.Getenv(apiKeyEnv)

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
