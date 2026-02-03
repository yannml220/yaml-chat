package chat

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/yannml220/chat_agent_app/internal/app/db"
	"github.com/yannml220/chat_agent_app/internal/pkg/utils/agent"
)


type ChatRepo interface {

	GetUserConversations(ctx context.Context , user_id string ) ( []Conversation  , error )

	GetUserConversationsWithMessages(ctx context.Context , userId string ) ( []ConversationWithMessages  , error )

	GetConceptConversationsWithMessages(ctx context.Context , conversationId string ) ( []ConversationWithMessages  , error )

	GetConceptConversations(ctx context.Context , concept_id string ) ( []Conversation , error )

	GetConversationById(ctx context.Context , id string ) ( *Conversation , error )

	GetConversationMessages(ctx context.Context , conversationId string) ( []Message , error ) 
	
	GetConversationMessagesByRole(ctx context.Context , conversationId string,role string) ( []Message , error ) 

	GetMessageById(ctx context.Context , id string ) ( *Message , error )

	UpdateConversation(ctx context.Context ,  id string ) error  
	
	UpdateConversationTitle(ctx context.Context ,title string ,  id string ) (error ) 

	DeleteConversationById(ctx context.Context , id string ) error 

	CreateMessage(ctx context.Context ,conversation_id string  ,message *Message ) (string,error ) 

	CreateUserConversation(ctx context.Context ,userId string , conversation *Conversation ) (string,error ) 

	CreateMessages(ctx context.Context , conversationId string  , messages []Message ) ([]string,error ) 

	DeleteMessageById(ctx context.Context , id string ) error 

	DeleteMessages(ctx context.Context , id []string ) error 

	DeleteUserConversations(ctx context.Context , userId string ) error 


	//DeleteAllMessagesBefore(ctx context.Context , anchorMessageId string ) error 

}




type ChatRepoImpl struct {
	Db *db.Db
}



func NewRepo( db *db.Db )  ChatRepo {
	return &ChatRepoImpl{
		Db : db,

	}

}



func  ( c * ChatRepoImpl ) GetUserConversations(ctx context.Context , userId string ) ( []Conversation  , error ){

	sql := "select id ,user_id ,meta_data,title,created_at from conversations where user_id = $1 ORDER BY created_at DESC" 

	rows, err := c.Db.GetDb().Query(ctx,sql,userId) 

	if err != nil {
		log.Print("here is the ERROR " , err ) 
		return nil, err
	}

	var userConversations []Conversation

	defer rows.Close()

	for rows.Next() {
		var id string
		var userId string
		var metaData map[string]interface{}
		var title *string
		var createdAt *time.Time
		
		err = rows.Scan(&id, &userId,&metaData,&title,&createdAt)

		if err != nil {
			return nil , err 
		}

		newConversation := Conversation {
			Id: id,
			UserId: userId,	
			MetaData: metaData,
			Title: title,
			CreatedAt: createdAt,
		}

		userConversations = append(userConversations, newConversation)
	}

	err = rows.Err()
	 if err != nil {
		 return nil , err 
	}

	if len(userConversations) == 0 {
		return nil, nil
	}

	return userConversations , nil


}


func (c *ChatRepoImpl) GetUserConversationsWithMessages(ctx context.Context, userId string) ([]ConversationWithMessages, error) {
    sql := `
        SELECT 
            c.id, c.user_id, c.meta_data, c.created_at,
            m.id, m.role, m.content, m.tool_calls, m.created_at
        FROM conversations c
        LEFT JOIN messages m ON c.id = m.conversation_id
        WHERE c.user_id = $1
        ORDER BY c.created_at DESC, m.created_at ASC`

    rows, err := c.Db.GetDb().Query(ctx, sql, userId)
    if err != nil {
        return nil, errors.New("here is the error")
    }
    defer rows.Close()

    convMap := make(map[string]*ConversationWithMessages)
    var orderedIds []string

    for rows.Next() {
        var cID, cUserID string
        var cMeta map[string]interface{}
        var cCreatedAt *time.Time
        
        var mID, mRole, mContent *string
        var mToolCalls []agent.ToolCall
        var mCreatedAt *time.Time

        err := rows.Scan(
            &cID, &cUserID, &cMeta, &cCreatedAt,
            &mID, &mRole, &mContent, &mToolCalls, &mCreatedAt,
        )
        if err != nil {
            return nil, err
        }

        if _, ok := convMap[cID]; !ok {
            orderedIds = append(orderedIds, cID)
            convMap[cID] = &ConversationWithMessages{
                Id:        cID,
                UserId:    cUserID,
                MetaData:  cMeta,
                CreatedAt: cCreatedAt,
                Messages:  []Message{},
            }
        }

        if mID != nil {
            convMap[cID].Messages = append(convMap[cID].Messages, Message{
                Id:        *mID,
                Role:      *mRole,
                Content:   *mContent,
                ToolCalls: mToolCalls,
                CreatedAt: mCreatedAt,
            })
        }
    }

    result := make([]ConversationWithMessages, 0, len(orderedIds))
    for _, id := range orderedIds {
        result = append(result, *convMap[id])
    }

    return result, nil
}





