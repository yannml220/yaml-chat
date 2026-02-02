import { useAuth } from "../../contexts/AuthContext";
import { Navigate, Outlet   } from "@tanstack/react-router";

export const ProtectedRoute = () =>{
	const { isLoggedIn,isLoading} = useAuth() ;

	if (isLoading) return null;

	if( !isLoggedIn ){
		return <Navigate to="/signin" replace />
	}

	return <Outlet/> ;

}


