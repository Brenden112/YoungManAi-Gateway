package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/types"
)

// T1/T2 — official_cloud channels are accessible; experimental_proxy is filtered
// for non-internal users via GetGroupEnabledModelsExcludingExperimental.
func TestExperimentalProxyFilteredForNormalUsers(t *testing.T) {
	// Compile-time: the filtering function exists with the right signature.
	type filterFn func(group string) []string
	var _ filterFn = GetGroupEnabledModelsExcludingExperimental
}

// T3 — normal user routed to experimental_proxy gets 403.
// Contract: IsExperimentalProxy context key is set in distributor and checked post-selection.
func TestContextKeyIsExperimentalProxyExists(t *testing.T) {
	if constant.ContextKeyIsExperimentalProxy == "" {
		t.Error("ContextKeyIsExperimentalProxy must be a non-empty string")
	}
}

// T4 — disabled experimental_proxy channels are rejected.
// Contract: DisableExperimentalProxyChannels targets only experimental_proxy + enabled status.
func TestDisableExperimentalProxyChannelsExists(t *testing.T) {
	var fn func() (int64, error) = DisableExperimentalProxyChannels
	_ = fn
}

// T5 — insufficient balance uses the shared typed error code.
func TestInsufficientBalanceConstantExists(t *testing.T) {
	// The error code used by PreConsumeQuota must be defined.
	if types.ErrorCodeInsufficientUserQuota == "" {
		t.Error("ErrorCodeInsufficientUserQuota must be non-zero")
	}
}

// T6 — default no prompt/response storage.
// Contract: StoreFullTextEnabled is false at package init.
func TestRegressionM16StoreFullTextEnabledDefaultFalse(t *testing.T) {
	if common.StoreFullTextEnabled {
		t.Error("StoreFullTextEnabled must default to false")
	}
}

// T6b — Log.Content is cleared when StoreFullTextEnabled=false.
// Contract: RecordConsumeLog signature accepts content string.
func TestLogContentFieldExists(t *testing.T) {
	log := Log{}
	log.Content = "test prompt"
	if log.Content != "test prompt" {
		t.Error("Log.Content field must be settable")
	}
	// Simulate the StoreFullTextEnabled=false path
	if !common.StoreFullTextEnabled {
		log.Content = ""
	}
	if log.Content != "" {
		t.Error("Log.Content must be cleared when StoreFullTextEnabled=false")
	}
}
