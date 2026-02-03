package agent

import (
	"bufio"
	"context"
	"encoding/json"

	//"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

)


const (
	RoleUser      = "user"
	RoleAssistant  ="assistant"
	RoleTool      = "tool"
	RoleSystem    = "system"
)



type ToolCall struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Payload	map[string]interface{} `json:"payload"`
}


type Message struct {
	Role       string `json:"role"`
	Content    string      `json:"content,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}


type ToolMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}


type Tool interface {
	Name() string
	Description() string
	PayloadSchema() map[string]interface{}
	Execute(ctx context.Context, payload map[string]interface{}) (string, error)
}


type ToolRegistry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}





func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}


func (r *ToolRegistry) Register(tool Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name()] = tool
}


func (r *ToolRegistry) GetToolsMetadata() []ToolMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metadata := make([]ToolMetadata, 0, len(r.tools))
	for _, tool := range r.tools {
		metadata = append(metadata, ToolMetadata{
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters :  tool.PayloadSchema(),
		})
	}
	return metadata
}


func (r *ToolRegistry) Execute(ctx context.Context, name string, payload map[string]interface{}) (string, error) {
	r.mu.RLock()
	tool, exists := r.tools[name]
	r.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("tool '%s' not found", name)
	}

	execCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result, err := tool.Execute(execCtx, payload)
	if err != nil {
		return "", fmt.Errorf("tool '%s' execution failed: %w", name, err)
	}

	return result, nil
}



type LLMResponse struct {
	Content   string     
	ToolCalls []ToolCall 
	IsToolUse bool      
}


type LLMClient interface {
	Chat(ctx context.Context, messages []Message, tools []ToolMetadata) (*LLMResponse, error)
	StreamFinalResponseWithStreaming(ctx context.Context, responseId string ,w *bufio.Writer ,messages []Message) (string, error) 
}


type AgentConfig struct {
	MaxIterations    int           
	GlobalTimeout    time.Duration
	EnableDebugLogs  bool        
}



func DefaultAgentConfig() AgentConfig {
	return AgentConfig{
		MaxIterations:   10,
		GlobalTimeout:   4 * time.Minute,
		EnableDebugLogs: true,
	}
}




type Agent struct {
	llmClient    LLMClient
	toolRegistry *ToolRegistry
	config       AgentConfig
}




func NewAgent(llmClient LLMClient, toolRegistry *ToolRegistry, config AgentConfig) *Agent {
	return &Agent{
		llmClient:    llmClient,
		toolRegistry: toolRegistry,
		config:       config,
	}
}



type FinalResult struct {
	Response  string   
	History []Message
	SessionHistory []Message
	ToolsUsed []string
	Error     error  
}




func (a *Agent) ProcessRequest(ctx context.Context, userMessage string, history []Message) (*FinalResult, error) {

	execCtx, cancel := context.WithTimeout(ctx, a.config.GlobalTimeout)
	defer cancel()

	sessionHistory := []Message{}

	sessionHistory = append( sessionHistory , Message{
		Role:      RoleUser,
		Content:   userMessage,
		CreatedAt: time.Now(),
	})

	messages := append(history, Message{
		Role:      RoleUser,
		Content:   userMessage,
		CreatedAt: time.Now(),
	})

	toolsUsed := make(map[string]bool)

	for i := 0; i < a.config.MaxIterations; i++ {
		if a.config.EnableDebugLogs {
			fmt.Printf("[Agent] Iteration %d/%d\n", i+1, a.config.MaxIterations)
		}

		select {
		case <-execCtx.Done():
			return nil, errors.New("agent timeout: request took too long")
		default:
		}

		toolsMetadata := a.toolRegistry.GetToolsMetadata()

		response, err := a.llmClient.Chat(execCtx, messages, toolsMetadata)
		if err != nil {
			return nil, fmt.Errorf("LLM request failed: %w", err)
		}

		if !response.IsToolUse {
			if a.config.EnableDebugLogs {
				fmt.Printf("[Agent] Final response received\n")
			}

			toolsList := make([]string, 0, len(toolsUsed))
			for tool := range toolsUsed {
				toolsList = append(toolsList, tool)
			}


			fmt.Printf("[Agent] Result  %d tool(s)\n", len(response.Content))

			sessionHistory = append(sessionHistory, Message{
				Role:      RoleAssistant,
				Content:   response.Content,
				CreatedAt: time.Now(),
			})


			return &FinalResult{
				Response:  response.Content,
				History : messages ,
				SessionHistory : sessionHistory ,
				ToolsUsed: toolsList,
			}, nil
		}

		if a.config.EnableDebugLogs {
			fmt.Printf("[Agent] LLM requested %d tool(s)\n", len(response.ToolCalls))
		}

		messages = append(messages, Message{
			Role:      RoleAssistant,
			ToolCalls: response.ToolCalls,
			CreatedAt: time.Now(),
		})

		sessionHistory = append(sessionHistory, Message{
			Role:      RoleAssistant,
			ToolCalls: response.ToolCalls,
			CreatedAt: time.Now(),
		})


		for _, toolCall := range response.ToolCalls {
			if a.config.EnableDebugLogs {
				fmt.Printf("[Agent] Executing tool: %s\n", toolCall.Name)
			}

			toolsUsed[toolCall.Name] = true

			result, err := a.toolRegistry.Execute(execCtx, toolCall.Name, toolCall.Payload)

			if err != nil {
				result = fmt.Sprintf("Erreur lors de l'exécution du tool '%s': %v", toolCall.Name, err)
				if a.config.EnableDebugLogs {
					fmt.Printf("[Agent] Tool execution error: %v\n", err)
				}
			}

			fmt.Printf("[Agent result] : %v\n", result)

			messages = append(messages, Message{
				Role:       RoleTool,
				ToolCallID: toolCall.ID,
				Content:    result,
				CreatedAt:  time.Now(),
			})
		}
	}

	return nil, fmt.Errorf("max iterations (%d) reached without final response", a.config.MaxIterations)
}



type StreamEventType string

const (
	EventTypeText       StreamEventType = "text"       
	EventTypeTool    StreamEventType = "tool"     
	EventTypeError    StreamEventType = "error"    
	EventTypeDone    StreamEventType = "done"    
)


type StreamEvent struct {
    Type    StreamEventType      `json:"type"`    
    Content interface{} `json:"content,omitempty"`
	Error *StreamError `json:"error,omitempty"`
	Timestamp string          `json:"timestamp"`
}



func SendSSE(w *bufio.Writer, event StreamEvent) error {
    data, _ := json.Marshal(event)
    _, err := fmt.Fprintf(w, "data: %s\n\n", data)
	if err != nil {
		return err
	}
    return w.Flush()
}


func SendDone(w *bufio.Writer) error {
	_, err := fmt.Fprint(w, "data: [DONE]\n\n")
	if err != nil {
		return err
	}
	return w.Flush()
}


type BaseStreamChunk struct {
	ID  string `json:"id"`
	Type  string `json:"type"`
	Model  string `json:"model"`
	Timestamp  int64 `json:"timestamp"`


}


type StreamChunk struct {
	BaseStreamChunk
	Delta string `json:"delta,omitempty"`
	Content string `json:"content,omitempty"`
    Role    string `json:"role,omitempty"`
}


type StreamError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}


type ErrorStreamChunk struct {
	BaseStreamChunk
	Error *StreamError `json:"error,omitempty"`

}


type DoneStreamChunk struct {
	BaseStreamChunk
	FinishReadon string `json:"finishReason"`

}



func TanstackSendSSE(w *bufio.Writer, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "data: %s\n\n", data)
	return w.Flush()
}




func (a *Agent) ProcessRequestWithStreaming(w *bufio.Writer,ctx context.Context, userMessage string, history []Message) (*FinalResult, error) {


	execCtx, cancel := context.WithTimeout(ctx, a.config.GlobalTimeout)
	defer cancel()

	responseId := fmt.Sprintf("msg_%d", time.Now().UnixNano())

	sessionHistory := []Message{}

	sessionHistory = append( sessionHistory , Message{
		Role:      RoleUser,
		Content:   userMessage,
		CreatedAt: time.Now(),
	})

	messages := append(history, Message{
		Role:      RoleUser,
		Content:   userMessage,
		CreatedAt: time.Now(),
	})

	toolsUsed := make(map[string]bool)

	var finalResponse  *FinalResult

	for i := 0; i < a.config.MaxIterations; i++ {

		if a.config.EnableDebugLogs {
			fmt.Printf("[Agent] Iteration %d/%d\n", i+1, a.config.MaxIterations)
		}

		select {
		case <-execCtx.Done():
			return nil ,errors.New("timeout")
		default:
		}

		toolsMetadata := a.toolRegistry.GetToolsMetadata()

		response, err := a.llmClient.Chat(execCtx, messages, toolsMetadata)
		if err != nil {
			return  nil , err
		}

		if !response.IsToolUse {
			if a.config.EnableDebugLogs {
				fmt.Printf("[Agent] Final response received\n")
			}

			toolsList := make([]string, 0, len(toolsUsed))

			for tool := range toolsUsed {
				toolsList = append(toolsList, tool)
			}

			fmt.Printf("[Agent] Result  %d tool(s)\n", len(response.Content))

			
			content , err := a.llmClient.StreamFinalResponseWithStreaming(ctx ,responseId ,w,messages )

			sessionHistory = append(sessionHistory, Message{
				Role:      RoleAssistant,
				Content:   response.Content,
				CreatedAt: time.Now(),
			})

			messages = append(messages, Message{
				Role:      RoleAssistant,
				Content:   response.Content,
				CreatedAt: time.Now(),
			})


			//fmt.Print("here is the streaming content :", content)

			if err != nil {
				return  nil , err
			}

			finalResponse =  &FinalResult{
				Response:  content,
				History : messages ,
				SessionHistory : sessionHistory ,
				ToolsUsed: toolsList,
			}
			return  finalResponse , nil
		}

		if a.config.EnableDebugLogs {
			fmt.Printf("[Agent] LLM requested %d tool(s)\n", len(response.ToolCalls))
		}

		messages = append(messages, Message{
			Role:      RoleAssistant,
			ToolCalls: response.ToolCalls,
			CreatedAt: time.Now(),
		})

		sessionHistory = append(sessionHistory, Message{
			Role:      RoleAssistant,
			ToolCalls: response.ToolCalls,
			CreatedAt: time.Now(),
		})


		for _, toolCall := range response.ToolCalls {
			if a.config.EnableDebugLogs {
				fmt.Printf("[Agent] Executing tool: %s\n", toolCall.Name)
			}

			toolsUsed[toolCall.Name] = true

			result, err := a.toolRegistry.Execute(execCtx, toolCall.Name, toolCall.Payload)

			if err != nil {
				result = fmt.Sprintf("Erreur lors de l'exécution du tool '%s': %v", toolCall.Name, err)
				if a.config.EnableDebugLogs {
					fmt.Printf("[Agent] Tool execution error: %v\n", err)
				}
			}

			fmt.Printf("[Agent result] : %v\n", result)

			messages = append(messages, Message{
				Role:       RoleTool,
				ToolCallID: toolCall.ID,
				Content:    result,
				CreatedAt:  time.Now(),
			})
		}
	}

	return nil , nil

}





