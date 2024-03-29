package main

import (
	"log"

	"github.com/harisquqo/quqo-challenge-1/infrastructure/config"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/routes"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}

	config.LoadConfiguration()
}

func main() {
	p, err := base.NewPersistence()

	if err != nil {
		log.Fatal(err)
	}

	router := routes.InitRouter(p)

    router.Run(":8080")
}