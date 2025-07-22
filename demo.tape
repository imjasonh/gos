# VHS documentation: https://github.com/charmbracelet/vhs

Output demo.gif

Require gos

Set FontSize 20
Set Width 1200
Set Height 600
Set Theme "Dracula"

Type "# Run it directly with shebang"
Enter
Type "./examples/hello.go foo bar"
Enter
Sleep 3s

Type "# Dependencies are automatically downloaded and cached by Go"
Enter
Sleep 2s
Type "./examples/hello.go foo bar"
Enter
Sleep 3s

Type "# Run tests in the same file"
Enter
Type "gos test examples/hello.go"
Enter
Sleep 3s
