import api from '@/API/api';
import { tokenData } from '@/lib/types';

const loginService = async (email: string, password: string) => {
	try {
		console.log(email);
		const response = await api
			.post<tokenData>('/login', { email, password })
			.then((resp) => resp.data);

		if (response.message) {
			throw new Error(response.message);
		}

		return response;
	} catch (error: any) {
		if (error.response) {
			throw new Error(
				'Unauthorized access. Please check your credentials.',
			);
		}

		throw new Error(error.message);
	}
};

export default loginService;
