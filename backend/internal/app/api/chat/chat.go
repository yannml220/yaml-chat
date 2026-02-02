package chat

import (
	"time"
	"github.com/yannml220/chat_agent_app/internal/pkg/utils/agent"
)




type Message struct {
	Id string  `json:"id"`
	ConversationId string `json:"conversation_id"`
	Role       string `json:"role"`
	Content    string      `json:"content,omitempty"`
	ToolCalls  []agent.ToolCall  `json:"tool_calls,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}



type Conversation struct {
	Id string  `json:"id"`
	UserId string `json:"user_id"`
	ConceptId string `json:"concept_id"`
	MetaData map[string]interface{} `json:"meta_data"`
	Title *string `json:"title"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}



type ConversationWithMessages struct {
	Id string  `json:"id"`
	UserId string `json:"user_id"`
	ConceptId string `json:"concept_id"`
	MetaData map[string]interface{} `json:"meta_data"`
	Messages []Message `json:"messages"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}






