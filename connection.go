package gomysql

import (
	"crypto/tls"
	"database/sql/driver"
	"fmt"
	"io"
	"net"
)

// implements from database/sql/driver.Conn
// and database/sql/driver.Queryer
type connection struct {
	conn               net.Conn
	cfg                *config
	buf                buffer
	seq                byte
	protocolVersion    byte
	serverVersion      string
	version            []byte
	connId             uint32
	serverCapabilities uint16
	serverStatus       uint16
	serverLanguage     uint8
	affectedRows       uint64
	insertId           uint64
	status             uint16
}

// implements from database/sql/driver.Conn
func (this *connection) Prepare(query string) (driver.Stmt, error) {
	logger.Infof("connection.Prepare(%s)", query)
	return &stmt{}, nil
}

// implements from database/sql/driver.Conn
func (this *connection) Close() error {
	logger.Info("connection.Close")
	return nil
}

// implements from database/sql/driver.Conn
func (this *connection) Begin() (driver.Tx, error) {
	logger.Info("connection.Begin")
	return &tx{}, nil
}

// implements from database/sql/driver.Queryer
func (this *connection) Query(query string, args []driver.Value) (driver.Rows, error) {
	logger.Infof("connection.Query(%s, %v)", query, args)
	if this.conn == nil {
		if this.cfg.debug {
			logger.Error(driver.ErrBadConn)
		}
		return nil, driver.ErrBadConn
	}

	if len(args) > 0 {
		if !this.cfg.interpolateParams {
			return nil, driver.ErrSkip
		}
		// try client-side prepare to reduce roundtrip
		/* unsupport now
		prepared, err := this.interpolateParams(query, args)
		if err != nil {
			return nil, err
		}
		query = prepared
		args = nil
		*/
	}

	if this.cfg.debug {
		logger.Infof("connection.Query %s", query)
	}

	return this.query(query)
}

func (this *connection) query(query string) (rw *rows, err error) {
	logger.Infof("connection.query << test ================ test >>")
	if len(query) > MAX_PACKET_SIZE {
		return nil, fmt.Errorf("query exceeds %d bytes", MAX_PACKET_SIZE)
	}

	p := this.newComPacket(COM_QUERY)
	p.WriteString(query)
	if err = this.sendPacket(p); err != nil {
		logger.Infof("connection.query  << test ================ test >>")
		return nil, err
	}

	if p, err = this.recvPacket(); err != nil {
		logger.Infof("connection.query  << test ================ test >>")
		return nil, err
	}

	rw = &rows{conn: this}
	switch p.FirstByte() {
	case OK:
		logger.Infof("connection.query  << test ================ test >>")
		if err = this.handleOkPacket(p); err != nil {
			logger.Error(err)
		}
		return rw, nil
	case ERR:
		logger.Infof("connection.query  <<  ================  >> ERR")
		return nil, p.ReadErr()
	case LOCAL_INFILE:
		p.ReadUint8()
		fn := string(p.Bytes())
		/*if err := cn.sendLocalFile(r, fn); err != nil {
			return nil, err
		}*/
		logger.Infoln("<<<<<=============>>>>> fn:%s", fn)
		return nil, fmt.Errorf("===================>> connect.query LOCAL_INFILE")
	default:
		logger.Infoln("==========================>> default")
		n, _ := p.ReadLCUint64()
		if rw.columns, err = this.readColumns(int(n)); err != nil {
			return nil, err
		}
		return rw, nil
	}
}

// Read Packets as Field Packets until EOF-Packet or an Error appears
// http://dev.mysql.com/doc/internals/en/com-query-response.html#packet-Protocol::ColumnDefinition41
func (this *connection) readColumns(n int) ([]column, error) {
	if n == 0 {
		return nil, nil
	}
	cols := make([]column, n)
	for i := range cols {
		if p, err := this.recvPacket(); err != nil {
			return nil, err
		} else {
			col := &cols[i]
			p.SkipLCBytes() // catalog
			p.SkipLCBytes() // schema Database [len coded string]

			if this.cfg.columnsWithAlias { // table Table [len coded string]
				col.tableName, _ = p.ReadLCString()
			} else {
				p.SkipLCBytes()
			}

			p.SkipLCBytes()                // org_table Original table [len coded string]
			col.name, _ = p.ReadLCString() // name
			p.SkipLCBytes()                // org_name
			p.ReadLCUint64()               // 0x0c Filler [uint8]
			col.charset = p.ReadUint16()   // Charset [charset, collation uint8]
			col.length = p.ReadUint32()    // Length [uint32]
			col.coltype = p.ReadUint8()    // Field type [uint8]
			col.flags = p.ReadUint16()     // Flags [uint16]
			col.decimals = p.ReadUint8()   // Decimals [uint8]
		}
	}
	p, err := this.recvPacket()
	if err != nil {
		return nil, err
	}
	if x := p.ReadUint8(); x != EOF {
		return nil, fmt.Errorf("readColumns: expected EOF, got %v", x)
	}
	return cols, nil
}

