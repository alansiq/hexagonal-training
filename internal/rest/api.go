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

	arm, err := h.appService.GetArm(r.Context(), hero.ArmID)
	if err != nil {
		return err
	}

	hero.Arm = arm

	return web.EncodeJSON(w, hero, http.StatusOK)
}

func (h *handler) HandleCreateHero(w http.ResponseWriter, r *http.Request) error {
	payload := &models.CreateHeroDto{}
	if err := web.DecodeJSON(r, payload); err != nil {
		return err
	}

	arm, err := h.appService.GetArm(r.Context(), payload.ArmID)
	if err != nil {
		return err
	}

	if arm == nil {
		return errors.New("arm doesnt exists")
	}

	err = h.appService.CreateHero(r.Context(), payload)
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, "hero crated", http.StatusOK)
}

func (h *handler) HandleCreateArm(w http.ResponseWriter, r *http.Request) error {
	payload := &models.CreateArmDTO{}
	if err := web.DecodeJSON(r, payload); err != nil {
		return err
	}

	err := h.appService.CreateArm(r.Context(), payload)
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, "arm crated", http.StatusOK)
}

func (h *handler) HandleGetArm(w http.ResponseWriter, r *http.Request) error {
	armID := chi.URLParam(r, "id")
	if armID == "" {
		return errors.New("arm id required")
	}

	armInt, err := strconv.Atoi(armID)
	if err != nil {
		return err
	}

	arm, err := h.appService.GetArm(r.Context(), armInt)
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, arm, http.StatusOK)
}

func (h *handler) HandleStats(w http.ResponseWriter, r *http.Request) error {
	stats, err := h.appService.Stats(r.Context())
	if err != nil {
		return err
	}

	return web.EncodeJSON(w, stats, http.StatusOK)
}
