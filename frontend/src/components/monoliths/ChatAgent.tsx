import { useLayoutEffect, useMemo, useRef, useState } from "react";
import Flex from "../base/Flex"
import { useTheme } from "@emotion/react";
import MarkdownRenderer from "./MarkdownRenderer";
import { ArrowUp, Ellipsis, HistoryIcon, Mic, Plus, SquarePen } from "lucide-react";
import { useResizeTextarea } from "../../hooks/useResizeTextarea";
import { useChat, fetchServerSentEvents, stream } from "@tanstack/ai-react";
import { useAuth } from "../../contexts/AuthContext";
import { useFetchConversationMessages, useFetchUserConversations } from "../../lib/api/chat/queries";
import { type UIMessage } from "@tanstack/ai";
import { useQueryClient } from "@tanstack/react-query";
import { getOrCreateDeviceId } from "../../lib/api/api";
import MyButton from "../base/Button";
import Spinner from "../base/Spinner";




interface ChatAgentProps extends React.HTMLAttributes<HTMLDivElement> {
	conversationId: string,
	inputLateralSpace?: string,


}

const ChatAgent = ({ conversationId, inputLateralSpace = "0px", ...RestProps }: ChatAgentProps) => {

	const { user } = useAuth()

	const { data: initialMessages, isPending: isFetchConversationMessagesPending, error: fetchConversationMessagesError } = useFetchConversationMessages(conversationId ?? "", user?.id ?? "");

	const theme = useTheme()

	const initialMessagesAsUiMessages = useMemo(() => {
		return initialMessages?.map((message) => ({
			id: message.id,
			role: message.role as "system" | "user" | "assistant",
			parts: [
				{
					type: "text" as const,
					content: message.content,
				},
			],
		})) as UIMessage[] | undefined


	}, [initialMessages, conversationId])




	return (

		<Flex
			direction="column"
			{...RestProps}

		>
			<Flex
				align="center"
				style={{
					padding: "0.5rem",
					borderBottom: "1px solid #d8d8dfc7",

				}}
			>
				<Flex
					style={{
						padding: "0.3rem",
						gap: "0.2rem",
					}}
				>
				</Flex>

				<Flex
					justify="flex-end"
					align="center"
					style={{
						padding: "0.3rem",
						gap: "0.2rem",
						flex: 1,
					}}
				>



					<MyButton
						style={{
							display: "flex",
							justifyContent: "center",
							alignItems: "center",
							padding: "0.4rem",
							color: "black",
						}}
						bgColor="transparent"
						hoverBgColor={`${theme.colors.background.superlightgrey}`}
					>
						<Ellipsis strokeWidth={1.8} color="black" size="18px" />
					</MyButton>

				</Flex>
			</Flex>
			{
				isFetchConversationMessagesPending ?

					<Spinner style={{
						margin: "auto",

					}} color="black" width={18} height={18} />
					:

					<ChatWindow
						key={conversationId}
						userId={user?.id ?? ""}
						conversationId={conversationId ?? ""}
						initialMessages={initialMessagesAsUiMessages ?? []}
						inputLateralSpace={inputLateralSpace}

					/>

			}

		</Flex>

	)

}


export default ChatAgent





const typingKeyframes = `
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 8px 12px;
  background: #f0f0f0;
  border-radius: 16px;
  width: fit-content;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #999;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.5;
  }
  30% {
    transform: translateY(-10px);
    opacity: 1;
  }
}
`

interface ChatWindowProps {
	userId: string,
	conversationId: string,
	hideWindow?: boolean,
	initialMessages: UIMessage[],

	inputLateralSpace?: string,

}


