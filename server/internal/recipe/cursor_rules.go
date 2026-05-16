package recipe

import (
	"fmt"
	"strings"
)

type CursorRulesRecipe struct{}

func NewCursorRulesRecipe() *CursorRulesRecipe {
	return &CursorRulesRecipe{}
}

func (r *CursorRulesRecipe) Type() GeneratorType { return GeneratorCursorRules }

func (r *CursorRulesRecipe) Build(ctx BuildContext) (*BuiltPrompt, error) {
	input, err := requireInput[CursorRulesInput](ctx.Input)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor rules input: %w", err)
	}

	systemPrompt := r.buildSystemPrompt(ctx.Quality)
	userPrompt := r.buildUserPrompt(input)

	return &BuiltPrompt{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	}, nil
}

func (r *CursorRulesRecipe) buildSystemPrompt(q QualityOptions) string {
	var b strings.Builder

	b.WriteString(`你是一位 Cursor IDE 配置专家，擅长根据项目的技术栈和编码规范生成专业、可执行、可直接保存为 .cursorrules 的配置文件。

## 输出结构要求

请严格按照以下结构输出，生成的内容应可直接保存为 .cursorrules 文件使用：

### AI 行为准则
- 定义 AI 在项目中的角色和行为边界。
- 明确 AI 应该优先做什么、避免做什么。
- 每条准则必须可执行、可检查，不要写空泛的描述。

### 项目技术栈
- 列出项目使用的语言、框架、运行时版本。
- 指定包管理器和构建工具。

### 代码风格
- 定义缩进、命名、注释等代码风格规则。
- 针对语言和框架生成专属规则，不要泛泛而谈。
- 每条规则给出具体示例。

### 文件组织
- 定义目录结构和文件命名规范。
- 说明模块划分原则。

### 命名规范
- 变量、函数、类、文件、目录的命名规则。
- 给出正例和反例。

### 禁止事项
- 明确列出 AI 不应该做的事情。
- 包括：不要引入未确认的依赖、不要重写无关文件、不要修改配置文件等。
- 每条禁止事项说明原因。

### 测试要求
- 定义测试覆盖率和测试风格要求。
- 指定测试框架和测试文件命名规范。

### 安全要求
- 列出安全编码规范。
- 包括：输入校验、认证授权、敏感数据处理等。
`)

	if q.IncludeAcceptanceCriteria {
		b.WriteString(`
### 验收标准
- 定义 .cursorrules 文件自身的验收条件。
- 列出每条规则的可检查标准。
`)
	}

	if q.IncludeRiskCheck {
		b.WriteString(`
### 风险与注意事项
- 指出规则可能带来的副作用。
- 标注需要根据项目实际情况调整的规则。
`)
	}

	// 质量模式
	b.WriteString(qualityModeDirective(q.QualityMode))

	// 输出格式
	b.WriteString(outputFormatDirective(q.OutputFormat))

	b.WriteString(`
## 质量要求

- 避免空泛规则，每条规则必须可执行、可检查。
- 针对语言和框架生成专属规则，不要写通用废话。
- 明确禁止事项，说明原因。
- 生成的内容应可直接保存为 .cursorrules 使用。

## 自检要求

输出前请自检：
1. 每条规则是否都有具体示例或可检查标准？
2. 是否针对用户指定的语言和框架生成了专属规则？
3. 禁止事项是否明确且合理？
4. 整体内容是否可直接保存为 .cursorrules 使用？
`)

	return b.String()
}

func (r *CursorRulesRecipe) buildUserPrompt(input CursorRulesInput) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("编程语言：%s\n", input.Language))
	if input.Framework != "" {
		b.WriteString(fmt.Sprintf("框架：%s\n", input.Framework))
	}
	if input.CodeStyle != "" {
		b.WriteString(fmt.Sprintf("代码风格偏好：%s\n", input.CodeStyle))
	}
	if len(input.Rules) > 0 {
		b.WriteString("额外规则要求：\n")
		for _, rule := range input.Rules {
			b.WriteString(fmt.Sprintf("- %s\n", rule))
		}
	}

	return b.String()
}
