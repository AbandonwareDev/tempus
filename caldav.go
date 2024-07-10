package main

import (
	"context"
	// "fmt"
	"github.com/emersion/go-ical"
	webdav "github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"github.com/teambition/rrule-go"

	"errors"
	"github.com/google/uuid"
	"strings"
	"time"

	"net/http"
	// "net/url"
	// "fmt"
)

// type TODO struct {
// 	Name string
// 	Desc string
// 	Time string //only parsed _time_
// 	// Priority int //TODO
// 	// Subtasks []TODO //TODO
// 	// Repeat //TODO
// 	// Alarm //TODO
// 	UID string //for editing events
// 	DueDateTime time.Time //for adding new events
// }

// TODO Is it safe to create global variables?
var clientWebDAV *webdav.Client
var client *caldav.Client // clientCalDAV
// var calendarObjects []caldav.CalendarObject
var ctx = context.Background()

// var authSession caldav.Client // clientCalDAV

func InitDAVclients(url, user, pass string) error {

	var err error
	authSession := webdav.HTTPClientWithBasicAuth(nil, user, pass)

	clientWebDAV, err = webdav.NewClient(authSession, url)
	if err != nil {
		// Handle error
		return err
	}
	client, err = caldav.NewClient(authSession, url)
	if err != nil {
		// Handle error
		return err
	}
	return nil
}

func (options *Options) InitDAVclients() error {

	err := InitDAVclients(options.URL, options.User, options.Password)
	if err != nil {
		// Handle error
		return err
	}
	return nil
}
func (m model) InitDAVclients() error {

	err := InitDAVclients(m.Creds.URL, m.Creds.Username, m.Creds.Password)
	if err != nil {
		// Handle error
		return err
	}
	return nil
}

func GetCalendars() ([]caldav.Calendar, error) {

	principal, err := clientWebDAV.FindCurrentUserPrincipal(ctx)
	if err != nil {
		// Handle error
		return nil, err
	}
	// fmt.Println("principal: ",principal) TODO log

	// Find the calendar home set
	calendarHomeSet, err := client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		return nil, err
	}

	// Find calendars in the calendar home set
	calendars, err := client.FindCalendars(ctx, calendarHomeSet)
	if err != nil {
		return nil, err
	}

	return calendars, nil
}

func GetTODOs(calendarPath string) (calendarObjects []caldav.CalendarObject, err error) {
	//TODO can we use calendar.Data.Component.Children to dont make this close to pointless request?
	date := time.Now()
	dateStart := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location()).AddDate(0, 0, -1) //TODO too complex - time.Now().Add(-92 * time.Hour),
	// dateEnd:= time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location()) //date +1 day
	dateEnd := date.AddDate(0, 1, 0) //TODO we hard limit pulling events only for this month, to prevent getting tooo much events

	calQuery := caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Comps: []caldav.CalendarCompRequest{{
				Name: "VTODO",
				Props: []string{
					"UID",
					"SUMMARY",
					"COMPLETED",
					"DESCRIPTION",
					"DUE",
					"PRIORITY",
					"RELATED-TO",
				},
			},
			},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{
				{
					Name:  "VTODO",
					Start: dateStart,
					End:   dateEnd,
					// FYI: server-side filtering is not supported for Nextcloud
				},
			},
		},
	}
	calendarObjects, err = client.QueryCalendar(ctx, calendarPath, &calQuery)
	if err != nil {
		return nil, err
	}
	return calendarObjects, nil
}

