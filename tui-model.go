package main

import (
	// "fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/emersion/go-webdav/caldav"
)

type model struct {
	Tabs []string
	// TabContent []string
	LoggedIn     bool
	ActiveWindow string

	TodayTab    list.Model
	TomorrowTab list.Model

	calendarList   list.Model
	calendarChoice string

	todoAddInputs []textinput.Model
	todoAddInputsTime []textinput.Model
	addTimeFocus bool
	loginInputs []textinput.Model
	focused     int
	btnFocus bool
	err         error

	Creds      Credentials
	Calendars  []caldav.Calendar
	CalObjects []caldav.CalendarObject

	errString string
}

type Credentials struct {
	URL          string
	Username     string
	Password     string
	CalendarName string
	CalendarPath string
}

func (m *model) CredentialsSave() (err error) {

	//TODO some proper error handler in case if we cant save
	err = storeCredentialsToKeyring(m.Creds.URL, m.Creds.Username, m.Creds.Password, m.Creds.CalendarPath)
	//TODO add skip flag
	if err != nil {
		return
	}
	return nil
}

func InitModel() model {
	var loginInputs1 []textinput.Model = make([]textinput.Model, 3)
	loginInputs1[url] = textinput.New()
	loginInputs1[url].Placeholder = "https://nextcloud.example/remote.php/dav"
	loginInputs1[url].Focus()
	loginInputs1[url].Width = 30
	loginInputs1[url].Prompt = ""
	// loginInputs1[url].Validate = urlValidator //TODO

	loginInputs1[login] = textinput.New()
	loginInputs1[login].Placeholder = "username"
	loginInputs1[login].Width = 30
	loginInputs1[login].Prompt = ""

	loginInputs1[pass] = textinput.New() 
	loginInputs1[pass].Placeholder = "MySecurePassword"
	loginInputs1[pass].Width = 30
	loginInputs1[pass].Prompt = ""
	loginInputs1[pass].EchoMode = 1 //list.EchoPassword

	var addTODOinputs  []textinput.Model = make([]textinput.Model, 5)
	//TODO NEW FEATURE - add config for default values in addTODO
	//TODO NEW FEATURE - addTODO - add repeat field (and underlying support in caldav.go) 
	//TODO NEW FEATURE - addTODO - multiline input support 
	//TODO NEW FEATURE - addTODO - priority picker
	// //TODO NEW FEATURE - addTODO -  add support for subtasks
	addTODOinputs[name] = textinput.New()
	addTODOinputs[name].Placeholder = "Do task"
	addTODOinputs[name].Focus()
	addTODOinputs[name].Width = 30
	addTODOinputs[name].Prompt = ""

	addTODOinputs[dueDate] = textinput.New()
	addTODOinputs[dueDate].Placeholder = "t/t1/12/12.08/12.08.2025" //TODO how
	// ok, i want to choose from  :
	//	no date
	//	today
	//	tomorrow
	//	custom - pick - https://github.com/EthanEFung/bubble-datepicker
	// + time
	addTODOinputs[dueDate].Width = 23
	addTODOinputs[dueDate].CharLimit = 10
	addTODOinputs[dueDate].Prompt = ""



	addTODOinputs[description] = textinput.New() 
	addTODOinputs[description].Placeholder = "I need to do stuff"
	// addTODOinputs[pass].CharLimit = 3
	addTODOinputs[description].Width = 30
	addTODOinputs[description].Prompt = ""

	addTODOinputs[priority] = textinput.New() 
	addTODOinputs[priority].Placeholder = "0-3 (unknown,low,medium,high)"
	// addTODOinputs[pass].CharLimit = 3
	addTODOinputs[priority].Width = 30
	addTODOinputs[priority].Prompt = ""
	// suggestions := []string{"0","1","2","a3"}
	// addTODOinputs[priority].SetSuggestions(suggestions)

	addTODOinputs[alarmOffset] = textinput.New() 
	addTODOinputs[alarmOffset].Placeholder = "1d/1h/1m"
	// addTODOinputs[pass].CharLimit = 3
	addTODOinputs[alarmOffset].Width = 30
	addTODOinputs[alarmOffset].Prompt = ""


var addTODOinputsTime  []textinput.Model = make([]textinput.Model, 2)
	addTODOinputsTime[dueTimeHour] = textinput.New()
	addTODOinputsTime[dueTimeHour].Placeholder = "14" //TODO how
	addTODOinputsTime[dueTimeHour].Width = 2
	addTODOinputsTime[dueTimeHour].CharLimit = 2
	addTODOinputsTime[dueTimeHour].Prompt = ""

	addTODOinputsTime[dueTimeMinute] = textinput.New()
	addTODOinputsTime[dueTimeMinute].Placeholder = "00" //TODO how
	addTODOinputsTime[dueTimeMinute].Width = 2
	addTODOinputsTime[dueTimeMinute].CharLimit = 2
	addTODOinputsTime[dueTimeMinute].Prompt = ""

	output := model{
		Tabs:        []string{"Today", "Tomorrow", "Add"},
		loginInputs: loginInputs1,
		focused:     0,
		err:         nil,
		todoAddInputs: addTODOinputs,
		todoAddInputsTime: addTODOinputsTime,
		// Creds:	Credentials{"test","test","test","test","test"},
		// TabContent: []string{"ERROR?", "Mascara Tab", "Foundation Tab"},
	}
	return output
}
