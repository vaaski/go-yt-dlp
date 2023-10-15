package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/buger/jsonparser"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

var ytDlpPath string = "yt-dlp"
var downloadPath string = "ytdl-download"
var DEFAULT_ARGS = [...]string{"--force-keyframes-at-cuts", "--embed-metadata", "--no-playlist", "--console-title"}
var PRESET_MAP = map[string][]string{
	"mp4":      {"--remux-video", "mp4"},
	"mp4-fast": {"-f", "b"},
	"mp3":      {"-x", "--audio-format", "mp3", "-o", "%(uploader)s - %(title)s.%(ext)s"},
}

// todo music search
// todo format flags
// todo record ascii cinema for readme
func main() {
	f, logErr := tea.LogToFile("debug.log", "debug")
	if logErr != nil {
		fmt.Println("fatal:", logErr)
		os.Exit(1)
	}
	defer f.Close()

	_, teaErr := tea.NewProgram(initialModel()).Run()
	if teaErr != nil {
		fmt.Printf("Could not start program :(\n%v\n", teaErr)
		os.Exit(1)
	}

	return
	argQuery := os.Args[1:]
	dynamicArgs := []string{}
	infoChannel := make(chan []byte)

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

	if fileExists("./yt-dlp") {
		ytDlpPath = "./yt-dlp"
	} else if fileExists("./yt-dlp.exe") {
		// todo test on windows
		ytDlpPath = "./yt-dlp.exe"
	}

	var url string
	if len(argQuery) > 0 {
		url = strings.Join(argQuery, " ")
	} else {
		println("enter a query, either a", formatBold("url"), "or something to", formatBold("search on youtube"))
		inputPrompt := textinput.New("query:")
		inputPrompt.InitialValue = getClipboardUrl()

		var inputErr error
		url, inputErr = inputPrompt.RunPrompt()
		if inputErr != nil {
			os.Exit(0)
		}
	}

	inputUrl := url
	if !validateUrl(url) {
		inputUrl = "https://youtube.com/search?q=" + url
		dynamicArgs = append(dynamicArgs, "-I", "1")
	}

	go resolveInfo(dynamicArgs, inputUrl, infoChannel)

	presets := maps.Keys(PRESET_MAP)
	slices.Sort(presets)
	slices.Reverse(presets)

	presetPicker := selection.New("preset:", presets)
	presetPicker.Filter = nil
	preset, presetErr := presetPicker.RunPrompt()
	if presetErr != nil {
		os.Exit(0)
	}

	cpuCount := runtime.NumCPU()
	dynamicArgs = append(dynamicArgs, "-N", fmt.Sprint(cpuCount))

	infoOut := <-infoChannel

	downloadArgs := append(DEFAULT_ARGS[:], PRESET_MAP[preset]...)
	downloadArgs = append(downloadArgs, dynamicArgs...)
	downloadArgs = append(downloadArgs, "-P", downloadPath)
	downloadArgs = append(downloadArgs, "--load-info-json", "-")

	downloadCmd := exec.Command(ytDlpPath, downloadArgs...)
	downloadCmd.Stdin = strings.NewReader(string(infoOut))
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	downloadErr := downloadCmd.Run()
	maybePanic(downloadErr)
}

func resolveInfo(dynamicArgs []string, inputUrl string, ch chan []byte) {
	infoArgs := append(DEFAULT_ARGS[:], "-J")
	infoArgs = append(infoArgs, dynamicArgs...)
	infoArgs = append(infoArgs, inputUrl)

	infoCmd := exec.Command(ytDlpPath, infoArgs...)
	infoOut, infoErr := infoCmd.Output()
	maybePanic(infoErr)

	//? if it's a playlist (like when you search for a song), just use the first entry
	firstEntry, _, _, entryErr := jsonparser.Get(infoOut, "entries", "[0]")
	if firstEntry != nil && entryErr == nil {
		infoOut = firstEntry
	}

	// title, titleErr := jsonparser.GetString(infoOut, "title")
	// if titleErr != nil {
	// 	title = "error extracting title"
	// }

	// println("title:", formatColor(title))

	ch <- infoOut
}
