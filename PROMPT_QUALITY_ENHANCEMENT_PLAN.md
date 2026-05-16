# DevPrompt AI Prompt 生成质量增强方案

## 1. 背景

DevPrompt AI 当前已经具备四类核心生成能力：

- 项目 Prompt 生成
- Cursor Rules 生成
- Claude Code Prompt 生成
- Prompt 优化

现有实现主要通过后端固定的 system prompt builder 完成，例如 `buildProjectPrompt`、`buildCursorRulesPrompt`、`buildClaudeCodePrompt`、`buildOptimizePrompt`。这种方式适合 MVP，但随着模板类型、目标工具、会员能力和运营需求增加，会逐渐出现以下问题：

- Prompt 逻辑写死在代码中，难以迭代和灰度。
- 不同目标工具缺少深度适配。
- 输出结构不够稳定，用户每次得到的结果质量波动较大。
- 缺少验收标准、风险检查、测试建议等高价值内容。
- 无法沉淀用户反馈，也无法判断哪个 Prompt 策略效果更好。
- 管理后台无法运营 Prompt Recipe，只能维护模板、模型和基础数据。

因此，本方案目标是将 DevPrompt AI 从“固定 system prompt 生成器”升级为“可版本化、可运营、可评估、可持续优化的 Prompt Recipe 引擎”。

## 2. 产品目标

### 2.1 用户价值

- 生成结果更专业、更稳定、更可直接复制使用。
- 支持按目标工具生成更适配的 Prompt。
- 支持选择输出格式、详细程度和质量策略。
- 生成结果包含明确的验收标准、风险提示和后续行动建议。
- 用户可以反馈生成质量，系统后续可以据此优化。

### 2.2 业务价值

- Prompt 质量成为产品护城河，而不是单纯调用模型 API。
- 支持 Pro / Team / Enterprise 分层能力设计。
- 通过反馈数据找到高转化、高复用、高满意度的 Recipe。
- 管理后台可以运营和灰度 Prompt Recipe。
- 为未来团队协作、API 调用、私有化部署打基础。

### 2.3 技术目标

- 将 Prompt 构建逻辑从 API handler 中解耦。
- 支持 Recipe 版本管理、默认版本、状态切换和回滚。
- 生成记录保存 Recipe 版本、质量模式、输出格式等上下文。
- 为后续 A/B 测试、质量评分、用户反馈提供数据基础。

## 3. 总体方案

当前链路：

```text
生成类型 -> 固定 system prompt -> 调用 AI -> 保存生成记录
```

目标链路：

```text
生成类型
  + 目标工具
  + 输出格式
  + 质量模式
  + Recipe 版本
  + 用户输入
  -> Prompt Recipe 引擎
  -> 构建结构化 system prompt / user prompt
  -> 调用 AI
  -> 质量自检
  -> 保存生成记录
  -> 用户反馈
  -> 后台分析与迭代
```

核心模块建议：

```text
server/internal/recipe/
├─ engine.go              Recipe 选择与构建入口
├─ types.go               Recipe、BuildContext、QualityOptions 类型
├─ builtin.go             内置默认 Recipe
├─ project.go             项目 Prompt Recipe
├─ cursor_rules.go        Cursor Rules Recipe
├─ claude_code.go         Claude Code Recipe
├─ optimize.go            Prompt 优化 Recipe
└─ tool_adapters.go       Cursor / Claude Code / GPT 等工具适配策略
```

## 4. 分期路线图

## Phase 1：内置 Recipe 引擎

目标：在不大规模改数据库的情况下，显著提升生成质量。

### 范围

- 新增 `server/internal/recipe` 包。
- 将四个 `build*Prompt()` 从 `generator.go` 中迁移到 Recipe 引擎。
- 为四类生成建立统一输出结构。
- 加入工具适配、验收标准、风险检查、自检要求。
- 保持现有 API 请求参数基本不变，减少前端改动。

### 交付物

- `recipe.Engine`
- 四类内置 Recipe
- `QualityOptions` 默认策略
- 后端生成服务调用 Recipe 引擎

### 验收标准

- 四类生成结果结构稳定。
- 项目 Prompt 输出包含目标、假设、架构、模块、步骤、验收标准、风险和最终可复制 Prompt。
- Cursor Rules 输出可直接保存为 `.cursorrules`。
- Claude Code 输出包含上下文、任务、约束、执行步骤和验收标准。
- Prompt 优化输出包含优化后 Prompt、优化说明、适用场景和风险提示。

