import axios, { AxiosError } from 'axios';
import { refreshToken } from '@/services/refreshToken';
import { logout } from '@/services/logout';
// import { AuthContext } from '@/context/AuthContext';

const api = axios.create({
  baseURL: 'http://localhost:3001',
  headers: {
    'Content-Type': 'application/json',
  },
});

// const data = useAuthContext();

export default api;

export const protectedApi = axios.create({
  baseURL: 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json',
  },
  // withCredentials: true,
});

const prefetchInt = protectedApi.interceptors.request.use((config) => {
  const token = sessionStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

const postfetchInt = protectedApi.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest?._retry) {
      originalRequest._retry = true;

      try {
        await refreshToken();
        return protectedApi.request(originalRequest);
      } catch (error) {
        logout();
        // redirect to login page
        return Promise.reject(error);
      }
    }
    return Promise.reject(error);
  }
);

// during developmnet ONLY
// remove interceptors
export const removeInterceptors = () => {
  protectedApi.interceptors.request.eject(prefetchInt);
  protectedApi.interceptors.response.eject(postfetchInt);
};
