import api, { protectedApi } from '@/API/api';
import { refreshTokenRes, tokenData } from '@/lib/types';

export async function refreshToken() {
  try {
    const response = await api.post<tokenData>('/refreshToken').then((res) => res.data);
    if (response.message) {
      throw new Error(response.message);
    }

    sessionStorage.setItem('accessToken', response.accessToken);

    // scheduleTokenRefresh(2343545); // time in unix
    return response;
  } catch (error) {
    console.error(error);
  }
}

function scheduleTokenRefresh(accessTokenExpiresAt: number) {
  const refreshTime = accessTokenExpiresAt - Date.now() - 5000; // Refresh 5 seconds before expiry
  if (refreshTime > 0) {
    setTimeout(refreshToken, refreshTime);
  }
}
