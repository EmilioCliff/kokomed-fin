import api from '@/API/api';
import { getLoanEventsType } from '@/lib/types';

const getLoanEvents = async () => {
	try {
		const response = await api
			.get<getLoanEventsType>('/helper/loanEvents')
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

export default getLoanEvents;
