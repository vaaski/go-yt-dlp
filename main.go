//go:generate go-winres make
//go:generate goreleaser --clean --snapshot
//go:generate go run mac-bundle/main.go

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vaaski/go-yt-dlp/src"

	tea "github.com/charmbracelet/bubbletea"
)

// todo record ascii cinema for readme
// todo download binaries automatically
func main() {
	// f, err := tea.LogToFile("debug.log", "debug")
	// if err != nil {
	// 	fmt.Println("fatal:", err)
	// 	os.Exit(1)
	// }
	// defer f.Close()

	updateFlag := flag.Bool("U", false, "Pass -U to auto-update")
	flag.Parse()

	if *updateFlag {
		src.AutoUpdate()
		os.Exit(0)
	}

	src.SetTermTitle("go-yt-dlp")

	_, teaErr := tea.NewProgram(src.InitialModel()).Run()
	if teaErr != nil {
		fmt.Printf("Could not start program :(\n%v\n", teaErr)
		os.Exit(1)
	}
}
