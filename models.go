package main

import (
	"github.com/charmbracelet/bubbles/spinner"
)

// ReviewItem represents a single issue or suggestion

type Severity string

const (
	SeverityHigh   Severity = "High"
	SeverityMedium Severity = "Medium"
	SeverityLow    Severity = "Low"
)

type ReviewItem struct {
	number     int
	title      string
	content    string
	codeBlocks []string
	severity   Severity
	checked    bool
}

// model is the state for the TUI
type model struct {
	spinner   spinner.Model
	loading   bool
	err       error
	review    string
	diff      string
	filePath  string
	apiKey    string
	quitting  bool
	items     []ReviewItem
	cursorPos int
	width     int
	height    int
}

// reviewMsg is sent when the review is complete
type reviewMsg struct {
	review string
	err    error
}

// OpenAI API structures
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}
