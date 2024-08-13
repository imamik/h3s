package logger

import (
	"fmt"
	"log"
	"strings"
)

var (
	logOnlyLatest = true
	logLive       = false
)

const (
	COLOR_GREEN = "\033[32m"
	COLOR_RED   = "\033[31m"
	COLOR_RESET = "\033[0m"
)

func normalize(s string, l int) string {
	if len(s) > l-3 {
		return s[:l-3] + "..."
	}
	return s + strings.Repeat(" ", l-len(s))
}

func LogError(err ...any) {
	log.Println(err...)
}

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
		log.Println(COLOR_GREEN + fmt.Sprint(logLine...) + COLOR_RESET)
	case Failure:
		log.Println(COLOR_RED + fmt.Sprint(logLine...) + COLOR_RESET)
	default:
		log.Println(logLine...)
	}
}
