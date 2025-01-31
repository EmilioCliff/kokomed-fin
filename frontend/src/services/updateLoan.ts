import api from '@/API/api';
import { updateLoanType } from '@/lib/types';
import { format } from 'date-fns';

export const updateLoan = async (data: updateLoanType) => {
	console.log(data);
	try {
		if (data.status && !data.disburseDate) {
			data.disburseDate = format(new Date(), 'yyyy-MM-dd');
		}

		const response = await api
			.patch<updateLoanType>(`/loan/${data.id}/disburse`, data)
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
