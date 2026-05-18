package model

import (
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
)

// validates A2.1.3 — BeforeSave encrypts Key; plaintext never equals ciphertext
func TestProviderAccountBeforeSaveEncryptsKey(t *testing.T) {
	pa := &ProviderAccount{
		Name:         "test-account",
		ProviderType: constant.ProviderTypeOfficialCloud,
		Key:          "sk-upstream-test-key-12345",
	}
	if err := pa.BeforeSave(nil); err != nil {
		t.Fatalf("BeforeSave error: %v", err)
	}
	if pa.EncryptedKey == "" {
		t.Error("EncryptedKey should be non-empty after BeforeSave")
	}
	if pa.EncryptedKey == pa.Key {
		t.Error("EncryptedKey must not equal plaintext Key")
	}
}

// validates A2.1.3 — AfterFind decrypts EncryptedKey back to original Key
func TestProviderAccountAfterFindDecryptsKey(t *testing.T) {
	original := "sk-upstream-test-key-12345"
	pa := &ProviderAccount{Key: original}
	if err := pa.BeforeSave(nil); err != nil {
		t.Fatalf("BeforeSave error: %v", err)
	}

	pa2 := &ProviderAccount{EncryptedKey: pa.EncryptedKey}
	if err := pa2.AfterFind(nil); err != nil {
		t.Fatalf("AfterFind error: %v", err)
	}
	if pa2.Key != original {
		t.Errorf("AfterFind Key = %q, want %q", pa2.Key, original)
	}
}

// validates A2.1.3 — empty Key does not overwrite existing EncryptedKey
func TestProviderAccountBeforeSaveEmptyKeyPreservesEncryptedKey(t *testing.T) {
	existing := "some-existing-ciphertext"
	pa := &ProviderAccount{EncryptedKey: existing} // Key is empty
	if err := pa.BeforeSave(nil); err != nil {
		t.Fatalf("BeforeSave error: %v", err)
	}
	if pa.EncryptedKey != existing {
		t.Errorf("BeforeSave overwrote EncryptedKey when Key was empty: got %q", pa.EncryptedKey)
	}
}

// validates A2.1.2 — ProviderType defaults to official_cloud via BeforeCreate
func TestProviderAccountBeforeCreateDefaultProviderType(t *testing.T) {
	pa := &ProviderAccount{Name: "test"}
	if err := pa.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate error: %v", err)
	}
	if pa.ProviderType != constant.ProviderTypeOfficialCloud {
		t.Errorf("ProviderType = %q, want %q", pa.ProviderType, constant.ProviderTypeOfficialCloud)
	}
}

// validates EncryptAES/DecryptAES round-trip in common package
func TestEncryptDecryptAESRoundTrip(t *testing.T) {
	cases := []string{
		"sk-short",
		"sk-a-longer-upstream-api-key-with-special-chars-!@#$%",
		"",
	}
	for _, plaintext := range cases {
		if plaintext == "" {
			continue // empty string edge case — skip
		}
		encrypted, err := common.EncryptAES(plaintext)
		if err != nil {
			t.Errorf("EncryptAES(%q) error: %v", plaintext, err)
			continue
		}
		if encrypted == plaintext {
			t.Errorf("EncryptAES(%q) returned plaintext unchanged", plaintext)
		}
		decrypted, err := common.DecryptAES(encrypted)
		if err != nil {
			t.Errorf("DecryptAES error: %v", err)
			continue
		}
		if decrypted != plaintext {
			t.Errorf("round-trip failed: got %q, want %q", decrypted, plaintext)
		}
	}
}

// validates MaskedKey redacts the key safely
func TestProviderAccountMaskedKey(t *testing.T) {
	pa := &ProviderAccount{Key: "sk-abcdefghijklmnop"}
	masked := pa.MaskedKey()
	if masked == pa.Key {
		t.Error("MaskedKey should not return the full plaintext key")
	}
	if masked == "" {
		t.Error("MaskedKey should not be empty for a non-empty key")
	}
}

func TestProviderAccountCredentialNotPersistedAsPlaintext(t *testing.T) {
	resetProviderAccountTestTables(t)
	withCryptoSecret(t, "aud018-test-secret")

	rawCredential := "sk-aud018-provider-account-runtime"
	account := &ProviderAccount{
		Name:         "aud018-account",
		ProviderType: constant.ProviderTypeOfficialCloud,
		Key:          rawCredential,
	}
	if err := DB.Create(account).Error; err != nil {
		t.Fatalf("create provider account: %v", err)
	}

	var encrypted string
	if err := DB.Table("provider_accounts").
		Select("encrypted_key").
		Where("id = ?", account.Id).
		Scan(&encrypted).Error; err != nil {
		t.Fatalf("load encrypted key: %v", err)
	}
	if encrypted == "" {
		t.Fatal("encrypted_key is empty")
	}
	if encrypted == rawCredential || strings.Contains(encrypted, rawCredential) {
		t.Fatal("provider account credential was persisted as plaintext")
	}
}

