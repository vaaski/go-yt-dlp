//go:build cgo

package clipboard

import (
	"net/url"

	"golang.design/x/clipboard"
)

func GetClipboardUrl() string {
	err := clipboard.Init()
	if err != nil {
		return ""
	}

	text := clipboard.Read(clipboard.FmtText)
	stringified := string(text)

	parsed, err := url.ParseRequestURI(stringified)

	if err == nil && parsed.Scheme != "" && parsed.Host != "" {
		return stringified
	} else {
		return ""
	}
}
