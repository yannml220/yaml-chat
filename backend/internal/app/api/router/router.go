package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yannml220/chat_agent_app/internal/app/api/auth"
	"github.com/yannml220/chat_agent_app/internal/app/api/chat"
	"github.com/yannml220/chat_agent_app/internal/app/api/graph"
	"github.com/yannml220/chat_agent_app/internal/app/api/middleware"
)

func RegisterAuthRoutes(router fiber.Router, authHandler auth.AuthHandler, am middleware.AuthMiddleware) *fiber.Router {

	authRouter := router.Group("/auth")

	// authRouter.Get("/verify/:token", authHandler.VerifyHandler)
	authRouter.Post("/token/refresh", authHandler.RotateTokenHandler)
	authRouter.Get("/provider/google/signin", authHandler.OauthSigninHandler)
	authRouter.Get("/provider/google/callback", authHandler.Oauth2CallbackHandler)
	authRouter.Get("/me", am.Run(), authHandler.GetUserInfosHandler)
	authRouter.Post("/signout", am.Run(), authHandler.SignoutHandler_)

	return &authRouter

}

func RegisterChatRoutes(router fiber.Router, chatHandler chat.ChatHandler, am middleware.AuthMiddleware) *fiber.Router {

	chatRouter := router.Group("/chat", am.Run())

	chatRouter.Get("/:conversation_id", chatHandler.GetConversationMessages)
	chatRouter.Get("/", chatHandler.GetUserConversations)
	//chatRouter.Post("/",chatHandler.Chat)
	chatRouter.Post("/", chatHandler.InitConversation)
	chatRouter.Post("/streaming", chatHandler.ChatWithStreaming)
	chatRouter.Delete("/:id", chatHandler.DeleteConversationById)

	return &chatRouter
}

func RegisterGraphRoutes(router fiber.Router, graphHandler graph.GraphHandler, am middleware.AuthMiddleware) *fiber.Router {

	graphRouter := router.Group("/graph", am.Run())

	graphRouter.Post("/", graphHandler.CreateGraphHandler)
	graphRouter.Delete("/:graph_id", graphHandler.DeleteGraphHandler)
	graphRouter.Get("/:graph_id", graphHandler.GetGraphByIdHandler)
	graphRouter.Get("/", graphHandler.GetConversationGraphsHandler)

	return &graphRouter
}
