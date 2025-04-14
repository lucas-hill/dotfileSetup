package selectlist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	name string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

type model struct {
	list     list.Model
	selected map[string]bool
	done     bool
}

// MultiSelect runs the multiselect prompt and returns selected item names
func MultiSelect(title string, options []string) ([]string, error) {
	items := make([]list.Item, len(options))
	for i, o := range options {
		items[i] = item{name: o}
	}

	l := list.New(items, list.NewDefaultDelegate(), 40, len(items)+2)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	m := model{
		list:     l,
		selected: map[string]bool{},
	}

	prog := tea.NewProgram(m)
	finalModel, err := prog.Run()
	if err != nil {
		return nil, err
	}

	if finalModel, ok := finalModel.(model); ok {
		var selected []string
		for k, v := range finalModel.selected {
			if v {
				selected = append(selected, k)
			}
		}
		return selected, nil
	}

	return nil, fmt.Errorf("unexpected model type")
}

// Init for tea.Model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles key events
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case " ":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected[i.name] = !m.selected[i.name]
			}

		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the list with checkboxes
func (m model) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.list.Title + "\n\n")

	for i, li := range m.list.Items() {
		it := li.(item)
		cursor := " " // no cursor
		if i == m.list.Index() {
			cursor = ">" // current selection
		}
		checked := " " // unchecked
		if m.selected[it.name] {
			checked = "x"
		}
		fmt.Fprintf(&b, "%s [%s] %s\n", cursor, checked, it.name)
	}

	b.WriteString("\n[space] to toggle • [enter] to confirm • [q] to quit\n")
	return b.String()
}
