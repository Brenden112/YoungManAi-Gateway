# Mission: B2B Multi-Provider AI API Gateway MVP

## Overview

二次开发 new-api，构建面向 B 端的 Multi-Provider AI API Gateway MVP。
第一阶段仅暴露 OpenAI-compatible API，后续扩展 Claude / Gemini / Responses API。

## Goals

1. 在 new-api 基础上叠加 Provider Adapter Framework，不破坏原有功能。
2. 支持四种 Provider 类型：official_cloud / aggregator / authorized_proxy / experimental_proxy。
3. experimental_proxy 默认高风险、仅内部可见、需手动启用。
4. 跑通完整 MVP 闭环：API Key → 路由 → 上游请求 → token 统计 → 余额扣费 → usage_log。
5. 管理员可手动充值、可一键禁用 experimental_proxy。
6. 普通用户不可见、不可调用 experimental_proxy。

## Out of Scope (v1)

- 自动支付 / 自动返佣
- 完整企业合同 / 多协议
- 复杂代理后台
- Claude-compatible / Gemini-compatible / Responses API（预留，不实现）
- 完整 Organization / Project 权限体系（仅建表，不实现业务逻辑）

## Execution Principles

| Rule | Description |
|------|-------------|
| Orchestrator | 只规划、分配、验收，不写业务代码 |
| Worker | 实现 feature，每个 feature 一个小闭环 |
| Validator | 检查 feature 是否满足 validation-contract |
| Feature scope | 每个 feature 最多修改 1 个大模块 |
| Assertions | 每个 feature 最多 3 个 validation assertion |
| Milestone | 尽量只涉及 1-2 类 worker |
| No skip tests | 不允许跳过测试 |
| Key security | 不允许明文保存上游 key 或用户 API Key |
| Privacy | 默认不保存完整 prompt / response |
| Fallback safety | official_cloud 失败不能 fallback 到 experimental_proxy |
| DB compat | 所有 migration 必须兼容 SQLite / MySQL / PostgreSQL |

## Milestones

### M0 — 仓库调研与 Mission 基础文件
**Goal**: 梳理调用链路，建立 .factory/ 知识库，确认二开文件清单。
**Workers**: repo-discovery-worker
**Acceptance**: .factory/library/ 三个文件存在且内容准确；features.json 生成完毕。

### M1 — Provider 类型与字段基础迁移
**Goal**: 在 Channel 模型上增加 provider_type / risk_level / available_scope / visibility / manual_enable_required 字段，DB migration 通过三库。
**Workers**: db-migration-worker, provider-worker
**Acceptance**: `go test ./model/...` 通过；三库 AutoMigrate 无报错；Channel CRUD 包含新字段。

### M2 — Provider Account 与 Channel 关联
**Goal**: 建立 ProviderAccount 模型（加密存储上游 key），Channel 关联 ProviderAccountId。
**Workers**: db-migration-worker, provider-worker
**Acceptance**: ProviderAccount 表存在；Channel.ProviderAccountId 可为空（向后兼容）；key 加密存储验证。

### M3 — OpenAI-compatible 接口回归
**Goal**: 确认 GET /v1/models 和 POST /v1/chat/completions 在新字段迁移后仍正常工作。
**Workers**: openai-api-worker, test-validation-worker
**Acceptance**: 两个接口返回 200；usage 字段存在；现有单元测试全部通过。

### M4 — 正规 Provider 最小接入
**Goal**: 现有 OpenAI channel 自动标记为 official_cloud；Channel 创建/更新 API 接受 provider_type 参数。
**Workers**: provider-worker
**Acceptance**: 创建 channel 时可指定 provider_type；查询 channel 返回 provider_type；旧数据默认 official_cloud。

### M5 — 模型映射与价格基础
**Goal**: 模型路由尊重 provider_type；新 provider_type 的 channel 可正常参与 channel 选择。
**Workers**: routing-security-worker, billing-worker
**Acceptance**: 指定 provider_type=official_cloud 的 channel 可被正常路由；pricing 不报错。

### M6 — KiroGatewayAdapter 骨架
**Goal**: 新增 ChannelTypeKiroGateway (type=58)，实现 Adaptor 接口骨架，注册到 GetAdaptor 工厂。
**Workers**: provider-worker
**Acceptance**: `go build ./...` 通过；GetAdaptor(58) 返回非 nil；KiroGateway channel 可创建。

### M7 — experimental_proxy 访问控制
**Goal**: Distribute 中间件拒绝非 internal 用户调用 experimental_proxy channel；internal 用户需 explicit enable 后才能调用。
**Workers**: routing-security-worker
**Acceptance**: 普通用户调用 experimental_proxy 模型返回 403；internal 用户在 channel enabled 时返回正常响应。

### M8 — experimental_proxy 路由与 fallback 隔离
**Goal**: Channel 选择逻辑不允许从 official_cloud 自动 fallback 到 experimental_proxy；experimental_proxy 可一键禁用。
**Workers**: routing-security-worker
**Acceptance**: 禁用 experimental_proxy channel 后所有调用返回 503/404；official_cloud 失败不触发 experimental_proxy。

