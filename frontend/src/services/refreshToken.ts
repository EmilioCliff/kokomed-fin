import api from "@/API/api";
import { tokenData } from "@/lib/types";

const refreshTokenService = async () => {
  try {
    const response = await api.get<tokenData>("/refreshToken").then((res) => res.data);

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

export default refreshTokenService;
