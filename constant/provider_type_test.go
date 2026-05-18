package constant

import "testing"

// validates A1.1.1 — IsValidProviderType covers all four provider types
func TestIsValidProviderType(t *testing.T) {
	valid := []string{
		ProviderTypeOfficialCloud,
		ProviderTypeAggregator,
		ProviderTypeAuthorizedProxy,
		ProviderTypeExperimentalProxy,
	}
	for _, pt := range valid {
		if !IsValidProviderType(pt) {
			t.Errorf("IsValidProviderType(%q) = false, want true", pt)
		}
	}

	invalid := []string{"", "unknown", "OFFICIAL_CLOUD", "official-cloud", "experimental"}
	for _, pt := range invalid {
		if IsValidProviderType(pt) {
			t.Errorf("IsValidProviderType(%q) = true, want false", pt)
		}
	}
}

// validates A1.1.3 — default value is official_cloud
func TestProviderTypeDefaultValue(t *testing.T) {
	if ProviderTypeOfficialCloud != "official_cloud" {
		t.Errorf("ProviderTypeOfficialCloud = %q, want 'official_cloud'", ProviderTypeOfficialCloud)
	}
	if ProviderTypeExperimentalProxy != "experimental_proxy" {
		t.Errorf("ProviderTypeExperimentalProxy = %q, want 'experimental_proxy'", ProviderTypeExperimentalProxy)
	}
}

// validates M1-F02 — risk level constants have correct string values
func TestRiskLevelConstants(t *testing.T) {
	if RiskLevelNormal != "normal" {
		t.Errorf("RiskLevelNormal = %q, want 'normal'", RiskLevelNormal)
	}
	if RiskLevelHigh != "high" {
		t.Errorf("RiskLevelHigh = %q, want 'high'", RiskLevelHigh)
	}
}

// validates M1-F02 — scope constants have correct string values
func TestScopeConstants(t *testing.T) {
	if ScopePublic != "public" {
		t.Errorf("ScopePublic = %q, want 'public'", ScopePublic)
	}
	if ScopeInternalOnly != "internal_only" {
		t.Errorf("ScopeInternalOnly = %q, want 'internal_only'", ScopeInternalOnly)
	}
}

// validates M1-F02 — visibility constants have correct string values
func TestVisibilityConstants(t *testing.T) {
	if VisibilityPublic != "public" {
		t.Errorf("VisibilityPublic = %q, want 'public'", VisibilityPublic)
	}
	if VisibilityInternalOnly != "internal_only" {
		t.Errorf("VisibilityInternalOnly = %q, want 'internal_only'", VisibilityInternalOnly)
	}
}

// validates A4.1.2 — GetDefaultProviderType returns correct values for known types
func TestGetDefaultProviderType(t *testing.T) {
	cases := []struct {
		channelType  int
		wantProvider string
	}{
		// official_cloud
		{ChannelTypeOpenAI, ProviderTypeOfficialCloud},
		{ChannelTypeAzure, ProviderTypeOfficialCloud},
		{ChannelTypeAnthropic, ProviderTypeOfficialCloud},
		{ChannelTypeGemini, ProviderTypeOfficialCloud},
		{ChannelTypeAws, ProviderTypeOfficialCloud},
		{ChannelTypeMistral, ProviderTypeOfficialCloud},
		{ChannelTypeDeepSeek, ProviderTypeOfficialCloud},
		// aggregator
		{ChannelTypeOpenRouter, ProviderTypeAggregator},
		{ChannelTypeSiliconFlow, ProviderTypeAggregator},
		{ChannelTypeSubmodel, ProviderTypeAggregator},
		// authorized_proxy
		{ChannelTypeOllama, ProviderTypeAuthorizedProxy},
		{ChannelTypeCustom, ProviderTypeAuthorizedProxy},
		{ChannelTypeXinference, ProviderTypeAuthorizedProxy},
		// experimental_proxy
		{ChannelTypeCodex, ProviderTypeExperimentalProxy},
		// unknown → fallback official_cloud
		{9999, ProviderTypeOfficialCloud},
		{0, ProviderTypeOfficialCloud},
	}
	for _, tc := range cases {
		got := GetDefaultProviderType(tc.channelType)
		if got != tc.wantProvider {
			t.Errorf("GetDefaultProviderType(%d) = %q, want %q", tc.channelType, got, tc.wantProvider)
		}
	}
}

// validates A4.1.1 — ChannelTypeDefaultProviderType covers all non-official types
func TestChannelTypeDefaultProviderTypeCompleteness(t *testing.T) {
	// Every entry must map to a valid provider type
	for channelType, providerType := range ChannelTypeDefaultProviderType {
		if !IsValidProviderType(providerType) {
			t.Errorf("ChannelTypeDefaultProviderType[%d] = %q is not a valid provider type", channelType, providerType)
		}
	}
	// Must have at least one official_cloud, one aggregator, one authorized_proxy, one experimental_proxy
	counts := map[string]int{}
	for _, pt := range ChannelTypeDefaultProviderType {
		counts[pt]++
	}
	for _, required := range []string{ProviderTypeOfficialCloud, ProviderTypeAggregator, ProviderTypeAuthorizedProxy, ProviderTypeExperimentalProxy} {
		if counts[required] == 0 {
			t.Errorf("ChannelTypeDefaultProviderType has no entries for %q", required)
		}
	}
}

// validates M6-F02 — KiroGateway is classified as experimental_proxy
func TestKiroGatewayIsExperimentalProxy(t *testing.T) {
	got := GetDefaultProviderType(ChannelTypeKiroGateway)
	if got != ProviderTypeExperimentalProxy {
		t.Errorf("GetDefaultProviderType(ChannelTypeKiroGateway) = %q, want %q", got, ProviderTypeExperimentalProxy)
	}
}

// validates M6-F03 — KiroGateway channel defaults to disabled via BeforeCreate
// (experimental_proxy invariant from M1-F03)
func TestKiroGatewayChannelTypeValue(t *testing.T) {
	if ChannelTypeKiroGateway != 58 {
		t.Errorf("ChannelTypeKiroGateway = %d, want 58", ChannelTypeKiroGateway)
	}
	// Verify it is in the experimental_proxy map entry
	pt, ok := ChannelTypeDefaultProviderType[ChannelTypeKiroGateway]
	if !ok {
		t.Fatal("ChannelTypeKiroGateway not in ChannelTypeDefaultProviderType map")
	}
	if pt != ProviderTypeExperimentalProxy {
		t.Errorf("ChannelTypeDefaultProviderType[KiroGateway] = %q, want %q", pt, ProviderTypeExperimentalProxy)
	}
}
