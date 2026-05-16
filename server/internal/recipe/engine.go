package recipe

import "fmt"

type Engine struct {
	builtinRecipes map[GeneratorType]Recipe
	dbRecipe       *DbRecipeProvider
}

func NewEngine() *Engine {
	e := &Engine{
		builtinRecipes: make(map[GeneratorType]Recipe),
	}
	e.registerBuiltinRecipes()
	return e
}

func NewEngineWithDB(dbProvider *DbRecipeProvider) *Engine {
	e := &Engine{
		builtinRecipes: make(map[GeneratorType]Recipe),
		dbRecipe:       dbProvider,
	}
	e.registerBuiltinRecipes()
	return e
}

func (e *Engine) registerBuiltinRecipes() {
	e.Register(NewProjectRecipe())
	e.Register(NewCursorRulesRecipe())
	e.Register(NewClaudeCodeRecipe())
	e.Register(NewOptimizeRecipe())
}

func (e *Engine) Register(r Recipe) {
	e.builtinRecipes[r.Type()] = r
}

func (e *Engine) Build(ctx BuildContext) (*BuiltPrompt, error) {
	// Try DB recipe first (if available and active)
	if e.dbRecipe != nil {
		dbResult, err := e.dbRecipe.Build(ctx)
		if err == nil && dbResult != nil {
			return dbResult, nil
		}
	}

	// Fall back to builtin recipe
	r, ok := e.builtinRecipes[ctx.Type]
	if !ok {
		return nil, fmt.Errorf("no recipe found for type: %s", ctx.Type)
	}
	return r.Build(ctx)
}