package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MediaButtonModel struct {
	Content string
	Style   lipgloss.Style

	MediaButtonSkip MediaControlsModel
}

func MediaButton(mediaButtonString string) MediaButtonModel {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 1)

	return MediaButtonModel{
		Content: mediaButtonString,
		Style:   style,
	}
}

func (m MediaButtonModel) Init() tea.Cmd {
	return nil
}

func (m MediaButtonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m MediaButtonModel) View() string {
	content := m.Content
	return m.Style.Render(content)
}
