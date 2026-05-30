package organizer

import (
	"github.com/fedora-oss/javinizer-go/internal/matcher"
	"github.com/fedora-oss/javinizer-go/internal/models"
)

type OperationStrategy interface {
	Plan(match matcher.MatchResult, movie *models.Movie, destDir string, forceUpdate bool) (*OrganizePlan, error)
	Execute(plan *OrganizePlan) (*OrganizeResult, error)
}