func  ( c * ChatRepoImpl ) GetConceptConversations(ctx context.Context , conceptId string ) ( []Conversation , error ){

	sql := "select id ,user_id ,meta_data,created_at from conversations where concept_id = $1" 

	rows, err := c.Db.GetDb().Query(ctx,sql,conceptId) 

	if err != nil {
		log.Print("here is the ERROR " , err ) 
		return nil, err
	}

	var conceptConversations []Conversation

	defer rows.Close()

	for rows.Next() {
		var id string
		var userId string
		var conceptId string
		//var metaData map[string]interface{}
		var metaDataBytes []byte
		var createdAt *time.Time
		
		err = rows.Scan(&id, &userId,&conceptId,&metaDataBytes,&createdAt)

		if err != nil {
			return nil , errors.New("here is the error") 
		}

		newConversation := Conversation {
			Id: id,
			UserId: userId,	
			ConceptId: conceptId,	
			CreatedAt: createdAt,
		}

		if len(metaDataBytes) > 0 {
            err = json.Unmarshal(metaDataBytes, &newConversation.MetaData)
            if err != nil {
                return nil, err
            }
        }

		conceptConversations = append(conceptConversations, newConversation)
	}

	err = rows.Err()
	 if err != nil {
		 return nil , err 
	}

	return conceptConversations , nil


}



func (c *ChatRepoImpl) GetConceptConversationsWithMessages(ctx context.Context, conceptId string) ([]ConversationWithMessages, error) {
    sql := `
        SELECT 
            c.id, c.user_id, c.meta_data, c.created_at,
            m.id, m.role, m.content, m.tool_calls, m.created_at
        FROM conversations c
        LEFT JOIN messages m ON c.id = m.conversation_id
        WHERE c.concept_id = $1
        ORDER BY c.created_at DESC, m.created_at ASC`

    rows, err := c.Db.GetDb().Query(ctx, sql, conceptId)
    if err != nil {
        return nil, errors.New("here is the error")
    }
    defer rows.Close()

    convMap := make(map[string]*ConversationWithMessages)
    var orderedIds []string

    for rows.Next() {
        var cID, cUserID string
        var cMeta map[string]interface{}
        var cCreatedAt *time.Time
        
        var mID, mRole, mContent *string
        var mToolCalls []agent.ToolCall
        var mCreatedAt *time.Time

        err := rows.Scan(
            &cID, &cUserID, &cMeta, &cCreatedAt,
            &mID, &mRole, &mContent, &mToolCalls, &mCreatedAt,
        )
        if err != nil {
            return nil, err
        }

        if _, ok := convMap[cID]; !ok {
            orderedIds = append(orderedIds, cID)
            convMap[cID] = &ConversationWithMessages{
                Id:        cID,
                UserId:    cUserID,
                MetaData:  cMeta,
                CreatedAt: cCreatedAt,
                Messages:  []Message{},
            }
        }

        if mID != nil {
            convMap[cID].Messages = append(convMap[cID].Messages, Message{
                Id:        *mID,
                Role:      *mRole,
                Content:   *mContent,
                ToolCalls: mToolCalls,
                CreatedAt: mCreatedAt,
            })
        }
    }

    result := make([]ConversationWithMessages, 0, len(orderedIds))
    for _, id := range orderedIds {
        result = append(result, *convMap[id])
    }

    return result, nil
}



