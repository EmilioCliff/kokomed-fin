import api from '@/API/api';
import { getLoanInstallmentsType } from '@/lib/types';

const getLoanInstallments = async (loanId: number) => {
	try {
		const response = await api
			.get<getLoanInstallmentsType>(`/loan/${loanId}/installments`)
			.then((res) => res.data);

		if (response.message) {
			throw new Error(response.message);
		}

		return response.data;
	} catch (error: any) {
		if (error.response) {
			throw new Error(error.response.data.message);
		}

		throw new Error(error.message);
	}
};

export default getLoanInstallments;
