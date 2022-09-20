package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mercadolibre/fury_cx-example/internal/domain"
	"github.com/mercadolibre/fury_go-core/pkg/rusty"
)

type WeaponDAO struct {
	getWeaponEndpoint    *rusty.Endpoint
	createWeaponEndpoint *rusty.Endpoint
}

func NewWeaponDAO(client HttpClient, uri string) (*WeaponDAO, error) {
	getWeaponEndpoint, err := rusty.NewEndpoint(client, rusty.URL(uri, "/{id}"), rusty.WithHeader("X-Admin-Id", "APP_CORE"))
	if err != nil {
		return nil, fmt.Errorf("error creating get weapon endpoint, %w", err)
	}

	createWeaponEndpoint, err := rusty.NewEndpoint(client, rusty.URL(uri, "/"), rusty.WithHeader("X-Admin-Id", "APP_CORE"))
	if err != nil {
		return nil, fmt.Errorf("error creating create weapon endpoint, %w", err)
	}

	return &WeaponDAO{
		getWeaponEndpoint:    getWeaponEndpoint,
		createWeaponEndpoint: createWeaponEndpoint,
	}, nil
}

func (cr *WeaponDAO) Get(ctx context.Context, weaponID int) (*domain.WeaponDTO, error) {
	resp, err := cr.getWeaponEndpoint.Get(ctx, rusty.WithParam("id", weaponID))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("weapon not found")
		}

		return nil, err
	}

	body := &domain.WeaponDTO{}
	if err := json.Unmarshal(resp.Body, body); err != nil {
		return nil, err
	}

	return body, nil
}

func (cr *WeaponDAO) Create(ctx context.Context, newHero *domain.CreateWeaponDTO) error {
	body, err := json.Marshal(newHero)
	if err != nil {
		return errors.New("error in new weapon marshal")
	}

	_, err = cr.createWeaponEndpoint.Post(ctx, rusty.WithBody(body))
	if err != nil {
		return err
	}

	return nil
}
