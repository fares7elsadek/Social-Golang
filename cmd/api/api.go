package main

import (
	"log"
	"net/http"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func(app *application) mount() http.Handler {
	
	r := chi.NewRouter()
	// base middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) 
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	

	r.Route("/v1",func (r chi.Router){
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}

func (app *application) run(handler http.Handler) error {

	srv := &http.Server{
		Addr: app.config.addr,
		Handler: handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server running on port %v",app.config.addr)

	return srv.ListenAndServe()
}