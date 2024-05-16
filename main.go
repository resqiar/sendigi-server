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

	// Prometheus Monitoring
	prometheus := fiberprometheus.New("sendigi-gateway")
	prometheus.RegisterAt(server, "/metrics")
	server.Use(prometheus.Middleware)

	// initiate routes
	routes.InitAPIRoutes(server)

	log.Fatal(server.Listen(":" + os.Getenv("PORT")))
}
