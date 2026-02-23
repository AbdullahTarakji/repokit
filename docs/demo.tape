Output docs/demo.gif

Set Shell "zsh"
Set FontSize 16
Set Width 1000
Set Height 650
Set Theme "Dracula"
Set TypingSpeed 50ms
Set Padding 20

# Scan repo with text output
Type "repokit scan --format text ."
Enter
Sleep 4s

# Scan with JSON output
Type "repokit scan --format json . | head -20"
Enter
Sleep 4s

# Auto-fix dry run
Type "repokit fix --dry-run ."
Enter
Sleep 4s

Sleep 1s
