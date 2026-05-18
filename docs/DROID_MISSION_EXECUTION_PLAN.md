# Droid Mission Execution Plan

本文档是可直接复制给 Droid Mission 的总提示词与执行方案，用于基于 new-api 二次开发一个面向 B 端的 Multi-Provider AI API Gateway MVP。

当前阶段只生成完整 Mission 执行方案，不直接写业务代码。确认后再从 M0 开始逐个 milestone 执行。

## 1. 直接复制给 Droid Mission 的总提示词

```text
我要使用 Droid Mission 模式，基于 new-api 二开一个面向 B 端的 Multi-Provider AI API Gateway MVP。

请严格按照 Mission 模式执行：先规划、再创建 validation-contract、再生成 features.json、再按 milestone 执行。当前阶段不要直接写业务代码，先生成完整 Mission 执行方案。

项目背景：
1. 主干基于 new-api 二次开发。
2. 第一阶段只暴露 OpenAI-compatible API。
3. 后续再扩展 Claude-compatible Messages API、Gemini-compatible API、OpenAI Responses API。
4. kiro-gateway 只作为 experimental_proxy 的一个 Provider Adapter 示例。
5. 系统不能写死 kiro-gateway，而要设计 Provider Adapter Framework。
6. Provider 类型包括 official_cloud、aggregator、authorized_proxy、experimental_proxy。
7. experimental_proxy 默认 risk_level=high、available_scope=internal_only、visibility=internal_only、manual_enable_required=true、enabled=false。
8. 第一版必须跑通：用户 API Key、GET /v1/models、POST /v1/chat/completions、Provider 路由、usage_log、token 统计、余额扣费、余额不足拒绝、管理员手动充值、internal 用户可调用 experimental_proxy、普通用户不能调用 experimental_proxy、experimental_proxy 可一键禁用。
9. 第一版只预留 organization/project，不做复杂企业权限。
10. 第一版不做自动支付、自动返佣、完整企业合同、多协议、复杂代理后台。

执行原则：
1. Orchestrator 只做规划、分配、验收，不直接写业务代码。
2. Worker 负责实现 feature。
3. Validator 负责检查 feature 是否满足 validation-contract。
4. 每个 feature 只能做一个小闭环。
5. 每个 feature 最多绑定 1～3 个 validation assertion。
6. 一个 milestone 尽量只涉及 1～2 类 worker。
7. 不允许一个 feature 同时修改 Provider、Billing、Admin UI、Routing 多个大模块。
8. 不允许跳过测试。
9. 不允许明文保存上游 key。
10. 不允许明文保存用户 API Key。
11. 默认不保存完整 prompt。
12. 默认不保存完整 response。
13. official_cloud 失败时不能默认 fallback 到 experimental_proxy。
14. experimental_proxy 不能展示给普通用户。
15. 普通用户默认不能调用 experimental_proxy。
16. internal 用户只有在 explicit enable 后才能调用 experimental_proxy。
17. 不破坏 new-api 原有基础功能。
18. 所有数据库设计必须尽量兼容 SQLite、MySQL、PostgreSQL，优先使用 GORM 抽象，避免数据库专属 SQL。

请先完成以下输出，不要开始写业务代码：

一、仓库调研计划
- 找出用户模块、token/API Key 模块、channel 模块、计费模块、日志模块、后台模块。
- 梳理 /v1/chat/completions 从入口到上游请求、计费、日志的完整调用链路。
- 输出适合二开的文件列表。

二、Mission 文件
请创建或规划以下文件：
- AGENTS.md
- mission.md
- validation-contract.md
- features.json
- .factory/services.yaml
- .factory/init.sh
- .factory/library/architecture.md
- .factory/library/environment.md
- .factory/library/user-testing.md

三、Worker Skills
请创建或规划以下 worker：
- repo-discovery-worker
- db-migration-worker
- provider-worker
- routing-security-worker
- openai-api-worker
- billing-worker
- admin-ui-worker
- deployment-worker
- test-validation-worker

四、Milestone 拆分
按 M0～M16 执行：
M0：仓库调研与 Mission 基础文件
M1：Provider 类型与字段基础迁移
M2：Provider Account 与 Channel 关联
M3：OpenAI-compatible 接口回归
M4：正规 Provider 最小接入
M5：模型映射与价格基础
M6：KiroGatewayAdapter 骨架
M7：experimental_proxy 访问控制
M8：experimental_proxy 路由与 fallback 隔离
M9：组织 / 项目基础表结构
M10：API Key 组织项目绑定与安全
M11：usage_log 与隐私策略
M12：token 统计与基础扣费
M13：余额不足与管理员手动充值
M14：管理后台 Provider / Channel 最小页面
M15：管理后台 API Key / usage_log / 余额页面
M16：Docker、文档和回归测试

五、Feature 格式
每个 feature 必须包含：
- id
- milestone
- worker
- description
- scope
- files_to_inspect
- files_to_modify
- preconditions
- implementation_steps
- verification_steps
- fulfills
- done_definition
- risk_notes
- status

请先输出：
1. mission.md 草案
2. validation-contract.md 草案
3. features.json 草案
4. 每个 milestone 的验收标准
5. 每个 worker 的职责说明
6. 推荐执行顺序

等待我确认后，再开始执行 M0。
```

## 2. mission.md 草案

```markdown
# mission.md

## Mission

基于 new-api 二次开发一个面向 B 端的 Multi-Provider AI API Gateway MVP。第一阶段只暴露 OpenAI-compatible API，支持多 Provider 元数据、Provider Account、Channel 路由、安全隔离、usage log、token 统计、余额扣费、余额不足拒绝、管理员手动充值和最小管理后台。

## Non-Goals

- 不实现自动支付。
- 不实现自动返佣。
- 不实现完整企业合同。
- 不实现复杂企业权限。
- 不实现多协议完整支持。
- 不实现复杂代理后台。
- 不把 kiro-gateway 写死到主流程。

## Provider Policy

Provider 类型包括：
- official_cloud
- aggregator
- authorized_proxy
- experimental_proxy

experimental_proxy 默认策略：
- risk_level = high
- available_scope = internal_only
- visibility = internal_only
- manual_enable_required = true
- enabled = false

official_cloud 失败时不能默认 fallback 到 experimental_proxy。普通用户默认不能看到或调用 experimental_proxy。internal 用户也必须在 explicit enable 后才能调用 experimental_proxy。

## Execution Rules

Orchestrator 只做规划、分配、验收，不直接写业务代码。Worker 负责实现 feature。Validator 负责检查 feature 是否满足 validation-contract。每个 feature 只能做一个小闭环，最多绑定 1～3 个 validation assertion。

所有数据库设计必须兼容 SQLite、MySQL、PostgreSQL，优先使用 GORM 抽象，避免数据库专属 SQL。

不得明文保存上游 key，不得明文保存用户 API Key。默认不保存完整 prompt 和完整 response。debug payload 日志如存在必须默认关闭，并带保留时间限制。
```

## 3. validation-contract.md 草案

