# B2B Multi-Provider API Gateway MVP Plan

## 1. 背景与主干选择

主干建议基于 `new-api` 二次开发。

原因：

- `new-api` 本身就是 AI API gateway/proxy，已有用户管理、计费、限流、管理后台、多上游聚合等基础能力。
- 项目能力与合法授权 API 网关、组织级鉴权、多模型管理、用量分析、成本核算等 B 端场景高度匹配。
- 第一阶段可以先复用 OpenAI-compatible API 与现有渠道能力，再逐步扩展 Provider Adapter、企业组织、支付、代理和多协议入口。

第一阶段只暴露：

- OpenAI-compatible API

后续逐步增加：

- Claude-compatible Messages API
- Gemini-compatible API
- OpenAI Responses API

## 2. 总体架构

系统按 7 层设计：

```text
客户端层
  ↓
统一协议入口层
  ↓
鉴权与企业权限层
  ↓
计费与额度层
  ↓
模型路由层
  ↓
Provider 适配层
  ↓
上游资源层
```

展开链路：

```text
国内 B 端客户 / 代理客户 / 自用用户
        ↓
API Gateway
        ↓
OpenAI-compatible / Claude-compatible / Gemini-compatible / Responses API
        ↓
用户鉴权 / 组织 / 项目 / API Key / IP 白名单
        ↓
余额 / 额度 / 模型倍率 / 企业合同价 / 代理归属
        ↓
路由策略
        ↓
Provider Adapter
        ↓
official_cloud / aggregator / authorized_proxy / experimental_proxy
        ↓
海外云厂商 / 海外聚合平台 / 内部备用通道
```

## 3. Provider Adapter Framework

核心原则：

- 不写死某个备用服务。
- 设计为 `provider adapter framework`。
- `kiro-gateway` 只作为其中一个 `experimental_proxy` adapter 示例。
- 所有 Provider 必须经过统一鉴权、计费、日志、风控、项目权限和路由策略。

### 3.1 Provider 类型

| 类型 | 说明 | 是否开放给 B 端 |
| --- | --- | --- |
| `official_cloud` | Azure、AWS Bedrock、Google Vertex、海外云厂商等 | 可以 |
| `aggregator` | OpenRouter、海外模型聚合平台等 | 可以，需评估 |
| `authorized_proxy` | 明确有授权、可转发的内部代理通道 | 谨慎开放 |
| `experimental_proxy` | 非标准备用通道，例如 kiro-gateway 或其他内部备用通道 | 默认仅内部 |

`experimental_proxy` 默认约束：

```yaml
risk_level: high
visibility: internal_only
manual_enable_required: true
available_scope: internal_only
```

### 3.2 Provider 字段

`providers` 至少包含：

| 字段 | 说明 |
| --- | --- |
| `id` | Provider ID |
| `name` | Provider 名称 |
| `type` | `official_cloud` / `aggregator` / `authorized_proxy` / `experimental_proxy` |
| `risk_level` | 风险等级 |
| `base_url` | 上游地址 |
| `auth_type` | 鉴权方式 |
| `protocol_type` | 协议类型 |
| `support_stream` | 是否支持流式 |
| `support_tools` | 是否支持工具调用 |
| `support_vision` | 是否支持视觉 |
| `support_embedding` | 是否支持 embedding |
| `support_usage_return` | 是否返回 usage |
| `enabled` | 是否启用 |
| `available_scope` | 可用范围 |
| `priority` | 优先级 |
| `weight` | 权重 |
| `timeout_ms` | 超时 |
| `max_retries` | 最大重试 |
| `daily_quota` | 日额度 |
| `monthly_quota` | 月额度 |
| `current_usage` | 当前用量 |
| `fail_count` | 失败次数 |
| `success_count` | 成功次数 |
| `last_error` | 最近错误 |
| `last_health_check_at` | 最近健康检查时间 |
| `created_at` | 创建时间 |
| `updated_at` | 更新时间 |

`available_scope` 枚举：

- `all_users`
- `enterprise_only`
- `selected_groups`
- `internal_only`
- `admin_only`

### 3.3 Provider Account

一个 Provider 下可能有多个上游账号、key、token 或 endpoint，因此需要独立表：

`provider_accounts` 字段建议：

