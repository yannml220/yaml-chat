import { AuthContext } from "../../contexts/AuthContext";
import { useLocation, useNavigate } from "@tanstack/react-router";
import {api} from "../../lib/api/api";
import {useQueryClient} from "@tanstack/react-query";
import {useEffect} from "react";
import {useFetchMe} from "../../lib/api/auth/queries";



export const AuthProvider = ({children}: {children:React.ReactNode}) =>{

	const queryClient = useQueryClient();
	const location = useLocation()
	const isLoginPage = location.pathname === '/signin'
	const navigate = useNavigate();

	const { data : user , isLoading : isFetchMePending  , error : fetchMeError } = useFetchMe(!isLoginPage)

	console.log("user :",user)
	console.log("user avatar :",user?.avatar)
	console.log("user name :",user?.name)
	console.log("user id :",user?.id)

	useEffect(() => {
		if (fetchMeError && !isLoginPage) {
		navigate({ to: "/signin", replace: true });
		}
	}, [fetchMeError, isLoginPage, navigate]);


	//console.log("user : ",user)

	const logout = () => {
		api.post('auth/signout').catch(() => {});
		queryClient.setQueryData(['auth-me'], null);
		navigate({ to: "/signin", replace: true });
  	};


	return (

		<AuthContext.Provider  value ={{
			user : user ?? null,
			isLoggedIn : !!user && !fetchMeError, 
			isLoading : isFetchMePending,
			logout ,
		}}>
			{!isFetchMePending && children}
		</AuthContext.Provider>

	)

}
