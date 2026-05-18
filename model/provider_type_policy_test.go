package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
)

func TestAUD021TokenAllowedProviderTypesParsingAndValidation(t *testing.T) {
	token := &Token{AllowedProviderTypes: "official_cloud, aggregator,official_cloud"}
	got := token.GetAllowedProviderTypes()
	if len(got) != 2 || got[0] != constant.ProviderTypeOfficialCloud || got[1] != constant.ProviderTypeAggregator {
		t.Fatalf("allowed provider types = %#v", got)
	}
	if err := token.ValidateAllowedProviderTypes(); err != nil {
		t.Fatalf("ValidateAllowedProviderTypes: %v", err)
	}
	bad := &Token{AllowedProviderTypes: "official_cloud,experimental"}
	if err := bad.ValidateAllowedProviderTypes(); err == nil {
		t.Fatal("expected invalid provider type rejection")
	}
}

func TestAUD021AllowedOfficialOnlySelectsOfficialCloud(t *testing.T) {
	withMemoryChannelCache(t)
	official := &Channel{Id: 2101, Type: constant.ChannelTypeOpenAI, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeOfficialCloud}
	aggregator := &Channel{Id: 2102, Type: constant.ChannelTypeOpenRouter, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeAggregator}
	experimental := &Channel{Id: 2103, Type: constant.ChannelTypeKiroGateway, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeExperimentalProxy}
	seedAUD021Memory("default", "aud021-official", official, aggregator, experimental)

	policy := NewProviderTypePolicy([]string{constant.ProviderTypeOfficialCloud}, false)
	for i := 0; i < 20; i++ {
		ch, err := GetRandomSatisfiedChannelWithProviderPolicy("default", "aud021-official", 0, policy)
		if err != nil {
			t.Fatalf("selection error: %v", err)
		}
		if ch == nil || ch.ProviderType != constant.ProviderTypeOfficialCloud {
			t.Fatalf("selected channel = %+v, want official_cloud", ch)
		}
	}
}

func TestAUD021AllowedAggregatorOnlyExcludesOfficialAndExperimental(t *testing.T) {
	withMemoryChannelCache(t)
	official := &Channel{Id: 2111, Type: constant.ChannelTypeOpenAI, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeOfficialCloud}
	aggregator := &Channel{Id: 2112, Type: constant.ChannelTypeOpenRouter, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeAggregator}
	experimental := &Channel{Id: 2113, Type: constant.ChannelTypeKiroGateway, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeExperimentalProxy}
	seedAUD021Memory("default", "aud021-aggregator", official, aggregator, experimental)

	policy := NewProviderTypePolicy([]string{constant.ProviderTypeAggregator}, true)
	for i := 0; i < 20; i++ {
		ch, err := GetRandomSatisfiedChannelWithProviderPolicy("default", "aud021-aggregator", 0, policy)
		if err != nil {
			t.Fatalf("selection error: %v", err)
		}
		if ch == nil || ch.ProviderType != constant.ProviderTypeAggregator {
			t.Fatalf("selected channel = %+v, want aggregator", ch)
		}
	}
}

func TestAUD021EmptyPolicyNormalCannotUseExperimental(t *testing.T) {
	withMemoryChannelCache(t)
	experimental := &Channel{Id: 2121, Type: constant.ChannelTypeKiroGateway, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeExperimentalProxy}
	seedAUD021Memory("default", "aud021-normal-empty", experimental)

	ch, err := GetRandomSatisfiedChannelWithProviderPolicy("default", "aud021-normal-empty", 0, NewProviderTypePolicy(nil, false))
	if err != nil {
		t.Fatalf("selection error: %v", err)
	}
	if ch != nil {
		t.Fatalf("normal empty policy selected experimental_proxy: %+v", ch)
	}
}

func TestAUD021InternalAllowExperimentalCanUseExperimental(t *testing.T) {
	withMemoryChannelCache(t)
	experimental := &Channel{Id: 2131, Type: constant.ChannelTypeKiroGateway, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeExperimentalProxy}
	seedAUD021Memory("internal", "aud021-internal-allow", experimental)

	ch, err := GetRandomSatisfiedChannelWithProviderPolicy("internal", "aud021-internal-allow", 0, NewProviderTypePolicy(nil, true))
	if err != nil {
		t.Fatalf("selection error: %v", err)
	}
	if ch == nil || ch.ProviderType != constant.ProviderTypeExperimentalProxy {
		t.Fatalf("selected channel = %+v, want experimental_proxy", ch)
	}
}

func TestAUD021InternalWithoutAllowExperimentalCannotUseExperimental(t *testing.T) {
	withMemoryChannelCache(t)
	experimental := &Channel{Id: 2141, Type: constant.ChannelTypeKiroGateway, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeExperimentalProxy}
	seedAUD021Memory("internal", "aud021-internal-deny", experimental)

	ch, err := GetRandomSatisfiedChannelWithProviderPolicy("internal", "aud021-internal-deny", 0, NewProviderTypePolicy(nil, false))
	if err != nil {
		t.Fatalf("selection error: %v", err)
	}
	if ch != nil {
		t.Fatalf("internal without allow_experimental selected experimental_proxy: %+v", ch)
	}
}