// Result Set Header Packet
// http://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::Resultset
/*func (this *connection) readResultSetHeaderPacket() (int, error) {
	p, err := this.recvPacket()
	if err != nil {
		return 0, err
	}
	switch p.FirstByte() {
	case OK:
	case ERR:
	case LOCAL_INFILE:
	default:
	}
}*/

func (this *connection) handshake() error {
	if challange, err := this.readInitPacket(); err != nil {
		if this.cfg.debug {
			logger.Errorf("handshake err:%v", err)
		}

		return err
	} else {
		if err := this.writeAuthPacket(challange, 0); err != nil {
			return err
		}

		logger.Infof("=========>> connection.seq:%d", this.seq)
		var p packet
		if this.seq, err = p.recv(this.conn, this.seq); err != nil {
			if err != io.EOF {
				return err
			}
		}

		switch p.FirstByte() {
		//switch p.ReadUint8() {
		case OK:
			logger.Info("================>> ok")
			//this.handleOkPacket(p)
		case ERR:
			logger.Info("================>> err")
			return p.ReadErr()
		default:
			logger.Info("================>> default")
			return fmt.Errorf("hello: expected OK or ERR, got %v", p.FirstByte())
		}

		return nil
	}
}

// Ok Packet
// http://dev.mysql.com/doc/internals/en/generic-response-packets.html#packet-OK_Packet
func (this *connection) handleOkPacket(p packet) error {
	logger.Infoln("================>> handleOkPacket", p)
	//return nil
	// 0x00 [1 byte]
	this.affectedRows = p.ReadUint64() // Affected rows [Length Coded Binary]
	this.insertId = p.ReadUint64()     // Insert id [Length Coded Binary]
	this.status = p.ReadUint16()       // server_status [2 bytes]
	logger.Infof("connect.handleOkPacket :: this.affectedRows:%d, this.insertId:%d, this.status:%d", this.affectedRows, this.insertId, this.status)
	// warning count [2 bytes]
	if this.cfg.strict && p.ReadUint16() > 0 {
		return this.getWarnings()
	}
	return nil
}

// Client Authentication Packet
// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeResponse
func (this *connection) writeAuthPacket(challange []byte, flags uint32) error {
	if this.cfg.debug {
		logger.Infof("connection.writeAuthPacket(%s, %d)", string(challange[:]), flags)
	}
	p := newPacket()
	// Adjust client flags based on server support
	flags |= CLIENT_PROTOCOL_41 |
		CLIENT_SECURE_CONNECTION |
		CLIENT_LONG_PASSWORD |
		CLIENT_TRANSACTIONS |
		CLIENT_LOCAL_FILES |
		CLIENT_LONG_FLAG

	if len(this.cfg.dbname) > 0 {
		flags |= CLIENT_CONNECT_WITH_DB
	}
	if this.cfg.clientFoundRows {
		flags |= CLIENT_FOUND_ROWS
	}

	p.WriteUint32(flags)            // ClientFlags [32 bit]
	p.WriteUint32(MAX_PACKET_SIZE)  // MaxPacketSize [32 bit] (none)
	p.WriteByte(this.cfg.collation) // Charset [1 byte]
	p.Write(make([]byte, 23))       // none

	// SSL Connection Request Packet
	// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::SSLRequest
	if this.cfg.tls != nil && this.serverCapabilities&CLIENT_SSL != 0 {
		// Send TLS / SSL request packet
		if err := this.sendPacket(p); err != nil {
			if this.cfg.debug {
				logger.Errorf("Send TLS/SSL request packet err:%v", err)
			}
			return err
		}
		// Switch to TLS
		tlsConn := tls.Client(this.conn, this.cfg.tls)
		if err := tlsConn.Handshake(); err != nil {
			if this.cfg.debug {
				logger.Errorf("Handshake with tls Connection err:%v", err)
			}
			return err
		}
		this.conn = tlsConn
		this.buf.rd = tlsConn
		//this.bufrd = bufio.NewReader(cn.netconn)
	}

	p.WriteString(this.cfg.user)
	p.WriteByte(0) // End of string

	if len(this.cfg.passwd) != 0 {
		token := passwordToken(this.cfg.passwd, challange)
		p.WriteByte(byte(len(token)))
		p.Write(token)
	} else {
		p.WriteByte(0) // End of password
	}
	if len(this.cfg.dbname) > 0 {
		p.WriteString(this.cfg.dbname)
		p.WriteByte(0) // End of database name
	}

	err := this.sendPacket(p)
	if err != nil && this.cfg.debug {
		logger.Errorf("auth err:%v", err)
	}
	return err
}

