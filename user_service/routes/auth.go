package routes

import (
	"notification_system/user_service/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/user/auth")
	auth.Post("/signup", controllers.Signup)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)
}
