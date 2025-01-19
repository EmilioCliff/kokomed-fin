import { protectedApi } from "@/API/api";
import { refreshTokenRes } from "@/lib/types";

export async function refreshToken() {
    try {
        const response = await protectedApi.post<refreshTokenRes>('/refreshToken').then((res) => res.data);
        console.log(response);
        if (response.status === 'Failure') {
            throw new Error(response.error);
        }

        // save to local storage

        return response.accessToken
    } catch (error) {
        console.error(error);
    }
}