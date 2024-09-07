package main

import (
	"github/GitKing-web/goSimpleBlogApp/config"
	"github/GitKing-web/goSimpleBlogApp/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := config.ConnectDb(); err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	app := fiber.New()

	routes.Routes(app)

	app.Listen(":3000")
}
