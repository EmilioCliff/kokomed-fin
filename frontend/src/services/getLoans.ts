import api from '@/API/api';
import { tableFilterType } from '@/lib/types';
import { getLoansType } from '@/lib/types';

export const getLoans = async (
	pageNumber: number,
	pageSize: number,
	filter: tableFilterType,
	search: string,
) => {
	try {
		let baseUrl = `/loan?limit=${pageSize}&page=${pageNumber + 1}`;

		if (search) {
			baseUrl = baseUrl + `&search=${encodeURIComponent(search)}`;
		}

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
