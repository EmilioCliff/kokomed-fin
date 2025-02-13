import api from '@/API/api';
import { getClientsType, tableFilterType } from '@/lib/types';

const getClients = async (
	pageNumber: number,
	pageSize: number,
	filter: tableFilterType,
	search: string,
) => {
	try {
		let baseUrl = `/client?limit=${pageSize}&page=${pageNumber + 1}`;

		if (search) {
			baseUrl = baseUrl + `&search=${encodeURIComponent(search)}`;
		}

		if (filter.options.length === 1) {
			const active = filter.options[0];
			baseUrl += `&active=${
				active.value === 'true'
					? encodeURIComponent(1)
					: encodeURIComponent(2)
			}`;
		}

		const response = await api
			.get<getClientsType>(baseUrl)
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

export default getClients;
