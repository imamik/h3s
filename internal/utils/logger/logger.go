// Package logger contains the functionality for constructing a logger, that can be used to log events in a structured way
package logger

import (
	"fmt"
	"log"
	"time"
)

const (
	errFormatMessage = "Error formatting log message"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	// LevelDebug represents debug level messages
	LevelDebug LogLevel = "DEBUG"
	// LevelInfo represents informational messages
	LevelInfo LogLevel = "INFO"
	// LevelWarning represents warning messages
	LevelWarning LogLevel = "WARNING"
	// LevelError represents error messages
	LevelError LogLevel = "ERROR"
)

// LogFields represents structured logging fields
type LogFields map[string]interface{}

// EventLogger is a struct that collects and logs events
type EventLogger struct {
	action        interface{}
	fields        LogFields
	parent        *EventLogger
	resource      LogResource
	level         LogLevel
	id            string
	events        []ResourceEvent
	logOnlyLatest bool
	logLive       bool
}

// ResourceEvent represents a single event in the event logger
type ResourceEvent struct {
	Time     time.Time
	Action   interface{}
	Fields   LogFields
	Resource LogResource
	ID       string
	Status   LogCrudStatus
	Level    LogLevel
	Err      []any
	Depth    int
	IsFirst  bool
	IsLast   bool
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
		fields:        make(LogFields),
		level:         LevelInfo,
	}
	e.AddEvent(Initialized)
	return &e
}

// WithFields adds fields to the logger
func (l *EventLogger) WithFields(fields LogFields) *EventLogger {
	for k, v := range fields {
		l.fields[k] = v
	}
	return l
}

// WithLevel sets the log level
func (l *EventLogger) WithLevel(level LogLevel) *EventLogger {
	l.level = level
	return l
}

// AddEvent adds an event to the event logger
func (l *EventLogger) AddEvent(status LogCrudStatus, err ...any) {
	level := l.level
	if len(err) > 0 {
		level = LevelError
	}

	event := ResourceEvent{
		Resource: l.resource,
		ID:       l.id,
		Action:   l.action,
		Status:   status,
		Err:      err,
		Depth:    0,
		IsFirst:  len(l.events) == 0,
		IsLast:   false,
		Time:     time.Now(),
		Fields:   l.fields,
		Level:    level,
	}
	l.events = append(l.events, event)
	if l.logLive {
		l.LogLatest()
	}
}

// LogEvents logs all events in the event logger
func (l *EventLogger) LogEvents() {
	if len(l.events) > 0 {
		l.events[len(l.events)-1].IsLast = true
	}
	switch {
	case l.parent != nil:
		for i := range l.events {
			l.events[i].Depth++
		}
		l.parent.AppendEvents(l.events)
	case l.logOnlyLatest || l.logLive:
		l.LogLatest()
	default:
		for _, event := range l.events {
			l.logEvent(event, true)
		}
	}
	l.events = []ResourceEvent{}
}

// getActionString converts the action interface to a string
func getActionString(action interface{}) (string, bool) {
	switch v := action.(type) {
	case LogCrudMethod:
		return string(v), true
	case string:
		return v, true
	default:
		return "", false
	}
}

// createLogData creates the log fields from a resource event
func createLogData(e ResourceEvent) LogFields {
	return LogFields{
		"timestamp": e.Time.Format(time.RFC3339),
		"level":     e.Level,
		"resource":  e.Resource,
		"action":    e.Action,
		"status":    e.Status,
	}
}

// logEvent logs the event with structured data
func (l *EventLogger) logEvent(e ResourceEvent, group bool) {
	actionStr, ok := getActionString(e.Action)
	if !ok {
		log.Println("Invalid method type")
		return
	}

	logData := createLogData(e)
	logData["action"] = actionStr

	if e.ID != "" {
		logData["id"] = e.ID
	}

	if len(e.Err) > 0 {
		logData["errors"] = e.Err
	}

	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			logData[k] = v
		}
	}

	if group {
		if e.IsFirst {
			logData["group"] = "start"
		} else if e.IsLast {
			logData["group"] = "end"
		}
	}

	msg := formatLogMessage(logData)
	log.Println(msg)
}

// extractRequiredField extracts a required field from LogFields with type assertion
func extractRequiredField[T any](data LogFields, key string) (T, bool) {
	value, ok := data[key]
	if !ok || value == nil {
		log.Printf("Error: %s field is missing or nil", key)
		return *new(T), false
	}
	cast, ok := value.(T)
	if !ok {
		log.Printf("Error: %s field is not of expected type", key)
		return *new(T), false
	}
	return cast, true
}

// formatBasicMessage formats the basic message with required fields
func formatBasicMessage(data LogFields) (string, bool) {
	resource, ok := extractRequiredField[LogResource](data, "resource")
	if !ok {
		return errFormatMessage, false
	}

	action, ok := extractRequiredField[string](data, "action")
	if !ok {
		return errFormatMessage, false
	}

	status, ok := extractRequiredField[LogCrudStatus](data, "status")
	if !ok {
		return errFormatMessage, false
	}

	level, _ := data["level"].(string)
	return fmt.Sprintf("[%s] %s %s %s", level, resource, action, status), true
}

// formatLogMessage formats the log message with all fields
func formatLogMessage(data LogFields) string {
	msg, ok := formatBasicMessage(data)
	if !ok {
		return errFormatMessage
	}

	if id, ok := data["id"]; ok {
		msg += fmt.Sprintf(" (%s)", id)
	}

	if errors, ok := data["errors"]; ok {
		msg += fmt.Sprintf(" - Errors: %v", errors)
	}

	// Add other fields
	skipFields := map[string]bool{
		"resource": true, "action": true, "status": true,
		"id": true, "errors": true, "level": true, "timestamp": true,
	}

	for k, v := range data {
		if !skipFields[k] {
			msg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	return msg
}

// LogLatest logs the latest event in the event logger
func (l *EventLogger) LogLatest() {
	lastIndex := len(l.events) - 1
	event := l.events[lastIndex]
	l.logEvent(event, false)
}

// AppendEvents appends a list of events to the event logger
func (l *EventLogger) AppendEvents(events []ResourceEvent) {
	l.events = append(l.events, events...)
}
