package src

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progressBar.Width = msg.Width - 4
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}

	case infoMsg:
		m.infoOut = msg
		m.title = getTitle(msg)
		if m.selectedPreset >= 0 && m.view == downloadView {
			return m, tea.Batch(startDownload(m), waitForDownloadLog(m.downloadLogChannel))
		}

	case downloadLogMsg:
		stringified := strings.TrimSpace(string(msg))
		var progressUpdate tea.Cmd

		if stringified != "" {
			if strings.HasPrefix(stringified, PROGRESS_PREFIX) {
				m.downloadProgress = parseProgressOutput(stringified)
				progressUpdate = m.progressBar.SetPercent(m.downloadProgress)
			} else {
				m.downloadLogs = append(m.downloadLogs, stringified)
			}
		}

		if !m.downloadDone {
			return m, tea.Batch(waitForDownloadLog(m.downloadLogChannel), progressUpdate)
		}

	case progress.FrameMsg:
		progressModel, cmd := m.progressBar.Update(msg)
		m.progressBar = progressModel.(progress.Model)
		return m, cmd

	case spinner.TickMsg:
		if m.infoOut == nil {
			var cmd tea.Cmd
			m.infoFetchSpinner, cmd = m.infoFetchSpinner.Update(msg)
			return m, cmd
		}

	case downloadFinishMsg:
		downloadedList = append(downloadedList, m.title)
		m.downloadDone = true
		SetTermTitle("go-yt-dlp")
	}

	if m.view == querySelect {
		return queryController(&m, msg)
	} else if m.view == presetSelect {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {

			case "up", "k":
				if m.presetCursor > 0 {
					m.presetCursor--
				} else {
					m.presetCursor = len(m.presets) - 1
				}

			case "down", "j":
				if m.presetCursor < len(m.presets)-1 {
					m.presetCursor++
				} else {
					m.presetCursor = 0
				}

			case " ", "enter":
				m.selectedPreset = m.presetCursor

				if PRESET_MAP[m.selectedPreset][0] == CUSTOM_PRESET {
					m.view = customPreset
					m.customPresetInput.Focus()
				} else {
					m.view = downloadView
					if m.infoOut != nil {
						return m, tea.Batch(startDownload(m), waitForDownloadLog(m.downloadLogChannel))
					}
				}
			}
		}

	} else if m.view == customPreset {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {

			case "enter":
				format := m.customPresetInput.Value()

				if format != "" {
					m.customPresetInput.Blur()
					PRESET_MAP[m.selectedPreset] = append(PRESET_MAP[m.selectedPreset], format)

					m.view = downloadView
					if m.infoOut != nil {
						return m, tea.Batch(startDownload(m), waitForDownloadLog(m.downloadLogChannel))
					}
				}
			}
		}

		m.customPresetInput, cmd = m.customPresetInput.Update(msg)
	} else if m.view == downloadView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case " ", "enter":
				lastProgressWidth := m.progressBar.Width

				m = InitialModel()
				m.progressBar.Width = lastProgressWidth
			}
		}
	}

	return m, cmd
}
