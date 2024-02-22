package main

import (
	"log"

	distributefilesystem "github.com/enteprise/etl-central/app/mainapp/code_editor/distribute_filesystem"
	dpconfig "github.com/enteprise/etl-central/app/mainapp/config"
	"github.com/enteprise/etl-central/app/mainapp/database"
)

func main() {
	dpconfig.LoadConfig()
	database.DBConnect()
	log.Println("🏃 Running")
	// CreateFiles()
	distributefilesystem.MoveCodeFilesToDB(database.DBConn)
}
