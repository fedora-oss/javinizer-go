//go:build !cgo

package database

import (
	"strings"
)

func isLocked(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "database is locked") || strings.Contains(err.Error(), "database table is locked"))
}