// func ParseDueDateTODOs_depricated(calObjs []caldav.CalendarObject, date time.Time) ([]TODO, error) {
// 	var output []TODO
//
// 	for _, calObj := range calObjs {
// 		// fmt.Println((*(*calObj.Data).Children[0]).Name)
// 		// TODO STATUS map[] COMPLETED
// 		for _, event := range (*calObj.Data).Children {
// 			// if (*event).Name == "VTODO" {
// 			// var notCompletedTODO, withDate, fromToday bool
// 			var notCompletedTODO, fromToday bool
// 			//TODO we can optimize there if we encounter wrong state to forcefully stop next analysis
// 			//TODO we can use event.Props.Get
// 			// notCompletedTODO
// 			if (*event).Props["COMPLETED"] == nil {
// 				if (*event).Props["STATUS"] == nil {
// 					notCompletedTODO = true
// 				} else {
// 					if (*event).Props["STATUS"][0].Value != "COMPLETED" {
// 						notCompletedTODO = true
// 					}
// 				}
// 			}
//
// 			// withTodayDate
// 			if (*event).Props["DUE"] != nil {
// 				// withDate = true
// 				// fromToday
// 				if strings.HasPrefix((*event).Props["DUE"][0].Value, date.Format("20060102")) {
// 					fromToday = true
// 				}
// 			}
//
// 			// if notCompletedTODO && withDate && fromToday {
// 			if notCompletedTODO  && fromToday {
//
// 				var tmpTODO TODO
//
// 				if (*event).Props["SUMMARY"] != nil {
// 					name := (*event).Props["SUMMARY"][0].Value
// 					tmpTODO.Name = name
// 				} else {
// 					tmpTODO.Name = "<EMPTY>"
// 				}
// 				if (*event).Props["DESCRIPTION"] != nil {
// 					description := (*event).Props["DESCRIPTION"][0].Value
// 					tmpTODO.Desc = description
// 				}
//
// 				var todoTime string
// 				due := (*event).Props["DUE"][0].Value
// 				index := strings.Index(due, "T")
// 				if index != -1 {
// 					str := due[index+1:]
// 					todoTime = str[:2] + ":" + str[2:4]
//
// 					tmpTODO.Time = todoTime
// 				}
// 				output = append(output, tmpTODO)
// 			}
// 			// }
// 		}
// 	}
//
// 	// //TODO 		sort: time, priority (if no time)
// 	// 			// it means, put DUE at first, other DUE;VALUE=DATE in priority order
// 	// //TODO 		color: priority, add time icon if have time
// 	// //TODO 		UID:7ed30f40-fce1-422c-be3b-0486dcfe8943
// 	// //TODO 		RELATED-TO:7ed30f40-fce1-422c-be3b-0486dcfe8943 # subtask
// 	// //TODO 		PRIORITY:1 #1-high, 5-mid, 9-low
// 	//TODO repeat function???
//
// 	//TODO on complete -repeat function
// 	// RRULE:FREQ=WEEKLY;INTERVAL=1
//
// 	//TODO if no repeat - mark as complted
// 	// STATUS:COMPLETED
// 	// COMPLETED:20240421T065323Z
// 	// PERCENT-COMPLETE:100
//
// 	//TODO support notifcations/alarms???
// 	// BEGIN:VTODO
// 	// ...
// 	// BEGIN:VALARM
// 	// TRIGGER;RELATED=END:-PT15M
// 	// ACTION:DISPLAY
// 	// DESCRIPTION:Default Tasks.org description
// 	// END:VALARM
// 	// END:VTODO
//
// 	return output, nil
// }

