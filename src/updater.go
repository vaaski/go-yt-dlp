package src

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/minio/selfupdate"
)

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func AutoUpdate() {
	baseUrl := "https://github.com/vaaski/go-yt-dlp/releases/latest/download/go-yt-dlp-"
	downloadUrl := baseUrl + runtime.GOOS + "-" + runtime.GOARCH

	if runtime.GOOS == "windows" {
		downloadUrl += ".exe"
	}

	resp, err := http.Get(downloadUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully updated go-yt-dlp to latest version")
	updateYTDLP()
	updateFFMPEG()
}

func updateYTDLP() {
	if !commandExists("yt-dlp") {
		InstallYTDLP()
		return
	}

	if commandExists("brew") {
		fmt.Println("Updating yt-dlp with brew...")
		child := exec.Command("brew", "upgrade", "yt-dlp")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr

		child.Run()
	}
}

func updateFFMPEG() {
	if !commandExists("ffmpeg") {
		InstallFFMPEG()
		return
	}

	if commandExists("brew") {
		fmt.Println("Updating ffmpeg with brew...")
		child := exec.Command("brew", "upgrade", "ffmpeg")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr

		child.Run()
	}
}
