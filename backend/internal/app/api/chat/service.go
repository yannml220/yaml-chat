package chat

import (
	"context"
	"github.com/yannml220/chat_agent_app/internal/pkg/types"
)



type ChatService interface {

	GetUserConversations(ctx context.Context , user_id string ) ( []Conversation  , error )

	GetUserConversationsWithMessages(ctx context.Context , userId string ) ( []ConversationWithMessages  , error )

	GetConceptConversationsWithMessages(ctx context.Context , conversationId string ) ( []ConversationWithMessages  , error )

	GetConceptConversations(ctx context.Context , concept_id string ) ( []Conversation , error )

	GetConversationById(ctx context.Context , id string ) ( *Conversation , error )

	GetConversationMessages(ctx context.Context , conversationId string) ( []Message , error ) 

	GetConversationMessagesByRole(ctx context.Context , conversationId string,role string) ( []Message , error ) 

	GetMessageById(ctx context.Context , id string ) ( *Message , error )

	UpdateConversation(ctx context.Context ,  id string ) error  

	UpdateConversationTitle(ctx context.Context ,title string,  id string ) error  

	DeleteConversationById(ctx context.Context , id string ) error 

	CreateMessage(ctx context.Context ,conversation_id string  ,message *Message ) (string,error ) 

	CreateMessages(ctx context.Context , conversation_id string  , messages []Message ) ([]string,error ) 

	CreateUserConversation(ctx context.Context ,userId string, conversation *Conversation ) (string,error ) 

	DeleteMessageById(ctx context.Context , id string ) error 

	DeleteMessages(ctx context.Context , id []string ) error 

	DeleteUserConversations(ctx context.Context , userId string ) error 


	//DeleteAllMessagesBefore(ctx context.Context , anchorMessageId string ) error 

}



type ChatServiceImpl struct {
	types.Service
	repo ChatRepo
}

func NewService( repo ChatRepo ) ChatService {

	return &ChatServiceImpl {
		types.Service{}, 
		repo ,
	}

}





func (c* ChatServiceImpl) GetUserConversations(ctx context.Context , userId string ) ( []Conversation  , error ){
	return c.repo.GetUserConversations(ctx,userId)
}

func (c* ChatServiceImpl) GetUserConversationsWithMessages(ctx context.Context , userId string ) ( []ConversationWithMessages  , error ){
	return c.repo.GetUserConversationsWithMessages(ctx,userId)
}

func (c* ChatServiceImpl) GetConceptConversationsWithMessages(ctx context.Context , conversationId string ) ( []ConversationWithMessages  , error ){
	return c.repo.GetUserConversationsWithMessages(ctx,conversationId)
}


func (c* ChatServiceImpl) GetConceptConversations(ctx context.Context , conceptId string ) ( []Conversation , error ){
	return c.repo.GetConceptConversations(ctx,conceptId)
}



func (c* ChatServiceImpl) GetConversationById(ctx context.Context , id string ) ( *Conversation , error ) {
	return c.repo.GetConversationById(ctx,id)
}
func (c* ChatServiceImpl) GetConversationMessages(ctx context.Context , conversationId string) ( []Message , error ) {
	return c.repo.GetConversationMessages(ctx,conversationId)
}

func (c* ChatServiceImpl) GetConversationMessagesByRole(ctx context.Context , conversationId string,role string) ( []Message , error ){
	return c.repo.GetConversationMessagesByRole(ctx,conversationId,role)
}




func (c* ChatServiceImpl) GetMessageById(ctx context.Context , id string ) ( *Message , error ){
	return c.repo.GetMessageById(ctx,id)
}

func (c* ChatServiceImpl) UpdateConversation(ctx context.Context ,  id string ) error  {
	return c.repo.UpdateConversation(ctx,id)
}

func (c* ChatServiceImpl) UpdateConversationTitle(ctx context.Context ,title string , id string ) error  {
	return c.repo.UpdateConversationTitle(ctx,title,id)
}

func (c* ChatServiceImpl) DeleteConversationById(ctx context.Context , id string ) error {
	return c.repo.DeleteConversationById(ctx,id)

}


func (c* ChatServiceImpl) CreateMessage(ctx context.Context ,conversationIid string  ,message *Message ) (string,error ) {
	return c.repo.CreateMessage(ctx,conversationIid,message)
}


func (c* ChatServiceImpl) CreateMessages(ctx context.Context , conversationId string  , messages []Message ) ([]string,error ) {
	return c.repo.CreateMessages(ctx,conversationId,messages)

}



func (c* ChatServiceImpl) CreateUserConversation(ctx context.Context ,userId string, conversation *Conversation ) (string,error ) {
	return c.repo.CreateUserConversation(ctx,userId,conversation)
}

func (c* ChatServiceImpl) DeleteMessageById(ctx context.Context , id string ) error  {
	return c.repo.DeleteMessageById(ctx,id)

}

func (c* ChatServiceImpl) DeleteMessages(ctx context.Context , ids []string ) error {
	return c.repo.DeleteMessages(ctx,ids)

}

func (c* ChatServiceImpl) DeleteUserConversations(ctx context.Context , userId string ) error  {
	return c.repo.DeleteUserConversations(ctx,userId)

}




/*func (c* ChatServiceImpl) Chat(ctx context.Context , userMessage string , history []agent.Message, ) ( []Conversation  , error ){
	return c.repo.GetUserConversations(ctx,userId)
}*/

