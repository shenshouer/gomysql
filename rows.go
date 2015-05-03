package gomysql

import (
	"database/sql/driver"
)

type column struct {
	tableName string
	name      string
	charset   uint16
	length    uint32
	coltype   byte
	flags     uint16
	decimals  byte
}

type rows struct {
	conn    *connection
	columns []column
}

func (this *rows) Columns() []string {
	c := make([]string, len(this.columns))
	for i, col := range this.columns {
		c[i] = col.name
	}
	return c
}

func (this *rows) Close() error {
	logger.Info("rows.Close")
	return nil
}

func (this *rows) Next(dest []driver.Value) error {
	logger.Infof("rows.Next", dest)
	return nil
}
