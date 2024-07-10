package main

import (
	"fmt"
	// "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"time"
	"errors"
	"strconv"
	"regexp"
	"strings"
	
)

const (
	name = iota
	dueDate
	
	description
	priority
	//TODO add support for subtasks
	alarmOffset
)

const (
	dueTimeMinute = iota
	dueTimeHour
)

func (m model) RenderAddTodo() string {

	return fmt.Sprintf(
		`%s

%s
%s

%s
%s %0-0s%s 

%s
%s

%s
%s

%s
%s


%s
	`,

		inputStyle.Width(30).Align(lipgloss.Center).Render("Create TODO"),
		inputStyle.Width(30).Foreground(lipgloss.AdaptiveColor{Dark: "50"}).Render("Name"),
		m.todoAddInputs[name].View(),
		inputStyle.Width(30).Foreground(lipgloss.AdaptiveColor{Dark: "50"}).Render("Due Time"),
		m.todoAddInputs[dueDate].View(),
		m.todoAddInputsTime[dueTimeHour].View(),
		m.todoAddInputsTime[dueTimeMinute].View(),
		inputStyle.Width(30).Foreground(lipgloss.AdaptiveColor{Dark: "50"}).Render("Description"),
		m.todoAddInputs[description].View(),
		inputStyle.Width(30).Foreground(lipgloss.AdaptiveColor{Dark: "50"}).Render("Priority"),
		//TODO add support for subtasks
		m.todoAddInputs[priority].View(),
		inputStyle.Width(30).Foreground(lipgloss.AdaptiveColor{Dark: "50"}).Render("Alarm offset"),
		m.todoAddInputs[alarmOffset].View(),
		buttonStyle.Render(" Add "),
	)
	// .Align(lipgloss.Center).BorderStyle(lipgloss.NormalBorder())

}

func (m *model) nextTODOInput() {

	m.todoAddInputsTime[0].Blur()
	m.todoAddInputsTime[1].Blur()
	m.todoAddInputs[m.focused].Blur()

	if m.focused == len(m.todoAddInputs)-1 {
		
		if m.btnFocus == true {
			m.btnFocus = false
			buttonStyle = lipgloss.NewStyle()
			m.focused = (m.focused + 1) % len(m.todoAddInputs)
			//TODO do smth with blinking cursor?
		} else {
			buttonStyle = buttonStyle.Background(lipgloss.Color("#7D56F4"))
			m.btnFocus = true
		}
	} else {
		// buttonStyle = lipgloss.NewStyle()
		m.focused = (m.focused + 1) % len(m.todoAddInputs)
	}


	m.todoAddInputs[m.focused].Focus()
}

// prevInput focuses the previous input field
func (m *model) prevTODOInput() {
	m.todoAddInputs[m.focused].Blur()
	if m.focused == len(m.todoAddInputs)-1 {
		if m.btnFocus == true {
			m.btnFocus = false
			buttonStyle = lipgloss.NewStyle()
			//TODO do smth with blinking cursor?
			// m.focused = (m.focused + 1) % len(m.loginInputs)
		} else {
			m.focused--
			// Wrap around
			if m.focused < 0 {
				m.focused = len(m.todoAddInputs) - 1
			}
		}
	} else {
		m.focused--
		// Wrap around
		if m.focused < 0 {
			m.focused = len(m.todoAddInputs) - 1
		}
	}
	m.todoAddInputs[m.focused].Focus()
}

