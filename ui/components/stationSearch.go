package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StationSearchModel struct {
	textInput textinput.Model
	Content   string
	Style     lipgloss.Style
}

func StationSearch(intialContent string) StationSearchModel {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).Width(40)

	ti := textinput.New()
	ti.Placeholder = "Search Station"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return StationSearchModel{
		Style:     style,
		textInput: ti,
	}
}

func (m StationSearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StationSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m StationSearchModel) View() string {
	// Render the content, then wrap it in the box style
	content := m.textInput.View()
	return m.Style.Render(content)
}
