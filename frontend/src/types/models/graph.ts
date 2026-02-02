

export interface IGraphData {
	modularity? : number ,
	nodes : {
		id : string ,
		pos? : string ,
		betweenness? : number
		page_rank? : number
		freq? : number
		gatewayness? : number
		modularity_community? : number
	}[],
	edges : {
		source : string ,
		target : string ,
		weight? : number ,
		pmi? : number ,
	}[],
}


export type IGraph = {
	id :string ,
	conversation_id : string ,
	graph_data? : IGraphData ,
	created_at : string ,
	udpated_at : string ,
	deleted_at : string ,
}


