package model

import (
	"crypto/hmac"
	"errors"
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/bytedance/gopkg/util/gopool"
	"gorm.io/gorm"
)

type Token struct {
	Id                 int     `json:"id"`
	UserId             int     `json:"user_id" gorm:"index"`
	Key                string  `json:"key" gorm:"type:varchar(128);index"`
	Status             int     `json:"status" gorm:"default:1"`
	Name               string  `json:"name" gorm:"index" `
	CreatedTime        int64   `json:"created_time" gorm:"bigint"`
	AccessedTime       int64   `json:"accessed_time" gorm:"bigint"`
	ExpiredTime        int64   `json:"expired_time" gorm:"bigint;default:-1"` // -1 means never expired
	RemainQuota        int     `json:"remain_quota" gorm:"default:0"`
	UnlimitedQuota     bool    `json:"unlimited_quota"`
	ModelLimitsEnabled bool    `json:"model_limits_enabled"`
	ModelLimits        string  `json:"model_limits" gorm:"type:text"`
	AllowIps           *string `json:"allow_ips" gorm:"default:''"`
	UsedQuota          int     `json:"used_quota" gorm:"default:0"` // used quota
	Group              string  `json:"group" gorm:"default:''"`
	CrossGroupRetry    bool    `json:"cross_group_retry"` // 跨分组重试，仅auto分组有效
	// M10-F01: org/project binding
	OrgId     *int `json:"org_id" gorm:"index"`
	ProjectId *int `json:"project_id" gorm:"index"`
	// M10-F02: secure key storage — hash for verification, prefix for display
	KeyHash   string `json:"-" gorm:"type:varchar(64);index"`
	KeyPrefix string `json:"key_prefix" gorm:"type:varchar(16)"`
	// M10-F04: experimental_proxy access control (default false — normal keys cannot access experimental)
	AllowExperimental    bool           `json:"allow_experimental" gorm:"default:false"`
	AllowedProviderTypes string         `json:"allowed_provider_types" gorm:"type:text"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

type TokenTenantScope struct {
	UserId    int
	OrgId     *int
	ProjectId *int
	Legacy    bool
}

func (token *Token) Clean() {
	token.Key = ""
}

func HashTokenKey(key string) string {
	key = NormalizeTokenKey(key)
	if key == "" {
		return ""
	}
	return common.GenerateHMAC(key)
}

func NormalizeTokenKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.TrimPrefix(key, "Bearer ")
	key = strings.TrimPrefix(key, "bearer ")
	return strings.TrimPrefix(key, "sk-")
}

func looksLikeTokenHash(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, r := range value {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}

func tokenKeyColumnName() string {
	if commonKeyCol != "" {
		return commonKeyCol
	}
	return "`key`"
}

func tokenHashForStoredValue(value string) string {
	value = NormalizeTokenKey(value)
	if value == "" {
		return ""
	}
	if looksLikeTokenHash(value) {
		return value
	}
	return HashTokenKey(value)
}

func tokenKeyPrefix(key string) string {
	key = NormalizeTokenKey(key)
	if len(key) <= 8 {
		return key
	}
	return key[:8]
}

func prepareTokenKeyForStorage(token *Token) {
	if token == nil {
		return
	}
	rawKey := NormalizeTokenKey(token.Key)
	if token.KeyHash == "" && rawKey != "" {
		token.KeyHash = HashTokenKey(rawKey)
	}
	if token.KeyPrefix == "" && rawKey != "" {
		token.KeyPrefix = tokenKeyPrefix(rawKey)
	}
	if token.KeyHash != "" {
		// Keep the deprecated legacy column non-secret for databases that still
		// have an existing uniqueness/index dependency on tokens.key.
		token.Key = token.KeyHash
	}
}

// BeforeCreate computes non-reversible key storage fields and removes plaintext.
func (token *Token) BeforeCreate(tx *gorm.DB) error {
	prepareTokenKeyForStorage(token)
	return nil
}

func intPtr(v int) *int {
	value := v
	return &value
}

func isUserInOrganization(userId int, org *Organization) (bool, error) {
	if org == nil {
		return false, nil
	}
	if org.OwnerId == userId {
		return true, nil
	}
	var count int64
	if err := DB.Model(&OrganizationMember{}).Where("org_id = ? AND user_id = ?", org.Id, userId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func ResolveTokenTenantScope(token *Token) (*TokenTenantScope, error) {
	if token == nil {
		return nil, errors.New("token is nil")
	}
	scope := &TokenTenantScope{
		UserId:    token.UserId,
		OrgId:     token.OrgId,
		ProjectId: token.ProjectId,
		Legacy:    token.OrgId == nil && token.ProjectId == nil,
	}
	if scope.Legacy {
		return scope, nil
	}

	var projectOrgId *int
	if token.ProjectId != nil {
		project, err := GetProjectById(*token.ProjectId)
		if err != nil {
			return nil, fmt.Errorf("invalid token project binding: %w", err)
		}
		if project.Status != common.UserStatusEnabled {
			return nil, errors.New("token project binding is disabled")
		}
		projectOrgId = intPtr(project.OrgId)
		if token.OrgId != nil && project.OrgId != *token.OrgId {
			return nil, errors.New("token project binding belongs to a different organization")
		}
		if scope.OrgId == nil {
			scope.OrgId = projectOrgId
		}
	}

	if scope.OrgId != nil {
		org, err := GetOrganizationById(*scope.OrgId)
		if err != nil {
			return nil, fmt.Errorf("invalid token organization binding: %w", err)
		}
		if org.Status != common.UserStatusEnabled {
			return nil, errors.New("token organization binding is disabled")
		}
		member, err := isUserInOrganization(token.UserId, org)
		if err != nil {
			return nil, err
		}
		if !member {
			return nil, errors.New("token user is not a member of the bound organization")
		}
	}

	return scope, nil
}

func ValidateTokenTenantBinding(token *Token) error {
	_, err := ResolveTokenTenantScope(token)
	return err
}

// DisableToken sets a token's status to disabled.
func DisableToken(id int) error {
	return DB.Model(&Token{}).Where("id = ?", id).Update("status", common.TokenStatusDisabled).Error
}

func MaskTokenKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 4 {
		return strings.Repeat("*", len(key))
	}
	if len(key) <= 8 {
		return key[:2] + "****" + key[len(key)-2:]
	}
	return key[:4] + "**********" + key[len(key)-4:]
}

func (token *Token) GetFullKey() string {
	return ""
}

func (token *Token) GetMaskedKey() string {
	return MaskTokenKey(token.KeyPrefix)
}

func (token *Token) GetIpLimits() []string {
	// delete empty spaces
	//split with \n
	ipLimits := make([]string, 0)
	if token.AllowIps == nil {
		return ipLimits
	}
	cleanIps := strings.ReplaceAll(*token.AllowIps, " ", "")
	if cleanIps == "" {
		return ipLimits
	}
	ips := strings.Split(cleanIps, "\n")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		ip = strings.ReplaceAll(ip, ",", "")
		if ip != "" {
			ipLimits = append(ipLimits, ip)
		}
	}
	return ipLimits
}

func GetAllUserTokens(userId int, startIdx int, num int) ([]*Token, error) {
	var tokens []*Token
	var err error
	err = DB.Where("user_id = ?", userId).Order("id desc").Limit(num).Offset(startIdx).Find(&tokens).Error
	return tokens, err
}

// GetAdminAllTokens returns a paginated list of all tokens across all users.
// M15-F01: supports optional filters for userId, orgId, projectId, allowExperimental.
func GetAdminAllTokens(userId int, orgId *int, projectId *int, allowExperimental *bool, startIdx int, num int) (tokens []*Token, total int64, err error) {
	tx := DB.Model(&Token{})
	if userId != 0 {
		tx = tx.Where("user_id = ?", userId)
	}
	if orgId != nil {
		tx = tx.Where("org_id = ?", *orgId)
	}
	if projectId != nil {
		tx = tx.Where("project_id = ?", *projectId)
	}
	if allowExperimental != nil {
		tx = tx.Where("allow_experimental = ?", *allowExperimental)
	}
	err = tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = tx.Order("id desc").Limit(num).Offset(startIdx).Find(&tokens).Error
	return tokens, total, err
}

// sanitizeLikePattern 校验并清洗用户输入的 LIKE 搜索模式。
// 规则：
//  1. 转义 ! 和 _（使用 ! 作为 ESCAPE 字符，兼容 MySQL/PostgreSQL/SQLite）
//  2. 连续的 % 合并为单个 %
//  3. 最多允许 2 个 %
//  4. 含 % 时（模糊搜索），去掉 % 后关键词长度必须 >= 2
//  5. 不含 % 时按精确匹配
func sanitizeLikePattern(input string) (string, error) {
	// 1. 先转义 ESCAPE 字符 ! 自身，再转义 _
	//    使用 ! 而非 \ 作为 ESCAPE 字符，避免 MySQL 中反斜杠的字符串转义问题
	input = strings.ReplaceAll(input, "!", "!!")
	input = strings.ReplaceAll(input, `_`, `!_`)

	// 2. 连续的 % 直接拒绝
	if strings.Contains(input, "%%") {
		return "", errors.New("搜索模式中不允许包含连续的 % 通配符")
	}

	// 3. 统计 % 数量，不得超过 2
	count := strings.Count(input, "%")
	if count > 2 {
		return "", errors.New("搜索模式中最多允许包含 2 个 % 通配符")
	}

	// 4. 含 % 时，去掉 % 后关键词长度必须 >= 2
	if count > 0 {
		stripped := strings.ReplaceAll(input, "%", "")
		if len(stripped) < 2 {
			return "", errors.New("使用模糊搜索时，关键词长度至少为 2 个字符")
		}
		return input, nil
	}

	// 5. 无 % 时，精确全匹配
	return input, nil
}

const searchHardLimit = 100

func SearchUserTokens(userId int, keyword string, token string, offset int, limit int) (tokens []*Token, total int64, err error) {
	// model 层强制截断
	if limit <= 0 || limit > searchHardLimit {
		limit = searchHardLimit
	}
	if offset < 0 {
		offset = 0
	}

	if token != "" {
		token = NormalizeTokenKey(token)
	}

	// 超量用户（令牌数超过上限）只允许精确搜索，禁止模糊搜索
	maxTokens := operation_setting.GetMaxUserTokens()
	hasFuzzy := strings.Contains(keyword, "%") || strings.Contains(token, "%")
	if hasFuzzy {
		count, err := CountUserTokens(userId)
		if err != nil {
			common.SysLog("failed to count user tokens: " + err.Error())
			return nil, 0, errors.New("获取令牌数量失败")
		}
		if int(count) > maxTokens {
			return nil, 0, errors.New("令牌数量超过上限，仅允许精确搜索，请勿使用 % 通配符")
		}
	}

	baseQuery := DB.Model(&Token{}).Where("user_id = ?", userId)

	// 非空才加 LIKE 条件，空则跳过（不过滤该字段）
	if keyword != "" {
		keywordPattern, err := sanitizeLikePattern(keyword)
		if err != nil {
			return nil, 0, err
		}
		baseQuery = baseQuery.Where("name LIKE ? ESCAPE '!'", keywordPattern)
	}
	if token != "" {
		tokenPattern, err := sanitizeLikePattern(token)
		if err != nil {
			return nil, 0, err
		}
		baseQuery = baseQuery.Where("key_prefix LIKE ? ESCAPE '!'", tokenPattern)
	}

	// 先查匹配总数（用于分页，受 maxTokens 上限保护，避免全表 COUNT）
	err = baseQuery.Limit(maxTokens).Count(&total).Error
	if err != nil {
		common.SysError("failed to count search tokens: " + err.Error())
		return nil, 0, errors.New("搜索令牌失败")
	}

	// 再分页查数据
	err = baseQuery.Order("id desc").Offset(offset).Limit(limit).Find(&tokens).Error
	if err != nil {
		common.SysError("failed to search tokens: " + err.Error())
		return nil, 0, errors.New("搜索令牌失败")
	}
	return tokens, total, nil
}

func ValidateUserToken(key string) (token *Token, err error) {
	if key == "" {
		return nil, ErrTokenNotProvided
	}
	token, err = GetTokenByKey(key, false)
	if err == nil {
		if token.Status == common.TokenStatusExhausted ||
			token.Status == common.TokenStatusExpired ||
			token.Status != common.TokenStatusEnabled {
			return token, ErrTokenInvalid
		}
		if token.ExpiredTime != -1 && token.ExpiredTime < common.GetTimestamp() {
			if !common.RedisEnabled {
				token.Status = common.TokenStatusExpired
				err := token.SelectUpdate()
				if err != nil {
					common.SysLog("failed to update token status" + err.Error())
				}
			}
			return token, ErrTokenInvalid
		}
		if !token.UnlimitedQuota && token.RemainQuota <= 0 {
			if !common.RedisEnabled {
				token.Status = common.TokenStatusExhausted
				err := token.SelectUpdate()
				if err != nil {
					common.SysLog("failed to update token status" + err.Error())
				}
			}
			return token, ErrTokenInvalid
		}
		if err := ValidateTokenTenantBinding(token); err != nil {
			common.SysLog(fmt.Sprintf("ValidateUserToken: invalid tenant binding for token %d: %v", token.Id, err))
			return token, ErrTokenInvalid
		}
		return token, nil
	}
	common.SysLog("ValidateUserToken: failed to get token: " + err.Error())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTokenInvalid
	}
	return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
}

func GetTokenByIds(id int, userId int) (*Token, error) {
	if id == 0 || userId == 0 {
		return nil, errors.New("id 或 userId 为空！")
	}
	token := Token{Id: id, UserId: userId}
	var err error = nil
	err = DB.First(&token, "id = ? and user_id = ?", id, userId).Error
	return &token, err
}

func GetTokenById(id int) (*Token, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	token := Token{Id: id}
	var err error = nil
	err = DB.First(&token, "id = ?", id).Error
	if shouldUpdateRedis(true, err) {
		gopool.Go(func() {
			if err := cacheSetToken(token); err != nil {
				common.SysLog("failed to update user status cache: " + err.Error())
			}
		})
	}
	return &token, err
}

func GetTokenByKey(key string, fromDB bool) (token *Token, err error) {
	keyHash := HashTokenKey(key)
	if keyHash == "" {
		return nil, gorm.ErrRecordNotFound
	}
	defer func() {
		// Update Redis cache asynchronously on successful DB read
		if shouldUpdateRedis(fromDB, err) && token != nil {
			gopool.Go(func() {
				if err := cacheSetToken(*token); err != nil {
					common.SysLog("failed to update user status cache: " + err.Error())
				}
			})
		}
	}()
	if !fromDB && common.RedisEnabled {
		// Try Redis first
		token, err := cacheGetTokenByKey(key)
		if err == nil {
			return token, nil
		}
		// Don't return error - fall through to DB
	}
	fromDB = true
	err = DB.Where("key_hash = ?", keyHash).First(&token).Error
	if err == nil && !hmac.Equal([]byte(token.KeyHash), []byte(keyHash)) {
		return nil, gorm.ErrRecordNotFound
	}
	return token, err
}

func (token *Token) Insert() error {
	var err error
	prepareTokenKeyForStorage(token)
	err = DB.Create(token).Error
	return err
}

// Update Make sure your token's fields is completed, because this will update non-zero values
func (token *Token) Update() (err error) {
	defer func() {
		if shouldUpdateRedis(true, err) {
			gopool.Go(func() {
				err := cacheSetToken(*token)
				if err != nil {
					common.SysLog("failed to update token cache: " + err.Error())
				}
			})
		}
	}()
	err = DB.Model(token).Select("name", "status", "expired_time", "remain_quota", "unlimited_quota",
		"model_limits_enabled", "model_limits", "allow_ips", "group", "cross_group_retry", "org_id", "project_id",
		"allow_experimental", "allowed_provider_types").Updates(token).Error
	return err
}

func (token *Token) SelectUpdate() (err error) {
	defer func() {
		if shouldUpdateRedis(true, err) {
			gopool.Go(func() {
				err := cacheSetToken(*token)
				if err != nil {
					common.SysLog("failed to update token cache: " + err.Error())
				}
			})
		}
	}()
	// This can update zero values
	return DB.Model(token).Select("accessed_time", "status").Updates(token).Error
}

func (token *Token) Delete() (err error) {
	cacheKey := token.KeyHash
	defer func() {
		if shouldUpdateRedis(true, err) {
			gopool.Go(func() {
				err := cacheDeleteToken(cacheKey)
				if err != nil {
					common.SysLog("failed to delete token cache: " + err.Error())
				}
			})
		}
	}()
	err = DB.Delete(token).Error
	return err
}

func (token *Token) IsModelLimitsEnabled() bool {
	return token.ModelLimitsEnabled
}

func (token *Token) GetModelLimits() []string {
	if token.ModelLimits == "" {
		return []string{}
	}
	return strings.Split(token.ModelLimits, ",")
}

func (token *Token) GetModelLimitsMap() map[string]bool {
	limits := token.GetModelLimits()
	limitsMap := make(map[string]bool)
	for _, limit := range limits {
		limitsMap[limit] = true
	}
	return limitsMap
}

func (token *Token) GetAllowedProviderTypes() []string {
	if token == nil || strings.TrimSpace(token.AllowedProviderTypes) == "" {
		return nil
	}
	parts := strings.Split(token.AllowedProviderTypes, ",")
	allowed := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		providerType := strings.TrimSpace(part)
		if providerType == "" {
			continue
		}
		if _, exists := seen[providerType]; exists {
			continue
		}
		seen[providerType] = struct{}{}
		allowed = append(allowed, providerType)
	}
	return allowed
}

func (token *Token) ValidateAllowedProviderTypes() error {
	for _, providerType := range token.GetAllowedProviderTypes() {
		if !constant.IsValidProviderType(providerType) {
			return fmt.Errorf("invalid allowed_provider_types value %q", providerType)
		}
	}
	return nil
}

func DisableModelLimits(tokenId int) error {
	token, err := GetTokenById(tokenId)
	if err != nil {
		return err
	}
	token.ModelLimitsEnabled = false
	token.ModelLimits = ""
	return token.Update()
}

func DeleteTokenById(id int, userId int) (err error) {
	// Why we need userId here? In case user want to delete other's token.
	if id == 0 || userId == 0 {
		return errors.New("id 或 userId 为空！")
	}
	token := Token{Id: id, UserId: userId}
	err = DB.Where(token).First(&token).Error
	if err != nil {
		return err
	}
	return token.Delete()
}

func IncreaseTokenQuota(tokenId int, key string, quota int) (err error) {
	if quota < 0 {
		return errors.New("quota 不能为负数！")
	}
	if common.RedisEnabled {
		gopool.Go(func() {
			err := cacheIncrTokenQuota(HashTokenKey(key), int64(quota))
			if err != nil {
				common.SysLog("failed to increase token quota: " + err.Error())
			}
		})
	}
	if common.BatchUpdateEnabled {
		addNewRecord(BatchUpdateTypeTokenQuota, tokenId, quota)
		return nil
	}
	return increaseTokenQuota(tokenId, quota)
}

func increaseTokenQuota(id int, quota int) (err error) {
	err = DB.Model(&Token{}).Where("id = ?", id).Updates(
		map[string]interface{}{
			"remain_quota":  gorm.Expr("remain_quota + ?", quota),
			"used_quota":    gorm.Expr("used_quota - ?", quota),
			"accessed_time": common.GetTimestamp(),
		},
	).Error
	return err
}

func DecreaseTokenQuota(id int, key string, quota int) (err error) {
	if quota < 0 {
		return errors.New("quota 不能为负数！")
	}
	if common.RedisEnabled {
		gopool.Go(func() {
			err := cacheDecrTokenQuota(HashTokenKey(key), int64(quota))
			if err != nil {
				common.SysLog("failed to decrease token quota: " + err.Error())
			}
		})
	}
	if common.BatchUpdateEnabled {
		addNewRecord(BatchUpdateTypeTokenQuota, id, -quota)
		return nil
	}
	return decreaseTokenQuota(id, quota)
}

func decreaseTokenQuota(id int, quota int) (err error) {
	err = DB.Model(&Token{}).Where("id = ?", id).Updates(
		map[string]interface{}{
			"remain_quota":  gorm.Expr("remain_quota - ?", quota),
			"used_quota":    gorm.Expr("used_quota + ?", quota),
			"accessed_time": common.GetTimestamp(),
		},
	).Error
	return err
}

// CountUserTokens returns total number of tokens for the given user, used for pagination
func CountUserTokens(userId int) (int64, error) {
	var total int64
	err := DB.Model(&Token{}).Where("user_id = ?", userId).Count(&total).Error
	return total, err
}

// BatchDeleteTokens 删除指定用户的一组令牌，返回成功删除数量
func BatchDeleteTokens(ids []int, userId int) (int, error) {
	if len(ids) == 0 {
		return 0, errors.New("ids 不能为空！")
	}

	tx := DB.Begin()

	var tokens []Token
	if err := tx.Where("user_id = ? AND id IN (?)", userId, ids).Find(&tokens).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Where("user_id = ? AND id IN (?)", userId, ids).Delete(&Token{}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	if common.RedisEnabled {
		gopool.Go(func() {
			for _, t := range tokens {
				_ = cacheDeleteToken(t.KeyHash)
			}
		})
	}

	return len(tokens), nil
}

func GetTokenKeysByIds(ids []int, userId int) ([]Token, error) {
	var tokens []Token
	err := DB.Select("id", "key_prefix").
		Where("user_id = ? AND id IN (?)", userId, ids).
		Find(&tokens).Error
	return tokens, err
}

// InvalidateUserTokensCache 清理指定用户所有令牌在 Redis 中的缓存，
// 配合 InvalidateUserCache 使用，可在用户被禁用/删除时立即阻断其令牌的请求。
// 下一次请求将从数据库重新加载令牌及用户状态，从而立即识别出被禁用的用户。
func InvalidateUserTokensCache(userId int) error {
	if !common.RedisEnabled {
		return nil
	}
	if userId <= 0 {
		return errors.New("userId 无效")
	}
	var tokens []Token
	if err := DB.Unscoped().
		Select("id", "key_hash").
		Where("user_id = ?", userId).
		Find(&tokens).Error; err != nil {
		return err
	}
	var firstErr error
	for _, t := range tokens {
		if t.KeyHash == "" {
			continue
		}
		if err := cacheDeleteToken(t.KeyHash); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func MigratePlaintextTokenKeys() error {
	var tokens []Token
	if err := DB.Unscoped().
		Where(tokenKeyColumnName() + " <> '' AND " + tokenKeyColumnName() + " IS NOT NULL").
		Find(&tokens).Error; err != nil {
		return err
	}
	for i := range tokens {
		tok := &tokens[i]
		legacyKey := NormalizeTokenKey(tok.Key)
		if legacyKey == "" {
			continue
		}
		if looksLikeTokenHash(legacyKey) && tok.KeyHash == legacyKey {
			continue
		}
		if tok.KeyHash == "" {
			tok.KeyHash = tokenHashForStoredValue(legacyKey)
		}
		if tok.KeyPrefix == "" && !looksLikeTokenHash(legacyKey) {
			tok.KeyPrefix = tokenKeyPrefix(legacyKey)
		}
		if err := DB.Unscoped().Model(&Token{}).Where("id = ?", tok.Id).Updates(map[string]interface{}{
			"key":        tok.KeyHash,
			"key_hash":   tok.KeyHash,
			"key_prefix": tok.KeyPrefix,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
