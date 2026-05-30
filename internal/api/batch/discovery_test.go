package batch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fedora-oss/javinizer-go/internal/config"
	"github.com/fedora-oss/javinizer-go/internal/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDiscoveryMatcher(t *testing.T) *matcher.Matcher {
	t.Helper()

	m, err := matcher.NewMatcher(&config.MatchingConfig{
		Extensions:   []string{".mp4", ".mkv"},
		RegexPattern: `(?i)([a-z]{2,10}-?\d{2,5}[a-z]?)`,
		RegexEnabled: true,
	})
	require.NoError(t, err)
	return m
}

func TestDiscoverSiblingParts_AdditionalScenarios(t *testing.T) {
	t.Run("empty input returns empty slice", func(t *testing.T) {
		cfg := config.DefaultConfig()
		files, metadata := discoverSiblingPartsWithMetadata(nil, newDiscoveryMatcher(t), cfg)
		assert.Empty(t, files)
		assert.Nil(t, metadata)
	})

	t.Run("non multipart input stays unchanged", func(t *testing.T) {
		cfg := config.DefaultConfig()
		dir := t.TempDir()
		cfg.API.Security.AllowedDirectories = []string{dir}

		filePath := filepath.Join(dir, "IPX-123.mp4")
		require.NoError(t, os.WriteFile(filePath, []byte("video"), 0o644))

		got, metadata := discoverSiblingPartsWithMetadata([]string{filePath}, newDiscoveryMatcher(t), cfg)
		assert.Equal(t, []string{filePath}, got)
		assert.NotNil(t, metadata)
		// Non-multipart files should still have metadata
		assert.Contains(t, metadata, filePath)
		assert.False(t, metadata[filePath].IsMultiPart)
	})

	t.Run("discovers missing multipart siblings from allowed directory", func(t *testing.T) {
		cfg := config.DefaultConfig()
		dir := t.TempDir()
		cfg.API.Security.AllowedDirectories = []string{dir}
		cfg.API.Security.MaxFilesPerScan = 100

		partA := filepath.Join(dir, "IPX-535-pt1.mp4")
		partB := filepath.Join(dir, "IPX-535-pt2.mp4")
		other := filepath.Join(dir, "IPX-999.mp4")
		require.NoError(t, os.WriteFile(partA, []byte("a"), 0o644))
		require.NoError(t, os.WriteFile(partB, []byte("b"), 0o644))
		require.NoError(t, os.WriteFile(other, []byte("c"), 0o644))

		got, metadata := discoverSiblingPartsWithMetadata([]string{partA}, newDiscoveryMatcher(t), cfg)
		assert.ElementsMatch(t, []string{partA, partB}, got)
		// Both parts should have multipart metadata
		assert.True(t, metadata[partA].IsMultiPart)
		assert.Equal(t, 1, metadata[partA].PartNumber)
		assert.True(t, metadata[partB].IsMultiPart)
		assert.Equal(t, 2, metadata[partB].PartNumber)
	})

	t.Run("disallowed directory skips sibling discovery", func(t *testing.T) {
		cfg := config.DefaultConfig()
		dir := t.TempDir()
		cfg.API.Security.AllowedDirectories = []string{filepath.Join(dir, "elsewhere")}

		partA := filepath.Join(dir, "IPX-535-pt1.mp4")
		partB := filepath.Join(dir, "IPX-535-pt2.mp4")
		require.NoError(t, os.WriteFile(partA, []byte("a"), 0o644))
		require.NoError(t, os.WriteFile(partB, []byte("b"), 0o644))

		got, metadata := discoverSiblingPartsWithMetadata([]string{partA}, newDiscoveryMatcher(t), cfg)
		assert.Equal(t, []string{partA}, got)
		// Since directory is disallowed, only the original file is returned
		// but it still has metadata
		assert.Contains(t, metadata, partA)
	})
}
