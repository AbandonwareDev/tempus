package main

import (
	// "github.com/studio-b12/gowebdav"
	"context"
	// "fmt"
	webdav "github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"

	// "log"
	"time"
	"strings"
)

type TODO struct {
	Name string
	Description string
	Time string
	// Priority int //TODO
	// Subtasks []TODO //TODO
	// Repeat //TODO
	// Alarm //TODO
}


var clientWebDAV *webdav.Client
var client *caldav.Client // clientCalDAV
// var calendarObjects []caldav.CalendarObject
var ctx = context.Background()
// var authSession caldav.Client // clientCalDAV

func (options *Options) InitDAVclients() error {
	
	var err error
	authSession := webdav.HTTPClientWithBasicAuth(nil, options.User, options.Password)

	clientWebDAV, err = webdav.NewClient(authSession, options.URL)
	if err != nil {
		// Handle error
		return err
	}
	client, err = caldav.NewClient(authSession, options.URL)
	if err != nil {
		// Handle error
		return err
	}
	return nil
}

func GetCalendars() ([]caldav.Calendar,error) {

	principal, err := clientWebDAV.FindCurrentUserPrincipal(ctx)
	if err != nil {
		// Handle error
		return nil,err
	}
	// fmt.Println("principal: ",principal) TODO log


	// Find the calendar home set
	calendarHomeSet, err := client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		return nil,err
	}

	// Find calendars in the calendar home set
	calendars, err := client.FindCalendars(ctx, calendarHomeSet)
	if err != nil {
		return nil,err
	}

	return calendars, nil
}

func GetTODOs(calendarPath string) (calendarObjects []caldav.CalendarObject,err error) {
	date := time.Now()
		dateStart:= time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location()).AddDate(0, 0, -1) //TODO too complex - time.Now().Add(-92 * time.Hour),
	// dateEnd:= time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location()) //date +1 day
	dateEnd:= date.AddDate(0,1,0) //TODO we hard limit pulling events only for this month, to prevent getting tooo much events

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
							Name: "VTODO",
							Start:   dateStart,
							End:   dateEnd,
							// FYI: server-side filtering is not supported for Nextcloud
						},
					},
				},
	}
	calendarObjects, err = client.QueryCalendar(ctx, calendarPath, &calQuery) 
	if err != nil {
		return nil,err
	}
	return calendarObjects,nil
}

func ParseDueDateTODOs(calObjs []caldav.CalendarObject,date time.Time) ([]TODO,error) {
		var output []TODO

	for _,calObj := range calObjs {
		// fmt.Println((*(*calObj.Data).Children[0]).Name)
		// TODO STATUS map[] COMPLETED
		for _,event := range (*calObj.Data).Children { 
			// if (*event).Name == "VTODO" {
				var notCompletedTODO, withDate, fromToday bool
				//TODO we can optimize there if we encounter wrong state to forcefully stop next analysis
				// notCompletedTODO
				if (*event).Props["COMPLETED"] == nil {
					if ((*event).Props["STATUS"] == nil) { notCompletedTODO = true } else {
						if ((*event).Props["STATUS"][0].Value != "COMPLETED") { notCompletedTODO = true }
					}
				}

				// withTodayDate
				if (*event).Props["DUE"] != nil {
					withDate = true
					// fromToday
					if strings.HasPrefix((*event).Props["DUE"][0].Value, date.Format("20060102")) { fromToday = true }
				}
				
				if notCompletedTODO && withDate && fromToday {

					var tmpTODO TODO
					
					if ((*event).Props["SUMMARY"] != nil) {
						name := (*event).Props["SUMMARY"][0].Value
						tmpTODO.Name = name
					} else {
						tmpTODO.Name = "<EMPTY>"
					}					
					if ((*event).Props["DESCRIPTION"] != nil) {
						description := (*event).Props["DESCRIPTION"][0].Value
						tmpTODO.Description = description
					}

					var todoTime string
					due := (*event).Props["DUE"][0].Value
					index := strings.Index(due, "T")
					if index != -1 {
							str := due[index+1:]
							todoTime = str[:2] + ":" + str[2:4]
							
							tmpTODO.Time = todoTime
						}
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
//TODO repeat function???

	return output,nil
}





//TODO on complete -repeat function
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

















