package src

import (
	"github.com/buger/jsonparser"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewMap string

const (
	querySelect  viewMap = "QuerySelect"
	presetSelect viewMap = "PresetSelect"
	customPreset viewMap = "CustomPreset"
	downloadView viewMap = "DownloadView"
)

var (
	accentColorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
	boldStyle        = lipgloss.NewStyle().Bold(true)
	defaultStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#9f9f9f"))

	downloadedList = []string{}
)

type model struct {
	view          viewMap
	queryInput    textinput.Model
	title         string
	downloadQuery string
	musicSearch   bool
	quitting      bool

	infoOut            []byte
	downloadLogs       []string
	downloadProgress   float64
	downloadLogChannel chan string
	downloadDone       bool
	progressBar        progress.Model

	presetCursor      int
	selectedPreset    int
	presets           []string
	customPresetInput textinput.Model
}

type downloadLogMsg string

func waitForDownloadLog(downloadChannel chan string) tea.Cmd {
	return func() tea.Msg {
		return downloadLogMsg(<-downloadChannel)
	}
}

func getTitle(infoOut []byte) string {
	title, titleErr := jsonparser.GetString(infoOut, "title")
	if titleErr != nil {
		title = "error extracting title"
	}

	return title
}

func InitialModel() model {
	presets := []string{}
	for _, preset := range PRESET_MAP {
		presets = append(presets, preset[0])
	}

	queryInput := textinput.New()
	queryInput.Placeholder = getClipboardUrl()
	queryInput.Focus()
	queryInput.CharLimit = 0
	queryInput.Width = 44 // length of a full youtube url

	return model{
		view:               querySelect,
		queryInput:         queryInput,
		downloadLogChannel: make(chan string),
		progressBar:        progress.New(progress.WithGradient("#0a0a0a", "#ff0000")),
		selectedPreset:     -1,
		presets:            presets,
		customPresetInput:  textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, setDownloadPath, setExecutablePath)
}
