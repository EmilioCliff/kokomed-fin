import api from '@/API/api';
import { commonresponse } from '@/lib/types';
import { BranchFormType } from '@/components/PAGES/branches/schema';

const addBranch = async (data: BranchFormType) => {
	try {
		const response = await api
			.post<commonresponse>('/branch', data)
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

export default addBranch;
