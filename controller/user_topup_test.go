package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupUserTopUpTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	gin.SetMode(gin.TestMode)
	oldDB := model.DB
	oldLogDB := model.LOG_DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	model.DB = db
	model.LOG_DB = db
	common.RedisEnabled = false
	if err := db.AutoMigrate(&model.User{}, &model.Log{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		model.DB = oldDB
		model.LOG_DB = oldLogDB
	})
	return db
}

func newTopUpContext(t *testing.T, role int, body any) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	payload, err := common.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/user/topup", bytes.NewReader(payload))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("id", 1)
	ctx.Set("username", "admin")
	ctx.Set("role", role)
	return ctx, recorder
}

func TestAdminManualTopUpSucceedsAndWritesManageLog(t *testing.T) {
	db := setupUserTopUpTestDB(t)
	if err := db.Create(&model.User{Id: 2, Username: "topup-user", Status: common.UserStatusEnabled, Quota: 100}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	ctx, recorder := newTopUpContext(t, common.RoleAdminUser, gin.H{"user_id": 2, "quota": 50})

	TopUp(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	var response tokenAPIResponse
	if err := common.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !response.Success {
		t.Fatalf("topup failed: %s", response.Message)
	}
	user, err := model.GetUserById(2, false)
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.Quota != 150 {
		t.Fatalf("quota = %d, want 150", user.Quota)
	}
	var count int64
	if err := db.Model(&model.Log{}).Where("user_id = ? AND type = ?", 2, model.LogTypeManage).Count(&count).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if count != 1 {
		t.Fatalf("manage log count = %d, want 1", count)
	}
}

func TestNormalUserManualTopUpRejected(t *testing.T) {
	setupUserTopUpTestDB(t)
	ctx, recorder := newTopUpContext(t, common.RoleCommonUser, gin.H{"user_id": 2, "quota": 50})

	TopUp(ctx)

	var response tokenAPIResponse
	if err := common.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Success {
		t.Fatal("normal user manual topup should be rejected")
	}
}

func TestAdminManualTopUpRejectsInvalidAmount(t *testing.T) {
	setupUserTopUpTestDB(t)
	ctx, recorder := newTopUpContext(t, common.RoleAdminUser, gin.H{"user_id": 2, "quota": -1})

	TopUp(ctx)

	var response tokenAPIResponse
	if err := common.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Success {
		t.Fatal("negative manual topup should be rejected")
	}
}
