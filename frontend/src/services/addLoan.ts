import api from "@/API/api";
import { Loan } from "@/components/PAGES/loans/schema";
import { LoanFormType } from "@/components/PAGES/loans/schema";

// should return loan
export const addLoan = async (data: LoanFormType) => {
  const response = await api
    .post<LoanFormType>("/loanForm", data)
    .then((resp) => resp.data);
  console.log(response);
  return response;
};
