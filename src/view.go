package src

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	var s string

	if m.quitting {
		s += "\n"

		if len(downloadedList) > 0 {
			s += defaultStyle.Render("Downloaded")
			s += "\n"

			for _, title := range downloadedList {
				s += accentColorStyle.Render(title)
				s += "\n"
			}
		} else {
			s += defaultStyle.Render("Nothing downloaded.")
			s += "\n"
		}

		return s
	}

	if m.view == querySelect {
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
	}

	if m.view == presetSelect || m.view == downloadView {
		s += defaultStyle.Render("Downloading: ")
		s += boldStyle.Render(m.title)
		s += "\n"
	}

	if m.view == presetSelect {
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

	if m.view == customPreset {
		s += boldStyle.Render("Enter a custom preset:")
		s += "\n"
		s += defaultStyle.Render("(see https://github.com/yt-dlp/yt-dlp#format-selection)")
		s += "\n\n"

		s += m.customPresetInput.View()
		s += "\n"
	}

	if m.view == downloadView {
		preset := PRESET_MAP[m.selectedPreset][0]

		s += defaultStyle.Render("Selected preset: ")
		s += boldStyle.Render(preset)
		s += "\n\n"

		s += defaultStyle.Render("Download logs:")
		s += "\n"
		s += strings.Join(m.downloadLogs, "\n")
		s += "\n"

		if m.downloadProgress >= 0 && !m.downloadDone {
			s += accentColorStyle.Render("Progress: ")
			s += boldStyle.Render(fmt.Sprintf("%.2f%%", m.downloadProgress*100))
			s += "\n"
		}

		if m.downloadDone {
			s += accentColorStyle.Render("Download finished.")
			s += "\n"
		}

	}

	s += "\n"

	if m.downloadDone {
		s += defaultStyle.Render("Press enter to reset, escape to quit.")
		s += "\n"
	} else {
		s += defaultStyle.Render("Press escape to quit.")
		s += "\n"
	}

	return s
}
