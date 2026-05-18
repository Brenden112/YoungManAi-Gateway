package middleware

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

func TestAUD021SetupRejectsPreferredDisallowedProviderWithoutCredentialLeak(t *testing.T) {
	channel := &model.Channel{
		Id:           22101,
		Type:         constant.ChannelTypeOpenRouter,
		Key:          "sk-aud021-do-not-log",
		Status:       common.ChannelStatusEnabled,
		Name:         "aud021-disallowed-preferred",
		ProviderType: constant.ProviderTypeAggregator,
	}
	ctx := newAUD021Context([]string{constant.ProviderTypeOfficialCloud}, false, "")

	apiErr := SetupContextForSelectedChannel(ctx, channel, "aud021-model")
	if apiErr == nil {
		t.Fatal("expected disallowed provider_type rejection")
	}
	if got := common.GetContextKeyString(ctx, constant.ContextKeyChannelKey); got != "" {
		t.Fatalf("credential context set on reject: %q", got)
	}
	if strings.Contains(apiErr.Error(), channel.Key) {
		t.Fatal("rejection error leaked credential")
	}
}

func TestAUD021SetupAllowsInternalExperimentalWhenExplicitlyAllowed(t *testing.T) {
	channel := &model.Channel{
		Id:           22102,
		Type:         constant.ChannelTypeKiroGateway,
		Key:          "sk-aud021-fake-runtime",
		Status:       common.ChannelStatusEnabled,
		Name:         "aud021-exp",
		ProviderType: constant.ProviderTypeExperimentalProxy,
	}
	ctx := newAUD021Context(nil, true, "internal")

	apiErr := SetupContextForSelectedChannel(ctx, channel, "aud021-model")
	if apiErr != nil {
		t.Fatalf("expected internal experimental request to pass: %v", apiErr)
	}
	if got := common.GetContextKeyString(ctx, constant.ContextKeyChannelKey); got != channel.Key {
		t.Fatalf("channel key context = %q, want runtime key", got)
	}
}

func TestAUD021SetupRejectsInternalExperimentalWhenTokenDisallows(t *testing.T) {
	channel := &model.Channel{
		Id:           22103,
		Type:         constant.ChannelTypeKiroGateway,
		Key:          "sk-aud021-fake-runtime",
		Status:       common.ChannelStatusEnabled,
		Name:         "aud021-exp-deny",
		ProviderType: constant.ProviderTypeExperimentalProxy,
	}
	ctx := newAUD021Context(nil, false, "internal")

	apiErr := SetupContextForSelectedChannel(ctx, channel, "aud021-model")
	if apiErr == nil {
		t.Fatal("expected internal request without allow_experimental to be rejected")
	}
	if got := common.GetContextKeyString(ctx, constant.ContextKeyChannelKey); got != "" {
		t.Fatalf("credential context set on reject: %q", got)
	}
}

func newAUD021Context(allowedProviders []string, allowExperimental bool, userGroup string) *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/v1/chat/completions", nil)
	common.SetContextKey(ctx, constant.ContextKeyTokenAllowedProviders, allowedProviders)
	common.SetContextKey(ctx, constant.ContextKeyTokenAllowExperimental, allowExperimental)
	if userGroup != "" {
		common.SetContextKey(ctx, constant.ContextKeyUserGroup, userGroup)
	}
	return ctx
}
