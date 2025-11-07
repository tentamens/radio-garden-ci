// Package ui holds the main ui of the TUI
package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"thomasjgriffin.dev/radio-garden-cli/internals/helpers"
	"thomasjgriffin.dev/radio-garden-cli/ui/components"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}

	WorldMap                components.WorldMapModel
	InteractionsContainer   components.InteractionsContainerModel
	StationDetailsComponent components.StationDetailsComponentModel

	width          int
	height         int
	placeID        string
	stationID      string
	stationDetails helpers.StationDetails

	cancelStream context.CancelFunc
}

func InitialModel() model {
	// initialize api client
	helpers.InitClient()

	return model{
		WorldMap:                components.WorldMap(""),
		InteractionsContainer:   components.InteractionsContainer(""),
		StationDetailsComponent: components.StationDetailsComponent(""),
		placeID:                 "p6yf8OtF",
		stationID:               "",
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

		// start stream on enter place
		case "enter":
			if m.cancelStream != nil {
				return m, nil
			}

			// add cancel stream to state to call at a later time
			ctx, cancel := context.WithCancel(context.Background())

			m.cancelStream = cancel

			return m, helpers.StreamMusic(m.stationID, ctx)

		// cancel the stream
		case "p":
			if m.cancelStream != nil {
				m.cancelStream()
				m.cancelStream = nil
			}

			return m, nil

		case "t":

		case "r":
			m.placeID = helpers.PickRandonPlace()
			stationID, err := helpers.PickStation(context.Background(), m.placeID)
			if err != nil {
				return m, nil
			}
			m.stationID = stationID
			m.stationDetails, err = helpers.GetStationDetails(m.stationID)
			if err != nil {
				return m, nil
			}

			updateCmd := func() tea.Msg {
				return components.NewStationDetailsMsg(m.stationDetails)
			}

			cmds = append(cmds, updateCmd)

		}
	}

	newStationSearchModel, stationSearchCmd := m.InteractionsContainer.Update(msg)
	m.InteractionsContainer = newStationSearchModel.(components.InteractionsContainerModel)
	newStationDetailsModel, stationDetailsCmd := m.StationDetailsComponent.Update(msg)
	m.StationDetailsComponent = newStationDetailsModel.(components.StationDetailsComponentModel)

	cmds = append(cmds, stationSearchCmd, stationDetailsCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			m.WorldMap.View(),
			lipgloss.JoinHorizontal(lipgloss.Center,
				m.InteractionsContainer.View(),
				m.StationDetailsComponent.View(),
			),
		),
	)
}
