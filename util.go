package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"

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
