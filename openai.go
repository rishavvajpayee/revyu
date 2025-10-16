package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func reviewDiff(apiKey, diff string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	prompt := fmt.Sprintf(`You are an expert code reviewer. Please review the following git diff and provide a detailed analysis.

For each point you make, please:
- Reference the specific file and approximate line numbers (e.g., "main.go:45-50")
- Include relevant code snippets using markdown code blocks with language syntax
- Be specific about what should be changed and why

Please structure your review with these sections:

1. **Summary**: Brief overview of what changed

2. **Quality Assessment**:
   - Code quality observations
   - Best practices compliance
   - Performance considerations
   Reference specific files and line numbers.

3. **Issues Found**:
   For each issue, provide:
   - File reference (e.g., "📄 main.go:42")
   - Description of the problem
   - Code snippet showing the issue
   - Severity (Critical/High/Medium/Low)

4. **Suggestions**:
   For each suggestion, provide:
   - File reference (e.g., "📄 utils.go:78")
   - What to change
   - Code snippet showing the recommended change
   - Explanation of why this is better

Use markdown code blocks with proper language syntax highlighting.
Use file references in the format: 📄 filename.ext:lineNumber

Here's the git diff:

%s

Please provide a comprehensive review with specific file references and code examples.`, diff)

	requestBody := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("OpenAI API call failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}
