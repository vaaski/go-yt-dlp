package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type View string

const (
	QuerySelect  View = "QuerySelect"
	PresetSelect View = "PresetSelect"
	DownloadView View = "DownloadView"
)

var (
	accentColorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
	boldStyle        = lipgloss.NewStyle().Bold(true)
	defaultStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#9f9f9f"))
)

type model struct {
	view          View
	textInput     textinput.Model
	title         string
	downloadQuery string
	quitting      bool

	infoOut []byte

	presetCursor   int
	selectedPreset string
	presets        []string
}

type infoMsg []byte

func fetchInfo(url string) tea.Cmd {
	return func() tea.Msg {
		infoArgs := append(DEFAULT_ARGS[:], "-J")

		if !validateUrl(url) {
			url = "https://youtube.com/search?q=" + url
			infoArgs = append(infoArgs, "-I", "1")
		}

		infoArgs = append(infoArgs, url)

		infoCmd := exec.Command(ytDlpPath, infoArgs...)
		infoOut, infoErr := infoCmd.Output()
		maybePanic(infoErr)

		//? if it's a playlist (like when you search for a song), just use the first entry
		firstEntry, _, _, entryErr := jsonparser.Get(infoOut, "entries", "[0]")
		if firstEntry != nil && entryErr == nil {
			infoOut = firstEntry
		}

		return infoMsg(infoOut)
	}
}

func startDownload(infoOut []byte, preset string) tea.Cmd {
	return func() tea.Msg {
		downloadArgs := append(DEFAULT_ARGS[:], PRESET_MAP[preset]...)
		downloadArgs = append(downloadArgs, "-P", downloadPath)

		cpuCount := runtime.NumCPU()
		downloadArgs = append(downloadArgs, "-N", fmt.Sprint(cpuCount))
		downloadArgs = append(downloadArgs, "--load-info-json", "-")

		downloadCmd := exec.Command(ytDlpPath, downloadArgs...)
		downloadCmd.Stdin = strings.NewReader(string(infoOut))
		downloadCmd.Stdout = os.Stdout
		downloadCmd.Stderr = os.Stderr

		downloadErr := downloadCmd.Run()
		maybePanic(downloadErr)
		return nil
	}
}

func getTitle(infoOut []byte) string {
	title, titleErr := jsonparser.GetString(infoOut, "title")
	if titleErr != nil {
		title = "error extracting title"
	}

	return title
}

func initialModel() model {
	presets := maps.Keys(PRESET_MAP)
	slices.Sort(presets)
	slices.Reverse(presets)

	ti := textinput.New()
	ti.Placeholder = getClipboardUrl()
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 44 // length of a full youtube url

	return model{
		view:      QuerySelect,
		presets:   presets,
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	case infoMsg:
		m.infoOut = msg
		m.title = getTitle(msg)
		if m.selectedPreset != "" && m.view == DownloadView {
			return m, startDownload(msg, m.selectedPreset)
		}
	}

	if m.view == QuerySelect {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.downloadQuery = m.textInput.Value()
				if m.downloadQuery == "" {
					m.downloadQuery = m.textInput.Placeholder
				}

				if m.downloadQuery != "" {
					m.textInput.Blur()
					m.view = PresetSelect
					m.title = m.downloadQuery
					return m, fetchInfo(m.downloadQuery)
				}
			}
		}

		m.textInput, cmd = m.textInput.Update(msg)
	} else if m.view == PresetSelect {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {

			case "up", "k":
				if m.presetCursor > 0 {
					m.presetCursor--
				}

			case "down", "j":
				if m.presetCursor < len(m.presets)-1 {
					m.presetCursor++
				}

			case " ", "enter":
				m.selectedPreset = m.presets[m.presetCursor]
				m.view = DownloadView
				if m.infoOut != nil {
					return m, startDownload(m.infoOut, m.selectedPreset)
				}
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "\nGoodbye!\n"
	}

	var s string

	if m.view == QuerySelect {
		s += defaultStyle.Render("enter either a ")
		s += boldStyle.Render("url")
		s += defaultStyle.Render(" or something to ")
		s += boldStyle.Render("search on youtube")
		s += "\n\n"

		s += m.textInput.View()
		s += "\n"
	}

	if m.view == PresetSelect || m.view == DownloadView {
		s += defaultStyle.Render("Downloading: ")
		s += boldStyle.Render(m.title)
		s += "\n"
	}

	if m.view == PresetSelect {
		s += defaultStyle.Render("Pick a preset:")
		s += "\n\n"

		for i, preset := range m.presets {
			if m.presetCursor == i {
				s += "> " + accentColorStyle.Render((preset))
			} else {
				s += "  " + preset
			}
			s += "\n"
		}
	}

	if m.view == DownloadView {
		s += defaultStyle.Render("Selected preset: ")
		s += boldStyle.Render(m.selectedPreset)
		s += "\n\n"
		s += defaultStyle.Render("this is download page")
		s += "\n"
	}

	s += "\n"
	s += defaultStyle.Render("Press escape to quit.")
	s += "\n"

	return s
}
