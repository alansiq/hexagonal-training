package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/mercadolibre/fury_cx-example/internal/core"
	"github.com/mercadolibre/fury_cx-example/internal/models"
	"github.com/mercadolibre/fury_go-core/pkg/web"
)

type handler struct {
	validator  *validator.Validate
	appService *core.AppService
}

func NewHandler(app *core.AppService) *handler {
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

	weapon, err := h.appService.GetWeapon(r.Context(), hero.WeaponID)
	if err != nil {
		return err
	}

	hero.Weapon = weapon

	return web.EncodeJSON(w, hero, http.StatusOK)
}

func (h *handler) HandleCreateHero(w http.ResponseWriter, r *http.Request) error {
	payload := &models.CreateHeroDto{}
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

func (h *handler) HandleCreateWeapon(w http.ResponseWriter, r *http.Request) error {
	payload := &models.CreateWeaponDTO{}
	if err := web.DecodeJSON(r, payload); err != nil {
		return err
	}

	err := h.appService.CreateWeapon(r.Context(), payload)
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, "weapon crated", http.StatusOK)
}

func (h *handler) HandleGetWeapon(w http.ResponseWriter, r *http.Request) error {
	weaponID := chi.URLParam(r, "id")
	if weaponID == "" {
		return errors.New("weapon id required")
	}

	weaponInt, err := strconv.Atoi(weaponID)
	if err != nil {
		return err
	}

	weapon, err := h.appService.GetWeapon(r.Context(), weaponInt)
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, weapon, http.StatusOK)
}

func (h *handler) HandleStats(w http.ResponseWriter, r *http.Request) error {
	stats, err := h.appService.Stats(r.Context())
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, stats, http.StatusOK)
}
