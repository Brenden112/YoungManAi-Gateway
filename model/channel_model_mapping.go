package model

// ChannelModelMapping stores per-channel, per-model configuration:
// public→provider name override, enabled flag, and price hints.
// It complements the existing JSON model_mapping field on Channel —
// entries here take precedence and are merged into the relay context.
type ChannelModelMapping struct {
	ID                int     `json:"id" gorm:"primaryKey;autoIncrement"`
	ChannelId         int     `json:"channel_id" gorm:"uniqueIndex:idx_channel_model;not null"`
	PublicModelName   string  `json:"public_model_name" gorm:"uniqueIndex:idx_channel_model;size:255;not null"`
	ProviderModelName string  `json:"provider_model_name" gorm:"size:255"`
	Enabled           bool    `json:"enabled" gorm:"default:true;not null"`
	InputPrice        float64 `json:"input_price"`  // per 1M tokens; 0 = use channel default
	OutputPrice       float64 `json:"output_price"` // per 1M tokens; 0 = use channel default
	CreatedAt         int64   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         int64   `json:"updated_at" gorm:"autoUpdateTime"`
}

// GetChannelModelMapping returns the mapping for a specific channel+model pair.
// Returns (nil, false) when no record exists.
func GetChannelModelMapping(channelId int, publicModelName string) (*ChannelModelMapping, bool) {
	var m ChannelModelMapping
	err := DB.Where("channel_id = ? AND public_model_name = ?", channelId, publicModelName).
		First(&m).Error
	if err != nil {
		return nil, false
	}
	return &m, true
}

// GetChannelModelMappings returns all mappings for a channel.
func GetChannelModelMappings(channelId int) ([]*ChannelModelMapping, error) {
	var mappings []*ChannelModelMapping
	err := DB.Where("channel_id = ?", channelId).Find(&mappings).Error
	return mappings, err
}

// UpsertChannelModelMapping creates or updates a mapping record.
func UpsertChannelModelMapping(m *ChannelModelMapping) error {
	return DB.Save(m).Error
}

// DeleteChannelModelMapping removes a mapping for a specific channel+model pair.
func DeleteChannelModelMapping(channelId int, publicModelName string) error {
	return DB.Where("channel_id = ? AND public_model_name = ?", channelId, publicModelName).
		Delete(&ChannelModelMapping{}).Error
}
