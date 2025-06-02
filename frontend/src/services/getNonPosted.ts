import api from '@/API/api';
import { getPaymentType } from '@/lib/types';

const getNonPosted = async (id: number) => {
	try {
		const response = await api
			.get<getPaymentType>(`/non-posted/${id}`)
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

export default getNonPosted;
