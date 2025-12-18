package main

import (
	"classy-branded-crew/pkg/utils"
	"log"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

func init() {

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Gagal memuat timezone Asia/Jakarta: %v", err)
	}

	time.Local = loc
	log.Printf("Timezone dikonfigurasi ke: %s", time.Now().Location().String())
	log.Printf("Waktu saat ini: %s", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Printf("Warning: .env.%s file not found, using system environment variables", env)
	}
	// db := database.InitDB()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://classybrandedcrew.com, http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        50,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return utils.ErrorResponse(c, 429, "Too many requests", nil)
		},
	}))

	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return utils.SuccessResponse(
			c, 200, "API is Healthy", map[string]string{"server_time": time.Now().Format("2006-01-02 15:04:05"), "timezone": time.Now().Location().String()},
		)
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":3000"
	} else if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("Server starting on port %s in %s mode", port, env)
	log.Fatal(app.Listen(port))

}
