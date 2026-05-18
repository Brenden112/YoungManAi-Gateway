package kiro

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/channel"
)

// validates M6-F01 — Adaptor struct satisfies channel.Adaptor interface (compile-time)
// If this file compiles, the interface is satisfied.
var _ channel.Adaptor = (*Adaptor)(nil)

// validates M6-F02 — GetChannelName returns the correct identifier
func TestKiroAdaptorGetChannelName(t *testing.T) {
	a := &Adaptor{}
	if got := a.GetChannelName(); got != "KiroGateway" {
		t.Errorf("GetChannelName() = %q, want 'KiroGateway'", got)
	}
}

// validates M6-F02 — GetModelList returns empty (no hardcoded models)
func TestKiroAdaptorGetModelList(t *testing.T) {
	a := &Adaptor{}
	if models := a.GetModelList(); len(models) != 0 {
		t.Errorf("GetModelList() = %v, want empty slice", models)
	}
}

// validates M6-F03 — all relay methods return ErrNotImplemented
// This ensures a KiroGateway channel cannot forward traffic even if somehow enabled.
func TestKiroAdaptorRejectsAllRelayMethods(t *testing.T) {
	a := &Adaptor{}

	if _, err := a.GetRequestURL(nil); err != ErrNotImplemented {
		t.Errorf("GetRequestURL: want ErrNotImplemented, got %v", err)
	}
	if err := a.SetupRequestHeader(nil, nil, nil); err != ErrNotImplemented {
		t.Errorf("SetupRequestHeader: want ErrNotImplemented, got %v", err)
	}
	if _, err := a.ConvertOpenAIRequest(nil, nil, nil); err != ErrNotImplemented {
		t.Errorf("ConvertOpenAIRequest: want ErrNotImplemented, got %v", err)
	}
	if _, err := a.ConvertRerankRequest(nil, 0, dto.RerankRequest{}); err != ErrNotImplemented {
		t.Errorf("ConvertRerankRequest: want ErrNotImplemented, got %v", err)
	}
	if _, err := a.ConvertEmbeddingRequest(nil, nil, dto.EmbeddingRequest{}); err != ErrNotImplemented {
		t.Errorf("ConvertEmbeddingRequest: want ErrNotImplemented, got %v", err)
	}
	if _, err := a.ConvertAudioRequest(nil, nil, dto.AudioRequest{}); err != ErrNotImplemented {
		t.Errorf("ConvertAudioRequest: want ErrNotImplemented, got %v", err)
	}
	if _, err := a.ConvertImageRequest(nil, nil, dto.ImageRequest{}); err != ErrNotImplemented {
		t.Errorf("ConvertImageRequest: want ErrNotImplemented, got %v", err)
	}
	if _, err := a.ConvertOpenAIResponsesRequest(nil, nil, dto.OpenAIResponsesRequest{}); err != ErrNotImplemented {
		t.Errorf("ConvertOpenAIResponsesRequest: want ErrNotImplemented, got %v", err)
	}
	if _, err := a.DoRequest(nil, nil, nil); err != ErrNotImplemented {
		t.Errorf("DoRequest: want ErrNotImplemented, got %v", err)
	}
	if _, apiErr := a.DoResponse(nil, nil, nil); apiErr == nil {
		t.Error("DoResponse: want non-nil error, got nil")
	}
	if _, err := a.ConvertClaudeRequest(nil, nil, nil); err != ErrNotImplemented {
		t.Errorf("ConvertClaudeRequest: want ErrNotImplemented, got %v", err)
	}
	if _, err := a.ConvertGeminiRequest(nil, nil, nil); err != ErrNotImplemented {
		t.Errorf("ConvertGeminiRequest: want ErrNotImplemented, got %v", err)
	}
}
