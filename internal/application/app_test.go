package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mercadolibre/fury_cx-example/internal/adapter/consumer/rest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mercadolibre/fury_cx-example/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAppService_GetHero(t *testing.T) {

	t.Run("get hero", func(t *testing.T) {
		heroID := 123
		hero := domain.HeroDto{
			ID:       123,
			Name:     "clark",
			Lastname: "kent",
			Age:      100,
			Level:    10,
			Type:     "human",
			WeaponID: 1,
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

		heroDAO, _ := rest.NewHeroClient(server.Client(), server.URL)
		srv := NewAppService(heroDAO, nil, nil)
		resp, err := srv.GetHero(context.Background(), heroID)
		assert.NoError(t, err)
		assert.Equal(t, hero, *resp)
	})
}
