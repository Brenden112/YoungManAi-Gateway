package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestAUD018SetupContextUsesProviderAccountCredential(t *testing.T) {
	db := setupAUD018MiddlewareDB(t)
	withAUD018CryptoSecret(t, "aud018-middleware-secret")

	account := seedAUD018MiddlewareProviderAccount(t, db, "sk-aud018-provider-runtime", common.ChannelStatusEnabled)
	channel := &model.Channel{
		Id:                2801,
		Type:              constant.ChannelTypeOpenAI,
		Key:               "sk-aud018-legacy-channel",
		Status:            common.ChannelStatusEnabled,
		Name:              "aud018-provider-account-channel",
		ProviderAccountId: &account.Id,
	}

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/v1/chat/completions", nil)
	apiErr := SetupContextForSelectedChannel(ctx, channel, "aud018-model")
	if apiErr != nil {
		t.Fatalf("SetupContextForSelectedChannel error: %v", apiErr)
	}
	got := common.GetContextKeyString(ctx, constant.ContextKeyChannelKey)
	if got != "sk-aud018-provider-runtime" {
		t.Fatalf("channel key context = %q, want provider account credential", got)
	}
	if got == channel.Key {
		t.Fatal("provider_account_id path used legacy channel key")
	}
}

func TestAUD018SetupContextRejectsProviderAccountDecryptFailure(t *testing.T) {
	db := setupAUD018MiddlewareDB(t)
	withAUD018CryptoSecret(t, "aud018-middleware-create-secret")
	account := seedAUD018MiddlewareProviderAccount(t, db, "sk-aud018-provider-decrypt-failure", common.ChannelStatusEnabled)
	common.CryptoSecret = "aud018-middleware-wrong-secret"

	channel := &model.Channel{
		Id:                2802,
		Type:              constant.ChannelTypeOpenAI,
		Key:               "sk-aud018-legacy-channel",
		Status:            common.ChannelStatusEnabled,
		Name:              "aud018-provider-account-channel",
		ProviderAccountId: &account.Id,
	}

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/v1/chat/completions", nil)
	apiErr := SetupContextForSelectedChannel(ctx, channel, "aud018-model")
	if apiErr == nil {
		t.Fatal("expected provider account decrypt failure")
	}
	if got := common.GetContextKeyString(ctx, constant.ContextKeyChannelKey); got != "" {
		t.Fatalf("credential context was set after decrypt failure: %q", got)
	}
}

func TestAUD018SetupContextRejectsDisabledProviderAccount(t *testing.T) {
	db := setupAUD018MiddlewareDB(t)
	withAUD018CryptoSecret(t, "aud018-middleware-secret")
	account := seedAUD018MiddlewareProviderAccount(t, db, "sk-aud018-provider-disabled", common.ChannelStatusManuallyDisabled)

	channel := &model.Channel{
		Id:                2803,
		Type:              constant.ChannelTypeOpenAI,
		Key:               "sk-aud018-legacy-channel",
		Status:            common.ChannelStatusEnabled,
		Name:              "aud018-provider-account-channel",
		ProviderAccountId: &account.Id,
	}

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/v1/chat/completions", nil)
	apiErr := SetupContextForSelectedChannel(ctx, channel, "aud018-model")
	if apiErr == nil {
		t.Fatal("expected disabled provider account rejection")
	}
	if got := common.GetContextKeyString(ctx, constant.ContextKeyChannelKey); got != "" {
		t.Fatalf("credential context was set for disabled provider account: %q", got)
	}
}

func setupAUD018MiddlewareDB(t *testing.T) *gorm.DB {
	t.Helper()
	oldDB := model.DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	model.DB = db
	if err := db.AutoMigrate(&model.ProviderAccount{}, &model.Channel{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		model.DB = oldDB
	})
	return db
}

func seedAUD018MiddlewareProviderAccount(t *testing.T, db *gorm.DB, credential string, status int) *model.ProviderAccount {
	t.Helper()
	account := &model.ProviderAccount{
		Name:         "aud018-provider-account",
		ProviderType: constant.ProviderTypeOfficialCloud,
		Status:       status,
		Key:          credential,
	}
	if err := db.Create(account).Error; err != nil {
		t.Fatalf("create provider account: %v", err)
	}
	return account
}

func withAUD018CryptoSecret(t *testing.T, secret string) {
	t.Helper()
	oldSecret := common.CryptoSecret
	common.CryptoSecret = secret
	t.Cleanup(func() {
		common.CryptoSecret = oldSecret
	})
}
