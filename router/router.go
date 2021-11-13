package router

import (
	"mo_fiber_1/handler"
	"mo_fiber_1/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	//Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	//User
	user := api.Group("/user")
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	//Products
	products := api.Group("/products")
	products.Get("", handler.GetAllProducts)
	products.Get("/:id", handler.GetProduct)
	products.Post("/", middleware.Protected(), handler.CreateProduct)
	products.Delete("/:id", handler.DeleteProduct)
}
