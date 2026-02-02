import type { IGraph } from "../../../types/models/graph";
import { api } from "../api";


export interface CreateGraphPayload {
	conversation_id: string;
	window_size: number;
}



export const createGraph = async (payload: CreateGraphPayload, userId: string): Promise<string | undefined> => {
	try {
		const resp = await api.post('graph', {
			json: payload,
			searchParams: { user_id: userId }
		}).json<{ data: { id: string } }>();

		return resp.data.id;
	} catch (error) {
		console.error("Error creating the graph:", error);
	}
}



export const deleteGraph = async (graphId: string, userId: string): Promise<string | undefined> => {
	try {
		const resp = await api.delete(`graph/${graphId}`, {
			searchParams: { user_id: userId }
		}).json<{ data: { id: string } }>();

		return resp.data.id;
	} catch (error) {
		console.error("Error deleting the graph:", error);
	}
}




export const fetchGraphById = async (graphId: string, userId: string): Promise<IGraph | undefined> => {
	try {
		const resp = await api.get(`graph/${graphId}`, {
			searchParams: { user_id: userId }
		}).json<{ data: { graph: IGraph } }>();

		return resp.data.graph;
	} catch (error) {
		console.error("Error fetching the graph by id:", error);
	}
}



export const fetchConversationGraphs = async (conversationId: string, userId: string): Promise<IGraph[] | undefined> => {
	try {
		const resp = await api.get(`graph`, {
			searchParams: { user_id: userId, conversation_id: conversationId },
		}).json<{ data: { graphs: IGraph[] } }>();

		console.log("returned graphs from api :", resp.data);
		return resp.data.graphs;
	} catch (error) {
		console.error("Error fetching conversation graphs:", error);
	}
}




