package chat

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yannml220/chat_agent_app/internal/app/api/auth"
	"github.com/yannml220/chat_agent_app/internal/pkg/types"
	"github.com/yannml220/chat_agent_app/internal/pkg/utils/agent"
)

type ChatRequest struct {
	Messages       []Message              `json:"messages"`
	ConversationId string                 `json:"conversation_id,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
}

type InitConversationRequest struct {
	Query          string `json:"query"`
}


type UpdateChunkRequest struct {
	RawText string `json:"raw_text"`
}

type EmptyBody struct {
}

type ChatHandler interface {
	//Chat(c *fiber.Ctx) error
	ChatWithStreaming(c *fiber.Ctx) error
	InitConversation(c *fiber.Ctx) error
	GetUserConversations(c *fiber.Ctx) error
	GetConversationMessages(c *fiber.Ctx) error
	DeleteConversationById(c *fiber.Ctx) error
	DeleteUserConversations(c *fiber.Ctx) error
}

type ChatHandlerImpl struct {
	types.Handler
	Service ChatService
	agent   *agent.Agent
}

func NewHandler(cs ChatService, agent *agent.Agent) ChatHandler {
	return &ChatHandlerImpl{
		Service: cs,
		agent:   agent,
	}

}

/*
func (ch *ChatHandlerImpl) Chat(c *fiber.Ctx) error {

	ctx , cancel := context.WithTimeout(c.Context(), 4 * time.Minute)

	defer cancel()

	body := new(ChatRequest)

	if err := ch.ValidateJsonBody(c,body) ; err != nil {

		log.Print("Error while validating the request body :",err.Error())
		return ch.JSONResponse(c,fiber.ErrBadRequest.Code , fiber.Map{

		})
	}


	if body.Message == "" {
		return ch.JSONResponse(c,fiber.ErrBadRequest.Code , fiber.Map{
			"error":   "invalid_request",
			"message": "Message is required",
		})
	}

	agentHistory := []agent.Message{}

	history := []Message{}

	conversationId :=  body.ConversationId

	userId :=  body.UserId


	if conversationId == "" {

		newConv := &Conversation {}

		newConversationId , err := ch.Service.CreateUserConversation(ctx ,userId,newConv)
		if err != nil {

			fmt.Printf("error creating the user conversation processing the user request: %v\n", err)
			return ch.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
				"error" :err,
			})
		}
		conversationId = newConversationId


	}else {

		conv , err := ch.Service.GetConversationById(ctx ,conversationId)
		if err != nil {
			fmt.Printf("error getting the conversation by id: %v\n", err)
			return ch.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
				"error" :err,
			})
		}
		conversationId = conv.Id

		history , err = ch.Service.GetConversationMessages(ctx ,conversationId)
		if err != nil {
			fmt.Printf("error getting the conversation messages: %v\n", err)
			return ch.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
				"error" :err,
			})
		}


		for _ , message := range history {
			agentHistory = append(agentHistory , agent.Message{
				Role : message.Role ,
				Content : message.Content ,
				ToolCalls : message.ToolCalls ,
				CreatedAt : *message.CreatedAt ,
			})
		}
	}

	result , err :=  ch.agent.ProcessRequest(ctx, body.Message,agentHistory )


	//fmt.Printf("history : %v\n",  result.History)
	//fmt.Printf("session history : %v\n",  result.SessionHistory)


	if err != nil {

		fmt.Printf("error processing the user request: %v\n", err)
		return ch.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
			"error" :err,
		})

	}

	newHistory :=  []Message{}

	for _ , message := range result.SessionHistory {
		newHistory = append(newHistory , Message{
			ConversationId : conversationId ,
			Role : message.Role ,
			Content : message.Content ,
			ToolCalls : message.ToolCalls ,
			CreatedAt : &message.CreatedAt ,
		})
	}

	_ , err =  ch.Service.CreateMessages(ctx , conversationId , newHistory)

	if err != nil {

		fmt.Printf("error creating adding the messages: %v\n", err)
		return ch.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
			"message": "error processing the user request",
			"error" :err,
		})
	}

	return ch.JSONResponse(c,fiber.StatusOK , fiber.Map{
		"content": result.Response,

	})

}

*/


func (ch *ChatHandlerImpl) InitConversation(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		userId = authCtx.UserId
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return ch.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	body := new(InitConversationRequest)

	if err := ch.ValidateJsonBody(c, body); err != nil {
		log.Print("Error while validating the request body :", err.Error())
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{})
	}

	query := body.Query

	if query == "" {
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "query is required",
		})
	}

	firstMessageContent := ""
	if len(query) > 200 {
		firstMessageContent = strings.TrimSpace(query[:200])
	} else {
		firstMessageContent = strings.TrimSpace(query)
	}

	newConv := &Conversation{
		Title: &firstMessageContent,
	}

	newConversationId, err := ch.Service.CreateUserConversation(ctx, userId, newConv)
	if err != nil {
		fmt.Printf("error creating the user conversation for init : %v\n", err)
		return ch.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error": err,
		})
	}

	return ch.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"id": newConversationId,
	})

}



func (ch *ChatHandlerImpl) ChatWithStreaming(c *fiber.Ctx) error {
	
	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		userId = authCtx.UserId
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return ch.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	body := new(ChatRequest)

	if err := ch.ValidateJsonBody(c, body); err != nil {

		log.Print("Error while validating the request body :", err.Error())
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{})
	}

	agentHistory := []agent.Message{}

	conversationId := body.ConversationId

	messages := body.Messages


	log.Println("conversation id :", conversationId)
	log.Println("user id :", userId)
	log.Print("messages  :", messages)

	if len(messages) == 0 || messages == nil {
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "Messages are required",
		})
	}

	if userId == "" {
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}


	for _, message := range messages[:len(messages)-1] {
		agentHistory = append(agentHistory, agent.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	c.Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")
	c.Set("Transfer-Encoding", "chunked")

	streamCtx, streamCancel := context.WithTimeout(context.Background(), 4*time.Minute)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {

		defer streamCancel()
		
		userQuery := messages[len(messages)-1].Content

		result, err := ch.agent.ProcessRequestWithStreaming(w, streamCtx, userQuery, agentHistory)

		if err != nil {

			fmt.Printf("error processing the user request: %v\n", err)

			agent.SendSSE(w, agent.StreamEvent{Type: "error", Timestamp: time.Now().String(), Error: &agent.StreamError{Message: err.Error()}})

			return

		}

		newHistory := []Message{}


		for _, message := range result.SessionHistory {
			
			newHistory = append(newHistory, Message{
				ConversationId: conversationId,
				Role:           message.Role,
				Content:        message.Content,
				ToolCalls:      message.ToolCalls,
				CreatedAt:      &message.CreatedAt,
			})
		}

		_, err = ch.Service.CreateMessages(streamCtx, conversationId, newHistory)

		if err != nil {
			fmt.Printf("error adding the messages: %v\n", err)
			agent.SendSSE(w, agent.StreamEvent{Type: "error", Timestamp: time.Now().String(), Error: &agent.StreamError{Message: err.Error()}})
			return
		}

		
	})

	return nil

}

func (ch *ChatHandlerImpl) GetUserConversations(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return ch.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	conversations, err := ch.Service.GetUserConversations(ctx, userId)

	if err != nil {

		log.Print("Error fetching the user conversations :", err.Error())

		return ch.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error": err,
		})
	}

	if conversations == nil {

		log.Print("No conversaton exist for that user")
		return ch.JSONResponse(c, fiber.StatusOK, fiber.Map{})
	}

	return ch.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"conversations": conversations,
	})

}

func (ch *ChatHandlerImpl) GetConversationMessages(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	conversationId := c.Params("conversation_id")
	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return ch.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	if conversationId == "" {
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "conversation id is required",
		})
	}

	// Verify conversation ownership
	conv, err := ch.Service.GetConversationById(ctx, conversationId)
	if err != nil {
		log.Printf("Error fetching conversation: %v", err)
		return ch.JSONResponse(c, fiber.StatusNotFound, fiber.Map{"error": "Conversation not found"})
	}

	if conv.UserId != authCtx.UserId {
		log.Printf("Ownership mismatch: conv.UserId=%s, auth.UserId=%s", conv.UserId, authCtx.UserId)
		return ch.JSONResponse(c, fiber.StatusForbidden, fiber.Map{"error": "Forbidden: You don't own this conversation"})
	}

	messages, err := ch.Service.GetConversationMessages(ctx, conversationId)

	if err != nil {

		log.Print("Error fetching the conversation messages", err.Error())
		return ch.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error": err,
		})
	}

	if messages == nil {

		log.Print("No message exist for that conversation")
		return ch.JSONResponse(c, fiber.StatusOK, fiber.Map{})
	}

	return ch.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"messages": messages,
	})

}

func (ch *ChatHandlerImpl) DeleteConversationById(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	conversationId := c.Params("conversation_id")
	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return ch.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	// Verify conversation ownership before deletion
	conv, err := ch.Service.GetConversationById(ctx, conversationId)
	if err != nil {
		log.Printf("Error fetching conversation: %v", err)
		return ch.JSONResponse(c, fiber.StatusNotFound, fiber.Map{"error": "Conversation not found"})
	}

	if conv.UserId != authCtx.UserId {
		log.Printf("Ownership mismatch: conv.UserId=%s, auth.UserId=%s", conv.UserId, authCtx.UserId)
		return ch.JSONResponse(c, fiber.StatusForbidden, fiber.Map{"error": "Forbidden: You don't own this conversation"})
	}

	err = ch.Service.DeleteConversationById(ctx, conversationId)

	return ch.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"id": conversationId,
	})

}

func (ch *ChatHandlerImpl) DeleteUserConversations(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	userId := c.Params("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		return ch.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return ch.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	_ = ch.Service.DeleteUserConversations(ctx, userId)

	return ch.JSONResponse(c, fiber.StatusOK, fiber.Map{})

}
