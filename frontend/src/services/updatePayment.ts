import api from '@/API/api';
import { updatePaymentType } from '@/lib/types';

const updatePayment = async (data: updatePaymentType) => {
	console.log(data);
	try {
		const response = await api
			.patch<updatePaymentType>(`/payment/${data.id}/assign`, data)
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

export default updatePayment;