| 字段 | 说明 |
| --- | --- |
| `id` | Account ID |
| `provider_id` | 所属 Provider |
| `account_name` | 账号名称 |
| `credential_encrypted` | 加密后的上游凭证 |
| `credential_type` | 凭证类型 |
| `region` | 区域 |
| `endpoint` | 独立 endpoint |
| `quota_limit` | 配额上限 |
| `quota_used` | 已用配额 |
| `status` | 状态 |
| `risk_level` | 风险等级 |
| `available_scope` | 可用范围 |
| `last_used_at` | 最近使用时间 |
| `last_error` | 最近错误 |
| `cooldown_until` | 冷却截止时间 |
| `created_at` | 创建时间 |
| `updated_at` | 更新时间 |

要求：

- 上游密钥必须加密存储。
- 不允许明文保存完整上游 key。

### 3.4 Channel 与 Provider 关系

定义：

- `Provider`：上游服务类型。
- `Provider Account`：该上游下的具体账号、key 或 endpoint。
- `Channel`：对外可用的一条路由通道。

示例：

```text
Provider: Azure OpenAI
Provider Account: azure-us-east-account-01
Channel: gpt-4.1-official-high-quality

Provider: OpenRouter
Provider Account: openrouter-main
Channel: claude-sonnet-router

Provider: Experimental Proxy
Provider Account: kiro-internal-01
Channel: claude-sonnet-internal-test
```

## 4. 数据库规划

所有数据库设计必须同时兼容 SQLite、MySQL 5.7.8+、PostgreSQL 9.6+。优先使用 GORM 抽象，避免数据库专属 SQL。

### 4.1 用户、组织与项目

核心表：

- `users`
- `organizations`
- `organization_members`
- `roles`
- `permissions`
- `projects`
- `project_members`

最小字段：

```text
users
- id
- email
- phone
- password_hash
- status
- user_type
- balance
- created_at

organizations
- id
- name
- owner_user_id
- plan_type
- contract_type
- status
- created_at

organization_members
- id
- organization_id
- user_id
- role
- status

projects
- id
- organization_id
- name
- status
- monthly_budget
- daily_budget
```

第一阶段 UI 可以不展示完整企业功能，但表结构要预留。

### 4.2 API Key

`api_keys` 字段建议：

```text
- id
- organization_id
- project_id
- user_id
- key_hash
- key_prefix
- name
- status
- allowed_models
- allowed_provider_types
- allowed_ip_list
- daily_quota
- monthly_quota
- expires_at
- last_used_at
- created_at
```

要求：

- 不保存完整 API Key 明文。
- 只保存 `key_hash` 与 `key_prefix`。

### 4.3 Provider 与路由

核心表：

- `providers`
- `provider_accounts`
- `channels`
- `model_mappings`
- `route_policies`
- `route_rules`

`model_mappings`：

```text
- id
- public_model_name
- internal_model_name
- provider_model_name
- protocol_type
- input_price
- output_price
- billing_multiplier
- enabled
```

`route_policies`：

```text
- id
- name
- organization_id
- project_id
- priority_mode
- allow_experimental
- fallback_enabled
```

`route_rules`：

```text
- id
- route_policy_id
- model_name
- channel_id
- priority
- weight
- enabled
```

要求：

- `allow_experimental` 默认必须为 `false`。
- 只有内部测试用户、测试项目或管理员手动配置后，才允许访问 `experimental_proxy`。

### 4.4 计费与日志

核心表：

- `usage_logs`
- `usage_events`
- `usage_daily`
- `billing_ledger`
- `orders`
- `payments`
- `refunds`

`usage_logs` 字段建议：

```text
- request_id
- user_id
- organization_id
- project_id
- api_key_id
- model
- public_model
- provider_id
- provider_type
- channel_id
- provider_account_id
- input_tokens
- output_tokens
- total_tokens
- cached_tokens
- reasoning_tokens
- cost
- charged_amount
- status
- error_code
- error_message_sanitized
- latency_ms
- ip
- user_agent
- created_at
```

默认不保存：

- 完整 prompt
- 完整 response

可选短期调试表：

`debug_payload_logs`

规则：

- 默认关闭。
- 项目级开启。
- 最多保存 1-7 天。
- 自动脱敏。
- 仅管理员和企业 owner 可见。

### 4.5 支付与代理分销

核心表：

- `invite_codes`
- `referral_relations`
- `agent_profiles`
- `commission_rules`
- `commission_records`
- `withdraw_requests`

