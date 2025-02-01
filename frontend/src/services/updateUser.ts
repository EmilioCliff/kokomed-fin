import api from '@/API/api';
import { updateUserType } from '@/lib/types';

const updateUser = async (data: updateUserType) => {
	console.log(data);
	try {
		const response = await api
			.patch<updateUserType>(`/user/${data.id}`, data)
			.then((resp) => resp.data);

		if (response.message) {
			throw new Error(response.message);
		}

		return response;
	} catch (error: any) {
		if (error.response) {
			throw new Error(error.response.data.message);
		}

		throw new Error(error.message);
	}
};

export default updateUser;
