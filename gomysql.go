// gomysql project gomysql.go
package gomysql

import (
	"database/sql"
)

func init() {
	sql.Register("mysql", &myDirver{})
}
