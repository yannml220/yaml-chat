package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yannml220/chat_agent_app/internal/app/api/auth"
	"github.com/yannml220/chat_agent_app/internal/app/api/chat"
	"github.com/yannml220/chat_agent_app/internal/app/api/graph"
	"github.com/yannml220/chat_agent_app/internal/app/api/middleware"
	"github.com/yannml220/chat_agent_app/internal/app/api/router"
	"github.com/yannml220/chat_agent_app/internal/app/api/user"
	"github.com/yannml220/chat_agent_app/internal/app/db"
	"github.com/yannml220/chat_agent_app/internal/pkg/initializers"
	"github.com/yannml220/chat_agent_app/internal/pkg/utils/agent"
)

type ApiServer struct {
	app  *fiber.App
	addr string
}

func newApiServer(addr string) *ApiServer {

	return &ApiServer{
		app:  fiber.New(),
		addr: addr,
	}

}

func (as *ApiServer) setup() {
	as.app.Use(cors.New(
		cors.Config{
			AllowOrigins:     "http://localhost:5173",
			AllowCredentials: true,
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Device-ID",
			ExposeHeaders:    "X-Conversation-Id",
		},
	))
	as.app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	as.app.Use(logger.New())

}

func (as *ApiServer) run() error {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	if port == "" {
		port = "5000"
	}
	if host == "" {
		host = "localhost"
	}

	if err := as.app.Listen(host + ":" + port); err != nil {
		return err

	}
	return nil
}

func main() {

	apiServer := newApiServer("")

	apiServer.setup()

	var err error

	if err = initializers.LoadEnvVariables(); err != nil {
		log.Fatal("couldn't load the env variables ", err)

	}

	var conn *db.Db = nil

	if conn, err = db.NewDb(); err != nil {
		log.Fatal("couldn't connected to the db")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	PingErr := conn.Ping(ctx)

	if PingErr != nil {
		log.Fatal("couldn't connected to the db")
	}

	api := apiServer.app.Group("/api")

	v1 := api.Group("/v1")

	userRepo := user.NewRepo(conn)
	userService := user.NewService(userRepo)

	authRepo := auth.NewRepo(conn)
	authService := auth.NewService(authRepo, userService)

	authMiddlware := middleware.NewAuthMiddleware(authService)

	// userHandler := user.NewHandler(userService)
	// router.RegisterUserRoutes(v1, userHandler, *authMiddlware)

	authHandler := auth.NewHandler(authService)
	router.RegisterAuthRoutes(v1, authHandler, *authMiddlware)

	ollamaClient, err := agent.NewOllamaClientInstance(agent.DefaultOllamaConfig())

	if err != nil {
		log.Fatal(err)
	}
	//kimi-k2.5:cloud
	//gpt-oss:20b-cloud

	ollamaClientStruct := agent.NewOllamaClientStruct(ollamaClient, "kimi-k2.5:cloud", 0.7, "your are a relevant assistant")

	toolRegistry := agent.NewToolRegistry()

	agent := agent.NewAgent(ollamaClientStruct, toolRegistry, agent.DefaultAgentConfig())

	chatRepo := chat.NewRepo(conn)
	chatService := chat.NewService(chatRepo)
	chatHandler := chat.NewHandler(chatService, agent)
	router.RegisterChatRoutes(v1, chatHandler, *authMiddlware)

	graphRepo := graph.NewRepo(conn)
	graphService := graph.NewService(graphRepo, chatService)
	graphHandler := graph.NewHandler(graphService)
	router.RegisterGraphRoutes(v1, graphHandler, *authMiddlware)

	if err := apiServer.run(); err != nil {
		log.Fatal(err)

	}

}
