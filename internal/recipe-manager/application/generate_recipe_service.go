package application

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/google/uuid"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
)

var (
	ErrInvalidGenerationMode   = errors.New("invalid generation mode")
	ErrPromptRequiredForMode   = errors.New("prompt is required when mode is prompt")
	ErrNilGeneratedRecipeDraft = errors.New("generated recipe draft is nil")
	ErrInvalidDraftAmounts     = errors.New("generated recipe contains invalid ingredient amounts")
	ErrInvalidDraftStructure   = errors.New("generated recipe draft has invalid structure")
	ErrInvalidBalancedRecipe   = errors.New("balanced recipe contains invalid numeric values")
)

type GenerateRecipeCommand struct {
	Mode   string
	Prompt string
	Pans   domain.Pans
}

func (c GenerateRecipeCommand) Validate() error {
	switch c.Mode {
	case "random":
	case "prompt":
		if c.Prompt == "" {
			return ErrPromptRequiredForMode
		}
	default:
		return ErrInvalidGenerationMode
	}

	return nil
}

type GeneratedRecipeDraft struct {
	Name        string
	Description string
	Author      string
	Dough       map[string]float64
	Topping     map[string]float64
}

func (d *GeneratedRecipeDraft) Validate() error {
	if d == nil {
		return ErrNilGeneratedRecipeDraft
	}
	if d.Name == "" {
		return errors.New("generated recipe name is required")
	}
	if len(d.Dough) == 0 {
		return errors.New("generated recipe dough is required")
	}
	if len(d.Topping) == 0 {
		return errors.New("generated recipe topping is required")
	}

	return nil
}

type RecipeGenerationClient interface {
	GenerateRecipe(context.Context, GenerateRecipeCommand) (*GeneratedRecipeDraft, error)
}

type GenerateRecipeService struct {
	generationClient RecipeGenerationClient
	calculator       CalculatorService
	balancer         BalancerService
}

func NewGenerateRecipeService(
	generationClient RecipeGenerationClient,
	calculator CalculatorService,
	balancer BalancerService,
) *GenerateRecipeService {
	return &GenerateRecipeService{
		generationClient: generationClient,
		calculator:       calculator,
		balancer:         balancer,
	}
}

func (s *GenerateRecipeService) Handle(ctx context.Context, command GenerateRecipeCommand) (*domain.RecipeAggregate, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	generatedDraft, err := s.generationClient.GenerateRecipe(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipe draft: %w", err)
	}
	if err := generatedDraft.Validate(); err != nil {
		return nil, err
	}
	if err := validateGeneratedDraftSemantics(*generatedDraft); err != nil {
		return nil, err
	}

	pans, err := s.calculator.TotalDoughWeightByPans(ctx, command.Pans)
	if err != nil {
		return nil, err
	}

	recipe := generatedDraft.toDomainRecipe()
	aggregate, err := s.balancer.Balance(ctx, recipe, *pans)
	if err != nil {
		return nil, err
	}
	if err := validateAggregateNumbers(*aggregate); err != nil {
		return nil, err
	}

	return aggregate, nil
}

func validateGeneratedDraftSemantics(draft GeneratedRecipeDraft) error {
	if !hasPositiveFiniteAmountCaseInsensitive(draft.Dough, "flour") || !hasPositiveFiniteAmountCaseInsensitive(draft.Dough, "water") {
		return fmt.Errorf("%w: dough must include positive finite 'flour' and 'water'", ErrInvalidDraftStructure)
	}
	if !allPositiveFinite(draft.Dough) || !allPositiveFinite(draft.Topping) {
		return fmt.Errorf("%w: every ingredient amount must be positive and finite", ErrInvalidDraftAmounts)
	}
	return nil
}

func hasPositiveFiniteAmountCaseInsensitive(ingredients map[string]float64, key string) bool {
	target := strings.ToLower(strings.TrimSpace(key))
	for k, v := range ingredients {
		if strings.ToLower(strings.TrimSpace(k)) == target {
			return isPositiveFinite(v)
		}
	}
	return false
}

func allPositiveFinite(ingredients map[string]float64) bool {
	for _, value := range ingredients {
		if !isPositiveFinite(value) {
			return false
		}
	}
	return true
}

func isPositiveFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value > 0
}

func validateAggregateNumbers(aggregate domain.RecipeAggregate) error {
	if !domainIngredientsFinite(aggregate.Recipe.Dough.Ingredients) || !domainIngredientsFinite(aggregate.Recipe.Topping.Ingredients) {
		return ErrInvalidBalancedRecipe
	}
	for _, dough := range aggregate.SplitIngredients.SplitDough {
		if !domainIngredientsFinite(dough.Ingredients) {
			return ErrInvalidBalancedRecipe
		}
	}
	for _, topping := range aggregate.SplitIngredients.SplitTopping {
		if !domainIngredientsFinite(topping.Ingredients) {
			return ErrInvalidBalancedRecipe
		}
	}
	return nil
}

func domainIngredientsFinite(ingredients []domain.Ingredient) bool {
	for _, ingredient := range ingredients {
		if math.IsNaN(ingredient.Amount) || math.IsInf(ingredient.Amount, 0) {
			return false
		}
	}
	return true
}

const defaultToppingReferenceArea = 1000.0

func (d GeneratedRecipeDraft) toDomainRecipe() domain.Recipe {
	return domain.Recipe{
		Uuid:        uuid.Nil,
		Name:        d.Name,
		Description: d.Description,
		Author:      d.Author,
		Dough: domain.Dough{
			Ingredients: toBakerPercentages(d.Dough),
		},
		Topping: domain.Topping{
			ReferenceArea: defaultToppingReferenceArea,
			Ingredients:   mapIngredients(d.Topping),
		},
	}
}

func toBakerPercentages(dough map[string]float64) []domain.Ingredient {
	flourAmount := findFlourAmountCaseInsensitive(dough)
	if flourAmount <= 0 {
		return mapIngredients(dough)
	}

	keys := make([]string, 0, len(dough))
	for name := range dough {
		keys = append(keys, name)
	}
	slices.Sort(keys)

	ingredients := make([]domain.Ingredient, 0, len(dough))
	for _, name := range keys {
		pct := (dough[name] / flourAmount) * 100
		ingredients = append(ingredients, domain.Ingredient{
			Name:   name,
			Amount: math.Round(pct*10) / 10,
		})
	}
	return ingredients
}

func findFlourAmountCaseInsensitive(ingredients map[string]float64) float64 {
	for k, v := range ingredients {
		if strings.ToLower(strings.TrimSpace(k)) == "flour" {
			return v
		}
	}
	return 0
}

func mapIngredients(input map[string]float64) []domain.Ingredient {
	keys := make([]string, 0, len(input))
	for name := range input {
		keys = append(keys, name)
	}
	slices.Sort(keys)

	ingredients := make([]domain.Ingredient, 0, len(input))
	for _, name := range keys {
		ingredients = append(ingredients, domain.Ingredient{
			Name:   name,
			Amount: input[name],
		})
	}

	return ingredients
}
