package main

import (
	"fmt"
	"os"
	"sync"
	"time"
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

	options.InitDAVclients()

	calendars, err := GetCalendars()
	errHandler(err, "Error getting calendars")

	for _,calendar := range calendars {
		fmt.Println(calendar.Name, "-", calendar.Path)
	}

	
	calendarObjects, err := GetTODOs(calendars[1].Path)
	errHandler(err, "Error getting TODOs")	

	today := time.Now()
	todos,err := ParseDueDateTODOs(calendarObjects ,today)
	
	fmt.Println(todos)

	fmt.Println("In total we have",len(calendarObjects), "todos")
}


