package llm

import (
	"bytes"
	"context"
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

func (s *LlmService) generateCommitMessages(
	ctx context.Context,
	diff string,
	numOfMessages int,
) ([]string, error) {

	prompt := fmt.Sprintf(
		GENERATE_COMMIT_MESSAGE_PROMPT,
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
					"model": s.modelName,
					"messages": []map[string]string{
						{"role": "user", "content": prompt},
					},
				})

				resp, err := s.makeAPIRequest(reqBody)
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

const GENERATE_COMMIT_MESSAGE_PROMPT = `
You are an expert at following the Conventional Commit specification. 

The expected git commit template is as follow :

    [type1/type2] Title
    - detail changes 1
    - detail changes 2

Where type and description is as follow

| Type     | Description                                                                                                                                |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| feat     | a new feature is introduced with the changes                                                                                               |
| fix      | a bug fix has occurred                                                                                                                     |
| chore    | changes that do not relate to a fix or feature and don't modify src or test files (for example updating dependencies)                      |
| refactor | refactored code that neither fixes a bug nor adds a feature                                                                                |
| docs     | updates to documentation such as the README or other markdown files                                                                        |
| style    | changes that do not affect the meaning of the code, likely related to code formatting such as white-space, missing semi-colons, and so on. |
| test     | including new or correcting previous tests                                                                                                 |
| perf     | performance improvements                                                                                                                   |
| ci       | continuous integration related                                                                                                             |
| build    | changes that affect the build system or external dependencies                                                                              |
| revert   | reverts a previous commit                                                                                                                  |

Given the git diff listed below, please generate a commit message for me following the rules below STRICTLY because the output will be consumed by another application

Rules: 

    1. Only reply with the raw generated commit message
    2. DON'T wrap the message in code tags
    3. DON'T give any explanation on the commit message
    4. Follow the template closely

Code diff:
%s
`
