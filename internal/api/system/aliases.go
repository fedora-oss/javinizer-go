package system

import (
	"github.com/fedora-oss/javinizer-go/internal/api/contracts"
	"github.com/fedora-oss/javinizer-go/internal/api/core"
)

type ServerDependencies = core.ServerDependencies

type ErrorResponse = contracts.ErrorResponse
type HealthResponse = contracts.HealthResponse
type ScraperOption = contracts.ScraperOption
type ScraperChoice = contracts.ScraperChoice
type ScraperInfo = contracts.ScraperInfo
type AvailableScrapersResponse = contracts.AvailableScrapersResponse
type ProxyTestRequest = contracts.ProxyTestRequest
type ProxyTestResponse = contracts.ProxyTestResponse
type UpdateConfigRequest = contracts.UpdateConfigRequest
type TranslationModelsRequest = contracts.TranslationModelsRequest
type TranslationModelsResponse = contracts.TranslationModelsResponse
