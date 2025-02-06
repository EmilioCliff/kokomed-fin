import api from '@/API/api';
import { commonresponse } from '@/lib/types';
import { ProductFormType } from '@/components/PAGES/products/schema';

const addProduct = async (data: ProductFormType) => {
	try {
		const response = await api
			.post<commonresponse>('/product', data)
			.then((resp) => resp.data);

		if (response.message) {
			throw new Error(response.message);
		}

		console.log(response);

		return response;
	} catch (error: any) {
		if (error.response) {
			throw new Error(error.response.data.message);
		}

		throw new Error(error.message);
	}
};

export default addProduct;
