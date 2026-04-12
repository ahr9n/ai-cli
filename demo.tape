Output demo.gif

Set FontSize 16
Set Width 1000
Set Height 600
Set Theme "Catppuccin Mocha"
Set Padding 20
Set TypingSpeed 40ms
Set Shell bash

# ── Setup: set default provider so commands are clean ──
Hide
Type "ai-cli default set ollama --model qwen3:4b"
Enter
Sleep 1s
Show

# ── 1. Single prompt ──
Type "# Ask a quick question"
Enter
Sleep 300ms
Type "ai-cli ollama 'What is Go in one sentence?'"
Enter
Sleep 5s
Enter

# ── 2. Code preset ──
Type "# Use a code preset"
Enter
Sleep 300ms
Type "ai-cli ollama -p code 'FizzBuzz in Go, no comments'"
Enter
Sleep 7s
Enter

# ── 3. Concise preset ──
Type "# Concise mode for short answers"
Enter
Sleep 300ms
Type "ai-cli ollama -p concise 'List 3 benefits of goroutines'"
Enter
Sleep 5s
Enter

# ── 4. Piped input ──
Type "# Pipe input from other commands"
Enter
Sleep 300ms
Type `echo 'func main() { print(hello) }' | ai-cli ollama 'What is wrong with this code?'`
Enter
Sleep 5s
Enter

# ── 5. JSON output ──
Type "# JSON output for scripting"
Enter
Sleep 300ms
Type "ai-cli ollama -o json 'Capital of Egypt?' | head -8"
Enter
Sleep 5s
Enter

# ── 6. List models ──
Type "# List available models"
Enter
Sleep 300ms
Type "ai-cli ollama --list-models"
Enter
Sleep 3s
Enter

# ── 7. Default provider ──
Type "# Check default provider"
Enter
Sleep 300ms
Type "ai-cli default show"
Enter
Sleep 1s

Type "# Use default directly — no provider needed"
Enter
Sleep 300ms
Type `ai-cli 'Who created Linux?'`
Enter
Sleep 5s
Enter

# ── 8. Interactive mode with session ──
Type "# Interactive chat with auto-save"
Enter
Sleep 300ms
Type "ai-cli ollama -i --session demo"
Enter
Sleep 1s

Type "What are channels in Go?"
Enter
Sleep 5s

Type "Give a short example"
Enter
Sleep 6s

Type "exit"
Enter
Sleep 1s
Enter

# ── 9. Session management ──
Type "# List saved sessions"
Enter
Sleep 300ms
Type "ai-cli chat list"
Enter
Sleep 2s

Type "# Export as markdown"
Enter
Sleep 300ms
Type "ai-cli chat export demo | head -12"
Enter
Sleep 2s

Type "# Clean up"
Enter
Sleep 300ms
Type "ai-cli chat delete demo"
Enter
Sleep 1s
Enter

# ── 10. Available providers ──
Type "ai-cli providers"
Enter
Sleep 3s
Enter

# ── 11. Version ──
Type "ai-cli version"
Enter
Sleep 1s

# ── Cleanup: clear default ──
Hide
Type "ai-cli default clear"
Enter
Sleep 500ms
Show
