package main

import (
	"strings"

	"github.com/buger/jsonparser"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	musicSearch   bool
	quitting      bool

	infoOut            []byte
	downloadLogs       []string
	downloadLogChannel chan string
	downloadDone       bool

	presetCursor   int
	selectedPreset int
	presets        []string
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

func initialModel() model {
	presets := []string{}
	for _, preset := range PRESET_MAP {
		presets = append(presets, preset[0])
	}

	ti := textinput.New()
	ti.Placeholder = getClipboardUrl()
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 44 // length of a full youtube url

	return model{
		view:               QuerySelect,
		presets:            presets,
		textInput:          ti,
		downloadLogChannel: make(chan string),
		selectedPreset:     -1,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, setDownloadPath, setExecutablePath)
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
		if m.selectedPreset >= 0 && m.view == DownloadView {
			return m, tea.Batch(startDownload(m), waitForDownloadLog(m.downloadLogChannel))
		}

	case downloadLogMsg:
		if strings.TrimSpace(string(msg)) != "" {
			m.downloadLogs = append(m.downloadLogs, string(msg))
		}

		if !m.downloadDone {
			return m, waitForDownloadLog(m.downloadLogChannel)
		}

	case downloadFinishMsg:
		m.downloadDone = true
		m.downloadLogs = append(m.downloadLogs, accentColorStyle.Render("Download finished."))
	}

	if m.view == QuerySelect {
		musicToggled := false

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "M":
				musicToggled = true
				m.musicSearch = !m.musicSearch

			case "enter":
				m.downloadQuery = m.textInput.Value()
				if m.downloadQuery == "" {
					m.downloadQuery = m.textInput.Placeholder
				}

				if m.downloadQuery != "" {
					m.textInput.Blur()
					m.view = PresetSelect
					m.title = m.downloadQuery
					return m, fetchInfo(m)
				}
			}
		}

		if !musicToggled {
			m.textInput, cmd = m.textInput.Update(msg)
		}
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
				m.selectedPreset = m.presetCursor
				m.view = DownloadView
				if m.infoOut != nil {
					return m, tea.Batch(startDownload(m), waitForDownloadLog(m.downloadLogChannel))
				}
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	var s string

	if m.quitting {
		if m.downloadDone {
			s += "\n"
			s += defaultStyle.Render("Downloaded")
			s += "\n"
			s += accentColorStyle.Render(getTitle(m.infoOut))
			// s += "\n\n"
			// s += defaultStyle.Render("To destination")
			// s += "\n"
			// s += accentColorStyle.Render("todo")
			s += "\n"
		} else {
			s += "\n"
			s += defaultStyle.Render("Nothing downloaded.")
			s += "\n"
		}

		return s
	}

	if m.view == QuerySelect {
		s += defaultStyle.Render("Enter either a ")
		s += boldStyle.Render("url")
		s += defaultStyle.Render(" or something to ")
		s += boldStyle.Render("search on youtube")

		if m.musicSearch {
			s += "\n"
			s += defaultStyle.Render("(music search enabled)")
		}

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
		preset := PRESET_MAP[m.selectedPreset][0]

		s += defaultStyle.Render("Selected preset: ")
		s += boldStyle.Render(preset)
		s += "\n\n"
		s += defaultStyle.Render("Download logs:")
		s += "\n"
		s += strings.Join(m.downloadLogs, "\n")
		s += "\n"
	}

	s += "\n"
	s += defaultStyle.Render("Press escape to quit.")
	s += "\n"

	return s
}