## Phase 2：前端质量控制项

目标：让用户可控地提升生成质量。

### 新增前端控制项

建议在生成器工作台增加：

- 详细程度：简洁 / 标准 / 专家
- 输出格式：Markdown / Checklist / JSON / 纯 Prompt
- 是否包含验收标准
- 是否包含风险检查
- 是否包含测试建议
- 是否包含部署建议

### 请求参数建议

```ts
type QualityMode = 'concise' | 'standard' | 'expert'
type OutputFormat = 'markdown' | 'checklist' | 'json' | 'plain'

interface QualityOptions {
  qualityMode: QualityMode
  outputFormat: OutputFormat
  includeAcceptanceCriteria: boolean
  includeRiskCheck: boolean
  includeTestPlan: boolean
  includeDeploymentNotes: boolean
}
```

### 后端请求结构建议

四类请求都可以增加可选字段：

```go
type QualityOptions struct {
    QualityMode               string `json:"qualityMode"`
    OutputFormat              string `json:"outputFormat"`
    IncludeAcceptanceCriteria bool   `json:"includeAcceptanceCriteria"`
    IncludeRiskCheck          bool   `json:"includeRiskCheck"`
    IncludeTestPlan           bool   `json:"includeTestPlan"`
    IncludeDeploymentNotes    bool   `json:"includeDeploymentNotes"`
}
```

### 验收标准

- 用户切换“简洁 / 标准 / 专家”后，输出长度和深度明显不同。
- 用户选择 JSON 后，输出应尽量保持可解析结构。
- 用户关闭风险检查时，结果中不再强制出现风险章节。
- 旧客户端不传字段时，后端使用默认策略，不影响兼容性。

## Phase 3：Recipe 版本化与后台运营

目标：让 Prompt 生成策略成为可运营资产。

### 新增数据表：`prompt_recipes`

```text
id
type                 project / cursor-rules / claude-code / optimize
target_tool          common / cursor / claude-code / gpt / gemini / deepseek / qwen
version              v1 / v2 / v3
name
description
system_prompt
user_template
output_schema
quality_rubric
status               draft / active / disabled
is_default
created_by
created_at
updated_at
```

### 后台新增页面

- Recipe 列表
- 新建 / 编辑 Recipe
- 设置默认 Recipe
- 启用 / 禁用 Recipe
- 查看 Recipe 使用次数、成功率、复制率、反馈评分
- Recipe 版本回滚

### 验收标准

- 管理员可以不发版修改 Recipe。
- 每个生成类型至少有一个默认 active Recipe。
- 生成记录中能看到使用的 Recipe ID 和版本。
- 禁用 Recipe 后不会再被新请求选中。

## Phase 4：质量评分与反馈闭环

目标：用真实用户行为优化 Prompt 质量。

### 新增反馈入口

在生成结果区域增加：

- 有用 / 无用
- 结果太泛泛
- 格式不好
- 不够准确
- 缺少上下文
- 可以直接使用
- 用户备注

### 新增行为事件

建议记录：

- 生成成功
- 生成失败
- 复制结果
- 导出 Markdown
- 再次生成
- 使用历史记录再次生成
- 用户提交反馈
- 用户点击升级

### 新增数据表：`generated_prompt_feedbacks`

```text
id
generated_prompt_id
user_id
rating              useful / not_useful
reason              too_generic / inaccurate / bad_format / missing_context / useful
comment
created_at
```

### 可计算指标

- Recipe 使用次数
- Recipe 成功率
- 平均生成耗时
- 平均 token 成本
- 复制率
- 导出率
- 再生成率
- 正向反馈率
- 负向反馈原因分布

## 5. 四类 Recipe 设计

## 5.1 项目 Prompt Recipe

### 适用场景

用户希望根据项目名称、项目类型、技术栈、功能需求和目标 AI 工具生成完整开发 Prompt。

### 推荐输出结构

```md
# 项目开发 Prompt

## 1. 任务目标
说明这个 Prompt 要让 AI 完成什么。

## 2. 项目背景与合理假设
列出用户已提供信息，以及系统基于输入做出的假设。

## 3. 技术栈与架构建议
说明技术栈选择、前后端架构、数据层、部署方式。

## 4. 功能模块拆解
按模块列出核心功能、边界和优先级。

## 5. 页面 / API / 数据模型建议
根据项目类型给出页面结构、接口设计或数据模型建议。

## 6. 开发步骤
给出可执行的开发顺序。

## 7. 验收标准
列出功能、质量、测试、性能方面的验收条件。

## 8. 风险与注意事项
指出需求模糊点、技术风险、安全风险、数据风险。

## 9. 可直接复制使用的最终 Prompt
给出一段面向目标 AI 工具的完整 Prompt。
```

