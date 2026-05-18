package model

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/gin-gonic/gin"
)

func TestSanitizeLogOtherRedactsAPIKey(t *testing.T) {
	other := SanitizeLogOther(map[string]interface{}{
		"api_key":       "AUD019_FAKE_API_KEY",
		"model":         "aud019-model",
		"provider_type": "official_cloud",
		"channel_id":    12,
		"tokens":        42,
		"cost":          3.14,
	})
	encoded := mustMarshalLogValue(t, other)
	assertNotContainsAny(t, encoded, "AUD019_FAKE_API_KEY")
	assertContainsAll(t, encoded, "aud019-model", "official_cloud", `"channel_id":12`, `"tokens":42`, `"cost":3.14`)
}

func TestSanitizeLogOtherRedactsBearerAuthorization(t *testing.T) {
	other := SanitizeLogOther(map[string]interface{}{
		"headers": map[string]interface{}{
			"Authorization": "Bearer AUD019_FAKE_BEARER_TOKEN",
		},
		"note": "Authorization: Bearer AUD019_INLINE_BEARER",
	})
	encoded := mustMarshalLogValue(t, other)
	assertNotContainsAny(t, encoded, "AUD019_FAKE_BEARER_TOKEN", "AUD019_INLINE_BEARER", "Bearer AUD019")
}

func TestSanitizeLogOtherRedactsPromptMessagesAndInput(t *testing.T) {
	other := SanitizeLogOther(map[string]interface{}{
		"prompt":   "AUD019_FAKE_PROMPT_BODY",
		"messages": []interface{}{"AUD019_FAKE_MESSAGE_BODY"},
		"input":    "AUD019_FAKE_INPUT_BODY",
	})
	encoded := mustMarshalLogValue(t, other)
	assertNotContainsAny(t, encoded, "AUD019_FAKE_PROMPT_BODY", "AUD019_FAKE_MESSAGE_BODY", "AUD019_FAKE_INPUT_BODY")
}

func TestSanitizeLogOtherRedactsResponseAndOutput(t *testing.T) {
	other := SanitizeLogOther(map[string]interface{}{
		"response": "AUD019_FAKE_RESPONSE_BODY",
		"output":   "AUD019_FAKE_OUTPUT_BODY",
	})
	encoded := mustMarshalLogValue(t, other)
	assertNotContainsAny(t, encoded, "AUD019_FAKE_RESPONSE_BODY", "AUD019_FAKE_OUTPUT_BODY")
}

func TestSanitizeLogOtherRedactsNestedSensitiveFields(t *testing.T) {
	other := SanitizeLogOther(map[string]interface{}{
		"admin_info": map[string]interface{}{
			"nested": map[string]interface{}{
				"credential": "AUD019_FAKE_CREDENTIAL",
				"tool_calls": []interface{}{
					map[string]interface{}{"arguments": "AUD019_FAKE_TOOL_ARGUMENTS"},
				},
			},
		},
	})
	encoded := mustMarshalLogValue(t, other)
	assertNotContainsAny(t, encoded, "AUD019_FAKE_CREDENTIAL", "AUD019_FAKE_TOOL_ARGUMENTS")
}

func TestSanitizeLogStringRedactsJSONString(t *testing.T) {
	input := `{"Authorization":"Bearer AUD019_JSON_BEARER","messages":[{"content":"AUD019_JSON_MESSAGE"}],"model":"aud019-model"}`
	got := SanitizeLogString(input)
	assertNotContainsAny(t, got, "AUD019_JSON_BEARER", "AUD019_JSON_MESSAGE", "Bearer AUD019")
	assertContainsAll(t, got, "aud019-model")
}

func TestSanitizeErrorMessageRedactsSecrets(t *testing.T) {
	msg := SanitizeErrorMessage("upstream failed api_key=AUD019_ERROR_KEY token=AUD019_ERROR_TOKEN credential=AUD019_ERROR_CREDENTIAL Authorization: Bearer AUD019_ERROR_BEARER")
	assertNotContainsAny(t, msg, "AUD019_ERROR_KEY", "AUD019_ERROR_TOKEN", "AUD019_ERROR_CREDENTIAL", "AUD019_ERROR_BEARER", "Bearer AUD019")
}

