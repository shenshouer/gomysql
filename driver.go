package gomysql

import (
	"crypto/tls"
	"database/sql/driver"
	"fmt"
	"net/url"
	"strconv"
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
		cfg := &config{ // default config
			socket:    "/var/run/mysqld/mysqld.sock",
			host:      "localhost",
			port:      3306,
			user:      "root",
			collation: defaultCollation,
		}
		// only support mysql and mysqls
		// maybe support other database later
		switch u.Scheme {
		case "mysql":
		case "mysqls":
			cfg.tls = &tls.Config{}
		default:
			return nil, fmt.Errorf("invalid scheme: %s", dsn)
		}

		// other option will be support
		for k, v := range u.Query() {
			switch k {
			case "debug":
				cfg.debug = true
			case "skip-verify":
				if cfg.tls != nil {
					cfg.tls.InsecureSkipVerify = true
				}
			case "allow-insecure-local-infile":
				cfg.allowLocalInfile = true
			case "charset":
				cfg.collation = collations[v[0]]
			case "socket":
				cfg.socket = v[0]
			case "strict":
				cfg.strict = true
			default:
				return nil, fmt.Errorf("invalid parameter: %s", k)
			}
		}

		if len(u.Host) > 0 {
			host_port := strings.SplitN(u.Host, ":", 2)
			cfg.host = host_port[0]

			if len(host_port) == 2 {
				cfg.port, err = strconv.Atoi(host_port[1])
				if err != nil {
					return nil, fmt.Errorf("invalid port: %s", dsn)
				}
			}
		}

		if u.User != nil {
			cfg.user = u.User.Username()
			if p, ok := u.User.Password(); ok {
				cfg.passwd = p
			}
		}

		// database name
		if len(u.Path) > 0 {
			path := strings.SplitN(u.Path, "/", 2)
			cfg.dbname = path[1]
		}

		if u.Host == "(unix)" {
			cfg.net = "unix"
		} else {
			cfg.net = "tcp"
		}

	}
	return &connection{}, nil
}