第一版只做：

- `invite_codes`
- `referral_relations`
- `agent_profiles`

不在第一版做自动提现。

## 5. 路由与备用 Provider 策略

### 5.1 请求路由流程

```text
1. 校验 API Key
2. 解析 organization / project / user
3. 检查余额和额度
4. 检查模型权限
5. 匹配 route_policy
6. 选择候选 channels
7. 过滤不可用 provider
8. 按优先级 / 权重选择 channel
9. 发起上游请求
10. 记录 usage
11. 扣费
12. 返回结果
```

### 5.2 路由策略

第一阶段支持 3 种：

| 策略 | 说明 |
| --- | --- |
| `quality_first` | 优先走稳定高质量上游 |
| `cost_first` | 优先走低成本上游 |
| `balance` | 按权重分配 |

后续可增加：

- `latency_first`
- `domestic_first`
- `official_only`
- `enterprise_dedicated`

### 5.3 备用 Provider 使用规则

默认规则：

- 普通用户：不可用 `experimental_proxy`
- 企业用户：默认不可用 `experimental_proxy`
- 内部用户：可用 `experimental_proxy`
- 测试项目：可用 `experimental_proxy`
- 管理员：可手动指定 `experimental_proxy`

控制字段：

- `allowed_provider_types`
- `allow_experimental`
- `available_scope`

### 5.4 故障切换规则

允许默认自动切换：

- `official_cloud` → `official_cloud`
- `official_cloud` → `aggregator`
- `aggregator` → `aggregator`
- `authorized_proxy` → `authorized_proxy`

不建议默认自动切换：

- `official_cloud` → `experimental_proxy`

例外条件：

- 项目显式开启 `allow_experimental = true`。
- 用户或项目明确知道该链路会使用非标准备用通道。

## 6. 开发路线图

### Phase 0：项目初始化与规则冻结

目标：先把方向、边界、仓库和文档定下来。

任务：

1. Fork `new-api`。
2. 单独保留 `kiro-gateway` 和其他备用 Provider 为 external services。
3. 建立 monorepo 或分仓结构。
4. 编写 `ARCHITECTURE.md`。
5. 编写 `PROVIDER_SPEC.md`。
6. 编写 `BILLING_SPEC.md`。
7. 编写 `ROUTING_SPEC.md`。
8. 编写 `LOGGING_POLICY.md`。
9. 编写 `SECURITY_POLICY.md`。
10. 编写 `COMPLIANCE_NOTES.md`。

建议结构：

```text
youngman-ai-gateway/
├── gateway/
├── provider-services/
│   ├── kiro-gateway/
│   ├── provider-adapter-sdk/
│   └── experimental-providers/
├── admin-extensions/
├── docs/
│   ├── ARCHITECTURE.md
│   ├── PROVIDER_SPEC.md
│   ├── BILLING_SPEC.md
│   ├── ROUTING_SPEC.md
│   ├── LOGGING_POLICY.md
│   └── ROADMAP.md
└── infra/
    ├── docker-compose.yml
    ├── nginx.conf
    ├── redis.conf
    └── backup/
```

验收：

- 能本地启动 `new-api`。
- 能本地启动 `kiro-gateway`。
- 能明确主链路和备用链路。
- 能确定第一批 Provider 类型。
- 能确定第一版只做 OpenAI-compatible API。

### Phase 1：OpenAI-compatible MVP

目标：跑通“用户调用 → 计量 → 扣费 → 日志”的闭环。

任务：

1. 部署 `new-api`。
2. 接入第一个海外正规上游。
3. 接入第二个海外正规上游。
4. 接入 `kiro-gateway` 作为 `experimental_proxy`。
5. 创建用户。
6. 创建 API Key。
7. 设置余额。
8. 设置模型倍率。
9. 调通 `POST /v1/chat/completions`。
10. 调通 `GET /v1/models`。
11. 记录 `usage_logs`。
12. 余额不足自动拒绝。
13. 请求失败记录错误。
14. 管理员手动加余额。

第一阶段模型接口：

```http
GET  /v1/models
POST /v1/chat/completions
POST /v1/embeddings
GET  /api/usage/token
```

验收：

