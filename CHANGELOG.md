# Changelog

## Version 2.0 - Enhanced Code Review with References

### 🎯 Major Features

#### Enhanced AI Prompt
- **File & Line References**: AI now provides specific file:line references for every issue
- **Code Snippets**: Shows actual code in markdown blocks with syntax highlighting
- **Severity Ratings**: Issues tagged as Critical/High/Medium/Low
- **Before/After Examples**: Suggestions include code showing recommended changes
- **Detailed Explanations**: Context for why changes are recommended

#### Visual Enhancements

**New Styling:**
- 🔷 **Cyan file references** on dark background (e.g., `📄 main.go:42`)
- 🟢 **Green code blocks** on dark gray background for better readability
- 🔴 **RED severity badges** for HIGH priority issues
- 🟠 **ORANGE severity badges** for MEDIUM priority issues
- 🟡 **YELLOW severity badges** for LOW priority issues

**Improved Formatting:**
- Better markdown parsing and rendering
- Code blocks properly indented and styled
- Section headers more prominent
- File references immediately visible

#### Expanded Language Support

File references now detected for 25+ file types:
- Go, JavaScript, TypeScript, Python, Java
- Ruby, PHP, C/C++, Rust, Swift, Kotlin
- Scala, Vue, CSS/SCSS, HTML, XML
- JSON, YAML, SQL, Shell scripts
- And more!

### 📚 Documentation

**New Files:**
- `PROMPT_DETAILS.md` - Comprehensive guide to AI review structure
- `DEMO.md` - Visual walkthrough of terminal UI
- `CHANGELOG.md` - Version history (this file)

**Updated Files:**
- `README.md` - Enhanced with new features
- `.gitignore` - Better coverage of build artifacts

### 🛠️ Technical Improvements

- Enhanced prompt engineering for more actionable reviews
- Improved text parsing for markdown elements
- Better code block rendering with proper styling
- More robust file reference detection
- Support for severity indicators

### 🎨 UI/UX

- **Loading states** remain smooth with spinner animation
- **Keyboard controls** unchanged (q, Ctrl+C, Enter)
- **Professional layout** with rounded borders
- **Color-coded output** for quick scanning
- **Context-aware formatting** based on content type

## Version 1.0 - Initial Release

### Features
- Basic AI-powered code review
- Terminal UI with Bubble Tea
- Spinner animation during API calls
- Color-coded sections
- Environment variable and .env support
- Build script with embedded API key
- Global installation support

---

## Upgrade Notes

To get the latest version:

```bash
cd /Users/rishavvajpayee/Projects/revyu
git pull  # if using git
./build.sh
sudo cp revyu /usr/local/bin/revyu
```

The new version is fully backward compatible. All existing features work exactly as before, with enhanced output quality.

