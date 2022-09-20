package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mercadolibre/fury_cx-example/internal/adapter/consumer/rest/dto"
	"net/http"

	"github.com/mercadolibre/fury_cx-example/internal/domain"
	"github.com/mercadolibre/fury_go-core/pkg/rusty"
)

type HeroClient struct {
	getHeroEndpoint      *rusty.Endpoint
	getAllHeroesEndpoint *rusty.Endpoint
	createHeroEndpoint   *rusty.Endpoint
}

func NewHeroClient(client HttpClient, uri string) (*HeroClient, error) {
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

	return &HeroClient{
		getHeroEndpoint:      getHeroEndpoint,
		getAllHeroesEndpoint: getAllHeroesEndpoint,
		createHeroEndpoint:   createHeroEndpoint,
	}, nil
}

func (cr *HeroClient) Get(ctx context.Context, heroID int) (*domain.Hero, error) {
	resp, err := cr.getHeroEndpoint.Get(ctx, rusty.WithParam("id", heroID))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("hero not found")
		}

		return nil, err
	}

	body := &dto.HeroDto{}
	if err := json.Unmarshal(resp.Body, body); err != nil {
		return nil, err
	}

	heroResponse := dto.DtoToHero(body)
	return heroResponse, nil
}

func (cr *HeroClient) Create(ctx context.Context, newHero *domain.Hero) error {
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

func (cr *HeroClient) ListAll(ctx context.Context) ([]domain.Hero, error) {
	resp, err := cr.getAllHeroesEndpoint.Get(ctx)
	if err != nil {
		return nil, err
	}

	var body []domain.HeroDto
	if err := json.Unmarshal(resp.Body, &body); err != nil {
		return nil, err
	}

	return body, nil
}
