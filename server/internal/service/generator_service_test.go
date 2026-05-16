package service

import (
	"strings"
	"testing"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/recipe"
)

func TestExtractQualityOptionsDefaults(t *testing.T) {
	opts := ExtractQualityOptions(nil, nil, nil, nil, nil, nil)

	if opts.QualityMode != "standard" {
		t.Fatalf("expected default quality mode standard, got %q", opts.QualityMode)
	}
	if opts.OutputFormat != "markdown" {
		t.Fatalf("expected default output format markdown, got %q", opts.OutputFormat)
	}
	if opts.IncludeAcceptanceCriteria || opts.IncludeRiskCheck || opts.IncludeTestPlan || opts.IncludeDeploymentNotes {
		t.Fatalf("expected optional quality flags to default to false, got %+v", opts)
	}
}

func TestNormalizeQualityOptionsRejectsFreeAdvancedOptions(t *testing.T) {
	user := &model.User{Role: "user", MembershipLevel: "free"}
	opts := recipe.QualityOptions{
		QualityMode:  "expert",
		OutputFormat: "markdown",
	}

	_, err := normalizeQualityOptions(user, GenerateProjectRequest{}, opts)
	if err == nil || !strings.Contains(err.Error(), "expert") {
		t.Fatalf("expected free-tier expert mode rejection, got %v", err)
	}
}

func TestNormalizeQualityOptionsRejectsFreeOptimizeExpert(t *testing.T) {
	user := &model.User{Role: "user", MembershipLevel: "free"}

	_, err := normalizeQualityOptions(user, GenerateOptimizeRequest{OptimizeLevel: "expert"}, recipe.DefaultQualityOptions())
	if err == nil || !strings.Contains(err.Error(), "Prompt") {
		t.Fatalf("expected free-tier optimize rejection, got %v", err)
	}
}

func TestNormalizeQualityOptionsAllowsPaidAdvancedOptions(t *testing.T) {
	user := &model.User{Role: "user", MembershipLevel: "pro"}
	opts := recipe.QualityOptions{
		QualityMode:               "expert",
		OutputFormat:              "json",
		IncludeAcceptanceCriteria: true,
		IncludeRiskCheck:          true,
		IncludeTestPlan:           true,
		IncludeDeploymentNotes:    true,
	}

	got, err := normalizeQualityOptions(user, GenerateOptimizeRequest{OptimizeLevel: "expert"}, opts)
	if err != nil {
		t.Fatalf("expected paid-tier options to pass, got %v", err)
	}
	if got.QualityMode != "expert" || got.OutputFormat != "json" {
		t.Fatalf("unexpected normalized options: %+v", got)
	}
}

func TestBuildRecipeInputNormalizesOptimizeLevel(t *testing.T) {
	input, err := buildRecipeInput(string(recipe.GeneratorOptimize), GenerateOptimizeRequest{
		RawPrompt:     "raw",
		TargetTool:    "gpt",
		OptimizeLevel: "EXPERT",
	})
	if err != nil {
		t.Fatalf("expected optimize input build success, got %v", err)
	}

	got, ok := input.(recipe.OptimizeInput)
	if !ok {
		t.Fatalf("expected recipe.OptimizeInput, got %T", input)
	}
	if got.OptimizeLevel != "expert" {
		t.Fatalf("expected normalized optimize level expert, got %q", got.OptimizeLevel)
	}
}

func TestMaxTokensForGenerateUsesOptimizeLevel(t *testing.T) {
	got := maxTokensForGenerate(GenerateOptimizeRequest{OptimizeLevel: "expert"}, recipe.DefaultQualityOptions())
	if got != expertGenerateMaxTokens {
		t.Fatalf("expected expert optimize max tokens %d, got %d", expertGenerateMaxTokens, got)
	}
}

func TestNormalizeJSONOutputSupportsCodeFence(t *testing.T) {
	content := "```json\n{\"a\":1,\"b\":[true,false]}\n```"

	got, err := normalizeJSONOutput(content)
	if err != nil {
		t.Fatalf("expected valid fenced json, got %v", err)
	}
	if !strings.Contains(got, "\"a\": 1") {
		t.Fatalf("expected pretty json output, got %q", got)
	}
}

func TestNormalizeJSONOutputRejectsTrailingContent(t *testing.T) {
	_, err := normalizeJSONOutput("{\"a\":1}\nextra")
	if err == nil {
		t.Fatal("expected trailing-content error")
	}
}
