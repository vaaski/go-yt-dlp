package src

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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
		if m.selectedPreset >= 0 && m.view == downloadView {
			return m, tea.Batch(startDownload(m), waitForDownloadLog(m.downloadLogChannel))
		}

	case downloadLogMsg:
		stringified := strings.TrimSpace(string(msg))

		if stringified != "" {
			if strings.HasPrefix(stringified, PROGRESS_PREFIX) {
				m.downloadProgress = parseProgressOutput(stringified)
			} else {
				m.downloadLogs = append(m.downloadLogs, stringified)
			}
		}

		if !m.downloadDone {
			return m, waitForDownloadLog(m.downloadLogChannel)
		}

	case downloadFinishMsg:
		downloadedList = append(downloadedList, m.title)
		m.downloadDone = true
	}

	if m.view == querySelect {
		musicToggled := false

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "M":
				musicToggled = true
				m.musicSearch = !m.musicSearch

			case "enter":
				m.downloadQuery = m.queryInput.Value()
				if m.downloadQuery == "" {
					m.downloadQuery = m.queryInput.Placeholder
				}

				if m.downloadQuery != "" {
					m.queryInput.Blur()
					m.view = presetSelect
					m.title = m.downloadQuery
					return m, fetchInfo(m)
				}
			}
		}

		if !musicToggled {
			m.queryInput, cmd = m.queryInput.Update(msg)
		}

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
				m = InitialModel()
			}
		}
	}

	return m, cmd
}
