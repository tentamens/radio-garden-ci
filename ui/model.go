// Package ui holds the main ui of the TUI
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"thomasjgriffin.dev/radio-garden-cli/ui/components"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}

	WorldMap      components.WorldMapModel
	StationSearch components.StationSearchModel

	width  int
	height int
}

func InitialModel() model {
	return model{
		WorldMap:      components.WorldMap(""),
		StationSearch: components.StationSearch(""),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		// exit program on these key pressed
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	newStationSearchModel, stationSearchCmd := m.StationSearch.Update(msg)
	m.StationSearch = newStationSearchModel.(components.StationSearchModel)
	cmds = append(cmds, stationSearchCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			m.WorldMap.View(),
			m.StationSearch.View(),
		),
	)
}
