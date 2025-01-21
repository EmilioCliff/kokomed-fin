import { protectedApi } from '@/API/api';
import api from '@/API/api';
import { getLoansType } from '@/lib/types';
import { Loan } from '@/components/PAGES/loans/schema';
import { tableFilterType } from '@/lib/types';

export const getLoans = async (
  pageNumber: number,
  pageSize: number,
  filter: tableFilterType,
  search: string
) => {
  try {
    // if filter has options then add them to the url
    if (filter.options.length > 0 || search) {
      let baseUrl = `/loans?`;

      if (search) {
        baseUrl = baseUrl + `client.fullName_like=${encodeURIComponent(search)}`;
      }
      // &loanOfficer.fullName_like=${encodeURIComponent(search)}

      if (filter.options.length > 0) {
        filter.options.forEach(({ value }) => {
          baseUrl += `&_limit=${pageSize}&_page=${
            pageNumber + 1
          }&status=${encodeURIComponent(value.toUpperCase())}`;
        });
      }

      console.log(baseUrl);

      const response = await api.get<Loan[]>(baseUrl).then((res) => res.data);
      return response;
    }

    const response = await api
      .get<Loan[]>(`/loans?_limit=${pageSize}&_page=${pageNumber + 1}`)
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
