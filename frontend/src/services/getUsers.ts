import api from '@/API/api';
import { getUsersType } from '@/lib/types';
import { tableFilterType } from '@/lib/types';

const getUsers = async (
	pageNumber: number,
	pageSize: number,
	filter: tableFilterType,
	search: string,
) => {
	try {
		let baseUrl = `/user?limit=${pageSize}&page=${pageNumber + 1}`;

		if (search) {
			baseUrl = baseUrl + `&search=${encodeURIComponent(search)}`;
		}

		if (filter.options.length > 0) {
			const roles = filter.options
				.map(({ value }) => value.toUpperCase())
				.join(',');
			baseUrl += `&role=${encodeURIComponent(roles)}`;
		}

		const response = await api
			.get<getUsersType>(baseUrl)
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

export default getUsers;
