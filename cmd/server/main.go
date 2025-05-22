package main

import (
	"github.com/RaphaCosil/messaging-api/internal/db"
	"github.com/RaphaCosil/messaging-api/internal/repository"
	"github.com/RaphaCosil/messaging-api/internal/service"
	"github.com/RaphaCosil/messaging-api/internal/router"
	"github.com/RaphaCosil/messaging-api/internal/handler"
)

func main() {
	database := db.NewPostgresConnection()

	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	chatRepo := repository.NewChatRepository(database)
	chatService := service.NewChatService(chatRepo)
	userChatRepo := repository.NewUserChatRepository(database)
	userChatService := service.NewUserChatService(userChatRepo)
	chatHandler := handler.NewChatHandler(chatService, userChatService)


	messageRepo := repository.NewMessageRepository(database)
	messageService := service.NewMessageService(messageRepo)
	
	wsHub := handler.NewHub()
	wsHandler := handler.NewWebSocketHandler(wsHub, messageService)
	go wsHandler.Hub.Run()

	
	r := router.SetupRouter(
		userHandler,
		chatHandler,
		wsHandler,
	)

	r.Run(":8080")
}
