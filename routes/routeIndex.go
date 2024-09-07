package routes

import (
	"github/GitKing-web/goSimpleBlogApp/controllers"
	// "github/GitKing-web/goSimpleBlogApp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/hello", controllers.Hello)
	app.Get("/posts", controllers.GetPosts)
	app.Get("/posts/:id", controllers.GetPost)
	app.Post("/register", controllers.HandleSignup)
	app.Post("/login", controllers.HandleLogin)
	app.Post("/post", controllers.CreatePost)
	app.Put("/posts/:id", controllers.UpdatePost)
	app.Delete("/posts/:id", controllers.DeletePost)
}
