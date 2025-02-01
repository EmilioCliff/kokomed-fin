import api from '@/API/api';
import { ClientFormType } from '@/components/PAGES/customers/schema';
import { commonresponse } from '@/lib/types';

const addClient = async (data: ClientFormType) => {
	try {
		const response = await api
			.post<commonresponse>('/client', data)
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

export default addClient;
