package logger

import (
	"fmt"
	"h3s/internal/utils/str"
	"log"
)

// EventLogger is a struct that collects and logs events
type EventLogger struct {
	action        interface{}
	parent        *EventLogger
	resource      LogResource
	id            string
	events        []ResourceEvent
	logOnlyLatest bool
	logLive       bool
}

// New creates a new event logger
func New(parent *EventLogger, resource LogResource, action interface{}, id string) *EventLogger {
	e := EventLogger{
		parent:        parent,
		events:        []ResourceEvent{},
		logOnlyLatest: false,
		logLive:       false,
		resource:      resource,
		action:        action,
		id:            id,
	}
	e.AddEvent(Initialized)
	return &e
}

// AddEvent adds an event to the event logger
func (l *EventLogger) AddEvent(status LogCrudStatus, err ...any) {
	event := ResourceEvent{
		Resource: l.resource,
		ID:       l.id,
		Action:   l.action,
		Status:   status,
		Err:      err,
		Depth:    0,
		IsFirst:  len(l.events) == 0,
		IsLast:   false,
	}
	l.events = append(l.events, event)
	if l.logLive {
		l.LogLatest()
	}
}

// AppendEvents appends a list of events to the event logger
func (l *EventLogger) AppendEvents(events []ResourceEvent) {
	l.events = append(l.events, events...)
}

// LogEvents logs all events in the event logger (or only the latest event if logOnlyLatest is set)
func (l *EventLogger) LogEvents() {
	// Set the last event to be the last event
	if len(l.events) > 0 {
		l.events[len(l.events)-1].IsLast = true
	}
	if l.parent != nil {
		// Increase the depth of the events
		for i := range l.events {
			l.events[i].Depth++
		}
		l.parent.AppendEvents(l.events)
	} else if l.logOnlyLatest || l.logLive {
		l.LogLatest()
	} else {
		for _, event := range l.events {
			l.logEvent(event, true)
		}
	}
	l.events = []ResourceEvent{}
}

func (l *EventLogger) LogLatest() {
	lastIndex := len(l.events) - 1
	event := l.events[lastIndex]
	l.logEvent(event, false)
}

// logEvent logs the event with the appropriate color
func (l *EventLogger) logEvent(e ResourceEvent, group bool) {
	// Build the action string
	var actionStr string
	switch v := e.Action.(type) {
	case LogCrudMethod:
		actionStr = string(v)
	case string:
		actionStr = v
	default:
		log.Println("Invalid method type")
		return
	}

	// Build the log line, composed of the action, resource, id, status and error with normalized lengths
	var logLine []any
	if group {
		// for every depth level, add the appropriate line
		for i := 0; i < e.Depth; i++ {
			logLine = append(logLine, "│ ")
		}
		if e.IsFirst && e.IsLast {
			logLine = append(logLine, "●─")
		} else if e.IsFirst {
			logLine = append(logLine, "┌─")
		} else if e.IsLast {
			logLine = append(logLine, "└─")
		} else {
			logLine = append(logLine, "├─")
		}
	}
	logLine = append(
		logLine,
		str.NormalizeLength(actionStr, 24),
		str.NormalizeLength(string(e.Resource), 24),
		str.NormalizeLength(e.ID, 48),
		str.NormalizeLength(string(e.Status), 24),
	)

	// Add the error messages to the log line
	for _, e := range e.Err {
		logLine = append(logLine, e)
	}

	// Log the line with the appropriate color (green for success, red for failure)
	switch e.Status {
	case Success:
		log.Println(ColorGreen + fmt.Sprint(logLine...) + ColorReset)
	case Failure:
		log.Println(ColorRed + fmt.Sprint(logLine...) + ColorReset)
	default:
		log.Println(ColorDefault + fmt.Sprint(logLine...) + ColorReset)
	}
}
