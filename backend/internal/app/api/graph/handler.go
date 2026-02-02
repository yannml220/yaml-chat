package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/yannml220/chat_agent_app/internal/app/api/auth"
	"github.com/yannml220/chat_agent_app/internal/pkg/types"
)

type CreateGraphRequest struct {
	ConversationId  string `json:"conversation_id" validate:"required"`
	WindowSize  int `json:"window_size" validate:"required"`
}

type UpdateGraphRequest struct {
	GraphId  string `json:"graph_id" validate:"required"`
	Graph  Graph `json:"graph" validate:"required"`
}


type GraphHandler interface {
	CreateGraphHandler(c *fiber.Ctx) error
	DeleteGraphHandler(c *fiber.Ctx) error
	GetGraphByIdHandler(c *fiber.Ctx) error
	GetConversationGraphsHandler(c *fiber.Ctx) error
}

type graphHandler struct {
	types.Handler
	Service GraphService
}

func NewHandler(gs GraphService) GraphHandler {
	return &graphHandler{
		Service: gs,
	}

}

func (gh *graphHandler) CreateGraphHandler(c *fiber.Ctx) error {

	ctx , cancel := context.WithTimeout(c.Context(), 5 * time.Second) 

	defer cancel()

	body := new(CreateGraphRequest)

	if err := gh.ValidateJsonBody(c,body) ; err != nil {

		log.Print("Error while validating the request body :",err.Error())
		return gh.JSONResponse(c,fiber.ErrBadRequest.Code , fiber.Map{

		})
	}

	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		return gh.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return gh.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}


	conversationId := body.ConversationId

	messages , err := gh.Service.(*GraphServiceImpl).ChatService.GetConversationMessagesByRole(ctx,conversationId,"assistant")

	if err != nil {
		fmt.Printf("error fetching the conversation messages by role : %v\n", err)
		return gh.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
			"error" : err,
		})
	}

	var  messageContentList []string

	for _, message := range messages {
		messageContentList = append(messageContentList , message.Content )
	}

	text := strings.Join(messageContentList , "\n" )

	payload := map[string]interface{}{
				"text":        text,
				"window_size":		body.WindowSize,
	}


	data, _ := json.Marshal(payload)
	response, err := http.Post("http://localhost:8000/nlp/compute-graph", "application/json", bytes.NewBuffer(data))
	if err != nil {

		return gh.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
			"error" : err,

		})

	}

	defer response.Body.Close()

	respBody, _ := io.ReadAll(response.Body)
	

	type ResponseBody struct {
		GraphData GraphData `json:"graph_data"`
	}


	var resp ResponseBody

	if err := json.Unmarshal([]byte(respBody) , &resp) ; err != nil {
		return gh.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
			"error" : err,
		})
	}

	newGraph := Graph{
		ConversationId: conversationId,
		Data: resp.GraphData,
	}

	 id, err := gh.Service.CreateGraph(ctx,conversationId,&newGraph)
	
	 if  err != nil {
		fmt.Printf("error creating the graph : %v\n", err)
		return gh.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
			"error" : err,
		})
	}

	return gh.JSONResponse(c,fiber.StatusOK , fiber.Map{
					"id" :  id, 
	})


}





func (gh *graphHandler) DeleteGraphHandler(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	graphId := c.Params("graph_id")
	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		return gh.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return gh.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	err := gh.Service.DeleteGraph(ctx, graphId)

	if  err != nil {
		fmt.Printf("error deleting the graph : %v\n", err)
		return gh.JSONResponse(c,fiber.ErrInternalServerError.Code , fiber.Map{
			"error" : err,
		})
	}

	return gh.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"id": graphId,
	})

}


func (gh *graphHandler) GetConversationGraphsHandler(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	conversationId := c.Query("conversation_id")
	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		return gh.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return gh.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	if conversationId == "" {
		return gh.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "conversation id is required",
		})
	}


	graphs, err := gh.Service.GetConversationGraphs(ctx, conversationId)

	if err != nil {

		log.Print("Error fetching the conversation graphs", err.Error())
		return gh.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error": err,
		})
	}

	if graphs == nil {

		log.Print("No graph exist for that conversation")
		return gh.JSONResponse(c, fiber.StatusOK, fiber.Map{})
	}

	return gh.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"graphs": graphs,
	})

}



func (gh *graphHandler) GetGraphByIdHandler(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	graphId := c.Params("graph_id")
	userId := c.Query("user_id")

	authCtx, ok := c.Locals("auth_context").(auth.AuthContext)
	if !ok {
		log.Print("Missing or invalid auth context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if userId == "" {
		return gh.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "user id is required",
		})
	}

	if userId != authCtx.UserId {
		log.Printf("User ID mismatch: request=%s, auth=%s", userId, authCtx.UserId)
		return gh.JSONResponse(c, fiber.StatusForbidden, fiber.Map{
			"error":   "forbidden",
			"message": "User ID mismatch",
		})
	}

	if graphId == "" {
		return gh.JSONResponse(c, fiber.ErrBadRequest.Code, fiber.Map{
			"error":   "invalid_request",
			"message": "graph id is required",
		})
	}


	graph, err := gh.Service.GetGraphById(ctx, graphId)

	if err != nil {

		log.Print("Error fetching the conversation graphs", err.Error())
		return gh.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error": err,
		})
	}

	if graph == nil {

		log.Print("No graph exist for with this id")
		return gh.JSONResponse(c, fiber.StatusOK, fiber.Map{})
	}

	return gh.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"graph": graph,
	})

}




