import api from '@/API/api';
import { commonresponse } from '@/lib/types';
import { ResetPassowordFormType } from '@/components/PAGES/login/schema';

const resetPasswordService = async (data: ResetPassowordFormType) => {
	try {
		const response = await api
			.patch<commonresponse>(`/user/reset-password/${data.token}`, {
				newPassword: data.newPassword,
			})
			.then((resp) => resp.data);

		if (response.message) {
			throw new Error(response.message);
		}

		return response;
	} catch (error: any) {
		console.log(error);
		if (error.response) {
			throw new Error(error.response.data.message);
		}

		throw new Error(error.message);
	}
};

export default resetPasswordService;
