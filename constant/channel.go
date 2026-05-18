package constant

const (
	ChannelTypeUnknown        = 0
	ChannelTypeOpenAI         = 1
	ChannelTypeMidjourney     = 2
	ChannelTypeAzure          = 3
	ChannelTypeOllama         = 4
	ChannelTypeMidjourneyPlus = 5
	ChannelTypeOpenAIMax      = 6
	ChannelTypeOhMyGPT        = 7
	ChannelTypeCustom         = 8
	ChannelTypeAILS           = 9
	ChannelTypeAIProxy        = 10
	ChannelTypePaLM           = 11
	ChannelTypeAPI2GPT        = 12
	ChannelTypeAIGC2D         = 13
	ChannelTypeAnthropic      = 14
	ChannelTypeBaidu          = 15
	ChannelTypeZhipu          = 16
	ChannelTypeAli            = 17
	ChannelTypeXunfei         = 18
	ChannelType360            = 19
	ChannelTypeOpenRouter     = 20
	ChannelTypeAIProxyLibrary = 21
	ChannelTypeFastGPT        = 22
	ChannelTypeTencent        = 23
	ChannelTypeGemini         = 24
	ChannelTypeMoonshot       = 25
	ChannelTypeZhipu_v4       = 26
	ChannelTypePerplexity     = 27
	ChannelTypeLingYiWanWu    = 31
	ChannelTypeAws            = 33
	ChannelTypeCohere         = 34
	ChannelTypeMiniMax        = 35
	ChannelTypeSunoAPI        = 36
	ChannelTypeDify           = 37
	ChannelTypeJina           = 38
	ChannelCloudflare         = 39
	ChannelTypeSiliconFlow    = 40
	ChannelTypeVertexAi       = 41
	ChannelTypeMistral        = 42
	ChannelTypeDeepSeek       = 43
	ChannelTypeMokaAI         = 44
	ChannelTypeVolcEngine     = 45
	ChannelTypeBaiduV2        = 46
	ChannelTypeXinference     = 47
	ChannelTypeXai            = 48
	ChannelTypeCoze           = 49
	ChannelTypeKling          = 50
	ChannelTypeJimeng         = 51
	ChannelTypeVidu           = 52
	ChannelTypeSubmodel       = 53
	ChannelTypeDoubaoVideo    = 54
	ChannelTypeSora           = 55
	ChannelTypeReplicate      = 56
	ChannelTypeCodex          = 57
	ChannelTypeKiroGateway    = 58
	ChannelTypeDummy          // this one is only for count, do not add any channel after this

)

var ChannelBaseURLs = []string{
	"",                                    // 0
	"https://api.openai.com",              // 1
	"https://oa.api2d.net",                // 2
	"",                                    // 3
	"http://localhost:11434",              // 4
	"https://api.openai-sb.com",           // 5
	"https://api.openaimax.com",           // 6
	"https://api.ohmygpt.com",             // 7
	"",                                    // 8
	"https://api.caipacity.com",           // 9
	"https://api.aiproxy.io",              // 10
	"",                                    // 11
	"https://api.api2gpt.com",             // 12
	"https://api.aigc2d.com",              // 13
	"https://api.anthropic.com",           // 14
	"https://aip.baidubce.com",            // 15
	"https://open.bigmodel.cn",            // 16
	"https://dashscope.aliyuncs.com",      // 17
	"",                                    // 18
	"https://api.360.cn",                  // 19
	"https://openrouter.ai/api",           // 20
	"https://api.aiproxy.io",              // 21
	"https://fastgpt.run/api/openapi",     // 22
	"https://hunyuan.tencentcloudapi.com", //23
	"https://generativelanguage.googleapis.com", //24
	"https://api.moonshot.cn",                   //25
	"https://open.bigmodel.cn",                  //26
	"https://api.perplexity.ai",                 //27
	"",                                          //28
	"",                                          //29
	"",                                          //30
	"https://api.lingyiwanwu.com",               //31
	"",                                          //32
	"",                                          //33
	"https://api.cohere.ai",                     //34
	"https://api.minimax.chat",                  //35
	"",                                          //36
	"https://api.dify.ai",                       //37
	"https://api.jina.ai",                       //38
	"https://api.cloudflare.com",                //39
	"https://api.siliconflow.cn",                //40
	"",                                          //41
	"https://api.mistral.ai",                    //42
	"https://api.deepseek.com",                  //43
	"https://api.moka.ai",                       //44
	"https://ark.cn-beijing.volces.com",         //45
	"https://qianfan.baidubce.com",              //46
	"",                                          //47
	"https://api.x.ai",                          //48
	"https://api.coze.cn",                       //49
	"https://api.klingai.com",                   //50
	"https://visual.volcengineapi.com",          //51
	"https://api.vidu.cn",                       //52
	"https://llm.submodel.ai",                   //53
	"https://ark.cn-beijing.volces.com",         //54
	"https://api.openai.com",                    //55
	"https://api.replicate.com",                 //56
	"https://chatgpt.com",                       //57
	"",                                          //58 KiroGateway — base URL configured per-channel
}