func  ( c * ChatRepoImpl ) GetConversationById(ctx context.Context , id string ) ( *Conversation , error ) {

	sql := "select id ,user_id,meta_data,created_at from conversations where id = $1" 

	conversation := Conversation{}
	
	if err := c.Db.GetDb().QueryRow(ctx,sql,id).Scan(&conversation.Id,&conversation.UserId ,&conversation.MetaData , &conversation.CreatedAt ) ; err != nil {
		
		return nil, err

	}
	return &conversation , nil


}




func  ( c * ChatRepoImpl ) GetConversationMessages(ctx context.Context , conversationId string) ( []Message , error ) {
	sql := "select id ,conversation_id,role, content,tool_calls,created_at from messages where conversation_id = $1 ORDER BY created_at ASC" 

	rows, err := c.Db.GetDb().Query(ctx,sql,conversationId) 

	if err != nil {
		log.Print("here is the ERROR " , err ) 
		return nil, err
	}

	var conversationMessages []Message

	defer rows.Close()

	for rows.Next() {
		var id string
		var conversationId string
		var role string
		var content string
		var toolCalls []agent.ToolCall
		var createdAt *time.Time
		
		err = rows.Scan(&id, &conversationId,&role,&content,&toolCalls,&createdAt)

		if err != nil {
			return nil , errors.New("here is the error") 
		}

		newMessage := Message {
			Id: id,
			ConversationId: conversationId,	
			Role: role,
			Content: content,
			ToolCalls: toolCalls,
			CreatedAt: createdAt,
		}

		conversationMessages = append(conversationMessages, newMessage)
	}

	err = rows.Err()
	 if err != nil {
		 return nil , err 
	}

	if len(conversationMessages) == 0 {
		return nil, nil
	}

	return conversationMessages , nil

}


func  ( c * ChatRepoImpl ) GetConversationMessagesByRole(ctx context.Context , conversationId string,role string) ( []Message , error ) {
	sql := "select id ,conversation_id,role, content,tool_calls,created_at from messages where conversation_id = $1 and role = $2  ORDER BY created_at ASC" 

	rows, err := c.Db.GetDb().Query(ctx,sql,conversationId,role) 

	if err != nil {
		log.Print("here is the ERROR " , err ) 
		return nil, err
	}

	var conversationMessages []Message

	defer rows.Close()

	for rows.Next() {
		var id string
		var conversationId string
		var role string
		var content string
		var toolCalls []agent.ToolCall
		var createdAt *time.Time
		
		err = rows.Scan(&id, &conversationId,&role,&content,&toolCalls,&createdAt)

		if err != nil {
			return nil , errors.New("here is the error") 
		}

		newMessage := Message {
			Id: id,
			ConversationId: conversationId,	
			Role: role,
			Content: content,
			ToolCalls: toolCalls,
			CreatedAt: createdAt,
		}

		conversationMessages = append(conversationMessages, newMessage)
	}

	err = rows.Err()
	 if err != nil {
		 return nil , err 
	}

	if len(conversationMessages) == 0 {
		return nil, nil
	}

	return conversationMessages , nil

}




func  ( c * ChatRepoImpl ) GetMessageById(ctx context.Context , id string ) ( *Message , error ) {

	sql := "select id ,conversation_id,role, content,tool_calls,created_at from messages where id = $1" 

	message := Message{}
	
	if err := c.Db.GetDb().QueryRow(ctx,sql,id).Scan(&message.Id,&message.ConversationId ,&message.Role, &message.Content,&message.ToolCalls,&message.CreatedAt) ; err != nil {
		
		return nil, err

	}
	return &message , nil

}




func  ( c * ChatRepoImpl ) CreateConceptConversation(ctx context.Context ,userId,conceptId , conversation *Conversation ) (string,error ) {

	sql := "insert into conversations (user_id,concept_id ,meta_data) values ($1, $2,$3) returning conversations.id"
	var ret string

	err := c.Db.GetDb().QueryRow(ctx , sql , userId ,conceptId,&conversation.MetaData).Scan(&ret)

	if err != nil {
		return "" , err 
	}
	return ret,nil 
}




