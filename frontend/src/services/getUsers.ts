import api from "@/API/api";
import { getUserType } from "@/lib/types";
export async function getUser() {
  try {
    const response = await api.get<getUserType>("/user?ID=12345").then((res) => res.data);
    if (response.status === "Failure") {
      // show error using toast
      throw new Error(response.error);
    }

    return response.data;
  } catch (error) {
    console.error(error);
  }
}

//   concurrent requests at once using Promise.all
// Promise.all([getUserAccount(), getUserPermissions()])
// .then(function ([acct, perm]) {
// // ...
// });
