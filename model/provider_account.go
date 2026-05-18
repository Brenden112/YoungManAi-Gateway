package model

import (
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"gorm.io/gorm"
)

// ProviderAccount stores upstream provider credentials.
// The raw API key is never persisted — only the AES-GCM ciphertext is stored.
type ProviderAccount struct {
	Id           int    `json:"id"`
	Name         string `json:"name" gorm:"type:varchar(128);not null"`
	ProviderType string `json:"provider_type" gorm:"type:varchar(32);default:'official_cloud'"`
	Status       int    `json:"status" gorm:"type:int;default:1"`
	// Key holds the plaintext key in memory only; gorm:"-" ensures it is never written to DB.
	Key string `json:"-" gorm:"-"`
	// EncryptedKey is the AES-256-GCM ciphertext stored in the database.
	EncryptedKey string `json:"-" gorm:"column:encrypted_key;type:text"`
	CreatedAt    int64  `json:"created_at" gorm:"autoCreateTime"`
}

// BeforeSave encrypts Key → EncryptedKey before any INSERT or UPDATE.
// If Key is empty the existing EncryptedKey is left unchanged.
func (pa *ProviderAccount) BeforeSave(tx *gorm.DB) error {
	if pa.Key != "" {
		encrypted, err := common.EncryptAES(pa.Key)
		if err != nil {
			return fmt.Errorf("ProviderAccount.BeforeSave: encrypt key: %w", err)
		}
		pa.EncryptedKey = encrypted
	}
	return nil
}

// AfterFind decrypts EncryptedKey → Key after any SELECT.
func (pa *ProviderAccount) AfterFind(tx *gorm.DB) error {
	if pa.EncryptedKey != "" {
		decrypted, err := common.DecryptAES(pa.EncryptedKey)
		if err != nil {
			return fmt.Errorf("ProviderAccount.AfterFind: decrypt key: %w", err)
		}
		pa.Key = decrypted
	}
	return nil
}

// MaskedKey returns a redacted version of the plaintext key safe for display.
func (pa *ProviderAccount) MaskedKey() string {
	if pa.Key == "" {
		return ""
	}
	if len(pa.Key) <= 8 {
		return strings.Repeat("*", len(pa.Key))
	}
	return pa.Key[:4] + strings.Repeat("*", 10) + pa.Key[len(pa.Key)-4:]
}

// BeforeCreate sets ProviderType default when omitted.
func (pa *ProviderAccount) BeforeCreate(tx *gorm.DB) error {
	if pa.ProviderType == "" {
		pa.ProviderType = constant.ProviderTypeOfficialCloud
	}
	if pa.Status == common.ChannelStatusUnknown {
		pa.Status = common.ChannelStatusEnabled
	}
	return nil
}

func (pa *ProviderAccount) IsEnabled() bool {
	return pa != nil && pa.Status == common.ChannelStatusEnabled
}

// ── CRUD ─────────────────────────────────────────────────────────────────────

func CreateProviderAccount(pa *ProviderAccount) error {
	return DB.Create(pa).Error
}

func GetProviderAccountById(id int) (*ProviderAccount, error) {
	var pa ProviderAccount
	err := DB.First(&pa, id).Error
	if err != nil {
		return nil, err
	}
	return &pa, nil
}

func GetAllProviderAccounts() ([]*ProviderAccount, error) {
	var accounts []*ProviderAccount
	err := DB.Order("id desc").Find(&accounts).Error
	return accounts, err
}

func UpdateProviderAccount(pa *ProviderAccount) error {
	return DB.Save(pa).Error
}

func DeleteProviderAccount(id int) error {
	return DB.Delete(&ProviderAccount{}, id).Error
}
