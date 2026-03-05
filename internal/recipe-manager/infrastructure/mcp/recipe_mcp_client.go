package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
)

type RecipeMCPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRecipeMCPClient(baseURL string, timeout time.Duration) *RecipeMCPClient {
	return &RecipeMCPClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

type generateRecipeRequest struct {
	Mode           string `json:"mode"`
	Prompt         string `json:"prompt,omitempty"`
	OutputContract struct {
		RequiredDoughIngredients []string `json:"requiredDoughIngredients,omitempty"`
	} `json:"outputContract"`
}

type generateRecipeResponse struct {
	Data struct {
		RecipeDraft struct {
			Name        string             `json:"name"`
			Description string             `json:"description"`
			Author      string             `json:"author"`
			Dough       map[string]float64 `json:"dough"`
			Topping     map[string]float64 `json:"topping"`
		} `json:"recipeDraft"`
	} `json:"data"`
}

func (c *RecipeMCPClient) GenerateRecipe(ctx context.Context, command application.GenerateRecipeCommand) (*application.GeneratedRecipeDraft, error) {
	requestBody := generateRecipeRequest{
		Mode:   command.Mode,
		Prompt: command.Prompt,
	}
	requestBody.OutputContract.RequiredDoughIngredients = []string{"flour", "water"}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/mcp/tools/generate_recipe", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("mcp generate_recipe failed: %s - %s", response.Status, string(responseBody))
	}

	var parsed generateRecipeResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, err
	}

	return &application.GeneratedRecipeDraft{
		Name:        parsed.Data.RecipeDraft.Name,
		Description: parsed.Data.RecipeDraft.Description,
		Author:      parsed.Data.RecipeDraft.Author,
		Dough:       parsed.Data.RecipeDraft.Dough,
		Topping:     parsed.Data.RecipeDraft.Topping,
	}, nil
}
