package main

import (
	"github.com/RaphaCosil/messaging-api/internal/db"
	"github.com/RaphaCosil/messaging-api/internal/handler/http"
	"github.com/RaphaCosil/messaging-api/internal/repository"
	"github.com/RaphaCosil/messaging-api/internal/service"
	"github.com/RaphaCosil/messaging-api/internal/router"
)

func main() {
	database := db.NewPostgresConnection()

	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	chatRepo := repository.NewChatRepository(database)
	chatService := service.NewChatService(chatRepo)
	userChatRepo := repository.NewUserChatRepository(database)
	userChatService := service.NewUserChatService(userChatRepo)
	chatHandler := http.NewChatHandler(chatService, userChatService)

	r := router.SetupRouter(
		userHandler,
		chatHandler,
	)

	r.Run(":8080")
}
