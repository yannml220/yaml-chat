import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import type { IGraph } from "../../../types/models/graph";
import { fetchConversationGraphs, fetchGraphById } from "./operations";



export function useFetchGraphById(graphId: string, userId: string): UseQueryResult<IGraph | undefined, Error> {
    return useQuery({
        queryKey: ['graph', { graph_id: graphId, user_id: userId }],
        queryFn: () => fetchGraphById(graphId, userId),
        enabled: !!graphId && !!userId
    })
}


export function useFetchConversationGraphs(conversationId: string, userId: string): UseQueryResult<IGraph[] | undefined, Error> {
    return useQuery({
        queryKey: ['graphs', { conversation_id: conversationId, user_id: userId }],
        queryFn: () => fetchConversationGraphs(conversationId, userId),
        enabled: !!conversationId && !!userId
    })
}
