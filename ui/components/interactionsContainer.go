package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InteractionsContainerModel struct {
	Content string
	Style   lipgloss.Style

	StationSearch StationSearchModel
	MediaControls MediaControlsModel
}

func InteractionsContainer(intialContent string) InteractionsContainerModel {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 1)

	return InteractionsContainerModel{
		Style:         style,
		StationSearch: StationSearch(""),
		MediaControls: MediaControls(""),
	}
}

func (m InteractionsContainerModel) Init() tea.Cmd {
	return nil
}

func (m InteractionsContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Non-interactive, so it just returns itself
	return m, nil
}

func (m InteractionsContainerModel) View() string {
	// Render the content, then wrap it in the box style
	content := lipgloss.NewStyle().Margin(0, 2).Render(
		lipgloss.JoinVertical(lipgloss.Center,
			m.StationSearch.View(),
			m.MediaControls.View(),
		))
	return m.Style.Render(content)
}