### 工具适配策略

- Cursor：强调代码库上下文、文件范围、渐进式修改、不要大范围重构。
- Claude Code：强调任务边界、执行步骤、文件路径、测试命令、验收标准。
- GPT / Gemini / DeepSeek / Qwen：强调结构化输出、方案完整度、可读性。

## 5.2 Cursor Rules Recipe

### 推荐输出结构

```md
# .cursorrules

## AI 行为准则
## 项目技术栈
## 代码风格
## 文件组织
## 命名规范
## 禁止事项
## 测试要求
## 安全要求
## 输出规范
```

### 质量要求

- 避免空泛规则。
- 每条规则必须可执行、可检查。
- 针对语言和框架生成专属规则。
- 明确禁止事项，如不要引入未确认依赖、不要重写无关文件。

## 5.3 Claude Code Recipe

### 推荐输出结构

```md
# Claude Code Task Prompt

## Objective
## Context
## Scope
## Constraints
## Implementation Steps
## Files To Inspect
## Acceptance Criteria
## Test Plan
## Safety Rules
## Final Instruction
```

### 质量要求

- 强调任务边界。
- 明确不要覆盖用户已有修改。
- 要求 Claude Code 先理解代码再修改。
- 给出测试和验证命令占位。

## 5.4 Prompt 优化 Recipe

### 推荐输出结构

```md
# 优化后的 Prompt

...

---

## 优化说明
## 关键改进点
## 适用场景
## 可能风险
## 可继续补充的信息
```

### 质量要求

- 不改变用户原始意图。
- 补充角色、目标、上下文、约束、输出格式。
- 对模糊需求给出假设。
- 专家模式下应给出更严格的评价标准。

## 6. 后端改造建议

## 6.1 新增类型

```go
type GeneratorType string

const (
    GeneratorProject     GeneratorType = "project"
    GeneratorCursorRules GeneratorType = "cursor-rules"
    GeneratorClaudeCode  GeneratorType = "claude-code"
    GeneratorOptimize    GeneratorType = "optimize"
)

type QualityOptions struct {
    QualityMode               string
    OutputFormat              string
    IncludeAcceptanceCriteria bool
    IncludeRiskCheck          bool
    IncludeTestPlan           bool
    IncludeDeploymentNotes    bool
}

type BuildContext struct {
    Type        GeneratorType
    TargetTool  string
    Input       any
    Quality     QualityOptions
    RecipeID    uint64
    Version     string
}

type BuiltPrompt struct {
    SystemPrompt string
    UserPrompt   string
    RecipeID     uint64
    RecipeVersion string
}
```

## 6.2 Recipe Engine 接口

```go
type Engine interface {
    Build(ctx BuildContext) (*BuiltPrompt, error)
}
```

## 6.3 GeneratorService 调整

建议将当前：

```go
Generate(userID, genType, title, input, systemPrompt)
```

逐步升级为：

```go
Generate(userID, genType, title, input, qualityOptions)
```

由 `GeneratorService` 内部调用 Recipe Engine 构建 system prompt 和 user prompt。

## 6.4 生成记录字段扩展

建议给 `generated_prompts` 增加：

```text
recipe_id
recipe_version
quality_mode
output_format
feedback_score
copied_at
exported_at
```

第一期可以先不加表字段，只在 `Input` JSON 中附加质量参数。进入 Phase 3 时再正式迁移字段。

## 7. 前端改造建议

## 7.1 生成器工作台

在每个 Tab 底部或右侧增加“生成质量”区域：

- 详细程度：简洁 / 标准 / 专家
- 输出格式：Markdown / 清单 / JSON / 纯 Prompt
- 高级开关：
  - 验收标准
  - 风险检查
  - 测试建议
  - 部署建议

### 默认值

```ts
const defaultQualityOptions = {
  qualityMode: 'standard',
  outputFormat: 'markdown',
  includeAcceptanceCriteria: true,
  includeRiskCheck: true,
  includeTestPlan: false,
  includeDeploymentNotes: false,
}
```

