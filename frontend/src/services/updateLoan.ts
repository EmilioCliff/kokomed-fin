import api, { protectedApi } from '@/API/api';
import { updateLoanType } from '@/lib/types';
import { format } from 'date-fns';

export const updateLoan = async (data: updateLoanType) => {
  if (!data.disburseDate) {
    data.disburseDate = format(new Date(), 'PPP');
  }
  const response = await api
    .post<updateLoanType>('/updateLoan', data)
    .then((resp) => resp.data);

  return response;
};
