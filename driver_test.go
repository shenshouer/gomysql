package gomysql

import (
	"database/sql"
	"testing"
)

const (
	dsn1 = "mysql://root:123456@localhost:3306/test?strict&debug"
	dsn2 = "mysqls://root:123456@localhost:3306/test?strict&debug"
)

func TestOpen(t *testing.T) {
	db, err := sql.Open("mysql", dsn2)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if rows, err := db.Query("select * from tbl1"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(rows)
	}
}
