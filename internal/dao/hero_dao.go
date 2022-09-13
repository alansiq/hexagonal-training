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

type HeroDAO struct {
	getHeroEndpoint      *rusty.Endpoint
	getAllHeroesEndpoint *rusty.Endpoint
	createHeroEndpoint   *rusty.Endpoint
}

func NewHeroDAO(client HttpClient, uri string) (*HeroDAO, error) {
	getHeroEndpoint, err := rusty.NewEndpoint(client, rusty.URL(uri, "/{id}"), rusty.WithHeader("X-Admin-Id", "APP_CORE"))
	if err != nil {
		return nil, fmt.Errorf("error creating get hero endpoint, %w", err)
	}

	getAllHeroesEndpoint, err := rusty.NewEndpoint(client, rusty.URL(uri, "/get_all"), rusty.WithHeader("X-Admin-Id", "APP_CORE"))
	if err != nil {
		return nil, fmt.Errorf("error creating get all heroes endpoint, %w", err)
	}

	createHeroEndpoint, err := rusty.NewEndpoint(client, rusty.URL(uri, "/"), rusty.WithHeader("X-Admin-Id", "APP_CORE"))
	if err != nil {
		return nil, fmt.Errorf("error creating create hero endpoint, %w", err)
	}

	return &HeroDAO{
		getHeroEndpoint:      getHeroEndpoint,
		getAllHeroesEndpoint: getAllHeroesEndpoint,
		createHeroEndpoint:   createHeroEndpoint,
	}, nil
}

func (cr *HeroDAO) GetHero(ctx context.Context, heroID int) (*models.HeroDto, error) {
	resp, err := cr.getHeroEndpoint.Get(ctx, rusty.WithParam("id", heroID))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("hero not found")
		}

		return nil, err
	}

	body := &models.HeroDto{}
	if err := json.Unmarshal(resp.Body, body); err != nil {
		return nil, err
	}

	return body, nil
}

func (cr *HeroDAO) CreateHero(ctx context.Context, newHero *models.CreateHeroDto) error {
	body, err := json.Marshal(newHero)
	if err != nil {
		return errors.New("error in new hero marshal")
	}

	_, err = cr.createHeroEndpoint.Post(ctx, rusty.WithBody(body))
	if err != nil {
		return err
	}

	return nil
}

func (cr *HeroDAO) GetHeroes(ctx context.Context) ([]models.HeroDto, error) {
	resp, err := cr.getAllHeroesEndpoint.Get(ctx)
	if err != nil {
		return nil, err
	}

	var body []models.HeroDto
	if err := json.Unmarshal(resp.Body, &body); err != nil {
		return nil, err
	}

	return body, nil
}
