//go:generate go-winres make
//go:generate goreleaser --clean --snapshot
//go:generate go run mac-bundle/main.go

package main

import (
	"fmt"
	"os"

	"github.com/vaaski/go-yt-dlp/src"

	tea "github.com/charmbracelet/bubbletea"
)

// todo record ascii cinema for readme
// todo download binaries automatically
// todo add a progress bar for downloads
// todo add a spinner for fetching info
func main() {
	src.SetTermTitle("go-yt-dlp")

	_, teaErr := tea.NewProgram(src.InitialModel()).Run()
	if teaErr != nil {
		fmt.Printf("Could not start program :(\n%v\n", teaErr)
		os.Exit(1)
	}
}
