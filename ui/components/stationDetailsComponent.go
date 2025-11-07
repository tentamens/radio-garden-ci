package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"thomasjgriffin.dev/radio-garden-cli/internals/helpers"
)

type StationDetailsComponentModel struct {
	Content        string
	Style          lipgloss.Style
	stationDetails helpers.StationDetails
}

type NewStationDetailsMsg helpers.StationDetails

func StationDetailsComponent(intialContent string) StationDetailsComponentModel {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).Width(34)

	return StationDetailsComponentModel{
		Style: style,
	}
}

func (m StationDetailsComponentModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StationDetailsComponentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case NewStationDetailsMsg:
		// Update the model's stationDetails with the data from the message
		m.stationDetails = helpers.StationDetails(msg)
	}

	return m, cmd
}

func (m StationDetailsComponentModel) View() string {
	// Render the content, then wrap it in the box style
	content := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			m.stationDetails.Title,
			m.stationDetails.City+" "+m.stationDetails.Country,
		),
	)

	return m.Style.Render(content)
}
