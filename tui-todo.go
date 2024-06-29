package main

import (
	
	"github.com/emersion/go-ical"
	"time"
	"github.com/charmbracelet/bubbles/list"
)


type TODO ical.Event

func (i TODO) Title() string       {
	out,err := i.Props.Get(ical.PropSummary).Text() 
	if err != nil {return "<EMPTY>"}
	return out
}
func (i TODO) Description() string {
	out,err := i.Props.Get(ical.PropDescription).Text() 
	if err != nil {return ""}
	return out
}
func (i TODO) FilterValue() string {
	out1,err1 := i.Props.Get(ical.PropSummary).Text()
	out2,err2 := i.Props.Get(ical.PropDescription).Text()
	if err1 != nil && err2 != nil {return ""}
	return out1+out2
}



func (m *model) GatherTodos() (err error) {
//TODO more modular approach
	calendarObjects, err := GetTODOs(m.Creds.CalendarPath)
	if err != nil {return}

	// var todayTodos []TODO
	
	today := time.Now() 
	todayTodosBuf, err := ParseDueDateTODOs(calendarObjects, today)
	tomorrow := time.Now().AddDate(0, 0, 1)
	tomorrowTodosBuf, err := ParseDueDateTODOs(calendarObjects, tomorrow)




	var todayTodos,tomorrowTodos []TODO
	for _,event := range todayTodosBuf {
		todayTodos = append(todayTodos,TODO(event))
	}
	for _,event := range tomorrowTodosBuf {
		tomorrowTodos = append(tomorrowTodos,TODO(event))
	}
	

	var itemsToday []list.Item
	var itemsTomorrow []list.Item
	for _, todo := range todayTodos {
		itemsToday = append(itemsToday, todo)
	}
	for _, todo := range tomorrowTodos {
		itemsTomorrow = append(itemsTomorrow, todo)
	}

	

	m.TodayTab = list.New(itemsToday, list.NewDefaultDelegate(), 0, 0)
	m.TodayTab.Title = "Today"
	m.TomorrowTab = list.New(itemsTomorrow, list.NewDefaultDelegate(), 0, 0)
	m.TomorrowTab.Title = "Tomorrow"
	
	return nil
}
