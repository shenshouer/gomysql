package gomysql

import (
	"crypto/tls"
	"database/sql/driver"
	"fmt"
	"net"
)

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
		var p packet
		if this.seq, err = p.recv(this.conn, this.seq); err != nil {
			return err
		}

		switch p.FirstByte() {
		case OK:
		case ERR:
			return p.ReadErr()
		default:
			return fmt.Errorf("hello: expected OK or ERR, got %v", p.FirstByte())
		}

		return nil
	}
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
	if s, err := p.ReadString('\x00'); err != nil {
		return nil, err
	} else {
		this.serverVersion = s[:len(s)-1]
	}
	if this.version, err = parseVersion(this.serverVersion); err != nil {
		logger.Errorf("warning: could not parse server version '%s'\n", this.serverVersion)
	}

	this.connId = p.ReadUint32()
	challange = p.Next(8)
	p.Next(1)
	this.serverCapabilities = p.ReadUint16()
	this.serverLanguage = p.ReadUint8()
	this.serverStatus = p.ReadUint16()
	p.Next(13)
	challange = append(challange, p.Next(12)...)
	p.Next(1)

	return challange, nil
}

func (this *connection) sendPacket(p packet) (err error) {
	err = p.send(this.conn, this.seq)
	this.seq += 1
	return err
}
