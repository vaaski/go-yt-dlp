package main

import (
	"net/url"
	"os"

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
