package model

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Ability struct {
	Group     string  `json:"group" gorm:"type:varchar(64);primaryKey;autoIncrement:false"`
	Model     string  `json:"model" gorm:"type:varchar(255);primaryKey;autoIncrement:false"`
	ChannelId int     `json:"channel_id" gorm:"primaryKey;autoIncrement:false;index"`
	Enabled   bool    `json:"enabled"`
	Priority  *int64  `json:"priority" gorm:"bigint;default:0;index"`
	Weight    uint    `json:"weight" gorm:"default:0;index"`
	Tag       *string `json:"tag" gorm:"index"`
}

type AbilityWithChannel struct {
	Ability
	ChannelType int `json:"channel_type"`
}

func GetAllEnableAbilityWithChannels() ([]AbilityWithChannel, error) {
	var abilities []AbilityWithChannel
	err := DB.Table("abilities").
		Select("abilities.*, channels.type as channel_type").
		Joins("join channels on abilities.channel_id = channels.id").
		Where("abilities.enabled = ? AND channels.status = ?", true, common.ChannelStatusEnabled).
		Scan(&abilities).Error
	return abilities, err
}

func GetGroupEnabledModels(group string) []string {
	return GetGroupEnabledModelsWithProviderPolicy(group, ProviderTypePolicyFromAllowExperimental(true))
}

func GetGroupEnabledModelsWithProviderPolicy(group string, policy ProviderTypePolicy) []string {
	var models []string
	var abilities []Ability
	if err := DB.Model(&Ability{}).
		Where(commonGroupCol+" = ? and enabled = ?", group, true).
		Find(&abilities).Error; err != nil {
		return models
	}
	channelIds := make([]int, 0, len(abilities))
	for _, ability := range abilities {
		channelIds = append(channelIds, ability.ChannelId)
	}
	var channels []Channel
	if len(channelIds) > 0 {
		DB.Where("id IN ?", channelIds).Find(&channels)
	}
	channelById := make(map[int]Channel, len(channels))
	for _, channel := range channels {
		channelById[channel.Id] = channel
	}
	seen := make(map[string]struct{})
	for _, ability := range abilities {
		channel, ok := channelById[ability.ChannelId]
		if !ok || !IsChannelRoutableWithProviderPolicy(channel, policy) {
			continue
		}
		if _, ok := seen[ability.Model]; ok {
			continue
		}
		seen[ability.Model] = struct{}{}
		models = append(models, ability.Model)
	}
	return models
}

// GetGroupEnabledModelsExcludingExperimental returns enabled models for a group,
// excluding any model that is ONLY served by experimental_proxy channels.
// Used for non-internal users who must not see experimental_proxy offerings.
func GetGroupEnabledModelsExcludingExperimental(group string) []string {
	return GetGroupEnabledModelsWithProviderPolicy(group, ProviderTypePolicyFromAllowExperimental(false))
}

func GetEnabledModels() []string {
	var models []string
	// Find distinct models
	enabledChannelIds := routableChannelIDsSubquery(true)
	DB.Table("abilities").
		Where("enabled = ? AND channel_id IN (?)", true, enabledChannelIds).
		Distinct("model").
		Pluck("model", &models)
	return models
}

func GetAllEnableAbilities() []Ability {
	var abilities []Ability
	enabledChannelIds := routableChannelIDsSubquery(true)
	DB.Where("enabled = ? AND channel_id IN (?)", true, enabledChannelIds).Find(&abilities)
	return abilities
}

func GetChannel(group string, model string, retry int, allowExperimental bool) (*Channel, error) {
	return GetChannelWithProviderPolicy(group, model, retry, ProviderTypePolicyFromAllowExperimental(allowExperimental))
}

