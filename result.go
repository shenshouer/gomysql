package gomysql

import (
	"fmt"
)

type result struct {
	rowsAffected int64
	lastInsertId int64
	warnings     uint16
	closed       bool
}

func (this *result) LastInsertId() (int64, error) {
	logger.Info("result.LastInsertId")
	return 0, nil
}

func (this *result) RowsAffected() (int64, error) {
	logger.Info("result.RowsAffected")
	return 0, nil
}

func (this *result) readOK(p packet) error {
	this.rowsAffected, this.lastInsertId, this.warnings = p.ReadOK()
	this.closed = true
	//return this.ReadWarnings()
	return fmt.Errorf("result.readOk")
}
