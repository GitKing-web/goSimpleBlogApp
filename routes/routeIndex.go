package routes

import (
	"github/GitKing-web/goSimpleBlogApp/controllers"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/hello", controllers.Hello)
	app.Get("/posts", controllers.GetPosts)
	app.Post("/register", controllers.HandleSignup)
	app.Post("/login", controllers.HandleLogin)
	app.Post("/post", controllers.CreatePost)
}
