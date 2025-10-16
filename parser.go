package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

// formatMarkdown converts markdown to styled terminal output
func formatMarkdown(markdown string, maxWidth int) string {
	var result strings.Builder
	lines := strings.Split(markdown, "\n")

	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#BD93F9")).
		MarginTop(1).
		MarginBottom(1)

	headingStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#50FA7B"))

	codeBlockStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50FA7B")).
		Background(lipgloss.Color("#282A36")).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))

	inCodeBlock := false
	var codeLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle code blocks
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				// End code block - render it
				for _, codeLine := range codeLines {
					result.WriteString("  ")
					result.WriteString(codeBlockStyle.Render(codeLine))
					result.WriteString("\n")
				}
				codeLines = []string{}
				inCodeBlock = false
				result.WriteString("\n")
			} else {
				// Start code block
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			codeLines = append(codeLines, line)
			continue
		}

		// Skip separator lines
		if strings.HasPrefix(trimmed, "---") || strings.HasPrefix(trimmed, "===") {
			result.WriteString(separatorStyle.Render(strings.Repeat("─", maxWidth)))
			result.WriteString("\n")
			continue
		}

		// Skip empty lines (but preserve some spacing)
		if trimmed == "" {
			result.WriteString("\n")
			continue
		}

		// Handle numbered section headers (1. Summary, 2. Quality Assessment, etc.)
		if len(trimmed) > 3 && trimmed[0] >= '1' && trimmed[0] <= '9' && trimmed[1] == '.' && trimmed[2] == ' ' {
			// Remove ** markers if present
			text := strings.TrimPrefix(trimmed[3:], "**")
			text = strings.TrimSuffix(text, "**")
			text = strings.TrimSuffix(text, ":")
			result.WriteString(sectionTitleStyle.Render("▸ " + text))
			result.WriteString("\n")
			continue
		}

		// Handle **bold** headers
		if strings.HasPrefix(trimmed, "**") && strings.HasSuffix(trimmed, "**") {
			text := strings.Trim(trimmed, "*")
			text = strings.TrimSuffix(text, ":")
			result.WriteString(headingStyle.Render("  • " + text))
			result.WriteString("\n")
			continue
		}

		// Handle ## headers
		if strings.HasPrefix(trimmed, "##") {
			text := strings.TrimPrefix(trimmed, "##")
			text = strings.TrimSpace(text)
			result.WriteString(sectionTitleStyle.Render("▸ " + text))
			result.WriteString("\n")
			continue
		}

		// Handle # headers
		if strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "##") {
			text := strings.TrimPrefix(trimmed, "#")
			text = strings.TrimSpace(text)
			result.WriteString(sectionTitleStyle.Render("▸ " + text))
			result.WriteString("\n")
			continue
		}

		// Handle bullet points
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			text := trimmed[2:]
			// Clean up inline markdown
			text = cleanInlineMarkdown(text)
			wrapped := wrapText(text, maxWidth-6)
			lines := strings.Split(wrapped, "\n")
			for i, l := range lines {
				if i == 0 {
					result.WriteString(normalStyle.Render("    • " + strings.TrimSpace(l)))
				} else {
					result.WriteString(normalStyle.Render("      " + strings.TrimSpace(l)))
				}
				result.WriteString("\n")
			}
			continue
		}

		// Handle file references with emoji
		if strings.Contains(trimmed, "📄") {
			fileRefStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8BE9FD")).
				Bold(true)
			result.WriteString(fileRefStyle.Render("  " + trimmed))
			result.WriteString("\n")
			continue
		}

		// Regular paragraphs
		text := cleanInlineMarkdown(trimmed)
		wrapped := wrapText(text, maxWidth-4)
		lines := strings.Split(wrapped, "\n")
		for _, l := range lines {
			result.WriteString(normalStyle.Render("  " + strings.TrimSpace(l)))
			result.WriteString("\n")
		}
	}

	return result.String()
}

// cleanInlineMarkdown removes or converts inline markdown syntax
func cleanInlineMarkdown(text string) string {
	// Remove bold markers but keep the text
	text = strings.ReplaceAll(text, "**", "")
	text = strings.ReplaceAll(text, "__", "")
	// Remove italic markers
	text = strings.ReplaceAll(text, "*", "")
	text = strings.ReplaceAll(text, "_", "")
	// Clean up inline code
	text = strings.ReplaceAll(text, "`", "")

	return text
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
