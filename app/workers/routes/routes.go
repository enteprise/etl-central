package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/enteprise/etl-central/app/mainapp/database"
	modelmain "github.com/enteprise/etl-central/app/mainapp/database/models"
	"github.com/enteprise/etl-central/app/mainapp/messageq"
	"github.com/enteprise/etl-central/app/mainapp/utilities"
	wrkerconfig "github.com/enteprise/etl-central/app/workers/config"
	runcodeworker "github.com/enteprise/etl-central/app/workers/runcode"
	"github.com/enteprise/etl-central/app/workers/runtask"
	"github.com/enteprise/etl-central/app/workers/secrets"
	"github.com/enteprise/etl-central/app/workers/workerhealth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
)

func Setup(port string) *fiber.App {

	app := fiber.New()

	start := time.Now()

	wrkerconfig.LoadConfig()

	// ------- DATABASE CONNECT ------
	database.RedisConnect()
	database.DBConnect()
	log.Println("🏃 ======== DATAPLANE WORKER ========")
	log.Println("🔢 Dataplane version:", wrkerconfig.Version)

	// -------- NATS Connect -------
	messageq.NATSConnect()

	// ------ Validate worker data ---------
	if wrkerconfig.WorkerGroup == "" {
		panic("Requires worker_group environment variable")
	}

	// Validate group name
	var isStringAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	if !isStringAlphaNumeric(wrkerconfig.WorkerGroup) {
		panic("Worker group - Only [a-z], [A-Z], [0-9] and _ are allowed")
	}

	if wrkerconfig.WorkerType == "" {
		panic("Requires worker_type environment variable")
	}

	if wrkerconfig.WorkerEnv == "" {
		panic("Requires worker_env environment variable")
	}

	// ----- Load platformID ------
	for i := 0; i < 50000; i++ {
		platform := modelmain.Platform{}
		database.DBConn.First(&platform)
		wrkerconfig.PlatformID = platform.ID

		if wrkerconfig.PlatformID != "" {

			/* Load encryption key */
			if os.Getenv("secret_encryption_key") == "" {
				utilities.Encryptphrase = platform.EncryptKey
			} else {
				utilities.Encryptphrase = os.Getenv("secret_encryption_key")
			}

			break
		} else {
			log.Printf("😩 Platform not setup - waiting for main app to start: try number. %d, retry in 5 seconds", i+1)
			time.Sleep(time.Second * 5)
		}

	}

	log.Println("🎯 Platform ID: ", wrkerconfig.PlatformID)

	e := modelmain.Environment{}
	database.DBConn.First(&e, "name = ?", wrkerconfig.WorkerEnv)

	// if e.Name != wrkerconfig.WorkerEnv  {
	// 	panic("Warning: Envrionment not found. Be sure environment is setup with mainapp - " + wrkerconfig.WorkerEnv)
	// }

	// For first time users create a development environment ID
	if e.ID == "" && wrkerconfig.WorkerEnv == "Development" {
		e = modelmain.Environment{
			ID:         uuid.New().String(),
			Name:       "Development",
			PlatformID: wrkerconfig.PlatformID,
			Active:     true,
		}

		err := database.DBConn.Create(&e).Error

		if err != nil {
			if wrkerconfig.Debug == "true" {
				log.Println(err)
			}
			panic("Failed to create a development environment on first use.")
		}

	}

	wrkerconfig.EnvName = e.Name
	wrkerconfig.EnvID = e.ID
	log.Println("🌳 Environment name and ID: ", wrkerconfig.EnvName, " - ", wrkerconfig.EnvID)

	// ------- LOAD secrets (must be loaded after environment id to load for this environment) ------
	secrets.MapSecrets()

	// Load a worker ID
	wrkerconfig.WorkerID = uuid.NewString()
	log.Println("👷 Worker Group and ID: ", wrkerconfig.WorkerGroup, " - ", wrkerconfig.WorkerID)

	//recover from panic
	app.Use(recover.New())

	// add timer field to response header
	app.Use(Timer())

	// wrkerconfig.Scheduler = gocron.NewScheduler(time.UTC)
	// wrkerconfig.Scheduler.StartAsync()

	if wrkerconfig.Debug == "true" {
		app.Use(logger.New(
			logger.Config{
				Format: "✨ Latency: ${latency} Time:${time} Status: ${status} Path:${path} \n",
			}))
		// Method:${method} -- bug in fiber, waiting for pull request
		// UA:${ua}
		// Host:${host}
		// Header:${header}
		// Query:${query}
	}

	// ------- LOAD DS Files ------

	// ------- LOAD packages ------
	runcodeworker.CodeLoadPackages(wrkerconfig.CodeLanguages, wrkerconfig.CodeLoadPackages, wrkerconfig.EnvID, wrkerconfig.WorkerGroup)

	// Runner
	// app.Post("/runner", runtask.Runtask())

	// Cancel running job
	app.Post("/runnercancel/:id", runtask.Canceltask())

	// ----- APP Version ----
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"Version": wrkerconfig.Version})
	})

	// Check healthz
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("Hello 👋! Healthy 🍏")
	})

	/* ---- Listen for tasks ------- */
	secrets.ListenSecretUpdates()
	runtask.ListenTasks()
	runcodeworker.ListenRunCode()
	runcodeworker.CodeLoadPackagesListen()
	runcodeworker.ListenDisributedStorageDownload()

	/* Every 5 seconds tell mainapp about my status
	Needs to be called after listen for tasks to avoid timing issues when accepting tasks
	*/
	workerhealth.WorkerHealthStart()
	log.Println("🚚 Submitting workers")
	workerhealth.WorkerLoad()

	stop := time.Now()
	// Do something with response
	log.Println("🐆 Worker start time:", fmt.Sprintf("%f", float32(stop.Sub(start))/float32(time.Millisecond))+"ms")

	return app
}

/* Add timer to header */
func Timer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// start timer
		start := time.Now()
		// next routes
		err := c.Next()
		// stop timer
		stop := time.Now()
		ms := float32(stop.Sub(start)) / float32(time.Millisecond)
		c.Append("Server-Timing", fmt.Sprintf("Dataplane;dur=%f", ms))

		return err
	}
}