func ParseDueDateTODOs(calObjs []caldav.CalendarObject, date time.Time) ([]ical.Event, error) {
	var output []ical.Event

	for _, calObj := range calObjs {

		// fmt.Println((*(*calObj.Data).Children[0]).Name)
		// TODO STATUS map[] COMPLETED
		// for _, event := range (*calObj.Data).Children {
		// event
		// for _, event1 := range (*calObj.Data).Children {
		// // fmt.Println("1:", event1)
		// tmp := ical.Event{Component: event1}
		// fmt.Println(tmp.Props.Get(ical.PropUID).Text())
		// // for _, event2 := range tmp.Events() {
		// // // for _, event2 := range event1.Events() {
		// // fmt.Println("1:", event2)
		// // }
		// }

		for _, eventComponent := range (*calObj.Data).Children {
			event := ical.Event{Component: eventComponent}
			// for _, event := range (*calObj.Data).Children {
			// fmt.Println("1:", event) //TODO rm me
			// if (*event).Name == "VTODO" {
			// var notCompletedTODO, withDate, fromToday bool
			var notCompletedTODO, fromToday bool
			//TODO we can optimize there if we encounter wrong state to forcefully stop next analysis
			//TODO we can use event.Props.Get
			// notCompletedTODO
			if event.Props["COMPLETED"] == nil {
				if event.Props["STATUS"] == nil {
					notCompletedTODO = true
				} else {
					if event.Props["STATUS"][0].Value != "COMPLETED" {
						notCompletedTODO = true
					}
				}
			}

			// withTodayDate
			if event.Props["DUE"] != nil {
				// withDate = true
				// fromToday
				if strings.HasPrefix(event.Props["DUE"][0].Value, date.Format("20060102")) {
					fromToday = true
				}
			}

			// if notCompletedTODO && withDate && fromToday {
			if notCompletedTODO && fromToday {

				tmpTODO := event

				if event.Props["SUMMARY"] == nil {
					tmpTODO.Props.SetText(ical.PropSummary, "<EMPTY>")
				}
				if event.Props["DESCRIPTION"] == nil {
					tmpTODO.Props.SetText(ical.PropDescription, "")
				}
				if event.Props["PRIORITY"] == nil {
					tmpTODO.Props.SetText("PRIORITY", "0")
				}

				// 				var todoTime string
				// 				due := (*event).Props["DUE"][0].Value
				// 				index := strings.Index(due, "T")
				// 				if index != -1 {
				// 					str := due[index+1:]
				// 					todoTime = str[:2] + ":" + str[2:4]
				//
				// 					tmpTODO.Time = todoTime
				// 				}

				output = append(output, tmpTODO)
			}
			// }
		}
	}

	// //TODO 		sort: time, priority (if no time)
	// 			// it means, put DUE at first, other DUE;VALUE=DATE in priority order
	// //TODO 		color: priority, add time icon if have time
	// //TODO 		UID:7ed30f40-fce1-422c-be3b-0486dcfe8943
	// //TODO 		RELATED-TO:7ed30f40-fce1-422c-be3b-0486dcfe8943 # subtask
	// //TODO 		PRIORITY:1 #1-high, 5-mid, 9-low
	//TODO on complete -repeat function
	//TODO repeate - RRULE:FREQ=DAILY;INTERVAL=1
	// RRULE:FREQ=WEEKLY;INTERVAL=1

	//TODO if no repeat - mark as complted
	// STATUS:COMPLETED
	// COMPLETED:20240421T065323Z
	// PERCENT-COMPLETE:100

	//TODO support notifcations/alarms???
	// BEGIN:VTODO
	// ...
	// BEGIN:VALARM
	// TRIGGER;RELATED=END:-PT15M
	// ACTION:DISPLAY
	// DESCRIPTION:Default Tasks.org description
	// END:VALARM
	// END:VTODO
	// fmt.Println(output) //TODO rm me
	// time.Sleep(3*time.Second)//TODO rm me

	return output, nil
}

type TodoInterface struct {
	name        string
	description string
	priority    int
	dueTime     time.Time
	alarmOffset string
	repeat      string
	//TODO repeat
	//TODO subtasks
}

