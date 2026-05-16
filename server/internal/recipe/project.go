package recipe

import (
	"fmt"
	"strings"
)

type ProjectRecipe struct{}

func NewProjectRecipe() *ProjectRecipe {
	return &ProjectRecipe{}
}

func (r *ProjectRecipe) Type() GeneratorType { return GeneratorProject }

func (r *ProjectRecipe) Build(ctx BuildContext) (*BuiltPrompt, error) {
	input, err := requireInput[ProjectInput](ctx.Input)
	if err != nil {
		return nil, fmt.Errorf("invalid project input: %w", err)
	}

	tool := input.TargetAiTool
	if tool == "" {
		tool = ctx.TargetTool
	}

	systemPrompt := r.buildSystemPrompt(tool, ctx.Quality)
	userPrompt := r.buildUserPrompt(input)

	return &BuiltPrompt{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	}, nil
}

func (r *ProjectRecipe) buildSystemPrompt(tool string, q QualityOptions) string {
	var b strings.Builder

	b.WriteString(`你是一位资深的软件架构师和技术顾问，擅长根据用户需求生成专业、结构化、可执行的开发 Prompt。

## 输出结构要求

请严格按照以下结构输出：

### 1. 任务目标
简要说明这个 Prompt 要让 AI 完成什么，用一句话概括核心目标。

### 2. 项目背景与合理假设
- 列出用户已提供的所有信息。
- 基于项目类型和技术栈，列出系统做出的合理假设。
- 对模糊点标注"[需确认]"。

### 3. 技术栈与架构建议
- 说明技术栈选择及理由。
- 描述前后端架构、数据层、部署方式。
- 如果有多个架构选择，给出推荐及理由。

### 4. 功能模块拆解
- 按模块列出核心功能、边界和优先级（P0/P1/P2）。
- 说明模块间的依赖关系。

### 5. 页面 / API / 数据模型建议
- 根据项目类型给出页面结构、接口设计或数据模型建议。
- 使用表格或列表清晰展示。

### 6. 开发步骤
- 给出可执行的开发顺序，标注每步的预计耗时和前置依赖。
- 优先实现核心功能，后续迭代完善。

### 7. 可直接复制使用的最终 Prompt
- 给出一段面向目标 AI 工具的完整 Prompt，用户可直接复制使用。
- 这段 Prompt 应包含上述所有关键信息，结构清晰，语言简洁。
`)

	if q.IncludeAcceptanceCriteria {
		b.WriteString(`
### 8. 验收标准
- 列出功能验收条件：每个核心功能如何判定完成。
- 列出质量验收条件：代码质量、测试覆盖率、性能指标。
- 列出安全验收条件：认证授权、数据保护、输入校验。
`)
	}

	if q.IncludeRiskCheck {
		b.WriteString(`
### 9. 风险与注意事项
- 指出需求模糊点及建议处理方式。
- 列出技术风险及缓解策略。
- 列出安全风险及防护建议。
- 列出数据风险及备份/恢复建议。
`)
	}

	if q.IncludeTestPlan {
		b.WriteString(`
### 10. 测试建议
- 单元测试：列出需要覆盖的核心函数和模块。
- 集成测试：列出关键业务流程的测试场景。
- E2E 测试：列出用户核心路径的测试用例。
`)
	}

	if q.IncludeDeploymentNotes {
		b.WriteString(`
### 11. 部署建议
- 环境配置：开发/测试/生产环境的差异。
- 部署流程：CI/CD 流程建议。
- 监控告警：关键指标和告警阈值。
`)
	}

	// 工具适配
	if adapter := toolAdapterPrompt(tool); adapter != "" {
		b.WriteString(adapter)
	}

	// 质量模式
	b.WriteString(qualityModeDirective(q.QualityMode))

	// 输出格式
	b.WriteString(outputFormatDirective(q.OutputFormat))

	b.WriteString(`
## 自检要求

输出前请自检：
1. 是否覆盖了用户提供的所有信息？
2. 每个章节是否都有实质内容，没有空泛描述？
3. 最终可复制的 Prompt 是否完整、可直接使用？
4. 技术建议是否与用户选择的技术栈一致？
`)

	return b.String()
}

func (r *ProjectRecipe) buildUserPrompt(input ProjectInput) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("项目名称：%s\n", input.ProjectName))
	b.WriteString(fmt.Sprintf("项目类型：%s\n", input.ProjectType))

	if len(input.TechStack) > 0 {
		b.WriteString(fmt.Sprintf("技术栈：%s\n", strings.Join(input.TechStack, "、")))
	}

	if len(input.Features) > 0 {
		b.WriteString("功能需求：\n")
		for _, f := range input.Features {
			b.WriteString(fmt.Sprintf("- %s\n", f))
		}
	}

	if input.TargetAiTool != "" {
		b.WriteString(fmt.Sprintf("目标 AI 工具：%s\n", input.TargetAiTool))
	}

	return b.String()
}
