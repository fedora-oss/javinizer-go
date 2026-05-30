package organizer

import (
	"github.com/fedora-oss/javinizer-go/internal/types"
)

type OperationMode = types.OperationMode

const (
	OperationModeOrganize              = types.OperationModeOrganize
	OperationModeInPlace               = types.OperationModeInPlace
	OperationModeInPlaceNoRenameFolder = types.OperationModeInPlaceNoRenameFolder
	OperationModeMetadataArtwork       = types.OperationModeMetadataArtwork
	OperationModePreview               = types.OperationModePreview
)

var ParseOperationMode = types.ParseOperationMode
