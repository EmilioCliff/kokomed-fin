import api from '@/API/api';
import { LoanFormType } from '@/components/PAGES/loans/schema';
import { commonresponse } from '@/lib/types';

const addLoan = async (data: LoanFormType) => {
	try {
		const response = await api
			.post<commonresponse>('/loan', data)
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

export default addLoan;
