package src

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/bodgit/sevenzip"
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

func downloadFile(url string, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	os.MkdirAll(filepath.Dir(filename), os.ModePerm)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	return err
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

func findExecutable(name string) (string, error) {
	globalPath, err := exec.LookPath(name)
	if err == nil {
		return globalPath, nil
	}

	currentExecutable, err := os.Executable()
	if err != nil {
		return "", err
	}

	executablePath := filepath.Dir(currentExecutable)

	adjacent := filepath.Join(executablePath)
	adjacentExe := filepath.Join(executablePath + ".exe")
	inBin := filepath.Join(binDir, name)
	inBinExe := filepath.Join(binDir, name+".exe")

	if fileExists(adjacent) {
		return adjacent, nil
	} else if fileExists(adjacentExe) {
		return adjacentExe, nil
	} else if fileExists(inBin) {
		return inBin, nil
	} else if fileExists(inBinExe) {
		return inBinExe, nil
	}

	return "", fmt.Errorf("could not find %s", name)
}

func executableExists(name string) bool {
	_, err := findExecutable(name)
	return err == nil
}

type Files map[string]string

func extractArchive(archive string, files Files) error {
	r, err := sevenzip.OpenReader(archive)
	maybePanic(err)
	defer r.Close()

	for _, file := range r.File {
		filename := filepath.Base(file.Name)
		outPath, isWanted := files[filename]

		if isWanted {
			err = extractFile(file, outPath)
			maybePanic(err)
		}
	}

	return nil
}

func extractFile(file *sevenzip.File, outPath string) error {
	rc, err := file.Open()
	maybePanic(err)
	defer rc.Close()

	out, err := os.Create(outPath)
	maybePanic(err)
	defer out.Close()

	_, err = io.Copy(out, rc)
	maybePanic(err)

	return nil
}

func SetTermTitle(title string) {
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
func scanLinesCR(data []byte, atEOF bool) (advance int, token []byte, err error) {
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

type Spinner struct {
	Frames []string
	FPS    time.Duration
}

var layerSpinner = Spinner{
	Frames: []string{
		"-",
		"=",
		"≡",
	},
	FPS: time.Second / 6,
}
