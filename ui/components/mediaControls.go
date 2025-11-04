package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"thomasjgriffin.dev/radio-garden-cli/internals/helpers"
)

type MediaControlsModel struct {
	Content string
	Style   lipgloss.Style

	MediaPauseButton  MediaButtonModel
	MediaNextButton   MediaButtonModel
	MediaRandomButton MediaButtonModel
}

func MediaControls(intialContent string) MediaControlsModel {
	arrowText, _ := helpers.LoadTextFile("assets/components/arrow.txt")
	pauseText, _ := helpers.LoadTextFile("assets/components/pausePlay.txt")
	randmonIconText, _ := helpers.LoadTextFile("assets/components/random.txt")

	return MediaControlsModel{
		MediaPauseButton:  MediaButton(pauseText),
		MediaNextButton:   MediaButton(arrowText),
		MediaRandomButton: MediaButton(randmonIconText),
	}
}

func (m MediaControlsModel) Init() tea.Cmd {
	return nil
}

func (m MediaControlsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Non-interactive, so it just returns itself
	return m, nil
}

func (m MediaControlsModel) View() string {
	// Render the content, then wrap it in the box style

	style := lipgloss.JoinHorizontal(lipgloss.Center,
		m.MediaRandomButton.View(),
		m.MediaPauseButton.View(),
		m.MediaNextButton.View(),
	)
	content := style
	return m.Style.Render(content)
}
