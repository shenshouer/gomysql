// see http://dev.mysql.com/doc/internals/en/client-server-protocol.html
// for the protocol definition.
package gomysql

// The capability flags are used by the client and server to indicate
// which features they support and want to use
// Protocol::CapabilityFlags:
// See http://dev.mysql.com/doc/internals/en/capability-flags.html for detail
const (
	CLIENT_LONG_PASSWORD                  = 1 << iota /* new more secure passwords */
	CLIENT_FOUND_ROWS                                 /* Found instead of affected rows */
	CLIENT_LONG_FLAG                                  /* Get all column flags */
	CLIENT_CONNECT_WITH_DB                            /* One can specify db on connect */
	CLIENT_NO_SCHEMA                                  /* Don't allow database.table.column */
	CLIENT_COMPRESS                                   /* Can use compression protocol */
	CLIENT_ODBC                                       /* Odbc client */
	CLIENT_LOCAL_FILES                                /* Can use LOAD DATA LOCAL */
	CLIENT_IGNORE_SPACE                               /* Ignore spaces before '(' */
	CLIENT_PROTOCOL_41                                /* New 4.1 protocol */
	CLIENT_INTERACTIVE                                /* This is an interactive client */
	CLIENT_SSL                                        /* Switch to SSL after handshake */
	CLIENT_IGNORE_SIGPIPE                             /* IGNORE sigpipes */
	CLIENT_TRANSACTIONS                               /* Client knows about transactions */
	CLIENT_RESERVED                                   /* Old flag for 4.1 protocol  */
	CLIENT_SECURE_CONNECTION                          /* New 4.1 authentication */
	CLIENT_MULTI_STATEMENTS                           /* Enable/disable multi-stmt support */
	CLIENT_MULTI_RESULTS                              /* Enable/disable multi-results */
	CLIENT_PS_MULTI_RESULTS                           /* Multi-results in PS-protocol */
	CLIENT_PLUGIN_AUTH                                /* auth plugins supports */
	CLIENT_CONNECT_ATTRS                              /* connection attributes supports */
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA             /* length of auth response data */
	CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS               /* can handle expired passwords */
	CLIENT_SESSION_TRACK                              /* expects the server to send sesson-state changes after a OK packet */
	CLIENT_DEPRECATE_EOF                              /* expects a OK (instead of EOF) after the resultset rows of a Text Resultset. */
)

// Text Protocol
// See http://dev.mysql.com/doc/internals/en/text-protocol.html for detail
const (
	COM_SLEEP = iota
	COM_QUIT
	COM_INIT_DB
	COM_QUERY
	COM_FIELD_LIST
	COM_CREATE_DB
	COM_DROP_DB
	COM_REFRESH
	COM_SHUTDOWN
	COM_STATISTICS
	COM_PROCESS_INFO
	COM_CONNECT
	COM_PROCESS_KILL
	COM_DEBUG
	COM_PING
	COM_TIME
	COM_DELAYED_INSERT
	COM_CHANGE_USER
	COM_BINLOG_DUMP
	COM_TABLE_DUMP
	COM_CONNECT_OUT
	COM_REGISTER_SLAVE
	COM_STMT_PREPARE
	COM_STMT_EXECUTE
	COM_STMT_SEND_LONG_DATA
	COM_STMT_CLOSE
	COM_STMT_RESET
	COM_SET_OPTION
	COM_STMT_FETCH
	COM_DAEMON
	COM_BINLOG_DUMP_GTID
	COM_RESET_CONNECTION
)