- 用户能创建 API Key。
- 用户能用 OpenAI SDK 调用。
- 系统能识别用户和项目。
- 系统能计算 token。
- 系统能扣费。
- 系统能记录日志。
- 余额不足能拒绝请求。
- 管理员能看到请求日志。
- internal 用户能走 `kiro-gateway`。
- 普通用户不能走 `experimental_proxy`。

### Phase 2：Provider 插件化改造

目标：让备用 Provider 不止 `kiro-gateway`，而是任何 Provider 都能按统一标准接入。

任务：

1. 定义 `ProviderAdapter` interface。
2. 实现 OpenAI-like Adapter。
3. 实现 Claude-like Adapter。
4. 实现 Gemini-like Adapter。
5. 实现 Generic HTTP Proxy Adapter。
6. 实现 Kiro Gateway Adapter。
7. 实现 Provider Health Check。
8. 实现 Provider Account Pool。
9. 实现 Channel 权重路由。
10. 实现 Provider 风险等级。

接口建议：

```go
type ProviderAdapter interface {
    Name() string
    Type() ProviderType
    SupportedProtocols() []ProtocolType
    ListModels(ctx context.Context) ([]ModelInfo, error)
    NormalizeRequest(ctx context.Context, req any) (any, error)
    SendRequest(ctx context.Context, req any) (any, error)
    StreamRequest(ctx context.Context, req any) (<-chan StreamChunk, error)
    NormalizeResponse(ctx context.Context, resp any) (any, error)
    ExtractUsage(ctx context.Context, resp any) (Usage, error)
    HealthCheck(ctx context.Context) error
    ClassifyError(err error) ProviderErrorType
}
```

错误分类：

- `auth_error`
- `rate_limit`
- `quota_exceeded`
- `timeout`
- `upstream_5xx`
- `network_error`
- `model_not_found`
- `content_filter`
- `unknown_error`

验收：

- 新增 Provider 不需要大改主流程。
- 每个 Provider 都能独立健康检查。
- 每个 Provider 都能独立计费。
- 每个 Provider 都能独立关闭。
- 每个 Provider 都能限制可见范围。
- `experimental_proxy` 默认不能被普通用户调用。

### Phase 3：企业组织与项目结构

目标：预留并逐步启用 B 端组织、项目、成员和项目级 Key。

任务：

1. 新增 `organizations` 表。
2. 新增 `organization_members` 表。
3. 新增 `projects` 表。
4. API Key 绑定 project。
5. `usage_logs` 绑定 organization/project。
6. 项目级余额或预算。
7. 项目级模型白名单。
8. 项目级 Provider 类型限制。
9. 项目级 IP 白名单。
10. 企业管理员角色。

第一版角色：

- `owner`
- `admin`
- `developer`
- `billing`
- `viewer`

验收：

- 一个企业可以有多个项目。
- 一个项目可以有多个 API Key。
- 每个项目可以限制模型。
- 每个项目可以限制 Provider 类型。
- 每个项目可以单独统计用量。

### Phase 4：计费系统增强

目标：从简单余额扣费升级为可运营账务系统。

任务：

1. `billing_ledger` 流水表。
2. `user_balance` 和 `organization_balance` 分离。
3. 支持模型输入价 / 输出价。
4. 支持缓存 token 价格。
5. 支持 reasoning token 价格。
6. 支持失败请求不扣费或部分扣费。
7. 支持企业合同价。
8. 支持代理客户价。
9. 支持余额预扣和最终结算。
10. 支持日账单和月账单。

计费流程：

```text
请求开始
  ↓
预估最大费用
  ↓
余额预冻结
  ↓
上游返回 usage
  ↓
计算实际费用
  ↓
释放多余冻结
  ↓
写 billing_ledger
```

`billing_ledger` 类型：

- `recharge`
- `consume`
- `refund`
- `adjustment`
- `commission`
- `contract_settlement`

验收：

- 所有余额变化都有流水。
- 每一条 usage 都能追溯到扣费记录。
- 管理员可以手动调账。
- 用户可以查看余额明细。
- 企业可以导出月度账单。

### Phase 5：支付系统

目标：支持 Stripe、支付宝、微信支付。

建议接入顺序：

1. 管理员手动充值。
2. Stripe。
3. 支付宝。
4. 微信支付。

任务：

