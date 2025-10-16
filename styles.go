package main

import "github.com/charmbracelet/lipgloss"

// Styles for the TUI
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
