import api from '@/API/api';
import { getFormDataType } from '@/lib/types';

const getFormData = async (
	product: boolean,
	client: boolean,
	user: boolean,
	branch: boolean,
) => {
	try {
		let baseUrl = `/helper/formData?`;

		if (product) {
			baseUrl = baseUrl + '&products=true';
		}

		if (client) {
			baseUrl = baseUrl + '&client=true';
		}

		if (user) {
			baseUrl = baseUrl + '&user=true';
		}

		if (branch) {
			baseUrl = baseUrl + '&branch=true';
		}

		const response = await api
			.get<getFormDataType>(baseUrl)
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

export default getFormData;