1. `orders` 表。
2. `payments` 表。
3. `payment_callbacks` 表。
4. 创建充值订单。
5. Stripe Checkout 或 Payment Link。
6. Stripe webhook 验签。
7. 支付成功后加余额。
8. 支付失败记录。
9. 幂等处理。
10. 支付宝接入。
11. 微信支付接入。

订单状态：

- `created`
- `pending`
- `paid`
- `failed`
- `cancelled`
- `refunded`
- `expired`

验收：

- 用户能创建充值订单。
- 支付成功能自动加余额。
- 重复回调不会重复加余额。
- 支付失败不会加余额。
- 管理员能查询订单。
- 用户能查看充值记录。

### Phase 6：代理分销第一版

目标：先记录代理关系，不急着自动返佣。

任务：

1. `invite_codes` 表。
2. `referral_relations` 表。
3. `agent_profiles` 表。
4. 用户注册时绑定邀请码。
5. 管理员查看代理名下用户。
6. 管理员查看代理名下消耗。
7. 管理员手动结算返佣。

第二版再做：

- 自动返佣。
- 冻结期。
- 提现申请。
- 提现审核。
- 退款扣回。
- 代理后台。
- 代理专属链接。

验收：

- 每个用户可以归属一个代理。
- 代理能看到下级客户数量和消耗。
- 管理员可以手动录入返佣。
- 返佣记录能追溯。

### Phase 7：企业合同制

目标：支持真正的 B 端客户。

任务：

1. `enterprise_contracts` 表。
2. 企业专属价格。
3. 企业专属模型白名单。
4. 企业专属 Provider 池。
5. 企业月度预算。
6. 企业月结账单。
7. 企业账单导出。
8. 企业 IP 白名单。
9. 企业审计日志。
10. 企业项目成本中心。

`enterprise_contracts` 字段：

```text
- organization_id
- contract_name
- billing_mode
- start_date
- end_date
- credit_limit
- monthly_commitment
- payment_terms
- dedicated_channels
- custom_prices
- sla_level
- status
```

验收：

- 企业可以按月结算。
- 企业可以拥有专属价格。
- 企业可以拥有专属渠道。
- 企业可以导出账单。
- 企业可以按项目查看成本。

### Phase 8：多协议升级

目标：从 OpenAI-compatible 扩展到 Claude、Gemini、Responses。

协议优先级：

1. OpenAI-compatible Chat Completions。
2. Claude-compatible Messages API。
3. OpenAI Responses API。
4. Gemini-compatible API。
5. Realtime / Live API。

关键设计：内部统一 Message IR。

```text
InternalMessageRequest
- model
- messages
- system
- tools
- tool_choice
- temperature
- max_tokens
- stream
- metadata
- multimodal_parts
```

转换链路：

```text
OpenAI request → Internal IR → Provider request
Claude request → Internal IR → Provider request
Gemini request → Internal IR → Provider request
Responses request → Internal IR → Provider request
```

任务：

1. 设计 Internal Message IR。
2. OpenAI Chat Completions → IR。
3. IR → OpenAI-like Provider。
4. Claude Messages → IR。
5. IR → Claude-like Provider。
6. Gemini generateContent → IR。
7. IR → Gemini Provider。
8. Responses API → IR。
9. 工具调用统一。
10. 流式输出统一。

验收：

- 同一个 Provider 可以被多个协议入口调用。
- 同一个模型可以暴露为多个协议格式。
- 工具调用能正常转发。
- 流式输出不丢 chunk。
- usage 能统一计费。

## 7. 备用 Provider 开发规范

### 7.1 统一包装

不管底层是什么，进入主系统都必须包装成 `GenericProviderAdapter` 或具体 `ProviderAdapter`。

必须提供：

- 统一请求格式。
- 统一响应格式。
- 统一错误分类。
- 统一 usage 提取。
- 统一健康检查。
- 统一开关。
- 统一日志。

禁止绕过：

- 鉴权。
- 计费。
- 日志。
- 风控。
- 项目权限。

### 7.2 默认策略

备用 Provider 默认：

- `internal_only`
- 不参与公开模型路由。
- 不参与企业正式流量。
- 不展示给普通用户。
- 不在价格页宣传。
- 不保存正文。
- 低并发。
- 可一键关闭。

### 7.3 风险隔离字段

建议字段：

```text
- risk_level
- data_sensitivity_allowed
- max_concurrency
- max_daily_requests
- max_daily_tokens
- allowed_projects
- manual_enable_required
- visible_to_user
```

