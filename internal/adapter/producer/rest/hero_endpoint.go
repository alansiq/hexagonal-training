package rest

import (
	"errors"
	"github.com/mercadolibre/fury_cx-example/internal/adapter/producer/rest/dto"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_cx-example/internal/application"
	"github.com/mercadolibre/fury_cx-example/internal/domain"
	"github.com/mercadolibre/fury_go-core/pkg/web"
)

type handler struct {
	validator  *validator.Validate
	appService *application.AppService
}

func NewHandler(app *application.AppService) *handler {
	return &handler{appService: app}
}

func (h *handler) HandleGetHero(w http.ResponseWriter, r *http.Request) error {
	heroID := chi.URLParam(r, "id")
	if heroID == "" {
		return errors.New("hero id required")
	}

	heroInt, err := strconv.Atoi(heroID)
	if err != nil {
		return err
	}

	hero, err := h.appService.GetHero(r.Context(), heroInt)
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, dto.HeroToDummy(hero), http.StatusOK)
}

func (h *handler) HandleCreateHero(w http.ResponseWriter, r *http.Request) error {
	payload := &domain.CreateHeroDto{}
	if err := web.DecodeJSON(r, payload); err != nil {
		return err
	}

	weapon, err := h.appService.GetWeapon(r.Context(), payload.WeaponID)
	if err != nil {
		return err
	}

	if weapon == nil {
		return errors.New("weapon doesnt exists")
	}

	err = h.appService.CreateHero(r.Context(), payload)
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, "hero crated", http.StatusOK)
}