```markdown
# validation-contract.md

## Discovery

### VAL-DISC-001
系统必须能够说明 new-api 的用户、API Key/token、channel、计费、日志、后台模块位置。

### VAL-DISC-002
系统必须能够说明 /v1/chat/completions 从入口、鉴权、路由、上游请求、usage、扣费、日志的完整链路。

## Mission Files

### VAL-MISSION-001
Mission 所需共享状态文件必须存在，包括 AGENTS.md、mission.md、validation-contract.md、features.json、.factory/services.yaml、.factory/init.sh。

### VAL-MISSION-002
services.yaml 必须包含本地启动、停止、测试、构建、healthcheck 命令。

## Provider

### VAL-PROVIDER-001
系统必须支持 official_cloud、aggregator、authorized_proxy、experimental_proxy 四种 provider_type。

### VAL-PROVIDER-002
provider_type 必须在数据库、后端模型、常量定义和后台配置中保持一致。

### VAL-PROVIDER-003
Provider 必须支持 risk_level 字段。

### VAL-PROVIDER-004
Provider 必须支持 available_scope 字段。

### VAL-PROVIDER-005
Provider 必须支持 manual_enable_required 和 enabled 字段。

### VAL-PROVIDER-006
一个 Provider 必须可以绑定多个 Provider Account。

### VAL-PROVIDER-007
Channel 必须可以关联 Provider。

### VAL-PROVIDER-008
Channel 必须可以关联 Provider Account。

### VAL-PROVIDER-009
旧 Channel 在迁移后不能失效。

## Experimental Proxy

### VAL-EXP-001
experimental_proxy 新建时默认 risk_level=high。

### VAL-EXP-002
experimental_proxy 新建时默认 available_scope=internal_only。

### VAL-EXP-003
experimental_proxy 新建时默认 manual_enable_required=true。

### VAL-EXP-004
experimental_proxy 新建时默认 enabled=false。

### VAL-EXP-005
普通用户默认不能看到 experimental_proxy。

### VAL-EXP-006
普通用户调用 experimental_proxy 必须被拒绝。

### VAL-EXP-007
internal 用户在显式启用后可以调用 experimental_proxy。

### VAL-EXP-008
disabled experimental_proxy 不能被任何用户调用。

### VAL-EXP-009
experimental_proxy 可在后台一键禁用。

## Routing

### VAL-ROUTE-001
系统至少可以配置两个非 experimental Provider。

### VAL-ROUTE-002
普通用户可以调用 official_cloud 或 aggregator Provider。

### VAL-ROUTE-003
official_cloud 失败时不能默认 fallback 到 experimental_proxy。

### VAL-ROUTE-004
allow_experimental=false 时，experimental_proxy 不能进入候选通道。

### VAL-ROUTE-005
只有 allow_experimental=true 且用户或项目具备权限时，experimental_proxy 才能进入候选通道。

## OpenAI-compatible API

### VAL-OAI-001
GET /v1/models 必须返回 OpenAI-compatible 模型列表。

### VAL-OAI-002
POST /v1/chat/completions 必须支持 OpenAI SDK 调用。

### VAL-OAI-003
如果原项目支持 stream，POST /v1/chat/completions 必须保持 stream 兼容。

## Model Mapping

### VAL-MODEL-001
系统必须支持 public_model_name 到 provider_model_name 的映射。

### VAL-MODEL-002
系统必须支持 input_price、output_price、billing_multiplier。

### VAL-MODEL-003
disabled model 不能被调用。

## Kiro Adapter

### VAL-KIRO-001
KiroGatewayAdapter 必须实现 ProviderAdapter 基础接口。

### VAL-KIRO-002
KiroGatewayAdapter 必须被标记为 experimental_proxy。

### VAL-KIRO-003
KiroGatewayAdapter 不得被写死到主请求流程中。

### VAL-KIRO-004
KiroGatewayAdapter 默认 disabled。

## Organization / Project

### VAL-ORG-001
系统必须支持 organizations 表。

### VAL-ORG-002
系统必须支持 organization_members 表。

### VAL-ORG-003
系统必须支持 projects 表。

### VAL-ORG-004
organization 可以拥有 project。

### VAL-ORG-005
user 可以归属 organization。

## API Key

### VAL-KEY-001
API Key 可以绑定 user_id。

### VAL-KEY-002
API Key 可以绑定 organization_id。

### VAL-KEY-003
API Key 可以绑定 project_id。

### VAL-KEY-004
API Key 只保存 key_hash 和 key_prefix。

### VAL-KEY-005
API Key 完整明文只展示一次。

### VAL-KEY-006
API Key 可以禁用。

### VAL-KEY-007
API Key 支持 allowed_models。

### VAL-KEY-008
API Key 支持 allowed_provider_types。

### VAL-KEY-009
普通 API Key 默认不允许 experimental_proxy。

## Usage Log

### VAL-USAGE-001
成功请求必须写入 usage_log。

### VAL-USAGE-002
失败请求必须写入 usage_log 或 error log。

### VAL-USAGE-003
usage_log 必须包含 request_id、user_id、organization_id、project_id、api_key_id。

### VAL-USAGE-004
usage_log 必须包含 model、public_model、provider_id、provider_type、channel_id、provider_account_id。

### VAL-USAGE-005
usage_log 必须包含 input_tokens、output_tokens、total_tokens、cost、charged_amount、status。

### VAL-USAGE-006
usage_log 必须包含 error_code、error_message_sanitized、latency_ms、ip、user_agent。

## Privacy

### VAL-PRIVACY-001
默认不保存完整 prompt。

### VAL-PRIVACY-002
默认不保存完整 response。

### VAL-PRIVACY-003
error_message 必须脱敏。

### VAL-PRIVACY-004
debug_payload_logs 如存在，必须默认关闭，并且必须有保留时间限制。

## Billing

### VAL-BILL-001
系统必须能够提取 input_tokens。

### VAL-BILL-002
系统必须能够提取 output_tokens。

### VAL-BILL-003
系统必须能够计算 total_tokens。

### VAL-BILL-004
系统必须根据模型价格和 token usage 计算 cost。

### VAL-BILL-005
成功请求后余额必须减少。

### VAL-BILL-006
失败请求不应错误扣费。

### VAL-BILL-007
余额不足时必须拒绝请求。

### VAL-BILL-008
余额不足时不能发起上游请求。

### VAL-BILL-009
余额不足请求必须写入日志。

### VAL-BILL-010
管理员可以手动增加用户余额。

### VAL-BILL-011
管理员手动充值必须有操作记录。

## Admin

### VAL-ADMIN-001
管理员可以查看 Provider 列表。

### VAL-ADMIN-002
Provider 列表显示 provider_type、risk_level、available_scope、enabled。

### VAL-ADMIN-003
管理员可以启用或禁用 Provider。

### VAL-ADMIN-004
管理员可以查看 Channel 与 Provider / Provider Account 的关系。

### VAL-ADMIN-005
管理员可以查看 API Key 列表。

### VAL-ADMIN-006
管理员可以禁用 API Key。

### VAL-ADMIN-007
管理员可以查看 usage_log。

### VAL-ADMIN-008
usage_log 页面显示 user、model、provider_type、tokens、cost、status、latency。

### VAL-ADMIN-009
管理员可以查看用户余额。

### VAL-ADMIN-010
管理员可以手动调整余额。

## Deploy

### VAL-DEPLOY-001
docker compose 可以启动 gateway、db、redis。

### VAL-DEPLOY-002
gateway healthcheck 必须通过。

### VAL-DEPLOY-003
db 和 redis healthcheck 必须通过。

### VAL-DEPLOY-004
.env.example 不包含真实密钥。

## Docs

### VAL-DOC-001
文档必须说明 kiro-gateway 是 external experimental provider。

### VAL-DOC-002
文档必须说明 experimental_proxy 默认 internal_only，普通用户不可调用。

### VAL-DOC-003
文档必须说明 official_cloud 不能默认 fallback 到 experimental_proxy。

### VAL-DOC-004
文档必须包含 OpenAI SDK 调用示例。

### VAL-DOC-005
文档必须包含 API Key 使用示例。

## Regression Tests

### VAL-TEST-001
存在普通用户调用 official provider 的测试。

### VAL-TEST-002
存在 internal 用户调用 experimental provider 的测试。

### VAL-TEST-003
存在普通用户调用 experimental provider 被拒绝的测试。

### VAL-TEST-004
存在 disabled experimental provider 被拒绝的测试。

### VAL-TEST-005
存在余额不足被拒绝且不调用上游的测试。

### VAL-TEST-006
存在默认不保存完整 prompt/response 的测试或检查。
```

## 4. features.json 草案

下面是给 Droid 写入 `features.json` 的结构化草案。每个 feature 都必须保持小闭环，且不得跨 Provider、Billing、Admin UI、Routing 多个大模块。

