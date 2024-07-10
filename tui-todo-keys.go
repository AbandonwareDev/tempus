package main

import (
	"github.com/emersion/go-ical"
	// "time"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "fmt"
	// "io"
	// "strings"
	// "errors"
	// "golang.org/x/term"
	// "github.com/muesli/reflow/truncate"
	// "strings"
)

// const (
// 	bullet   = "•"
// 	ellipsis = "…"
// )


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

// type DefaultDelegate struct {
// 	ShowDescription bool
// 	Styles          list.DefaultItemStyles
// 	UpdateFunc      func(tea.Msg, *list.Model) tea.Cmd
// 	ShortHelpFunc   func() []key.Binding
// 	FullHelpFunc    func() [][]key.Binding
// 	// contains filtered or unexported fields
// }
var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyleList         = lipgloss.NewStyle().PaddingLeft(4)
	// selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	// paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	// helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
// func (d DefaultDelegate) Height() int                               { return 1 }
// func (d DefaultDelegate) Spacing() int                              { return 0 }
// func (d DefaultDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
// func (d DefaultDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
// 	i, ok := listItem.(item)
// 	if !ok {
// 		return
// 	}
// 
// 	str := fmt.Sprintf("%d. %s", index+1, i)
// 
// 	// fn := itemStyleList.Foreground(TODO(listItem).PriorityColor).Render
// 	fn := itemStyle.Render
// 	// fn := itemStyleList.Foreground(lipgloss.Color(m.SelectedItem().(TODO).PriorityColor())).Render
// 	if index == m.Index() {
// 		fn = func(s ...string) string {
// 			// return selectedItemStyle.Render("> " + s)
// 			return selectedItemStyle.Render(s[index])
// 		}
// 	}
// 
// 	fmt.Fprintf(w, fn(str))
// }
// type filteredItems []filteredItem
// func (d DefaultDelegate) Render(w io.Writer, m *list.Model, index int, item list.Item) {
// 	var (
// 		title, desc  string
// 		matchedRunes []int
// 		s            = &d.Styles
// 	)
// 	width, height, _ := term.GetSize(0)
// 
// 	if i, ok := item.(list.DefaultItem); ok {
// 		title = i.Title()
// 		desc = i.Description()
// 	} else {
// 		return
// 	}
// 
// 	if width <= 0 {
// 		// short-circuit
// 		return
// 	}
// 
// 	// Prevent text from exceeding list width
// 	textwidth := uint(width - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
// 	title = truncate.StringWithTail(title, textwidth, ellipsis)
// 	if d.ShowDescription {
// 		var lines []string
// 		for i, line := range strings.Split(desc, "\n") {
// 			if i >= height-1 {
// 				break
// 			}
// 			lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
// 		}
// 		desc = strings.Join(lines, "\n")
// 	}
// 
// 	// Conditions
// 	var (
// 		isSelected  = index == m.Index()
// 		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
// 		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
// 	)
// 
// 	if m.IsFiltered && index < len(m.filteredItems) {
// 		// Get indices of matched characters
// 		matchedRunes = m.MatchesForItem(index)
// 	}
// 
// 	if emptyFilter {
// 		title = s.DimmedTitle.Render(title)
// 		desc = s.DimmedDesc.Render(desc)
// 	} else if isSelected && m.FilterState() != Filtering {
// 		if isFiltered {
// 			// Highlight matches
// 			unmatched := s.SelectedTitle.Inline(true)
// 			matched := unmatched.Copy().Inherit(s.FilterMatch)
// 			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
// 		}
// 		title = s.SelectedTitle.Render(title)
// 		desc = s.SelectedDesc.Render(desc)
// 	} else {
// 		if isFiltered {
// 			// Highlight matches
// 			unmatched := s.NormalTitle.Inline(true)
// 			matched := unmatched.Copy().Inherit(s.FilterMatch)
// 			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
// 		}
// 		title = s.NormalTitle.Render(title)
// 		desc = s.NormalDesc.Render(desc)
// 	}
// 
// 	if d.ShowDescription {
// 		fmt.Fprintf(w, "%s\n%s", title, desc)
// 		return
// 	}
// 	fmt.Fprintf(w, "%s", title)
// }
// 
// 





















// buttons
func (m *model) newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	// d := DefaultDelegate{}
	m.ActiveWindow = "addTODO"
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
				if err != nil {
					return errHandler_tui(err, "can't complete item")
				}
				index := ml.Index()
				ml.RemoveItem(index)
				if len(ml.Items()) == 0 {
					keys.choose.SetEnabled(false)
				}
				return ml.NewStatusMessage(statusMessageStyle("You completed " + title+"!"))

			case key.Matches(msg, keys.remove):
				err := m.DelTodo(todoIcal)
				if err != nil {
					return errHandler_tui(err, "can't delete item")
				}
				index := ml.Index()
				ml.RemoveItem(index)
				if len(ml.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return ml.NewStatusMessage(statusMessageStyle("Deleted " + title))

			// case key.Matches(msg, keys.add):
				//TODO IMPORTANT go to add todo form
				// todoInfo := TodoInterface{
				// 	name: "testName1",
				// 	description: "description",
				// 	priority: 3,
				// 	dueTime: time.Now(),
				// 	alarmOffset: "1h",
				// }
				// task, err := CreateTodo(todoInfo)
				// if err != nil {return m.errHandler(err,"Failed to creating todo")}
				// err = m.UploadTodo(task)
				// if err != nil {return m.errHandler(err,"Failed to uploade todo")}
				// err = m.UpdateTodos(task)
				// if err != nil {return m.errHandler(err,"Failed updating todo list")} //Useless

				// m.ActiveWindow = "addTODO"
				// fmt.Println("ADSASDASDASD")
				// return  errHandler_tui(nil,"Failed updating todo list")
				
				// return nil
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove, keys.add}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}




// 	d.Render = func(w io.Writer, m list.Model, index int, listItem list.Item) {
// 
// 		i, ok := listItem.(item)
// 		if !ok {
// 			return
// 		}
// 
// 		str := fmt.Sprintf("%d. %s", index+1, i)
// 
// 		fn := itemStyle.Render
// 		if index == m.Index() {
// 			fn = func(s string) string {
// 				return lipgloss.NewStyle().PaddingLeft(2).Foreground(listItem.PriorityColor)
// 				// .Render("> " + s)
// 			}
// 		}
// 
// 		fmt.Fprintf(w, fn(str))
// 	}
	
	return d
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
	add key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "complete"),
		),
		remove: key.NewBinding(
			key.WithKeys("backspace", "backspace"),
			key.WithHelp("backspace", "delete"),
		),
		add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Add TODO"),
		),
		//TODO PRIORITY FEATURE - edit TODO button
	}
}
