import api from '@/API/api';
import { PaymentFormType } from '@/components/PAGES/payments/schema';
import { commonresponse } from '@/lib/types';

const addPayment = async (data: PaymentFormType) => {
	try {
		const response = await api
			.post<commonresponse>('/payment/callback', {
				TransAmount: String(data.TransAmount),
				TransID: data.TransID,
				BillRefNumber: data.BillRefNumber,
				MSISDN: data.MSISDN,
				FirstName: data.FirstName,
				App: data.App,
			})
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

export default addPayment;
