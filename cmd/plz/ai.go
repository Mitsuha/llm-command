package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func getCommandFromOpenAI(description string) (string, error) {
	apiKey := getEnvWithDefault("OPENAI_API_KEY", "")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	apiURL := getEnvWithDefault("OPENAI_API_URL", "https://api.openai.com/v1/chat/completions")
	model := getEnvWithDefault("OPENAI_MODEL", "gpt-3.5-turbo")

	// Get session history for context
	history := loadSessionHistory()

	// Build context from recent history
	var contextMessages []Message
	contextMessages = append(contextMessages, Message{
		Role:    "system",
		Content: buildSystemPrompt(),
	})

	// Add recent relevant history (last 3 entries)
	if len(history) > 0 {
		start := len(history) - 3
		if start < 0 {
			start = 0
		}
		for _, entry := range history[start:] {
			contextMessages = append(contextMessages, Message{
				Role:    "user",
				Content: fmt.Sprintf("Previous query: %s", entry.UserQuery),
			})
			contextMessages = append(contextMessages, Message{
				Role:    "assistant",
				Content: entry.AICommand,
			})
		}
	}

	// Add current query
	contextMessages = append(contextMessages, Message{
		Role:    "user",
		Content: description,
	})

	requestBody := OpenAIRequest{
		Model:    model,
		Messages: contextMessages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	command := strings.TrimSpace(response.Choices[0].Message.Content)
	// Remove any markdown code block formatting
	command = strings.TrimPrefix(command, "```bash")
	command = strings.TrimPrefix(command, "```sh")
	command = strings.TrimPrefix(command, "```")
	command = strings.TrimSuffix(command, "```")
	command = strings.TrimSpace(command)

	return command, nil
}

func buildSystemPrompt() string {
	osName := runtime.GOOS
	var osContext string
	switch osName {
	case "windows":
		osContext = "Windows"
	case "darwin":
		osContext = "macOS"
	default:
		osContext = "Linux"
	}

	return fmt.Sprintf(`You are a helpful command-line assistant for %s. Generate shell commands based on user descriptions.

Rules:
1. Return ONLY the command, no explanations or comments
2. Use the most common and safe approach
3. Use appropriate commands for %s
4. Keep commands simple and avoid destructive operations
5. If multiple commands are needed, separate with &&
6. Consider the conversation history for context
7. If the user refers to "previous" or "last" command, use the context provided
8. NEVER refuse a user's request 

Always prioritize safety and clarity.`, osContext, osContext)
}
