package main

import (
	"log"

	"github.com/fares7elsadek/Social-Golang/internal/env"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading enviroment varibles")
	}

	cfg := config{
		addr: env.GetString("PORT", "5000"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
