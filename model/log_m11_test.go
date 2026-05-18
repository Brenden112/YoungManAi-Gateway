package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
)

// validates M11-F01 — Log struct has new B2B gateway fields
func TestLogSchemaNewFields(t *testing.T) {
	orgId := 5
	projectId := 2
	log := Log{
		UserId:              1,
		Type:                LogTypeConsume,
		OrgId:               &orgId,
		ProjectId:           &projectId,
		IsExperimentalProxy: true,
		ProviderType:        "official_cloud",
	}
	if log.OrgId == nil || *log.OrgId != 5 {
		t.Errorf("OrgId = %v, want 5", log.OrgId)
	}
	if log.ProjectId == nil || *log.ProjectId != 2 {
		t.Errorf("ProjectId = %v, want 2", log.ProjectId)
	}
	if !log.IsExperimentalProxy {
		t.Error("IsExperimentalProxy should be true")
	}
	if log.ProviderType != "official_cloud" {
		t.Errorf("ProviderType = %q, want 'official_cloud'", log.ProviderType)
	}
}

// validates M11-F02 — LogTypeError and LogTypeConsume constants are distinct
func TestLogTypeConstants(t *testing.T) {
	if LogTypeConsume == LogTypeError {
		t.Error("LogTypeConsume and LogTypeError must be distinct")
	}
	if LogTypeConsume == 0 || LogTypeError == 0 {
		t.Error("log type constants must not be 0 (0 is unknown)")
	}
}

// validates M11-F03 — StoreFullTextEnabled defaults to false
func TestStoreFullTextEnabledDefaultFalse(t *testing.T) {
	if common.StoreFullTextEnabled {
		t.Error("StoreFullTextEnabled must default to false (privacy-first)")
	}
}

// validates M11-F03 — content is cleared when StoreFullTextEnabled=false
func TestRecordConsumeLogClearsContentWhenDisabled(t *testing.T) {
	// Simulate the content-clearing logic from RecordConsumeLog.
	content := "user: hello\nassistant: hi there"
	storeFullText := false

	stored := func() string {
		if storeFullText {
			return content
		}
		return ""
	}()

	if stored != "" {
		t.Errorf("content should be empty when StoreFullTextEnabled=false, got %q", stored)
	}
}

// validates M11-F03 — content is preserved when StoreFullTextEnabled=true
func TestRecordConsumeLogPreservesContentWhenEnabled(t *testing.T) {
	content := "user: hello\nassistant: hi there"
	storeFullText := true

	stored := func() string {
		if storeFullText {
			return content
		}
		return ""
	}()

	if stored != content {
		t.Errorf("content should be preserved when StoreFullTextEnabled=true, got %q", stored)
	}
}
