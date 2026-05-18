package model

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
)

type ProviderTypePolicy struct {
	AllowedProviderTypes []string
	AllowExperimental    bool
}

func NewProviderTypePolicy(allowedProviderTypes []string, allowExperimental bool) ProviderTypePolicy {
	policy := ProviderTypePolicy{
		AllowedProviderTypes: normalizeProviderTypes(allowedProviderTypes),
		AllowExperimental:    allowExperimental,
	}
	return policy
}

func ProviderTypePolicyFromAllowExperimental(allowExperimental bool) ProviderTypePolicy {
	return NewProviderTypePolicy(nil, allowExperimental)
}

func (policy ProviderTypePolicy) Allows(channel Channel) bool {
	providerType, ok := channel.EffectiveProviderType()
	if !ok {
		return false
	}
	if !constant.IsValidProviderType(providerType) {
		return false
	}
	if providerType == constant.ProviderTypeExperimentalProxy && !policy.AllowExperimental {
		return false
	}
	if len(policy.AllowedProviderTypes) == 0 {
		return providerType != constant.ProviderTypeExperimentalProxy || policy.AllowExperimental
	}
	for _, allowed := range policy.AllowedProviderTypes {
		if providerType == allowed {
			return true
		}
	}
	return false
}

func (channel Channel) EffectiveProviderType() (string, bool) {
	providerType := strings.TrimSpace(channel.ProviderType)
	if providerType == "" {
		providerType = constant.GetDefaultProviderType(channel.Type)
	}
	if !constant.IsValidProviderType(providerType) {
		common.SysError("AUD-021 route guard rejected channel with invalid provider_type")
		return "", false
	}
	return providerType, true
}

func IsChannelAllowedByProviderPolicy(channel Channel, policy ProviderTypePolicy) bool {
	return policy.Allows(channel)
}

func normalizeProviderTypes(providerTypes []string) []string {
	if len(providerTypes) == 0 {
		return nil
	}
	normalized := make([]string, 0, len(providerTypes))
	seen := make(map[string]struct{}, len(providerTypes))
	for _, providerType := range providerTypes {
		providerType = strings.TrimSpace(providerType)
		if providerType == "" {
			continue
		}
		if _, ok := seen[providerType]; ok {
			continue
		}
		seen[providerType] = struct{}{}
		normalized = append(normalized, providerType)
	}
	return normalized
}
