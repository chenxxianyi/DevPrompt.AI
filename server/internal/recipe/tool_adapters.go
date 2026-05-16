package recipe

import "strings"

func toolAdapterPrompt(tool string) string {
	switch strings.ToLower(tool) {
	case "cursor":
		return cursorAdapterPrompt
	case "claude code", "claude-code":
		return claudeCodeAdapterPrompt
	case "gpt", "chatgpt", "openai":
		return gptAdapterPrompt
	case "gemini":
		return geminiAdapterPrompt
	case "deepseek":
		return deepseekAdapterPrompt
	case "qwen", "通义千问":
		return qwenAdapterPrompt
	default:
		return ""
	}
}

const cursorAdapterPrompt = `
## 工具适配：Cursor

生成内容需要适配 Cursor IDE 的使用方式，请遵循以下原则：

- 强调代码库上下文的重要性，要求 AI 先阅读相关文件再动手修改。
- 指定需要操作的文件范围，避免大范围重构。
- 渐进式修改，每次只修改必要的部分。
- 不要删除用户已有的注释和文档。
- 修改后提供完整的文件内容，而不是 diff 片段。
- 每个功能点尽量在单个文件内完成，减少跨文件修改。
`

const claudeCodeAdapterPrompt = `
## 工具适配：Claude Code

生成内容需要适配 Claude Code CLI 工具的使用方式，请遵循以下原则：

- 明确任务边界，说明要做什么、不要做什么。
- 提供具体的文件路径和模块名称。
- 给出可执行的步骤，每步都有明确输入和预期输出。
- 包含测试命令和验证方式。
- 定义验收标准，让 AI 能自检。
- 不要覆盖用户已有的修改，先理解再行动。
`

const gptAdapterPrompt = `
## 工具适配：GPT / ChatGPT

生成内容需要适配 GPT 系列模型的使用方式，请遵循以下原则：

- 强调结构化输出，使用清晰的标题和分段。
- 确保方案完整度，不要省略中间步骤。
- 提高可读性，使用列表、表格等格式。
- 指定输出格式和长度要求。
- 在复杂任务中分解为子任务，逐步执行。
`

const geminiAdapterPrompt = `
## 工具适配：Gemini

生成内容需要适配 Google Gemini 的使用方式，请遵循以下原则：

- 强调结构化输出，使用 Markdown 分级标题。
- 提供完整的技术方案，包含选型理由。
- 注重可读性和逻辑连贯性。
- 指定输出格式，避免冗余信息。
`

const deepseekAdapterPrompt = `
## 工具适配：DeepSeek

生成内容需要适配 DeepSeek 的使用方式，请遵循以下原则：

- 强调代码实现的完整性和正确性。
- 提供详细的实现步骤和代码示例。
- 注重逻辑严密性，避免模糊表述。
- 指定输出格式和详细程度。
`

const qwenAdapterPrompt = `
## 工具适配：Qwen / 通义千问

生成内容需要适配通义千问的使用方式，请遵循以下原则：

- 强调结构化输出，使用清晰的标题和分段。
- 提供完整的技术方案和实现路径。
- 注重中文表达的准确性和专业性。
- 指定输出格式和详细程度。
`

func qualityModeDirective(mode string) string {
	switch mode {
	case "concise":
		return `## 质量模式：简洁

请以简洁方式输出：
- 每个章节控制在 2-3 句话以内。
- 只列出最关键的要点，省略显而易见的内容。
- 代码示例仅提供核心片段。
- 避免重复表述。
`
	case "expert":
		return `## 质量模式：专家

请以专家深度输出：
- 每个章节提供详尽的分析和论证。
- 考虑边界情况、性能影响和可维护性。
- 提供多种方案对比和选型理由。
- 代码示例完整可运行，包含错误处理。
- 包含架构图描述、数据流、模块依赖关系。
- 指出常见的陷阱和反模式。
`
	default:
		return `## 质量模式：标准

请以标准详细度输出：
- 每个章节提供足够的信息让读者理解核心思路。
- 给出关键代码示例，不需要完整实现。
- 提供明确的步骤和优先级。
- 包含必要的注意事项。
`
	}
}

func outputFormatDirective(format string) string {
	switch format {
	case "checklist":
		return `## 输出格式：清单

请使用 Checklist 格式输出，每个要点使用 - [ ] 或 - [x] 标记。确保每条可独立检查和验证。
`
	case "json":
		return `## 输出格式：JSON

请尽量使用 JSON 结构输出。如果整体无法用 JSON 表达，请在主要结构部分使用 JSON，说明部分使用 Markdown。确保 JSON 格式合法可解析。
`
	case "plain":
		return `## 输出格式：纯 Prompt

请直接输出可复制使用的 Prompt 文本，不包含额外说明和注释。内容应可直接粘贴到目标 AI 工具中使用。
`
	default:
		return ""
	}
}
