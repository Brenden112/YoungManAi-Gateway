// Package kiro provides the KiroGateway channel adapter skeleton.
// KiroGateway is classified as experimental_proxy and is disabled by default.
// All relay methods return ErrNotImplemented until a real implementation is wired in.
package kiro

import (
	"fmt"
	"io"
	"net/http"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

// ErrNotImplemented is returned by all relay methods until the adapter is fully implemented.
var ErrNotImplemented = fmt.Errorf("KiroGateway: adapter not yet implemented")

// Adaptor is the KiroGateway channel adapter skeleton.
// It satisfies the channel.Adaptor interface but rejects all relay calls,
// ensuring that even if a KiroGateway channel is somehow enabled it cannot
// forward traffic until the implementation is complete.
type Adaptor struct{}

func (a *Adaptor) Init(_ *relaycommon.RelayInfo) {}

func (a *Adaptor) GetRequestURL(_ *relaycommon.RelayInfo) (string, error) {
	return "", ErrNotImplemented
}

func (a *Adaptor) SetupRequestHeader(_ *gin.Context, _ *http.Header, _ *relaycommon.RelayInfo) error {
	return ErrNotImplemented
}

func (a *Adaptor) ConvertOpenAIRequest(_ *gin.Context, _ *relaycommon.RelayInfo, _ *dto.GeneralOpenAIRequest) (any, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) ConvertRerankRequest(_ *gin.Context, _ int, _ dto.RerankRequest) (any, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) ConvertEmbeddingRequest(_ *gin.Context, _ *relaycommon.RelayInfo, _ dto.EmbeddingRequest) (any, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) ConvertAudioRequest(_ *gin.Context, _ *relaycommon.RelayInfo, _ dto.AudioRequest) (io.Reader, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) ConvertImageRequest(_ *gin.Context, _ *relaycommon.RelayInfo, _ dto.ImageRequest) (any, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) ConvertOpenAIResponsesRequest(_ *gin.Context, _ *relaycommon.RelayInfo, _ dto.OpenAIResponsesRequest) (any, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) DoRequest(_ *gin.Context, _ *relaycommon.RelayInfo, _ io.Reader) (any, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) DoResponse(_ *gin.Context, _ *http.Response, _ *relaycommon.RelayInfo) (any, *types.NewAPIError) {
	return nil, types.NewError(ErrNotImplemented, types.ErrorCodeBadResponseStatusCode)
}

func (a *Adaptor) ConvertClaudeRequest(_ *gin.Context, _ *relaycommon.RelayInfo, _ *dto.ClaudeRequest) (any, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) ConvertGeminiRequest(_ *gin.Context, _ *relaycommon.RelayInfo, _ *dto.GeminiChatRequest) (any, error) {
	return nil, ErrNotImplemented
}

func (a *Adaptor) GetModelList() []string { return nil }

func (a *Adaptor) GetChannelName() string { return "KiroGateway" }
