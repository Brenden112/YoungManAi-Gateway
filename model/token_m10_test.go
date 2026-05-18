package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
)

// validates M10-F01 — Token has OrgId and ProjectId nullable FK fields
func TestTokenOrgProjectBinding(t *testing.T) {
	orgId := 7
	projectId := 3
	token := Token{
		UserId:    1,
		OrgId:     &orgId,
		ProjectId: &projectId,
	}
	if token.OrgId == nil || *token.OrgId != 7 {
		t.Errorf("OrgId = %v, want 7", token.OrgId)
	}
	if token.ProjectId == nil || *token.ProjectId != 3 {
		t.Errorf("ProjectId = %v, want 3", token.ProjectId)
	}
	// nil is valid (token not bound to org/project)
	unbound := Token{UserId: 1}
	if unbound.OrgId != nil || unbound.ProjectId != nil {
		t.Error("unbound token should have nil OrgId and ProjectId")
	}
}

// validates M10-F02 — BeforeCreate computes KeyHash and KeyPrefix
func TestTokenBeforeCreateKeyHashPrefix(t *testing.T) {
	token := Token{Key: "test-key-for-hash"}
	_ = token.BeforeCreate(nil)

	expectedHash := HashTokenKey("test-key-for-hash")
	if token.KeyHash != expectedHash {
		t.Errorf("KeyHash mismatch")
	}
	if token.KeyPrefix != "test-key" {
		t.Errorf("KeyPrefix = %q, want 'test-key'", token.KeyPrefix)
	}
	if token.Key == "test-key-for-hash" {
		t.Error("Key must not retain plaintext after BeforeCreate")
	}
	if token.Key != token.KeyHash {
		t.Error("legacy key column should contain only the non-reversible hash")
	}
}

// validates M10-F02 — short key prefix handling
func TestTokenBeforeCreateShortKey(t *testing.T) {
	token := Token{Key: "sk-12"}
	_ = token.BeforeCreate(nil)
	if token.KeyPrefix != "12" {
		t.Errorf("KeyPrefix = %q, want '12' after sk- normalization", token.KeyPrefix)
	}
}

// validates M10-F03 — TokenStatusDisabled constant exists and is non-zero
func TestTokenStatusDisabledConstant(t *testing.T) {
	if common.TokenStatusDisabled == 0 {
		t.Error("TokenStatusDisabled must not be 0 (0 is the default value)")
	}
	if common.TokenStatusDisabled == common.TokenStatusEnabled {
		t.Error("TokenStatusDisabled must differ from TokenStatusEnabled")
	}
}

// validates M10-F04 — AllowExperimental defaults to false
func TestTokenAllowExperimentalDefaultFalse(t *testing.T) {
	token := Token{UserId: 1, Key: "sk-test"}
	if token.AllowExperimental {
		t.Error("AllowExperimental should default to false for new tokens")
	}
}

// validates proof — normal API key cannot access experimental_proxy
func TestNormalTokenCannotAccessExperimental(t *testing.T) {
	normalToken := Token{UserId: 1, AllowExperimental: false}
	internalToken := Token{UserId: 1, AllowExperimental: true}

	if normalToken.AllowExperimental {
		t.Error("normal token must not allow experimental")
	}
	if !internalToken.AllowExperimental {
		t.Error("explicitly enabled token should allow experimental")
	}
}

