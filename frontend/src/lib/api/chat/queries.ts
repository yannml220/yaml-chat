import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import type { IConversation, IMessage } from "../../../types/models/chat";
import { fetchConversationMessages, fetchUserConversations } from "./operations";



export function useFetchUserConversations(userId: string): UseQueryResult<IConversation[] | undefined, Error> {

	return useQuery({
		queryKey: ['conversations', { user_id: userId }],
		queryFn: () => fetchUserConversations(userId ?? "")

	})

}


export function useFetchConversationMessages(conversationId: string, userId: string): UseQueryResult<IMessage[] | undefined, Error> {

	return useQuery({
		queryKey: ['messages', { conversation_id: conversationId, user_id: userId }],
		queryFn: () => fetchConversationMessages(conversationId ?? "", userId ?? "")

	})

}
