// Package dialect provides database-specific SQL helpers so that repository
// implementations remain dialect-agnostic.  Add a new Dialect implementation
// when a new database backend is introduced; no other package needs to change.
package dialect

import (
	"fmt"
	"time"
)

// Dialect abstracts the small set of SQL idioms that differ across database
// engines supported by the application (SQLite, PostgreSQL, MySQL).
//
// Rules for implementors:
//   - Methods must be pure functions — no I/O, no state mutation.
//   - All time arguments are assumed to be in UTC.
type Dialect interface {
	// BetweenDateTimeExpr returns the WHERE clause fragment for a half-open
	// [start, end] range on the given column using the engine's datetime type.
	//
	// Example (SQLite):
	//   "datetime(created_at) BETWEEN datetime('2024-01-01 00:00:00') AND datetime('2024-12-31 23:59:59')"
	BetweenDateTimeExpr(column string, start, end time.Time) (clause string, args []interface{})

	// BeforeDateTimeExpr returns a WHERE clause fragment for "column < date".
	BeforeDateTimeExpr(column string, date time.Time) (clause string, args []interface{})
}

// --- SQLite ---

// SQLiteDialect implements Dialect for SQLite.
// SQLite stores timestamps as TEXT; the datetime() function normalises
// different literal formats before comparison.
type SQLiteDialect struct{}

const sqliteTimeFmt = "2006-01-02 15:04:05"

func (SQLiteDialect) BetweenDateTimeExpr(column string, start, end time.Time) (string, []interface{}) {
	clause := fmt.Sprintf("datetime(%s) BETWEEN datetime(?) AND datetime(?)", column)
	args := []interface{}{start.UTC().Format(sqliteTimeFmt), end.UTC().Format(sqliteTimeFmt)}
	return clause, args
}

func (SQLiteDialect) BeforeDateTimeExpr(column string, date time.Time) (string, []interface{}) {
	clause := fmt.Sprintf("datetime(%s) < datetime(?)", column)
	args := []interface{}{date.UTC().Format(sqliteTimeFmt)}
	return clause, args
}

// --- PostgreSQL ---

// PostgreSQLDialect implements Dialect for PostgreSQL.
// PostgreSQL has a native TIMESTAMPTZ type so standard ISO-8601 strings work
// directly; no wrapper function is needed.
type PostgreSQLDialect struct{}

func (PostgreSQLDialect) BetweenDateTimeExpr(column string, start, end time.Time) (string, []interface{}) {
	clause := fmt.Sprintf("%s BETWEEN ? AND ?", column)
	args := []interface{}{start.UTC(), end.UTC()}
	return clause, args
}

func (PostgreSQLDialect) BeforeDateTimeExpr(column string, date time.Time) (string, []interface{}) {
	clause := fmt.Sprintf("%s < ?", column)
	args := []interface{}{date.UTC()}
	return clause, args
}

// --- MySQL ---

// MySQLDialect implements Dialect for MySQL / MariaDB.
// MySQL DATETIME columns store values without a timezone, so we pass UTC
// strings formatted as "YYYY-MM-DD HH:MM:SS".
type MySQLDialect struct{}

const mysqlTimeFmt = "2006-01-02 15:04:05"

func (MySQLDialect) BetweenDateTimeExpr(column string, start, end time.Time) (string, []interface{}) {
	clause := fmt.Sprintf("%s BETWEEN ? AND ?", column)
	args := []interface{}{start.UTC().Format(mysqlTimeFmt), end.UTC().Format(mysqlTimeFmt)}
	return clause, args
}

func (MySQLDialect) BeforeDateTimeExpr(column string, date time.Time) (string, []interface{}) {
	clause := fmt.Sprintf("%s < ?", column)
	args := []interface{}{date.UTC().Format(mysqlTimeFmt)}
	return clause, args
}

// --- Factory ---

// DBType enumerates the supported database backends.
type DBType string

const (
	DBTypeSQLite   DBType = "sqlite"
	DBTypePostgres DBType = "postgres"
	DBTypeMySQL    DBType = "mysql"
)

// For returns the Dialect that matches the given DBType.
// Returns an error for unknown types so callers fail fast at startup.
func For(dbType DBType) (Dialect, error) {
	switch dbType {
	case DBTypeSQLite, "": // empty defaults to SQLite for backward-compat
		return SQLiteDialect{}, nil
	case DBTypePostgres:
		return PostgreSQLDialect{}, nil
	case DBTypeMySQL:
		return MySQLDialect{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type %q: valid options are sqlite, postgres, mysql", dbType)
	}
}
