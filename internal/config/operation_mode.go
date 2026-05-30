package config

import (
	"github.com/fedora-oss/javinizer-go/internal/types"
)

func GetOperationMode(mode string) types.OperationMode {
	if mode == "" {
		return types.OperationModeOrganize
	}
	parsed, err := types.ParseOperationMode(mode)
	if err != nil {
		return types.OperationModeOrganize
	}
	return parsed
}

func (o *OutputConfig) GetOperationMode() types.OperationMode {
	return GetOperationMode(string(o.OperationMode))
}
