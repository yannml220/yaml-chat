import { useMutation, useQueryClient } from "@tanstack/react-query";
import { initConversation, deleteConversationById } from "./operations";

export const useInitConversation = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: ({ query, userId }: { query: string; userId: string }) => initConversation(query, userId),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['conversations', { user_id: variables.userId }] });
        }
    });
};

export const useDeleteConversation = (userId: string) => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (conversationId: string) => deleteConversationById(conversationId, userId),
        onSuccess: (_, conversationId) => {
            queryClient.invalidateQueries({ queryKey: ['conversations', { user_id: userId }] });
            queryClient.invalidateQueries({ queryKey: ['messages', { conversation_id: conversationId, user_id: userId }] });
        }
    });
};
