package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
)

const DIST_FOLDER = "./dist"

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("creating mac bundle(s)...")

	distContents, err := os.ReadDir(DIST_FOLDER)
	maybePanic(err)

	darwinArchFolders := []string{}

	for _, file := range distContents {
		if file.IsDir() && strings.Contains(file.Name(), "darwin") {
			darwinArchFolders = append(darwinArchFolders, file.Name())
		}
	}

	for _, currentArch := range darwinArchFolders {
		fmt.Println("creating bundle for", currentArch)

		currentContents, err := os.ReadDir(filepath.Join(DIST_FOLDER, currentArch))
		maybePanic(err)

		appFolder := filepath.Join(DIST_FOLDER, currentArch, "go-yt-dlp.app")
		currentExecutable := filepath.Join(DIST_FOLDER, currentArch, currentContents[0].Name())

		cp.Copy("./mac-bundle/Contents", filepath.Join(appFolder, "Contents"))
		cp.Copy(currentExecutable, filepath.Join(appFolder, "Contents/MacOS/go-yt-dlp"))
	}

}
