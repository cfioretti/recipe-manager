package mcp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cfioretti/recipe-manager/internal/recipe-manager/application"
)

func TestRecipeMCPClient_GenerateRecipe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/mcp/tools/generate_recipe", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			var payload map[string]any
			assert.NoError(t, json.Unmarshal(body, &payload))
			assert.Equal(t, "prompt", payload["mode"])
			outputContract, ok := payload["outputContract"].(map[string]any)
			assert.True(t, ok)
			requiredDoughIngredients, ok := outputContract["requiredDoughIngredients"].([]any)
			assert.True(t, ok)
			assert.Equal(t, []any{"flour", "water"}, requiredDoughIngredients)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"data": {
					"recipeDraft": {
						"name": "AI Pizza",
						"description": "generated",
						"author": "mcp",
						"dough": {"flour": 1000, "water": 700},
						"topping": {"mozzarella": 300}
					}
				}
			}`))
		}))
		defer server.Close()

		client := NewRecipeMCPClient(server.URL, 2*time.Second)
		result, err := client.GenerateRecipe(context.Background(), application.GenerateRecipeCommand{
			Mode:   "prompt",
			Prompt: "high hydration",
		})

		assert.NoError(t, err)
		assert.Equal(t, "AI Pizza", result.Name)
		assert.Equal(t, 1000.0, result.Dough["flour"])
	})

	t.Run("error status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "bad request", http.StatusBadRequest)
		}))
		defer server.Close()

		client := NewRecipeMCPClient(server.URL, 2*time.Second)
		result, err := client.GenerateRecipe(context.Background(), application.GenerateRecipeCommand{
			Mode: "random",
		})

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}
