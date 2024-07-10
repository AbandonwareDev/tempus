package main

import (
	// "errors"
	"github.com/charmbracelet/bubbles/list"
	"github.com/emersion/go-ical"
	"strings"
	"time"
	"sort"
	"fmt"
	"strconv"
	"github.com/charmbracelet/lipgloss"
)

type TODO ical.Event

func (i TODO) Title() string {

		
		priority := "? "
		prioInt, err := i.Props.Get("PRIORITY").Int()
		var prioStr string
		if err != nil {
			prioStr, err = i.Props.Get("PRIORITY").Text()
			if err == nil {
			prioInt, err = strconv.Atoi(prioStr)
			}
		}
		
		if err == nil {
			switch prioInt{
				case 9:
					priority = "" //Low - Blue
				case 5:
					priority = "❕ " //Mid - Yellow
				case 1:
					priority = "❗ " //High - Red
			}
		} 

		time1, err := i.Props.DateTime("DUE",nil)
		// today := time.Now().Truncate(24 * time.Hour)
		var dueTime,additionalZero string
		if (time1.Hour() != 0) || (time1.Minute() != 0) {
			
			if time1.Minute() <=9 {
				additionalZero = "0"
			}
			dueTime = fmt.Sprint(time1.Hour())+":"+additionalZero+fmt.Sprint(time1.Minute()) + " " //TODO move to sprintf
			// if time1.Minute() == 0 {
			// 	dueTime += "0 "
			// } else {dueTime += " "}
		}

// 		timer, err := i.Props.Get("DUE").Int() //TODO
			// currentDUE = currentDUE.Add(time.Hour * time.Duration(offset))
		//TODO if both times exists
		// if
			// timer+"->"+dueTime +" |"

	out, err := i.Props.Get(ical.PropSummary).Text()
	if err != nil {
		return "<EMPTY>"
	}
	return priority + dueTime + out
}
func (i TODO) UID() string {
	out, err := i.Props.Get(ical.PropUID).Text()
	if err != nil {
		return "<EMPTY>"
	}
	return out
}

func (i TODO) PriorityColor() lipgloss.Color {
	clr, err := i.Props.Get("PRIORITY").Int() 
	c := lipgloss.Color("255") //White
	if err != nil {
			return c
		}
	switch clr{
		case 9:
			c = lipgloss.Color("27") //Blue
		case 5:
			c = lipgloss.Color("220") //Yellow
		case 1:
			// c = lipgloss.Color("#FF9999") //Red
			c = lipgloss.Color("161") //Red
	}
	return c
}
func (i TODO) Description() string {
	out, err := i.Props.Get(ical.PropDescription).Text()
	if err != nil {
		return ""
	}
	return out
}
func (i TODO) FilterValue() string {
	// index
	out1, err1 := i.Props.Get(ical.PropSummary).Text()
	out2, err2 := i.Props.Get(ical.PropDescription).Text()
	if err1 != nil && err2 != nil {
		return ""
	}
	return out1 + out2
}



// type itemDelegate struct{}
// 
// func (d itemDelegate) Height() int                               { return 1 }
// func (d itemDelegate) Spacing() int                              { return 0 }
// func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
// func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
// 	i, ok := listItem.(item)
// 	if !ok {
// 		return
// 	}
// 
// 	str := fmt.Sprintf("%d. %s", index+1, i)
// 
// 	fn := itemStyle.Render
// 	if index == m.Index() {
// 		fn = func(s string) string {
// 			return selectedItemStyle.Render("> " + s)
// 		}
// 	}
// 
// 	fmt.Fprintf(w, fn(str))
// }



func (m *model) GatherTodos() (err error) {
	//TODO more modular approach
	m.CalObjects, err = GetTODOs(m.Creds.CalendarPath)
	if err != nil {
		return
	}

	calendarObjects := m.CalObjects //TODO rm me
	// var todayTodos []TODO

	today := time.Now()
	todayTodosBuf, err := ParseDueDateTODOs(calendarObjects, today)
	// if err != nil {
	// 		return
	// 	}
	tomorrow := time.Now().AddDate(0, 0, 1)
	tomorrowTodosBuf, err := ParseDueDateTODOs(calendarObjects, tomorrow)
	// if err != nil {
	// 		return
	// 	}

	todayTodosBuf,err = SortTodos_Default(todayTodosBuf)
	// if err != nil {
	// 		return
	// 	}
	tomorrowTodosBuf,err = SortTodos_Default(tomorrowTodosBuf)

	// if err != nil {
	// 		return
	// 	}

	//TODO fmt - we can combine things below
	var todayTodos, tomorrowTodos []TODO
	for _, event := range todayTodosBuf {
		todayTodos = append(todayTodos, TODO(event))
	}
	for _, event := range tomorrowTodosBuf {
		tomorrowTodos = append(tomorrowTodos, TODO(event))
	}

	var itemsToday []list.Item
	var itemsTomorrow []list.Item
	for _, todo := range todayTodos {
		itemsToday = append(itemsToday, todo)
	}
	for _, todo := range tomorrowTodos {
		itemsTomorrow = append(itemsTomorrow, todo)
	}

	delegateKeys := newDelegateKeyMap()
	delegate := m.newItemDelegate(delegateKeys)
	// m.TodayTab = list.New(itemsToday, list.NewDefaultDelegate(), 0, 0)
	m.TodayTab = list.New(itemsToday, delegate, 0, 0)
	m.TodayTab.FilterInput.Reset() //TODO debug
	// m.TodayTab.Title = "Today"
	// 	m.TodayTab.SetShowStatusBar(false)
	// 	m.TodayTab.SetFilteringEnabled(false)
	// 	m.TodayTab.Styles.Title = titleStyle
	// 	m.TodayTab.Styles.PaginationStyle = paginationStyle
	// 	m.TodayTab.Styles.HelpStyle = helpStyle

	// m.TomorrowTab = list.New(itemsTomorrow, list.NewDefaultDelegate(), 0, 0)
	m.TomorrowTab = list.New(itemsTomorrow, delegate, 0, 0)
	m.TodayTab.FilterInput.Reset()//TODO debug
	m.TomorrowTab.Title = "Tomorrow"
	// m.TomorrowTab.SetShowStatusBar(false)
	// 	m.TomorrowTab.SetFilteringEnabled(false)
	// 	m.TomorrowTab.Styles.Title = titleStyle
	// 	m.TomorrowTab.Styles.PaginationStyle = paginationStyle
	// 	m.TomorrowTab.Styles.HelpStyle = helpStyle

	return nil
}


