package components

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MediaControlsModel struct {
	Content string
	Style   lipgloss.Style
}

func MediaControls(intialContent string) MediaControlsModel {
	asciiMapBytes, err := os.ReadFile("assets/worldMaps/defaultWorldMap.txt")
	if err != nil {
		log.Fatalf("could not read map file: %v", err)
	}

	asciiMap := string(asciiMapBytes)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 1)

	return MediaControlsModel{
		Content: asciiMap,
		Style:   style,
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

	style := lipgloss.JoinHorizontal(lipgloss.Center)

	content := m.Content
	return m.Style.Render(content)
}
