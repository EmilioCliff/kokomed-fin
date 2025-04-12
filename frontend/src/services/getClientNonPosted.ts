import api from '@/API/api';
import { getClientNonPostedType, getPaymentsType } from '@/lib/types';

export const getClientNonPosted = async ({
	id,
	phoneNumber,
}: {
	id: number;
	phoneNumber: string;
}) => {
	try {
		const response = await api
			.post<getClientNonPostedType>('/non-posted/clients', {
				id: id,
				phoneNumber: phoneNumber,
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

export const getClientPayment = async (
	id: number,
	phoneNumber: string,
	pageNumber: number,
	pageSize: number,
) => {
	try {
		const response = await api
			.post<getPaymentsType>(
				`/helper/client-payments?limit=${pageSize}&page=${
					pageNumber + 1
				}`,
				{
					id: id,
					phoneNumber: phoneNumber,
				},
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
