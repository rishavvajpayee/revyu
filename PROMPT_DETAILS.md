# 📝 Enhanced Prompt Details

This document explains the enhanced prompt used in Revyu and what you can expect in the AI code review output.

## Prompt Structure

The AI is instructed to provide a comprehensive code review with:

### 1. Specific Requirements

For **every** point made, the AI must:
- ✅ Reference the exact file and line numbers (e.g., `main.go:45-50`)
- ✅ Include relevant code snippets in markdown code blocks
- ✅ Explain what should change and why

### 2. Review Sections

#### Summary
A brief overview of what changed in the commit/diff.

#### Quality Assessment
Analysis of:
- Code quality observations
- Best practices compliance  
- Performance considerations
- Each point references specific files and line numbers

#### Issues Found
For each issue discovered:
- **File reference**: `📄 main.go:42`
- **Description**: Clear explanation of the problem
- **Code snippet**: Shows the problematic code
- **Severity level**: Critical/High/Medium/Low

Example:
```
📄 auth.go:156

Missing error handling in authentication flow

func authenticate(token string) User {
    user := validateToken(token)
    return user  // No error check!
}

Severity: High
```

#### Suggestions
For each improvement recommendation:
- **File reference**: `📄 utils.go:78`
- **What to change**: Clear directive
- **Code snippet**: Shows recommended implementation
- **Explanation**: Why this is better

Example:
```
📄 database.go:234

Add connection pooling for better performance

// Recommended:
pool := &sql.DB{
    MaxOpenConns: 25,
    MaxIdleConns: 5,
    ConnMaxLifetime: 5 * time.Minute,
}

This improves performance by reusing connections and prevents connection exhaustion.
```

## Visual Formatting

The terminal UI renders these elements with special styling:

### File References
```
📄 main.go:42-45
```
- **Color**: Cyan text on dark background
- **Style**: Bold, padded
- **Purpose**: Instantly identify which file needs attention

### Code Blocks
```go
func example() error {
    return nil
}
```
- **Color**: Green text on dark background
- **Style**: Monospace, indented
- **Purpose**: Easy to read code examples

### Severity Badges
- 🔴 **HIGH**: Red background - Critical issues requiring immediate attention
- 🟠 **MEDIUM**: Orange background - Important issues to address soon
- 🟡 **LOW**: Yellow background - Nice-to-have improvements

### Section Headers
```
▸ Issues Found
▸ Suggestions  
```
- **Color**: Orange
- **Style**: Bold with arrow prefix
- **Purpose**: Clear navigation through review sections

## Benefits

### For Developers
1. **Jump directly to issues** - File:line references let you navigate instantly
2. **See the actual code** - No need to switch contexts to understand the issue
3. **Prioritize fixes** - Severity levels help you tackle critical issues first
4. **Learn best practices** - Explanations help you understand the "why"

### For Teams
1. **Consistent reviews** - Same structure every time
2. **Actionable feedback** - Specific, not vague
3. **Educational** - Team members learn from each review
4. **Time-saving** - No back-and-forth asking "which file?" or "where exactly?"

## Example Full Review

```
🔍 Revyu - AI-Powered Code Review
Reviewing: all changed files

✅ Review Complete

▸ 1. Summary

  This PR adds user authentication with JWT tokens and password hashing.
  Changes span 3 files with 145 additions and 23 deletions.

▸ 2. Quality Assessment

  - Good separation of concerns with dedicated auth package
  - Password hashing implemented correctly with bcrypt
  - JWT implementation follows best practices

  📄 auth/jwt.go:12-25
  
  Token expiration set appropriately at 24 hours.

▸ 3. Issues Found

  📄 auth/handler.go:42

  Missing rate limiting on login endpoint
  
  func LoginHandler(w http.ResponseWriter, r *http.Request) {
      // No rate limiting implemented
      user := authenticate(r)
  }
  
   HIGH  Security vulnerability

  📄 models/user.go:67

  Password field exposed in JSON serialization
  
  type User struct {
      ID       int    `json:"id"`
      Email    string `json:"email"`
      Password string `json:"password"`  // Should be omitted
  }
  
   MEDIUM  Data exposure risk

▸ 4. Suggestions

  📄 auth/handler.go:42

  Add rate limiting middleware
  
  // Recommended:
  import "golang.org/x/time/rate"
  
  var loginLimiter = rate.NewLimiter(5, 10)
  
  func LoginHandler(w http.ResponseWriter, r *http.Request) {
      if !loginLimiter.Allow() {
          http.Error(w, "Too many requests", 429)
          return
      }
      user := authenticate(r)
  }
  
  Prevents brute force attacks on the login endpoint.

  📄 models/user.go:67

  Use JSON tag to omit password field
  
  // Recommended:
  type User struct {
      ID       int    `json:"id"`
      Email    string `json:"email"`
      Password string `json:"-"`  // Omit from JSON
  }
  
  Ensures password hash never appears in API responses.
```

## Tips for Best Results

1. **Make focused commits** - Smaller diffs get more detailed reviews
2. **Use descriptive file names** - Helps AI provide better context
3. **Include tests** - AI will review your test coverage too
4. **Review regularly** - Catch issues early before they compound

## Supported File Types

File references are automatically detected for:
- Go (`.go`)
- JavaScript/TypeScript (`.js`, `.ts`, `.jsx`, `.tsx`)
- Python (`.py`)
- Java (`.java`)
- Ruby (`.rb`)
- PHP (`.php`)
- C/C++ (`.c`, `.cpp`, `.h`, `.hpp`)
- Rust (`.rs`)
- Swift (`.swift`)
- Kotlin (`.kt`)
- Scala (`.scala`)
- Vue (`.vue`)
- CSS/SCSS (`.css`, `.scss`)
- HTML/XML (`.html`, `.xml`)
- Config files (`.json`, `.yaml`, `.yml`)
- SQL (`.sql`)
- Shell (`.sh`, `.bash`)

