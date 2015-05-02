package gomysql

import (
	"database/sql/driver"
	"fmt"
	"net"
	"net/url"
	"strings"
)

// implements from database/sql/driver.Driver
type myDirver struct {
}

// implements from database/sql/driver.Driver
// @param dsn is Data Source Name
func (this *myDirver) Open(dsn string) (driver.Conn, error) {
	logger.Infof("driver.Open(%s)", dsn)
	return this.open(dsn)
}

func (this *myDirver) open(dsn string) (*connection, error) {
	name := strings.Trim(dsn, " ")
	if len(name) == 0 {
		return nil, fmt.Errorf("Invaild Data Source Name!")
	}

	if u, err := url.Parse(dsn); err != nil {
		return nil, fmt.Errorf("Parse DSN err:%v", err)
	} else {
		fmt.Println(u.Scheme)
		fmt.Println(u.User)
		fmt.Println(u.User.Username())
		p, _ := u.User.Password()
		fmt.Println(p)
		fmt.Println(u.Host)
		host, port, _ := net.SplitHostPort(u.Host)
		fmt.Println(host)
		fmt.Println(port)
		fmt.Println(u.Path)
		fmt.Println(u.Fragment)
		fmt.Println(u.RawQuery)
		m, _ := url.ParseQuery(u.RawQuery)
		fmt.Println(m)
		for k, v := range m {
			fmt.Printf("k=%s, v=%s \n", k, v)
		}
	}
	return &connection{}, nil
}
