package main

import (
	
	"github.com/emersion/go-ical"
	// "time"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "strings"
	// "errors"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	// titleStyle = lipgloss.NewStyle().
	// 		Foreground(lipgloss.Color("#FFFDF5")).
	// 		Background(lipgloss.Color("#25A065")).
	// 		Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

// buttons
func (m model) newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, ml *list.Model) tea.Cmd {
		var title string
		// var todoUID string
		var todoIcal ical.Event

		if i, ok := ml.SelectedItem().(TODO); ok {
			title = i.Title()
			// todoUID = i.UID()
			todoIcal = ical.Event(i)
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				err := m.CompleteTodo(todoIcal)
				if err != nil {return errHandler_tui(err,"can't complete item")}
				index := ml.Index()
				ml.RemoveItem(index)
				if len(ml.Items()) == 0 {
					keys.choose.SetEnabled(false)
				}
				return ml.NewStatusMessage(statusMessageStyle("You chose " + title))

			case key.Matches(msg, keys.remove):
				err := m.DelTodo(todoIcal)
				if err != nil {return errHandler_tui(err,"can't delete item")}
				index := ml.Index()
				ml.RemoveItem(index)
				if len(ml.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return ml.NewStatusMessage(statusMessageStyle("Deleted " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}


type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}


func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}
