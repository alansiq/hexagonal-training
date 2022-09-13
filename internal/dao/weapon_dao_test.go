package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mercadolibre/fury_cx-example/internal/models"
	"github.com/mercadolibre/fury_cx-example/test/bugs/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWeaponDAO_Get(t *testing.T) {
	t.Run("get weapon", func(t *testing.T) {
		weaponID := 123
		weapon := models.WeaponDTO{
			ID:   weaponID,
			Name: "knife",
		}
		weaponBytes, _ := json.Marshal(weapon)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%d", weaponID) {
				t.Errorf("Expected to request '/123', got: %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected 'GET' method, got: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(weaponBytes)
		}))
		defer server.Close()

		s, _ := NewWeaponDAO(server.Client(), server.URL)
		resp, err := s.Get(context.Background(), weaponID)
		assert.NoError(t, err)
		assert.Equal(t, weapon, *resp)
	})

	t.Run("get weapon alternative", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		m := mocks.NewMockHttpClient(mockCtrl)
		weaponID := 123
		weapon := models.WeaponDTO{
			ID:   weaponID,
			Name: "knife",
		}

		weaponBytes, _ := json.Marshal(weapon)
		r := ioutil.NopCloser(bytes.NewReader(weaponBytes))
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				StatusCode: 201,
				Body:       r,
				Request:    &http.Request{},
			}, nil).
			Times(1)

		s, _ := NewWeaponDAO(m, "http://localhost:8000")
		resp, err := s.Get(context.Background(), weaponID)
		assert.NoError(t, err)
		assert.Equal(t, weapon, *resp)
	})
}
