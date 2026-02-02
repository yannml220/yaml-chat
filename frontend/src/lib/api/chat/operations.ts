import type { IConversation, IMessage } from "../../../types/models/chat";
import { api } from "../api";


export const fetchUserConversations = async (userId: string): Promise<IConversation[] | undefined> => {
	try {
		const resp = await api.get('chat', {
			searchParams: { user_id: userId }
		}).json<{ data: { conversations: IConversation[] } }>();

		console.log("returned conversations from api :", resp.data);
		return resp.data.conversations;
	} catch (error) {
		console.error("Error fetching user conversations:", error);
	}
}



export const fetchConversationMessages = async (conversationId: string, userId: string): Promise<IMessage[] | undefined> => {
	try {
		const resp = await api.get(`chat/${conversationId}`, {
			searchParams: { user_id: userId }
		}).json<{ data: { messages: IMessage[] } }>();

		console.log("returned messages from api :", resp.data);
		return resp.data.messages;
	} catch (error) {
		console.error("Error fetching conversation messages:", error);
	}
}
