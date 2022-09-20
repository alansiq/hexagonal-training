package application

import (
	"context"
	"github.com/mercadolibre/fury_cx-example/internal/domain"
)

type RepositoryId string

type HeroClient interface {
	Get(ctx context.Context, heroId int) (*domain.Hero, error)
	Create(ctx context.Context, newHero *domain.CreateHeroDto) error
}

type WeaponClient interface {
	Get(ctx context.Context, weaponID int) (*domain.WeaponDTO, error)
	Create(ctx context.Context, newHero *domain.CreateWeaponDTO) error
}

type RepositoryClient interface {
	Get(ctx context.Context, key RepositoryId) (interface{}, error)
	Save(ctx context.Context, key RepositoryId, newObject interface{}) error
}
