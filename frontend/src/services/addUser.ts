import api from '@/API/api';
import { UserFormType } from '@/components/PAGES/users/schema';
import { commonresponse } from '@/lib/types';

const addUser = async (data: UserFormType) => {
	try {
		const response = await api
			.post<commonresponse>('/user', data)
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

export default addUser;
