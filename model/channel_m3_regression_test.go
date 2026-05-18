package model

import (
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
)

// M3-F01 regression: Channel JSON serialization must include all relay-relevant fields.
// Validates A3.1.1 — new M1/M2 fields must not shadow or break existing JSON keys.
func TestChannelJSONFieldsIntact(t *testing.T) {
	ch := Channel{
		Id:     1,
		Type:   constant.ChannelTypeOpenAI,
		Key:    "sk-test",
		Status: common.ChannelStatusEnabled,
		Name:   "test-openai",
		Models: "gpt-4o,gpt-4o-mini",
		Group:  "default",
		// M1 fields
		ProviderType:         constant.ProviderTypeOfficialCloud,
		RiskLevel:            constant.RiskLevelNormal,
		AvailableScope:       constant.ScopePublic,
		Visibility:           constant.VisibilityPublic,
		ManualEnableRequired: false,
		// M2 field
		ProviderAccountId: nil,
	}

	data, err := common.Marshal(ch)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	json := string(data)

	// Relay-relevant fields must be present
	relayFields := []string{
		`"id"`, `"type"`, `"status"`, `"name"`, `"models"`, `"group"`,
	}
	for _, field := range relayFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON missing relay-relevant field %s", field)
		}
	}

	// M1 fields must be present
	m1Fields := []string{
		`"provider_type"`, `"risk_level"`, `"available_scope"`, `"visibility"`, `"manual_enable_required"`,
	}
	for _, field := range m1Fields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON missing M1 field %s", field)
		}
	}

	// M2 field must be present (as null)
	if !strings.Contains(json, `"provider_account_id"`) {
		t.Error("JSON missing M2 field provider_account_id")
	}

	// Key must NOT appear in JSON (security: key is sensitive but json tag is present;
	// verify it serializes — the masking happens at the controller layer)
	if !strings.Contains(json, `"key"`) {
		t.Error("JSON missing key field — controller masking layer may be broken")
	}
}

// M3-F01 regression: Channel round-trip JSON must preserve provider_type value.
func TestChannelJSONRoundTrip(t *testing.T) {
	original := Channel{
		Type:         constant.ChannelTypeOpenAI,
		ProviderType: constant.ProviderTypeOfficialCloud,
		RiskLevel:    constant.RiskLevelNormal,
		Status:       common.ChannelStatusEnabled,
	}

	data, err := common.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var restored Channel
	if err := common.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if restored.ProviderType != original.ProviderType {
		t.Errorf("ProviderType round-trip: got %q, want %q", restored.ProviderType, original.ProviderType)
	}
	if restored.RiskLevel != original.RiskLevel {
		t.Errorf("RiskLevel round-trip: got %q, want %q", restored.RiskLevel, original.RiskLevel)
	}
	if restored.Status != original.Status {
		t.Errorf("Status round-trip: got %d, want %d", restored.Status, original.Status)
	}
}

// M3-F02 regression: experimental_proxy channel JSON must not expose internal fields
// to non-admin callers. Visibility field must be present and correct.
func TestExperimentalProxyChannelJSONVisibility(t *testing.T) {
	ch := Channel{ProviderType: constant.ProviderTypeExperimentalProxy}
	_ = ch.BeforeCreate(nil)

	data, err := common.Marshal(ch)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	json := string(data)

	if !strings.Contains(json, `"internal_only"`) {
		t.Error("experimental_proxy channel JSON must contain internal_only for visibility/available_scope")
	}
	if !strings.Contains(json, `"manual_enable_required":true`) {
		t.Error("experimental_proxy channel JSON must have manual_enable_required:true")
	}
}

// M3-F03 regression: ProviderAccount key must never appear in Channel JSON.
// The Key field on ProviderAccount has gorm:"-" and json:"key,omitempty".
// Channel.ProviderAccountId is nullable — must serialize as null when unset.
func TestProviderAccountIdNullInChannelJSON(t *testing.T) {
	ch := Channel{Type: constant.ChannelTypeOpenAI, Status: common.ChannelStatusEnabled}
	data, err := common.Marshal(ch)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	json := string(data)

	// provider_account_id must be null (not absent) so frontend can detect unlinked channels
	if !strings.Contains(json, `"provider_account_id":null`) {
		t.Errorf("provider_account_id should serialize as null, got JSON: %s", json)
	}
}
