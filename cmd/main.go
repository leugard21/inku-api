package main

import (
	"log"

	"github.com/leugard21/inku-api/cmd/api"
	"github.com/leugard21/inku-api/configs"
	"github.com/leugard21/inku-api/db"
)

func main() {
	storage, err := db.NewPostgresStorage(configs.Envs)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":"+configs.Envs.Port, storage)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
