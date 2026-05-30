package scrape

import (
	"fmt"
	"testing"

	"github.com/fedora-oss/javinizer-go/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Register scraper defaults for NormalizeScraperConfigs
	_ "github.com/fedora-oss/javinizer-go/internal/scraper/dmm"
)

// TestApplyFlagOverrides_ActressDB tests actress-db flag overrides
func TestApplyFlagOverrides_ActressDB(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		flagVal  bool
		expected bool
	}{
		{"enable actress-db", "actress-db", true, true},
		{"disable no-actress-db", "no-actress-db", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			cfg := config.DefaultConfig()
			cfg.Metadata.ActressDatabase.Enabled = !tt.expected

			err := cmd.Flags().Set(tt.flagName, fmt.Sprintf("%t", tt.flagVal))
			require.NoError(t, err)

			ApplyFlagOverrides(cmd, cfg)

			assert.Equal(t, tt.expected, cfg.Metadata.ActressDatabase.Enabled,
				"Flag %s should set ActressDatabase.Enabled to %t", tt.flagName, tt.expected)
		})
	}
}

// TestApplyFlagOverrides_GenreReplacement tests genre-replacement flag overrides
func TestApplyFlagOverrides_GenreReplacement(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		flagVal  bool
		expected bool
	}{
		{"enable genre-replacement", "genre-replacement", true, true},
		{"disable no-genre-replacement", "no-genre-replacement", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			cfg := config.DefaultConfig()
			cfg.Metadata.GenreReplacement.Enabled = !tt.expected

			err := cmd.Flags().Set(tt.flagName, fmt.Sprintf("%t", tt.flagVal))
			require.NoError(t, err)

			ApplyFlagOverrides(cmd, cfg)

			assert.Equal(t, tt.expected, cfg.Metadata.GenreReplacement.Enabled,
				"Flag %s should set GenreReplacement.Enabled to %t", tt.flagName, tt.expected)
		})
	}
}

// Note: DMM-specific flag tests (scrape-actress, browser, browser-timeout) have been
// removed since ScraperSettings.Extra was removed as part of the plugin system migration.
// These tests will be restored once DMM-specific CLI flags are reimplemented via
// the concrete DMMConfig type.