```json
{
  "milestones": [
    {
      "id": "M0",
      "name": "仓库调研与 Mission 基础文件",
      "goal": "只读仓库，建立 Mission 共享上下文，不修改业务逻辑。",
      "features": [
        {
          "id": "M0-F01-repo-structure-scan",
          "milestone": "M0",
          "worker": "repo-discovery-worker",
          "description": "扫描 new-api 仓库结构，识别用户、API Key/token、channel、计费、日志、后台模块位置。",
          "scope": ["read-only", "documentation"],
          "files_to_inspect": ["README", "router", "controller", "model", "service", "middleware", "web", "admin"],
          "files_to_modify": ["docs/REPO_STRUCTURE.md"],
          "preconditions": ["仓库可以打开"],
          "implementation_steps": ["读取目录结构", "识别核心模块", "输出适合二开的文件列表"],
          "verification_steps": ["检查 docs/REPO_STRUCTURE.md 是否存在", "确认核心模块路径完整"],
          "fulfills": ["VAL-DISC-001"],
          "done_definition": "文档能说明 new-api 的核心模块位置。",
          "risk_notes": "禁止修改业务代码。",
          "status": "pending"
        },
        {
          "id": "M0-F02-openai-request-flow",
          "milestone": "M0",
          "worker": "repo-discovery-worker",
          "description": "梳理 /v1/chat/completions 从入口到上游请求、计费和日志的完整链路。",
          "scope": ["read-only", "documentation"],
          "files_to_inspect": ["router", "middleware", "controller", "relay", "service", "model"],
          "files_to_modify": ["docs/OPENAI_REQUEST_FLOW.md"],
          "preconditions": ["M0-F01 完成"],
          "implementation_steps": ["定位路由", "定位鉴权", "定位 channel 选择", "定位上游请求", "定位扣费", "定位日志"],
          "verification_steps": ["文档包含入口、鉴权、路由、上游请求、usage、扣费、日志节点"],
          "fulfills": ["VAL-DISC-002"],
          "done_definition": "能用链路图说明请求全过程。",
          "risk_notes": "禁止修改业务代码。",
          "status": "pending"
        },
        {
          "id": "M0-F03-create-mission-files",
          "milestone": "M0",
          "worker": "repo-discovery-worker",
          "description": "创建 Mission 基础文件和 .factory 共享上下文。",
          "scope": ["documentation", "mission-config"],
          "files_to_inspect": ["docs/REPO_STRUCTURE.md", "docs/OPENAI_REQUEST_FLOW.md"],
          "files_to_modify": ["AGENTS.md", "mission.md", "validation-contract.md", "features.json", ".factory/services.yaml", ".factory/init.sh", ".factory/library/architecture.md", ".factory/library/environment.md", ".factory/library/user-testing.md"],
          "preconditions": ["M0-F01 完成", "M0-F02 完成"],
          "implementation_steps": ["创建 Mission 文件", "写入项目边界", "写入执行规则", "写入禁止事项"],
          "verification_steps": ["检查所有文件存在", "检查 services.yaml 有启动、停止、测试、构建、healthcheck 命令"],
          "fulfills": ["VAL-MISSION-001", "VAL-MISSION-002"],
          "done_definition": "Mission 文件齐全，后续 Worker 可以读取统一上下文。",
          "risk_notes": "不修改业务逻辑。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M1",
      "name": "Provider 类型与字段基础迁移",
      "goal": "新增 Provider 类型、风险字段和 experimental_proxy 默认值。",
      "features": [
        {
          "id": "M1-F01-provider-type-enum",
          "milestone": "M1",
          "worker": "db-migration-worker",
          "description": "新增 official_cloud、aggregator、authorized_proxy、experimental_proxy 四类 provider_type。",
          "scope": ["database", "backend-model"],
          "files_to_inspect": ["channel model", "provider model", "migration", "constants"],
          "files_to_modify": ["model", "migration", "constants"],
          "preconditions": ["M0 完成"],
          "implementation_steps": ["找到现有 channel/provider 数据结构", "增加 provider_type 字段", "增加四种枚举值", "保持旧 channel 兼容"],
          "verification_steps": ["migration 可运行", "旧 channel 默认映射为合理默认类型", "四种枚举值可保存"],
          "fulfills": ["VAL-PROVIDER-001", "VAL-PROVIDER-002"],
          "done_definition": "provider_type 在 DB、后端模型、常量定义中一致。",
          "risk_notes": "不能破坏旧 channel。",
          "status": "pending"
        },
        {
          "id": "M1-F02-provider-risk-scope-fields",
          "milestone": "M1",
          "worker": "db-migration-worker",
          "description": "新增 risk_level、available_scope、manual_enable_required、enabled 字段。",
          "scope": ["database", "backend-model"],
          "files_to_inspect": ["provider model", "migration"],
          "files_to_modify": ["model", "migration"],
          "preconditions": ["M1-F01 完成"],
          "implementation_steps": ["新增 risk_level", "新增 available_scope", "新增 manual_enable_required", "新增 enabled", "设置兼容默认值"],
          "verification_steps": ["字段可读写", "available_scope 支持 all_users、enterprise_only、selected_groups、internal_only、admin_only"],
          "fulfills": ["VAL-PROVIDER-003", "VAL-PROVIDER-004", "VAL-PROVIDER-005"],
          "done_definition": "Provider 风险和可见范围字段可读写。",
          "risk_notes": "默认值不能导致现有 Provider 全部不可用。",
          "status": "pending"
        },
        {
          "id": "M1-F03-experimental-default-policy",
          "milestone": "M1",
          "worker": "routing-security-worker",
          "description": "experimental_proxy 新建时自动应用 high risk、internal_only、manual_enable_required、disabled 默认值。",
          "scope": ["backend-policy"],
          "files_to_inspect": ["provider create/update logic"],
          "files_to_modify": ["provider create/update logic", "tests"],
          "preconditions": ["M1-F01 完成", "M1-F02 完成"],
          "implementation_steps": ["识别 provider_type=experimental_proxy", "套用默认值", "增加测试"],
          "verification_steps": ["创建 experimental_proxy 后检查默认字段", "创建 official_cloud 后不套用 experimental 默认值"],
          "fulfills": ["VAL-EXP-001", "VAL-EXP-002", "VAL-EXP-003", "VAL-EXP-004"],
          "done_definition": "experimental_proxy 默认不可公开使用。",
          "risk_notes": "不能影响 official_cloud 正常启用。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M2",
      "name": "Provider Account 与 Channel 关联",
      "goal": "建立 Provider、Provider Account、Channel 三者关系，并保持旧 channel 兼容。",
      "features": [
        {
          "id": "M2-F01-provider-accounts-table",
          "milestone": "M2",
          "worker": "db-migration-worker",
          "description": "新增 provider_accounts 表，支持一个 Provider 下多个账号、key、endpoint。",
          "scope": ["database", "security"],
          "files_to_inspect": ["model", "migration", "encryption service if exists"],
          "files_to_modify": ["model", "migration", "encryption service if exists"],
          "preconditions": ["M1 完成"],
          "implementation_steps": ["新增 provider_accounts model", "包含 credential_encrypted 等字段", "写 migration", "禁止明文保存 credential"],
          "verification_steps": ["一个 provider 可创建多个 provider_account", "credential 字段不保存明文"],
          "fulfills": ["VAL-PROVIDER-006"],
          "done_definition": "Provider Account 可持久化，凭证字段使用加密或加密占位接口。",
          "risk_notes": "没有加密实现时，必须阻断真实 key 保存。",
          "status": "pending"
        },
        {
          "id": "M2-F02-channel-provider-link",
          "milestone": "M2",
          "worker": "provider-worker",
          "description": "Channel 增加 provider_id 和 provider_account_id 关联。",
          "scope": ["database", "backend-model"],
          "files_to_inspect": ["channel model", "channel service"],
          "files_to_modify": ["channel model", "migration", "channel service"],
          "preconditions": ["M2-F01 完成"],
          "implementation_steps": ["Channel 增加 provider_id", "Channel 增加 provider_account_id", "保持旧 channel 可用"],
          "verification_steps": ["channel 可关联 provider", "channel 可关联 provider_account"],
          "fulfills": ["VAL-PROVIDER-007", "VAL-PROVIDER-008"],
          "done_definition": "路由层能拿到 channel 对应的 provider 和 provider_account。",
          "risk_notes": "不要强制旧数据必须立刻填写 provider_account_id。",
          "status": "pending"
        },
        {
          "id": "M2-F03-legacy-channel-compatibility",
          "milestone": "M2",
          "worker": "test-validation-worker",
          "description": "验证旧 Channel 在迁移后不失效。",
          "scope": ["test", "compatibility"],
          "files_to_inspect": ["channel model", "routing tests"],
          "files_to_modify": ["tests"],
          "preconditions": ["M2-F02 完成"],
          "implementation_steps": ["准备旧 channel 测试数据", "运行迁移", "验证旧 channel 可读取和调用"],
          "verification_steps": ["旧 channel 不报错", "旧 channel 可被路由"],
          "fulfills": ["VAL-PROVIDER-009"],
          "done_definition": "Provider 结构改造不破坏 new-api 原有 channel。",
          "risk_notes": "兼容性优先于新字段强约束。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M3",
      "name": "OpenAI-compatible 接口回归",
      "goal": "证明 Provider 改造没有破坏原有 OpenAI-compatible API。",
      "features": [
        {
          "id": "M3-F01-verify-models-endpoint",
          "milestone": "M3",
          "worker": "openai-api-worker",
          "description": "验证 GET /v1/models 保持 OpenAI-compatible。",
          "scope": ["test"],
          "files_to_inspect": ["router", "controller", "tests"],
          "files_to_modify": ["tests", "docs if needed"],
          "preconditions": ["M2 完成"],
          "implementation_steps": ["编写 GET /v1/models 测试", "检查返回格式"],
          "verification_steps": ["curl /v1/models 成功", "OpenAI SDK 可读取模型列表"],
          "fulfills": ["VAL-OAI-001"],
          "done_definition": "/v1/models 返回 OpenAI-compatible 模型列表。",
          "risk_notes": "不重构主接口。",
          "status": "pending"
        },
        {
          "id": "M3-F02-verify-chat-completions",
          "milestone": "M3",
          "worker": "openai-api-worker",
          "description": "验证 POST /v1/chat/completions 非流式调用。",
          "scope": ["test"],
          "files_to_inspect": ["relay", "controller", "tests"],
          "files_to_modify": ["tests"],
          "preconditions": ["M3-F01 完成"],
          "implementation_steps": ["编写 chat completions 非流式测试", "检查 choices 返回"],
          "verification_steps": ["OpenAI SDK 可调用", "返回包含 choices"],
          "fulfills": ["VAL-OAI-002"],
          "done_definition": "非流式 OpenAI-compatible 主链路可用。",
          "risk_notes": "不要改动协议格式。",
          "status": "pending"
        },
        {
          "id": "M3-F03-verify-chat-streaming-if-supported",
          "milestone": "M3",
          "worker": "openai-api-worker",
          "description": "如果原项目支持 stream，验证 stream 不被破坏。",
          "scope": ["test"],
          "files_to_inspect": ["relay", "controller", "tests"],
          "files_to_modify": ["tests"],
          "preconditions": ["M3-F02 完成"],
          "implementation_steps": ["检查原项目是否支持 stream", "支持则添加 stream 测试", "不支持则记录原因"],
          "verification_steps": ["stream chunk 格式兼容，或文档说明不支持"],
          "fulfills": ["VAL-OAI-003"],
          "done_definition": "stream 能力状态明确且不被误改。",
          "risk_notes": "不强行新增原项目没有的 stream 能力。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M4",
      "name": "正规 Provider 最小接入",
      "goal": "接入两个非 experimental Provider，并让普通用户可调用。",
      "features": [
        {
          "id": "M4-F01-first-official-provider",
          "milestone": "M4",
          "worker": "provider-worker",
          "description": "配置并验证第一个 official_cloud Provider。",
          "scope": ["provider", "routing"],
          "files_to_inspect": ["provider config", "channel config", "routing"],
          "files_to_modify": ["provider config", "channel config", "tests"],
          "preconditions": ["M3 完成"],
          "implementation_steps": ["创建 official_cloud provider", "创建 provider_account", "创建 channel", "绑定模型映射", "完成一次调用"],
          "verification_steps": ["provider enabled=true", "普通用户可调用", "日志能看到 provider_type=official_cloud"],
          "fulfills": ["VAL-ROUTE-002"],
          "done_definition": "普通用户能通过 official_cloud 完成一次请求。",
          "risk_notes": "不要使用 experimental_proxy。",
          "status": "pending"
        },
        {
          "id": "M4-F02-second-normal-provider",
          "milestone": "M4",
          "worker": "provider-worker",
          "description": "配置第二个 official_cloud 或 aggregator Provider。",
          "scope": ["provider", "routing"],
          "files_to_inspect": ["provider config", "channel config", "routing"],
          "files_to_modify": ["provider config", "channel config", "tests"],
          "preconditions": ["M4-F01 完成"],
          "implementation_steps": ["创建第二个非 experimental provider", "创建对应 account/channel", "加入候选路由"],
          "verification_steps": ["系统至少存在两个非 experimental Provider", "路由可选择其中一个"],
          "fulfills": ["VAL-ROUTE-001"],
          "done_definition": "多正规 Provider 基础路由可用。",
          "risk_notes": "不要让 experimental_proxy 进入候选。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M5",
      "name": "模型映射与价格基础",
      "goal": "建立模型公开名、上游名、价格和倍率基础。",
      "features": [
        {
          "id": "M5-F01-model-mapping-table",
          "milestone": "M5",
          "worker": "db-migration-worker",
          "description": "新增或扩展 model_mappings 表。",
          "scope": ["database", "model"],
          "files_to_inspect": ["model", "migration", "setting model config"],
          "files_to_modify": ["model", "migration"],
          "preconditions": ["M4 完成"],
          "implementation_steps": ["增加 public_model_name", "增加 provider_model_name", "增加 protocol_type", "增加 enabled"],
          "verification_steps": ["public model 可映射到 provider model"],
          "fulfills": ["VAL-MODEL-001"],
          "done_definition": "模型名映射可以持久化。",
          "risk_notes": "保持旧模型配置可用。",
          "status": "pending"
        },
        {
          "id": "M5-F02-model-price-fields",
          "milestone": "M5",
          "worker": "db-migration-worker",
          "description": "模型映射支持 input_price、output_price、billing_multiplier。",
          "scope": ["database", "billing-prep"],
          "files_to_inspect": ["model", "migration", "billing config"],
          "files_to_modify": ["model", "migration"],
          "preconditions": ["M5-F01 完成"],
          "implementation_steps": ["增加 input_price", "增加 output_price", "增加 billing_multiplier"],
          "verification_steps": ["价格字段可读写"],
          "fulfills": ["VAL-MODEL-002"],
          "done_definition": "模型价格字段可用于后续计费。",
          "risk_notes": "价格字段类型避免浮点误差，优先考虑整数微单位或 decimal。",
          "status": "pending"
        },
        {
          "id": "M5-F03-disabled-model-rejection",
          "milestone": "M5",
          "worker": "routing-security-worker",
          "description": "disabled model 不可被调用。",
          "scope": ["routing", "security"],
          "files_to_inspect": ["routing logic", "model resolution"],
          "files_to_modify": ["routing logic", "tests"],
          "preconditions": ["M5-F01 完成"],
          "implementation_steps": ["在模型解析或路由前检查 enabled", "disabled 时拒绝"],
          "verification_steps": ["禁用模型调用失败", "启用模型调用成功"],
          "fulfills": ["VAL-MODEL-003"],
          "done_definition": "模型开关生效。",
          "risk_notes": "错误信息不要泄露上游信息。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M6",
      "name": "KiroGatewayAdapter 骨架",
      "goal": "定义 ProviderAdapter，并实现 KiroGatewayAdapter skeleton，但不开放流量。",
      "features": [
        {
          "id": "M6-F01-provider-adapter-interface",
          "milestone": "M6",
          "worker": "provider-worker",
          "description": "定义 ProviderAdapter interface。",
          "scope": ["adapter", "architecture"],
          "files_to_inspect": ["relay/channel", "relay/common", "provider packages"],
          "files_to_modify": ["provider adapter package", "docs/PROVIDER_SPEC.md"],
          "preconditions": ["M5 完成"],
          "implementation_steps": ["定义 Name、Type、SupportedProtocols、SendRequest、StreamRequest、ExtractUsage、HealthCheck、ClassifyError 等方法", "写文档"],
          "verification_steps": ["interface 可被普通 provider 和 kiro adapter 实现"],
          "fulfills": ["VAL-KIRO-001"],
          "done_definition": "ProviderAdapter 抽象存在。",
          "risk_notes": "不要大改主流程。",
          "status": "pending"
        },
        {
          "id": "M6-F02-kiro-gateway-adapter-skeleton",
          "milestone": "M6",
          "worker": "provider-worker",
          "description": "实现 KiroGatewayAdapter skeleton。",
          "scope": ["adapter"],
          "files_to_inspect": ["provider adapter package"],
          "files_to_modify": ["kiro adapter package", "tests"],
          "preconditions": ["M6-F01 完成"],
          "implementation_steps": ["实现基础方法", "Type 返回 experimental_proxy", "默认 disabled"],
          "verification_steps": ["KiroGatewayAdapter 可注册", "Type=experimental_proxy", "默认 disabled"],
          "fulfills": ["VAL-KIRO-002", "VAL-KIRO-004"],
          "done_definition": "Kiro adapter 只是一个可注册 adapter。",
          "risk_notes": "不接入真实用户流量。",
          "status": "pending"
        },
        {
          "id": "M6-F03-no-hardcoded-kiro-flow",
          "milestone": "M6",
          "worker": "test-validation-worker",
          "description": "验证 Kiro 没有被写死到主请求流程。",
          "scope": ["test", "architecture-check"],
          "files_to_inspect": ["relay", "routing", "provider adapter package"],
          "files_to_modify": ["tests", "docs/PROVIDER_SPEC.md"],
          "preconditions": ["M6-F02 完成"],
          "implementation_steps": ["检查主流程是否通过 ProviderAdapter 调用", "确认没有硬编码 Kiro 分支"],
          "verification_steps": ["主流程只依赖 ProviderAdapter", "Kiro 不出现在默认路由"],
          "fulfills": ["VAL-KIRO-003"],
          "done_definition": "kiro-gateway 不污染主流程。",
          "risk_notes": "禁止快速 hack。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M7",
      "name": "experimental_proxy 访问控制",
      "goal": "普通用户不能看到或调用 experimental_proxy，internal 用户显式启用后才能调用。",
      "features": [
        {
          "id": "M7-F01-internal-user-detection",
          "milestone": "M7",
          "worker": "routing-security-worker",
          "description": "增加或复用 normal/internal/admin 用户识别。",
          "scope": ["security", "routing"],
          "files_to_inspect": ["user model", "middleware", "auth service"],
          "files_to_modify": ["user permission helpers", "tests"],
          "preconditions": ["M6 完成"],
          "implementation_steps": ["检查现有角色系统", "复用或增加 user_type", "提供权限判断函数"],
          "verification_steps": ["normal 用户不具备 internal 权限", "internal 用户可被识别", "admin 不受影响"],
          "fulfills": ["VAL-EXP-007"],
          "done_definition": "路由层可识别 internal 用户。",
          "risk_notes": "不要重写认证系统。",
          "status": "pending"
        },
        {
          "id": "M7-F02-hide-experimental-from-normal-user",
          "milestone": "M7",
          "worker": "routing-security-worker",
          "description": "普通用户默认不能看到 experimental_proxy。",
          "scope": ["visibility", "security"],
          "files_to_inspect": ["model list", "provider list"],
          "files_to_modify": ["provider/model list filter", "tests"],
          "preconditions": ["M7-F01 完成"],
          "implementation_steps": ["在 Provider/模型列表过滤 experimental_proxy", "后台或 internal 例外"],
          "verification_steps": ["普通用户模型列表不显示 experimental_proxy", "管理员可见"],
          "fulfills": ["VAL-EXP-005"],
          "done_definition": "experimental_proxy 不对普通用户可见。",
          "risk_notes": "不要隐藏管理员必要审计信息。",
          "status": "pending"
        },
        {
          "id": "M7-F03-block-normal-user-experimental-call",
          "milestone": "M7",
          "worker": "routing-security-worker",
          "description": "普通用户调用 experimental_proxy 必须被拒绝。",
          "scope": ["routing", "security"],
          "files_to_inspect": ["routing candidate builder"],
          "files_to_modify": ["routing logic", "tests"],
          "preconditions": ["M7-F02 完成"],
          "implementation_steps": ["在路由候选过滤阶段加入 provider_type 检查", "normal 用户过滤 experimental_proxy"],
          "verification_steps": ["normal 用户调用 experimental_proxy 返回 403 或明确错误"],
          "fulfills": ["VAL-EXP-006"],
          "done_definition": "普通用户不能使用 experimental_proxy。",
          "risk_notes": "错误信息不要泄露内部 Provider 细节。",
          "status": "pending"
        },
        {
          "id": "M7-F04-allow-internal-user-explicit-experimental-call",
          "milestone": "M7",
          "worker": "routing-security-worker",
          "description": "internal 用户在 explicit enable 后可以调用 experimental_proxy。",
          "scope": ["routing", "security"],
          "files_to_inspect": ["routing logic", "provider policy"],
          "files_to_modify": ["routing logic", "tests"],
          "preconditions": ["M7-F03 完成"],
          "implementation_steps": ["检查 enabled", "检查 allow_experimental", "检查 internal/admin 权限"],
          "verification_steps": ["internal 用户显式开启后可调用", "disabled 时任何用户都不能调用"],
          "fulfills": ["VAL-EXP-007", "VAL-EXP-008"],
          "done_definition": "experimental_proxy 的访问控制闭环完成。",
          "risk_notes": "enabled=false 必须最高优先级拒绝。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M8",
      "name": "experimental_proxy 路由与 fallback 隔离",
      "goal": "experimental_proxy 不参与默认 fallback，所有使用必须可审计。",
      "features": [
        {
          "id": "M8-F01-exclude-experimental-when-not-allowed",
          "milestone": "M8",
          "worker": "routing-security-worker",
          "description": "allow_experimental=false 时排除 experimental_proxy。",
          "scope": ["routing"],
          "files_to_inspect": ["route policy", "candidate builder"],
          "files_to_modify": ["routing logic", "tests"],
          "preconditions": ["M7 完成"],
          "implementation_steps": ["检查 route_policy", "构造候选通道时排除 experimental_proxy"],
          "verification_steps": ["allow_experimental=false 时 experimental_proxy 不在候选列表"],
          "fulfills": ["VAL-ROUTE-004"],
          "done_definition": "默认路由不会进入 experimental_proxy。",
          "risk_notes": "不要改变官方 Provider 默认路由。",
          "status": "pending"
        },
        {
          "id": "M8-F02-no-official-to-experimental-fallback",
          "milestone": "M8",
          "worker": "routing-security-worker",
          "description": "official_cloud 失败不能自动 fallback 到 experimental_proxy。",
          "scope": ["routing", "fallback"],
          "files_to_inspect": ["fallback logic", "retry logic"],
          "files_to_modify": ["fallback logic", "tests"],
          "preconditions": ["M8-F01 完成"],
          "implementation_steps": ["检查 fallback 逻辑", "过滤 experimental_proxy", "添加测试"],
          "verification_steps": ["official_cloud 失败后不会自动走 experimental_proxy"],
          "fulfills": ["VAL-ROUTE-003"],
          "done_definition": "experimental_proxy 不参与默认故障切换。",
          "risk_notes": "禁止 silent fallback 到高风险 provider。",
          "status": "pending"
        },
        {
          "id": "M8-F03-explicit-experimental-candidate",
          "milestone": "M8",
          "worker": "routing-security-worker",
          "description": "只有 allow_experimental=true 且具备权限时，experimental_proxy 才能进入候选。",
          "scope": ["routing", "security"],
          "files_to_inspect": ["route policy", "permission helper"],
          "files_to_modify": ["routing logic", "tests"],
          "preconditions": ["M8-F02 完成"],
          "implementation_steps": ["检查 route_policy.allow_experimental", "检查用户或项目权限"],
          "verification_steps": ["allow_experimental=true 且 internal 用户时才进入候选"],
          "fulfills": ["VAL-ROUTE-005"],
          "done_definition": "experimental_proxy 候选路径显式可控。",
          "risk_notes": "不要通过模型名绕过权限。",
          "status": "pending"
        },
        {
          "id": "M8-F04-experimental-log-tags",
          "milestone": "M8",
          "worker": "billing-worker",
          "description": "experimental_proxy 请求记录 provider_type、channel_id、provider_account_id。",
          "scope": ["logging", "audit"],
          "files_to_inspect": ["usage log writer"],
          "files_to_modify": ["usage log writer", "tests"],
          "preconditions": ["M8-F03 完成"],
          "implementation_steps": ["扩展 usage_log 写入字段", "确保 experimental 请求写入 provider_type/channel_id/provider_account_id"],
          "verification_steps": ["internal 用户调用 experimental_proxy 后日志字段完整"],
          "fulfills": ["VAL-USAGE-004"],
          "done_definition": "experimental_proxy 使用可审计。",
          "risk_notes": "不要记录完整 prompt/response。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M9",
      "name": "组织 / 项目基础表结构",
      "goal": "预留 B 端 organization/project，但不做复杂企业权限。",
      "features": [
        {
          "id": "M9-F01-organizations-table",
          "milestone": "M9",
          "worker": "db-migration-worker",
          "description": "新增 organizations 表。",
          "scope": ["database"],
          "files_to_inspect": ["model", "migration"],
          "files_to_modify": ["model", "migration", "tests"],
          "preconditions": ["M8 完成"],
          "implementation_steps": ["新增 organizations model", "字段包含 id、name、owner_user_id、plan_type、contract_type、status、created_at", "创建 migration"],
          "verification_steps": ["organization 可创建", "owner_user_id 可查询"],
          "fulfills": ["VAL-ORG-001"],
          "done_definition": "组织结构可被后续引用。",
          "risk_notes": "不做复杂企业权限。",
          "status": "pending"
        },
        {
          "id": "M9-F02-organization-members-table",
          "milestone": "M9",
          "worker": "db-migration-worker",
          "description": "新增 organization_members 表。",
          "scope": ["database"],
          "files_to_inspect": ["model", "migration"],
          "files_to_modify": ["model", "migration", "tests"],
          "preconditions": ["M9-F01 完成"],
          "implementation_steps": ["新增 organization_members", "字段包含 organization_id、user_id、role、status"],
          "verification_steps": ["用户可加入 organization", "role 可查询"],
          "fulfills": ["VAL-ORG-002", "VAL-ORG-005"],
          "done_definition": "组织成员关系已预留。",
          "risk_notes": "不改变现有用户登录。",
          "status": "pending"
        },
        {
          "id": "M9-F03-projects-table",
          "milestone": "M9",
          "worker": "db-migration-worker",
          "description": "新增 projects 表。",
          "scope": ["database"],
          "files_to_inspect": ["model", "migration"],
          "files_to_modify": ["model", "migration", "tests"],
          "preconditions": ["M9-F02 完成"],
          "implementation_steps": ["新增 projects model", "字段包含 organization_id、name、status、monthly_budget、daily_budget", "创建 migration"],
          "verification_steps": ["organization 下可创建 project", "project status 可用"],
          "fulfills": ["VAL-ORG-003", "VAL-ORG-004"],
          "done_definition": "项目结构可被 API Key 和 usage_log 引用。",
          "risk_notes": "预算字段先预留，不做复杂预算系统。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M10",
      "name": "API Key 组织项目绑定与安全",
      "goal": "API Key 支持 user/org/project 归属、安全存储、禁用、模型和 Provider 类型限制。",
      "features": [
        {
          "id": "M10-F01-api-key-org-project-binding",
          "milestone": "M10",
          "worker": "db-migration-worker",
          "description": "API Key 绑定 user_id、organization_id、project_id。",
          "scope": ["database", "api-key"],
          "files_to_inspect": ["token/API Key model", "auth middleware"],
          "files_to_modify": ["token/API Key model", "migration", "tests"],
          "preconditions": ["M9 完成"],
          "implementation_steps": ["找到现有 token/API Key 表", "增加 organization_id", "增加 project_id", "保持 user_id 兼容"],
          "verification_steps": ["API Key 可绑定 user/org/project", "旧 API Key 不失效"],
          "fulfills": ["VAL-KEY-001", "VAL-KEY-002", "VAL-KEY-003"],
          "done_definition": "API Key 具备 B 端项目归属能力。",
          "risk_notes": "不要破坏旧 token 鉴权。",
          "status": "pending"
        },
        {
          "id": "M10-F02-api-key-hash-prefix-once",
          "milestone": "M10",
          "worker": "routing-security-worker",
          "description": "API Key 只保存 key_hash/key_prefix，完整明文只展示一次。",
          "scope": ["security", "api-key"],
          "files_to_inspect": ["token creation", "token auth"],
          "files_to_modify": ["token creation", "token auth", "tests"],
          "preconditions": ["M10-F01 完成"],
          "implementation_steps": ["检查现有 token 存储方式", "保存 key_hash/key_prefix", "创建时返回完整 key 一次"],
          "verification_steps": ["数据库不保存完整 key", "创建后只展示一次完整 key"],
          "fulfills": ["VAL-KEY-004", "VAL-KEY-005"],
          "done_definition": "API Key 符合安全存储要求。",
          "risk_notes": "迁移旧 key 需要兼容策略。",
          "status": "pending"
        },
        {
          "id": "M10-F03-api-key-disable",
          "milestone": "M10",
          "worker": "routing-security-worker",
          "description": "API Key 支持禁用。",
          "scope": ["security", "api-key"],
          "files_to_inspect": ["token auth", "token model"],
          "files_to_modify": ["token model", "token auth", "tests"],
          "preconditions": ["M10-F02 完成"],
          "implementation_steps": ["增加 status 字段或复用现有状态", "鉴权时检查 status"],
          "verification_steps": ["禁用 key 后调用失败"],
          "fulfills": ["VAL-KEY-006"],
          "done_definition": "API Key 可被禁用。",
          "risk_notes": "禁用错误需明确且不泄露敏感信息。",
          "status": "pending"
        },
        {
          "id": "M10-F04-api-key-model-provider-limits",
          "milestone": "M10",
          "worker": "routing-security-worker",
          "description": "API Key 支持 allowed_models 和 allowed_provider_types。",
          "scope": ["security", "routing"],
          "files_to_inspect": ["token model", "routing logic"],
          "files_to_modify": ["token model", "routing/auth checks", "tests"],
          "preconditions": ["M10-F03 完成"],
          "implementation_steps": ["增加 allowed_models", "增加 allowed_provider_types", "在鉴权或路由阶段校验", "普通 API Key 默认不允许 experimental_proxy"],
          "verification_steps": ["未授权模型不可调用", "普通 API Key 默认不能访问 experimental_proxy"],
          "fulfills": ["VAL-KEY-007", "VAL-KEY-008", "VAL-KEY-009"],
          "done_definition": "API Key 可控制模型和 Provider 类型访问范围。",
          "risk_notes": "allowed_provider_types 不能覆盖 experimental_proxy 安全边界。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M11",
      "name": "usage_log 与隐私策略",
      "goal": "请求可审计，默认不保存敏感正文。",
      "features": [
        {
          "id": "M11-F01-usage-log-schema",
          "milestone": "M11",
          "worker": "db-migration-worker",
          "description": "扩展 usage_logs 字段。",
          "scope": ["database", "logging"],
          "files_to_inspect": ["log model", "quota log model", "usage log writer"],
          "files_to_modify": ["model", "migration", "tests"],
          "preconditions": ["M10 完成"],
          "implementation_steps": ["增加 request_id、user_id、organization_id、project_id、api_key_id", "增加 provider_id、provider_type、channel_id、provider_account_id", "增加 tokens、cost、status、latency、ip、user_agent"],
          "verification_steps": ["migration 可执行", "字段可读写"],
          "fulfills": ["VAL-USAGE-003", "VAL-USAGE-004", "VAL-USAGE-005", "VAL-USAGE-006"],
          "done_definition": "usage_log 足够支撑账单、审计、统计。",
          "risk_notes": "跨库字段类型保持兼容。",
          "status": "pending"
        },
        {
          "id": "M11-F02-write-success-and-failed-logs",
          "milestone": "M11",
          "worker": "billing-worker",
          "description": "成功和失败请求都写入日志。",
          "scope": ["logging"],
          "files_to_inspect": ["usage log writer", "request handling"],
          "files_to_modify": ["usage log writer", "tests"],
          "preconditions": ["M11-F01 完成"],
          "implementation_steps": ["成功请求写 usage_log", "失败请求写 usage_log 或 error log"],
          "verification_steps": ["成功请求有日志", "失败请求有日志"],
          "fulfills": ["VAL-USAGE-001", "VAL-USAGE-002"],
          "done_definition": "请求结果均可审计。",
          "risk_notes": "失败日志也不能包含完整 prompt/response。",
          "status": "pending"
        },
        {
          "id": "M11-F03-no-prompt-response-logging",
          "milestone": "M11",
          "worker": "routing-security-worker",
          "description": "默认不保存完整 prompt/response，并脱敏 error_message。",
          "scope": ["privacy", "logging"],
          "files_to_inspect": ["log writer", "error handling"],
          "files_to_modify": ["log writer", "error sanitization", "tests"],
          "preconditions": ["M11-F02 完成"],
          "implementation_steps": ["检查现有日志逻辑", "关闭 prompt/response 持久化", "error_message 脱敏", "debug_payload_logs 默认关闭"],
          "verification_steps": ["usage_log 不包含完整 prompt", "usage_log 不包含完整 response", "error_message 不泄露 key"],
          "fulfills": ["VAL-PRIVACY-001", "VAL-PRIVACY-002", "VAL-PRIVACY-003", "VAL-PRIVACY-004"],
          "done_definition": "默认日志满足隐私要求。",
          "risk_notes": "不要在测试 fixture 中引入真实密钥。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M12",
      "name": "token 统计与基础扣费",
      "goal": "形成 token → cost → balance deduction 闭环。",
      "features": [
        {
          "id": "M12-F01-token-usage-extraction",
          "milestone": "M12",
          "worker": "billing-worker",
          "description": "提取 input/output/total tokens。",
          "scope": ["billing", "usage"],
          "files_to_inspect": ["OpenAI-compatible response handling", "usage log writer"],
          "files_to_modify": ["usage extraction", "tests"],
          "preconditions": ["M11 完成"],
          "implementation_steps": ["识别 OpenAI-compatible usage 字段", "提取 prompt_tokens", "提取 completion_tokens", "计算 total_tokens"],
          "verification_steps": ["成功请求后 usage_log 有 input/output/total tokens"],
          "fulfills": ["VAL-BILL-001", "VAL-BILL-002", "VAL-BILL-003"],
          "done_definition": "token 统计可用于扣费。",
          "risk_notes": "缺失 usage 时需要明确 fallback 策略。",
          "status": "pending"
        },
        {
          "id": "M12-F02-cost-calculation",
          "milestone": "M12",
          "worker": "billing-worker",
          "description": "根据模型价格和 token usage 计算 cost。",
          "scope": ["billing"],
          "files_to_inspect": ["model mapping", "billing service"],
          "files_to_modify": ["billing service", "tests"],
          "preconditions": ["M12-F01 完成"],
          "implementation_steps": ["读取 input_price/output_price/billing_multiplier", "计算 cost"],
          "verification_steps": ["cost 等于价格公式"],
          "fulfills": ["VAL-BILL-004"],
          "done_definition": "费用计算可复现。",
          "risk_notes": "避免浮点误差。",
          "status": "pending"
        },
        {
          "id": "M12-F03-balance-deduction",
          "milestone": "M12",
          "worker": "billing-worker",
          "description": "成功请求后扣余额，失败请求不错误扣费。",
          "scope": ["billing", "balance"],
          "files_to_inspect": ["quota service", "billing service"],
          "files_to_modify": ["billing service", "balance service", "tests"],
          "preconditions": ["M12-F02 完成"],
          "implementation_steps": ["成功请求扣余额", "失败请求不扣费或按规则处理", "写 charged_amount"],
          "verification_steps": ["成功请求后余额减少", "失败请求不错误扣费"],
          "fulfills": ["VAL-BILL-005", "VAL-BILL-006"],
          "done_definition": "基础按量计费闭环完成。",
          "risk_notes": "扣费需要考虑并发一致性。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M13",
      "name": "余额不足与管理员手动充值",
      "goal": "防止余额不足继续调用上游，MVP 支持管理员手动充值。",
      "features": [
        {
          "id": "M13-F01-insufficient-balance-rejection",
          "milestone": "M13",
          "worker": "billing-worker",
          "description": "余额不足时拒绝请求且不调用上游。",
          "scope": ["billing", "routing"],
          "files_to_inspect": ["billing precheck", "relay request flow"],
          "files_to_modify": ["billing precheck", "tests"],
          "preconditions": ["M12 完成"],
          "implementation_steps": ["上游请求前检查余额", "余额不足返回明确错误", "写日志", "不发起上游请求"],
          "verification_steps": ["余额为 0 的用户调用付费模型被拒绝", "无上游调用", "日志记录 insufficient_balance"],
          "fulfills": ["VAL-BILL-007", "VAL-BILL-008", "VAL-BILL-009"],
          "done_definition": "余额不足不会产生亏损调用。",
          "risk_notes": "余额预检不能绕过模型权限检查。",
          "status": "pending"
        },
        {
          "id": "M13-F02-admin-manual-recharge",
          "milestone": "M13",
          "worker": "billing-worker",
          "description": "管理员可以手动增加用户余额，并记录操作。",
          "scope": ["billing", "admin-api"],
          "files_to_inspect": ["admin controller", "quota service", "operation log"],
          "files_to_modify": ["admin recharge API", "balance service", "operation log", "tests"],
          "preconditions": ["M13-F01 完成"],
          "implementation_steps": ["增加或复用管理员充值接口", "更新用户余额", "写操作记录"],
          "verification_steps": ["管理员可加余额", "余额变化可查", "操作有记录"],
          "fulfills": ["VAL-BILL-010", "VAL-BILL-011"],
          "done_definition": "MVP 可通过手动充值运营。",
          "risk_notes": "只允许管理员调用。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M14",
      "name": "管理后台 Provider / Channel 最小页面",
      "goal": "管理员能查看 Provider/Channel 风险和关联关系，并一键关闭风险 Provider。",
      "features": [
        {
          "id": "M14-F01-provider-admin-page",
          "milestone": "M14",
          "worker": "admin-ui-worker",
          "description": "Provider 列表显示 provider_type、risk_level、available_scope、enabled。",
          "scope": ["admin-ui"],
          "files_to_inspect": ["web/default", "provider/channel admin APIs"],
          "files_to_modify": ["admin provider page", "i18n files", "tests if present"],
          "preconditions": ["M13 完成"],
          "implementation_steps": ["找到后台 Provider/Channel 页面", "增加 provider_type/risk_level/available_scope/enabled 显示"],
          "verification_steps": ["管理员能看到 Provider 列表和关键字段"],
          "fulfills": ["VAL-ADMIN-001", "VAL-ADMIN-002"],
          "done_definition": "管理员能识别正规 Provider 和 experimental_proxy。",
          "risk_notes": "前端新增文案需要 i18n。",
          "status": "pending"
        },
        {
          "id": "M14-F02-provider-enable-disable",
          "milestone": "M14",
          "worker": "admin-ui-worker",
          "description": "管理员可以启用/禁用 Provider。",
          "scope": ["admin-ui", "admin-api"],
          "files_to_inspect": ["provider admin page", "provider admin API"],
          "files_to_modify": ["provider admin page", "provider admin API if needed", "tests"],
          "preconditions": ["M14-F01 完成"],
          "implementation_steps": ["增加 enabled 开关", "路由层读取 enabled 状态"],
          "verification_steps": ["禁用 Provider 后不可被路由", "experimental_proxy 可一键关闭"],
          "fulfills": ["VAL-ADMIN-003", "VAL-EXP-009"],
          "done_definition": "风险 Provider 可以快速下线。",
          "risk_notes": "禁用状态必须影响路由。",
          "status": "pending"
        },
        {
          "id": "M14-F03-channel-provider-admin",
          "milestone": "M14",
          "worker": "admin-ui-worker",
          "description": "Channel 页面显示 Provider / Provider Account 关联。",
          "scope": ["admin-ui"],
          "files_to_inspect": ["channel admin page", "channel admin API"],
          "files_to_modify": ["channel admin page", "i18n files", "tests if present"],
          "preconditions": ["M14-F02 完成"],
          "implementation_steps": ["Channel 列表显示 provider_id、provider_account_id、provider_type"],
          "verification_steps": ["管理员可以看到 Channel 关联的 Provider 和 Account"],
          "fulfills": ["VAL-ADMIN-004"],
          "done_definition": "Channel 与 Provider 的关系可被后台管理。",
          "risk_notes": "不要在页面显示明文凭证。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M15",
      "name": "管理后台 API Key / usage_log / 余额页面",
      "goal": "管理员可以完成 MVP 运营闭环。",
      "features": [
        {
          "id": "M15-F01-api-key-admin-page",
          "milestone": "M15",
          "worker": "admin-ui-worker",
          "description": "管理员查看并禁用 API Key。",
          "scope": ["admin-ui", "api-key"],
          "files_to_inspect": ["API Key admin APIs", "web/default"],
          "files_to_modify": ["API Key admin page", "i18n files", "tests if present"],
          "preconditions": ["M14 完成"],
          "implementation_steps": ["API Key 页面显示 key_prefix/status", "支持 disable"],
          "verification_steps": ["管理员能查看用户 API Key", "管理员能禁用 API Key"],
          "fulfills": ["VAL-ADMIN-005", "VAL-ADMIN-006"],
          "done_definition": "API Key 可被后台管控。",
          "risk_notes": "不得显示完整 API Key。",
          "status": "pending"
        },
        {
          "id": "M15-F02-usage-log-admin-page",
          "milestone": "M15",
          "worker": "admin-ui-worker",
          "description": "管理员查看 usage_log。",
          "scope": ["admin-ui", "usage-log"],
          "files_to_inspect": ["usage log APIs", "web/default"],
          "files_to_modify": ["usage log admin page", "i18n files", "tests if present"],
          "preconditions": ["M15-F01 完成"],
          "implementation_steps": ["日志页面增加 user、model、provider_type、tokens、cost、status、latency"],
          "verification_steps": ["管理员可以查看完整 usage_log"],
          "fulfills": ["VAL-ADMIN-007", "VAL-ADMIN-008"],
          "done_definition": "管理员能审计请求和成本。",
          "risk_notes": "日志页面不得显示完整 prompt/response。",
          "status": "pending"
        },
        {
          "id": "M15-F03-balance-admin-page",
          "milestone": "M15",
          "worker": "admin-ui-worker",
          "description": "管理员查看和调整用户余额。",
          "scope": ["admin-ui", "billing"],
          "files_to_inspect": ["balance APIs", "web/default"],
          "files_to_modify": ["balance admin page", "i18n files", "tests if present"],
          "preconditions": ["M15-F02 完成"],
          "implementation_steps": ["余额页面显示用户余额", "支持手动调整"],
          "verification_steps": ["管理员可查看余额", "管理员可调整余额"],
          "fulfills": ["VAL-ADMIN-009", "VAL-ADMIN-010"],
          "done_definition": "管理员能完成手动充值运营。",
          "risk_notes": "余额调整必须有操作记录。",
          "status": "pending"
        }
      ]
    },
    {
      "id": "M16",
      "name": "Docker、文档和回归测试",
      "goal": "MVP 可启动、可调用、可回归验证。",
      "features": [
        {
          "id": "M16-F01-docker-compose-local",
          "milestone": "M16",
          "worker": "deployment-worker",
          "description": "Docker Compose 启动 gateway、db、redis。",
          "scope": ["deployment"],
          "files_to_inspect": ["docker-compose.yml", "Dockerfile", ".env.example"],
          "files_to_modify": ["docker-compose.yml", ".env.example", ".factory/services.yaml"],
          "preconditions": ["M15 完成"],
          "implementation_steps": ["检查 docker-compose.yml", "增加 db/redis/gateway", "增加 healthcheck", "更新 .env.example"],
          "verification_steps": ["docker compose up 成功", "gateway/db/redis healthcheck 通过", ".env.example 无真实密钥"],
          "fulfills": ["VAL-DEPLOY-001", "VAL-DEPLOY-002", "VAL-DEPLOY-003", "VAL-DEPLOY-004"],
          "done_definition": "新环境可以本地一键启动。",
          "risk_notes": "不得提交真实密钥。",
          "status": "pending"
        },
        {
          "id": "M16-F02-provider-policy-docs",
          "milestone": "M16",
          "worker": "deployment-worker",
          "description": "编写 experimental_proxy 策略文档。",
          "scope": ["docs"],
          "files_to_inspect": ["docs", "provider policy"],
          "files_to_modify": ["docs/PROVIDER_POLICY.md"],
          "preconditions": ["M16-F01 完成"],
          "implementation_steps": ["说明 kiro-gateway 是 external experimental provider", "说明默认 internal_only", "说明禁止默认 fallback"],
          "verification_steps": ["文档包含 experimental_proxy 默认约束和风险边界"],
          "fulfills": ["VAL-DOC-001", "VAL-DOC-002", "VAL-DOC-003"],
          "done_definition": "后续维护者不会误把 experimental_proxy 开放给普通用户。",
          "risk_notes": "文档必须和代码策略一致。",
          "status": "pending"
        },
        {
          "id": "M16-F03-openai-sdk-docs",
          "milestone": "M16",
          "worker": "deployment-worker",
          "description": "编写 OpenAI SDK 和 API Key 调用文档。",
          "scope": ["docs"],
          "files_to_inspect": ["docs", "API routes"],
          "files_to_modify": ["docs/OPENAI_COMPATIBLE_API.md"],
          "preconditions": ["M16-F02 完成"],
          "implementation_steps": ["写 curl 示例", "写 OpenAI SDK 示例", "写 API Key 使用方法", "写余额不足错误说明"],
          "verification_steps": ["文档包含 /v1/models 和 /v1/chat/completions 示例"],
          "fulfills": ["VAL-DOC-004", "VAL-DOC-005"],
          "done_definition": "用户可以照文档调用 MVP。",
          "risk_notes": "示例 key 必须是假数据。",
          "status": "pending"
        },
        {
          "id": "M16-F04-regression-tests",
          "milestone": "M16",
          "worker": "test-validation-worker",
          "description": "编写 MVP 回归测试脚本。",
          "scope": ["test", "regression"],
          "files_to_inspect": ["tests", ".factory/services.yaml"],
          "files_to_modify": ["tests", "scripts if needed"],
          "preconditions": ["M16-F03 完成"],
          "implementation_steps": ["普通用户调用 official provider 测试", "internal 用户调用 experimental provider 测试", "普通用户调用 experimental 被拒绝测试", "disabled experimental 被拒绝测试", "余额不足测试", "不保存 prompt/response 检查"],
          "verification_steps": ["所有测试可本地执行且输出清晰"],
          "fulfills": ["VAL-TEST-001", "VAL-TEST-002", "VAL-TEST-003", "VAL-TEST-004", "VAL-TEST-005", "VAL-TEST-006"],
          "done_definition": "MVP 有可重复回归测试。",
          "risk_notes": "测试不得依赖真实上游 key。",
          "status": "pending"
        }
      ]
    }
  ]
}
```