func CreateTodo(info TodoInterface) (event ical.Event, err error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return
	}
	event = *ical.NewEvent()
	event.Name = ical.CompToDo //VTODO
	event.Props.SetText(ical.PropUID, uid.String())
	event.Props.SetText(ical.PropSummary, info.name)
	event.Props.SetText(ical.PropDescription, info.description)
	if !info.dueTime.IsZero() {
		hour := info.dueTime.Hour()
		minute := info.dueTime.Minute()
		if (hour != 0 || minute != 0 ) {
			event.Props.SetDateTime(ical.PropDue, info.dueTime)
		} else {
			event.Props.SetDate(ical.PropDue, info.dueTime)
		}
		// }
	} // 'zero' time is `time.Time{}`

	switch info.priority {
	case 0: //No priority
		event.Props.SetText(ical.PropPriority, "0")
	case 1: //Light
		event.Props.SetText(ical.PropPriority, "9")
	case 2: //Medium
		event.Props.SetText(ical.PropPriority, "5")
	case 3: //Urgent
		event.Props.SetText(ical.PropPriority, "1")
	default:
		err = errors.New("Wrong priority, expecting 0-3")
		return
	}
	//TODO repeat - RRULE:FREQ=DAILY;INTERVAL=1
	if info.alarmOffset != "" {

		// alarmComponent :=  ical.Component{Name:ical.CompAlarm}
		alarmComponent := ical.NewComponent(ical.CompAlarm)

		alarmComponent.Props.SetText(ical.PropAction, "DISPLAY")
		// alarmComponent.Props.Add(&ical.Prop{Name:ical.PropAction,Value:"DISPLAY",})
		alarmComponent.Props.Add(&ical.Prop{Name: ical.PropDescription, Value: "Default Alarm Tempus description"})
		alarmComponent.Props.SetText(ical.PropDescription, "Default Alarm Tempus description")
		// ACTION:DISPLAY
		// DESCRIPTION:Default Tasks.org description
		var value string
		if info.alarmOffset == "0" {
			// TRIGGER;RELATED=END:PT0S
			value = "PT0S"

		}

		if strings.HasSuffix(info.alarmOffset, "h") {
			offset, _ := strings.CutSuffix(info.alarmOffset, "h")
			value = "-PT" + offset + "H"
			// TRIGGER;RELATED=END:-PT1H

		}
		//TODO stop next if
		if strings.HasSuffix(info.alarmOffset, "m") {
			offset, _ := strings.CutSuffix(info.alarmOffset, "m")
			value = "-PT" + offset + "M"
			// TRIGGER;RELATED=END:-PT10M
		}

		if strings.HasSuffix(info.alarmOffset, "d") {
			offset, _ := strings.CutSuffix(info.alarmOffset, "d")
			value = "-P" + offset + "D"
			// TRIGGER;RELATED=END:-P1D
		}
		alarmComponent.Props.SetText("TRIGGER;RELATED=END", value)
		event.Children = append(event.Children, alarmComponent)
	}

	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	event.Props.SetDateTime(ical.PropCreated, time.Now().UTC()) //TODO if it don't exist already (in case if we edit todo)
	event.Props.SetDateTime(ical.PropLastModified, time.Now().UTC())
	//TODO add a function to verify event (exist in ical lib)
	return
}

func (m model) UploadTodo(event ical.Event) (err error) {
	// event.Props.SetDateTime(ical.PropDateTimeStart, startDateTime)

	// TODO Alarm component properties
	// PropAction  = "ACTION"
	// PropRepeat  = "REPEAT"
	// PropTrigger = "TRIGGER"}

	// calendar, err := client.GetCalendarObject(ctx, m.Creds.CalendarPath) //makes error on nextcloud
	// if err != nil {return err}

	calendar := ical.NewCalendar()
	calendar.Props.SetText(ical.PropProductID, "+//Casual//Tempus//EN")
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Component.Children = append(calendar.Component.Children, event.Component)

	todoGUID, err := event.Props.Get(ical.PropUID).Text()
	if err != nil {
		return err
	}
	//TODO check GUID uniq and regenerate if needed

	var buf strings.Builder
	encoder := ical.NewEncoder(&buf)
	err = encoder.Encode(calendar)
	if err != nil {
		return err
	}
	_, err = client.PutCalendarObject(ctx, m.Creds.CalendarPath+todoGUID+".isc", calendar)
	if err != nil {
		return err
	}
	return nil

}

