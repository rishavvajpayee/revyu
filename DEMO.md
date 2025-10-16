# 🎨 Revyu Terminal UI Demo

## Enhanced Code Review with File References & Code Snippets

When you run `revyu`, you'll experience a modern, beautiful terminal interface with detailed code analysis:

### 1. Loading State
```
🔍 Revyu - AI-Powered Code Review
Reviewing: main.go

⠋ Analyzing git diff with AI...
  This may take a few moments
```

### 2. Review Display

The review is beautifully formatted with:

#### Header Section
- **Purple title** with tool name
- **Gray subtitle** showing the file being reviewed
- **Green success indicator** when complete

#### Color-Coded Sections
- **Orange section headers** (Summary, Quality Assessment, Issues, Suggestions)
- **White content text** with proper indentation
- **Cyan file references** on dark background (📄 main.go:45)
- **Green code blocks** on dark background with syntax highlighting
- **Severity badges** with color coding:
  - 🔴 RED background for HIGH severity
  - 🟠 ORANGE background for MEDIUM severity
  - 🟡 YELLOW background for LOW severity

### Example Output Sections

#### Issues Section
```
▸ Issues Found

  📄 main.go:42

  Potential null pointer dereference

  func processData(data *Data) {
      fmt.Println(data.Name)  // No nil check
  }

   HIGH  Critical issue
```

#### Suggestions Section
```
▸ Suggestions

  📄 utils.go:78

  Add input validation before processing

  // Recommended change:
  func validateInput(input string) error {
      if len(input) == 0 {
          return errors.New("empty input")
      }
      return nil
  }

  This prevents downstream errors and improves reliability.
```

#### Interactive Features
- **Keyboard controls**: Press `q` or `Ctrl+C` to quit anytime
- **Smooth animations**: Dot spinner during loading
- **Rounded border box**: Professional layout
- **Auto-formatted output**: Parses markdown and structures content

### 3. Error Handling

If something goes wrong, you'll see:
- **Red error messages** that are clear and actionable
- **Helpful hints** for resolving issues
- **Clean exit** with instructions

## Color Scheme

- 🟣 **Primary (Purple)**: #7D56F4 - Headers and branding
- 🟠 **Secondary (Orange)**: #FFA500 - Section headers
- 🟢 **Success (Green)**: #04B575 - Success messages and code
- 🔴 **Error (Red)**: #FF0000 - Error messages
- ⚪ **Text (White)**: #FFFFFF - Main content
- ⚫ **Muted (Gray)**: #626262 - Subtitles and hints

## User Experience

1. **Fast startup** - Immediate feedback with loading spinner
2. **Clear progress** - You always know what's happening
3. **Easy to read** - Color coding makes scanning easy
4. **Professional look** - Rounded borders and proper spacing
5. **Keyboard friendly** - Simple controls (Enter, q, Ctrl+C)

## Technical Details

Built with:
- **Bubble Tea** - Terminal UI framework
- **Lipgloss** - Styling and layout
- **Bubbles/Spinner** - Loading animations

The UI is:
- **Responsive** - Adapts to your terminal size
- **Async** - Non-blocking API calls
- **Interactive** - Real-time updates
- **Accessible** - Clear keyboard controls