## 7.2 生成结果区域

新增：

- 质量标签：使用的模式、格式、目标工具。
- 反馈按钮：有用 / 无用。
- 负反馈原因选择。
- “生成更详细版本”快捷按钮。
- “生成更简洁版本”快捷按钮。

## 7.3 会员分层建议

Free：

- 标准模式
- Markdown 输出
- 基础 Recipe

Pro：

- 专家模式
- 风险检查
- 测试建议
- 更多目标工具适配

Team：

- 团队共享 Recipe
- 团队模板库
- API 调用
- 团队反馈分析

Enterprise：

- 私有 Recipe
- 私有模型配置
- SLA
- 审计日志
- 私有化部署

## 8. 管理后台改造建议

## 8.1 Prompt Recipe 管理

新增菜单：`Prompt Recipe 管理`

页面能力：

- Recipe 列表
- 按生成类型筛选
- 按目标工具筛选
- 创建 Recipe
- 编辑 Recipe
- 启用 / 禁用
- 设置默认版本
- 查看效果数据

## 8.2 Recipe 效果分析

建议展示：

- 使用次数
- 成功次数
- 失败次数
- 成功率
- 平均耗时
- 平均 token
- 复制率
- 导出率
- 正反馈率
- 负反馈原因 Top 5

## 9. 风险与处理策略

| 风险 | 说明 | 处理策略 |
| --- | --- | --- |
| 输出过长 | 专家模式可能消耗更多 token | 增加 MaxTokens 策略和会员限制 |
| JSON 不稳定 | 模型可能输出非法 JSON | 增加格式修复和解析校验 |
| Recipe 失控 | 后台修改 Prompt 可能造成质量下降 | 支持草稿、预览、灰度、回滚 |
| 成本增加 | 更复杂 Prompt 会增加 token | 建立 token 成本看板 |
| 反馈稀疏 | 用户不一定主动反馈 | 用复制、导出、再生成等行为作为隐式反馈 |
| 工具适配维护成本高 | Cursor / Claude Code 等工具会变化 | 使用 adapter 分层，减少主 Recipe 变更 |

## 10. MVP 推荐实施清单

建议第一轮只做高收益、低风险的内容：

1. 新增 `server/internal/recipe` 包。
2. 将四类固定 prompt builder 迁移为内置 Recipe。
3. 给四类输出统一加入：
   - 输出结构
   - 验收标准
   - 风险检查
   - 自检要求
   - 最终可复制 Prompt
4. 前端增加两个质量控制项：
   - 详细程度
   - 输出格式
5. 后端请求结构兼容旧参数，新增质量参数为可选。
6. 生成记录 `Input` 中保存质量参数。
7. 生成结果区域增加“有用 / 无用”反馈按钮。

## 11. 推荐里程碑

### 第 1 周

- 完成 Recipe Engine 基础结构。
- 完成四类内置 Recipe。
- 完成后端生成链路接入。

### 第 2 周

- 前端生成器增加质量控制项。
- 优化生成结果展示。
- 完成基础测试和回归。

### 第 3 周

- 增加反馈接口和反馈表。
- 后台增加基础反馈查看。
- 统计复制率、导出率、再生成率。

### 第 4 周

- 增加 Recipe 数据表。
- 后台支持 Recipe 管理。
- 支持默认版本和回滚。

## 12. 最终验收标准

产品侧：

- 用户能明显感知生成结果更完整、更专业。
- 用户可以选择生成质量和输出格式。
- 结果中包含可执行的验收标准和风险提示。
- 用户可以对生成结果做反馈。

研发侧：

- Prompt 构建逻辑不再散落在 API handler。
- 新增一种生成策略不需要大改主链路。
- Recipe 可以版本化、回滚和统计。
- 生成记录可追踪使用的 Recipe 与质量参数。

运营侧：

- 后台能看到不同 Recipe 的使用与反馈数据。
- 可以根据反馈迭代 Recipe。
- 可以将高级 Recipe 作为会员权益。

## 13. 结论

Prompt 质量增强不应该只靠“把 system prompt 写长”，而应该建设一套可运营的 Recipe 体系。短期通过内置 Recipe 和质量控制项提升结果稳定性，中期通过反馈和数据指标找到高质量策略，长期通过版本化 Recipe 和后台运营形成真正的产品护城河。

推荐优先执行 Phase 1 和 Phase 2 的核心部分，这部分改动相对可控，但用户感知最明显。
