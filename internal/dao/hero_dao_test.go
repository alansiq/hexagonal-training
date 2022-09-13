package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mercadolibre/fury_cx-example/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHeroDAO_GetHero(t *testing.T) {
	t.Run("get hero", func(t *testing.T) {
		heroID := 123
		hero := models.HeroDto{
			ID:       123,
			Name:     "clark",
			Lastname: "kent",
			Age:      100,
			Level:    10,
			Type:     "human",
			ArmID:    1,
		}
		HeroBytes, _ := json.Marshal(hero)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%d", heroID) {
				t.Errorf("Expected to request '/123', got: %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected 'GET' method, got: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(HeroBytes)
		}))
		defer server.Close()

		s, _ := NewHeroDAO(server.Client(), server.URL)
		resp, err := s.GetHero(context.Background(), heroID)
		assert.NoError(t, err)
		assert.Equal(t, hero, *resp)
	})
}
