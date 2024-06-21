package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"

	// "strconv"
)

var waitGroup sync.WaitGroup

func errHandler(err error, message string) {
	if err != nil {
		fmt.Printf("\n\n%s: %s\n", message, err)
		os.Exit(1)
	}
}

func main() {
	options, err := ParseOptions()
	errHandler(err, "Error parsing options")



	// debugKeyring()

	err = storeCredentialsToKeyring("https://pkg.go.dev/encoding/base64#Encoding.EncodeToString","casual","h>ÕÇ³PlÇqQ+çCQ{ð±;KÐ¼7¸Âhð~Ümy)v")
	errHandler(err, "Error parsing options")

// 	url1,login1,password1,err := getCredentialsFromKeyring()
// 	errHandler(err, "Error parsing options")
// m.loginInputs[url].Value()
// m.loginInputs[login].Value()
// m.loginInputs[pass].Value()

	// if true {os.Exit(1)}
















 
	options.InitDAVclients()

	calendars, err := GetCalendars()
	errHandler(err, "Error getting calendars")

	for _, calendar := range calendars {
		fmt.Println(calendar.Name, "-", calendar.Path)
	}

	calendarObjects, err := GetTODOs(calendars[1].Path)
	errHandler(err, "Error getting TODOs")

	today := time.Now()
	todayTodos, err := ParseDueDateTODOs(calendarObjects, today)
	tomorrow := time.Now().AddDate(0,0,1)
	tomorrowTodos, err := ParseDueDateTODOs(calendarObjects, tomorrow)

	// fmt.Println(todos)

	fmt.Println("In total we have", len(calendarObjects), "todos")

	var itemsToday []list.Item
	var itemsTomorrow []list.Item

	for _,todo := range todayTodos {
		itemsToday = append(itemsToday,todo)
	}

	for _,todo := range tomorrowTodos {
		itemsTomorrow = append(itemsTomorrow,todo)
	}

	m := InitModel()
	m.TodayTab = list.New(itemsToday, list.NewDefaultDelegate(), 0, 0)
	m.TodayTab.Title = "Today"
	m.TomorrowTab = list.New(itemsTomorrow, list.NewDefaultDelegate(), 0, 0)
	m.TomorrowTab.Title = "Tomorrow"
	// m.LoggedIn = true
	m.ActiveWindow = "login"
	
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	fmt.Println(m.loginInputs[login].Value())
}
