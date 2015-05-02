package gomysql

import (
	"database/sql/driver"
)

type connection struct {
}

func (this *connection) Prepare(query string) (driver.Stmt, error) {
	logger.Infof("connection.Prepare(%s)", query)
	return &stmt{}, nil
}

func (this *connection) Close() error {
	logger.Info("connection.Close")
	return nil
}

func (this *connection) Begin() (driver.Tx, error) {
	logger.Info("connection.Begin")
	return &tx{}, nil
}
