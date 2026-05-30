package dialect_test

import (
	"testing"
	"time"

	"github.com/fedora-oss/javinizer-go/internal/database/dialect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testStart = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	testEnd   = time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	testDate  = time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
)

// ---- For() factory ----

func TestFor_KnownTypes(t *testing.T) {
	cases := []struct {
		input dialect.DBType
	}{
		{dialect.DBTypeSQLite},
		{dialect.DBTypePostgres},
		{dialect.DBTypeMySQL},
		{""}, // empty string → SQLite (backward-compat)
	}
	for _, tc := range cases {
		d, err := dialect.For(tc.input)
		require.NoError(t, err, "type=%q", tc.input)
		assert.NotNil(t, d)
	}
}

func TestFor_UnknownType_ReturnsError(t *testing.T) {
	_, err := dialect.For("mssql")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database type")
}

// ---- SQLiteDialect ----

func TestSQLiteDialect_BetweenDateTimeExpr(t *testing.T) {
	d := dialect.SQLiteDialect{}
	clause, args := d.BetweenDateTimeExpr("created_at", testStart, testEnd)

	assert.Contains(t, clause, "datetime(created_at)")
	assert.Contains(t, clause, "BETWEEN")
	assert.Contains(t, clause, "datetime(?)")
	require.Len(t, args, 2)
	assert.Equal(t, "2024-01-01 00:00:00", args[0])
	assert.Equal(t, "2024-12-31 23:59:59", args[1])
}

func TestSQLiteDialect_BeforeDateTimeExpr(t *testing.T) {
	d := dialect.SQLiteDialect{}
	clause, args := d.BeforeDateTimeExpr("created_at", testDate)

	assert.Contains(t, clause, "datetime(created_at)")
	assert.Contains(t, clause, "<")
	require.Len(t, args, 1)
	assert.Equal(t, "2024-06-15 12:00:00", args[0])
}

// ---- PostgreSQLDialect ----

func TestPostgreSQLDialect_BetweenDateTimeExpr(t *testing.T) {
	d := dialect.PostgreSQLDialect{}
	clause, args := d.BetweenDateTimeExpr("created_at", testStart, testEnd)

	assert.Contains(t, clause, "created_at")
	assert.Contains(t, clause, "BETWEEN")
	require.Len(t, args, 2)
	assert.Equal(t, testStart.UTC(), args[0])
	assert.Equal(t, testEnd.UTC(), args[1])
}

func TestPostgreSQLDialect_BeforeDateTimeExpr(t *testing.T) {
	d := dialect.PostgreSQLDialect{}
	clause, args := d.BeforeDateTimeExpr("created_at", testDate)

	assert.Contains(t, clause, "created_at <")
	require.Len(t, args, 1)
	assert.Equal(t, testDate.UTC(), args[0])
}

// ---- MySQLDialect ----

func TestMySQLDialect_BetweenDateTimeExpr(t *testing.T) {
	d := dialect.MySQLDialect{}
	clause, args := d.BetweenDateTimeExpr("created_at", testStart, testEnd)

	assert.Contains(t, clause, "created_at")
	assert.Contains(t, clause, "BETWEEN")
	require.Len(t, args, 2)
	assert.Equal(t, "2024-01-01 00:00:00", args[0])
	assert.Equal(t, "2024-12-31 23:59:59", args[1])
}

func TestMySQLDialect_BeforeDateTimeExpr(t *testing.T) {
	d := dialect.MySQLDialect{}
	clause, args := d.BeforeDateTimeExpr("created_at", testDate)

	assert.Contains(t, clause, "created_at <")
	require.Len(t, args, 1)
	assert.Equal(t, "2024-06-15 12:00:00", args[0])
}

// ---- Timezone safety ----

func TestDialects_NormalizeToUTC(t *testing.T) {
	// All dialects must produce UTC-based values regardless of input timezone.
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh") // UTC+7
	localTime := time.Date(2024, 6, 15, 19, 0, 0, 0, loc)
	expectedUTC := "2024-06-15 12:00:00"

	sqlite := dialect.SQLiteDialect{}
	_, args := sqlite.BeforeDateTimeExpr("ts", localTime)
	assert.Equal(t, expectedUTC, args[0], "SQLite must normalize to UTC")

	mysql := dialect.MySQLDialect{}
	_, mysqlArgs := mysql.BeforeDateTimeExpr("ts", localTime)
	assert.Equal(t, expectedUTC, mysqlArgs[0], "MySQL must normalize to UTC")

	pg := dialect.PostgreSQLDialect{}
	_, pgArgs := pg.BeforeDateTimeExpr("ts", localTime)
	assert.Equal(t, localTime.UTC(), pgArgs[0], "PostgreSQL must normalize to UTC time.Time")
}