func TestRecordConsumeLogSanitizesOtherBeforePersisting(t *testing.T) {
	resetLogSanitizeTables(t)
	oldStoreFullText := common.StoreFullTextEnabled
	common.StoreFullTextEnabled = false
	t.Cleanup(func() { common.StoreFullTextEnabled = oldStoreFullText })

	ctx := newLogSanitizeContext()
	common.SetContextKey(ctx, constant.ContextKeyChannelProviderType, constant.ProviderTypeOfficialCloud)
	RecordConsumeLog(ctx, 1, RecordConsumeLogParams{
		ChannelId:        99,
		PromptTokens:     7,
		CompletionTokens: 8,
		ModelName:        "aud019-model",
		TokenName:        "aud019-token",
		Quota:            123,
		Content:          "AUD019_FAKE_PROMPT_CONTENT",
		TokenId:          5,
		UseTimeSeconds:   2,
		Group:            "default",
		Other: map[string]interface{}{
			"api_key":       "AUD019_PERSISTED_API_KEY",
			"messages":      []interface{}{"AUD019_PERSISTED_MESSAGE"},
			"response":      "AUD019_PERSISTED_RESPONSE",
			"model":         "aud019-model",
			"provider_type": "official_cloud",
			"channel_id":    99,
			"tokens":        15,
			"cost":          1.23,
		},
	})

	var log Log
	if err := LOG_DB.Order("id desc").First(&log).Error; err != nil {
		t.Fatalf("load log: %v", err)
	}
	if log.Content != "" {
		t.Fatalf("content stored while StoreFullTextEnabled=false: %q", log.Content)
	}
	assertNotContainsAny(t, log.Other, "AUD019_PERSISTED_API_KEY", "AUD019_PERSISTED_MESSAGE", "AUD019_PERSISTED_RESPONSE", "AUD019_FAKE_PROMPT_CONTENT")
	assertContainsAll(t, log.Other, "aud019-model", "official_cloud", `"channel_id":99`, `"tokens":15`, `"cost":1.23`)
}

func TestRecordErrorLogSanitizesContentAndOther(t *testing.T) {
	resetLogSanitizeTables(t)
	ctx := newLogSanitizeContext()
	RecordErrorLog(ctx, 1, 2, "aud019-model", "aud019-token",
		"failed with api_key=AUD019_ERROR_LOG_KEY credential=AUD019_ERROR_LOG_CREDENTIAL",
		3, 4, false, "default", map[string]interface{}{
			"Authorization": "Bearer AUD019_ERROR_LOG_BEARER",
			"input":         "AUD019_ERROR_LOG_INPUT",
		})

	var log Log
	if err := LOG_DB.Order("id desc").First(&log).Error; err != nil {
		t.Fatalf("load log: %v", err)
	}
	assertNotContainsAny(t, log.Content, "AUD019_ERROR_LOG_KEY", "AUD019_ERROR_LOG_CREDENTIAL")
	assertNotContainsAny(t, log.Other, "AUD019_ERROR_LOG_BEARER", "AUD019_ERROR_LOG_INPUT", "Bearer AUD019")
}

func TestDebugPayloadLogsDefaultClosed(t *testing.T) {
	if common.StoreFullTextEnabled {
		t.Fatal("payload/full-text logging must default to disabled")
	}
}

func newLogSanitizeContext() *gin.Context {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	ctx.Set("username", "aud019-user")
	return ctx
}

func resetLogSanitizeTables(t *testing.T) {
	t.Helper()
	LOG_DB.Exec("DELETE FROM logs")
	t.Cleanup(func() {
		LOG_DB.Exec("DELETE FROM logs")
	})
}

func mustMarshalLogValue(t *testing.T, value any) string {
	t.Helper()
	bytes, err := common.Marshal(value)
	if err != nil {
		t.Fatalf("marshal value: %v", err)
	}
	return string(bytes)
}

func assertNotContainsAny(t *testing.T, haystack string, needles ...string) {
	t.Helper()
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			t.Fatalf("expected %q not to contain %q", haystack, needle)
		}
	}
}

func assertContainsAll(t *testing.T, haystack string, needles ...string) {
	t.Helper()
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			t.Fatalf("expected %q to contain %q", haystack, needle)
		}
	}
}