示例：

```yaml
experimental_proxy:
  risk_level: high
  data_sensitivity_allowed: low
  max_concurrency: 2
  max_daily_requests: 500
  available_scope: internal_only
```

## 8. 管理后台规划

第一版：

- 用户管理。
- 企业管理。
- 项目管理。
- API Key 管理。
- Provider 管理。
- Channel 管理。
- 模型管理。
- 价格倍率。
- 请求日志。
- 余额管理。
- 手动充值。
- 错误日志。

第二版：

- 订单管理。
- 支付管理。
- 代理管理。
- 返佣管理。
- 企业合同。
- 账单导出。
- 渠道健康状态。
- Provider 账号池状态。

第三版：

- SLA 面板。
- 成本利润面板。
- 模型利润率。
- 上游消耗统计。
- 异常调用告警。
- 渠道自动熔断记录。

## 9. 用户后台规划

第一版：

- API Key 创建。
- 余额查看。
- 用量查看。
- 充值记录。
- 调用日志。
- 模型列表。
- API 文档。

第二版：

- 项目管理。
- 成员管理。
- 项目预算。
- IP 白名单。
- 模型白名单。
- Webhook 告警。

第三版：

- 企业账单。
- 成本分析。
- 成员消耗排行。
- 项目消耗排行。
- 代理后台。

## 10. 部署架构

第一阶段使用 Docker Compose：

```text
nginx / caddy
gateway-new-api
mysql / postgres
redis
provider-service-kiro
provider-service-extra
worker
admin-web
```

组件职责：

- `gateway-new-api`：主 API。
- `worker`：日志聚合、账单统计、健康检查。
- `redis`：限流、队列、缓存。
- `mysql/postgres`：主数据库。
- `nginx/caddy`：反向代理、HTTPS、基础限流。

第二阶段升级：

- API 多副本。
- 数据库主从。
- Redis 持久化。
- 日志系统。
- 对象存储。
- 监控告警。

## 11. 监控与告警

第一版基础指标：

- 请求数。
- 成功率。
- 失败率。
- 平均延迟。
- P95 延迟。
- 每个 Provider 成功率。
- 每个 Channel 成功率。
- 每个模型消耗。
- 每个用户消耗。
- 每个企业消耗。
- 余额不足次数。
- 上游限流次数。
- 超时次数。

告警规则：

- 某 Provider 连续失败 N 次 → 自动熔断。
- 某 Channel 成功率低于阈值 → 降权。
- 某用户消耗异常增长 → 通知管理员。
- 某企业余额低于阈值 → 通知客户。
- 某备用 Provider 被外部用户调用 → 高危告警。

## 12. 安全设计

第一版必须做：

- API Key hash 存储。
- 上游密钥加密存储。
- 管理员 2FA 预留。
- IP 白名单。
- 请求限流。
- 项目限额。
- 用户限额。
- 日志脱敏。
- 操作审计。

禁止保存：

- 完整上游 key 明文。
- 完整用户 prompt。
- 完整用户 response。
- 支付回调敏感原文长期保存。

## 13. 第一版范围冻结

第一版明确只做：

- `new-api` 部署。
- OpenAI-compatible API。
- 2 个正规海外 Provider。
- 1 个 `kiro-gateway` `experimental_proxy` Provider。
- Provider 类型字段。
- Provider 风险等级。
- 用户注册。
- API Key。
- 余额。
- 按量计费。
- 调用日志。
- 管理员手动充值。
- 企业/项目表结构预留。
- `experimental_proxy` `internal_only`。

第一版不做：

- 自动支付。
- 自动返佣。
- 完整企业合同。
- 完整多协议。
- 复杂代理后台。
- 复杂组织权限。

## 14. AI Agent 开发任务拆分

### 任务 1：阅读项目结构

目标：梳理用户、token、渠道、计费、日志、后台相关模块。

输出：

1. 项目目录说明。
2. 核心数据表说明。
3. 请求调用链路。
4. 扣费逻辑位置。
5. 渠道路由逻辑位置。
6. 适合二开的文件列表。

### 任务 2：设计 Provider 插件体系

目标：在不破坏原有渠道逻辑的基础上，设计 Provider Adapter 抽象。

输出：

1. `ProviderAdapter` interface。
2. Provider 类型枚举。
3. Provider 配置表。
4. `ProviderAccount` 表。
5. Channel 与 Provider 的关系。
6. `KiroGatewayAdapter` 示例。

