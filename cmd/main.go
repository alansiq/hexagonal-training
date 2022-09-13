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

	armServer := armMockServer()
	defer armServer.Close()
	armDAO, err := dao.NewArmDAO(armServer.Client(), armServer.URL)
	if err != nil {
		return err
	}

	appService := core.NewAppService(heroDAO, armDAO, kvs.NewKvs())
	handler := rest.NewHandler(appService)
	app.Router.Get("/hero/{id}", handler.HandleGetHero)
	app.Router.Post("/hero", handler.HandleCreateHero)
	app.Router.Get("/arm/{id}", handler.HandleGetArm)
	app.Router.Post("/arm", handler.HandleCreateArm)
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
		ArmID:    111,
	}
	HeroBytes, _ := json.Marshal(hero)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(HeroBytes)
	}))
}

func armMockServer() *httptest.Server {
	armID := 111
	arm := models.ArmDTO{
		ID:   armID,
		Name: "knife",
	}
	armBytes, _ := json.Marshal(arm)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(armBytes)
	}))
}
