package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// This variable can be set at build time using -ldflags
var buildTimeAPIKey string

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginTop(1).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4")).
			MarginBottom(1)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#44475A"))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000")).
			MarginTop(1).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			MarginTop(1)

	contentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			MarginLeft(2).
			MarginBottom(0)

	severityHighStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF0000")).
				Background(lipgloss.Color("#4a0000"))

	severityMediumStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFA500")).
				Background(lipgloss.Color("#4a3000"))

	severityLowStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFD700")).
				Background(lipgloss.Color("#3a3000"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)
)

// ReviewItem represents a single issue or suggestion
type ReviewItem struct {
	number     int
	title      string
	content    string
	codeBlocks []string
	severity   string
	checked    bool
}

// Model for the TUI
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

type reviewMsg struct {
	review string
	err    error
}

func initialModel(apiKey, filePath, diff string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	return model{
		spinner:   s,
		loading:   true,
		filePath:  filePath,
		diff:      diff,
		apiKey:    apiKey,
		items:     []ReviewItem{},
		cursorPos: 0,
		width:     120,
		height:    40,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.getReview(),
	)
}

func (m model) getReview() tea.Cmd {
	return func() tea.Msg {
		review, err := reviewDiff(m.apiKey, m.diff)
		return reviewMsg{review: review, err: err}
	}
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		wordLen := len(word)
		if lineLen+wordLen+1 > width && lineLen > 0 {
			result.WriteString("\n  ")
			lineLen = 2
		} else if i > 0 {
			result.WriteString(" ")
			lineLen++
		}
		result.WriteString(word)
		lineLen += wordLen
	}

	return result.String()
}

