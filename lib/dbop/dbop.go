package dbop

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/patrick/test-grpc-docker-gitactions/lib/env"

	//used by gorm for ms sql connection
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	msdb *gorm.DB
	once sync.Once
)

const (
	ForShare  = "FOR SHARE"
	ForUpdate = "FOR UPDATE"
)

// MSDB is a singleton wrapper to the MSSQL gorm database object.
func MSDB() *gorm.DB {
	var err error
	// var connectionString string
	if msdb == nil {
		msdb, err = NewMSDB(env.DBConnectionString)
		if err != nil {
			panic(err)
		}
	}

	return msdb
}

func getEnvInt(key string) *int {
	val := GetValue(key)
	n, err := strconv.Atoi(val)
	if err != nil {
		return nil
	}
	return &n
}

func CloseDB() {
	if msdb != nil {
		msdb.Close()
	}
}

func getEnvBool(key string) *bool {
	val := GetValue(key)
	b, err := strconv.ParseBool(val)
	if err != nil {
		return nil
	}
	return &b
}

type DBOptions struct {
	Host         *string
	Port         *int
	User         *string
	Password     *string
	Database     *string
	SslMode      *string
	LogMode      *bool
	MaxOpenConns *int
	MaxIdleConns *int
}

//NewMSDB MS SQL database connection
func NewMSDB(connectionString string) (dbT *gorm.DB, err error) {
	dbT, err = gorm.Open("mssql", connectionString)
	if err != nil {
		return nil, err
	}

	// default = 20 (Go's default is 0 == unlimited)
	maxOpenConns := 20
	if maxOpenConns == 0 {
		maxOpenConns = 20
	}
	dbT.DB().SetMaxOpenConns(maxOpenConns)

	maxIdleConns := 39
	if maxIdleConns != 0 {
		dbT.DB().SetMaxIdleConns(maxIdleConns)
	}

	// so it doesn't reuse stale connections
	dbT.DB().SetConnMaxLifetime(30 * time.Minute)
	dbT.LogMode(true)

	return dbT, nil
}

// Reconnect pings the database to re-establish
// a connection.
func Reconnect() error {
	if msdb == nil {
		return errors.Errorf("db is nil")
	}

	return msdb.DB().Ping()
}

// IsConnectionError returns true if the supplied error
// is a connection related error based on PostgreSQL
// connection exceptions class. See:
// http://www.postgresql.org/docs/9.4/static/errcodes-appendix.html#ERRCODES-TABLE
// for details.
func IsConnectionError(err error) bool {
	return pqErrorCode(err) == "08"
}

func InsufficientResources(err error) bool {
	return pqErrorCode(err) == "53"
}

func pqErrorCode(err error) pq.ErrorCode {
	if err != nil {
		pqErr, ok := err.(*pq.Error)

		if ok {
			return pqErr.Code[0:2]
		}
	}
	return ""
}

// IsSerializabilityError returns true if the supplied error
// is due to a serializability failure in the DB.
func IsSerializabilityError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "could not serialize access due to concurrent update")
}

// Serializable begins a transaction with isolation level
// set to SERIALIZABLE.
func Serializable() *gorm.DB {
	return MSDB().Begin().Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
}

// RepeatableRead begins a transaction with isolation level
// set to REPEATABLE READ.
func RepeatableRead() *gorm.DB {
	return MSDB().Begin().Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
}

// ReadCommitted begins a transaction with isolation level
// set to READ COMMITTED.
func ReadCommitted() *gorm.DB {
	return MSDB().Begin().Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
}

// ReadUncomitted begins a transaction with isolation level
// set to READ UNCOMMITTED.
func ReadUncomitted() *gorm.DB {
	return MSDB().Begin().Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
}

// Begin a transaction.
func Begin() *gorm.DB {
	return MSDB().Begin()
}

var dVal sync.Map

// GetValue by key
func GetValue(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if v, _ := dVal.Load(key); v != nil {
			return v.(string)
		} else {
			return ""
		}
	}
	return value
}

// GetBool returns a true if the environment variable is
// "true", "t", "yes" or "y" with case insensitive.
func GetBool(key string, defaultValue bool) bool {
	strval := GetValue(key)
	if strval == "" {
		return defaultValue
	} else if strings.EqualFold(strval, "true") ||
		strings.EqualFold(strval, "t") ||
		strings.EqualFold(strval, "yes") ||
		strings.EqualFold(strval, "1") ||
		strings.EqualFold(strval, "y") {
		return true
	}

	return false
}

// GetInt returns the integer value stored in the enviroment variable.
// When the key is missing, or it cannot be parsed to int, it returns 0
// this function DOES NOT return a default value when a missing key, or non-int value is encountered!
func GetInt(key string) int {
	v, err := strconv.Atoi(GetValue(key))
	if err != nil {
		return 0
	}

	return v
}
