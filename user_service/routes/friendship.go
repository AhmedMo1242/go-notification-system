package routes

import (
	"notification_system/user_service/controllers"

	"github.com/gofiber/fiber/v2"
)

func FriendshipRoutes(app *fiber.App) {
	friend := app.Group("/user/friend")

	friend.Post("/send", controllers.SendFriendRequest)
	friend.Get("/requests", controllers.ViewFriendRequests)
	friend.Post("/accept", controllers.AcceptFriendRequest)
	friend.Get("/list", controllers.ViewFriends)
	friend.Post("/unfriend", controllers.Unfriend)
	friend.Post("/unfollow", controllers.Unfollow)
	friend.Post("/follow-again", controllers.FollowAgain)
}
