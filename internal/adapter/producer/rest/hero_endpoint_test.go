package rest

import (
	"encoding/json"
	"fmt"
	rest2 "github.com/mercadolibre/fury_cx-example/internal/adapter/consumer/rest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mercadolibre/fury_cx-example/internal/application"
	"github.com/mercadolibre/fury_cx-example/internal/domain"
	"github.com/mercadolibre/fury_go-core/pkg/web"
	"github.com/stretchr/testify/assert"
)

func Test_handler_HandleGetHero(t *testing.T) {

	t.Run("get hero", func(t *testing.T) {
		weaponID := 456
		expectedWeapon := domain.WeaponDTO{
			ID:   weaponID,
			Name: "knife",
		}
		heroID := 123
		expectedHero := domain.HeroDto{
			ID:       123,
			Name:     "clark",
			Lastname: "kent",
			Age:      100,
			Level:    10,
			Type:     "human",
			WeaponID: weaponID,
			Weapon:   &expectedWeapon,
		}

		HeroBytes, _ := json.Marshal(expectedHero)
		heroServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%d", heroID) {
				t.Errorf("Expected to request '/123', got: %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected 'GET' method, got: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(HeroBytes)
		}))
		defer heroServer.Close()

		weaponBytes, _ := json.Marshal(expectedWeapon)
		weaponServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%d", weaponID) {
				t.Errorf("Expected to request '/456', got: %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected 'GET' method, got: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(weaponBytes)
		}))
		defer weaponServer.Close()

		weaponDAO, _ := rest2.NewWeaponDAO(weaponServer.Client(), weaponServer.URL)

		heroDAO, _ := rest2.NewHeroClient(heroServer.Client(), heroServer.URL)
		srv := application.NewAppService(heroDAO, weaponDAO, nil)

		//handler := NewHandler(srv)

		req := httptest.NewRequest(http.MethodGet, "/hero/123", nil)
		resp := httptest.NewRecorder()

		handler := web.New()
		handler.Get("/hero/{id}", NewHandler(srv).HandleGetHero)
		handler.ServeHTTP(resp, req)

		assert.NotNil(t, resp.Result())
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		heroResult := domain.HeroDto{}
		json.Unmarshal(resp.Body.Bytes(), &heroResult)
		assert.Equal(t, expectedHero, heroResult)
	})
}
