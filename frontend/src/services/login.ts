import api, { protectedApi } from '@/API/api';
import { tokenData } from '@/lib/types';

export const login = async ({ name, password }: { name: string; password: string }) => {
  const response = await api
    .post<tokenData>('/login', {
      name: name,
      password: password,
    })
    .then((resp) => resp.data);

  // save the accessToken
  sessionStorage.setItem('accessToken', response.accessToken);

  // save the refreshToken
  document.cookie = `refreshToken=${response.refreshToken}; HttpOnly; Secure; SameSite=Strict; Path=/`;

  console.log('successful login');
  return;
};
