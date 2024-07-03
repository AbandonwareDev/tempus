package main

import (
	
	"github.com/emersion/go-ical"
	"time"
	"github.com/charmbracelet/bubbles/list"
	"strings"
	"errors"
)


type TODO ical.Event

func (i TODO) Title() string       {
	out,err := i.Props.Get(ical.PropSummary).Text() 
	if err != nil {return "<EMPTY>"}
	return out
}
func (i TODO) UID() string       {
	out,err := i.Props.Get(ical.PropUID).Text() 
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
	m.CalObjects, err = GetTODOs(m.Creds.CalendarPath)
	if err != nil {return}

	calendarObjects := m.CalObjects //TODO rm me
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



	delegateKeys  := newDelegateKeyMap()
	delegate := m.newItemDelegate(delegateKeys)
	// m.TodayTab = list.New(itemsToday, list.NewDefaultDelegate(), 0, 0)
	m.TodayTab = list.New(itemsToday, delegate, 0, 0)
	m.TodayTab.Title = "Today"
	// m.TomorrowTab = list.New(itemsTomorrow, list.NewDefaultDelegate(), 0, 0)
	m.TomorrowTab = list.New(itemsTomorrow, delegate, 0, 0)
	m.TomorrowTab.Title = "Tomorrow"
	
	return nil
}


func (m *model) UpdateTodos(todo ical.Event) (err error) {
	today := time.Now() 
	tomorrow := time.Now().AddDate(0, 0, 1)
	errorI := 0
	if strings.HasPrefix(todo.Props["DUE"][0].Value, today.Format("20060102")) {m.TodayTab.InsertItem(-1,TODO(todo))} else {errorI += 1}
	if strings.HasPrefix(todo.Props["DUE"][0].Value, tomorrow.Format("20060102")) {m.TomorrowTab.InsertItem(-1,TODO(todo)) } else {errorI += 1}
	
	if errorI == 2 {return errors.New("don't match today and tomorrow")}
	
	
	

	return nil
}
