import { useMutation, useQueryClient } from "@tanstack/react-query";
import { initConversation } from "./operations";

export const useInitConversation = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: ({ query, userId }: { query: string; userId: string }) => initConversation(query, userId),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['conversations', { user_id: variables.userId }] });
        }
    });
};
