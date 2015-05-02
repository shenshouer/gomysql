package gomysql

import (
	"database/sql/driver"
)

type rows struct {
}

func (this *rows) Columns() []string {
	logger.Info("rows.Columns")
	return nil
}

func (this *rows) Close() error {
	logger.Info("rows.Close")
	return nil
}

func (this *rows) Next(dest []driver.Value) error {
	logger.Infof("rows.Next", dest)
	return nil
}
