package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
)

// validates A7.3.1 — experimental_proxy BeforeCreate sets high-risk defaults
func TestChannelBeforeCreateExperimentalProxy(t *testing.T) {
	ch := &Channel{ProviderType: constant.ProviderTypeExperimentalProxy}
	if err := ch.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate returned error: %v", err)
	}
	if ch.RiskLevel != constant.RiskLevelHigh {
		t.Errorf("RiskLevel = %q, want %q", ch.RiskLevel, constant.RiskLevelHigh)
	}
	if ch.AvailableScope != constant.ScopeInternalOnly {
		t.Errorf("AvailableScope = %q, want %q", ch.AvailableScope, constant.ScopeInternalOnly)
	}
	if ch.Visibility != constant.VisibilityInternalOnly {
		t.Errorf("Visibility = %q, want %q", ch.Visibility, constant.VisibilityInternalOnly)
	}
	if !ch.ManualEnableRequired {
		t.Error("ManualEnableRequired = false, want true for experimental_proxy")
	}
}

// validates A1.1.3 — empty ProviderType defaults to official_cloud with normal-risk defaults
func TestChannelBeforeCreateOfficialCloud(t *testing.T) {
	ch := &Channel{}
	if err := ch.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate returned error: %v", err)
	}
	if ch.ProviderType != constant.ProviderTypeOfficialCloud {
		t.Errorf("ProviderType = %q, want %q", ch.ProviderType, constant.ProviderTypeOfficialCloud)
	}
	if ch.RiskLevel != constant.RiskLevelNormal {
		t.Errorf("RiskLevel = %q, want %q", ch.RiskLevel, constant.RiskLevelNormal)
	}
	if ch.AvailableScope != constant.ScopePublic {
		t.Errorf("AvailableScope = %q, want %q", ch.AvailableScope, constant.ScopePublic)
	}
	if ch.Visibility != constant.VisibilityPublic {
		t.Errorf("Visibility = %q, want %q", ch.Visibility, constant.VisibilityPublic)
	}
	if ch.ManualEnableRequired {
		t.Error("ManualEnableRequired = true, want false for official_cloud")
	}
}

// validates security invariant — ManualEnableRequired always forced true for experimental_proxy
func TestChannelBeforeCreateExperimentalProxyManualEnableForced(t *testing.T) {
	ch := &Channel{
		ProviderType:         constant.ProviderTypeExperimentalProxy,
		ManualEnableRequired: false, // caller tries to set false — hook must override
	}
	_ = ch.BeforeCreate(nil)
	if !ch.ManualEnableRequired {
		t.Error("ManualEnableRequired must always be true for experimental_proxy regardless of caller input")
	}
}

// validates that explicit non-empty values are preserved by BeforeCreate
func TestChannelBeforeCreatePreservesExplicitValues(t *testing.T) {
	ch := &Channel{
		ProviderType:   constant.ProviderTypeAggregator,
		RiskLevel:      constant.RiskLevelHigh, // explicitly set high for an aggregator
		AvailableScope: constant.ScopePublic,
		Visibility:     constant.VisibilityPublic,
	}
	_ = ch.BeforeCreate(nil)
	if ch.RiskLevel != constant.RiskLevelHigh {
		t.Errorf("BeforeCreate overwrote explicit RiskLevel: got %q", ch.RiskLevel)
	}
	if ch.ProviderType != constant.ProviderTypeAggregator {
		t.Errorf("BeforeCreate overwrote explicit ProviderType: got %q", ch.ProviderType)
	}
}

// validates A7.3.1, A7.3.2 — experimental_proxy defaults to Status=ManuallyDisabled
func TestChannelBeforeCreateExperimentalProxyStatusDisabled(t *testing.T) {
	ch := &Channel{ProviderType: constant.ProviderTypeExperimentalProxy}
	_ = ch.BeforeCreate(nil)
	if ch.Status != common.ChannelStatusManuallyDisabled {
		t.Errorf("Status = %d, want %d (ChannelStatusManuallyDisabled)", ch.Status, common.ChannelStatusManuallyDisabled)
	}
}

