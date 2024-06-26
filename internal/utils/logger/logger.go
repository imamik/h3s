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

	log.Println(logLine...)
}
