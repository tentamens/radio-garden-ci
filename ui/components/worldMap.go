package components

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WorldMapModel struct {
	Content string
	Style   lipgloss.Style
}

func WorldMap(intialContent string) WorldMapModel {
	asciiMapBytes, err := os.ReadFile("assets/worldMaps/defaultWorldMap.txt")
	if err != nil {
		log.Fatalf("could not read map file: %v", err)
	}

	asciiMap := string(asciiMapBytes)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 1)

	return WorldMapModel{
		Content: asciiMap,
		Style:   style,
	}
}

func (m WorldMapModel) Init() tea.Cmd {
	return nil
}

func (m WorldMapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Non-interactive, so it just returns itself
	return m, nil
}

func (m WorldMapModel) View() string {
	// Render the content, then wrap it in the box style
	content := m.Content
	return m.Style.Render(content)
}