## 5. 每个 Milestone 的验收标准

### M0：仓库调研与 Mission 基础文件

- 完成 `docs/REPO_STRUCTURE.md`。
- 完成 `docs/OPENAI_REQUEST_FLOW.md`。
- Mission 文件存在：`AGENTS.md`、`mission.md`、`validation-contract.md`、`features.json`、`.factory/services.yaml`、`.factory/init.sh`。
- `.factory/services.yaml` 包含本地启动、停止、测试、构建、healthcheck 命令。
- 不修改任何业务逻辑。

### M1：Provider 类型与字段基础迁移

- 支持 `official_cloud`、`aggregator`、`authorized_proxy`、`experimental_proxy`。
- Provider 支持 `risk_level`、`available_scope`、`manual_enable_required`、`enabled`。
- `experimental_proxy` 新建默认 high risk、internal_only、manual enable、disabled。
- 旧 Channel 不因 Provider 字段增加而失效。

### M2：Provider Account 与 Channel 关联

- 一个 Provider 可绑定多个 Provider Account。
- Provider Account 凭证不得明文保存。
- Channel 可关联 Provider 和 Provider Account。
- 旧 Channel 迁移后仍可读取和路由。

### M3：OpenAI-compatible 接口回归

- `GET /v1/models` 返回 OpenAI-compatible 模型列表。
- `POST /v1/chat/completions` 支持 OpenAI SDK 非流式调用。
- 如果原项目支持 stream，stream 兼容性不被破坏。

