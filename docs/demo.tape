Output docs/demo.gif

Set Shell "zsh"
Set FontSize 14
Set Width 900
Set Height 600
Set Theme "Dracula"

Type "# RepoKit — GitHub Repository Health Auditor"
Enter
Sleep 500ms

Type "repokit scan --format text ."
Enter
Sleep 3s

Type "repokit scan --format json . | head -20"
Enter
Sleep 3s

Type "repokit fix --dry-run ."
Enter
Sleep 3s

Sleep 1s
