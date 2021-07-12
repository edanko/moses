package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/edanko/moses/cmd/server/routes"
	"github.com/edanko/moses/service/nest"
	"github.com/edanko/moses/service/profile"
	"github.com/edanko/moses/service/remnant"
	"github.com/edanko/moses/service/spacing"

	repo "github.com/edanko/moses/repository/mongo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	db, err := repo.NewMongoDB(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DATABASE"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Database connected!")

	remnantCollection := db.Collection("remnants")
	remnantRepo := repo.NewRemnantRepo(remnantCollection)
	remnantService := remnant.NewService(remnantRepo)

	profileCollection := db.Collection("profiles")
	profileRepo := repo.NewProfileRepo(profileCollection)
	profileService := profile.NewService(profileRepo)

	spacingCollection := db.Collection("spacing")
	spacingRepo := repo.NewSpacingRepo(spacingCollection)
	spacingService := spacing.NewService(spacingRepo)

	nestCollection := db.Collection("nests")
	nestRepo := repo.NewNestRepo(nestCollection)
	nestService := nest.NewService(nestRepo, remnantService, profileService, spacingService)

	prod, _ := strconv.ParseBool(os.Getenv("PROD"))

	app := fiber.New(fiber.Config{
		Prefork: prod,
	})

	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: 1,
	}))
	app.Use(cors.New())

	api := app.Group("/api/v1")
	routes.RemnantRouter(api, remnantService)
	routes.ProfileRouter(api, profileService)
	routes.NestRouter(api, nestService)

	if err := app.Listen(os.Getenv("API_PORT")); err != nil {
		log.Panic(err)
	}
}
