import api from '@/API/api';
import { PaymentFormType } from '@/components/PAGES/payments/schema';
import { commonresponse } from '@/lib/types';
import { format } from 'date-fns';

const addPayment = async (data: PaymentFormType) => {
	try {
		if (data.DatePaid === '') {
			const date = new Date();
			data.DatePaid = format(date, 'yyyy-MM-dd');
		}

		const response = await api
			.post<commonresponse>('/payment/callback', {
				TransAmount: String(data.TransAmount),
				TransID: data.TransID,
				BillRefNumber: data.BillRefNumber,
				MSISDN: data.MSISDN,
				FirstName: data.FirstName,
				App: data.App,
				DatePaid: data.DatePaid,
				Email: data.Email,
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
