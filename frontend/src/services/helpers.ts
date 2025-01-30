import api from "@/API/api";
import { DashboardData } from "@/components/PAGES/dashboard/schema";
import { commonDataResponse } from "@/lib/types";

interface loanFormData {
  products?: commonDataResponse[];
  loanOfficers?: commonDataResponse[];
  clients?: commonDataResponse[];
}

export const getDashboardData = async () => {
  try {
    const response = await api
      .get<DashboardData>("/helper/dashboard")
      .then((res) => res.data);

    return response;
  } catch (error: any) {
    throw new Error(error.message);
  }
};

export const getLoanFormData = async () => {
  try {
    const response = await api
      .get<loanFormData>("/helper/loanForm")
      .then((res) => res.data);

    return response;
  } catch (error) {
    console.error(error);
  }
};
