import { ArrowUp, PanelLeft, Plus } from "lucide-react"
import Avatar from "../components/base/Avatar"
import Flex from "../components/base/Flex"
import Img from "../components/base/Img"
import ChatAgent, {ChatWindow} from "../components/monoliths/ChatAgent"
import { useAuth } from "../contexts/AuthContext"
import MyButton from "../components/base/Button"
import { useTheme } from "@emotion/react"
import { useFetchUserConversations } from "../lib/api/chat/queries"
import { useRef, useState } from "react"
import { Link, useNavigate, useParams } from '@tanstack/react-router'
import styled from "@emotion/styled"
import {useInitConversation} from "../lib/api/chat/mutations"
import {useResizeTextarea} from "../hooks/useResizeTextarea"



interface HomeProps {

}



const Home = ({ }: HomeProps) => {

	const { user } = useAuth()
	const params = useParams({ strict: false }) as any
	const conversationId = params?.conversationId || "";
	console.log("CONV ID :", conversationId)
	const { data: conversations, isPending: isFetchUserConversationsPending, error: fetchUserConversationsError } = useFetchUserConversations(user?.id ?? "");

	const theme = useTheme()


	return (
		<Flex
			isBorderBox
			style={{
				backgroundColor: "rgb(243 244 246)",
				minHeight: "100%",
				height: "100vh",
				overflowX: "hidden",
				overflowY: "hidden",
				position: "fixed",
				top: 0,
				left: 0,
				right: 0,
				bottom: 0,

			}}
		>
			<Flex
				style={{
					width: "260px",
					borderRight: "1px solid lightgrey",
					flexDirection: "column",

				}}
			>
				<Flex
					style={{
						height: "60px",
						padding: "0.5rem",
						alignItems: "center",
						justifyContent: "flex-end",
					}}
				>
					<MyButton
						onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
						}}
						hoverBgColor={`${theme.colors.background.superlightgrey}`}
						style={{
							height: "fit-content",
							width: "fit-content",
							display: "flex",
							justifyContent: "center",
							alignItems: "center",
							gap: "0.3rem",
							padding: "0.5rem",
							color: "black",
						}}
					>
						<PanelLeft size="18px" color="black" strokeWidth={1.2} />
					</MyButton>
				</Flex>

				<Flex
					style={{
						padding: "0.5rem 1rem",
						//flex : 1,
					}}
				>
					
					<ConversationLink
						style={{
							height: "fit-content",
							borderRadius: "13px",
							width: "100%",
							display: "flex",
							alignItems: "center",
							gap: "0.8rem",
							padding: "0.5rem 1rem",
							color: "black",
						}}
						to="/"
					>
						<Plus size="18px" color="black" strokeWidth={1.2} />
						<span
							style={{
								fontWeight: 500,
								fontSize: "0.9rem",
							}}
						>
							New conversation
						</span>

					</ConversationLink>


				</Flex>


				<div
				style={{
					padding: "1rem",
					height: "150px",
				}}
				>
					<Flex
					align="center"
					style={{
						//minHeight: "70px",
					}}
					>
						<span
							style={{
								fontSize: "0.85rem",
								padding: "0 1rem",
								fontWeight:590,
							}}
						>
							Tags
						</span>

						<MyButton
						onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
						}}
						hoverBgColor={`${theme.colors.background.superlightgrey}`}
						style={{
							borderRadius:"50px",
							marginLeft:"auto",
							height: "fit-content",
							width: "fit-content",
							display: "flex",
							justifyContent: "center",
							alignItems: "center",
							gap: "0.3rem",
							padding: "0.3rem",
							color: "black",
							}}
						>
							<Plus size="18px" color="grey" strokeWidth={1.2} />
						</MyButton>

					</Flex>


				</div>

									

					<>

						<Flex
							style={{
								flexDirection: "column",
								padding: "1rem",
								gap: "1rem",
								flex: 1,
								overflowY:"auto",
								scrollbarWidth: "thin" ,
								scrollbarColor: `${theme.colors.background.ligth} transparent` ,
					
							}}
						>	
							<div
								style={{
									fontWeight:590,
									fontSize: "0.85rem",
									padding: "0 1rem",
								}}
							>
								Discussions
							</div>

							<div
							style={{
								padding:"0 0.5rem",
							}}

							>
								
							{


								conversations?.map((conversation) => (
									<ConversationLink
										key={conversation.id}
										to="/chat/$conversationId"
										params={{ conversationId: conversation.id }}
										isSelected={conversation.id === conversationId}
									>
										<div
											style={{
												display: "inline-block",
												wordBreak: "break-word",
												width: "100%",
												whiteSpace: "nowrap",
												overflow: "hidden",
												textOverflow: "ellipsis",
												padding: "0.5rem 0",
												borderRadius: "13px",

											}}
										>
											<span style={{
												marginRight: "auto",
												whiteSpace: "nowrap",
												color: theme.colors.font.dark,
												fontSize: "0.80rem",
												fontWeight: 490,
												width: "100%",
												lineHeight: "1.25rem",
												borderRadius: "13px",
											}}>{conversation.title}</span>
										</div>
									</ConversationLink>
								))
							}
							</div>

							
						</Flex>
					</>

			</Flex>

			<Flex
				style={{
					flex: 1,
					flexDirection: "column",
				}}
			>
				
				<Flex
				style={{
					alignItems: "center",
					height: "60px",
					borderBottom: "1px solid lightgrey",
					padding: "0 2rem",
				}}
				>

					<Avatar style={{
						height: "35px",
						width: "35px",
						borderRadius: "50%",
						marginLeft: "auto",
					}} >
						<Img style={{
							borderRadius: "50%",

						}}
							src={user?.avatar ?? ""}
							height="100%" width="100%"
						/>
					</Avatar>

				</Flex>

				<Flex
				style={{
					flex:1,

				}}
				>
					
				{
					conversationId === "" ? 
						<Flex
						style={{
							flexDirection:"column",
							padding:"2rem",
							gap:"1.5rem",
							margin:"auto",
							flex:1,
							maxWidth:"700px",
						}}
						>
							<Flex
							direction="column"
							style={{
								gap:"0.4rem",
								padding:"0.4rem 0.8rem",

							}}
							>
								<Flex
								align="center"
								style={{
									flex:1,
								}}
								>
									<span
									style={{
										fontSize:"1.6rem",
										fontWeight:400,
									}}
									>
										{`Hey ${user?.name.split(" ")[0]} ,`}
									</span>

								</Flex>

								<Flex
								align="center"
								style={{
									flex:1,
								}}
								>
									<span
									style={{
										fontSize:"2.0rem",
										fontWeight:500,
									}}
									>
										What should we learn today ?
									</span>

								</Flex>

							</Flex>

							<WelcomeChatInput 
							userId={user?.id ?? ""}
							/>
																				
							
						</Flex>
					:
						
					<>

						<Flex
							style={{
								flex: 1,
								justifyContent: "center",
								background: "white",
							}}
						>
							<Flex
								style={{
									//maxWidth:"700px",
									flex: 1,
									position: "relative",
									paddingBottom: "1rem",
								}}
							>
								<ChatAgent
									conversationId={conversationId}
									inputLateralSpace="1.5rem"
									style={{
										//maxWidth:"700px",
										position: "absolute",
										//right: 15,
										bottom: "10px",
										height: "100%",
										width: "100%",
										zIndex: 100,
									}}
								/>
							</Flex>
							<Flex
								style={{
									flexDirection: "column",
									//maxWidth:"700px",
									flex: 1,
									position: "relative",
									paddingBottom: "1rem",
									borderLeft: "1px solid lightgrey",
								}}
							>
								{/* Placeholder or secondary content */}
							</Flex>

						</Flex>

					</>

				}

				</Flex>


			</Flex>

		</Flex>
	)

}



