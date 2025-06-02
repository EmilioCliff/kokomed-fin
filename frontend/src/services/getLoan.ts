import api from '@/API/api';
import { getLoanType } from '@/lib/types';

const getLoan = async (id: number) => {
	try {
		const response = await api
			.get<getLoanType>(`/loan/${id}`)
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

export default getLoan;
