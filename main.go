package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// todo record ascii cinema for readme
// todo download binaries automatically
// todo add a progress bar for downloads
// todo add a spinner for fetching info
func main() {
	_, teaErr := tea.NewProgram(initialModel()).Run()
	if teaErr != nil {
		fmt.Printf("Could not start program :(\n%v\n", teaErr)
		os.Exit(1)
	}
}
