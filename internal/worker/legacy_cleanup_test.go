package worker

import (
	"reflect"
	"testing"

	"github.com/fedora-oss/javinizer-go/internal/api/contracts"
	"github.com/stretchr/testify/assert"
)

func TestBatchJob_HasNoLegacyOverrideFields(t *testing.T) {
	jobType := reflect.TypeOf(BatchJob{})

	_, hasMoveToFolderOverride := jobType.FieldByName("MoveToFolderOverride")
	assert.False(t, hasMoveToFolderOverride,
		"BatchJob should not have MoveToFolderOverride field — removed in LGCY-01")

	_, hasRenameFolderInPlaceOverride := jobType.FieldByName("RenameFolderInPlaceOverride")
	assert.False(t, hasRenameFolderInPlaceOverride,
		"BatchJob should not have RenameFolderInPlaceOverride field — removed in LGCY-01")

	_, hasOperationModeOverride := jobType.FieldByName("OperationModeOverride")
	assert.True(t, hasOperationModeOverride,
		"BatchJob must have OperationModeOverride field — the replacement for legacy boolean overrides")
}

func TestBatchJobSlim_HasNoLegacyOverrideFields(t *testing.T) {
	slimType := reflect.TypeOf(BatchJobSlim{})

	_, hasMoveToFolderOverride := slimType.FieldByName("MoveToFolderOverride")
	assert.False(t, hasMoveToFolderOverride,
		"BatchJobSlim should not have MoveToFolderOverride field — removed in LGCY-01")

	_, hasRenameFolderInPlaceOverride := slimType.FieldByName("RenameFolderInPlaceOverride")
	assert.False(t, hasRenameFolderInPlaceOverride,
		"BatchJobSlim should not have RenameFolderInPlaceOverride field — removed in LGCY-01")

	_, hasOperationModeOverride := slimType.FieldByName("OperationModeOverride")
	assert.True(t, hasOperationModeOverride,
		"BatchJobSlim must have OperationModeOverride field — the replacement for legacy boolean overrides")
}

func TestBatchScrapeRequest_HasNoLegacyBooleanFields(t *testing.T) {
	reqType := reflect.TypeOf(contracts.BatchScrapeRequest{})

	_, hasMoveToFolder := reqType.FieldByName("MoveToFolder")
	assert.False(t, hasMoveToFolder,
		"BatchScrapeRequest should not have MoveToFolder field — removed in LGCY-01")

	_, hasRenameFolderInPlace := reqType.FieldByName("RenameFolderInPlace")
	assert.False(t, hasRenameFolderInPlace,
		"BatchScrapeRequest should not have RenameFolderInPlace field — removed in LGCY-01")

	_, hasOperationMode := reqType.FieldByName("OperationMode")
	assert.True(t, hasOperationMode,
		"BatchScrapeRequest must have OperationMode field — the replacement for legacy boolean fields")
}

func TestNFOComparisonRequest_HasNoMergeStrategy(t *testing.T) {
	reqType := reflect.TypeOf(contracts.NFOComparisonRequest{})

	_, hasMergeStrategy := reqType.FieldByName("MergeStrategy")
	assert.False(t, hasMergeStrategy,
		"NFOComparisonRequest should not have MergeStrategy field — removed in DEAD-03")

	_, hasPreset := reqType.FieldByName("Preset")
	assert.True(t, hasPreset,
		"NFOComparisonRequest must have Preset field — replacement for MergeStrategy")

	_, hasScalarStrategy := reqType.FieldByName("ScalarStrategy")
	assert.True(t, hasScalarStrategy,
		"NFOComparisonRequest must have ScalarStrategy field — replacement for MergeStrategy")

	_, hasArrayStrategy := reqType.FieldByName("ArrayStrategy")
	assert.True(t, hasArrayStrategy,
		"NFOComparisonRequest must have ArrayStrategy field — replacement for MergeStrategy")
}
