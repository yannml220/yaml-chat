import type { IUser } from "../../../types/models/user";
import { api } from "../api";

export const fetchMe = async (): Promise<IUser | null> => {
	const resp = await api.get('auth/me').json<{ data: IUser }>();
	console.log("returned user from fetch me : ", resp.data)
	return resp.data
}
