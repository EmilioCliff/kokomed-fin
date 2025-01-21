import { protectedApi } from '@/API/api';
import { refreshTokenRes, tokenData } from '@/lib/types';

export async function refreshToken() {
  try {
    const response = await protectedApi
      .post<refreshTokenRes>('/refreshToken')
      .then((res) => res.data);
    if (response.status === 'Failure') {
      throw new Error(response.error);
    }

    sessionStorage.setItem('accessToken', response.accessToken);

    scheduleTokenRefresh(2343545); // time in unix

    return response.accessToken;
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