### M9 — 组织 / 项目基础表结构
**Goal**: 建立 Organization / Project stub 表，仅建表不实现业务逻辑。
**Workers**: db-migration-worker
**Acceptance**: 两张表存在；AutoMigrate 三库通过；不影响现有功能。

### M10 — API Key 组织项目绑定与安全
**Goal**: Token 增加可选 OrgId / ProjectId 字段；确认 key 以 hash 形式存储。
**Workers**: db-migration-worker, routing-security-worker
**Acceptance**: Token 表有新字段；key 存储验证（DB 中无明文 sk- 前缀完整 key）。

### M11 — usage_log 与隐私策略
**Goal**: Log 增加 provider_type 字段；默认不写入 prompt_text / completion_text；隐私策略可配置。
**Workers**: billing-worker
**Acceptance**: Log 表有 provider_type；`STORE_FULL_TEXT_ENABLED=false` 时 log.Other 无完整 prompt；单测通过。

### M12 — token 统计与基础扣费
**Goal**: 请求完成后 user.UsedQuota / token.UsedQuota 正确增加；log.Quota > 0。
**Workers**: billing-worker
**Acceptance**: 发起一次 chat/completions 后，DB 中对应 user 和 token 的 used_quota 增加；log 记录存在。

### M13 — 余额不足与管理员手动充值
**Goal**: RemainQuota=0 的 token 请求被拒绝（402/429）；管理员可通过 API 手动充值。
**Workers**: billing-worker, admin-ui-worker
**Acceptance**: 零额度 token 请求返回 402/429；POST /api/user/topup 成功后 user.Quota 增加；log type=manage 记录存在。

### M14 — 管理后台 Provider / Channel 最小页面
**Goal**: Channel 列表支持按 provider_type 过滤；experimental_proxy channel 在普通用户视图中不可见。
**Workers**: admin-ui-worker
**Acceptance**: 管理员可按 provider_type 筛选 channel；前端 channel 列表 API 对非管理员隐藏 experimental_proxy。

### M15 — 管理后台 API Key / usage_log / 余额页面
**Goal**: 管理员可查看 API Key 列表、usage_log、用户余额；支持手动充值操作。
**Workers**: admin-ui-worker
**Acceptance**: 三个管理页面可访问；usage_log 不展示完整 prompt；余额充值操作成功。

### M16 — Docker、文档和回归测试
**Goal**: docker-compose 可一键启动；smoke test 脚本覆盖 MVP 核心路径；go test ./... 全部通过。
**Workers**: deployment-worker, test-validation-worker
**Acceptance**: `docker compose up` 启动成功；smoke test 脚本 exit 0；`go test ./...` 无 FAIL。

## Workers

| Worker | Responsibility |
|--------|---------------|
| repo-discovery-worker | 仓库调研、调用链梳理、知识库文件生成 |
| db-migration-worker | GORM 模型变更、AutoMigrate、三库兼容性 |
| provider-worker | Channel 类型扩展、Adaptor 实现、Provider 注册 |
| routing-security-worker | Distribute 中间件、访问控制、fallback 隔离 |
| openai-api-worker | OpenAI-compatible 接口实现与回归 |
| billing-worker | 计费会话、quota 扣减、usage_log、隐私策略 |
| admin-ui-worker | 管理后台 API 和前端页面 |
| deployment-worker | Docker、compose、CI 配置 |
| test-validation-worker | 单元测试、集成测试、smoke test |

## Recommended Execution Order

```
M0 → M1 → M2 → M3 → M4 → M5 → M6 → M7 → M8
                                          ↓
M16 ← M15 ← M14 ← M13 ← M12 ← M11 ← M10 ← M9
```

M3 is a regression gate — if M3 fails, stop and fix before proceeding.
M7/M8 are security gates — must pass before M9+.

---

## Mandatory Feature Completion Rule

After completing **any** feature, the Worker MUST update both files before returning the final handoff:

- `docs/DEVELOPMENT_LOG.md`
- `.factory/mission-state.json`

**Do NOT claim a feature is completed unless both files are updated.**

### On Completion
- Add to `completed_features`; remove from `pending_features`
- Update `validation_status` for fulfilled assertions
- Append change record to `docs/DEVELOPMENT_LOG.md`
- Update `current_milestone`, `current_feature`, `next_recommended_action`

### On Blocked
- Add to `blocked_features`; keep out of `completed_features`
- Record `blocker`, `why_it_blocks_the_mission`, `minimal_fix_path`
- Do NOT continue to next feature automatically

### On Failed
- Mark `status: failed`; record failure reason and failed commands with exit codes
- Do NOT continue to next feature automatically

### Required handoff fields
```json
{
  "development_log_updated": true,
  "mission_state_updated": true
}
```

This rule applies to all milestones M0–M16 and all workers.
Full handoff templates are in `.factory/mission-state.json`.
