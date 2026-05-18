package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
)

func withMemoryChannelCache(t *testing.T) {
	t.Helper()

	oldMemoryCacheEnabled := common.MemoryCacheEnabled
	channelSyncLock.RLock()
	oldChannelsIDM := channelsIDM
	oldGroup2Model2Channels := group2model2channels
	channelSyncLock.RUnlock()
	common.MemoryCacheEnabled = true
	t.Cleanup(func() {
		common.MemoryCacheEnabled = oldMemoryCacheEnabled
		channelSyncLock.Lock()
		channelsIDM = oldChannelsIDM
		group2model2channels = oldGroup2Model2Channels
		channelSyncLock.Unlock()
	})
}

// validates M8-F01 — experimental_proxy channels are filtered from candidate list
// when allowExperimental=false (in-memory path, no DB required).
func TestGetRandomSatisfiedChannelFiltersExperimental(t *testing.T) {
	withMemoryChannelCache(t)

	// Build a minimal in-memory cache with one official_cloud and one experimental_proxy channel.
	officialCh := &Channel{Id: 1, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeOfficialCloud}
	experimentalCh := &Channel{Id: 2, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeExperimentalProxy}

	channelSyncLock.Lock()
	channelsIDM = map[int]*Channel{1: officialCh, 2: experimentalCh}
	group2model2channels = map[string]map[string][]int{
		"default": {"gpt-4o": {1, 2}},
	}
	channelSyncLock.Unlock()

	// allowExperimental=false: only official channel should be eligible.
	for i := 0; i < 20; i++ {
		ch, err := GetRandomSatisfiedChannel("default", "gpt-4o", 0, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ch == nil {
			t.Fatal("expected a channel, got nil")
		}
		if ch.ProviderType == constant.ProviderTypeExperimentalProxy {
			t.Errorf("iteration %d: experimental_proxy channel returned when allowExperimental=false", i)
		}
	}
}

// validates M8-F02 — official_cloud failure cannot fallback to experimental_proxy
// because experimental_proxy is not in the candidate pool when allowExperimental=false.
func TestGetRandomSatisfiedChannelNoFallbackToExperimental(t *testing.T) {
	withMemoryChannelCache(t)

	// Only experimental_proxy channels available.
	experimentalCh := &Channel{Id: 3, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeExperimentalProxy}

	channelSyncLock.Lock()
	channelsIDM = map[int]*Channel{3: experimentalCh}
	group2model2channels = map[string]map[string][]int{
		"default": {"experimental-model": {3}},
	}
	channelSyncLock.Unlock()

	// allowExperimental=false: no channels should be returned (nil, not error).
	ch, err := GetRandomSatisfiedChannel("default", "experimental-model", 0, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch != nil {
		t.Errorf("expected nil channel when only experimental_proxy available and allowExperimental=false, got %+v", ch)
	}
}

// validates M8-F03 — explicit opt-in: allowExperimental=true returns experimental_proxy channels.
func TestGetRandomSatisfiedChannelAllowsExperimentalWhenOptIn(t *testing.T) {
	withMemoryChannelCache(t)

	experimentalCh := &Channel{Id: 4, Status: common.ChannelStatusEnabled, ProviderType: constant.ProviderTypeExperimentalProxy}

	channelSyncLock.Lock()
	channelsIDM = map[int]*Channel{4: experimentalCh}
	group2model2channels = map[string]map[string][]int{
		"internal": {"kiro-model": {4}},
	}
	channelSyncLock.Unlock()

	ch, err := GetRandomSatisfiedChannel("internal", "kiro-model", 0, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch == nil {
		t.Fatal("expected experimental_proxy channel when allowExperimental=true, got nil")
	}
	if ch.ProviderType != constant.ProviderTypeExperimentalProxy {
		t.Errorf("expected experimental_proxy, got %q", ch.ProviderType)
	}
}

// validates AUD-017 — normal users cannot route to disabled experimental_proxy channels.
func TestAUD017NormalUserCannotCallDisabledExperimentalProxy(t *testing.T) {
	truncateTables(t)

	seedAUD017Ability(t, 17, common.ChannelStatusManuallyDisabled, constant.ProviderTypeExperimentalProxy, "default", "aud017-disabled-normal", 100)

	ch, err := GetChannel("default", "aud017-disabled-normal", 0, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch != nil {
		t.Fatalf("normal selection returned disabled experimental_proxy channel: %+v", ch)
	}
}

// validates AUD-017 — internal users with experimental opt-in still cannot route disabled experimental_proxy.
func TestAUD017InternalUserCannotCallDisabledExperimentalProxy(t *testing.T) {
	truncateTables(t)

	seedAUD017Ability(t, 18, common.ChannelStatusManuallyDisabled, constant.ProviderTypeExperimentalProxy, "internal", "aud017-disabled-internal", 100)

	ch, err := GetChannel("internal", "aud017-disabled-internal", 0, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch != nil {
		t.Fatalf("internal selection returned disabled experimental_proxy channel: %+v", ch)
	}
}

// validates AUD-017 — admin ordinary path does not make disabled experimental_proxy routable.
func TestAUD017AdminOrdinaryPathCannotCallDisabledExperimentalProxy(t *testing.T) {
	truncateTables(t)

	seedAUD017Ability(t, 19, common.ChannelStatusManuallyDisabled, constant.ProviderTypeExperimentalProxy, "default", "aud017-disabled-admin", 100)

	ch, err := GetChannel("default", "aud017-disabled-admin", 0, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch != nil {
		t.Fatalf("admin ordinary selection returned disabled experimental_proxy channel: %+v", ch)
	}
}

// validates AUD-017 — DB ability selection excludes disabled experimental_proxy channels.
func TestAUD017DBAbilitySelectionDoesNotReturnDisabledExperimentalProxy(t *testing.T) {
	truncateTables(t)

	seedAUD017Ability(t, 20, common.ChannelStatusManuallyDisabled, constant.ProviderTypeExperimentalProxy, "default", "aud017-db-ability", 100)

	models := GetGroupEnabledModels("default")
	for _, modelName := range models {
		if modelName == "aud017-db-ability" {
			t.Fatalf("GetGroupEnabledModels returned model served only by disabled experimental_proxy")
		}
	}
	abilities := GetAllEnableAbilities()
	for _, ability := range abilities {
		if ability.ChannelId == 20 {
			t.Fatalf("GetAllEnableAbilities returned disabled experimental_proxy ability: %+v", ability)
		}
	}
}

// validates AUD-017 — fallback skips disabled experimental_proxy and selects an enabled channel.
func TestAUD017FallbackDoesNotReturnDisabledExperimentalProxy(t *testing.T) {
	truncateTables(t)

	seedAUD017Ability(t, 21, common.ChannelStatusManuallyDisabled, constant.ProviderTypeExperimentalProxy, "default", "aud017-fallback", 100)
	seedAUD017Ability(t, 22, common.ChannelStatusEnabled, constant.ProviderTypeOfficialCloud, "default", "aud017-fallback", 10)

	ch, err := GetChannel("default", "aud017-fallback", 0, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch == nil {
		t.Fatal("expected enabled official_cloud fallback channel")
	}
	if ch.Id != 22 {
		t.Fatalf("fallback returned channel #%d, want enabled official_cloud #22", ch.Id)
	}
}

func TestAUD017LegacyMemoryCandidateDoesNotReturnDisabledExperimentalProxy(t *testing.T) {
	withMemoryChannelCache(t)

	disabledExperimental := &Channel{Id: 23, Status: common.ChannelStatusManuallyDisabled, ProviderType: constant.ProviderTypeExperimentalProxy}

	channelSyncLock.Lock()
	channelsIDM = map[int]*Channel{23: disabledExperimental}
	group2model2channels = map[string]map[string][]int{
		"internal": {"aud017-legacy": {23}},
	}
	channelSyncLock.Unlock()

	ch, err := GetRandomSatisfiedChannel("internal", "aud017-legacy", 0, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch != nil {
		t.Fatalf("legacy memory selection returned disabled experimental_proxy channel: %+v", ch)
	}
}

func seedAUD017Ability(t *testing.T, channelID int, status int, providerType string, group string, modelName string, priority int64) {
	t.Helper()

	if err := DB.Create(&Channel{
		Id:           channelID,
		Type:         constant.ChannelTypeOpenAI,
		Key:          "sk-aud017",
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
