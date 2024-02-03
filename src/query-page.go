package src

import tea "github.com/charmbracelet/bubbletea"

func queryView(m *model) string {
	var s string

	if m.musicSearch {
		s += defaultStyle.Render("Enter something to ")
		s += boldStyle.Render("search on ")
		s += accentColorStyle.Render("YouTube Music")
	} else {
		s += defaultStyle.Render("Enter either a ")
		s += boldStyle.Render("URL")
		s += defaultStyle.Render(" or something to ")
		s += boldStyle.Render("search on ")
		s += accentColorStyle.Render("YouTube")
	}

	s += "\n\n"

	s += m.queryInput.View()
	s += "\n\n"
	s += defaultStyle.Render("Press Tab to toggle YouTube Music search.")

	return s
}

func queryController(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var musicToggled bool

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.downloadQuery = m.queryInput.Value()
			if m.downloadQuery == "" {
				m.downloadQuery = m.queryInput.Placeholder
			}

			if m.downloadQuery != "" {
				m.queryInput.Blur()
				m.view = presetSelect
				m.title = m.downloadQuery
				return m, tea.Batch(fetchInfo(*m), m.infoFetchSpinner.Tick)
			}
		case tea.KeyTab:
			musicToggled = true
			m.musicSearch = !m.musicSearch
		}
	}

	if !musicToggled {
		m.queryInput, cmd = m.queryInput.Update(msg)
		return m, cmd
	}

	return m, nil
}
