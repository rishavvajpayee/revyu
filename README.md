# Revyu 🔍

An AI-powered Git diff reviewer that uses OpenAI to provide intelligent code review feedback with a beautiful terminal UI.

## Features

- 🎨 **Beautiful Terminal UI** - Modern, colorful interface with smooth animations
- 🤖 **AI-Powered Reviews** - Uses OpenAI's GPT-4 for intelligent code analysis
- ⚡ **Fast & Responsive** - Live loading spinner and progress indicators
- 📝 **Flexible Input** - Review specific files or all changed files
- 🔍 **Comprehensive Analysis** including:
  - Summary of changes
  - Quality assessment
  - Issue detection (bugs, security concerns)
  - Improvement suggestions
- 🎯 **Smart Formatting** - Automatically formatted and color-coded output

## Prerequisites

- Go 1.16 or higher
- Git repository
- OpenAI API key

## Installation

### Option 1: Build with Embedded API Key (Recommended)

This creates a self-contained binary with the API key baked in:

1. Clone or navigate to the project directory:
```bash
cd /Users/rishavvajpayee/Projects/revyu
```

2. Create a `.env` file with your OpenAI API key:
```bash
cp .env.example .env
```

3. Edit `.env` and add your OpenAI API key:
```
OPENAI_API_KEY=sk-your-actual-api-key-here
```

4. Build with embedded key and code signing:

**Using Makefile (Recommended):**
```bash
make build          # Build and sign the binary
make install        # Build, sign, and install to /usr/local/bin
make help          # Show all available commands
```

**Using build script:**
```bash
./build.sh
```

5. Install globally (if not using `make install`):
```bash
sudo cp revyu /usr/local/bin/revyu
```

### Option 2: Build without Embedded Key

If you prefer to use environment variables:

1. Build normally:
```bash
go build -o revyu
```

2. Set environment variable (add to `~/.zshrc` or `~/.bashrc`):
```bash
export OPENAI_API_KEY="sk-your-actual-api-key-here"
```

3. Install globally:
```bash
sudo cp revyu /usr/local/bin/revyu
```

## Usage

### Review a specific file
```bash
./revyu path/to/your/file.go
```

### Review all changed files
```bash
./revyu .
```

## Makefile Commands

The Makefile provides convenient shortcuts for building, installing, and managing the binary:

```bash
make build          # Build and code sign the binary
make install        # Build and install to /usr/local/bin
make verify         # Verify the code signature of the local binary
make verify-installed # Verify the code signature of the installed binary
make clean          # Remove built binaries
make uninstall      # Remove the installed binary from /usr/local/bin
make rebuild        # Clean and rebuild
make reinstall      # Clean, rebuild, and reinstall
make test           # Run tests
make help           # Show all available commands
```

**Note:** The Makefile automatically handles:
- Loading environment variables from `.env`
- Building with `CGO_ENABLED=1`
- Code signing with runtime hardening and timestamp
- Installing with proper permissions

## How It Works

1. **Git Diff**: The tool runs `git diff` to get the changes in your repository
2. **AI Analysis**: Sends the diff to OpenAI's GPT-4 for comprehensive review
3. **Review Output**: Displays a structured review with:
   - Summary of changes
   - Quality assessment
   - Identified issues
   - Suggestions for improvement

## Example

```bash
# Make some changes to your code
echo "func example() {}" >> main.go

# Review the changes
revyu main.go
```

You'll see a beautiful terminal UI with:
- 🔄 **Loading spinner** while AI analyzes your code
- 🎨 **Color-coded sections** for easy reading
- 📊 **Structured output** with clear sections:
  - Summary
  - Quality Assessment
  - Issues
  - Suggestions
- ⌨️ **Interactive controls** - Press Enter or 'q' to quit

## UI Features

The terminal UI includes:
- **Purple themed** interface with professional styling
- **Animated spinner** during API calls
- **Section headers** in orange for easy navigation
- **File references** highlighted in cyan with dark background (e.g., 📄 main.go:42)
- **Code blocks** with syntax-aware styling and dark background
- **Severity badges** color-coded (🔴 HIGH, 🟠 MEDIUM, 🟡 LOW)
- **Error messages** in red with helpful hints
- **Success indicators** in green
- **Boxed layout** with rounded borders
- **Keyboard shortcuts** for navigation

### Enhanced Code Review Output

The AI now provides:
- **Specific file and line references** for each issue
- **Code snippets** showing problematic code
- **Before/after examples** in suggestions
- **Severity ratings** for prioritizing fixes
- **Detailed explanations** with context

## Development

### Dependencies

The project uses these amazing Go libraries:
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Style definitions and layout
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - TUI components (spinner)
- **[godotenv](https://github.com/joho/godotenv)** - Environment variable loading

To run the application without building:
```bash
go run main.go <filename|.>
```

To install dependencies:
```bash
go mod download
```

To tidy dependencies:
```bash
go mod tidy
```

### UI Demo & Prompt Details

- **[DEMO.md](DEMO.md)** - Visual walkthrough of the terminal interface
- **[PROMPT_DETAILS.md](PROMPT_DETAILS.md)** - Complete guide to what the AI provides in reviews

These documents show you exactly what to expect from Revyu's code reviews, including file references, code snippets, and severity ratings.

## Environment Variables

- `OPENAI_API_KEY` (required): Your OpenAI API key

## Notes

- The tool only reviews **uncommitted changes** (working directory vs. staged/committed)
- Make sure you're in a Git repository before running the tool
- The quality of review depends on OpenAI's API response
- API calls to OpenAI may incur costs based on your usage

## License

MIT

## Contributing

Feel free to open issues or submit pull requests for improvements!

