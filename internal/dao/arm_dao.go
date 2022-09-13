package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mercadolibre/fury_cx-example/internal/models"
	"github.com/mercadolibre/fury_go-core/pkg/rusty"
)

type ArmDAO struct {
	getArmEndpoint    *rusty.Endpoint
	createArmEndpoint *rusty.Endpoint
}

func NewArmDAO(client HttpClient, uri string) (*ArmDAO, error) {
	getArmEndpoint, err := rusty.NewEndpoint(client, rusty.URL(uri, "/{id}"), rusty.WithHeader("X-Admin-Id", "APP_CORE"))
	if err != nil {
		return nil, fmt.Errorf("error creating get arm endpoint, %w", err)
	}

	createArmEndpoint, err := rusty.NewEndpoint(client, rusty.URL(uri, "/"), rusty.WithHeader("X-Admin-Id", "APP_CORE"))
	if err != nil {
		return nil, fmt.Errorf("error creating create arm endpoint, %w", err)
	}

	return &ArmDAO{
		getArmEndpoint:    getArmEndpoint,
		createArmEndpoint: createArmEndpoint,
	}, nil
}

func (cr *ArmDAO) Get(ctx context.Context, armID int) (*models.ArmDTO, error) {
	resp, err := cr.getArmEndpoint.Get(ctx, rusty.WithParam("id", armID))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("arm not found")
		}

		return nil, err
	}

	body := &models.ArmDTO{}
	if err := json.Unmarshal(resp.Body, body); err != nil {
		return nil, err
	}

	return body, nil
}

func (cr *ArmDAO) Create(ctx context.Context, newHero *models.CreateArmDTO) error {
	body, err := json.Marshal(newHero)
	if err != nil {
		return errors.New("error in new arm marshal")
	}

	_, err = cr.createArmEndpoint.Post(ctx, rusty.WithBody(body))
	if err != nil {
		return err
	}

	return nil
}
