import api, { protectedApi } from '@/API/api';
import { DashboardData } from '@/components/PAGES/dashboard/schema';
import { commonDataResponse } from '@/lib/types';

interface loanFormData {
  products?: commonDataResponse[];
  loanOfficers?: commonDataResponse[];
  clients?: commonDataResponse[];
}

export const getDashboardData = async () => {
  try {
    const response = await protectedApi
      .get<DashboardData>('/helper/dashboard')
      .then((res) => res.data);
    return response;
  } catch (error) {
    console.error(error);
  }
};

export const getLoanFormData = async () => {
  try {
    const response = await protectedApi
      .get<loanFormData>('/helper/loanForm')
      .then((res) => res.data);

    return response;
  } catch (error) {
    console.error(error);
  }
};
