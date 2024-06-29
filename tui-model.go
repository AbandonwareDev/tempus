package main

import (
	// "fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/emersion/go-webdav/caldav"
	"github.com/charmbracelet/bubbles/textinput"
	
)



type model struct {
	Tabs []string
	// TabContent []string
	LoggedIn     bool
	ActiveWindow string

	TodayTab    list.Model
	TomorrowTab list.Model
	
	calendarList list.Model
	calendarChoice   string
	
	loginInputs []textinput.Model
	focused     int
	err         error

	Creds Credentials 
	Calendars []caldav.Calendar
	CalObjects []caldav.CalendarObject

	errString	string
}

type Credentials struct {
	URL string
	Username string
	Password string
	CalendarName string
	CalendarPath string
}


func (m *model) CredentialsSave() (err error) {

	//TODO some proper error handler in case if we cant save 
	err = storeCredentialsToKeyring(m.Creds.URL,m.Creds.Username,m.Creds.Password,m.Creds.CalendarPath)
	//TODO add skip flag
	if err != nil {return}
	return nil
}


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

	inputs[pass] = textinput.New() //TODO make pass hidden with "github.com/erikgeiser/promptkit/textinput"
	inputs[pass].Placeholder = "MySecurePassword"
	// inputs[pass].CharLimit = 3
	inputs[pass].Width = 30
	inputs[pass].Prompt = ""
	// inputs[pass].Validate = passValidator

	output := model{
		Tabs:        []string{"Today", "Tomorrow", "Add"},
		loginInputs: inputs,
		focused:     0,
		err:         nil,
		// Creds:	Credentials{"test","test","test","test","test"},
		// TabContent: []string{"ERROR?", "Mascara Tab", "Foundation Tab"},
	}
	return output
}
