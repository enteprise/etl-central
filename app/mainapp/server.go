package main

import (
	"log"
	"os"

	"github.com/enteprise/etl-central/app/mainapp/routes"
	_ "go.uber.org/automaxprocs"
)

func main() {

	port := os.Getenv("DP_PORT")
	if port == "" {
		port = "9000"
	}

	app := routes.Setup(port)

	log.Fatal(app.Listen("0.0.0.0:" + port))
}