export const ChatWindow = ({ userId, hideWindow = false, inputLateralSpace = "0px", conversationId, initialMessages }: ChatWindowProps) => {

	const theme = useTheme()
	const scrollRef = useRef<HTMLDivElement>(null);
	const messagesEndRef = useRef<HTMLDivElement>(null);
	const textareaRef = useRef<HTMLTextAreaElement>(null);

	const [inputVal, setInputVal] = useState<string>('');
	// Track the active conversation ID for when we receive it from metadata
	const [activeConversationId, setActiveConversationId] = useState<string>(conversationId);
	useResizeTextarea(textareaRef, inputVal, 220, true, 44)



	console.log("voici les initial messages :\n", initialMessages)

	// Custom SSE connection that intercepts metadata events for navigation
	const connectionWithNavigation = stream((messages, data) => {
		const baseConnection = fetchServerSentEvents(
			() => `http://localhost:5000/api/v1/chat/streaming`,
			() => ({
				method: "POST",
				credentials: "include",
				headers: {
					'X-Device-ID': getOrCreateDeviceId(),
				},
				body: {
					message: inputVal,
					user_id: userId,
					conversation_id: conversationId,
				},
			})
		);

		// Return async generator that intercepts events
		return (async function* () {
			// Get the stream from base connection
			const stream = baseConnection.connect(messages, data);

			// Intercept and filter events
			for await (const chunk of stream) {
				// Type cast to check for custom metadata events
				const anyChunk = chunk as any;
				console.log('Received chunk:', anyChunk?.type, anyChunk);

				// Handle metadata events for navigation
				if (anyChunk?.type === 'metadata' && anyChunk?.content?.conversation_id) {
					const newConversationId = anyChunk.content.conversation_id;
					// Only update if we don't already have a conversation ID (new chat)
					if (!conversationId || conversationId === '') {
						console.log('Setting active conversation:', newConversationId);
						// Update state to show the message window
						setActiveConversationId(newConversationId);
						// Update URL without triggering component remount
						window.history.replaceState(null, '', `/chat/${newConversationId}`);
					}
					// Don't yield metadata events to chat UI
					continue;
				}

				// Yield all other events to the chat system
				console.log('Yielding chunk to chat UI:', anyChunk?.type);
				yield chunk;
			}
			console.log('Stream ended');
		})();
	});

	const { clear, messages, sendMessage, isLoading } = useChat({
		initialMessages: initialMessages,
		connection: connectionWithNavigation,
	});



	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		if (inputVal.trim() && !isLoading) {
			sendMessage(inputVal);
			setInputVal('');

			const queryClient = useQueryClient();
			const queryKey = ['conversations', { user_id: userId }]
			queryClient.invalidateQueries({ queryKey })
		}
	}



	useLayoutEffect(() => {
		messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
	}, [messages]);


	// Show message window if we have an active conversation (either from prop or from streaming)
	const shouldShowWindow = !hideWindow || activeConversationId !== '';

	return (
		<>
			{
				shouldShowWindow &&


				<Flex
					direction="column"
					style={{
						padding: "1rem 4rem",
						flex: 1,
						overflowY: "auto",
						gap: "0.7rem",
					}}
					ref={scrollRef}
				>

					<style>{typingKeyframes}</style>

					{
						messages?.map((message: UIMessage) => {
							const firstPart = message?.parts[0]

							if (firstPart?.type !== "text") return null

							return message?.role === "user" ? (
								<UserBubble>
									{firstPart?.content}
								</UserBubble>
							) : (
								<AssistantBubble content={firstPart?.content} />
							)
						})

					}

					{isLoading && (
						<Flex style={{ gap: "0.3rem", alignItems: "center" }}>
							<div className="typing-indicator">
								<span></span>
								<span></span>
								<span></span>
							</div>
						</Flex>
					)}
					<div ref={messagesEndRef} />
				</Flex>


			}


			<Flex
				direction="column"
				style={{
					boxSizing: "border-box",
					background: "inherit",
					display: "flex",
					paddingBottom: "0.5rem",
					paddingLeft: inputLateralSpace,
					paddingRight: inputLateralSpace,
				}}
			>
				<Flex
					direction="column"
					style={{
						padding: "0 0.5rem 0.5rem 0.5rem",
						borderRadius: "10px",
						border: `1px solid ${theme.colors.background.ligth} `,
						width: "100%",
						background: "white",
						boxShadow: "rgba(31, 34, 37, 0.12) .4px .4px .4px .4px",


					}}
				>


					<textarea
						onKeyDown={(e) => {
							//if (e.key === 'Enter' && e.shiftKey) {
							//e.preventDefault(); 
							// return;
							//}

							if (e.key === 'Enter' && !e.shiftKey) {
								e.preventDefault();
								handleSubmit(e as any);
							}
						}}

						ref={textareaRef}
						value={inputVal}
						onChange={(e) => setInputVal(e.currentTarget.value)}

						style={{
							width: "100%",
							height: "44px !important",
							textDecoration: "none",
							padding: "0.5rem",
							resize: "none",
							background: "inherit",
							fontSize: "0.9rem",
							fontFamily: "inherit",
							border: "none",
							boxSizing: "border-box",
							outline: 'none',
							lineHeight: "1.6",
							display: "block",
						}}
					/>



					<Flex
						style={{
							gap: "0.2rem",
						}}
					>
						<MyButton
							style={{
								display: "flex",
								alignItems: "center",
								padding: "0.4rem",
								color: "black",
							}}
							bgColor="transparent"
							hoverBgColor={`${theme.colors.background.superlightgrey}`}
						>
							<Plus strokeWidth={1.8} color="black" size="18px" />
						</MyButton>

						<Flex
							style={{
								marginLeft: "auto",
								gap: "0.2rem",
							}}

						>
							{
								inputVal !== '' ?
									<MyButton
										style={{
											borderRadius: "50px",
											display: "flex",
											justifyContent: "center",
											alignItems: "center",
											padding: "0.4rem",
											background: "black",
										}}
										bgColor="transparent"
										hoverBgColor={`${theme.colors.background.superlightgrey}`}
									>
										<ArrowUp strokeWidth={1.8} color="white" size="18px" />
									</MyButton>
									:

									<MyButton
										style={{
											display: "flex",
											justifyContent: "center",
											alignItems: "center",
											padding: "0.4rem",
											color: "black",
										}}
										bgColor="transparent"
										hoverBgColor={`${theme.colors.background.superlightgrey}`}
									>
										<Mic strokeWidth={1.8} color="black" size="18px" />
									</MyButton>

							}

						</Flex>

					</Flex>

				</Flex>
			</Flex>


		</>


	)



}








