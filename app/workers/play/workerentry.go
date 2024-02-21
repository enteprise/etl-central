package main

import (
	"log"

	"github.com/enteprise/etl-central/app/mainapp/database"
	wrkerconfig "github.com/enteprise/etl-central/app/workers/config"
)

func main() {
	wrkerconfig.LoadConfig()
	database.DBConnect()
	log.Println("🏃 Running")
	// CreateFiles()
	// distfilesystem.DownloadFiles()
}