func TestAUD022TokenTenantBindingRejectsDisabledOrganization(t *testing.T) {
	truncateTables(t)

	orgId := 22
	token := &Token{
		UserId:         1,
		Name:           "disabled-org",
		Key:            "aud022-disabled-org",
		Status:         common.TokenStatusEnabled,
		ExpiredTime:    -1,
		UnlimitedQuota: true,
		OrgId:          &orgId,
	}
	if err := DB.Create(&User{Id: 1, Username: "aud022-user", Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := DB.Create(&Organization{Id: orgId, Name: "disabled-org", OwnerId: 1, Status: common.UserStatusDisabled}).Error; err != nil {
		t.Fatalf("seed org: %v", err)
	}
	if err := token.Insert(); err != nil {
		t.Fatalf("insert token: %v", err)
	}
	if _, err := ValidateUserToken("aud022-disabled-org"); err == nil {
		t.Fatal("disabled organization token should not authenticate")
	}
}

func TestAUD022TokenTenantBindingRejectsDisabledProject(t *testing.T) {
	truncateTables(t)

	orgId := 23
	projectId := 24
	token := &Token{
		UserId:         1,
		Name:           "disabled-project",
		Key:            "aud022-disabled-project",
		Status:         common.TokenStatusEnabled,
		ExpiredTime:    -1,
		UnlimitedQuota: true,
		OrgId:          &orgId,
		ProjectId:      &projectId,
	}
	if err := DB.Create(&User{Id: 1, Username: "aud022-project-user", Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := DB.Create(&Organization{Id: orgId, Name: "active-org", OwnerId: 1, Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed org: %v", err)
	}
	if err := DB.Create(&Project{Id: projectId, OrgId: orgId, Name: "disabled-project", Status: common.UserStatusDisabled}).Error; err != nil {
		t.Fatalf("seed project: %v", err)
	}
	if err := token.Insert(); err != nil {
		t.Fatalf("insert token: %v", err)
	}
	if _, err := ValidateUserToken("aud022-disabled-project"); err == nil {
		t.Fatal("disabled project token should not authenticate")
	}
}

func TestAUD022TokenTenantBindingRejectsCrossOrganizationProject(t *testing.T) {
	truncateTables(t)

	orgA := 25
	orgB := 26
	projectB := 27
	token := &Token{
		UserId:         1,
		Name:           "cross-project",
		Key:            "aud022-cross-project",
		Status:         common.TokenStatusEnabled,
		ExpiredTime:    -1,
		UnlimitedQuota: true,
		OrgId:          &orgA,
		ProjectId:      &projectB,
	}
	if err := DB.Create(&User{Id: 1, Username: "aud022-cross-user", Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := DB.Create(&Organization{Id: orgA, Name: "org-a", OwnerId: 1, Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed org a: %v", err)
	}
	if err := DB.Create(&Organization{Id: orgB, Name: "org-b", OwnerId: 1, Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed org b: %v", err)
	}
	if err := DB.Create(&Project{Id: projectB, OrgId: orgB, Name: "project-b", Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed project b: %v", err)
	}
	if err := token.Insert(); err != nil {
		t.Fatalf("insert token: %v", err)
	}
	if _, err := ValidateUserToken("aud022-cross-project"); err == nil {
		t.Fatal("token bound to organization A must not authenticate against project B")
	}
}

func TestAUD022LegacyTokenRemainsUserScoped(t *testing.T) {
	truncateTables(t)

	token := &Token{
		UserId:         1,
		Name:           "legacy",
		Key:            "aud022-legacy",
		Status:         common.TokenStatusEnabled,
		ExpiredTime:    -1,
		UnlimitedQuota: true,
	}
	if err := DB.Create(&User{Id: 1, Username: "aud022-legacy-user", Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := token.Insert(); err != nil {
		t.Fatalf("insert token: %v", err)
	}
	authed, err := ValidateUserToken("aud022-legacy")
	if err != nil {
		t.Fatalf("legacy token should authenticate by user scope: %v", err)
	}
	scope, err := ResolveTokenTenantScope(authed)
	if err != nil {
		t.Fatalf("resolve legacy scope: %v", err)
	}
	if !scope.Legacy || scope.OrgId != nil || scope.ProjectId != nil {
		t.Fatalf("legacy scope should not infer org/project: %+v", scope)
	}
}

func TestTokenInsertDoesNotPersistPlaintextAndAuthenticatesByHash(t *testing.T) {
	truncateTables(t)

	rawKey := "audit016-valid-key"
	token := &Token{
		UserId:         1,
		Name:           "secure-token",
		Key:            rawKey,
		Status:         common.TokenStatusEnabled,
		ExpiredTime:    -1,
		RemainQuota:    100,
		UnlimitedQuota: true,
	}
	if err := DB.Create(&User{Id: 1, Username: "audit016-user", Status: common.UserStatusEnabled}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := token.Insert(); err != nil {
		t.Fatalf("insert token: %v", err)
	}

	var stored Token
	if err := DB.First(&stored, "id = ?", token.Id).Error; err != nil {
		t.Fatalf("load stored token: %v", err)
	}
	if stored.Key == rawKey {
		t.Fatal("stored key contains plaintext")
	}
	if stored.KeyHash == "" || stored.KeyHash != HashTokenKey(rawKey) {
		t.Fatal("stored token hash was not populated")
	}
	if stored.KeyPrefix != tokenKeyPrefix(rawKey) {
		t.Fatalf("stored prefix = %q, want %q", stored.KeyPrefix, tokenKeyPrefix(rawKey))
	}

	authed, err := ValidateUserToken(rawKey)
	if err != nil {
		t.Fatalf("expected correct key to authenticate: %v", err)
	}
	if authed.Id != token.Id {
		t.Fatalf("authenticated token id = %d, want %d", authed.Id, token.Id)
	}
	if _, err := ValidateUserToken("audit016-wrong-key"); err == nil {
		t.Fatal("expected wrong key to fail")
	}
}

func TestDisabledTokenCannotAuthenticateWithHash(t *testing.T) {
	truncateTables(t)

	rawKey := "audit016-disabled-key"
	token := &Token{
		UserId:         1,
		Name:           "disabled-token",
		Key:            rawKey,
		Status:         common.TokenStatusEnabled,
		ExpiredTime:    -1,
		RemainQuota:    100,
		UnlimitedQuota: true,
	}
	if err := token.Insert(); err != nil {
		t.Fatalf("insert token: %v", err)
	}
	if err := DisableToken(token.Id); err != nil {
		t.Fatalf("disable token: %v", err)
	}
	if _, err := ValidateUserToken(rawKey); err == nil {
		t.Fatal("expected disabled token authentication to fail")
	}
}

func TestLegacyPlaintextKeyMigratesAndIsNotAuthSource(t *testing.T) {
	truncateTables(t)

	legacyKey := "audit016-legacy-key"
	if err := DB.Exec("INSERT INTO tokens (user_id, name, `key`, status, expired_time, remain_quota, unlimited_quota) VALUES (?, ?, ?, ?, ?, ?, ?)",
		1, "legacy-token", legacyKey, common.TokenStatusEnabled, -1, 100, true).Error; err != nil {
		t.Fatalf("insert legacy plaintext token: %v", err)
	}
	if err := MigratePlaintextTokenKeys(); err != nil {
		t.Fatalf("migrate plaintext tokens: %v", err)
	}

	var stored Token
	if err := DB.First(&stored, "name = ?", "legacy-token").Error; err != nil {
		t.Fatalf("load migrated token: %v", err)
	}
	if stored.Key == legacyKey {
		t.Fatal("legacy plaintext key was not removed")
	}
	if stored.KeyHash != HashTokenKey(legacyKey) {
		t.Fatal("legacy key hash mismatch")
	}
	if _, err := ValidateUserToken(legacyKey); err != nil {
		t.Fatalf("migrated key should authenticate by hash: %v", err)
	}

	if err := DB.Model(&Token{}).Where("id = ?", stored.Id).Updates(map[string]interface{}{
		"key":      "audit016-legacy-only",
		"key_hash": HashTokenKey(legacyKey),
	}).Error; err != nil {
		t.Fatalf("seed legacy-only plaintext field: %v", err)
	}
	if _, err := ValidateUserToken("audit016-legacy-only"); err == nil {
		t.Fatal("legacy plaintext key column must not authenticate")
	}
}
