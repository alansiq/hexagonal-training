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

func TestArmDAO_Get(t *testing.T) {
	t.Run("get arm", func(t *testing.T) {
		armID := 123
		arm := models.ArmDTO{
			ID:   armID,
			Name: "knife",
		}
		armBytes, _ := json.Marshal(arm)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%d", armID) {
				t.Errorf("Expected to request '/123', got: %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected 'GET' method, got: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(armBytes)
		}))
		defer server.Close()

		s, _ := NewArmDAO(server.Client(), server.URL)
		resp, err := s.Get(context.Background(), armID)
		assert.NoError(t, err)
		assert.Equal(t, arm, *resp)
	})

	t.Run("get arm alternative", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		m := mocks.NewMockHttpClient(mockCtrl)
		armID := 123
		arm := models.ArmDTO{
			ID:   armID,
			Name: "knife",
		}

		armBytes, _ := json.Marshal(arm)
		r := ioutil.NopCloser(bytes.NewReader(armBytes))
		m.
			EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				StatusCode: 201,
				Body:       r,
				Request:    &http.Request{},
			}, nil).
			Times(1)

		s, _ := NewArmDAO(m, "http://localhost:8000")
		resp, err := s.Get(context.Background(), armID)
		assert.NoError(t, err)
		assert.Equal(t, arm, *resp)
	})
}
