import api from '@/API/api';
import { getLoansType, tableFilterType } from '@/lib/types';

export const getClientLoans = async (
	id: number,
	pageNumber: number,
	pageSize: number,
	filter: tableFilterType,
) => {
	try {
		let baseUrl = `/loan/client/${id}?limit=${pageSize}&page=${
			pageNumber + 1
		}`;

		if (filter.options.length > 0) {
			const statuses = filter.options
				.map(({ value }) => value.toUpperCase())
				.join(',');
			baseUrl += `&status=${encodeURIComponent(statuses)}`;
		}

		const response = await api
			.get<getLoansType>(baseUrl)
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
