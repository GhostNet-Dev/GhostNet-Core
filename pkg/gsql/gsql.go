package gsql

import (
	"database/sql"	
)

// GSqler sql 접근을 위한 interface
type GSqler interface {
	OpenSQL()
}

// GSql sql Instance
type GSql struct {
	db *sql.DB
}

var (
	gSql GSql
)

// NewGSql sql instance를 생성한다.
func NewGSql(sqlType string) *GSql {
	switch sqlType {
		case "postgres":
		case "sqlite3":
			gSql.db = sqlite
	}

	return &gSql
}