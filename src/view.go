package src

import (
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
		s += queryView(&m)
	}

	if m.view == presetSelect || m.view == downloadView {
		if m.infoOut == nil {
			s += m.infoFetchSpinner.View()
			s += " "
			s += defaultStyle.Render("Resolving: ")
		} else {
			s += defaultStyle.Render("Downloading: ")
		}

		s += boldStyle.Render(m.title)
		s += "\n"
	}

	if m.view == presetSelect {
		s += presetView(&m)
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

		if m.downloadProgress > 0 && !m.downloadDone {
			s += m.progressBar.View()
			s += "\n\n"
		}

		if m.downloadDone {
			s += accentColorStyle.Render("Download finished.")
			s += "\n\n"
		}

		s += defaultStyle.Render("Download logs:")
		s += "\n"
		s += strings.Join(m.downloadLogs, "\n")
		s += "\n"

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
