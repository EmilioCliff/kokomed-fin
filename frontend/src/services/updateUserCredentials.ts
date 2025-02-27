import api from '@/API/api';

const updateUserCredentials = async (password: string, token: string) => {
	try {
		if (token === '') {
			throw new Error('Missing Access Token');
		}
		const response = await api
			.patch(`/user/reset-password/${token}`, {
				newPassword: password,
			})
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

export default updateUserCredentials;
