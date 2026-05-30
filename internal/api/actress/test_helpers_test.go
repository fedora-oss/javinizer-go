package actress

import (
	"github.com/fedora-oss/javinizer-go/internal/api/testkit"
	"github.com/fedora-oss/javinizer-go/internal/database"
)

func newMockActressRepo() *database.ActressRepository {
	return testkit.NewMockActressRepo()
}
