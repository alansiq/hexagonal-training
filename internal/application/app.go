package application

import (
	"context"
	"errors"
	"github.com/mercadolibre/fury_cx-example/internal/adapter/consumer/rest/dto"
	"github.com/mercadolibre/fury_cx-example/internal/domain"
)

type AppService struct {
	heroClient   HeroClient
	weaponClient WeaponClient
	stats        RepositoryClient
}

func NewAppService(heroClient HeroClient, weaponClient WeaponClient, stats RepositoryClient) *AppService {
	return &AppService{
		heroClient:   heroClient,
		weaponClient: weaponClient,
		stats:        stats,
	}
}

func (a *AppService) GetHero(ctx context.Context, heroID int) (*domain.Hero, error) {
	hero, err := a.heroClient.Get(ctx, heroID)
	if err != nil {
		return nil, err
	}

	weapon, err := a.weaponClient.Get(ctx, hero.WeaponId)
	if err != nil {
		return nil, err
	}

	hero.Weapon = weapon

	return hero, nil
}

func (a *AppService) CreateHero(ctx context.Context, newHero *dto.CreateHeroDto) error {
	err := a.heroClient.Create(ctx, newHero)
	if err != nil {
		return nil
	}

	stats, err := a.stats.Get(ctx, "stats")
	if err != nil {
		return err
	}

	if stats == nil {
		stats := domain.Stats{}
		if newHero.Type == "human" {
			stats.Humans++
		} else {
			stats.Ogres++
		}
		a.stats.Save(ctx, "stats", stats)
	} else {
		if v, ok := stats.(domain.Stats); ok {
			if newHero.Type == "human" {
				v.Humans++
			} else {
				v.Ogres++
			}
			a.stats.Save(ctx, "stats", v)
		}
	}

	return nil
}

func (a *AppService) CreateWeapon(ctx context.Context, newWeapon *domain.Weapon) error {
	return a.weaponClient.Create(ctx, newWeapon)
}

func (a *AppService) GetWeapon(ctx context.Context, weaponID int) (*domain.Weapon, error) {
	return a.weaponClient.Get(ctx, weaponID)
}

func (a *AppService) Stats(ctx context.Context) (*domain.Stats, error) {
	stats, err := a.stats.Get(ctx, "stats")
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return &domain.Stats{}, nil
	}
	if v, ok := stats.(domain.Stats); ok {
		return &v, nil
	}

	return nil, errors.New("error getting stats")
}
