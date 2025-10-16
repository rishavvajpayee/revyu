package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

// Update handles messages and updates the model state
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
				for i := range m.items {
					m.items[i].checked = true
				}
			case "n":
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

// View renders the TUI
func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder
	maxWidth := m.width - 10
	if maxWidth > 110 {
		maxWidth = 110
	}

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
		// Fallback: show formatted review if no items parsed
		s.WriteString(formatMarkdown(m.review, maxWidth))
		s.WriteString("\n")
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
