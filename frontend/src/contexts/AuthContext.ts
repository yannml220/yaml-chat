import { createContext , useContext } from "react";
import { type IUser} from "../types/models/user";



interface AuthContextType {
	user : IUser | null ,
	isLoggedIn : boolean ,
	isLoading : boolean ,
	logout : ()=>void,

} 


export const AuthContext = createContext<AuthContextType | null>({
	user : null ,
	isLoggedIn : false ,
	isLoading : false ,
	logout : ()=>{},
	

});


export const useAuth  = ()=>{
    const context = useContext(AuthContext);
    if (!context) throw new Error("AuthContext used outside of bounds");
    return context ;

}





