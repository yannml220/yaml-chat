package graph

import (
	"context"

	"github.com/yannml220/chat_agent_app/internal/app/api/chat"
	"github.com/yannml220/chat_agent_app/internal/pkg/types"
)

type GraphService interface {
	CreateGraph(ctx context.Context, conversationId string, graph *Graph) (string, error)
	GetGraphById(ctx context.Context, id string) (*Graph, error)
	UpdateGraph(ctx context.Context, graph *Graph, id string) error
	DeleteGraph(ctx context.Context, id string) error
	GetConversationGraphs(ctx context.Context, conversationId string) ([]Graph, error)
}

type GraphServiceImpl struct {
	types.Service
	repo        GraphRepo
	ChatService chat.ChatService
}

func NewService(repo GraphRepo, chatService chat.ChatService) GraphService {

	return &GraphServiceImpl{
		types.Service{},
		repo,
		chatService,
	}

}

func (gs GraphServiceImpl) CreateGraph(ctx context.Context, conversationId string, graph *Graph) (string, error) {
	return gs.repo.CreateGraph(ctx, conversationId, graph)
}

func (gs GraphServiceImpl) GetGraphById(ctx context.Context, id string) (*Graph, error) {
	return gs.repo.GetGraphById(ctx, id)
}

func (gs GraphServiceImpl) GetConversationGraphs(ctx context.Context, conversationId string) ([]Graph, error) {
	return gs.repo.GetConversationGraphs(ctx, conversationId)
}

func (gs GraphServiceImpl) UpdateGraph(ctx context.Context, graph *Graph, id string) error {
	return gs.repo.UpdateGraph(ctx, graph, id)
}

func (gs GraphServiceImpl) DeleteGraph(ctx context.Context, id string) error {
	return gs.repo.DeleteGraph(ctx, id)
}