// func (m model) DelTodo(delUID string) (err error) {
func (m model) DelTodo(todo ical.Event) (err error) {

	delUID, err := todo.Props.Get(ical.PropUID).Text()
	// if err != nil {return}

	// calendar, err := client.GetCalendarObject(ctx, m.Creds.CalendarPath+delUID+".isc")
	// if err != nil {return}
	//
	// var newEvents []*ical.Component
	// for _, component := range calendar.Data.Component.Children {
	// 	if component.Name == ical.CompEvent {
	// 		var uid string
	// 		uid, err = component.Props.Text(ical.PropUID)
	// 		if err != nil {return}
	// 		if uid != delUID {
	// 			newEvents = append(newEvents, component)
	// 		}
	// 	}
	// }
	//
	// calendar.Data.Component.Children = newEvents
	// var buf strings.Builder
	// encoder := ical.NewEncoder(&buf)
	// err = encoder.Encode(calendar.Data)
	// if err != nil {return}
	//
	// _, err = client.PutCalendarObject(ctx, m.Creds.CalendarPath, calendar.Data)
	// if err != nil {return}

	// req := http.Request{
	// 	Method: "DELETE",
	// 	URL: url.Parse(m.Creds.URL + m.Creds.CalendarPath + delUID + ".isc"),
	//
	// }

	client := &http.Client{}

	parts := strings.Split(m.Creds.URL, "/")
	baseURL := parts[0] + "//" + parts[2]
	// Create request
	req, err := http.NewRequest("DELETE", baseURL+m.Creds.CalendarPath+delUID+".isc", nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(m.Creds.Username, m.Creds.Password)

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	// resp.Body.Close()
	defer resp.Body.Close()

	// Read Response Body
	// respBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	//     // fmt.Println(err)
	//     return
	// }

	if resp.Status == "204 No Content" {
		return nil
	}
	if resp.Status == "404 Not Found" {
		// Nextcloud Tasks create differend GUID for Vtodo and path
		// TODO code fmt - Highway to for-if hell

		for _, calObj := range m.CalObjects {
			for _, eventComponent := range (*calObj.Data).Children {
				event := ical.Event{Component: eventComponent}
				if event.Name == "VTODO" {
					// eventUID,_ := (*event).Props["UID"][0].Text()

					eventUID, err := event.Props.Get(ical.PropUID).Text()

					if err != nil {
						eventUID = ""
					}
					if eventUID == delUID {
						//verify that it's indeed that event by comparing Name
						sum1, err := todo.Props.Get("SUMMARY").Text()
						if err != nil {
							return err
						}
						sum2, err := event.Props.Get("SUMMARY").Text()
						if err != nil {
							return err
						}
						if sum1 == sum2 {
							req2, err := http.NewRequest("DELETE", baseURL+calObj.Path, nil)
							if err != nil {
								return err
							}
							req2.SetBasicAuth(m.Creds.Username, m.Creds.Password)

							// Fetch Request
							resp2, err := client.Do(req2)
							if err != nil {
								return err
							}

							defer resp2.Body.Close()
							if resp2.Status == "204 No Content" {
								return nil
							} else {
								return errors.New("Can't delete, response status: " + resp2.Status + ".")
							}

							//TODO exit for loop
							// return errors.New("test")
						}
					}
				}
			}

		}

		// return errors.New("test")
	}
	time.Sleep(2 * time.Second) //TODO DEBUG RM ME
	// Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))

	// return nil
	return errors.New("Can't delete, response status: " + resp.Status + ".")
}

func (m model) EditTodo(todo ical.Event) (err error) {
	//TODO is there proper edit function ???
	// uid,err := todo.Props.Get(ical.PropUID).Text()
	// if err != nil {return}

	err = m.DelTodo(todo)
	if err != nil {
		return
	}

	err = m.UploadTodo(todo)
	if err != nil {
		//TODO panic, we deleted but cant upload event. Need to save somehow. Or at least write debug string with parameters so we can re-add manually in that case
		return
	}

	return nil
}

// TODO
func (m model) CompleteTodo(todo ical.Event) (err error) {

	//if it repitable - update DUE time, otherwise complete it
	if todo.Props["RRULE"] != nil {
		var rOptions *rrule.ROption
		rOptions, err = todo.Props.RecurrenceRule()
		if err != nil {
			return
		}
		offset := rOptions.Interval
		var currentDUE time.Time
		currentDUE, err = todo.Props.DateTime("DUE", nil)
		if err != nil {
			return
		}
		switch rOptions.Freq {
		case rrule.DAILY:
			currentDUE = currentDUE.AddDate(0, 0, offset)
		case rrule.WEEKLY:
			currentDUE = currentDUE.AddDate(0, 0, 7*offset)
		case rrule.MONTHLY:
			currentDUE = currentDUE.AddDate(0, offset, 0)
		case rrule.YEARLY:
			currentDUE = currentDUE.AddDate(offset, 0, 0)
		case rrule.HOURLY: //TODO didn't debug
			currentDUE = currentDUE.Add(time.Hour * time.Duration(offset))
		case rrule.MINUTELY: //TODO didn't debug
			currentDUE = currentDUE.Add(time.Minute * time.Duration(offset))
		case rrule.SECONDLY: //TODO didn't debug
			currentDUE = currentDUE.Add(time.Second * time.Duration(offset))
			//TODO case if nothing changed
		}

		todo.Props.SetDateTime("DUE", currentDUE)
	} else {
		todo.Props.SetText(ical.PropStatus, "COMPLETED")
		todo.Props.SetText(ical.PropPercentComplete, "100")
		todo.Props.SetDateTime(ical.PropCompleted, time.Now())
	}

	err = m.EditTodo(todo)
	if err != nil {
		return
	}

	return nil
}
