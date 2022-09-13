package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mercadolibre/fury_cx-example/internal/core"
	"github.com/mercadolibre/fury_cx-example/internal/dao"
	"github.com/mercadolibre/fury_cx-example/internal/models"
	"github.com/mercadolibre/fury_go-core/pkg/web"
	"github.com/stretchr/testify/assert"
)

func Test_handler_HandleGetHero(t *testing.T) {

	t.Run("get hero", func(t *testing.T) {
		armID := 456
		expectedArm := models.ArmDTO{
			ID:   armID,
			Name: "knife",
		}
		heroID := 123
		expectedHero := models.HeroDto{
			ID:       123,
			Name:     "clark",
			Lastname: "kent",
			Age:      100,
			Level:    10,
			Type:     "human",
			ArmID:    armID,
			Arm:      &expectedArm,
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

		armBytes, _ := json.Marshal(expectedArm)
		armServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%d", armID) {
				t.Errorf("Expected to request '/456', got: %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected 'GET' method, got: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(armBytes)
		}))
		defer armServer.Close()

		armDAO, _ := dao.NewArmDAO(armServer.Client(), armServer.URL)

		heroDAO, _ := dao.NewHeroDAO(heroServer.Client(), heroServer.URL)
		srv := core.NewAppService(heroDAO, armDAO, nil)

		//handler := NewHandler(srv)

		req := httptest.NewRequest(http.MethodGet, "/hero/123", nil)
		resp := httptest.NewRecorder()

		handler := web.New()
		handler.Get("/hero/{id}", NewHandler(srv).HandleGetHero)
		handler.ServeHTTP(resp, req)

		assert.NotNil(t, resp.Result())
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		heroResult := models.HeroDto{}
		json.Unmarshal(resp.Body.Bytes(), &heroResult)
		assert.Equal(t, expectedHero, heroResult)
	})
}
