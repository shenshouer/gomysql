package gomysql

import (
	"github.com/shenshouer/logging"
)

func NewLogger() logging.Logger {
	return logging.NewSimpleLogger()
}

var logger logging.Logger = NewLogger()
