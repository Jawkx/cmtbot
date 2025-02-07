package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (s *LlmService) generateCommitMessages(
	diff string,
	numOfMessages int,
) ([]string, error) {

	prompt := fmt.Sprintf(
		GENERATE_COMMIT_MESSAGE_PROMPT,
		diff,
	)

	messages := make([]string, 0, numOfMessages)
	for i := 0; i < numOfMessages; i++ {
		reqBody, err := json.Marshal(map[string]interface{}{
			"model":    s.modelName,
			"messages": []map[string]string{{"role": "user", "content": prompt}},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		resp, err := s.makeAPIRequest(reqBody)
		if err != nil {
			return nil, fmt.Errorf("API request failed: %w", err)
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
			return nil, fmt.Errorf("failed to decode response body: %w", err)
		}

		if len(result.Choices) > 0 {
			messages = append(messages, result.Choices[0].Message.Content)
		}
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("no commit messages generated")
	}

	return messages, nil
}

func (s *LlmService) makeAPIRequest(body []byte) (*http.Response, error) {
	apiKey := os.Getenv(s.apiKeyEnv)

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
