package main

import (
	// "fmt"
	// "os"
	// "time"
	// "io"
	// "strings"

	// "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	// textblink "github.com/charmbracelet/bubbles/textinput"
	// "github.com/erikgeiser/promptkit/textinput"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	// "github.com/emersion/go-ical"
	// "github.com/emersion/go-webdav/caldav"

	// "errors"
)
//TODO add new TODos
//TODO edit TODOs
//TODO add search
//TODO add custom filter (search that saves with filter)
// TODO add tabs for days/searcf/filter






type errMsg struct {message string}

func (m model) errHandler(desc string,err error) (tea.Cmd) {
	if err != nil {
		output := desc+": "+err.Error()
		return func() tea.Msg { return errMsg{output} }
		// return errMsg{output}
		// m.Send(output)
	}
	
	return nil
}














func (m model) Init() tea.Cmd {
	return textinput.Blink
	// return nil
}



//TODO add changing calendar



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//TODO separate to funcs
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.ActiveWindow == "calendarChoose" {
			switch keypress := msg.String(); keypress {
					case "q", "ctrl+c":
						// m.quitting = true
						return m, tea.Quit
			
					case "enter":
						i, ok := m.calendarList.SelectedItem().(item)
						if ok {
							m.calendarChoice = string(i) //TODO remove from model - calendarChoice
						}
						for _,calendar := range m.Calendars {
							if calendar.Name == m.calendarChoice {
								m.Creds.CalendarPath = calendar.Path
							}
						}
						
						err := m.CredentialsSave()
						if err != nil { return m, m.errHandler("Failed to save credentials",err)}
						m.CalendarToTodo()
						return m, nil
					}
		}
		
		if m.ActiveWindow == "login" {
			switch keypress := msg.String(); keypress {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				if m.focused == len(m.loginInputs)-1 {
					//TODO check that we have all fields not empty and notificate about it
					//TODO submit
					for i := range m.loginInputs {
								m.loginInputs[i], cmd = m.loginInputs[i].Update(msg)
							}
					// fmt.Println(m.loginInputs[url].Value())//DEBUG
					m.Creds.URL = m.loginInputs[url].Value()
					m.Creds.Username = m.loginInputs[login].Value()
					m.Creds.Password = m.loginInputs[pass].Value()
					// fmt.Println(m.Creds.URL)//DEBUG
						
					// time.Sleep(1 * time.Second) ///DEBUG
					// m.Creds.
					// return m, nil
					err := m.LoginToCalendar()
					if err != nil {return m, m.errHandler("Failed to authenticate",err)}
					// try login -> choose calendar -> store -> move to getting stuff
					return m, nil
					// return m, tea.Quit
				}
				m.nextInput()
			case "shift+tab", "up":
				m.prevInput()
			case "tab", "down":
				m.nextInput()
			}
			for i := range m.loginInputs {
				m.loginInputs[i].Blur()
			}
			m.loginInputs[m.focused].Focus()

		}

		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		// case "right", "l", "n", "tab":
		//TODO add to help
		case "tab":
			if m.LoggedIn {
				h, v := docStyle.GetFrameSize()
				width, height, _ := term.GetSize(0)
				switch m.ActiveWindow {
				case "today":
					m.ActiveWindow = "tomorrow"
					m.TomorrowTab.SetSize(width-h, height-v)
				case "tomorrow":
					m.ActiveWindow = "today"
					m.TodayTab.SetSize(width-h, height-v)
				}
				return m, nil
			}

		// case "a" {
		// 	TODO add new element
		// 	return m, tea.Quit
		// }
		case "t":
			// TODO add new element
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		
		switch m.ActiveWindow {
		// case "login": //TODO
		// 	m.TodayTab.SetSize(msg.Width-h, msg.Height-v)
		// case "calendarChoose":  //TODO
		// 	m.calendarList.SetWidth(msg.Width)
		case "today":
			h, v := docStyle.GetFrameSize()
			m.TodayTab.SetSize(msg.Width-h, msg.Height-v)
		case "tomorrow":
			h, v := docStyle.GetFrameSize()
			m.TomorrowTab.SetSize(msg.Width-h, msg.Height-v)
		}


	case errMsg:
		m.errString = msg.message
		m.ActiveWindow = ""
	}

	
	// text input
	switch m.ActiveWindow {
	case "login":
		for i := range m.loginInputs {
			m.loginInputs[i], cmd = m.loginInputs[i].Update(msg)
		}
	case "calendarChoose":
			m.calendarList, cmd = m.calendarList.Update(msg)
	case "today":
		m.TodayTab, cmd = m.TodayTab.Update(msg)
	case "tomorrow":
		m.TomorrowTab, cmd = m.TomorrowTab.Update(msg)
	
	case "": //exit on key press
		switch tmp := msg.(type) { //TODO debug and optimize
			case tea.KeyMsg:
			_ = tmp
			if m.LoggedIn {
				m.ActiveWindow = "today"
			} else {m.ActiveWindow = "login"}
		}
	}
	// m.TodayTab, cmd = m.TodayTab.Update(msg)
	// m.TomorrowTab, cmd = m.TomorrowTab.Update(msg)

	return m, cmd
}

func (m model) View() string {

	var tabOutput string

	switch m.ActiveWindow {
	case "login":
		width, height, _ := term.GetSize(0)
		width -= 2
		height -= 2
		loginStyle = loginStyle.
			// Width(30).
			// Height(height/5).
			MarginTop(height / 5).
			MarginLeft(width/2 - 20)
			// MarginRight(width/3)
		tabOutput = loginStyle.Render(m.RenderLogin())
		// w, h := lipgloss.Size(tabOutput)
	case "calendarChoose": 
		width, height, _ := term.GetSize(0)
		width -= 2
		height -= 2
		loginStyle = loginStyle.
			// Width(30).
			// Height(height/5).
			MarginTop(height / 5).
			MarginLeft(width/2 - 20)
			// MarginRight(width/3)
		tabOutput = loginStyle.Render(m.RenderCalendarChooser())
		// w, h := lipgloss.Size(tabOutput)

	case "today":
		tabOutput = docStyle.Render(m.TodayTab.View())
	case "tomorrow":
		tabOutput = docStyle.Render(m.TomorrowTab.View())
	case "":
		width, height, _ := term.GetSize(0)
		width -= 2
		height -= 2
		loginStyle = loginStyle.
			Width(width / 3).
			Height(1).
			MarginTop(height / 2).
			MarginLeft(width/3 + 2).
			MarginRight(width / 3)
		tabOutput = loginStyle.Render("ERROR - " + m.errString)
	}
	// if m.activeTab == 0 {
	// 	tabOutput = docStyle.Render(m.TodayTab.View())
	// }
	// if m.activeTab == 1 {
	// 	tabOutput = docStyle.Render(m.TomorrowTab.View())
	// }

	return tabOutput
}



