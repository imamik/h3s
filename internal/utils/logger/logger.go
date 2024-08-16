package logger

import (
	"fmt"
	"log"
	"strings"
)

var (
	logOnlyLatest = true  // logOnlyLatest determines if only the latest event should be logged
	logLive       = false // logLive determines if events should be logged live or stacked until explicitly logged
)

const (
	ColorGreen = "\033[32m" // ColorGreen is the color code to print green in the terminal
	ColorRed   = "\033[31m" // ColorRed is the color code to print red in the terminal
	ColorReset = "\033[0m"  // ColorReset is the color code to reset the terminal color
)

// normalize returns a string with a fixed length
func normalize(s string, l int) string {
	if len(s) > l-3 {
		return s[:l-3] + "..."
	}
	return s + strings.Repeat(" ", l-len(s))
}

// LogError prints an error message to the console
func LogError(err ...any) {
	log.Println(err...)
}

// NewEventLogger creates a new logger
func NewEventLogger(resource LogResource, action interface{}, id string) (AddEventFunc, LogFunc) {
	var events []ResourceEvent
	addEvent := AddEventFunc(func(status LogCrudStatus, err ...any) {
		events = append(events, ResourceEvent{
			Resource: resource,
			ID:       id,
			Action:   action,
			Status:   status,
			Err:      err,
		})
		if logLive {
			LogResourceEvent(resource, action, id, status, err...)
		}
	})
	logEvents := LogFunc(func() {
		if logOnlyLatest {
			latest := events[len(events)-1]
			LogResourceEvent(latest.Resource, latest.Action, latest.ID, latest.Status, latest.Err...)
		} else {
			for _, e := range events {
				LogResourceEvent(e.Resource, e.Action, e.ID, e.Status, e.Err...)
			}
		}
	})
	addEvent(Initialized)
	return addEvent, logEvents
}

// LogResourceEvent logs a resource event
func LogResourceEvent(
	resource LogResource,
	action interface{},
	id string,
	status LogCrudStatus,
	err ...any,
) {
	var actionStr string

	switch v := action.(type) {
	case LogCrudMethod:
		actionStr = string(v)
	case string:
		actionStr = v
	default:
		log.Println("Invalid method type")
		return
	}

	var logLine []any
	logLine = append(
		logLine,
		normalize(actionStr, 24),
		normalize(string(resource), 24),
		normalize(id, 48),
		normalize(string(status), 24),
	)

	for _, e := range err {
		logLine = append(logLine, e)
	}

	switch status {
	case Success:
		log.Println(ColorGreen + fmt.Sprint(logLine...) + ColorReset)
	case Failure:
		log.Println(ColorRed + fmt.Sprint(logLine...) + ColorReset)
	default:
		log.Println(logLine...)
	}
}
