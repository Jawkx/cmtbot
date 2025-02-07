package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"
)

type LlmService struct {
	apiBase   string
	apiKeyEnv string
	modelName string
}

func NewLlmService(apiBase, apiKeyEnv, modelName string) *LlmService {
	return &LlmService{
		apiBase,
		apiKeyEnv,
		modelName,
	}
}

func (s *LlmService) GenerateCommitMessages(
	diff string,
	numOfMessages int,
) ([]string, error) {

	prompt := fmt.Sprintf(
		GENERATE_COMMIT_MESSAGE_PROMPT,
		diff,
	)

	messages := make([]string, 0, numOfMessages)
	errChan := make(chan error, numOfMessages)
	resultChan := make(chan string, numOfMessages)

	for i := 0; i < numOfMessages; i++ {
		go func() {
			reqBody, err := json.Marshal(map[string]interface{}{
				"model":    s.modelName,
				"messages": []map[string]string{{"role": "user", "content": prompt}},
			})
			if err != nil {
				errChan <- fmt.Errorf("failed to marshal request body: %w", err)
				return
			}

			resp, err := s.makeAPIRequest(reqBody)
			if err != nil {
				errChan <- fmt.Errorf("API request failed: %w", err)
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
				errChan <- fmt.Errorf("failed to decode response body: %w", err)
				return
			}

			if len(result.Choices) > 0 {
				resultChan <- result.Choices[0].Message.Content
				errChan <- nil // Signal success
				return
			}
			errChan <- fmt.Errorf("no choices returned")
		}()
	}

	for i := 0; i < numOfMessages; i++ {
		err := <-errChan
		if err != nil {
			return nil, err
		}
		message := <-resultChan
		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("no commit messages generated")
	}

	return messages, nil
}

func (s *LlmService) makeAPIRequest(body []byte) (*http.Response, error) {
	apiKey := os.Getenv(s.apiKeyEnv)
	log.Printf("Making API request to: %s", s.apiBase)
	log.Printf("Using model: %s", s.modelName)
	req, err := http.NewRequest("POST", s.apiBase, bytes.NewBuffer(body))
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