// validates that explicit Status is preserved even for experimental_proxy
func TestChannelBeforeCreateExperimentalProxyExplicitStatusPreserved(t *testing.T) {
	ch := &Channel{
		ProviderType: constant.ProviderTypeExperimentalProxy,
		Status:       common.ChannelStatusEnabled, // admin explicitly enables at creation
	}
	_ = ch.BeforeCreate(nil)
	if ch.Status != common.ChannelStatusEnabled {
		t.Errorf("BeforeCreate overwrote explicit Status: got %d, want %d", ch.Status, common.ChannelStatusEnabled)
	}
}

// validates that non-experimental channels are NOT defaulted to disabled
func TestChannelBeforeCreateOfficialCloudStatusNotDisabled(t *testing.T) {
	ch := &Channel{ProviderType: constant.ProviderTypeOfficialCloud}
	_ = ch.BeforeCreate(nil)
	if ch.Status == common.ChannelStatusManuallyDisabled {
		t.Error("official_cloud channel should not be defaulted to disabled status")
	}
}

// validates A2.2.2 — ProviderAccountId is nil by default (backward-compatible)
func TestChannelProviderAccountIdNilByDefault(t *testing.T) {
	ch := &Channel{}
	if ch.ProviderAccountId != nil {
		t.Errorf("ProviderAccountId should be nil by default, got %v", *ch.ProviderAccountId)
	}
}

// validates GetProviderAccount returns (nil, nil) when ProviderAccountId is nil
func TestChannelGetProviderAccountNilId(t *testing.T) {
	ch := &Channel{}
	pa, err := ch.GetProviderAccount()
	if err != nil {
		t.Fatalf("GetProviderAccount error: %v", err)
	}
	if pa != nil {
		t.Error("GetProviderAccount should return nil when ProviderAccountId is nil")
	}
}

// validates ProviderAccountId pointer can be set and read back
func TestChannelProviderAccountIdCanBeSet(t *testing.T) {
	id := 42
	ch := &Channel{ProviderAccountId: &id}
	if ch.ProviderAccountId == nil {
		t.Fatal("ProviderAccountId should not be nil after being set")
	}
	if *ch.ProviderAccountId != 42 {
		t.Errorf("ProviderAccountId = %d, want 42", *ch.ProviderAccountId)
	}
}

// validates A2.2.2 — legacy channels (ProviderAccountId=nil) are unaffected by BeforeCreate.
// Relay-relevant fields Key, Status, Type must not be modified by the hook.
// This is the core backward-compatibility regression test for M2-F03.
func TestChannelLegacyCompatibilityNilProviderAccountId(t *testing.T) {
	ch := &Channel{
		Type:   1, // ChannelTypeOpenAI — legacy channel
		Key:    "sk-legacy-upstream-key",
		Status: common.ChannelStatusEnabled,
		Name:   "legacy-openai-channel",
		// ProviderAccountId intentionally omitted — nil by default
	}

	if err := ch.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate error on legacy channel: %v", err)
	}

	// ProviderAccountId must remain nil — BeforeCreate must not touch it
	if ch.ProviderAccountId != nil {
		t.Errorf("BeforeCreate set ProviderAccountId on legacy channel: got %d", *ch.ProviderAccountId)
	}
	// Relay path uses channel.Key directly — must be preserved
	if ch.Key != "sk-legacy-upstream-key" {
		t.Errorf("BeforeCreate modified Key: got %q", ch.Key)
	}
	// Status must be preserved for non-experimental channels
	if ch.Status != common.ChannelStatusEnabled {
		t.Errorf("BeforeCreate modified Status: got %d, want %d", ch.Status, common.ChannelStatusEnabled)
	}
	// Type must be preserved
	if ch.Type != 1 {
		t.Errorf("BeforeCreate modified Type: got %d, want 1", ch.Type)
	}
}
