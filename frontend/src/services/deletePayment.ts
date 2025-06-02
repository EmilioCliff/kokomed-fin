import api from '@/API/api';
import { commonresponse, getSimulatedPaymentType } from '@/lib/types';

export const deletePayment = async (id: number, description: string) => {
	try {
		const response = await api
			.post<commonresponse>(`/payment/${id}/delete`, {
				description: description,
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

export const simulateDeletePayment = async (id: number) => {
	try {
		const response = await api
			.get<getSimulatedPaymentType>(`/payment/${id}/simulate-delete`)
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