// Data Type
// See http://dev.mysql.com/doc/internals/en/com-query-response.html for detail
const (
	MYSQL_TYPE_DECIMAL = iota
	MYSQL_TYPE_TINY
	MYSQL_TYPE_SHORT
	MYSQL_TYPE_LONG
	MYSQL_TYPE_FLOAT
	MYSQL_TYPE_DOUBLE
	MYSQL_TYPE_NULL
	MYSQL_TYPE_TIMESTAMP
	MYSQL_TYPE_LONGLONG
	MYSQL_TYPE_INT24
	MYSQL_TYPE_DATE
	MYSQL_TYPE_TIME
	MYSQL_TYPE_DATETIME
	MYSQL_TYPE_YEAR
	MYSQL_TYPE_NEWDATE
	MYSQL_TYPE_VARCHAR
	MYSQL_TYPE_BIT
	MYSQL_TYPE_TIMESTAMP2
	MYSQL_TYPE_DATETIME2
	MYSQL_TYPE_TIME2
)
const (
	MYSQL_TYPE_NEWDECIMAL = iota + 0xf6
	MYSQL_TYPE_ENUM
	MYSQL_TYPE_SET
	MYSQL_TYPE_TINY_BLOB
	MYSQL_TYPE_MEDIUM_BLOB
	MYSQL_TYPE_LONG_BLOB
	MYSQL_TYPE_BLOB
	MYSQL_TYPE_VAR_STRING
	MYSQL_TYPE_STRING
	MYSQL_TYPE_GEOMETRY
)

// Filed Flag
const (
	NOT_NULL_FLAG           = 1 << iota /* Field can't be NULL */
	PRI_KEY_FLAG                        /* Field is part of a primary key */
	UNIQUE_KEY_FLAG                     /* Field is part of a unique key */
	MULTIPLE_KEY_FLAG                   /* Field is part of a key */
	BLOB_FLAG                           /* Field is a blob */
	UNSIGNED_FLAG                       /* Field is unsigned */
	ZEROFILL_FLAG                       /* Field is zerofill */
	BINARY_FLAG                         /* Field is binary   */
	ENUM_FLAG                           /* field is an enum */
	AUTO_INCREMENT_FLAG                 /* field is a autoincrement field */
	TIMESTAMP_FLAG                      /* Field is a timestamp */
	SET_FLAG                            /* field is a set */
	NO_DEFAULT_VALUE_FLAG               /* Field doesn't have default value */
	ON_UPDATE_NOW_FLAG                  /* Field is set to NOW on UPDATE */
	NUM_FLAG                            /* Field is num (for clients) */
	PART_KEY_FLAG                       /* Intern; Part of some key */
	GROUP_FLAG                          /* Intern: Group field */
	UNIQUE_FLAG                         /* Intern: Used by sql_yacc */
	BINCMP_FLAG                         /* Intern: Used by sql_yacc */
	GET_FIXED_FIELDS_FLAG               /* Used to get fields in item tree */
	FIELD_IN_PART_FUNC_FLAG             /* Field part of partition func */
	FIELD_IN_ADD_INDEX                  /* Intern: Field used in ADD INDEX */
	FIELD_IS_RENAMED                    /* Intern: Field is being renamed */
)

// Status Flags
// See http://dev.mysql.com/doc/internals/en/status-flags.html for detail
const (
	STATUS_IN_TRANS = 1 << iota
	STATUS_AUTOCOMMIT
	STATUS_RESERVED // Not in documentation
	STATUS_MORE_RESULTS_EXISTS
	STATUS_NO_GOOD_INDEX_USED
	STATUS_NO_INDEX_USED
	STATUS_CURSOR_EXISTS
	STATUS_LAST_ROW_SENT
	STATUS_DB_DROPPED
	STATUS_NO_BACKSLASH_ESCAPES
	STATUS_METADATA_CHANGED
	STATUS_QUERY_WAS_SLOW
	STATUS_PS_OUT_PARAMS
	STATUS_IN_TRANS_READONLY
	STATUS_SESSION_STATE_CHANGED
)

// Generic Response Packets
// See http://dev.mysql.com/doc/internals/en/generic-response-packets.html for detail
const (
	OK           = 0x00
	EOF          = 0xfe
	LOCAL_INFILE = 0xfb
	ERR          = 0xff
)

const (
	MAX_PACKET_SIZE = 1<<24 - 1
	MAX_DATA_CHUNK  = 1 << 19
)

// COM_STMT_EXECUTE Response Flags
// See http://dev.mysql.com/doc/internals/en/com-stmt-execute.html for detail
const (
	CURSOR_TYPE_NO_CURSOR  = 0
	CURSOR_TYPE_READ_ONLY  = 1
	CURSOR_TYPE_FOR_UPDATE = 2
	CURSOR_TYPE_SCROLLABLE = 4
)

const (
	SERVER_MORE_RESULTS_EXISTS = 8
)
