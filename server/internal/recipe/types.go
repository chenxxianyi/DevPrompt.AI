package recipe

import "fmt"

type GeneratorType string

const (
	GeneratorProject     GeneratorType = "project"
	GeneratorCursorRules GeneratorType = "cursor-rules"
	GeneratorClaudeCode  GeneratorType = "claude-code"
	GeneratorOptimize    GeneratorType = "optimize"
)

type ProjectInput struct {
	ProjectName  string
	ProjectType  string
	TechStack    []string
	Features     []string
	TargetAiTool string
}

type CursorRulesInput struct {
	Language  string
	Framework string
	CodeStyle string
	Rules     []string
}

type ClaudeCodeInput struct {
	Task         string
	Context      string
	Requirements []string
}

type OptimizeInput struct {
	RawPrompt     string
	TargetTool    string
	OptimizeLevel string
}

type QualityOptions struct {
	QualityMode               string // concise / standard / expert
	OutputFormat              string // markdown / checklist / json / plain
	IncludeAcceptanceCriteria bool
	IncludeRiskCheck          bool
	IncludeTestPlan           bool
	IncludeDeploymentNotes    bool
}

func DefaultQualityOptions() QualityOptions {
	return QualityOptions{
		QualityMode:  "standard",
		OutputFormat: "markdown",
	}
}

func IsValidQualityMode(mode string) bool {
	switch mode {
	case "concise", "standard", "expert":
		return true
	default:
		return false
	}
}

func IsValidOutputFormat(format string) bool {
	switch format {
	case "markdown", "checklist", "json", "plain":
		return true
	default:
		return false
	}
}

func IsValidOptimizeLevel(level string) bool {
	switch level {
	case "", "basic", "professional", "expert":
		return true
	default:
		return false
	}
}

type BuildContext struct {
	Type       GeneratorType
	TargetTool string
	Input      any
	Quality    QualityOptions
}

type BuiltPrompt struct {
	SystemPrompt  string
	UserPrompt    string
	RecipeID      uint64
	RecipeVersion string
}

type Recipe interface {
	Type() GeneratorType
	Build(ctx BuildContext) (*BuiltPrompt, error)
}

func requireInput[T any](input any) (T, error) {
	var zero T

	if typed, ok := input.(T); ok {
		return typed, nil
	}
	if typedPtr, ok := input.(*T); ok && typedPtr != nil {
		return *typedPtr, nil
	}

	return zero, fmt.Errorf("unexpected input type %T", input)
}
