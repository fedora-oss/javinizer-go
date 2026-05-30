package system

import (
	"testing"

	"github.com/fedora-oss/javinizer-go/internal/api/core"
	"github.com/fedora-oss/javinizer-go/internal/api/testkit"
	"github.com/fedora-oss/javinizer-go/internal/config"
)

func createTestDeps(t *testing.T, cfg *config.Config, configFile string) *core.ServerDependencies {
	return testkit.CreateTestDeps(t, cfg, configFile)
}