### 任务 3：新增 experimental_proxy 类型

目标：支持 `kiro-gateway` 和其他备用 Provider 作为 `experimental_proxy`，但默认 `internal_only`。

验收：

1. 管理员可以新增 `experimental_proxy`。
2. 普通用户默认不可见。
3. internal 用户可调用。
4. 日志中标记 `provider_type`。
5. 可一键禁用。

### 任务 4：组织与项目结构

目标：预留 B 端组织、项目、成员、项目级 Key。

验收：

1. `users` 不再是唯一主体。
2. API Key 可绑定 project。
3. `usage_logs` 可绑定 organization/project。
4. 项目可设置预算。
5. 项目可设置模型白名单。

### 任务 5：计费流水系统

目标：将余额变化全部流水化。

验收：

1. 充值写 ledger。
2. 消耗写 ledger。
3. 退款写 ledger。
4. 管理员调整写 ledger。
5. 每条 `usage_log` 能关联 `billing_ledger`。

### 任务 6：支付系统

目标：先接 Stripe，再接支付宝和微信。

验收：

1. 创建订单。
2. 支付回调验签。
3. 幂等处理。
4. 支付成功加余额。
5. 用户能查充值记录。
6. 管理员能查订单。

### 任务 7：代理邀请码

目标：先做轻量代理体系。

验收：

1. 代理可以生成邀请码。
2. 新用户注册绑定代理。
3. 管理员能看代理客户。
4. 管理员能看代理客户消耗。
5. 管理员能手动记录返佣。

### 任务 8：Claude Messages API

目标：新增 Claude-compatible Messages API 入口。

验收：

1. 支持 `/v1/messages` 或 `/anthropic/v1/messages`。
2. 支持 `system` / `messages` / `max_tokens` / `stream`。
3. 能转成 Internal IR。
4. 能路由到 Claude-like Provider。
5. usage 能计费。

### 任务 9：Responses API

目标：新增 OpenAI Responses API 入口。

验收：

1. 支持 `/v1/responses`。
2. 支持 `input` / `instructions` / `model` / `stream` / `tools`。
3. 能转成 Internal IR。
4. 能返回 Responses 格式。
5. 流式事件兼容。

### 任务 10：Gemini-compatible API

目标：新增 Gemini `generateContent` 兼容入口。

验收：

1. 支持 `generateContent`。
2. 支持 `streamGenerateContent`。
3. 支持 `contents` / `parts`。
4. 支持多模态字段预留。
5. usage 能计费。

## 15. 第一条 Mission

可直接给 Droid Mission / Claude Code / Codex 使用：

```text
目标：
基于 new-api 二开一个面向 B 端的多 Provider API 网关 MVP。当前要求只完成 OpenAI-compatible API、Provider 类型扩展、experimental_proxy 隔离、组织/项目表结构预留、基础计费日志闭环。

背景：
kiro-gateway 已经跑通，但它只能作为 experimental_proxy 的一个 provider。系统不能写死 kiro-gateway，而要设计 ProviderAdapter 抽象，未来可接入 official_cloud、aggregator、authorized_proxy、experimental_proxy 多种 provider。

第一阶段范围：
1. 阅读 new-api 项目结构。
2. 找到用户、token、渠道、计费、日志相关模块。
3. 设计 ProviderAdapter interface。
4. 新增 provider_type 字段：official_cloud、aggregator、authorized_proxy、experimental_proxy。
5. 新增 risk_level、available_scope 字段。
6. experimental_proxy 默认 internal_only。
7. 新增 organization、project、organization_member 基础表。
8. api_key 和 usage_log 预留 organization_id、project_id。
9. 新增 KiroGatewayAdapter 作为示例 experimental_proxy。
10. 验证普通用户不能调用 experimental_proxy，internal 用户可以调用。
11. 输出数据库迁移、配置说明和测试用例。

验收：
- 普通用户可调用 official_cloud provider。
- internal 用户可调用 kiro-gateway experimental_proxy。
- 所有请求都记录 usage_log。
- usage_log 包含 user_id、organization_id、project_id、provider_type、channel_id、model、tokens、cost、status。
- experimental_proxy 可在后台一键关闭。
- 不破坏 new-api 原有基础功能。
```
