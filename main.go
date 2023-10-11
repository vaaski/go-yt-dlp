package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"golang.design/x/clipboard"
)

var DEFAULT_ARGS = [...]string{"--force-keyframes-at-cuts", "-P", "ytdl-download", "--embed-metadata", "--no-playlist", "--console-title"}
var PRESET_MAP = map[string][]string{
	"mp4":      {"--remux-video", "mp4"},
	"mp4-fast": {"-f", "b"},
	"mp3":      {"-x", "--audio-format", "mp3", "-o", "%(uploader)s - %(title)s.%(ext)s"},
}

func validateUrl(inputUrl string) bool {
	_, err := url.ParseRequestURI(inputUrl)
	return err == nil
}

func getClipboardUrl() string {
	err := clipboard.Init()
	maybePanic(err)

	text := clipboard.Read(clipboard.FmtText)
	stringified := string(text)
	validUrl := validateUrl(stringified)

	if validUrl {
		return stringified
	} else {
		return ""
	}
}

func main() {
	dynamicArgs := []string{}

	inputPrompt := textinput.New("url:")
	inputPrompt.InitialValue = getClipboardUrl()

	url, inputErr := inputPrompt.RunPrompt()
	if inputErr != nil {
		os.Exit(0)
	}

	inputUrl := url
	if !validateUrl(url) {
		inputUrl = "https://youtube.com/search?q=" + url
		dynamicArgs = append(dynamicArgs, "-I", "1")
	}

	cpuCount := runtime.NumCPU()
	dynamicArgs = append(dynamicArgs, "-N", fmt.Sprint(cpuCount))

	println("Input url:", inputUrl)

	infoArgs := append(DEFAULT_ARGS[:], "-J")
	infoArgs = append(infoArgs, dynamicArgs...)
	infoArgs = append(infoArgs, inputUrl)

	println("finalArgs:", strings.Join(infoArgs, " "))
	infoCmd := exec.Command("yt-dlp", infoArgs...)
	infoOut, err := infoCmd.Output()
	maybePanic(err)

	os.WriteFile("info.json", infoOut, 0644)

	title, titleErr := jsonparser.GetString(infoOut, "entries", "[0]", "title")
	if titleErr != nil {

		title, titleErr = jsonparser.GetString(infoOut, "title")
		if titleErr != nil {
			title = "error extracting title"
		}
	}

	println("title:", title)

	presets := make([]string, len(PRESET_MAP))
	i := 0
	for k := range PRESET_MAP {
		presets[i] = k
		i++
	}

	presetPicker := selection.New("Select preset", presets)
	presetPicker.Filter = nil
	preset, presetErr := presetPicker.RunPrompt()
	if presetErr != nil {
		os.Exit(0)
	}

	println("Selected preset:", strings.Join(PRESET_MAP[preset], " "))
}
