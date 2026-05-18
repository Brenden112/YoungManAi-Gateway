package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/QuantumNous/new-api/common"
)

const (
	logRedactedValue       = "[REDACTED]"
	logMaxStringValueBytes = 1024
)

var (
	logBearerPattern    = regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._~+/=-]+`)
	logKVSecretPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(api[_-]?key\s*[:=]\s*)[^\s'",}&]+`),
		regexp.MustCompile(`(?i)(authorization\s*[:=]\s*)[^\n\r,}]+`),
		regexp.MustCompile(`(?i)(token\s*[:=]\s*)[^\s'",}&]+`),
		regexp.MustCompile(`(?i)(access[_-]?token\s*[:=]\s*)[^\s'",}&]+`),
		regexp.MustCompile(`(?i)(refresh[_-]?token\s*[:=]\s*)[^\s'",}&]+`),
		regexp.MustCompile(`(?i)(id[_-]?token\s*[:=]\s*)[^\s'",}&]+`),
		regexp.MustCompile(`(?i)(credential\s*[:=]\s*)[^\s'",}&]+`),
		regexp.MustCompile(`(?i)(secret\s*[:=]\s*)[^\s'",}&]+`),
		regexp.MustCompile(`(?i)(password\s*[:=]\s*)[^\s'",}&]+`),
		regexp.MustCompile(`(?i)(prompt\s*[:=]\s*)[^\n\r,}]+`),
		regexp.MustCompile(`(?i)(messages?\s*[:=]\s*)[^\n\r,}]+`),
		regexp.MustCompile(`(?i)(input\s*[:=]\s*)[^\n\r,}]+`),
		regexp.MustCompile(`(?i)(responses?\s*[:=]\s*)[^\n\r,}]+`),
		regexp.MustCompile(`(?i)(outputs?\s*[:=]\s*)[^\n\r,}]+`),
		regexp.MustCompile(`(?i)(tools?\s*[:=]\s*)[^\n\r,}]+`),
		regexp.MustCompile(`(?i)(tool_calls?\s*[:=]\s*)[^\n\r,}]+`),
		regexp.MustCompile(`(?i)(headers?\s*[:=]\s*)[^\n\r,}]+`),
	}
)

func SanitizeLogOther(other map[string]interface{}) map[string]interface{} {
	if other == nil {
		return nil
	}
	sanitized, ok := sanitizeLogValue(other, "").(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return sanitized
}

func SanitizeLogPayload(payload any) any {
	return sanitizeLogValue(payload, "")
}

func SanitizeLogString(value string) string {
	var decoded any
	if err := common.Unmarshal([]byte(value), &decoded); err == nil {
		encoded, marshalErr := common.Marshal(sanitizeLogValue(decoded, ""))
		if marshalErr == nil {
			return string(encoded)
		}
	}
	return truncateLogString(redactSensitiveText(value))
}

func SanitizeErrorMessage(message string) string {
	return SanitizeLogString(common.MaskSensitiveInfo(message))
}

func sanitizeLogValue(value any, key string) any {
	if isSensitiveLogKey(key) {
		return logRedactedValue
	}
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]interface{}:
		out := make(map[string]interface{}, len(typed))
		for k, v := range typed {
			out[k] = sanitizeLogValue(v, k)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(typed))
		for i, v := range typed {
			out[i] = sanitizeLogValue(v, key)
		}
		return out
	case string:
		return sanitizeLogStringValue(typed)
	case fmt.Stringer:
		return sanitizeLogStringValue(typed.String())
	default:
		return typed
	}
}

func sanitizeLogStringValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) > 1 && ((strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) || (strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]"))) {
		var decoded any
		if err := common.Unmarshal([]byte(trimmed), &decoded); err == nil {
			encoded, marshalErr := common.Marshal(sanitizeLogValue(decoded, ""))
			if marshalErr == nil {
				return string(encoded)
			}
		}
	}
	return truncateLogString(redactSensitiveText(value))
}

func isSensitiveLogKey(key string) bool {
	normalized := normalizeLogKey(key)
	if strings.Contains(normalized, "credential") ||
		strings.Contains(normalized, "secret") ||
		strings.Contains(normalized, "authorization") ||
		strings.Contains(normalized, "api_key") ||
		strings.Contains(normalized, "apikey") {
		return true
	}
	switch normalized {
	case "api_key", "apikey", "api-key", "x-api-key", "key",
		"authorization", "auth", "bearer",
		"token", "access_token", "refresh_token", "id_token", "auth_token",
		"secret", "client_secret", "credential", "credentials", "password",
		"prompt", "prompts", "messages", "message", "input", "inputs",
		"response", "responses", "output", "outputs", "completion", "completions",
		"tool", "tools", "tool_calls", "tool_call", "arguments", "parameters",
		"headers", "header", "request", "request_body", "response_body", "body", "payload", "raw", "metadata":
		return true
	default:
		return false
	}
}

func normalizeLogKey(key string) string {
	key = strings.TrimSpace(strings.ToLower(key))
	key = strings.ReplaceAll(key, "-", "_")
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.ReplaceAll(key, ".", "_")
	return key
}

func redactSensitiveText(value string) string {
	value = logBearerPattern.ReplaceAllString(value, "Bearer "+logRedactedValue)
	for _, pattern := range logKVSecretPatterns {
		value = pattern.ReplaceAllString(value, "${1}"+logRedactedValue)
	}
	return value
}

func truncateLogString(value string) string {
	if len(value) <= logMaxStringValueBytes {
		return value
	}
	return value[:logMaxStringValueBytes] + "...[TRUNCATED]"
}

func logOtherToJSON(other map[string]interface{}) string {
	sanitized := SanitizeLogOther(other)
	if sanitized == nil {
		return ""
	}
	bytes, err := common.Marshal(sanitized)
	if err != nil {
		return ""
	}
	return string(bytes)
}
