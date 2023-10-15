package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// todo music search
// todo record ascii cinema for readme
func main() {
	_, teaErr := tea.NewProgram(initialModel()).Run()
	if teaErr != nil {
		fmt.Printf("Could not start program :(\n%v\n", teaErr)
		os.Exit(1)
	}
}
