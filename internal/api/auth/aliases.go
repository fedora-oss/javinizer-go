package auth

import (
	"github.com/fedora-oss/javinizer-go/internal/api/contracts"
	"github.com/fedora-oss/javinizer-go/internal/api/core"
)

type ServerDependencies = core.ServerDependencies

type ErrorResponse = contracts.ErrorResponse
type AuthStatusResponse = contracts.AuthStatusResponse
type AuthCredentialsRequest = contracts.AuthCredentialsRequest
