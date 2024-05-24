package main

import (
	"log"
	"os"
	"sendigi-server/configs"
	"sendigi-server/routes"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// database connection
	if err := configs.InitDBCon(); err != nil {
		log.Fatal(err)
	}

	// redis connection
	configs.InitRedis()
	amqpCon := configs.InitRabbitMQ()
	defer amqpCon.Close()

	server := fiber.New()

	// Setup CORS
	server.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_CLIENT_URL"),
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	// Initialize sessions
	configs.InitGoogleConfig()
	configs.InitSession()
	configs.InitStateSession()

	// Compression
	server.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	// Init logger
	logFile, err := os.OpenFile("./logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening log file", err)
	}
	defer logFile.Close()
	server.Use(logger.New(logger.Config{
		Format:        "[${time}] ${ip}:${port} ${status} - ${method} ${path} | ${latency} \n",
		DisableColors: true,
		Output:        logFile,
	}))

	// Prometheus Monitoring
	prometheus := fiberprometheus.New("sendigi-gateway")
	prometheus.RegisterAt(server, "/metrics")
	server.Use(prometheus.Middleware)

	// initiate routes
	routes.InitAPIRoutes(server)

	log.Fatal(server.Listen(":" + os.Getenv("PORT")))
}
