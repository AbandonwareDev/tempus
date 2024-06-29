package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)


var loginStyle = lipgloss.NewStyle().Width(40).Align(lipgloss.Center).BorderStyle(lipgloss.NormalBorder())

var inputStyle = lipgloss.NewStyle()

const (
	url = iota
	login
	pass
)


func (m model) RenderLogin() string {

	return fmt.Sprintf(
		`%s

	
%s
%s

%s
%s

%s
%s


%s
	`,

		inputStyle.Width(30).Align(lipgloss.Center).Render("Login"),
		inputStyle.Width(30).Foreground(lipgloss.AdaptiveColor{Dark: "50"}).Render("WebDAV server URL"),
		m.loginInputs[url].View(),
		inputStyle.Width(30).Foreground(lipgloss.AdaptiveColor{Dark: "50"}).Render("Login"),
		m.loginInputs[login].View(),
		inputStyle.Width(30).Foreground(lipgloss.AdaptiveColor{Dark: "50"}).Render("Password"),
		m.loginInputs[pass].View(), //TODO hide
		inputStyle.Render("Continue ->"),
	)
	// .Align(lipgloss.Center).BorderStyle(lipgloss.NormalBorder())

}

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.loginInputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.loginInputs) - 1
	}
}


func (m *model) LoginToCalendar() (err error) {
	err = m.InitDAVclients()
	if err != nil { return }
	
	m.Calendars, err = GetCalendars()
	if err != nil { return }

	items := []list.Item{}

	for _,calendar := range m.Calendars { 
		// fmt.Println(calendar.Name)
		items = append(items,item(calendar.Name))
	}

	// m.LoggedIn = true TODO after calendar choose

	const defaultWidth = 20	
	const listHeight = 24
		l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
		l.Title = "Which To-Do list do we need?"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle

	m.calendarList = l
	m.ActiveWindow = "calendarChoose"

	
	return nil
}