### M4：正规 Provider 最小接入

- 系统至少存在两个非 experimental Provider。
- 普通用户可调用 `official_cloud` 或 `aggregator`。
- 日志能识别 `provider_type=official_cloud` 或 `provider_type=aggregator`。
- 不使用 experimental_proxy。

### M5：模型映射与价格基础

- 支持 `public_model_name` 到 `provider_model_name` 的映射。
- 支持 `input_price`、`output_price`、`billing_multiplier`。
- disabled model 不可被调用。
- 价格字段避免浮点误差。

### M6：KiroGatewayAdapter 骨架

- `ProviderAdapter` 基础接口存在。
- `KiroGatewayAdapter` 实现基础接口。
- `KiroGatewayAdapter` 标记为 `experimental_proxy`。
- `KiroGatewayAdapter` 默认 disabled。
- 主请求流程不得硬编码 Kiro 分支。

### M7：experimental_proxy 访问控制

- 路由层可识别 normal/internal/admin 用户。
- 普通用户不能看到 experimental_proxy。
- 普通用户调用 experimental_proxy 被拒绝。
- internal 用户只有显式开启后才能调用。
- disabled experimental_proxy 任何用户都不能调用。

### M8：experimental_proxy 路由与 fallback 隔离

- `allow_experimental=false` 时 experimental_proxy 不进入候选通道。
- official_cloud 失败不能默认 fallback 到 experimental_proxy。
- 只有 `allow_experimental=true` 且具备权限时 experimental_proxy 才能进入候选。
- experimental_proxy 请求日志包含 provider/account/channel 标识。

