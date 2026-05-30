package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fedora-oss/javinizer-go/internal/config"
	"github.com/fedora-oss/javinizer-go/internal/database/dialect"
	gormMySQL "gorm.io/driver/mysql"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the GORM database connection together with its SQL dialect.
// Repositories receive a *DB so they can build engine-correct SQL without
// importing gorm directly or branching on the Type string themselves.
type DB struct {
	*gorm.DB
	dsn     string
	Dialect dialect.Dialect // never nil after New() returns without error
}

var sqliteMemoryDSNCounter atomic.Uint64

// parseLogLevel converts a log level string to a GORM logger.LogLevel
// Normalizes input by trimming whitespace and converting to lowercase
// Returns logger.Silent for invalid values with a warning
func parseLogLevel(level string) logger.LogLevel {
	// Normalize input: trim whitespace and convert to lowercase for case-insensitive comparison
	normalized := strings.ToLower(strings.TrimSpace(level))

	switch normalized {
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	case "silent", "":
		return logger.Silent
	default:
		// Invalid log level provided - warn and default to silent
		log.Printf("Warning: invalid database log_level '%s', defaulting to 'silent'. Valid options: silent, error, warn, info\n", level)
		return logger.Silent
	}
}

// New opens a database connection for the backend specified in cfg.Database.Type.
//
// Supported backends:
//   - "sqlite" or "" (default) – uses gorm.io/driver/sqlite with CGO.
//   - "postgres"               – uses gorm.io/driver/postgres (pure-Go).
//   - "mysql"                  – uses gorm.io/driver/mysql   (pure-Go).
//
// Connection-pool settings (MaxOpenConns, MaxIdleConns, ConnMaxLifetime) are
// applied after the initial open so callers do not need to encode them in the
// DSN string.
func New(cfg *config.Config) (*DB, error) {
	dbType := dialect.DBType(strings.ToLower(strings.TrimSpace(cfg.Database.Type)))

	// Resolve the correct SQL dialect first so we can fail fast before
	// attempting any network/file connection.
	dbDialect, err := dialect.For(dbType)
	if err != nil {
		return nil, fmt.Errorf("database configuration: %w", err)
	}

	// Build the GORM dialector for the selected backend.
	var dialector gorm.Dialector
	switch dbType {
	case dialect.DBTypeSQLite, "": // empty string → sqlite (backward-compat)
		// SQLite requires CGO and a file path (or :memory:).
		dialector = sqlite.Open(normalizeSQLiteDSN(cfg.Database.DSN))

	case dialect.DBTypePostgres:
		// PostgreSQL: DSN is a libpq-style connection string.
		// Example: "host=localhost user=jav password=secret dbname=jav sslmode=disable"
		if cfg.Database.DSN == "" {
			return nil, fmt.Errorf("database type %q requires a non-empty DSN", dbType)
		}
		dialector = gormPostgres.Open(cfg.Database.DSN)

	case dialect.DBTypeMySQL:
		// MySQL/MariaDB: DSN follows go-sql-driver/mysql format.
		// Example: "jav:secret@tcp(localhost:3306)/jav?parseTime=True&loc=UTC"
		if cfg.Database.DSN == "" {
			return nil, fmt.Errorf("database type %q requires a non-empty DSN", dbType)
		}
		dialector = gormMySQL.Open(cfg.Database.DSN)

	default:
		// Should be unreachable because dialect.For() already validated the type,
		// but kept for defence-in-depth.
		return nil, fmt.Errorf("unsupported database type: %q", dbType)
	}

	// Configure GORM logger (independent from application-level logging).
	logLevel := parseLogLevel(cfg.Database.LogLevel)

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			// Always store timestamps in UTC for cross-timezone consistency.
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("open %s database: %w", dbType, err)
	}

	// Apply connection-pool settings when they are non-zero.
	// This is a no-op for SQLite in most configurations, but harmless.
	if cfg.Database.MaxOpenConns > 0 || cfg.Database.MaxIdleConns > 0 || cfg.Database.ConnMaxLifetime > 0 {
		sqlDB, err := gormDB.DB()
		if err != nil {
			return nil, fmt.Errorf("retrieve sql.DB handle for pool config: %w", err)
		}
		if cfg.Database.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		}
		if cfg.Database.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		}
		if cfg.Database.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
		}
	}

	return &DB{
		DB:      gormDB,
		dsn:     cfg.Database.DSN,
		Dialect: dbDialect,
	}, nil
}

// AutoMigrate runs startup database migrations.
//
// Kept for backward compatibility in tests and existing call sites.
// New runtime paths should call RunMigrationsOnStartup directly.
func (db *DB) AutoMigrate() error {
	return db.RunMigrationsOnStartup(context.Background())
}

// SqliteTimeFormat is used to format time.Time values for SQLite datetime comparisons.
// SQLite stores timestamps as TEXT in inconsistent formats (RFC3339 with T/Z, with
// fractional seconds, etc.) and GORM binds time.Time as "2006-01-02 15:04:05" (space,
// no TZ). Direct TEXT comparison between these formats produces wrong results because
// 'T' > ' ' and fractional seconds alter lexicographic order. Wrapping both sides in
// datetime() normalizes to a consistent format before comparison.
const SqliteTimeFormat = "2006-01-02 15:04:05"

func normalizeSQLiteDSN(dsn string) string {
	normalized := strings.ToLower(strings.TrimSpace(dsn))
	if normalized != ":memory:" {
		return dsn
	}
	// `:memory:` is scoped per SQLite connection. Goose migration checks and applies
	// can use multiple connections, so convert to a unique shared-cache memory URI.
	next := sqliteMemoryDSNCounter.Add(1)
	return fmt.Sprintf("file:javinizer_mem_%d_%d?mode=memory&cache=shared&_busy_timeout=5000", time.Now().UnixNano(), next)
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// isRecordNotFound returns true if the error is a GORM record-not-found error.
func isRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
