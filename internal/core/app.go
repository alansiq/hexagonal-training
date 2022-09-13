package core

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_cx-example/internal/dao"
	"github.com/mercadolibre/fury_cx-example/internal/models"
	"github.com/mercadolibre/fury_cx-example/pkg/kvs"
)

type AppService struct {
	heroDAO *dao.HeroDAO
	armDAO  *dao.ArmDAO
	stats   *kvs.Kvs
}

func NewAppService(heroDAO *dao.HeroDAO, armDAO *dao.ArmDAO, stats *kvs.Kvs) *AppService {
	return &AppService{
		heroDAO: heroDAO,
		armDAO:  armDAO,
		stats:   stats,
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

func (a *AppService) CreateArm(ctx context.Context, newArm *models.CreateArmDTO) error {
	return a.armDAO.Create(ctx, newArm)
}

func (a *AppService) GetArm(ctx context.Context, armID int) (*models.ArmDTO, error) {
	return a.armDAO.Get(ctx, armID)
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