### M9：组织 / 项目基础表结构

- 支持 `organizations` 表。
- 支持 `organization_members` 表。
- 支持 `projects` 表。
- organization 可以拥有 project。
- user 可以归属 organization。

### M10：API Key 组织项目绑定与安全

- API Key 可绑定 user/org/project。
- API Key 只保存 `key_hash` 和 `key_prefix`。
- 完整明文 API Key 只展示一次。
- API Key 可禁用。
- API Key 支持 `allowed_models` 和 `allowed_provider_types`。
- 普通 API Key 默认不允许 experimental_proxy。

### M11：usage_log 与隐私策略

- 成功请求必须写 usage_log。
- 失败请求必须写 usage_log 或 error log。
- usage_log 包含请求、用户、组织、项目、key、provider、channel、tokens、cost、status、latency、ip、user_agent。
- 默认不保存完整 prompt。
- 默认不保存完整 response。
- error_message 必须脱敏。

### M12：token 统计与基础扣费

- 可提取 input/output/total tokens。
- 可根据模型价格和 token usage 计算 cost。
- 成功请求后余额减少。
- 失败请求不错误扣费。
- usage_log 写入 cost 和 charged_amount。

### M13：余额不足与管理员手动充值

- 余额不足时拒绝请求。
- 余额不足时不能发起上游请求。
- 余额不足请求写入日志。
- 管理员可以手动增加用户余额。
- 管理员手动充值有操作记录。