func GetChannelWithProviderPolicy(group string, model string, retry int, policy ProviderTypePolicy) (*Channel, error) {
	var abilities []Ability
	err := DB.Model(&Ability{}).
		Where(commonGroupCol+" = ? and model = ? and enabled = ?", group, model, true).
		Order("priority DESC, weight DESC").
		Find(&abilities).Error
	if err != nil {
		return nil, err
	}
	if len(abilities) == 0 {
		return nil, nil
	}

	channelIds := make([]int, 0, len(abilities))
	for _, ability := range abilities {
		channelIds = append(channelIds, ability.ChannelId)
	}
	var channels []Channel
	if err = DB.Where("id IN ?", channelIds).Find(&channels).Error; err != nil {
		return nil, err
	}
	channelById := make(map[int]Channel, len(channels))
	for _, channel := range channels {
		channelById[channel.Id] = channel
	}

	filtered := make([]Ability, 0, len(abilities))
	for _, ability := range abilities {
		channel, ok := channelById[ability.ChannelId]
		if !ok {
			continue
		}
		if IsChannelAllowedByProviderPolicy(channel, policy) && IsChannelRoutableWithProviderPolicy(channel, policy) {
			filtered = append(filtered, ability)
		}
	}
	if len(filtered) == 0 {
		return nil, nil
	}

	uniquePriorities := make([]int, 0)
	seenPriority := map[int]struct{}{}
	for _, ability := range filtered {
		priority := 0
		if ability.Priority != nil {
			priority = int(*ability.Priority)
		}
		if _, ok := seenPriority[priority]; !ok {
			seenPriority[priority] = struct{}{}
			uniquePriorities = append(uniquePriorities, priority)
		}
	}
	if retry >= len(uniquePriorities) {
		retry = len(uniquePriorities) - 1
	}
	targetPriority := uniquePriorities[retry]

	candidates := make([]Ability, 0, len(filtered))
	weightSum := uint(0)
	for _, ability := range filtered {
		priority := 0
		if ability.Priority != nil {
			priority = int(*ability.Priority)
		}
		if priority == targetPriority {
			candidates = append(candidates, ability)
			weightSum += ability.Weight + 10
		}
	}
	if len(candidates) == 0 {
		return nil, errors.New("database consistency error: no channel candidate after provider policy filtering")
	}

	weight := common.GetRandomInt(int(weightSum))
	selectedChannelID := candidates[0].ChannelId
	for _, ability := range candidates {
		weight -= int(ability.Weight) + 10
		if weight <= 0 {
			selectedChannelID = ability.ChannelId
			break
		}
	}
	channel, ok := channelById[selectedChannelID]
	if !ok {
		return nil, fmt.Errorf("database consistency error: channel #%d not found", selectedChannelID)
	}
	return &channel, nil
}

func routableChannelIDsSubquery(allowExperimental bool) *gorm.DB {
	query := DB.Table("channels").Select("id").Where("status = ?", common.ChannelStatusEnabled)
	if !allowExperimental {
		query = query.Where("(provider_type <> ? OR provider_type IS NULL OR provider_type = '')", constant.ProviderTypeExperimentalProxy)
	}
	return query
}

func IsChannelRoutable(channel Channel, allowExperimental bool) bool {
	return IsChannelRoutableWithProviderPolicy(channel, ProviderTypePolicyFromAllowExperimental(allowExperimental))
}

func IsChannelRoutableWithProviderPolicy(channel Channel, policy ProviderTypePolicy) bool {
	if channel.Status != common.ChannelStatusEnabled {
		return false
	}
	return policy.Allows(channel)
}