func (c *ChatRepoImpl) CreateUserConversation(ctx context.Context, userId string, conversation *Conversation) (string, error) {

	sql := "insert into conversations (user_id,title,meta_data) values ($1, $2, $3) returning conversations.id"
	var ret string

	err := c.Db.GetDb().QueryRow(ctx, sql, userId, conversation.Title, &conversation.MetaData).Scan(&ret)

	if err != nil {
		return "", err
	}
	return ret, nil
}



func  ( c * ChatRepoImpl ) UpdateConversationTitle(ctx context.Context ,title string ,  id string ) (error ) {

	query := `UPDATE conversations SET title = coalesce($2,title), updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := c.Db.GetDb().Exec(ctx, query, id,&title)
	if err != nil {
		return err
	}

	return nil

}


func  ( c * ChatRepoImpl ) UpdateConversation(ctx context.Context ,  id string ) (error ) {

	query := `UPDATE conversations SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := c.Db.GetDb().Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil

}




func  ( c * ChatRepoImpl ) DeleteConversationById(ctx context.Context , id string ) error {

	sql := "delete from conversations where id = $1" 

	if _ , err := c.Db.GetDb().Exec(ctx , sql , &id ) ; err != nil {
		return err

	}
	return nil 

}




func  ( c * ChatRepoImpl ) DeleteUserConversations(ctx context.Context , userId string ) error {

	sql := "DELETE FROM conversations WHERE user_id = $1" 

	if _ , err := c.Db.GetDb().Exec(ctx , sql , &userId ) ; err != nil {
		return err

	}
	return nil

}



func  ( c * ChatRepoImpl ) DeleteConceptConversations(ctx context.Context , conceptId string ) error {

	sql := "DELETE FROM conversations WHERE concept_id = $1" 

	if _ , err := c.Db.GetDb().Exec(ctx , sql , &conceptId ) ; err != nil {
		return err

	}
	return nil

}




func  ( c * ChatRepoImpl ) CreateMessage(ctx context.Context ,conversationId string  ,message *Message ) (string,error ) {

	sql := "insert into messages (conversation_id,role, content,tool_calls) values ($1, $2,$3,$4) returning messages.id"
	var ret string

	err := c.Db.GetDb().QueryRow(ctx , sql , conversationId, &message.Role, &message.Content,&message.ToolCalls ).Scan(&ret)

	if err != nil {
		return "" , err 
	}
	return ret,nil 

}




func  ( c * ChatRepoImpl ) CreateMessages(ctx context.Context , conversationId string  , messages []Message ) ([]string,error ) {

    batch := &pgx.Batch{}
    
    sql := "insert into messages (conversation_id,role, content,tool_calls,created_at) values ($1, $2,$3,$4,$5) returning messages.id"
    
    for _, message := range messages {
        batch.Queue(sql,conversationId, message.Role,message.Content, message.ToolCalls,message.CreatedAt)
    }
    
    results := c.Db.GetDb().SendBatch(ctx, batch)
    defer results.Close()
    
    ids := make([]string, 0, len(messages))
    for range messages {
        var id string
        err := results.QueryRow().Scan(&id)
        if err != nil {
            return nil, err
        }
        ids = append(ids, id)
    }
    
    return ids, nil


}




func  ( c * ChatRepoImpl ) DeleteMessageById(ctx context.Context , id string ) error {

	sql := "delete from messages where id = $1" 

	if _ , err := c.Db.GetDb().Exec(ctx , sql , &id ) ; err != nil {
		return err

	}
	return nil 

}



func  ( c * ChatRepoImpl ) DeleteMessages(ctx context.Context , ids []string ) error {
	batch := &pgx.Batch{}
	for _, id := range ids {
		batch.Queue("DELETE FROM messages WHERE id = $1", id)
	}

	br := c.Db.GetDb().SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec() 
	
	if err != nil {
        return  err
    }
	return nil
}




/*func  ( c * ChatRepoImpl ) DeleteAllMessagesBefore(ctx context.Context , anchorMessageId string ) error {

}*/