interface userBubbleProps {
	children?: React.ReactNode;
}

const UserBubble = ({ children }: userBubbleProps) => {

	return (
		<Flex
			style={{
				alignSelf: 'flex-end',
				maxWidth: '80%',
				padding: "0.8rem 1rem",
				background: "rgb(31 34 37 / 3%)",
				borderRadius: "15px",
				boxShadow: "rgba(31, 34, 37, 0.10) 0px 0px 0px 1px inset, rgba(0, 0, 0, 0.04) 0px 2px 8px -2px, rgba(0, 0, 0, 0.04) 0px 2px 4px -2px",
				whiteSpace: "pre-wrap",
				wordBreak: "break-word",
				lineHeight: "18px",
				fontWeight: 500,
				fontFamily: "Inter, system-ui, ui-sans-serif, -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji",

			}}
		>
			{children}
		</Flex>
	)
}






interface assistantBubbleProps {
	content?: string
}


const AssistantBubble = ({ content }: assistantBubbleProps) => {

	return (
		<Flex
			direction="column"
			style={{
				padding: "0.5rem 0",
			}}
		>
			<MarkdownRenderer
				content={content ?? ""}
				highlightTerms={[]}
				className="prose prose-slate max-w-none"
				style={{
					fontFamily: "Inter, system-ui, ui-sans-serif, -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji, Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji",
					//backgroundColor :`${theme.colors.background.ligther}`,
					height: "100%",
					fontSize: "0.90rem",
					fontWeight: 415,
					lineHeight: "23px",
					padding: "0.5rem",
				}}
			/>
		</Flex>
	)
}




/*
	<MyButton
						onClick={(e) => { e.preventDefault(); e.stopPropagation(); setCurrentConversationId(null) }}
						style={{
							display: "flex",
							justifyContent: "center",
							alignItems: "center",
							padding: "0.4rem",
							color: "black",
						}}
						bgColor="transparent"
						hoverBgColor={`${theme.colors.background.superlightgrey}`}
					>
						<SquarePen strokeWidth={1.4} color="black" size="18px" />
					</MyButton>

*/
