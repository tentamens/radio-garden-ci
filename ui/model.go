// Package ui holds the main ui of the TUI
package ui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"thomasjgriffin.dev/radio-garden-cli/internals/helpers"
	"thomasjgriffin.dev/radio-garden-cli/ui/components"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}

	WorldMap              components.WorldMapModel
	InteractionsContainer components.InteractionsContainerModel

	width  int
	height int

	cancelStream context.CancelFunc
}

func InitialModel() model {
	// initialize api client
	helpers.InitClient()

	return model{
		WorldMap:              components.WorldMap(""),
		InteractionsContainer: components.InteractionsContainer(""),
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

			return m, helpers.StreamMusic("AHWRUh8V", ctx)

		// cancel the stream
		case "p":
			if m.cancelStream != nil {
				fmt.Println("cancling stream")
				m.cancelStream()
				m.cancelStream = nil
			}

			return m, nil

		case "t":
			fmt.Println(helpers.PickRandonStation(), " random station picked")
		}

	}

	newStationSearchModel, stationSearchCmd := m.InteractionsContainer.Update(msg)
	m.InteractionsContainer = newStationSearchModel.(components.InteractionsContainerModel)
	cmds = append(cmds, stationSearchCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			m.WorldMap.View(),
			m.InteractionsContainer.View(),
		),
	)
}