func (channel *Channel) AddAbilities(tx *gorm.DB) error {
	models_ := strings.Split(channel.Models, ",")
	groups_ := strings.Split(channel.Group, ",")
	abilitySet := make(map[string]struct{})
	abilities := make([]Ability, 0, len(models_))
	for _, model := range models_ {
		for _, group := range groups_ {
			key := group + "|" + model
			if _, exists := abilitySet[key]; exists {
				continue
			}
			abilitySet[key] = struct{}{}
			ability := Ability{
				Group:     group,
				Model:     model,
				ChannelId: channel.Id,
				Enabled:   channel.Status == common.ChannelStatusEnabled,
				Priority:  channel.Priority,
				Weight:    uint(channel.GetWeight()),
				Tag:       channel.Tag,
			}
			abilities = append(abilities, ability)
		}
	}
	if len(abilities) == 0 {
		return nil
	}
	// choose DB or provided tx
	useDB := DB
	if tx != nil {
		useDB = tx
	}
	for _, chunk := range lo.Chunk(abilities, 50) {
		err := useDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&chunk).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (channel *Channel) DeleteAbilities() error {
	return DB.Where("channel_id = ?", channel.Id).Delete(&Ability{}).Error
}

// UpdateAbilities updates abilities of this channel.
// Make sure the channel is completed before calling this function.
func (channel *Channel) UpdateAbilities(tx *gorm.DB) error {
	isNewTx := false
	// 如果没有传入事务，创建新的事务
	if tx == nil {
		tx = DB.Begin()
		if tx.Error != nil {
			return tx.Error
		}
		isNewTx = true
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
	}

	// First delete all abilities of this channel
	err := tx.Where("channel_id = ?", channel.Id).Delete(&Ability{}).Error
	if err != nil {
		if isNewTx {
			tx.Rollback()
		}
		return err
	}

	// Then add new abilities
	models_ := strings.Split(channel.Models, ",")
	groups_ := strings.Split(channel.Group, ",")
	abilitySet := make(map[string]struct{})
	abilities := make([]Ability, 0, len(models_))
	for _, model := range models_ {
		for _, group := range groups_ {
			key := group + "|" + model
			if _, exists := abilitySet[key]; exists {
				continue
			}
			abilitySet[key] = struct{}{}
			ability := Ability{
				Group:     group,
				Model:     model,
				ChannelId: channel.Id,
				Enabled:   channel.Status == common.ChannelStatusEnabled,
				Priority:  channel.Priority,
				Weight:    uint(channel.GetWeight()),
				Tag:       channel.Tag,
			}
			abilities = append(abilities, ability)
		}
	}

	if len(abilities) > 0 {
		for _, chunk := range lo.Chunk(abilities, 50) {
			err = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&chunk).Error
			if err != nil {
				if isNewTx {
					tx.Rollback()
				}
				return err
			}
		}
	}

	// 如果是新创建的事务，需要提交
	if isNewTx {
		return tx.Commit().Error
	}

	return nil
}

func UpdateAbilityStatus(channelId int, status bool) error {
	return DB.Model(&Ability{}).Where("channel_id = ?", channelId).Select("enabled").Update("enabled", status).Error
}

func UpdateAbilityStatusByTag(tag string, status bool) error {
	return DB.Model(&Ability{}).Where("tag = ?", tag).Select("enabled").Update("enabled", status).Error
}

func UpdateAbilityByTag(tag string, newTag *string, priority *int64, weight *uint) error {
	ability := Ability{}
	if newTag != nil {
		ability.Tag = newTag
	}
	if priority != nil {
		ability.Priority = priority
	}
	if weight != nil {
		ability.Weight = *weight
	}
	return DB.Model(&Ability{}).Where("tag = ?", tag).Updates(ability).Error
}

var fixLock = sync.Mutex{}

func FixAbility() (int, int, error) {
	lock := fixLock.TryLock()
	if !lock {
		return 0, 0, errors.New("已经有一个修复任务在运行中，请稍后再试")
	}
	defer fixLock.Unlock()

	// truncate abilities table
	if common.UsingSQLite {
		err := DB.Exec("DELETE FROM abilities").Error
		if err != nil {
			common.SysLog(fmt.Sprintf("Delete abilities failed: %s", err.Error()))
			return 0, 0, err
		}
	} else {
		err := DB.Exec("TRUNCATE TABLE abilities").Error
		if err != nil {
			common.SysLog(fmt.Sprintf("Truncate abilities failed: %s", err.Error()))
			return 0, 0, err
		}
	}
	var channels []*Channel
	// Find all channels
	err := DB.Model(&Channel{}).Find(&channels).Error
	if err != nil {
		return 0, 0, err
	}
	if len(channels) == 0 {
		return 0, 0, nil
	}
	successCount := 0
	failCount := 0
	for _, chunk := range lo.Chunk(channels, 50) {
		ids := lo.Map(chunk, func(c *Channel, _ int) int { return c.Id })
		// Delete all abilities of this channel
		err = DB.Where("channel_id IN ?", ids).Delete(&Ability{}).Error
		if err != nil {
			common.SysLog(fmt.Sprintf("Delete abilities failed: %s", err.Error()))
			failCount += len(chunk)
			continue
		}
		// Then add new abilities
		for _, channel := range chunk {
			err = channel.AddAbilities(nil)
			if err != nil {
				common.SysLog(fmt.Sprintf("Add abilities for channel %d failed: %s", channel.Id, err.Error()))
				failCount++
			} else {
				successCount++
			}
		}
	}
	InitChannelCache()
	return successCount, failCount, nil
}
