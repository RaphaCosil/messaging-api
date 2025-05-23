package router

import (
	"github.com/RaphaCosil/messaging-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	chatHandler *handler.ChatHandler,
	wsHandler *handler.WebSocketHandler,
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

	r.GET("/ws", func(c *gin.Context) {
        wsHandler.HandleConnection(
			c.Writer,
			c.Request,
		)
    })
	return r
}
