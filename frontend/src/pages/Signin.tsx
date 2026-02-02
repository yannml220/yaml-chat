import React, { useRef, useState}   from 'react'
import PageContent from '../components/layouts/PageContent'
import { useTheme } from '@emotion/react'
import { FcGoogle } from "react-icons/fc";
import Flex from '../components/base/Flex'
import {useQueryClient} from '@tanstack/react-query';
import {api} from '../lib/api/api';
import MyButton from '../components/base/Button';
import {BsApple, BsWindows, BsWindowSplit} from 'react-icons/bs';



function Signin() {

	const [ isLoading , setIsLoading ] = useState(false)

  	const theme = useTheme();

	const queryClient = useQueryClient()



	const oauthSigninHandler = async (e: React.MouseEvent<HTMLButtonElement>) => {
	  e.preventDefault();
	  setIsLoading(true)
	  try {
	 
		const resp = await api.get("auth/provider/google/signin").json<{ data: { url: string } }>();
		queryClient.invalidateQueries({queryKey : ['auth-me']})
	  	setIsLoading(false)
		if (resp.data?.url) {
			window.location.href = resp.data.url;
		}
	  } catch (error) {
		console.error("Erreur lors de l'initialisation du login Google:", error);
	  }
	};

    


  return (
   <PageContent display ="flex" justify = "center" align ="center" > 
   		<Flex
		style={{
			flexDirection:"column",
			gap:"0.5rem",
			padding :"2rem",
			width :"400px",
		}}
		>
			<MyButton
			onClick={(e)=>{oauthSigninHandler(e)}}
			hoverBgColor={theme.colors.background.ligther}
			style={{
				position:"relative",
				display:"flex",
				alignItems:"center",
				borderRadius :"10px",
				padding:"0.6rem 0",
				border : "1px solid lightgrey",
			}}
			>
				<Flex
				style={{
					alignItems:"center",
					position:"absolute",
					padding :"0 1rem",
					left : 10 ,
				}}
				>
					<FcGoogle size="22px"/>
				</Flex>
				<Flex
				style={{
					flex:1,
					justifyContent:"center",
					alignItems:"center"
				}}
				>
					<span
					style={{
						fontSize : "0.9rem",
						fontWeight :550,
					}}
					>
						Continue with Google
					</span>
				</Flex>
			</MyButton>
			
			<MyButton
			onClick={(e)=>{oauthSigninHandler(e)}}
			hoverBgColor={theme.colors.background.ligther}
			style={{
				position:"relative",
				display:"flex",
				alignItems:"center",
				borderRadius :"10px",
				padding:"0.6rem 0",
				border : "1px solid lightgrey",
			}}
			>
				<Flex
				style={{
					alignItems:"center",
					position:"absolute",
					padding :"0 1rem",
					left : 10 ,
				}}
				>
					<BsApple size="22px"/>
				</Flex>
				<Flex
				style={{
					flex:1,
					justifyContent:"center",
					alignItems:"center"
				}}
				>
					<span
					style={{
						fontSize : "0.9rem",
						fontWeight :550,
					}}
					>
						Continue with Apple
					</span>
				</Flex>
			</MyButton>


		</Flex>
    </PageContent>
 
  )
}

export default Signin




