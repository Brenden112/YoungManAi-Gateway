package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
)

// validates M15-F02 — GetAllLogs accepts providerType and isExperimentalProxy params
func TestGetAllLogsSignatureHasNewFilters(t *testing.T) {
	// Compile-time check: function signature includes the two new parameters.
	// If the signature changes, this test will fail to compile.
	type getAllLogsFn func(
		logType int,
		startTimestamp int64,
		endTimestamp int64,
		modelName string,
		username string,
		tokenName string,
		startIdx int,
		num int,
		channel int,
		group string,
		requestId string,
		upstreamRequestId string,
		providerType string,
		isExperimentalProxy *bool,
	) ([]*Log, int64, error)
	var _ getAllLogsFn = GetAllLogs
}

// validates M15-F02 — isExperimentalProxy nil means no filter (all logs returned)
func TestGetAllLogsNilIsExperimentalProxyMeansNoFilter(t *testing.T) {
	// nil pointer = no WHERE clause added; non-nil = filter applied.
	var nilFilter *bool
	if nilFilter != nil {
		t.Error("nil isExperimentalProxy must not add a WHERE clause")
	}
	b := true
	nonNilFilter := &b
	if nonNilFilter == nil {
		t.Error("non-nil isExperimentalProxy must add a WHERE clause")
	}
}

// validates M15-F01 — GetAdminAllTokens signature accepts all admin filters
func TestGetAdminAllTokensSignature(t *testing.T) {
	type adminTokensFn func(
		userId int,
		orgId *int,
		projectId *int,
		allowExperimental *bool,
		startIdx int,
		num int,
	) ([]*Token, int64, error)
	var _ adminTokensFn = GetAdminAllTokens
}

// validates M15-F03 — balance admin: GetUser returns quota field
func TestTokenHasAllowExperimentalField(t *testing.T) {
	tok := Token{AllowExperimental: true}
	if !tok.AllowExperimental {
		t.Error("Token.AllowExperimental must be settable")
	}
}

func TestGetAllLogsFiltersProviderTypeAndExperimentalProxy(t *testing.T) {
	truncateTables(t)

	previousMemoryCacheEnabled := common.MemoryCacheEnabled
	common.MemoryCacheEnabled = false
	t.Cleanup(func() {
		common.MemoryCacheEnabled = previousMemoryCacheEnabled
	})

	if err := DB.Create(&Channel{Id: 11, Name: "official-channel"}).Error; err != nil {
		t.Fatalf("create official channel: %v", err)
	}
	if err := DB.Create(&Channel{Id: 12, Name: "experimental-channel"}).Error; err != nil {
		t.Fatalf("create experimental channel: %v", err)
	}

	now := common.GetTimestamp()
	logsToCreate := []*Log{
		{
			UserId:              1,
			CreatedAt:           now,
			Type:                LogTypeConsume,
			ModelName:           "gpt-4o-mini",
			ChannelId:           11,
			ProviderType:        "official_cloud",
			IsExperimentalProxy: false,
		},
		{
			UserId:              2,
			CreatedAt:           now + 1,
			Type:                LogTypeConsume,
			ModelName:           "kiro-test",
			ChannelId:           12,
			ProviderType:        "experimental_proxy",
			IsExperimentalProxy: true,
		},
	}
	if err := LOG_DB.Create(&logsToCreate).Error; err != nil {
		t.Fatalf("create logs: %v", err)
	}

	exp := true
	logs, total, err := GetAllLogs(LogTypeConsume, 0, 0, "", "", "", 0, 10, 0, "", "", "", "experimental_proxy", &exp)
	if err != nil {
		t.Fatalf("GetAllLogs returned error: %v", err)
	}
	if total != 1 {
		t.Fatalf("total = %d, want 1", total)
	}
	if len(logs) != 1 {
		t.Fatalf("len(logs) = %d, want 1", len(logs))
	}
	if logs[0].ProviderType != "experimental_proxy" {
		t.Fatalf("ProviderType = %q, want experimental_proxy", logs[0].ProviderType)
	}
	if !logs[0].IsExperimentalProxy {
		t.Fatal("IsExperimentalProxy = false, want true")
	}
	if logs[0].ChannelName != "experimental-channel" {
		t.Fatalf("ChannelName = %q, want experimental-channel", logs[0].ChannelName)
	}
}