func TestAUD021OfficialFailureDoesNotFallbackToExperimental(t *testing.T) {
	truncateTables(t)
	clearAUD021Tables(t)
	seedAUD021Ability(t, 2151, common.ChannelStatusEnabled, constant.ProviderTypeOfficialCloud, constant.ChannelTypeOpenAI, "default", "aud021-fallback", 100)
	seedAUD021Ability(t, 2152, common.ChannelStatusEnabled, constant.ProviderTypeExperimentalProxy, constant.ChannelTypeKiroGateway, "default", "aud021-fallback", 10)

	policy := NewProviderTypePolicy([]string{constant.ProviderTypeOfficialCloud}, false)
	first, err := GetChannelWithProviderPolicy("default", "aud021-fallback", 0, policy)
	if err != nil {
		t.Fatalf("first selection error: %v", err)
	}
	if first == nil || first.Id != 2151 {
		t.Fatalf("first selection = %+v, want official channel", first)
	}
	second, err := GetChannelWithProviderPolicy("default", "aud021-fallback", 1, policy)
	if err != nil {
		t.Fatalf("retry selection error: %v", err)
	}
	if second == nil || second.Id != 2151 {
		t.Fatalf("retry/fallback selected %+v, want official channel and never experimental", second)
	}
}

func TestAUD021FallbackRetrySkipsDisallowedProviderType(t *testing.T) {
	truncateTables(t)
	clearAUD021Tables(t)
	seedAUD021Ability(t, 2161, common.ChannelStatusEnabled, constant.ProviderTypeOfficialCloud, constant.ChannelTypeOpenAI, "default", "aud021-retry", 100)
	seedAUD021Ability(t, 2162, common.ChannelStatusEnabled, constant.ProviderTypeAggregator, constant.ChannelTypeOpenRouter, "default", "aud021-retry", 10)

	policy := NewProviderTypePolicy([]string{constant.ProviderTypeAggregator}, false)
	ch, err := GetChannelWithProviderPolicy("default", "aud021-retry", 0, policy)
	if err != nil {
		t.Fatalf("selection error: %v", err)
	}
	if ch == nil || ch.ProviderType != constant.ProviderTypeAggregator {
		t.Fatalf("selected channel = %+v, want aggregator", ch)
	}
}

func TestAUD021DBAbilitySelectionDoesNotReturnDisallowedProviderType(t *testing.T) {
	truncateTables(t)
	clearAUD021Tables(t)
	seedAUD021Ability(t, 2171, common.ChannelStatusEnabled, constant.ProviderTypeOfficialCloud, constant.ChannelTypeOpenAI, "default", "aud021-db", 100)

	policy := NewProviderTypePolicy([]string{constant.ProviderTypeAggregator}, false)
	ch, err := GetChannelWithProviderPolicy("default", "aud021-db", 0, policy)
	if err != nil {
		t.Fatalf("selection error: %v", err)
	}
	if ch != nil {
		t.Fatalf("DB ability selection returned disallowed provider type: %+v", ch)
	}
}

func TestAUD021EmptyProviderTypeInfersExperimentalByChannelTypeAndRejects(t *testing.T) {
	withMemoryChannelCache(t)
	experimental := &Channel{Id: 2181, Type: constant.ChannelTypeKiroGateway, Status: common.ChannelStatusEnabled, ProviderType: ""}
	seedAUD021Memory("default", "aud021-empty-provider-type", experimental)

	ch, err := GetRandomSatisfiedChannelWithProviderPolicy("default", "aud021-empty-provider-type", 0, NewProviderTypePolicy(nil, false))
	if err != nil {
		t.Fatalf("selection error: %v", err)
	}
	if ch != nil {
		t.Fatalf("empty provider_type experimental channel was allowed: %+v", ch)
	}
}

func TestAUD021DisabledChannelStillUnavailable(t *testing.T) {
	withMemoryChannelCache(t)
	disabled := &Channel{Id: 2191, Type: constant.ChannelTypeOpenAI, Status: common.ChannelStatusManuallyDisabled, ProviderType: constant.ProviderTypeOfficialCloud}
	seedAUD021Memory("default", "aud021-disabled", disabled)

	ch, err := GetRandomSatisfiedChannelWithProviderPolicy("default", "aud021-disabled", 0, NewProviderTypePolicy([]string{constant.ProviderTypeOfficialCloud}, false))
	if err != nil {
		t.Fatalf("selection error: %v", err)
	}
	if ch != nil {
		t.Fatalf("disabled channel was selected: %+v", ch)
	}
}

func seedAUD021Memory(group, modelName string, channels ...*Channel) {
	channelSyncLock.Lock()
	defer channelSyncLock.Unlock()
	channelsIDM = make(map[int]*Channel, len(channels))
	ids := make([]int, 0, len(channels))
	for _, channel := range channels {
		channelsIDM[channel.Id] = channel
		ids = append(ids, channel.Id)
	}
	group2model2channels = map[string]map[string][]int{
		group: {modelName: ids},
	}
}

func seedAUD021Ability(t *testing.T, channelID int, status int, providerType string, channelType int, group string, modelName string, priority int64) {
	t.Helper()
	if err := DB.Create(&Channel{
		Id:           channelID,
		Type:         channelType,
		Key:          "sk-aud021-placeholder",
		Status:       status,
		Name:         modelName,
		ProviderType: providerType,
		Models:       modelName,
		Group:        group,
		Priority:     &priority,
	}).Error; err != nil {
		t.Fatalf("seed channel: %v", err)
	}
	if err := DB.Create(&Ability{
		Group:     group,
		Model:     modelName,
		ChannelId: channelID,
		Enabled:   true,
		Priority:  &priority,
		Weight:    100,
	}).Error; err != nil {
		t.Fatalf("seed ability: %v", err)
	}
}

func clearAUD021Tables(t *testing.T) {
	t.Helper()
	if err := DB.Exec("DELETE FROM abilities").Error; err != nil {
		t.Fatalf("clear abilities: %v", err)
	}
	if err := DB.Exec("DELETE FROM channels").Error; err != nil {
		t.Fatalf("clear channels: %v", err)
	}
}
