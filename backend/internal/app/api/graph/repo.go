package graph

import (
	"context"
	"errors"
	"time"

	"github.com/yannml220/chat_agent_app/internal/app/db"
)

type GraphRepo interface {
	CreateGraph(ctx context.Context, conversationId string,graph *Graph) ( string, error)
	GetGraphById(ctx context.Context, id string) (*Graph, error)
	GetConversationGraphs(ctx context.Context , conversationId string) ( []Graph , error )
	UpdateGraph(ctx context.Context, graph *Graph, id string) error
	DeleteGraph(ctx context.Context, id string) error
}

type graphRepo struct {
	Db *db.Db
}

func NewRepo(db *db.Db) GraphRepo {
	return &graphRepo{
		Db: db,
	}

}

func (gr *graphRepo) GetGraphById(ctx context.Context, id string) (*Graph, error){
	sql := "select id , conversation_id ,data , created_at from graphs where id = $1"
	graph := Graph{}

	if err := gr.Db.GetDb().QueryRow(ctx, sql, id).Scan(&graph.Id, &graph.ConversationId, &graph.Data, &graph.CreatedAt); err != nil {

		return nil, err

	}

	return &graph, nil

}




func  (gr *graphRepo) GetConversationGraphs(ctx context.Context , conversationId string) ( []Graph , error ) {

	sql := "select graphs.id , graphs.conversation_id, graphs.data ,graphs.created_at from graphs inner join conversations on graphs.conversation_id = conversations.id  where conversations.id = $1" 

	rows, err := gr.Db.GetDb().Query(ctx,sql,conversationId) 

	if err != nil {
		return nil, err
	}

	var graphs []Graph

	defer rows.Close()

	for rows.Next() {
		var id string
		var conversationId string
		var data GraphData
		var createdAt time.Time


		err = rows.Scan(&id,&conversationId ,&data,&createdAt)
		
		if err != nil {
			return nil , errors.New("here is the error") 
		}

		newGraph := Graph {
			Id: id,
			ConversationId: conversationId,	
			Data: data,
			CreatedAt: &createdAt,
		}

		graphs = append(graphs, newGraph)
	}

	err = rows.Err()
	 if err != nil {
		 return nil , err 
	}


	return graphs , nil


}








func (gr *graphRepo) CreateGraph(ctx context.Context, conversationId string,graph *Graph) (string, error) {
	sql := "insert into graphs (conversation_id,data) values ($1, $2 )  returning id"
	var ret string

	err := gr.Db.GetDb().QueryRow(ctx, sql, conversationId, graph.Data).Scan(&ret)

	if err != nil {
		return "", err
	}

	return ret, nil

}



func (gr *graphRepo) UpdateGraph(ctx context.Context, graph *Graph, id string) error {

	sql := `
        UPDATE graphs 
        SET 
            data = COALESCE($2::jsonb, data),
            updated_at = CURRENT_TIMESTAMP 
        WHERE id = $1`

	_, err := gr.Db.GetDb().Exec(ctx, sql,
		id,
		graph.Data,
	)

	return err
}


func (gr *graphRepo) DeleteGraph(ctx context.Context, id string) error {
	sql := "delete from graphs where id = $1"

	if _, err := gr.Db.GetDb().Exec(ctx, sql, id); err != nil {

		return err
	}
	return nil
}
