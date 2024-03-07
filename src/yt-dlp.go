package src

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/buger/jsonparser"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	ytDlpPath             = "yt-dlp"
	downloadPath          = "ytdl-download"
	youtubeSearchUrl      = "https://youtube.com/search?q="
	youtubeMusicSearchUrl = "https://music.youtube.com/search?q="
	// https://github.com/yt-dlp/yt-dlp/issues/6007#issuecomment-1769137538
	youtubeMusicSearchPostfix = "#songs"

	PROGRESS_PREFIX = "[[DL]]"
	CUSTOM_PRESET   = "custom"
	DEFAULT_ARGS    = []string{"--force-keyframes-at-cuts", "--embed-metadata", "--no-playlist"}
	PROGRESS_ARGS   = []string{"--progress-template", PROGRESS_PREFIX + "%(progress)j", "--console-title"}
	PRESET_MAP      = [][]string{
		{"mp4-fast", "-f", "b"},
		{"mp4", "--remux-video", "mp4"},
		{"mp3", "-x", "--audio-format", "mp3", "-o", "%(uploader)s - %(title)s.%(ext)s"},
		{"wav", "-x", "--audio-format", "wav"},
		{CUSTOM_PRESET, "-f"},
	}
)

func setDownloadPath() tea.Msg {
	cwd, _ := os.Getwd()
	downloadPath = filepath.Join(cwd, downloadPath)

	return nil
}

func setExecutablePath() tea.Msg {
	path, err := findExecutable("yt-dlp")
	maybePanic(err)

	ytDlpPath = path
	return nil
}

type infoMsg []byte

func fetchInfo(m model) tea.Cmd {
	return func() tea.Msg {
		infoArgs := append(DEFAULT_ARGS[:], "-J")

		if !validateUrl(m.downloadQuery) {
			infoArgs = append(infoArgs, "-I", "1")

			if m.musicSearch {
				m.downloadQuery = youtubeMusicSearchUrl + m.downloadQuery + youtubeMusicSearchPostfix
			} else {
				m.downloadQuery = youtubeSearchUrl + m.downloadQuery
			}
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
		downloadArgs = append(downloadArgs, PROGRESS_ARGS...)
		downloadArgs = append(downloadArgs, "--load-info-json", "-")

		downloadCmd := exec.Command(ytDlpPath, downloadArgs...)
		downloadCmd.Stdin = strings.NewReader(string(m.infoOut))
		stdout, downloadErr := downloadCmd.StdoutPipe()
		maybePanic(downloadErr)

		downloadErr = downloadCmd.Start()
		maybePanic(downloadErr)

		scanner := bufio.NewScanner(stdout)
		scanner.Split(scanLinesCR)

		for scanner.Scan() {
			m.downloadLogChannel <- scanner.Text()
		}

		downloadErr = downloadCmd.Wait()
		maybePanic(downloadErr)

		close(m.downloadLogChannel)

		return downloadFinishMsg(true)
	}
}