export default Home




const ConversationLink = styled(Link) <{ isSelected?: boolean }>`
  padding: 0 1rem;
  display: flex;
  justify-content: flex-start;
  font-family: system-ui, Segoe UI, Roboto, Helvetica, Arial, sans-serif;
  font-size: .875rem;
  line-height: 1.4;
  font-weight: 430;
  color: black;
  border-radius: 13px;
  text-decoration: none;
  background: ${props => props.isSelected ? props.theme.colors.background.superlightgrey : 'transparent'};
  
  &:hover {
    background: ${props => props.theme.colors.background.superlightgrey};
  }
`;








interface WelcomeChatInputProps {
	userId: string;
}

const WelcomeChatInput = ({ userId }: WelcomeChatInputProps) => {
	const theme = useTheme();
	const navigate = useNavigate();
	const [inputVal, setInputVal] = useState('');
	const textareaRef = useRef<HTMLTextAreaElement>(null);
	const { mutate: initChat, isPending } = useInitConversation();

	useResizeTextarea(textareaRef, inputVal, 220, true, 44);

	const handleSubmit = (e: React.FormEvent) => {
		console.log("test log")
		e.preventDefault();
		if (inputVal.trim() && !isPending) {
			initChat({ query: inputVal, userId }, {
				onSuccess: (convId) => {
					if (convId) {
						navigate({
							to: '/chat/$conversationId',
							params: { conversationId: convId },
							search: { q: inputVal, init: true } as any
						});
					}
				}
			});
		}
	};

	return (
		<Flex
			direction="column"
			style={{
				boxSizing: "border-box",
				background: "inherit",
				display: "flex",
				paddingBottom: "0.5rem",
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
						if (e.key === 'Enter' && !e.shiftKey) {
							e.preventDefault();
							handleSubmit(e as any);
						}
					}}
					ref={textareaRef}
					value={inputVal}
					onChange={(e) => setInputVal(e.currentTarget.value)}
					placeholder="Ask Anything..."
					style={{
						width: "100%",
						height: "44px",
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
				<Flex style={{ gap: "0.2rem" }}>

					<Flex style={{ marginLeft: "auto", gap: "0.2rem" }}>

						<MyButton
							onClick={handleSubmit}
							disabled={isPending || !inputVal.trim()}
							style={{
								borderRadius: "50px",
								display: "flex",
								justifyContent: "center",
								alignItems: "center",
								padding: "0.4rem",
								background: isPending || !inputVal.trim() ? "grey" : "black",
							}}
							bgColor="transparent"
						>
							<ArrowUp strokeWidth={1.8} color="white" size="18px" />
						</MyButton>

					</Flex>

				</Flex>

			</Flex>
		</Flex>
	);
};
