package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/buger/jsonparser"
	"golang.design/x/clipboard"
)

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func validateUrl(inputUrl string) bool {
	parsed, err := url.ParseRequestURI(inputUrl)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func setTermTitle(title string) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", "title", title)
		if err := cmd.Run(); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		print("\033]0;" + title + "\007")
	}
}

// takes the output of yt-dlp --progress-template and returns
// the progress as a float. returns -1 if there's an error
func parseProgressOutput(line string) float64 {
	truncated := line[len(PROGRESS_PREFIX):]

	downloaded, err := jsonparser.GetFloat([]byte(truncated), "downloaded_bytes")
	if err != nil {
		return -1
	}

	total, err := jsonparser.GetFloat([]byte(truncated), "total_bytes")
	if err != nil {
		return -1
	}

	return downloaded / total
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// required to parse the output of yt-dlp --progress-template
// because it uses \r instead of \n for newlines
func ScanLinesCR(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a carriage return-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
