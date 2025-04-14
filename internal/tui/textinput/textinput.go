package textinput

import (
	"strings"

	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Label   string
	Input   ti.Model
	Done    bool
	Entered string
	Quit    bool
}

func New(label string, placeholder string, initial string) Model {
	tiModel := ti.New()
	tiModel.Placeholder = placeholder
	tiModel.SetValue(initial)
	tiModel.Focus()
	tiModel.CharLimit = 256
	tiModel.Width = 50

	return Model{
		Label:   label,
		Input:   tiModel,
		Done:    false,
		Entered: "",
	}
}

func (m Model) Init() tea.Cmd {
	return ti.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.Done = true
			m.Entered = strings.TrimSpace(m.Input.Value())
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.Quit = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Done || m.Quit {
		return ""
	}
	return m.Label + "\nInput: " + m.Input.View() + "\n(Press Enter to confirm)"
}