var ChannelTypeNames = map[int]string{
	ChannelTypeUnknown:        "Unknown",
	ChannelTypeOpenAI:         "OpenAI",
	ChannelTypeMidjourney:     "Midjourney",
	ChannelTypeAzure:          "Azure",
	ChannelTypeOllama:         "Ollama",
	ChannelTypeMidjourneyPlus: "MidjourneyPlus",
	ChannelTypeOpenAIMax:      "OpenAIMax",
	ChannelTypeOhMyGPT:        "OhMyGPT",
	ChannelTypeCustom:         "Custom",
	ChannelTypeAILS:           "AILS",
	ChannelTypeAIProxy:        "AIProxy",
	ChannelTypePaLM:           "PaLM",
	ChannelTypeAPI2GPT:        "API2GPT",
	ChannelTypeAIGC2D:         "AIGC2D",
	ChannelTypeAnthropic:      "Anthropic",
	ChannelTypeBaidu:          "Baidu",
	ChannelTypeZhipu:          "Zhipu",
	ChannelTypeAli:            "Ali",
	ChannelTypeXunfei:         "Xunfei",
	ChannelType360:            "360",
	ChannelTypeOpenRouter:     "OpenRouter",
	ChannelTypeAIProxyLibrary: "AIProxyLibrary",
	ChannelTypeFastGPT:        "FastGPT",
	ChannelTypeTencent:        "Tencent",
	ChannelTypeGemini:         "Gemini",
	ChannelTypeMoonshot:       "Moonshot",
	ChannelTypeZhipu_v4:       "ZhipuV4",
	ChannelTypePerplexity:     "Perplexity",
	ChannelTypeLingYiWanWu:    "LingYiWanWu",
	ChannelTypeAws:            "AWS",
	ChannelTypeCohere:         "Cohere",
	ChannelTypeMiniMax:        "MiniMax",
	ChannelTypeSunoAPI:        "SunoAPI",
	ChannelTypeDify:           "Dify",
	ChannelTypeJina:           "Jina",
	ChannelCloudflare:         "Cloudflare",
	ChannelTypeSiliconFlow:    "SiliconFlow",
	ChannelTypeVertexAi:       "VertexAI",
	ChannelTypeMistral:        "Mistral",
	ChannelTypeDeepSeek:       "DeepSeek",
	ChannelTypeMokaAI:         "MokaAI",
	ChannelTypeVolcEngine:     "VolcEngine",
	ChannelTypeBaiduV2:        "BaiduV2",
	ChannelTypeXinference:     "Xinference",
	ChannelTypeXai:            "xAI",
	ChannelTypeCoze:           "Coze",
	ChannelTypeKling:          "Kling",
	ChannelTypeJimeng:         "Jimeng",
	ChannelTypeVidu:           "Vidu",
	ChannelTypeSubmodel:       "Submodel",
	ChannelTypeDoubaoVideo:    "DoubaoVideo",
	ChannelTypeSora:           "Sora",
	ChannelTypeReplicate:      "Replicate",
	ChannelTypeCodex:          "Codex",
	ChannelTypeKiroGateway:    "KiroGateway",
}

func GetChannelTypeName(channelType int) string {
	if name, ok := ChannelTypeNames[channelType]; ok {
		return name
	}
	return "Unknown"
}

type ChannelSpecialBase struct {
	ClaudeBaseURL string
	OpenAIBaseURL string
}

var ChannelSpecialBases = map[string]ChannelSpecialBase{
	"glm-coding-plan": {
		ClaudeBaseURL: "https://open.bigmodel.cn/api/anthropic",
		OpenAIBaseURL: "https://open.bigmodel.cn/api/coding/paas/v4",
	},
	"glm-coding-plan-international": {
		ClaudeBaseURL: "https://api.z.ai/api/anthropic",
		OpenAIBaseURL: "https://api.z.ai/api/coding/paas/v4",
	},
	"kimi-coding-plan": {
		ClaudeBaseURL: "https://api.kimi.com/coding",
		OpenAIBaseURL: "https://api.kimi.com/coding/v1",
	},
	"doubao-coding-plan": {
		ClaudeBaseURL: "https://ark.cn-beijing.volces.com/api/coding",
		OpenAIBaseURL: "https://ark.cn-beijing.volces.com/api/coding/v3",
	},
}

// Provider type constants for the B2B gateway layer.
// These classify upstream channels by trust level and availability scope.
const (
	ProviderTypeOfficialCloud     = "official_cloud"
	ProviderTypeAggregator        = "aggregator"
	ProviderTypeAuthorizedProxy   = "authorized_proxy"
	ProviderTypeExperimentalProxy = "experimental_proxy"
)

var validProviderTypes = map[string]bool{
	ProviderTypeOfficialCloud:     true,
	ProviderTypeAggregator:        true,
	ProviderTypeAuthorizedProxy:   true,
	ProviderTypeExperimentalProxy: true,
}

// IsValidProviderType returns true if pt is one of the four recognised provider types.
func IsValidProviderType(pt string) bool {
	return validProviderTypes[pt]
}