// parseReviewIntoItems extracts issues and suggestions from the review
func parseReviewIntoItems(review string) []ReviewItem {
	items := []ReviewItem{}
	lines := strings.Split(review, "\n")
	itemNum := 1

	var currentItem *ReviewItem
	inIssuesSection := false
	inSuggestionsSection := false
	inCodeBlock := false
	var currentCodeBlock []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect sections
		if strings.Contains(trimmed, "**Issues Found**") || strings.Contains(trimmed, "3. Issues") {
			inIssuesSection = true
			inSuggestionsSection = false
			continue
		} else if strings.Contains(trimmed, "**Suggestions**") || strings.Contains(trimmed, "4. Suggestions") {
			inIssuesSection = false
			inSuggestionsSection = true
			continue
		} else if strings.HasPrefix(trimmed, "**") && strings.HasSuffix(trimmed, "**") {
			inIssuesSection = false
			inSuggestionsSection = false
			continue
		}

		// Skip non-issue/suggestion sections
		if !inIssuesSection && !inSuggestionsSection {
			continue
		}

		// Detect file references as item titles
		if strings.Contains(trimmed, "📄") || (strings.Contains(trimmed, ":") && len(trimmed) > 0) {
			isFileRef := false
			fileExtensions := []string{".go", ".js", ".ts", ".py", ".java", ".vue", ".jsx", ".tsx"}
			for _, ext := range fileExtensions {
				if strings.Contains(trimmed, ext) {
					isFileRef = true
					break
				}
			}

			if isFileRef {
				// Save previous item
				if currentItem != nil {
					items = append(items, *currentItem)
				}
				// Start new item
				currentItem = &ReviewItem{
					number:     itemNum,
					title:      trimmed,
					content:    "",
					codeBlocks: []string{},
					severity:   "Low",
				}
				itemNum++
				continue
			}
		}

		// Handle code blocks
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				// End of code block - save accumulated code
				if currentItem != nil && len(currentCodeBlock) > 0 {
					currentItem.codeBlocks = append(currentItem.codeBlocks, strings.Join(currentCodeBlock, "\n"))
					currentCodeBlock = []string{}
				}
				inCodeBlock = false
			} else {
				// Start of code block
				inCodeBlock = true
				currentCodeBlock = []string{}
			}
			continue
		}

		// Accumulate code block content
		if inCodeBlock {
			currentCodeBlock = append(currentCodeBlock, line)
			continue
		}

		// Detect severity
		if currentItem != nil && (strings.Contains(trimmed, "Severity:") || strings.Contains(trimmed, "severity:")) {
			lower := strings.ToLower(trimmed)
			if strings.Contains(lower, "critical") || strings.Contains(lower, "high") {
				currentItem.severity = "High"
			} else if strings.Contains(lower, "medium") {
				currentItem.severity = "Medium"
			} else {
				currentItem.severity = "Low"
			}
			continue
		}

		// Add content to current item
		if currentItem != nil && trimmed != "" {
			if currentItem.content != "" {
				currentItem.content += " "
			}
			currentItem.content += trimmed
		}
	}

	// Add last item
	if currentItem != nil {
		items = append(items, *currentItem)
	}

	return items
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		if !m.loading {
			switch msg.String() {
			case "up", "k":
				if m.cursorPos > 0 {
					m.cursorPos--
				}
			case "down", "j":
				if m.cursorPos < len(m.items)-1 {
					m.cursorPos++
				}
			case " ", "x":
				if len(m.items) > 0 && m.cursorPos < len(m.items) {
					m.items[m.cursorPos].checked = !m.items[m.cursorPos].checked
				}
			case "a":
				// Mark all as checked
				for i := range m.items {
					m.items[i].checked = true
				}
			case "n":
				// Mark all as unchecked
				for i := range m.items {
					m.items[i].checked = false
				}
			case "enter":
				m.quitting = true
				return m, tea.Quit
			}
		}

	case reviewMsg:
		m.loading = false
		m.review = msg.review
		m.err = msg.err
		if msg.err == nil {
			m.items = parseReviewIntoItems(msg.review)
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder
	maxWidth := m.width - 10
	if maxWidth > 110 {
		maxWidth = 110
	}

	// Header
	s.WriteString(titleStyle.Render("🔍 Revyu - AI-Powered Code Review"))
	s.WriteString("\n")

	target := m.filePath
	if target == "." {
		target = "all changed files"
	}
	s.WriteString(subtitleStyle.Render(fmt.Sprintf("Reviewing: %s", target)))
	s.WriteString("\n\n")

	if m.loading {
		s.WriteString(m.spinner.View())
		s.WriteString(" Analyzing git diff with AI...\n")
		s.WriteString(subtitleStyle.Render("  This may take a few moments"))
		return s.String()
	}

	if m.err != nil {
		s.WriteString(errorStyle.Render("❌ Error"))
		s.WriteString("\n")
		s.WriteString(contentStyle.Render(m.err.Error()))
		s.WriteString("\n\n")
		s.WriteString(subtitleStyle.Render("Press 'q' to quit"))
		return s.String()
	}

	// Success header
	s.WriteString(successStyle.Render("✅ Review Complete"))
	s.WriteString("\n")
	s.WriteString(separatorStyle.Render(strings.Repeat("─", maxWidth)))
	s.WriteString("\n")

	// Show summary
	checkedCount := 0
	for _, item := range m.items {
		if item.checked {
			checkedCount++
		}
	}
	s.WriteString(subtitleStyle.Render(fmt.Sprintf("Found %d issues/suggestions  •  %d completed", len(m.items), checkedCount)))
	s.WriteString("\n\n")

	// Display items as interactive checklist
	if len(m.items) == 0 {
		// Fallback: show raw review if no items parsed
		s.WriteString(contentStyle.Render(wrapText(m.review, maxWidth-4)))
		s.WriteString("\n\n")
	} else {
		// Display each item with checkbox
		for i, item := range m.items {
			// Cursor indicator
			cursor := "  "
			if i == m.cursorPos {
				cursor = "▶ "
			}

			// Checkbox
			checkbox := "[ ]"
			if item.checked {
				checkbox = "[✓]"
			}

			// Severity badge
			var severityBadge string
			switch item.severity {
			case "High":
				severityBadge = severityHighStyle.Render(" ⚠ HIGH ")
			case "Medium":
				severityBadge = severityMediumStyle.Render(" ● MED ")
			case "Low":
				severityBadge = severityLowStyle.Render(" ○ LOW ")
			}

			// Item header with number, checkbox, and severity
			headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8F8F2"))
			if i == m.cursorPos {
				headerStyle = headerStyle.Background(lipgloss.Color("#44475A"))
			}

			itemHeader := fmt.Sprintf("%s%s #%d ", cursor, checkbox, item.number)
			s.WriteString(headerStyle.Render(itemHeader))
			s.WriteString(severityBadge)
			s.WriteString("\n")

			// File reference
			fileStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8BE9FD")).
				MarginLeft(4)
			s.WriteString(fileStyle.Render(item.title))
			s.WriteString("\n")

			// Content (wrapped)
			if item.content != "" {
				wrappedContent := wrapText(item.content, maxWidth-6)
				contentLines := strings.Split(wrappedContent, "\n")
				for _, line := range contentLines {
					s.WriteString(contentStyle.Render("    " + line))
					s.WriteString("\n")
				}
			}

			// Code blocks
			if len(item.codeBlocks) > 0 {
				s.WriteString("\n")
				codeStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#50FA7B")).
					Background(lipgloss.Color("#282A36")).
					Padding(0, 1).
					MarginLeft(4)

				for _, codeBlock := range item.codeBlocks {
					codeLines := strings.Split(codeBlock, "\n")
					for _, codeLine := range codeLines {
						s.WriteString(codeStyle.Render(codeLine))
						s.WriteString("\n")
					}
				}
			}
			s.WriteString("\n")
		}
	}

	// Footer with instructions
	s.WriteString(separatorStyle.Render(strings.Repeat("─", maxWidth)))
	s.WriteString("\n")

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")).
		Italic(true)

	s.WriteString(instructionStyle.Render("↑/↓: Navigate  •  Space/X: Toggle  •  A: Check all  •  N: Uncheck all  •  Enter/Q: Quit"))
	s.WriteString("\n")

	return boxStyle.Render(s.String())
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

func main() {
	// Get OpenAI API key with priority order:
	// 1. Build-time embedded key
	// 2. Environment variable
	// 3. .env file
	apiKey := buildTimeAPIKey

	if apiKey == "" {
		// Try environment variable first
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if apiKey == "" {
		// Try loading from .env file as fallback
		err := godotenv.Load()
		if err == nil {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
	}

	if apiKey == "" {
		fmt.Println(errorStyle.Render("❌ Error: OPENAI_API_KEY not found"))
		fmt.Println(contentStyle.Render("Please either:"))
		fmt.Println(contentStyle.Render("  1. Set OPENAI_API_KEY environment variable"))
		fmt.Println(contentStyle.Render("  2. Create a .env file with your OpenAI API key"))
		fmt.Println(contentStyle.Render("  3. Build with embedded key using: ./build.sh"))
		os.Exit(1)
	}

	// Get file path from command line arguments
	if len(os.Args) < 2 {
		fmt.Println(titleStyle.Render("🔍 Revyu - AI-Powered Code Review TESTING"))
		fmt.Println()
		fmt.Println(subtitleStyle.Render("Usage:"))
		fmt.Println(contentStyle.Render("  revyu <filename>  - Review git diff for a specific file"))
		fmt.Println(contentStyle.Render("  revyu .           - Review git diff for all tracked files"))
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Get git diff
	diff, err := getGitDiff(filePath)
	if err != nil {
		fmt.Println(errorStyle.Render("❌ Error getting git diff"))
		fmt.Println(contentStyle.Render(err.Error()))
		os.Exit(1)
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println(subtitleStyle.Render("No changes detected in git diff"))
		os.Exit(0)
	}

	// Run the TUI
	p := tea.NewProgram(initialModel(apiKey, filePath, diff))
	if _, err := p.Run(); err != nil {
		fmt.Println(errorStyle.Render("Error running program: " + err.Error()))
		os.Exit(1)
	}
}

// getGitDiff executes git diff command and returns the output
func getGitDiff(filePath string) (string, error) {
	var cmd *exec.Cmd

	if filePath == "." {
		// Get diff for all files
		cmd = exec.Command("git", "diff")
	} else {
		// Get diff for specific file
		cmd = exec.Command("git", "diff", filePath)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %v\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// reviewDiff sends the diff to OpenAI for review
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
