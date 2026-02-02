


export type IConversation = {
	id :string ,
	user_id : string ,
	concept_id? : string  ,
	meta_data? : string  ,
	title? : string  ,
	created_at : string ,
	udpated_at : string ,
	deleted_at : string ,

}


export type IToolCall = {
	id :string,
	name :string,
	payload :Record<string,any>,

}

export type IMessage = {
	id :string ,
	conversation_id : string ,
	role : string  ,
	content : string  ,
	tool_calls? : IToolCall[] ,
	created_at : string ,
	udpated_at : string ,
	deleted_at : string ,

}




