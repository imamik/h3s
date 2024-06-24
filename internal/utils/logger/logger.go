package logger

import (
	"log"
	"strings"
)

func normalize(s string, l int) string {
	if len(s) > l-3 {
		return s[:l-3] + "..."
	}
	return s + strings.Repeat(" ", l-len(s))
}

func LogError(err ...any) {
	log.Println(err)
}

func LogResourceEvent(
	resource LogResource,
	method LogCrudMethod,
	id string,
	status LogCrudStatus,
	err ...any,
) {
	var logLine []any
	logLine = append(
		logLine,
		normalize(string(method), 10),
		normalize(string(resource), 20),
		normalize(id, 48),
		normalize(string(status), 16),
	)

	for _, e := range err {
		logLine = append(logLine, e)
	}

	log.Println(logLine...)
}
