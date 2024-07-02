package main

import (
	"fmt"
	"os"
	"net/http"
	"crypto/tls"
	// "sync"
	// "time"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/emersion/go-webdav/caldav"
	// "slices"
	

	// "strconv"
)

// var waitGroup sync.WaitGroup

func errHandler(err error, message string) {
	if err != nil {
		fmt.Printf("\n\n%s: %s\n", message, err)
		os.Exit(1)
	}
}

func main() {
    http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	options, err := ParseOptions()
	errHandler(err, "Error parsing options")
	m := InitModel()

	// debugKeyring()

	// err = storeCredentialsToKeyring("https://pkg.go.dev/encoding/base64#Encoding.EncodeToString", "casual", "h>ÕÇ³PlÇqQ+çCQ{ð±;KÐ¼7¸Âhð~Ümy)v")
	// errHandler(err, "Error parsing options:")

	// if we provide login credentials via cli
	var calendars []caldav.Calendar
	if options.URL != "" {
		err = options.InitDAVclients()
		errHandler(err, "Unexpected error (we couldn't initiate WebDAV/CalDAV client)")
		calendars, err = GetCalendars()
		errHandler(err, "Error getting calendars (incorrect url/login/password)")

		var found bool
		// var calPath string
		if options.Calendar != "" {
			for _,calendar := range calendars { 
				if calendar.Name == options.Calendar {
					found = true
					// calPath = calendar.Path
					m.Creds.CalendarPath = calendar.Path
				}
			}
			if ! found {
				fmt.Println("we don't have calendar ", options.Calendar, ". We have:")
				for _,calendar := range calendars { 
					fmt.Println(calendar.Name)
				}
				os.Exit(1)
			}


			m.LoginToCalendar()
			m.GatherTodos()


// 			calendarObjects, err := GetTODOs(calPath)
// 			errHandler(err, "Error getting TODOs")
// 
// 			today := time.Now() //TODO move to tui and remove it
// 			todayTodos, err := ParseDueDateTODOs(calendarObjects, today)
// 			tomorrow := time.Now().AddDate(0, 0, 1)
// 			tomorrowTodos, err := ParseDueDateTODOs(calendarObjects, tomorrow)
// //TODO remove it
// 			fmt.Println("In total we have", len(calendarObjects), "todos")
// //TODO remove it
// 			var itemsToday []list.Item
// 			var itemsTomorrow []list.Item
// 			for _, todo := range todayTodos {
// 				itemsToday = append(itemsToday, todo)
// 			}
// 			for _, todo := range tomorrowTodos {
// 				itemsTomorrow = append(itemsTomorrow, todo)
// 			}
//TODO remove it
			m.GatherTodos()

			
			// m.TodayTab = list.New(itemsToday, list.NewDefaultDelegate(), 0, 0)
			// m.TodayTab.Title = "Today"
			// m.TomorrowTab = list.New(itemsTomorrow, list.NewDefaultDelegate(), 0, 0)
			// m.TomorrowTab.Title = "Tomorrow"

			m.LoggedIn = true
			m.ActiveWindow = "today"
			
		} else {
			//TODO go to calendars page
			// m.LoggedIn = true
			m.ActiveWindow = "calendarChoose"

			// items := []list.Item{
			// 		item("Ramen"),
			// 		item("Tomato Soup"),
			// 		item("Hamburgers"),
			// 		item("Cheeseburgers"),
			// 		item("Currywurst"),
			// 		item("Okonomiyaki"),
			// 		item("Pasta"),
			// 		item("Fillet Mignon"),
			// 		item("Caviar"),
			// 		item("Just Wine"),
			// 	}
		}
		
	} else {
	//TODO I'm on a highway to (IfElse) hell!
	//TODO probably need to do more careful debug
		creds, err := getCredentialsFromKeyring()
		m.Creds = creds
		if err != nil {
			//we don't have saved credentials, go to login page
			m.ActiveWindow = "login"
		} else {
			err = m.LoginToCalendar()
			if err != nil {
				//we have problem to auth! Go to relogin
				//TODO say that there is a problem with login
				m.ActiveWindow = "login"
			} else {
				err = m.GatherTodos()
				if err != nil {
					//we have problem to auth! Go to relogin
					//TODO say that there is a problem with login
					m.ActiveWindow = "login"
				} else {
					m.LoggedIn = true
					m.ActiveWindow = "today"
				}
			}
			
		}
		
	}


	//DEBUG stuff
	// task, err := CreateTodo("testName","description",3,time.Now())
	// errHandler(err,"test fail")
	// err = m.UploadTodo(task)
	// errHandler(err,"test fail2")
	//DEBUG stuff

	
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// fmt.Println(m.)


}
