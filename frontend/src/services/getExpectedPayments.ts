import api from '@/API/api';
import { getExpectedPaymentsType } from '@/lib/types';

const getExpectedPayments = async (
	pageNumber: number,
	pageSize: number,
	search: string,
) => {
	try {
		let baseUrl = `/loan/expected-payments?limit=${pageSize}&page=${
			pageNumber + 1
		}`;

		if (search) {
			baseUrl = baseUrl + `&search=${encodeURIComponent(search)}`;
		}

		const response = await api
			.get<getExpectedPaymentsType>(baseUrl)
			.then((res) => res.data);

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

export default getExpectedPayments;