func TestResolveActiveCredentialUsesProviderAccountCredential(t *testing.T) {
	resetProviderAccountTestTables(t)
	withCryptoSecret(t, "aud018-test-secret")

	account := seedProviderAccount(t, "aud018-account", "sk-aud018-provider", common.ChannelStatusEnabled)
	channel := &Channel{
		Id:                1801,
		Key:               "sk-aud018-legacy-channel",
		ProviderAccountId: &account.Id,
	}

	key, index, apiErr := channel.ResolveActiveCredential()
	if apiErr != nil {
		t.Fatalf("ResolveActiveCredential error: %v", apiErr)
	}
	if index != 0 {
		t.Fatalf("key index = %d, want 0 for provider account credential", index)
	}
	if key != "sk-aud018-provider" {
		t.Fatalf("credential = %q, want provider account credential", key)
	}
}

func TestResolveActiveCredentialDoesNotFallbackToChannelKeyWhenProviderAccountExists(t *testing.T) {
	resetProviderAccountTestTables(t)
	withCryptoSecret(t, "aud018-test-secret")

	account := seedProviderAccount(t, "aud018-account", "sk-aud018-provider-primary", common.ChannelStatusEnabled)
	channel := &Channel{
		Id:                1802,
		Key:               "sk-aud018-legacy-should-not-be-used",
		ProviderAccountId: &account.Id,
	}

	key, _, apiErr := channel.ResolveActiveCredential()
	if apiErr != nil {
		t.Fatalf("ResolveActiveCredential error: %v", apiErr)
	}
	if key == channel.Key {
		t.Fatal("provider_account_id path fell back to legacy channel key")
	}
}

func TestResolveActiveCredentialRejectsProviderAccountDecryptFailure(t *testing.T) {
	resetProviderAccountTestTables(t)
	withCryptoSecret(t, "aud018-create-secret")

	account := seedProviderAccount(t, "aud018-account", "sk-aud018-provider-decrypt", common.ChannelStatusEnabled)
	withCryptoSecret(t, "aud018-wrong-secret")

	channel := &Channel{
		Id:                1803,
		Key:               "sk-aud018-legacy-channel",
		ProviderAccountId: &account.Id,
	}
	key, _, apiErr := channel.ResolveActiveCredential()
	if apiErr == nil {
		t.Fatal("expected decrypt failure to reject provider account credential")
	}
	if key != "" {
		t.Fatal("decrypt failure returned a credential")
	}
	if strings.Contains(apiErr.Error(), "sk-aud018-provider-decrypt") || strings.Contains(apiErr.Error(), channel.Key) {
		t.Fatal("credential leaked in decrypt failure error")
	}
}

func TestResolveActiveCredentialRejectsDisabledProviderAccount(t *testing.T) {
	resetProviderAccountTestTables(t)
	withCryptoSecret(t, "aud018-test-secret")

	account := seedProviderAccount(t, "aud018-disabled-account", "sk-aud018-disabled-provider", common.ChannelStatusManuallyDisabled)
	channel := &Channel{
		Id:                1804,
		Key:               "sk-aud018-legacy-channel",
		ProviderAccountId: &account.Id,
	}

	key, _, apiErr := channel.ResolveActiveCredential()
	if apiErr == nil {
		t.Fatal("expected disabled provider account to be rejected")
	}
	if key != "" {
		t.Fatal("disabled provider account returned a credential")
	}
	if strings.Contains(apiErr.Error(), "sk-aud018-disabled-provider") || strings.Contains(apiErr.Error(), channel.Key) {
		t.Fatal("credential leaked in disabled provider account error")
	}
}

func TestResolveActiveCredentialLegacyChannelStillCompatible(t *testing.T) {
	channel := &Channel{
		Id:  1805,
		Key: "sk-aud018-legacy-compatible",
	}

	key, index, apiErr := channel.ResolveActiveCredential()
	if apiErr != nil {
		t.Fatalf("ResolveActiveCredential error: %v", apiErr)
	}
	if index != 0 {
		t.Fatalf("key index = %d, want 0", index)
	}
	if key != channel.Key {
		t.Fatalf("legacy credential = %q, want channel key", key)
	}
}

func resetProviderAccountTestTables(t *testing.T) {
	t.Helper()
	DB.Exec("DELETE FROM channels")
	DB.Exec("DELETE FROM provider_accounts")
	t.Cleanup(func() {
		DB.Exec("DELETE FROM channels")
		DB.Exec("DELETE FROM provider_accounts")
	})
}

func seedProviderAccount(t *testing.T, name string, credential string, status int) *ProviderAccount {
	t.Helper()
	account := &ProviderAccount{
		Name:         name,
		ProviderType: constant.ProviderTypeOfficialCloud,
		Status:       status,
		Key:          credential,
	}
	if err := DB.Create(account).Error; err != nil {
		t.Fatalf("create provider account: %v", err)
	}
	return account
}

func withCryptoSecret(t *testing.T, secret string) {
	t.Helper()
	oldSecret := common.CryptoSecret
	common.CryptoSecret = secret
	t.Cleanup(func() {
		common.CryptoSecret = oldSecret
	})
}
