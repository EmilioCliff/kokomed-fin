import { protectedApi } from '@/API/api';
import api from '@/API/api';
import { getLoansType } from '@/lib/types';
import { Loan } from '@/data/schema';

interface getLoansJson {
  first: number;
  prev: number;
  next: number;
  last: number;
  pages: number;
  items: number;
  data: Loan[];
}

export const getLoans = async (pageNumber: number, pageSize: number) => {
  try {
    const response = await api
      .get<getLoansJson>(`/loans?_per_page=${pageSize}&_page=${pageNumber + 1}`)
      .then((res) => res.data);

    // console.log(response);

    // if i get a metadata use the z to parse it then return parsed data

    // if (response.status === 'Failure') {
    //   // show error using toast
    //   throw new Error(response.error);
    // }

    return response;
  } catch (error) {
    console.error(error);
  }
};
