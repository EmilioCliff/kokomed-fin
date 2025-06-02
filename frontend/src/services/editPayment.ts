import api from '@/API/api';
import { EditPaymentFormType } from '@/components/PAGES/payments/schema';
import { commonresponse, getSimulatedPaymentType } from '@/lib/types';

export const editPayment = async (data: EditPaymentFormType) => {
	try {
		const response = await api
			.post<commonresponse>(`/payment/${data.id}/update`, data)
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

export const simulateEditPayment = async (data: EditPaymentFormType) => {
	try {
		const response = await api
			.post<getSimulatedPaymentType>(
				`/payment/${data.id}/simulate-update`,
				data,
			)
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
