package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

// This variable can be set at build time using -ldflags
var buildTimeAPIKey string

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
