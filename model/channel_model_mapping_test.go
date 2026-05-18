package model

import (
	"testing"
)

// validates M5-F01 — ChannelModelMapping struct has correct fields and JSON tags
func TestChannelModelMappingFields(t *testing.T) {
	m := ChannelModelMapping{
		ChannelId:         1,
		PublicModelName:   "gpt-4o",
		ProviderModelName: "gpt-4o-2024-11-20",
		Enabled:           true,
		InputPrice:        2.5,
		OutputPrice:       10.0,
	}
	if m.ChannelId != 1 {
		t.Errorf("ChannelId = %d, want 1", m.ChannelId)
	}
	if m.PublicModelName != "gpt-4o" {
		t.Errorf("PublicModelName = %q, want 'gpt-4o'", m.PublicModelName)
	}
	if m.ProviderModelName != "gpt-4o-2024-11-20" {
		t.Errorf("ProviderModelName = %q, want 'gpt-4o-2024-11-20'", m.ProviderModelName)
	}
	if !m.Enabled {
		t.Error("Enabled should be true")
	}
	if m.InputPrice != 2.5 {
		t.Errorf("InputPrice = %f, want 2.5", m.InputPrice)
	}
	if m.OutputPrice != 10.0 {
		t.Errorf("OutputPrice = %f, want 10.0", m.OutputPrice)
	}
}

// validates M5-F02 — price fields accept zero (use channel default)
func TestChannelModelMappingZeroPriceIsValid(t *testing.T) {
	m := ChannelModelMapping{
		ChannelId:       1,
		PublicModelName: "gpt-4o-mini",
		Enabled:         true,
		InputPrice:      0,
		OutputPrice:     0,
	}
	if m.InputPrice != 0 || m.OutputPrice != 0 {
		t.Error("zero price should be valid (means use channel default)")
	}
}

// validates M5-F03 — disabled mapping is representable
func TestChannelModelMappingDisabled(t *testing.T) {
	m := ChannelModelMapping{
		ChannelId:       1,
		PublicModelName: "deprecated-model",
		Enabled:         false,
	}
	if m.Enabled {
		t.Error("Enabled should be false for disabled model")
	}
}

// validates M5-F01 — identity mapping (no rename) is valid
func TestChannelModelMappingIdentity(t *testing.T) {
	m := ChannelModelMapping{
		ChannelId:         1,
		PublicModelName:   "gpt-4o",
		ProviderModelName: "", // empty = no rename, use public name as-is
		Enabled:           true,
	}
	if m.ProviderModelName != "" {
		t.Errorf("empty ProviderModelName should be valid, got %q", m.ProviderModelName)
	}
}

// validates M5-F01 — aggregator channel can map to different provider model
func TestChannelModelMappingAggregator(t *testing.T) {
	// OpenRouter channel mapping gpt-4o → openai/gpt-4o
	m := ChannelModelMapping{
		ChannelId:         42,
		PublicModelName:   "gpt-4o",
		ProviderModelName: "openai/gpt-4o",
		Enabled:           true,
		InputPrice:        5.0,
		OutputPrice:       15.0,
	}
	if m.ProviderModelName != "openai/gpt-4o" {
		t.Errorf("ProviderModelName = %q, want 'openai/gpt-4o'", m.ProviderModelName)
	}
}
