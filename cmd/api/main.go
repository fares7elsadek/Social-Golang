package main

import (
	"log"

	"github.com/fares7elsadek/Social-Golang/internal/env"
	"github.com/fares7elsadek/Social-Golang/internal/repository/postgres"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading enviroment varibles")
	}

	db,err := postgres.Connect(env.GetString("DB_URL", "postgres://postgres:password@localhost:5432/social_golang"))
	if err != nil {
		log.Fatal("Error connecting to database")
	}else{
		log.Println("Connected to database successfully")
	}

	

	cfg := config{
		addr: env.GetString("PORT", "5000"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount(db)
	log.Fatal(app.run(mux))
}
