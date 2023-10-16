package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/buger/jsonparser"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	ytDlpPath             string = "yt-dlp"
	downloadPath          string = "ytdl-download"
	youtubeSearchUrl             = "https://youtube.com/search?q="
	youtubeMusicSearchUrl        = "https://music.youtube.com/search?q="

	DEFAULT_ARGS = [...]string{"--force-keyframes-at-cuts", "--embed-metadata", "--no-playlist", "--console-title"}
	PRESET_MAP   = [][]string{
		{"mp4-fast", "-f", "b"},
		{"mp4", "--remux-video", "mp4"},
		{"mp3", "-x", "--audio-format", "mp3", "-o", "%(uploader)s - %(title)s.%(ext)s"},
	}
)

func setDownloadPath() tea.Msg {
	executablePath, _ := os.Executable()
	executableFolder := path.Join(executablePath, "..")
	if strings.HasPrefix(executableFolder, "/var/folders") {
		// the path for the executable is in some temp folder when using `go run .`
		// so we use the current working directory instead
		cwd, _ := os.Getwd()
		downloadPath = path.Join(cwd, downloadPath)
	} else {
		downloadPath = path.Join(executableFolder, downloadPath)
	}

	return nil
}

func setExecutablePath() tea.Msg {
	if fileExists("./yt-dlp") {
		ytDlpPath = "./yt-dlp"
	} else if fileExists("./yt-dlp.exe") {
		ytDlpPath = "./yt-dlp.exe"
	}

	return nil
}

type infoMsg []byte

func fetchInfo(m model) tea.Cmd {
	return func() tea.Msg {
		infoArgs := append(DEFAULT_ARGS[:], "-J")

		if !validateUrl(m.downloadQuery) {
			if m.musicSearch {
				m.downloadQuery = youtubeMusicSearchUrl + m.downloadQuery
			} else {
				m.downloadQuery = youtubeSearchUrl + m.downloadQuery
			}
			infoArgs = append(infoArgs, "-I", "1")
		}

		infoArgs = append(infoArgs, m.downloadQuery)

		infoCmd := exec.Command(ytDlpPath, infoArgs...)
		infoOut, infoErr := infoCmd.Output()
		maybePanic(infoErr)

		//? if it's a playlist (like when you search for a song), just use the first entry
		firstEntry, _, _, entryErr := jsonparser.Get(infoOut, "entries", "[0]")
		if firstEntry != nil && entryErr == nil {
			infoOut = firstEntry
		}

		return infoMsg(infoOut)
	}
}

type downloadFinishMsg bool

func startDownload(m model) tea.Cmd {
	return func() tea.Msg {
		downloadArgs := append(DEFAULT_ARGS[:], PRESET_MAP[m.selectedPreset][1:]...)
		downloadArgs = append(downloadArgs, "-P", downloadPath)

		cpuCount := runtime.NumCPU()
		downloadArgs = append(downloadArgs, "-N", fmt.Sprint(cpuCount))
		downloadArgs = append(downloadArgs, "--load-info-json", "-")

		downloadCmd := exec.Command(ytDlpPath, downloadArgs...)
		downloadCmd.Stdin = strings.NewReader(string(m.infoOut))
		stdout, downloadErr := downloadCmd.StdoutPipe()
		maybePanic(downloadErr)

		downloadErr = downloadCmd.Start()
		maybePanic(downloadErr)

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m.downloadLogChannel <- scanner.Text()
		}
		close(m.downloadLogChannel)

		downloadErr = downloadCmd.Wait()
		maybePanic(downloadErr)

		return downloadFinishMsg(true)
	}
}
