package src

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	windowsYTDLPUrl  = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	windowsYTDLPFile = "yt-dlp.exe"
)

var (
	homeDir, _ = os.UserHomeDir()
	confDir    = filepath.Join(homeDir, ".go-yt-dlp")
	binDir     = filepath.Join(confDir, "bin")
)

func InstallYTDLP() {
	if executableExists("yt-dlp") {
		return
	}

	if runtime.GOOS == "windows" {
		fmt.Println("Installing yt-dlp for Windows...")

		outputFile := filepath.Join(binDir, windowsYTDLPFile)
		err := downloadFile(windowsYTDLPUrl, outputFile)
		if err != nil {
			panic(err)
		}

		fmt.Println("yt-dlp installed successfully to", outputFile)
		return
	}

	if commandExists("brew") {
		fmt.Println("Installing yt-dlp with brew...")
		child := exec.Command("brew", "install", "yt-dlp")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr

		err := child.Run()
		if err != nil {
			panic(err)
		}

		return
	}

	panic("yt-dlp is not available and cannot be installed")
}

func InstallFFMPEG() {
	if commandExists("ffmpeg") {
		return
	}

	if commandExists("brew") {
		fmt.Println("Installing ffmpeg with brew...")
		child := exec.Command("brew", "install", "ffmpeg")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr

		err := child.Run()
		if err != nil {
			panic(err)
		}

		return
	}

	panic("ffmpeg is not available and cannot be installed")
}
