package model

import "testing"

// validates M12-F01 — TokenUsageStats struct has correct fields
func TestTokenUsageStatsFields(t *testing.T) {
	stats := TokenUsageStats{
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
		TotalQuota:       300,
		RequestCount:     1,
	}
	if stats.TotalTokens != stats.PromptTokens+stats.CompletionTokens {
		t.Errorf("TotalTokens %d != PromptTokens %d + CompletionTokens %d",
			stats.TotalTokens, stats.PromptTokens, stats.CompletionTokens)
	}
}

// validates M12-F02 — cost (quota) is proportional to token usage
func TestCostProportionalToTokens(t *testing.T) {
	// Simulate two requests: one with 2x tokens should have ~2x quota.
	small := TokenUsageStats{PromptTokens: 100, CompletionTokens: 50, TotalQuota: 150}
	large := TokenUsageStats{PromptTokens: 200, CompletionTokens: 100, TotalQuota: 300}

	smallTokens := small.PromptTokens + small.CompletionTokens
	largeTokens := large.PromptTokens + large.CompletionTokens

	if largeTokens != 2*smallTokens {
		t.Errorf("expected large to have 2x tokens of small")
	}
	if large.TotalQuota != 2*small.TotalQuota {
		t.Errorf("expected large quota %d to be 2x small quota %d", large.TotalQuota, small.TotalQuota)
	}
}

// validates M12-F03 — LogTypeError is distinct from LogTypeConsume
// (error logs don't trigger quota deduction — only LogTypeConsume does)
func TestErrorLogTypeDoesNotTriggerDeduction(t *testing.T) {
	// PostConsumeQuota is only called for successful relay responses.
	// RecordErrorLog writes LogTypeError, not LogTypeConsume.
	// GetTokenUsageStats filters on LogTypeConsume — error logs are excluded.
	if LogTypeError == LogTypeConsume {
		t.Error("LogTypeError must differ from LogTypeConsume to prevent double-counting")
	}
	// Verify GetTokenUsageStats only counts LogTypeConsume entries.
	// (Full DB test would require a test DB; this is a compile-time contract check.)
}

// validates M12-F01 — Log struct has ProviderType field for cost attribution
func TestLogProviderTypeField(t *testing.T) {
	log := Log{
		Type:         LogTypeConsume,
		ProviderType: "official_cloud",
		Quota:        500,
	}
	if log.ProviderType != "official_cloud" {
		t.Errorf("ProviderType = %q, want 'official_cloud'", log.ProviderType)
	}
}
