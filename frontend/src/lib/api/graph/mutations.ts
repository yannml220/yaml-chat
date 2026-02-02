import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createGraph, deleteGraph, type CreateGraphPayload } from "./operations";



export function useCreateGraph(userId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (payload: CreateGraphPayload) => createGraph(payload, userId),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({
                queryKey: ['graphs', { conversation_id: variables.conversation_id, user_id: userId }]
            });
        }
    });
}


export function useDeleteGraph(userId: string, conversationId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (graphId: string) => deleteGraph(graphId, userId),
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: ['graphs', { conversation_id: conversationId, user_id: userId }]
            });
        }
    });
}
