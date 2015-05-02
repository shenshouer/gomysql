package gomysql

import (
	"database/sql/driver"
)

type stmt struct {
}

func (this *stmt) Close() error {
	logger.Info("stmt.Close")
	return nil
}

func (this *stmt) NumInput() int {
	logger.Info("stmt.NumInput")
	return 0
}

func (this *stmt) Exec(args []driver.Value) (driver.Result, error) {
	logger.Infof("stmt.Exec %v", args)
	return &result{}, nil
}

func (this *stmt) Query(args []driver.Value) (driver.Rows, error) {
	logger.Infof("stmt.Query %v", args)
	return &rows{}, nil
}
