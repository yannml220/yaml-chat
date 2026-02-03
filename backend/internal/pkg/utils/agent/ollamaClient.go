package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/ollama/ollama/api"
)



type OllamaClientStruct struct {
	Client *api.Client
	ModelName string     
	Temperature float32 
	SystemPrompt string

}


func NewOllamaClientStruct( client *api.Client , modelName string , temperature float32 , systemPrompt string ) *OllamaClientStruct {
	return &OllamaClientStruct {
		Client : client ,
		ModelName : modelName ,
		Temperature : temperature ,
		SystemPrompt : systemPrompt ,
	}
}




type OllamaConfig struct {
    BaseURL *string        `json:"base_url"`
    Timeout time.Duration `json:"timeout"`

    ModelName   string  `json:"model_name"` 
    Temperature float32 `json:"temperature"`
    SystemPrompt string `json:"system_prompt"`
    
    MaxRetries int `json:"max_retries"`
}



func DefaultOllamaConfig() OllamaConfig {
	baseURL :=  "http://127.0.0.1:11434"
	
	return OllamaConfig{
		BaseURL:      &baseURL,
		//Timeout:      90 * time.Second, 
	}
}



func NewOllamaClientInstance(cfg OllamaConfig) (*api.Client, error) {
    targetURL :=  "http://127.0.0.1:11434"
    if cfg.BaseURL != nil && *cfg.BaseURL != "" {
        targetURL = *cfg.BaseURL
    }

    u, err := url.Parse(targetURL)
    if err != nil {

		fmt.Printf("error creating the ollama instance  : %v\n", err)
        return nil, err
    }

    return api.NewClient(u, &http.Client{Timeout: cfg.Timeout}), nil
}



/*func toOllamaMessages(msgs []Message) []api.Message {
	out := make([]api.Message, len(msgs))
	for i, m := range msgs {
		out[i] = api.Message{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return out
} */


func toOllamaMessages(msgs []Message) []api.Message {
    out := make([]api.Message, 0, len(msgs))
    
    for _, m := range msgs {
        ollamaMsg := api.Message{
            Role:    m.Role,
            Content: m.Content,
        }
        
        if len(m.ToolCalls) > 0 {
            ollamaMsg.ToolCalls = make([]api.ToolCall, len(m.ToolCalls))
            for i, tc := range m.ToolCalls {
                ollamaMsg.ToolCalls[i] = api.ToolCall{
                    Function: api.ToolCallFunction{
                        Name:      tc.Name,
                        Arguments: tc.Payload, 
                    },
                }
            }
        }
        
        out = append(out, ollamaMsg)
    }
    
    return out
}





func toOllamaTools(tools []ToolMetadata) []api.Tool {
	out := make([]api.Tool, len(tools))
	for i, t := range tools {
		bytes, _ := json.Marshal(t.Parameters)

		tool := api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        t.Name,
				Description: t.Description,
			},
		}

		json.Unmarshal(bytes, &tool.Function.Parameters)

		out[i] = tool
	}
	return out
}




func toAgentToolCalls(ollamaToolCalls []api.ToolCall) []ToolCall {
	if ollamaToolCalls == nil {
		return nil
	}

	out := make([]ToolCall, len(ollamaToolCalls))
	for i, tc := range ollamaToolCalls {
		out[i] = ToolCall{
			ID:      "", 
			Name:    tc.Function.Name,
			Payload: tc.Function.Arguments,
		}
	}
	return out
}




func (oc *OllamaClientStruct) Chat(ctx context.Context,messages []Message,tools []ToolMetadata) (*LLMResponse, error){
	
	req := &api.ChatRequest{
		Model:    oc.ModelName,
		Messages: toOllamaMessages(messages), 
		Tools:    toOllamaTools(tools),      
		Stream:   new(bool),
	}

	var finalResponse *api.ChatResponse

	err := oc.Client.Chat(ctx, req, func(resp api.ChatResponse) error {
		finalResponse = &resp
        return nil
	})

	if err != nil {
		
		fmt.Printf("error calling Chat function : %v\n", err)
		return  nil , err

	}


	fmt.Printf("the result of the ollama chat method :\n %v", finalResponse.Message.Content)

	
	llmResponse := LLMResponse{
		Content : finalResponse.Message.Content,
		ToolCalls : toAgentToolCalls(finalResponse.Message.ToolCalls),
		IsToolUse : len(finalResponse.Message.ToolCalls) > 0 ,
	}

	return &llmResponse , nil

} 


	
func (oc *OllamaClientStruct) StreamFinalResponseWithStreaming(ctx context.Context,responseId string ,w *bufio.Writer ,messages []Message) (string, error) {
	
	isStreaming := true

		req := &api.ChatRequest{
		Model:    oc.ModelName,
		Messages: toOllamaMessages(messages), 
		Stream:   &isStreaming,
	}

	ollamaMessages := toOllamaMessages(messages) 
	fmt.Printf("========== MESSAGES ENVOYÉS À OLLAMA ==========\n")
    for i, msg := range ollamaMessages {
        fmt.Printf("[%d] Role: %s\n", i, msg.Role)
        fmt.Printf("    Content: %s\n", msg.Content)
        fmt.Printf("    Length: %d\n", len(msg.Content))
        fmt.Println("---")
    }
    fmt.Printf("===============================================\n")


	/*req := &api.ChatRequest{
        Model:  oc.ModelName,
        Stream: &isStreaming,
        Messages: []api.Message{
            {
                Role:    "user",
                Content: "Explique context engineering en 2 phrases",
            },
        },
    }*/

	var fullContent strings.Builder

	

	err := oc.Client.Chat(ctx, req, func(resp api.ChatResponse) error {
		if resp.Message.Content != "" {
			//fmt.Printf("Envoi du chunk: %s\n", resp.Message.Content)

			fmt.Printf("=== CHUNK ===\n")
			fmt.Printf("Content: '%s'\n", resp.Message.Content)
			fmt.Printf("Done: %v\n", resp.Done)
			fmt.Printf("=============\n")

			fullContent.WriteString(resp.Message.Content)
			TanstackSendSSE(w, StreamChunk{
			BaseStreamChunk: BaseStreamChunk{
				ID:        responseId,
				Type:      "content",
				Model:     oc.ModelName,
				Timestamp: time.Now().UnixMilli(),
			},
			Delta:   resp.Message.Content,
			Content: fullContent.String(),
			Role:    "assistant",
			})        
		}
        return nil
	})

	if err != nil {

        TanstackSendSSE(w, ErrorStreamChunk{
			BaseStreamChunk: BaseStreamChunk{
				ID:        responseId,
				Type:      "error", 
				Model:     oc.ModelName,
				Timestamp: time.Now().UnixMilli(),
			},
			Error: &StreamError{
				Message: err.Error(),
			},
		})
		
		fmt.Printf("error calling Chat function : %v\n", err)
		return  "" , err 
	}
	SendDone(w)
	/*TanstackSendSSE(w, DoneStreamChunk{
			BaseStreamChunk: BaseStreamChunk{
				ID:        responseId,
				Type:      "done", 
				Model:     oc.ModelName,
				Timestamp: time.Now().UnixMilli(),
			},
			FinishReadon: "stop",
	})*/

	return  fullContent.String() , nil

} 






