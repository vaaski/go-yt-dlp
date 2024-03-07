//go:generate go-winres make
//go:generate goreleaser --clean --snapshot
//go:generate go run mac-bundle/main.go

package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/vaaski/go-yt-dlp/src"
)

// todo record ascii cinema for readme
func main() {
	// f, err := tea.LogToFile("debug.log", "debug")
	// if err != nil {
	// 	fmt.Println("fatal:", err)
	// 	os.Exit(1)
	// }
	// defer f.Close()

	updateFlag := flag.Bool("U", false, "Run auto-update.")
	wtFlag := flag.Bool("wt", false, "Do not try to open in Windows Terminal. Meant for internal use.")
	flag.Parse()

	if *updateFlag {
		src.AutoUpdate()
		os.Exit(0)
	}

	if !*wtFlag {
		src.OpenInWindowsTerminal()
	}

	src.InstallYTDLP()
	src.InstallFFMPEG()

	src.SetTermTitle("go-yt-dlp")

	_, teaErr := tea.NewProgram(src.InitialModel()).Run()
	if teaErr != nil {
		fmt.Printf("Could not start program :(\n%v\n", teaErr)
		os.Exit(1)
	}
}
