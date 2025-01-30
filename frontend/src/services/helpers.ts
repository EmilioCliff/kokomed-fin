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

export const createUser = async () => {
  const data = await api
    .patch(
      "/user/reset-password/eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjM1YTFmMDI4LWRlZDEtMTFlZi05YWYwLTAwMTU1ZDU0NTc3MiIsInVzZXJfaWQiOjEsImVtYWlsIjoiYWxpY2VAZXhhbXBsZS5jb20iLCJyb2xlIjoiQURNSU4iLCJicmFuY2hfaWQiOjEsImlzcyI6Imtva29tZWRMb2FuQXBwIiwiZXhwIjoxNzM4MjUzNTYwLCJpYXQiOjE3MzgyMTc1NjB9.dWBIkS2ChhPYRxIH4LKU0axoQgnS8UYQnF3i4rm0jGxEEKHm-UlvkvMqbZw74VPF2dDsCrf1Ob468GdlckkCmabxghtEAHGNJVERHCeqRJvpl_euvxB1i946Q2p17YKIJQ_Wv0F2W-2x0Ph7T4f9UCeZhPqgViynp0wQhLU5uqWOtyD-A4YDR7nRaY2TSfi0Ve1gnee1KdpKLIYpSxLXsZK8VP3N3V12oyB_7VliAKj75zDwuOdvwIflVoqbC9ICd_V8I21-GKX1sop-zRJFbmPnSj88cE52rC1Q5Br0moqNXE3SWCFu60ywGp2SWI8DfB7vZQG8NPCw8dXDehbE6g",
      {
        email: "alice@example.com",
        new_password: "secret",
      }
    )
    .then((resp) => resp.data);

  console.log(data);
};
