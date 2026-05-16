package recipe

import (
	"fmt"
	"strings"
)

type OptimizeRecipe struct{}

func NewOptimizeRecipe() *OptimizeRecipe {
	return &OptimizeRecipe{}
}

func (r *OptimizeRecipe) Type() GeneratorType { return GeneratorOptimize }

func (r *OptimizeRecipe) Build(ctx BuildContext) (*BuiltPrompt, error) {
	input, err := requireInput[OptimizeInput](ctx.Input)
	if err != nil {
		return nil, fmt.Errorf("invalid optimize input: %w", err)
	}

	tool := input.TargetTool
	if tool == "" {
		tool = ctx.TargetTool
	}

	level := input.OptimizeLevel
	if level == "" {
		level = ctx.Quality.QualityMode
		if level == "" {
			level = "standard"
		}
	}

	systemPrompt := r.buildSystemPrompt(tool, level, ctx.Quality)
	userPrompt := r.buildUserPrompt(input)

	return &BuiltPrompt{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	}, nil
}

func (r *OptimizeRecipe) buildSystemPrompt(tool string, level string, q QualityOptions) string {
	var b strings.Builder

	levelDesc := "基础优化"
	if level == "professional" || level == "expert" {
		levelDesc = "专业级优化"
	}
	if level == "expert" {
		levelDesc = "专家级优化"
	}

	b.WriteString(fmt.Sprintf(`你是一位 Prompt Engineering 专家，擅长对用户提供的原始 Prompt 进行%s。

## 核心原则

1. 不改变用户原始意图，只增强表达和结构。
2. 补充缺失的角色定义、目标、上下文、约束和输出格式。
3. 对模糊需求给出合理假设，标注"[需确认]"。
4. 优化后的 Prompt 应比原始 Prompt 更清晰、更完整、更可执行。

## 输出结构要求

### 优化后的 Prompt
直接输出优化后的完整 Prompt，用户应可直接复制使用。

---

### 优化说明
- 简要说明做了哪些优化。
- 每个优化点用一句话概括。

### 关键改进点
- 列出具体的改进项，对比优化前后的差异。
- 说明每项改进的价值。

### 适用场景
- 说明优化后的 Prompt 适合在哪些场景使用。
- 指出可能不适用的场景。

`, levelDesc))

	if q.IncludeRiskCheck {
		b.WriteString(`### 可能风险
- 指出优化可能带来的副作用。
- 说明哪些假设可能不成立。
- 提醒用户需要根据实际情况调整的部分。
`)
	}

	b.WriteString(`### 可继续补充的信息
- 列出如果用户提供更多信息，可以进一步优化的方向。
- 包括：角色背景、输出格式要求、约束条件、示例等。
`)

	if q.IncludeAcceptanceCriteria {
		b.WriteString(`
### 验收标准
- 优化后的 Prompt 是否保留了原始意图？
- 是否补充了角色、目标、上下文、约束、输出格式？
- 是否可直接复制使用？
`)
	}

	// 工具适配
	if adapter := toolAdapterPrompt(tool); adapter != "" {
		b.WriteString(adapter)
	}

	// 专家模式额外要求
	if level == "expert" {
		b.WriteString(`
## 专家级优化额外要求

- 使用 Chain-of-Thought 结构，引导 AI 逐步推理。
- 添加 Few-shot 示例（如果适用）。
- 定义输出格式模板。
- 设置自检和验证步骤。
- 考虑边界情况和异常处理。
- 给出更严格的评价标准。
`)
	}

	// 输出格式
	b.WriteString(outputFormatDirective(q.OutputFormat))

	b.WriteString(`
## 自检要求

输出前请自检：
1. 优化后的 Prompt 是否保留了用户的原始意图？
2. 是否补充了角色、目标、上下文、约束、输出格式？
3. 优化说明是否清晰，用户能理解改了什么？
4. 优化后的 Prompt 是否可直接复制使用？
`)

	return b.String()
}

func (r *OptimizeRecipe) buildUserPrompt(input OptimizeInput) string {
	var b strings.Builder

	b.WriteString("原始 Prompt：\n")
	b.WriteString(input.RawPrompt)
	b.WriteString("\n")

	if input.TargetTool != "" {
		b.WriteString(fmt.Sprintf("目标 AI 工具：%s\n", input.TargetTool))
	}
	if input.OptimizeLevel != "" {
		b.WriteString(fmt.Sprintf("优化级别：%s\n", input.OptimizeLevel))
	}

	return b.String()
}
