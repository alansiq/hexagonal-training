package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/mercadolibre/fury_go-core/pkg/web"
	"net/http"
	"strconv"
)

func (h *handler) HandleCreateWeapon(w http.ResponseWriter, r *http.Request) error {
	payload := &domain.CreateWeaponDTO{}
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
