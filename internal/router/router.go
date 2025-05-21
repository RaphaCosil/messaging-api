package router

import (
	"github.com/RaphaCosil/messaging-api/internal/handler/http"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *http.UserHandler,
	chatHandler *http.ChatHandler,
) *gin.Engine {
	r := gin.Default()

	userGroup := r.Group("/user")
	{
		userGroup.GET("", userHandler.FindAll)
		userGroup.GET("/:id", userHandler.FindByID)
		userGroup.POST("", userHandler.Create)
		userGroup.PUT("/:id", userHandler.Update)
		userGroup.DELETE("/:id", userHandler.Delete)
		userGroup.GET("/username/:username", userHandler.FindByUsername)
		userGroup.GET("/chat/:chat_id", userHandler.FindByChatID)
	}

	chatGroup := r.Group("/chat")
	{
		chatGroup.GET("", chatHandler.FindAll)
		chatGroup.GET("/:id", chatHandler.FindByID)
		chatGroup.POST("", chatHandler.Create)
		chatGroup.PUT("/:id", chatHandler.Update)
		chatGroup.DELETE("/:id", chatHandler.Delete)
		chatGroup.GET("/user/:user_id", chatHandler.AddUserToChat)
		chatGroup.DELETE("/user/:user_id", chatHandler.RemoveUserFromChat)
		chatGroup.GET("/user/:user_id/chat/:chat_id/access", chatHandler.UserHasAccessToChat)
	}

	return r
}
