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
	maybePanic(err)
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	maybePanic(err)

	fmt.Println("Successfully updated go-yt-dlp to latest version")
	updateYTDLP()
	updateFFMPEG()
}

func updateYTDLP() {
	if !executableExists("yt-dlp") {
		InstallYTDLP()
		return
	}

	if runtime.GOOS == "windows" {
		path, err := findExecutable("yt-dlp")
		maybePanic(err)

		fmt.Println("Updating yt-dlp for Windows...")
		child := exec.Command(path, "-U")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr

		err = child.Run()
		maybePanic(err)
		return
	}

	if commandExists("brew") {
		fmt.Println("Updating yt-dlp with brew...")
		child := exec.Command("brew", "upgrade", "yt-dlp")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr

		err := child.Run()
		maybePanic(err)

		return
	}

	fmt.Println("cannot update yt-dlp on this platform")
}

func updateFFMPEG() {
	if !executableExists("ffmpeg") {
		InstallFFMPEG()
		return
	}

	if commandExists("brew") {
		fmt.Println("Updating ffmpeg with brew...")
		child := exec.Command("brew", "upgrade", "ffmpeg")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr

		err := child.Run()
		maybePanic(err)

		return
	}

	fmt.Println("cannot update ffmpeg on this platform")
}
