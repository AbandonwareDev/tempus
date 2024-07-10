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
//TODO NEW FEATURE IMPORTANT - list today/tomorrow -  edit TODOs
//TODO NEW FEATURE -  add search
//TODO NEW FEATURE -  add custom filter (search that saves with filter)
// TODO NEW FEATURE -  add tabs for days/searcf/filter?




type errMsg struct {message string}


//TODO fix multiple errHandlers, rm m.errHandler, rename CLI errHandler (in main.go). rename errHandler_tui to errHandler. Rename ebery call 
func (m model) errHandler(err error,desc string) (tea.Cmd) {
	if err != nil {
		output := desc+": "+err.Error()
		return func() tea.Msg { return errMsg{output} }
		// return errMsg{output}
		// m.Send(output)
	}
	
	return nil
}
func errHandler_tui(err error,desc string) (tea.Cmd) {
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


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//TODO separate to funcs
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.ActiveWindow {
		case "tomorrow": //TODO rm me  debug
			case "a":
				m.ActiveWindow = "addTODO"
				return m, nil
		case "today": //TODO rm me  debug
			switch keypress := msg.String(); keypress {
				case "q", "ctrl+c":
					// m.quitting = true
					return m, tea.Quit
				case "a":
					m.ActiveWindow = "addTODO"
					return m, nil
		
				// case "y":
					// m."2797322749061742597"
// 				case "h":
// 
// 					todoInfo := TodoInterface{
// 						name: "testName1",
// 						description: "description",
// 						priority: 3,
// 						dueTime: time.Now(),
// 						alarmOffset: "1h",
// 					}
// 					task, err := CreateTodo(todoInfo)
// 					if err != nil {return m, m.errHandler(err,"test fail")}
// 					err = m.UploadTodo(task)
// 					if err != nil {return m, m.errHandler(err,"test fail2")}
// 					m.UpdateTodos(task)
// 					// fmt.Println(&task)
// 					// fmt.Println(err)
// 					// time.Sleep(2*time.Second)
// 					return m, nil
			}
		case "calendarChoose":
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
						if err != nil { return m, m.errHandler(err,"Failed to save credentials")}
						m.CalendarToTodo()
						return m, nil
					}
		
		
			case "login":
				switch keypress := msg.String(); keypress {
				// case "ctrl+c", "q"://are you stupid?
				// 	return m, tea.Quit
				case "enter":
					if m.focused == len(m.loginInputs)-1 {
						//TODO check that we have all fields not empty and notificate about it
						for i := range m.loginInputs {
									m.loginInputs[i], cmd = m.loginInputs[i].Update(msg)
								}
						m.Creds.URL = m.loginInputs[url].Value()
						m.Creds.Username = m.loginInputs[login].Value()
						m.Creds.Password = m.loginInputs[pass].Value()

						err := m.LoginToCalendar()
						if err != nil {return m, m.errHandler(err,"Failed to authenticate")}
						// try login -> choose calendar -> store -> move to getting stuff
						return m, nil
						// return m, tea.Quit
					}
					m.nextLoginInput()
				case "shift+tab", "up":
					// buttonStyle = lipgloss.NewStyle()
					m.prevLoginInput()
				case "tab", "down":
					// buttonStyle = buttonStyle.Background(lipgloss.Color("#7D56F4"))
					m.nextLoginInput()
				}
				for i := range m.loginInputs {
					m.loginInputs[i].Blur()
				}
				m.loginInputs[m.focused].Focus()
	
		
		
			case "addTODO":
				switch keypress := msg.String(); keypress {
				// case "ctrl+c", "q": //are you stupid?
				// 	return m, tea.Quit
				case "enter":
					if m.focused == len(m.todoAddInputs)-1 {
						//TODO check that we have all fields not empty or contains wrong value and notify about it (add textinput validator in tui-model.go)
						
						for i := range m.todoAddInputs {
									m.todoAddInputs[i], cmd = m.todoAddInputs[i].Update(msg)
								}

						err := m.AddTODOtoList()
						if err != nil {return m, m.errHandler(err,"Failed to add TODO")}
						return m, nil
					}
					m.nextTODOInput()
					m.addTimeFocus = false
					// return m, nil
				case "shift+tab", "up":
					// buttonStyle = lipgloss.NewStyle()
					m.addTimeFocus = false
					m.prevTODOInput()
					// return m, nil
				case "right":
				//TODO BUG - cant move left-right, probably just move return in IF statement
				//TODO BUG - focus todoAddInputsTime only if  cursor at last character
					// if focus Due Time
					if m.focused == 1 {
						m.addTimeFocus = true
						m.todoAddInputs[m.focused].Blur()	
						if m.todoAddInputsTime[1].Focused() {
							m.todoAddInputsTime[1].Blur()
							m.todoAddInputsTime[0].Focus()	
						} else {
							m.todoAddInputsTime[0].Blur()
							m.todoAddInputsTime[1].Focus()	
						}						
					}
					return m, nil
				case "left":
				//TODO BUG - cant move left-right, probably just move return in IF statement
				//TODO BUG - focus todoAddInputs only if cursor at zero character 
					if m.focused == 1 {
						m.todoAddInputs[m.focused].Blur()	
						if m.todoAddInputsTime[0].Focused() {
							m.todoAddInputsTime[0].Blur()
							m.todoAddInputsTime[1].Focus()	
						} else {
							m.todoAddInputsTime[0].Blur()
							m.todoAddInputsTime[1].Blur()	
							m.todoAddInputs[m.focused].Focus()
						}						
					}
					return m, nil
				case "tab", "down":
					// buttonStyle = buttonStyle.Background(lipgloss.Color("#7D56F4"))
					m.nextTODOInput()
					m.addTimeFocus = false
					// return m, nil
				}
				// for i := range m.todoAddInputs {
				// 	m.todoAddInputs[i].Blur()
				// }
				// for i := range m.todoAddInputsTime {
				// 	m.todoAddInputsTime[i].Blur()
				// }
				// if !m.addTimeFocus {m.todoAddInputs[m.focused].Focus()}
	
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
		// case "t":
		// 	// TODO add new element
		// 	return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		
		switch m.ActiveWindow {
		//TODO BUG - addTODO/login - output render breaks a bit if errorHandler appeared on screen
		// case "login": //TODO
		// 	m.TodayTab.SetSize(msg.Width-h, msg.Height-v)
		// case "calendarChoose":  //TODO
		// case "addTODO":  //TODO
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
	case "addTODO":
		if (m.focused == 1) && (!m.todoAddInputs[m.focused].Focused()) {
			for i := range m.todoAddInputsTime {
				m.todoAddInputsTime[i], cmd = m.todoAddInputsTime[i].Update(msg)
			}			
		} else {
			for i := range m.todoAddInputs {
				m.todoAddInputs[i], cmd = m.todoAddInputs[i].Update(msg)
			}
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
			// TODO exit on key press
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
	case "addTODO":
		width, height, _ := term.GetSize(0)
		width -= 2
		height -= 2
		loginStyle = loginStyle.
			// Width(30).
			// Height(height/5).
			MarginTop(height / 5).
			MarginLeft(width/2 - 20)
			// MarginRight(width/3)
		tabOutput = loginStyle.Render(m.RenderAddTodo())
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



