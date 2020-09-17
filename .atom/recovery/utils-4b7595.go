package database

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

const (
	MYSQL_ER_NO_REFERENCED_ROW   uint16 = 1216
	MYSQL_ER_ROW_IS_REFERENCED   uint16 = 1217
	MYSQL_ER_ROW_IS_REFERENCED_2 uint16 = 1451
	MYSQL_ER_NO_REFERENCED_ROW_2 uint16 = 1452
)

var foreignKeyErrors map[uint16]bool = map[uint16]bool{
	MYSQL_ER_NO_REFERENCED_ROW:   true,
	MYSQL_ER_ROW_IS_REFERENCED:   true,
	MYSQL_ER_NO_REFERENCED_ROW_2: true,
	MYSQL_ER_ROW_IS_REFERENCED_2: true,
}

// IsForeignKeyError checks whether or not the passed error struct is an
// instance of MySQL ERROR 1452 or ERROR 1216
func IsForeignKeyError(err error) bool {
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if exists := foreignKeyErrors[driverErr.Number]; exists {
			return true
		}
	}
	return false
}

func uint8ToTime() (*time.Time, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return time.Parse("2006-01-02 15:04:05", string(t))
}
