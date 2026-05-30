package token

import (
	"github.com/gin-gonic/gin"
	"github.com/fedora-oss/javinizer-go/internal/api/core"
)

func RegisterRoutes(protected *gin.RouterGroup, writeProtected *gin.RouterGroup, deps *core.ServerDependencies) {
	protected.GET("/tokens", listTokens(deps))
	writeProtected.POST("/tokens", createToken(deps))
	writeProtected.DELETE("/tokens/:id", revokeToken(deps))
	writeProtected.POST("/tokens/:id/regenerate", regenerateToken(deps))
}