### M14：管理后台 Provider / Channel 最小页面

- 管理员可查看 Provider 列表。
- Provider 列表显示 provider_type、risk_level、available_scope、enabled。
- 管理员可启用/禁用 Provider。
- 管理员可查看 Channel 与 Provider / Provider Account 的关系。
- experimental_proxy 可一键禁用。

### M15：管理后台 API Key / usage_log / 余额页面

- 管理员可查看 API Key 列表。
- 管理员可禁用 API Key。
- 管理员可查看 usage_log。
- usage_log 页面显示 user、model、provider_type、tokens、cost、status、latency。
- 管理员可查看和调整用户余额。

### M16：Docker、文档和回归测试

- docker compose 可启动 gateway、db、redis。
- gateway/db/redis healthcheck 通过。
- `.env.example` 不包含真实密钥。
- 文档说明 kiro-gateway 是 external experimental provider。
- 文档说明 experimental_proxy 默认 internal_only，普通用户不可调用。
- 文档说明 official_cloud 不能默认 fallback 到 experimental_proxy。
- 文档包含 OpenAI SDK 和 API Key 使用示例。
- 回归测试覆盖 official provider、experimental provider、普通用户拒绝、disabled experimental 拒绝、余额不足拒绝、默认不保存 prompt/response。

## 6. Worker 职责说明

### repo-discovery-worker

负责只读调研仓库，输出模块位置、请求链路、二开文件清单和 Mission 基础文件。不得修改业务逻辑。

### db-migration-worker

负责数据库模型、GORM migration、跨 SQLite/MySQL/PostgreSQL 兼容性、旧数据兼容和 schema 级测试。不得写数据库专属 SQL，除非提供跨库 fallback。

### provider-worker

负责 Provider、Provider Account、Channel 关系、ProviderAdapter 抽象和 provider adapter skeleton。不得处理 billing、admin UI 或复杂路由策略。

### routing-security-worker

负责路由候选过滤、experimental_proxy 权限隔离、API Key 安全、模型禁用检查和隐私策略。不得实现扣费 UI 或部署。

### openai-api-worker

负责 OpenAI-compatible API 回归测试，覆盖 `/v1/models`、`/v1/chat/completions`、stream 兼容性。不得改协议格式。

### billing-worker

负责 usage_log 写入、token 提取、cost 计算、余额扣费、余额不足拒绝和管理员手动充值接口。不得实现复杂支付或返佣。

### admin-ui-worker

负责管理后台最小页面：Provider/Channel、API Key、usage_log、余额。前端新增文案必须走 i18n，不得展示明文 key 或完整 prompt/response。

### deployment-worker

负责 Docker Compose、本地环境、`.env.example`、部署文档、Provider 策略文档和 OpenAI SDK 示例。不得提交真实密钥。

### test-validation-worker

负责 validation-contract 对应的测试、架构检查和回归脚本。发现 feature 未满足契约时必须阻断，不允许继续执行后续 feature。

## 7. 推荐执行顺序

1. 第 1 轮：只让 Droid 生成 `mission.md`、`validation-contract.md`、`features.json`，不写代码。
2. 第 2 轮：执行 M0。
3. 第 3 轮：执行 M1。
4. 第 4 轮：执行 M2。
5. 第 5 轮：执行 M3。
6. 第 6 轮：执行 M4。
7. 第 7 轮：执行 M5。
8. 第 8 轮：执行 M6。
9. 第 9 轮：执行 M7。
10. 第 10 轮：执行 M8。
11. 第 11 轮：执行 M9。
12. 第 12 轮：执行 M10。
13. 第 13 轮：执行 M11。
14. 第 14 轮：执行 M12。
15. 第 15 轮：执行 M13。
16. 第 16 轮：执行 M14。
17. 第 17 轮：执行 M15。
18. 第 18 轮：执行 M16。

关键边界：

```text
先 Provider 元数据
再 Provider Account / Channel
再 OpenAI-compatible 回归
再正规 Provider
再模型映射和价格
再 Kiro adapter
再 experimental_proxy 权限隔离
再 fallback 隔离
再组织 / 项目 / API Key
再 usage_log
再 billing
再后台
最后部署和测试
```

## 8. 每个 Milestone 的启动提示词

### 启动 M0

```text
确认 Mission 规划。现在开始执行 M0：仓库调研与 Mission 基础文件。

范围限制：
1. 只读仓库。
2. 只创建文档和 Mission 文件。
3. 不修改业务逻辑。
4. 不新增数据库字段。
5. 不改变 API 行为。

需要完成：
- M0-F01-repo-structure-scan
- M0-F02-openai-request-flow
- M0-F03-create-mission-files

完成后返回结构化 handoff。
```

