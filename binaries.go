package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
)

type platformDownload struct {
	url      string
	filename string
}

var (
	ytDlpPath string = "yt-dlp"

	YT_DLP_DOWNLOADS = map[string]platformDownload{
		"darwin": {
			filename: "yt-dlp",
			url:      "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos",
		},
		"linux": {
			filename: "yt-dlp",
			url:      "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux",
		},
		"windows": {
			filename: "yt-dlp.exe",
			url:      "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe",
		},
	}
)

func DownloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func checkExistingYTDLP() (string, error) {
	if _, err := exec.LookPath("yt-dlp"); err == nil {
		return "yt-dlp", nil
	} else if fileExists("./yt-dlp") {
		return "./yt-dlp", nil
	} else if fileExists("./yt-dlp.exe") {
		return "./yt-dlp.exe", nil
	}

	return "", fmt.Errorf("yt-dlp not found")
}

func ensureYTDLP() (string, error) {
	existing, err := checkExistingYTDLP()
	if err == nil {
		fmt.Println("existing", existing)
		return existing, nil
	}

	platform := YT_DLP_DOWNLOADS[runtime.GOOS]
	if platform == (platformDownload{}) {
		fmt.Println("Unsupported platform")
		return "", fmt.Errorf("unsupported platform")
	}

	currentFolder := getCurrentFolder()
	destination := path.Join(currentFolder, platform.filename)

	fmt.Println("downloading", platform.filename, "from", platform.url, "to", destination)

	err = DownloadFile(platform.url, destination)
	if err != nil {
		fmt.Println("error downloading yt-dlp:", err)
		return "", err
	}

	return destination, nil
}
