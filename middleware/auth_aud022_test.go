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

func setupAUD022MiddlewareDB(t *testing.T) *gorm.DB {
	t.Helper()
	oldDB := model.DB
	oldLogDB := model.LOG_DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	model.DB = db
	model.LOG_DB = db
	common.RedisEnabled = false
	common.LogConsumeEnabled = true
	if err := db.AutoMigrate(&model.User{}, &model.Token{}, &model.Organization{}, &model.OrganizationMember{}, &model.Project{}, &model.Log{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		model.DB = oldDB
		model.LOG_DB = oldLogDB
	})
	return db
}

func newAUD022Context() *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/v1/chat/completions", nil)
	return ctx
}

func TestAUD022SetupContextUsesTokenTenantAndIgnoresSpoofedContext(t *testing.T) {
	db := setupAUD022MiddlewareDB(t)
	orgId := 101
	projectId := 201
	if err := db.Create(&model.Organization{Id: orgId, Name: "aud022-org", OwnerId: 7, Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed org: %v", err)
	}
	if err := db.Create(&model.Project{Id: projectId, OrgId: orgId, Name: "aud022-project", Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed project: %v", err)
	}
	token := &model.Token{
		Id:                   301,
		UserId:               7,
		Name:                 "tenant-token",
		RemainQuota:          99,
		ModelLimitsEnabled:   true,
		ModelLimits:          "gpt-a",
		AllowedProviderTypes: constant.ProviderTypeOfficialCloud,
		OrgId:                &orgId,
		ProjectId:            &projectId,
	}
	ctx := newAUD022Context()
	common.SetContextKey(ctx, constant.ContextKeyTokenOrgId, 999)
	common.SetContextKey(ctx, constant.ContextKeyTokenProjectId, 999)

	if err := SetupContextForToken(ctx, token); err != nil {
		t.Fatalf("setup token context: %v", err)
	}
	if got := common.GetContextKeyInt(ctx, constant.ContextKeyTokenOrgId); got != orgId {
		t.Fatalf("org context = %d, want token org %d", got, orgId)
	}
	if got := common.GetContextKeyInt(ctx, constant.ContextKeyTokenProjectId); got != projectId {
		t.Fatalf("project context = %d, want token project %d", got, projectId)
	}
	if !ctx.GetBool("token_model_limit_enabled") {
		t.Fatal("model limit should come from authenticated token")
	}
	modelLimitValue, ok := ctx.Get("token_model_limit")
	if !ok {
		t.Fatal("model limit map missing from context")
	}
	modelLimit, ok := modelLimitValue.(map[string]bool)
	if !ok || !modelLimit["gpt-a"] {
		t.Fatal("model limit map missing token-bound model")
	}
	if got := common.GetContextKeyStringSlice(ctx, constant.ContextKeyTokenAllowedProviders); len(got) != 1 || got[0] != constant.ProviderTypeOfficialCloud {
		t.Fatalf("allowed providers = %v, want official token scope", got)
	}

	model.RecordConsumeLog(ctx, token.UserId, model.RecordConsumeLogParams{TokenId: token.Id, ModelName: "gpt-a", Quota: 1})
	var log model.Log
	if err := db.First(&log, "token_id = ?", token.Id).Error; err != nil {
		t.Fatalf("load usage log: %v", err)
	}
	if log.OrgId == nil || *log.OrgId != orgId {
		t.Fatalf("log org = %v, want %d", log.OrgId, orgId)
	}
	if log.ProjectId == nil || *log.ProjectId != projectId {
		t.Fatalf("log project = %v, want %d", log.ProjectId, projectId)
	}
}

func TestAUD022LegacyTokenClearsSpoofedTenantContext(t *testing.T) {
	setupAUD022MiddlewareDB(t)
	ctx := newAUD022Context()
	common.SetContextKey(ctx, constant.ContextKeyTokenOrgId, 999)
	common.SetContextKey(ctx, constant.ContextKeyTokenProjectId, 999)
	token := &model.Token{Id: 302, UserId: 8, Name: "legacy-token", UnlimitedQuota: true}

	if err := SetupContextForToken(ctx, token); err != nil {
		t.Fatalf("setup legacy token context: %v", err)
	}
	if got := common.GetContextKeyInt(ctx, constant.ContextKeyTokenOrgId); got != 0 {
		t.Fatalf("legacy token should not retain spoofed org, got %d", got)
	}
	if got := common.GetContextKeyInt(ctx, constant.ContextKeyTokenProjectId); got != 0 {
		t.Fatalf("legacy token should not retain spoofed project, got %d", got)
	}
}
