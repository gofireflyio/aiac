package libaiac

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	spinner  spinner.Model
	quitting bool
	input    string
	err      error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil

	case string:
		if len(msg) > 0 {
			m.quitting = true
			m.input = msg
			return m, tea.Quit
		}
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Generating code...\n\n", m.spinner.View())
	if m.quitting {
		return m.input
	}
	return str
}
