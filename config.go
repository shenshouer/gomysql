package gomysql

import (
	"crypto/tls"
	"time"
)

type config struct {
	host             string        // db host
	port             int           // db post
	user             string        // db username
	passwd           string        // db password
	net              string        // db protocol unix or tcp; According to the DSN Host
	dbname           string        // database name in database server
	socket           string        // the socket file for the protocol of unix
	tls              *tls.Config   // tls connection
	timeout          time.Duration // dail timeout
	collation        uint8         // charset
	debug            bool          // debug mode
	strict           bool          // strict mode
	allowLocalInfile bool          //
	clientFoundRows  bool
}
