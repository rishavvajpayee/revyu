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
- 🌍 **Cross-Platform** - Works on macOS, Linux, and Windows

## Prerequisites

- Go 1.19 or higher
- Git repository
- OpenAI API key

## Installation

> 📖 **For detailed platform-specific instructions, see [INSTALL.md](INSTALL.md)**

### Quick Install (Cross-Platform)

> ⚠️ **Important**: Create a `.env` file with your OpenAI API key to embed it in the binary. Without it, you'll need to set the `OPENAI_API_KEY` environment variable every time you use Revyu.

**macOS / Linux:**
```bash
# Clone the repository
git clone https://github.com/yourusername/revyu.git
cd revyu

# Create .env file with your API key (REQUIRED for embedded key)
cp .env.example .env
# Then edit .env and replace with your actual API key
# Or create it directly:
# echo "OPENAI_API_KEY=sk-your-api-key-here" > .env

# Build and install
make install
```

**Windows (PowerShell):**
```powershell
# Clone the repository
git clone https://github.com/yourusername/revyu.git
cd revyu

# Create .env file with your API key (REQUIRED for embedded key)
"OPENAI_API_KEY=sk-your-api-key-here" | Out-File -FilePath .env -Encoding ASCII

# Build and install
.\build.ps1
make install
```

### Platform-Specific Notes

- **macOS**: Automatic code signing applied; installs to `/usr/local/bin`
- **Linux**: Installs to `/usr/local/bin` (may require sudo)
- **Windows**: Installs to `%USERPROFILE%\bin` (ensure it's in PATH)

The Makefile automatically detects your OS and uses appropriate paths and commands.

### Alternative: Build Without Embedded Key

If you prefer NOT to embed the API key in the binary, you can build without a `.env` file and set the environment variable each time:

```bash
# Build without .env file
go build -o revyu

# Then ALWAYS set this environment variable before using revyu:
export OPENAI_API_KEY="sk-your-api-key-here"  # macOS/Linux
# or
$env:OPENAI_API_KEY="sk-your-api-key-here"    # Windows PowerShell

# Make it permanent by adding to your shell profile:
# macOS/Linux: Add to ~/.bashrc, ~/.zshrc, or ~/.profile
# Windows: Use setx or add to PowerShell profile
```

**Note:** This method requires setting the environment variable every time you open a new terminal, unless you add it to your shell profile.

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
make help           # Show all available commands
make build          # Build the binary (OS-aware)
make install        # Build and install (OS-aware paths)
make verify         # Verify code signature (macOS only)
make clean          # Remove built binaries
make uninstall      # Remove the installed binary
make rebuild        # Clean and rebuild
make reinstall      # Clean, rebuild, and reinstall
make test           # Run tests
```

**The Makefile automatically:**
- Detects your operating system (macOS, Linux, or Windows)
- Uses appropriate install paths for each platform
- Loads environment variables from `.env`
- Applies code signing on macOS
- Handles permissions correctly per platform

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

