package llm

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

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
