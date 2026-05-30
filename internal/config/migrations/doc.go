package migrations

import "github.com/fedora-oss/javinizer-go/internal/config"

func init() {
	config.RegisterMigration(config.NewLegacyMigration())
}
