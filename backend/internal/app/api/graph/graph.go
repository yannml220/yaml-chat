package graph

import "time"



type Node struct {
	ID string `json:"id"`
    ModularityCommunity int `json:"modularity_community"`
	Freq int `json:"freq"`
    Pos string `json:"pos"`
	Betweenness float32 `json:"betweenness"`
	Gatewayness float32 `json:"gatewayness"`
    PageRank float32 `json:"page_rank"`
}


type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Weight int    `json:"weight"`
	Pmi float32    `json:"pmi"`
}




type GraphData struct {
	Modularity float32 `json:"modularity"`
    Nodes []Node `json:"nodes"`
    Edges []Edge `json:"edges"`
}




type Graph struct {
	Id             string     `json:"id"`
	ConversationId string     `json:"conversation_id"`
	Data           GraphData     `json:"data"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}


