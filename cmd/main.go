package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/mercadolibre/fury_cx-example/internal/core"
	"github.com/mercadolibre/fury_cx-example/internal/dao"
	"github.com/mercadolibre/fury_cx-example/internal/models"
	"github.com/mercadolibre/fury_cx-example/internal/rest"
	"github.com/mercadolibre/fury_cx-example/pkg/kvs"
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/mercadolibre/fury_go-toolkit-otel/pkg/otel"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error in run: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	shutdown, err := otel.Start(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdown()

	app, err := fury.NewWebApplication()
	if err != nil {
		return err
	}

	heroServer := heroMockServer()
	defer heroServer.Close()
	heroDAO, err := dao.NewHeroDAO(heroServer.Client(), heroServer.URL)
	if err != nil {
		return err
	}

	weaponServer := weaponMockServer()
	defer weaponServer.Close()
	weaponDAO, err := dao.NewWeaponDAO(weaponServer.Client(), weaponServer.URL)
	if err != nil {
		return err
	}

	appService := core.NewAppService(heroDAO, weaponDAO, kvs.NewKvs("dummy"))
	handler := rest.NewHandler(appService)
	app.Router.Get("/hero/{id}", handler.HandleGetHero)
	app.Router.Post("/hero", handler.HandleCreateHero)
	app.Router.Get("/weapon/{id}", handler.HandleGetWeapon)
	app.Router.Post("/weapon", handler.HandleCreateWeapon)
	app.Router.Get("/stats", handler.HandleStats)

	return app.Run()
}

func heroMockServer() *httptest.Server {
	heroID := 123
	hero := models.HeroDto{
		ID:       heroID,
		Name:     "clark",
		Lastname: "kent",
		Age:      100,
		Level:    10,
		Type:     "human",
		WeaponID: 111,
	}
	HeroBytes, _ := json.Marshal(hero)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(HeroBytes)
	}))
}

func weaponMockServer() *httptest.Server {
	weaponID := 111
	weapon := models.WeaponDTO{
		ID:   weaponID,
		Name: "knife",
	}
	weaponBytes, _ := json.Marshal(weapon)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(weaponBytes)
	}))
}
