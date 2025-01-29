import api, { protectedApi } from '@/API/api';
import { tokenData } from '@/lib/types';
import { LoginForm } from '@/components/PAGES/login/schema';

export const login = async (details: LoginForm) => {
  try {
    const response = await api
      .post<tokenData>('/login', {
        email: details.email,
        password: details.password,
      })
      .then((resp) => resp.data);

    if (response.message) {
      throw new Error(response.message || 'An unknown error occurred.');
    }

    sessionStorage.setItem('accessToken', response.accessToken);

    // save the refreshToken
    // document.cookie = `refreshToken=${response.data.refreshToken}; HttpOnly; Secure; SameSite=Strict; Path=/`;

    return response;
  } catch (error: any) {
    if (error.response) {
      throw new Error('Unauthorized access. Please check your credentials.');
    }

    throw new Error(error.message);
  }
};