### 启动 M1

```text
现在开始执行 M1：Provider 类型与字段基础迁移。

范围限制：
1. 只做 Provider 类型、风险字段和 experimental_proxy 默认值。
2. 不接入真实 Kiro 流量。
3. 不改变原有 OpenAI-compatible API 行为。
4. 必须兼容原有 channel。

需要完成：
- M1-F01-provider-type-enum
- M1-F02-provider-risk-scope-fields
- M1-F03-experimental-default-policy

完成后运行 migration、lint、test，并返回结构化 handoff。
```

### 启动 M2

```text
现在开始执行 M2：Provider Account 与 Channel 关联。

范围限制：
1. 只做 Provider、Provider Account、Channel 三者关系。
2. 不做路由策略。
3. 不接入 Kiro。
4. 必须保证旧 Channel 不失效。

需要完成：
- M2-F01-provider-accounts-table
- M2-F02-channel-provider-link
- M2-F03-legacy-channel-compatibility

完成后返回结构化 handoff。
```

### 启动 M3

```text
现在开始执行 M3：OpenAI-compatible 接口回归。

范围限制：
1. 只验证原有 OpenAI-compatible API。
2. 不新增 Provider。
3. 不修改计费逻辑。
4. 不修改 experimental_proxy。

需要完成：
- M3-F01-verify-models-endpoint
- M3-F02-verify-chat-completions
- M3-F03-verify-chat-streaming-if-supported

完成后必须证明 /v1/models 和 /v1/chat/completions 没有被破坏。
```

### 启动 M4

```text
现在开始执行 M4：正规 Provider 最小接入。

范围限制：
1. 只处理 official_cloud / aggregator。
2. 不处理 experimental_proxy。
3. 不接入 Kiro。
4. 不做复杂 billing ledger。

需要完成：
- M4-F01-first-official-provider
- M4-F02-second-normal-provider

完成后必须证明普通用户可以调用至少一个 official_cloud Provider，系统至少存在两个非 experimental Provider。
```

### 启动 M5

```text
现在开始执行 M5：模型映射与价格基础。

范围限制：
1. 只做 model mapping 和价格字段。
2. 不做真正扣费。
3. 不做前端复杂页面。

需要完成：
- M5-F01-model-mapping-table
- M5-F02-model-price-fields
- M5-F03-disabled-model-rejection

完成后必须证明 public_model_name 可以映射到 provider_model_name，disabled model 不可调用。
```

### 启动 M6

```text
现在开始执行 M6：KiroGatewayAdapter 骨架。

范围限制：
1. Kiro 只能作为 experimental_proxy adapter。
2. 不允许写死 Kiro 到主流程。
3. 不允许默认开放 Kiro 流量。
4. Kiro adapter 默认 disabled。

需要完成：
- M6-F01-provider-adapter-interface
- M6-F02-kiro-gateway-adapter-skeleton
- M6-F03-no-hardcoded-kiro-flow

完成后必须证明 KiroGatewayAdapter 只是一个 adapter，不污染主流程。
```

### 启动 M7

```text
现在开始执行 M7：experimental_proxy 访问控制。

范围限制：
1. 只做可见性和访问控制。
2. 不做 fallback。
3. 不做 billing。
4. 不做后台页面。

需要完成：
- M7-F01-internal-user-detection
- M7-F02-hide-experimental-from-normal-user
- M7-F03-block-normal-user-experimental-call
- M7-F04-allow-internal-user-explicit-experimental-call

完成后必须证明普通用户不能看到或调用 experimental_proxy，internal 用户显式启用后才能调用，disabled experimental_proxy 任何用户都不能调用。
```

### 启动 M8

```text
现在开始执行 M8：experimental_proxy 路由与 fallback 隔离。

范围限制：
1. 只做路由候选集和 fallback 隔离。
2. 不做用户角色系统。
3. 不做后台页面。
4. 不做计费扣费。

需要完成：
- M8-F01-exclude-experimental-when-not-allowed
- M8-F02-no-official-to-experimental-fallback
- M8-F03-explicit-experimental-candidate
- M8-F04-experimental-log-tags

完成后必须证明 official_cloud 失败不会默认 fallback 到 experimental_proxy，allow_experimental=false 时 experimental_proxy 不会进入候选通道。
```

### 启动 M9

```text
现在开始执行 M9：组织 / 项目基础表结构。

范围限制：
1. 只做 organizations、organization_members、projects 表。
2. 不做完整企业权限。
3. 不做企业后台。
4. 不做合同计费。

需要完成：
- M9-F01-organizations-table
- M9-F02-organization-members-table
- M9-F03-projects-table

完成后必须证明 organization 可以拥有 project，user 可以归属 organization。
```

### 启动 M10

```text
现在开始执行 M10：API Key 组织项目绑定与安全。

范围限制：
1. 只做 API Key 归属、安全存储、禁用和访问限制。
2. 不做复杂企业权限。
3. 不做支付。
4. 不做代理。

需要完成：
- M10-F01-api-key-org-project-binding
- M10-F02-api-key-hash-prefix-once
- M10-F03-api-key-disable
- M10-F04-api-key-model-provider-limits

完成后必须证明 API Key 可以绑定 user/org/project，只保存 key_hash/key_prefix，普通 API Key 默认不能访问 experimental_proxy。
```

### 启动 M11

```text
现在开始执行 M11：usage_log 与隐私策略。

范围限制：
1. 只做 usage_log 字段、成功/失败日志和隐私策略。
2. 不做扣费。
3. 不做余额不足。
4. 不做后台 UI。

需要完成：
- M11-F01-usage-log-schema
- M11-F02-write-success-and-failed-logs
- M11-F03-no-prompt-response-logging

完成后必须证明成功和失败请求都有日志，默认不保存完整 prompt/response，error_message 已脱敏。
```

### 启动 M12

```text
现在开始执行 M12：token 统计与基础扣费。

范围限制：
1. 只做 token 提取、cost 计算、成功扣余额。
2. 不做余额不足拒绝。
3. 不做管理员充值。
4. 不做支付系统。

需要完成：
- M12-F01-token-usage-extraction
- M12-F02-cost-calculation
- M12-F03-balance-deduction

完成后必须证明成功请求后 usage_log 有 token 和 cost，用户余额减少，失败请求不错误扣费。
```

### 启动 M13

```text
现在开始执行 M13：余额不足与管理员手动充值。

范围限制：
1. 只做余额不足拒绝和管理员手动充值。
2. 不做 Stripe。
3. 不做支付宝。
4. 不做微信支付。
5. 不做自动返佣。

需要完成：
- M13-F01-insufficient-balance-rejection
- M13-F02-admin-manual-recharge

完成后必须证明余额不足不调用上游，管理员可以手动给用户充值，并且有操作记录。
```

### 启动 M14

```text
现在开始执行 M14：管理后台 Provider / Channel 最小页面。

范围限制：
1. 只做 Provider 和 Channel 管理页面的最小改造。
2. 不做复杂统计图。
3. 不做代理后台。
4. 不做企业合同后台。

需要完成：
- M14-F01-provider-admin-page
- M14-F02-provider-enable-disable
- M14-F03-channel-provider-admin

完成后必须证明管理员能查看 Provider 类型、风险等级、可见范围，并能一键禁用 experimental_proxy。
```

### 启动 M15

```text
现在开始执行 M15：管理后台 API Key / usage_log / 余额页面。

范围限制：
1. 只做 API Key、usage_log、余额管理页面。
2. 不做完整企业后台。
3. 不做代理后台。
4. 不做复杂图表。

需要完成：
- M15-F01-api-key-admin-page
- M15-F02-usage-log-admin-page
- M15-F03-balance-admin-page

完成后必须证明管理员可以查看/禁用 API Key，查看 usage_log，查看和调整用户余额。
```

### 启动 M16

```text
现在开始执行 M16：Docker、文档和回归测试。

范围限制：
1. 只做本地部署、文档、测试脚本。
2. 不做生产 Kubernetes。
3. 不做复杂监控系统。

需要完成：
- M16-F01-docker-compose-local
- M16-F02-provider-policy-docs
- M16-F03-openai-sdk-docs
- M16-F04-regression-tests

完成后必须证明 docker compose 可启动，OpenAI SDK 示例可用，回归测试覆盖 official provider、experimental provider、普通用户拒绝、disabled experimental 拒绝、余额不足拒绝、默认不保存 prompt/response。
```

## 9. Feature Handoff 模板

每完成一个 feature，Droid 必须按以下格式返回：

```json
{
  "feature_id": "",
  "milestone": "",
  "status": "completed | blocked | failed",
  "files_inspected": [],
  "files_modified": [],
  "implementation_summary": "",
  "validation_assertions_fulfilled": [],
  "commands_run": [
    {
      "command": "",
      "exit_code": 0,
      "observation": ""
    }
  ],
  "tests_added": [],
  "manual_checks": [],
  "risks_or_todos": [],
  "breaking_changes": [],
  "next_recommended_feature": ""
}
```

如果某个 feature 失败，不要继续往下跑，必须返回：

```json
{
  "feature_id": "",
  "status": "blocked",
  "blocker": "",
  "why_it_blocks_the_mission": "",
  "minimal_fix_path": "",
  "files_related": [],
  "recommended_next_action": ""
}
```

## 10. 最终执行顺序

第 1 轮：只让 Droid 生成 mission.md、validation-contract.md、features.json，不写代码。
第 2 轮：执行 M0。
第 3 轮：执行 M1。
第 4 轮：执行 M2。
第 5 轮：执行 M3。
第 6 轮：执行 M4。
第 7 轮：执行 M5。
第 8 轮：执行 M6。
第 9 轮：执行 M7。
第 10 轮：执行 M8。
第 11 轮：执行 M9。
第 12 轮：执行 M10。
第 13 轮：执行 M11。
第 14 轮：执行 M12。
第 15 轮：执行 M13。
第 16 轮：执行 M14。
第 17 轮：执行 M15。
第 18 轮：执行 M16。

最关键的边界：

```text
先 Provider 元数据
再 Provider Account / Channel
再 OpenAI-compatible 回归
再正规 Provider
再模型映射和价格
再 Kiro adapter
再 experimental_proxy 权限隔离
再 fallback 隔离
再组织 / 项目 / API Key
再 usage_log
再 billing
再后台
最后部署和测试
```

这样 Droid Mission 会更稳，不容易在一个 milestone 里同时改数据库、路由、计费、后台和部署。
