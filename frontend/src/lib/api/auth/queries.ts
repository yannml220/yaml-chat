import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import type {IUser} from "../../../types/models/user";
import {fetchMe} from "./operations";


export function useFetchMe(enabled : boolean) : UseQueryResult<IUser | null,Error> {

	return useQuery({
	  queryKey: ['auth-me'],
	  queryFn: ()=> fetchMe(),
	  retry: false, 
	  //refetchOnWindowFocus: true,
	  staleTime: 1000 * 60 * 7,
	  /*refetchInterval: (query) => {
		return query.state.data ? 1000 * 60 * 10 : false
      },*/
	  throwOnError: false,
	  enabled
	})
}

