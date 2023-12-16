package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	cp "github.com/otiai10/copy"
)

const DIST_FOLDER = "./dist"

func main() {
	fmt.Println("creating mac bundle(s)...")

	distContents, err := os.ReadDir(DIST_FOLDER)
	if err != nil {
		panic(err)
	}

	darwinArchFolders := []string{}

	for _, file := range distContents {
		if file.IsDir() && strings.Contains(file.Name(), "darwin") {
			darwinArchFolders = append(darwinArchFolders, file.Name())
		}
	}

	for _, currentArch := range darwinArchFolders {
		fmt.Println("creating bundle for", currentArch)

		appFolder := path.Join(DIST_FOLDER, currentArch, "go-yt-dlp.app")
		currentExecutable := path.Join(DIST_FOLDER, currentArch, "go-yt-dlp")

		cp.Copy("./mac-bundle/Contents", path.Join(appFolder, "Contents"))
		cp.Copy(currentExecutable, path.Join(appFolder, "Contents/MacOS/go-yt-dlp"))
	}

}