// ChannelTypeDefaultProviderType maps each channel type to its default provider_type.
// Types absent from this map default to official_cloud.
var ChannelTypeDefaultProviderType = map[int]string{
	// ── official_cloud — direct upstream API ─────────────────────────────────
	ChannelTypeOpenAI:      ProviderTypeOfficialCloud,
	ChannelTypeAzure:       ProviderTypeOfficialCloud,
	ChannelTypeAnthropic:   ProviderTypeOfficialCloud,
	ChannelTypeGemini:      ProviderTypeOfficialCloud,
	ChannelTypeAws:         ProviderTypeOfficialCloud,
	ChannelTypeVertexAi:    ProviderTypeOfficialCloud,
	ChannelTypeBaidu:       ProviderTypeOfficialCloud,
	ChannelTypeBaiduV2:     ProviderTypeOfficialCloud,
	ChannelTypeZhipu:       ProviderTypeOfficialCloud,
	ChannelTypeZhipu_v4:    ProviderTypeOfficialCloud,
	ChannelTypeAli:         ProviderTypeOfficialCloud,
	ChannelTypeXunfei:      ProviderTypeOfficialCloud,
	ChannelType360:         ProviderTypeOfficialCloud,
	ChannelTypeTencent:     ProviderTypeOfficialCloud,
	ChannelTypeMoonshot:    ProviderTypeOfficialCloud,
	ChannelTypePerplexity:  ProviderTypeOfficialCloud,
	ChannelTypeLingYiWanWu: ProviderTypeOfficialCloud,
	ChannelTypeCohere:      ProviderTypeOfficialCloud,
	ChannelTypeMiniMax:     ProviderTypeOfficialCloud,
	ChannelTypeMistral:     ProviderTypeOfficialCloud,
	ChannelTypeDeepSeek:    ProviderTypeOfficialCloud,
	ChannelTypeVolcEngine:  ProviderTypeOfficialCloud,
	ChannelTypeXai:         ProviderTypeOfficialCloud,
	ChannelTypeReplicate:   ProviderTypeOfficialCloud,
	ChannelTypeKling:       ProviderTypeOfficialCloud,
	ChannelTypeJimeng:      ProviderTypeOfficialCloud,
	ChannelTypeVidu:        ProviderTypeOfficialCloud,
	ChannelTypeSora:        ProviderTypeOfficialCloud,
	ChannelTypeDoubaoVideo: ProviderTypeOfficialCloud,
	ChannelTypePaLM:        ProviderTypeOfficialCloud,
	ChannelTypeMokaAI:      ProviderTypeOfficialCloud,

	// ── aggregator — multi-provider aggregator ────────────────────────────────
	ChannelTypeOpenRouter:     ProviderTypeAggregator,
	ChannelTypeSiliconFlow:    ProviderTypeAggregator,
	ChannelTypeSubmodel:       ProviderTypeAggregator,
	ChannelTypeOpenAIMax:      ProviderTypeAggregator,
	ChannelTypeOhMyGPT:        ProviderTypeAggregator,
	ChannelTypeAILS:           ProviderTypeAggregator,
	ChannelTypeAIProxy:        ProviderTypeAggregator,
	ChannelTypeAPI2GPT:        ProviderTypeAggregator,
	ChannelTypeAIGC2D:         ProviderTypeAggregator,
	ChannelTypeAIProxyLibrary: ProviderTypeAggregator,
	ChannelTypeFastGPT:        ProviderTypeAggregator,
	ChannelTypeCoze:           ProviderTypeAggregator,

	// ── authorized_proxy — self-hosted / local ────────────────────────────────
	ChannelTypeOllama:     ProviderTypeAuthorizedProxy,
	ChannelTypeCustom:     ProviderTypeAuthorizedProxy,
	ChannelTypeXinference: ProviderTypeAuthorizedProxy,
	ChannelTypeDify:       ProviderTypeAuthorizedProxy,
	ChannelCloudflare:     ProviderTypeAuthorizedProxy,
	ChannelTypeJina:       ProviderTypeAuthorizedProxy,
	ChannelTypeSunoAPI:    ProviderTypeAuthorizedProxy,

	// ── experimental_proxy — unofficial / experimental ────────────────────────
	ChannelTypeCodex:       ProviderTypeExperimentalProxy,
	ChannelTypeKiroGateway: ProviderTypeExperimentalProxy,
}

// GetDefaultProviderType returns the default provider_type for a given channel type.
// Falls back to official_cloud for unknown or unlisted types.
func GetDefaultProviderType(channelType int) string {
	if pt, ok := ChannelTypeDefaultProviderType[channelType]; ok {
		return pt
	}
	return ProviderTypeOfficialCloud
}

// Risk level constants for Channel.RiskLevel.
const (
	RiskLevelNormal = "normal"
	RiskLevelHigh   = "high"
)

// Scope constants for Channel.AvailableScope.
const (
	ScopePublic       = "public"
	ScopeInternalOnly = "internal_only"
)

// Visibility constants for Channel.Visibility.
const (
	VisibilityPublic       = "public"
	VisibilityInternalOnly = "internal_only"
)