func (m *model) AddTODOtoList() (err error) {

	//TODO BUG - we can get another due date to timezon differs from UTC
	//TODO BUG - need to fix `PRIORITY;VALUE=TEXT:0` - shouldn't be TEXT thing
	now := time.Now()
	// time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	// _, offset := now.Zone() // substract timezone offset because after `now.Truncate(24 * time.Hour)` it makes this -> 2024-07-08 03:00:00 +0300 MSK. so start of a day still in UTC??!? wtf 
	// timezoneOffset := time.Duration(-offset) * time.Second
	// now = now.Add(timezoneOffset) 
	now = time.Date(now.Year(), now.Month(), now.Day(),0,0,0,0, now.Location())
	// now := time.Now().UTC()
	todoInfo := TodoInterface{
		name: m.todoAddInputs[name].Value(),
		description: m.todoAddInputs[description].Value(),
		// priority: priorityTmp,
		// dueTime: time.Now(),
		// alarmOffset: m.todoAddInputs[alarmOffset].Value(),
	}

	
	var dueDateTime time.Time

	m.todoAddInputs[dueDate].Value()


	//t\d
	tdigitRegexp := regexp.MustCompile(`t\d{1,2}`) //tXX
	digitRegexp := regexp.MustCompile(`\d{1,2}`) // XX
	dayMonthRegexp := regexp.MustCompile(`\d{1,2}\.\d{1,2}`) // xx.xx
	fullDateRegexp := regexp.MustCompile(`\d{1,2}\.\d{1,2}\.\d{2,4}`) //xx.xx.xxxx

	//Date 
	date := m.todoAddInputs[dueDate].Value()
	switch {
		// case "":
			//no due date
		case date == "t":
			//if t = today (without time)
			// now := time.Now()
			dueDateTime = now
			// dueDateTime = now.Truncate(24 * time.Hour)
			// dueDateTime = now.Truncate(24 * time.Hour).UTC()
			// return errors.New(fmt.Sprint(offset))
		case tdigitRegexp.MatchString(date):
			//if tX = today + x days (without time) //TODO notify user about this feature somehow
			// now := time.Now()
			addDays,err := strconv.Atoi(strings.TrimPrefix(date,"t"))
			
			dueDateTime = now
			// dueDateTime = now.Truncate(24 * time.Hour)
			if err == nil {dueDateTime = dueDateTime.AddDate(0,0, addDays)
			} else {return errors.New("Can't convert tXX date. Examples: t1, t7, t14")} //TODO return error?

			
		case fullDateRegexp.MatchString(date):
			//if 12.08.2024 /  = you know
			// if 12.08.24  = you know
			sep := strings.Split(date, ".")
			year,err1 := strconv.Atoi(sep[2])
			monthInt,err2 := strconv.Atoi(sep[1])
			day,err3 := strconv.Atoi(sep[0])

			if year <=2000 {year+=2000}
			// year := time.Year(yearInt)
			month := time.Month(monthInt)
			// day := time.Day(dayInt)
			
			//TODO if it's not full year???
			if err1== nil&&err2== nil&&err3 == nil {
					dueDateTime = time.Date(year, month,day, 0, 0, 0, 0, time.Local)	

			}
		case dayMonthRegexp.MatchString(date):
		//if 12.08 = 12th day of 8th month of this year
			currentYear := now.Year()
			sep := strings.Split(date, ".")
				monthInt,err2 := strconv.Atoi(sep[1])
				month := time.Month(monthInt)
				day,err3 := strconv.Atoi(sep[0])
			if err2== nil&&err3 == nil {
			dueDateTime = time.Date(currentYear, month,day, 0, 0, 0, 0, time.Local)
			}
		case digitRegexp.MatchString(date):
		//if 12 = 12th day of this month  of this year
			// now := time.Now()
				currentMonth := now.Month()
				currentYear := now.Year()
				intDate,err := strconv.Atoi(date) 
				if err != nil {_=err} //TODO return error?
			dueDateTime = time.Date(currentYear, currentMonth, intDate, 0, 0, 0, 0, time.Local)
		// case
		// default:
	}
	

	

	

	//Time
	//if hour!="" add to date
	//if monute!="" add to date
	hour,err1 :=  strconv.Atoi(m.todoAddInputsTime[1].Value())
	minute,err2  :=  strconv.Atoi(m.todoAddInputsTime[0].Value())
	if (hour != 0  && err1 ==nil ) || (err2 == nil && minute != 0) {
		if dueDateTime.IsZero() { //set to today
			dueDateTime = now
			// dueDateTime = now.Truncate(24 * time.Hour)
		}
		dueDateTime = dueDateTime.Add(time.Hour * time.Duration(hour) + time.Minute * time.Duration(minute))
	}

	if !dueDateTime.IsZero() {
		todoInfo.dueTime = dueDateTime
		// todoInfo.dueTime = dueDateTime.Add(timezoneOffset)
		// return errors.New(fmt.Sprint(todoInfo.dueTime)) //TODO RM ME DEBUG
	}


	todoInfo.priority = 0
	var priorityTmp int
	if m.todoAddInputs[priority].Value() != "" {
		var err error
		priorityTmp,err =  strconv.Atoi(m.todoAddInputs[priority].Value())
		if err != nil {return err} //TODO debug.  do like in next line 
		// if err != nil {priorityTmp = 0}
	} else {}

	alarmOffsetTmp := m.todoAddInputs[alarmOffset].Value()
	
	

	if priorityTmp != 0 {
		todoInfo.priority = priorityTmp
	}
	if !dueDateTime.IsZero() && alarmOffsetTmp != "" {
		todoInfo.alarmOffset = alarmOffsetTmp
	}



	
	task, err := CreateTodo(todoInfo)
	// if err != nil {return m.errHandler(err,"Failed to creating todo")}
	if err != nil {return err}
	err = m.UploadTodo(task)
	// if err != nil {return m.errHandler(err,"Failed to uploade todo")}
	if err != nil {return err}
	err = m.UpdateTodos(task)
	// if err != nil {return m.errHandler(err,"Failed updating todo list")} //Useless
	if err != nil {return err}
	//TODO deal with errors - append comment on what broken?
	
	m.ActiveWindow = "today" //TODO remove to previus window


	//TODO clean previous input
	for i,_ := range m.todoAddInputsTime {
		m.todoAddInputsTime[i].Blur()
		m.todoAddInputsTime[i].Reset()
	}
	for i,_ := range m.todoAddInputs {
		m.todoAddInputs[i].Blur()
		m.todoAddInputs[i].Reset()
	}
	
	m.focused = 0
	m.todoAddInputs[m.focused].Focus()
	// m.todoAddInputsTime.Reset()
	// m.todoAddInputs.Reset()
	
	
	
	return nil
}
