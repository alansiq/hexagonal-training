package core

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_cx-example/internal/dao"
	"github.com/mercadolibre/fury_cx-example/internal/models"
	"github.com/mercadolibre/fury_cx-example/pkg/kvs"
)

type AppService struct {
	heroDAO   *dao.HeroDAO
	weaponDAO *dao.WeaponDAO
	stats     *kvs.Kvs
}

func NewAppService(heroDAO *dao.HeroDAO, weaponDAO *dao.WeaponDAO, stats *kvs.Kvs) *AppService {
	return &AppService{
		heroDAO:   heroDAO,
		weaponDAO: weaponDAO,
		stats:     stats,
	}
}

func (a *AppService) GetHero(ctx context.Context, heroID int) (*models.HeroDto, error) {
	return a.heroDAO.GetHero(ctx, heroID)
}

func (a *AppService) CreateHero(ctx context.Context, newHero *models.CreateHeroDto) error {
	err := a.heroDAO.CreateHero(ctx, newHero)
	if err != nil {
		return nil
	}

	stats, err := a.stats.Get(ctx, "stats")
	if err != nil {
		return err
	}

	if stats == nil {
		stats := models.Stats{}
		if newHero.Type == "human" {
			stats.Humans++
		} else {
			stats.Ogres++
		}
		a.stats.Save(ctx, "stats", stats)
	} else {
		if v, ok := stats.(models.Stats); ok {
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

func (a *AppService) CreateWeapon(ctx context.Context, newWeapon *models.CreateWeaponDTO) error {
	return a.weaponDAO.Create(ctx, newWeapon)
}

func (a *AppService) GetWeapon(ctx context.Context, weaponID int) (*models.WeaponDTO, error) {
	return a.weaponDAO.Get(ctx, weaponID)
}

func (a *AppService) Stats(ctx context.Context) (*models.Stats, error) {
	stats, err := a.stats.Get(ctx, "stats")
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return &models.Stats{}, nil
	}
	if v, ok := stats.(models.Stats); ok {
		return &v, nil
	}

	return nil, errors.New("error getting stats")
}
