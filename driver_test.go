package gomysql

import (
	"database/sql"
	"testing"
)

const (
	dsn = "mysql://gopher2:secret@localhost:3306/test?strict&client-multi-results"
)

func TestOpen(t *testing.T) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if rows, err := db.Query("select * from tt"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(rows)
	}
}