//for sorting
	type ByTime []ical.Event
	func (a ByTime) Len() int           { return len(a) }
	func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
	func (a ByTime) Less(i, j int) bool { 
		time1,_ := a[i].Props.DateTime("DUE",nil)
		//TODO is there better way?
		// if time1.Hour() == 0 || time1.Minute() == 0 {
		// 	time1 = time.Date(time1.Year(), time1.Month(), time1.Day(), 23, 59, 0, 0, time.UTC)
		// }
		// timeInt1,_ := strconv.Atoi(fmt.Sprint("%d%d",time1.Hour(),time1.Minute()))
		s1:=fmt.Sprintf("%d%d",time1.Hour(),time1.Minute())
		time2,_ := a[j].Props.DateTime("DUE",nil)
		// timeInt2,_ := strconv.Atoi(fmt.Sprint("%d%d",time2.Hour(),time2.Minute()))
		s2 := fmt.Sprintf("%d%d",time2.Hour(),time2.Minute())

		timeInt1, _ := strconv.Atoi(s1)
		timeInt2, _ := strconv.Atoi(s2)

		if timeInt1 == 0 {
			timeInt1 = 9999
		}
		if timeInt2 == 0 {
			timeInt2 = 9999
		}
		// if time2.Hour() == 0 || time2.Minute() == 0 {
		// 			time2 = time.Date(time2.Year(), time2.Month(), time2.Day(), 23, 59, 0, 0, time.UTC)
		// 		}
		// return time2.After(time1)
		return timeInt1<timeInt2
	}
	type ByPriority []ical.Event
	func (a ByPriority) Len() int           { return len(a) }
	func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
	func (a ByPriority) Less(i, j int) bool { 
		prior1,err := a[i].Props.Get("PRIORITY").Int()

		// in case we have `PRIORITY;VALUE=TEXT:0`
		if err != nil {
			prioStr, err := a[i].Props.Get("PRIORITY").Text()
			if err == nil {
			prior1, err = strconv.Atoi(prioStr)
			}
		}
		

		//if no priority
		if prior1 == 0 {
			prior1 = 11
		}	
			
		prior2,_ := a[j].Props.Get("PRIORITY").Int()
		if prior2 == 0 {
			prior2 = 11
		}
		
		return prior1 < prior2
	}
func SortTodos_Default(inputEvents []ical.Event)  ([]ical.Event, error) {
	//TODO sort other alphabetically 
	var timeEvents,priorityEvents,output []ical.Event
	
	for _,event := range inputEvents {
		time,err := event.Props.DateTime("DUE",nil)
		if (err == nil) && (time.Hour() != 0 || time.Minute() != 0) {
			timeEvents = append(timeEvents,event)
		} else {
			priorityEvents = append(priorityEvents, event)
		}
	}
	
	sort.Sort(ByPriority(priorityEvents))
	sort.Sort(ByTime(timeEvents)) 

	output = timeEvents
	output = append(output,priorityEvents...)

	return output,nil
}



func (m *model) UpdateTodos(todo ical.Event) (err error) {
	today := time.Now()
	tomorrow := time.Now().AddDate(0, 0, 1)
	errorI := 0
	if todo.Props["DUE"] != nil {
	if strings.HasPrefix(todo.Props["DUE"][0].Value, today.Format("20060102")) {
		m.TodayTab.InsertItem(-1, TODO(todo))
		//TODO makey another type of "TodayTab"?, like []slice of tabs
		//TODO sort after update
	} else {
		errorI += 1
	}
	
		if strings.HasPrefix(todo.Props["DUE"][0].Value, tomorrow.Format("20060102")) {
			m.TomorrowTab.InsertItem(-1, TODO(todo))
			//TODO sort after update
		} else {
			errorI += 1
		}
	} else {
				errorI += 1
			}

	if errorI == 2 {
		// return errors.New("don't match today and tomorrow") // debug???
		return nil
	}
	return nil	
}
