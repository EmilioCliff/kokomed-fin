import api from '@/API/api';
import { tableFilterType } from '@/lib/types';
import { getPaymentsType } from '@/lib/types';

const getPayments = async (
	pageNumber: number,
	pageSize: number,
	filter: tableFilterType,
	from: string,
	to: string,
	search: string,
) => {
	try {
		let baseUrl = `/non-posted/all?limit=${pageSize}&page=${
			pageNumber + 1
		}`;
		baseUrl = baseUrl + `&from=${encodeURIComponent(from)}`;
		baseUrl = baseUrl + `&to=${encodeURIComponent(to)}`;

		if (search) {
			baseUrl = baseUrl + `&search=${encodeURIComponent(search)}`;
		}

		if (filter.options.length > 0) {
			const statuses = filter.options
				.map(({ value }) => value.toUpperCase())
				.join(',');
			baseUrl += `&source=${encodeURIComponent(statuses)}`;
		}

		const response = await api
			.get<getPaymentsType>(baseUrl)
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

export default getPayments;
