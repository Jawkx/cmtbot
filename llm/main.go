package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jawkx/cmtbot/git"
)

type LlmService struct {
	apiBase   string
	apiKeyEnv string
	modelName string
	prompt    string
}

func NewLlmService(apiBase, apiKeyEnv, modelName, prompt string) *LlmService {
	return &LlmService{
		apiBase,
		apiKeyEnv,
		modelName,
		prompt,
	}
}

func (s *LlmService) GenerateCommitMessages(
	diff string,
	diffFiles []git.StagedFile,
	numOfMessages int,
) ([]string, error) {

	var promptStringBuilder strings.Builder
	promptStringBuilder.WriteString(s.prompt)

	promptStringBuilder.WriteString("## Changed Files:\n")
	for _, file := range diffFiles {

		fileContent, _ := getFileContent(file.Path)

		promptStringBuilder.WriteString(fmt.Sprintf("### FilePath: %s \n", file.Path))
		promptStringBuilder.WriteString("``` \n")
		promptStringBuilder.WriteString(fileContent)
		promptStringBuilder.WriteString("\n ```")
	}

	promptStringBuilder.WriteString("## Diff: \n ```")
	promptStringBuilder.WriteString(diff)
	promptStringBuilder.WriteString("\n ```")

	prompt := promptStringBuilder.String()

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

	for idx, message := range messages {
		messages[idx] = strings.TrimSpace(message)
	}

	return messages, nil
}
