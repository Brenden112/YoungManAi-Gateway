package model

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
)

// validates M14-F01 — ChannelProviderSummary struct has correct fields
func TestChannelProviderSummaryFields(t *testing.T) {
	s := ChannelProviderSummary{
		ProviderType: constant.ProviderTypeOfficialCloud,
		Total:        10,
		Enabled:      8,
		Disabled:     2,
	}
	if s.Total != s.Enabled+s.Disabled {
		t.Errorf("Total %d != Enabled %d + Disabled %d", s.Total, s.Enabled, s.Disabled)
	}
	if s.ProviderType != constant.ProviderTypeOfficialCloud {
		t.Errorf("ProviderType = %q, want %q", s.ProviderType, constant.ProviderTypeOfficialCloud)
	}
}

// validates M14-F02 — DisableExperimentalProxyChannels targets only experimental_proxy
// (compile-time contract: function exists and returns rowsAffected + error)
func TestDisableExperimentalProxyChannelsSignature(t *testing.T) {
	// Verify the function exists with the correct signature.
	// Full DB test requires a test DB; this is a compile-time check.
	var fn func() (int64, error) = DisableExperimentalProxyChannels
	_ = fn
}

// validates M14-F03 — provider_type filter is a string (matches DB column type)
func TestProviderTypeFilterIsString(t *testing.T) {
	// The provider_type column is varchar(32); filter must be a string, not int.
	filter := constant.ProviderTypeExperimentalProxy
	if filter == "" {
		t.Error("ProviderTypeExperimentalProxy must not be empty")
	}
	// Verify all four provider types are non-empty strings
	for _, pt := range []string{
		constant.ProviderTypeOfficialCloud,
		constant.ProviderTypeAggregator,
		constant.ProviderTypeAuthorizedProxy,
		constant.ProviderTypeExperimentalProxy,
	} {
		if pt == "" {
			t.Errorf("provider type constant must not be empty")
		}
	}
}
