package recipe

import (
	"fmt"
	"strings"
)

type ClaudeCodeRecipe struct{}

func NewClaudeCodeRecipe() *ClaudeCodeRecipe {
	return &ClaudeCodeRecipe{}
}

func (r *ClaudeCodeRecipe) Type() GeneratorType { return GeneratorClaudeCode }

func (r *ClaudeCodeRecipe) Build(ctx BuildContext) (*BuiltPrompt, error) {
	input, err := requireInput[ClaudeCodeInput](ctx.Input)
	if err != nil {
		return nil, fmt.Errorf("invalid claude code input: %w", err)
	}

	systemPrompt := r.buildSystemPrompt(ctx.Quality)
	userPrompt := r.buildUserPrompt(input)

	return &BuiltPrompt{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	}, nil
}

func (r *ClaudeCodeRecipe) buildSystemPrompt(q QualityOptions) string {
	var b strings.Builder

	b.WriteString(`你是一位 Claude Code CLI 工具的使用专家，擅长生成结构化、可执行的任务 Prompt，让 Claude Code 能高效完成开发任务。

## 输出结构要求

请严格按照以下结构输出，生成的内容应可直接复制到 Claude Code 中使用：

### Objective
- 用一句话明确任务目标。
- 说明完成后的预期结果。

### Context
- 提供任务所需的背景信息。
- 说明项目当前状态和相关代码位置。
- 列出已知的约束和前提条件。

### Scope
- 明确任务范围：要做什么。
- 明确不在范围内：不要做什么。
- 标注任务边界，避免 AI 越界修改。

### Constraints
- 技术约束：语言版本、框架版本、兼容性要求。
- 业务约束：不能违反的业务规则。
- 代码约束：不要覆盖用户已有修改、不要引入未确认依赖。

### Implementation Steps
- 给出可执行的步骤，每步都有明确输入和预期输出。
- 步骤之间有清晰的依赖关系。
- 每步标注需要操作的文件或模块。

### Files To Inspect
- 列出 AI 在动手修改前应先阅读的文件。
- 说明每个文件的作用和需要关注的重点。

### Final Instruction
- 给出一段可直接复制使用的完整 Prompt。
- 这段 Prompt 应包含上述所有关键信息，结构清晰。
`)

	if q.IncludeAcceptanceCriteria {
		b.WriteString(`
### Acceptance Criteria
- 列出功能验收条件：每个核心功能如何判定完成。
- 列出质量验收条件：代码质量、测试覆盖率。
- 使用可量化的标准，避免模糊表述。
`)
	}

	if q.IncludeTestPlan {
		b.WriteString(`
### Test Plan
- 列出需要编写的测试用例。
- 给出测试命令占位，如：` + "`go test ./...`" + `、` + "`npm test`" + `。
- 说明如何验证修改没有引入回归。
`)
	}

	if q.IncludeRiskCheck {
		b.WriteString(`
### Safety Rules
- 不要覆盖用户已有的修改。
- 修改前先理解代码逻辑。
- 不要删除现有测试用例。
- 不要修改无关文件。
- 每次修改后运行测试验证。
`)
	}

	// 质量模式
	b.WriteString(qualityModeDirective(q.QualityMode))

	// 输出格式
	b.WriteString(outputFormatDirective(q.OutputFormat))

	b.WriteString(`
## 质量要求

- 强调任务边界，明确做什么、不做什么。
- 给出具体文件路径和模块名称。
- 步骤可执行，每步有明确输入和预期输出。
- 包含测试和验证命令。

## 自检要求

输出前请自检：
1. Objective 是否一句话说清了目标？
2. Scope 是否明确划定了边界？
3. Implementation Steps 是否可执行？
4. Final Instruction 是否完整、可直接复制使用？
5. 是否避免了空泛描述，每条都有实质内容？
`)

	return b.String()
}

func (r *ClaudeCodeRecipe) buildUserPrompt(input ClaudeCodeInput) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("任务描述：%s\n", input.Task))
	if input.Context != "" {
		b.WriteString(fmt.Sprintf("上下文信息：%s\n", input.Context))
	}
	if len(input.Requirements) > 0 {
		b.WriteString("具体要求：\n")
		for _, req := range input.Requirements {
			b.WriteString(fmt.Sprintf("- %s\n", req))
		}
	}

	return b.String()
}
