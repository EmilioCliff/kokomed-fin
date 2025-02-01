import api from '@/API/api';
import { updateClientType } from '@/lib/types';

const updateClient = async (data: updateClientType) => {
	console.log(data);

	try {
		const response = await api
			.patch<updateClientType>(`/client/${data.id}`, data)
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

export default updateClient;
