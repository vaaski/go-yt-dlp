package src

import tea "github.com/charmbracelet/bubbletea"

func presetView(m *model) string {
	var s string

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

	return s
}

func presetController(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
					return m, tea.Batch(startDownload(*m), waitForDownloadLog(m.downloadLogChannel))
				}
			}
		}
	}

	return m, nil
}
