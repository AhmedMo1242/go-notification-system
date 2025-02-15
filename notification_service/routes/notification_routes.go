package routes

import (
	"notification_system/notification_service/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	notifications := app.Group("/notifications")

	notifications.Get("/", controllers.GetNotifications)
	notifications.Put("/:notification_id/read", controllers.MarkAsRead)
}
