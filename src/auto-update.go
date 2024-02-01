package src

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/minio/selfupdate"
)

func AutoUpdate() error {
	baseUrl := "https://github.com/vaaski/go-yt-dlp/releases/latest/download/go-yt-dlp-"
	downloadUrl := baseUrl + runtime.GOOS + "-" + runtime.GOARCH

	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully updated to latest version")
	return err
}
