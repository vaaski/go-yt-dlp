package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/buger/jsonparser"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

var DEFAULT_ARGS = [...]string{"--force-keyframes-at-cuts", "-P", "ytdl-download", "--embed-metadata", "--no-playlist", "--console-title"}
var PRESET_MAP = map[string][]string{
	"mp4":      {"--remux-video", "mp4"},
	"mp4-fast": {"-f", "b"},
	"mp3":      {"-x", "--audio-format", "mp3", "-o", "%(uploader)s - %(title)s.%(ext)s"},
}

func main() {
	dynamicArgs := []string{}

	println("enter a query, either a", formatBold("url"), "or something to", formatBold("search on youtube"))

	inputPrompt := textinput.New("query:")
	inputPrompt.InitialValue = getClipboardUrl()

	url, inputErr := inputPrompt.RunPrompt()
	if inputErr != nil {
		os.Exit(0)
	}

	presets := maps.Keys(PRESET_MAP)
	slices.Sort(presets)
	slices.Reverse(presets)

	presetPicker := selection.New("preset:", presets)
	presetPicker.Filter = nil
	preset, presetErr := presetPicker.RunPrompt()
	if presetErr != nil {
		os.Exit(0)
	}

	inputUrl := url
	if !validateUrl(url) {
		inputUrl = "https://youtube.com/search?q=" + url
		dynamicArgs = append(dynamicArgs, "-I", "1")
	}

	cpuCount := runtime.NumCPU()
	dynamicArgs = append(dynamicArgs, "-N", fmt.Sprint(cpuCount))

	infoArgs := append(DEFAULT_ARGS[:], "-J")
	infoArgs = append(infoArgs, dynamicArgs...)
	infoArgs = append(infoArgs, inputUrl)

	infoCmd := exec.Command("yt-dlp", infoArgs...)
	infoOut, infoErr := infoCmd.Output()
	maybePanic(infoErr)

	//? if it's a playlist (like when you search for a song), just use the first entry
	firstEntry, _, _, entryErr := jsonparser.Get(infoOut, "entries", "[0]")
	if firstEntry != nil && entryErr == nil {
		infoOut = firstEntry
	}

	title, titleErr := jsonparser.GetString(infoOut, "title")
	if titleErr != nil {
		title = "error extracting title"
	}

	println("title:", formatColor(title))

	downloadArgs := append(DEFAULT_ARGS[:], PRESET_MAP[preset]...)
	downloadArgs = append(downloadArgs, dynamicArgs...)
	downloadArgs = append(downloadArgs, "--load-info-json", "-")

	downloadCmd := exec.Command("yt-dlp", downloadArgs...)
	downloadCmd.Stdin = strings.NewReader(string(infoOut))
	downloadCmd.Stdout = os.Stdout
	downloadErr := downloadCmd.Run()
	maybePanic(downloadErr)
}
