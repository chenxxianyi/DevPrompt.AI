package recipe

import (
	"encoding/json"
	"fmt"

	"devprompt-ai/internal/repository"
)

type DbRecipeProvider struct {
	recipeRepo *repository.PromptRecipeRepository
}

func NewDbRecipeProvider(recipeRepo *repository.PromptRecipeRepository) *DbRecipeProvider {
	return &DbRecipeProvider{recipeRepo: recipeRepo}
}

func (p *DbRecipeProvider) Build(ctx BuildContext) (*BuiltPrompt, error) {
	recipe, err := p.recipeRepo.FindDefaultByType(string(ctx.Type))
	if err != nil {
		return nil, err
	}
	if recipe == nil {
		return nil, fmt.Errorf("no active DB recipe for type: %s", ctx.Type)
	}

	userPrompt := buildUserPromptFromTemplate(recipe.UserTemplate, ctx.Input)

	return &BuiltPrompt{
		SystemPrompt:  recipe.SystemPrompt,
		UserPrompt:    userPrompt,
		RecipeID:      recipe.ID,
		RecipeVersion: recipe.Version,
	}, nil
}

func buildUserPromptFromTemplate(template string, input any) string {
	if template == "" {
		raw, _ := json.Marshal(input)
		return string(raw)
	}
	return template
}
