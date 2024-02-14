package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	borcolor   lipgloss.Color
	inputstyle lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.borcolor = lipgloss.Color("36")
	s.inputstyle = lipgloss.NewStyle().BorderForeground(s.borcolor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type model struct {
	index     int
	questions []string
	width     int
	height    int
	input     textinput.Model
	style     *Styles
}

func New(questions []string) *model {
	inputwindow := textinput.New()
	styles := DefaultStyles()
	inputwindow.Placeholder = "Your answer"
	return &model{questions: questions, input: inputwindow, style: styles}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "loading"
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.questions[m.index],
		m.input.View())
}

func main() {
	questions := []string{"Whats ur name", "how are you"}
	m := New(questions)
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return
	}
	defer f.Close()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return
	}
}
