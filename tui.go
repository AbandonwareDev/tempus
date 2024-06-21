package main


import (
	"fmt"
	// "os"
	// "strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	// textblink "github.com/charmbracelet/bubbles/textinput"
	// "github.com/erikgeiser/promptkit/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var loginStyle = lipgloss.NewStyle().Width(40).Align(lipgloss.Center).BorderStyle(lipgloss.NormalBorder())

var inputStyle = lipgloss.NewStyle()

func (i TODO) Title() string       { return i.Name }
func (i TODO) Description() string { return i.Desc }
func (i TODO) FilterValue() string { return i.Name }

type model struct {
	Tabs       []string
	// TabContent []string
	LoggedIn bool
	ActiveWindow  string
	
	TodayTab list.Model
	TomorrowTab list.Model

	loginInputs  []textinput.Model
	focused int
	err     error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
	// return nil
}

const (
	url = iota
	login
	pass
)
//TODO add changing calendar 

func InitModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[url] = textinput.New()
	inputs[url].Placeholder = "https://nextcloud.example/remote.php/dav"
	inputs[url].Focus()
	// inputs[url].CharLimit = 20
	inputs[url].Width = 30
	inputs[url].Prompt = ""
	// inputs[url].Validate = urlValidator

	inputs[login] = textinput.New()
	inputs[login].Placeholder = "username"
	// inputs[login].CharLimit = 5
	inputs[login].Width = 30
	inputs[login].Prompt = ""
	// inputs[login].Validate = loginValidator

	inputs[pass] = textinput.New()
	inputs[pass].Placeholder = "MySecurePassword"
	// inputs[pass].CharLimit = 3
	inputs[pass].Width = 30
	inputs[pass].Prompt = ""
	// inputs[pass].Validate = passValidator

	output := model{
		Tabs: []string{"Today", "Tomorrow", "Add"},
		loginInputs: inputs,
		focused: 0,
		err:     nil,
		// TabContent: []string{"ERROR?", "Mascara Tab", "Foundation Tab"},
	}
	return output
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.ActiveWindow == "login" {
			switch keypress := msg.String(); keypress {
				case "ctrl+c", "q":
					return m, tea.Quit
				case "enter":
					if m.focused == len(m.loginInputs)-1 {
						//TODO submit

						// try login -> store -> move to getting stuff
						return m, tea.Quit
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
		h, v := docStyle.GetFrameSize()
		switch m.ActiveWindow {
			case "login":
				m.TodayTab.SetSize(msg.Width-h, msg.Height-v)
			case "today":
				m.TodayTab.SetSize(msg.Width-h, msg.Height-v)
			case "tomorrow":
				m.TomorrowTab.SetSize(msg.Width-h, msg.Height-v)
		}

	}


	var cmd tea.Cmd
	// text input
	switch m.ActiveWindow {
		case "login":
			for i := range m.loginInputs {
				m.loginInputs[i], cmd = m.loginInputs[i].Update(msg)
			}
		case "today":
			m.TodayTab, cmd = m.TodayTab.Update(msg)
		case "tomorrow":
			m.TomorrowTab, cmd = m.TomorrowTab.Update(msg)
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
					width -=2
					height -=2
					loginStyle = loginStyle.    
					    // Width(30).
					    // Height(height/5).
   						MarginTop(height/5).
   						MarginLeft(width/2-20)
   						// MarginRight(width/3)
					tabOutput = loginStyle.Render(m.RenderLogin())
					// w, h := lipgloss.Size(tabOutput)

					
				case "today":
					tabOutput = docStyle.Render(m.TodayTab.View())
				case "tomorrow":
					tabOutput = docStyle.Render(m.TomorrowTab.View())
				case "":
					width, height, _ := term.GetSize(0)
					width -=2
					height -=2
					loginStyle = loginStyle.    
					    Width(width/3).
					    Height(1).
   						MarginTop(height/2).
   						MarginLeft(width/3+2).
   						MarginRight(width/3)
					tabOutput = loginStyle.Render("ERROR")
			}
	// if m.activeTab == 0 {
	// 	tabOutput = docStyle.Render(m.TodayTab.View())
	// }
	// if m.activeTab == 1 {
	// 	tabOutput = docStyle.Render(m.TomorrowTab.View())
	// }
	
	return tabOutput
}

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