func (this *connection) readInitPacket() (challange []byte, err error) {
	var p packet
	if this.seq, err = p.recv(this.conn, this.seq); err != nil {
		return nil, err
	}
	this.protocolVersion = p.ReadUint8()
	if s, err := p.ReadString('\x00'); err != nil { // server version [null terminated string]
		return nil, err
	} else {
		this.serverVersion = s[:len(s)-1]
	}
	if this.version, err = parseVersion(this.serverVersion); err != nil {
		logger.Errorf("warning: could not parse server version '%s'\n", this.serverVersion)
	}

	this.connId = p.ReadUint32()             // connection id [4 bytes]
	challange = p.Next(8)                    // first part of the password cipher [8 bytes]
	p.Next(1)                                // (filler) always 0x00 [1 byte]
	this.serverCapabilities = p.ReadUint16() // capability flags (lower 2 bytes) [2 bytes]
	this.serverLanguage = p.ReadUint8()      // character set [1 byte]
	this.serverStatus = p.ReadUint16()       // status flags [2 bytes]
	// capability flags (upper 2 bytes) [2 bytes]
	// length of auth-plugin-data [1 byte]
	// reserved (all [00]) [10 bytes]
	p.Next(13)

	// second part of the password cipher [mininum 13 bytes],
	// where len=MAX(13, length of auth-plugin-data - 8)
	//
	// The web documentation is ambiguous about the length. However,
	// according to mysql-5.7/sql/auth/sql_authentication.cc line 538,
	// the 13th byte is "\0 byte, terminating the second part of
	// a scramble". So the second part of the password cipher is
	// a NULL terminated string that's at least 13 bytes with the
	// last byte being NULL.
	//
	// The official Python library uses the fixed length 12
	// which seems to work but technically could have a hidden bug.
	challange = append(challange, p.Next(12)...)
	p.Next(1)

	return challange, nil
}

func (this *connection) newComPacket(com byte) (p packet) {
	this.seq = 0
	p = newPacket()
	p.WriteByte(com)
	return p
}

func (this *connection) sendPacket(p packet) (err error) {
	err = p.send(this.conn, this.seq)
	this.seq += 1
	return err
}

func (this *connection) recvPacket() (p packet, err error) {
	this.seq, err = p.recv(this.buf.rd, this.seq)
	return p, err
}

func (this *connection) getWarnings() error {
	rows, err := this.Query("SHOW WARNINGS", nil)
	if err != nil {
		return err
	}
	var warnings = MySQLWarnings{}
	var values = make([]driver.Value, 3)

	for {
		err = rows.Next(values)
		switch err {
		case nil:
			warning := MySQLWarning{}

			if raw, ok := values[0].([]byte); ok {
				warning.Level = string(raw)
			} else {
				warning.Level = fmt.Sprintf("%s", values[0])
			}
			if raw, ok := values[1].([]byte); ok {
				warning.Code = string(raw)
			} else {
				warning.Code = fmt.Sprintf("%s", values[1])
			}
			if raw, ok := values[2].([]byte); ok {
				warning.Message = string(raw)
			} else {
				warning.Message = fmt.Sprintf("%s", values[0])
			}

			warnings = append(warnings, warning)

		case io.EOF:
			return warnings

		default:
			rows.Close()
			return nil
		}
	}
}
