package gomysql

import (
	"fmt"
)

// MySQLWarning is an error type which represents a single MySQL warning.
// Warnings are returned in groups only. See MySQLWarnings
type MySQLWarning struct {
	Level   string
	Code    string
	Message string
}

// MySQLWarnings is an error type which represents a group of one or more MySQL
// warnings
type MySQLWarnings []MySQLWarning

func (mws MySQLWarnings) Error() string {
	var msg string
	for i, warning := range mws {
		if i > 0 {
			msg += "\r\n"
		}
		msg += fmt.Sprintf(
			"%s %s: %s",
			warning.Level,
			warning.Code,
			warning.Message,
		)
	}
	return msg
}
