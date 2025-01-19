import axios, { AxiosError } from 'axios';
import { refreshToken } from '@/services/refreshToken';
import { useAuthContext } from '@/context/AuthContext';

const api = axios.create({
    baseURL: 'http://localhost:3000',
    headers: {
        'Content-Type': 'application/json',
    },
});

const data = useAuthContext();

export default api;

export const protectedApi = axios.create({
    baseURL: 'http://localhost:3000',
    headers: {
        'Content-Type': 'application/json',
    },
    // withCredentials: true,
});

const prefetchInt = protectedApi.interceptors.request.use(
    (config) => {
      const accessToken = localStorage.getItem('accessToken') || sessionStorage.getItem('accessToken');
      if (accessToken) {
        config.headers.Authorization = `Bearer ${accessToken}`;
      }
  
      return config;
    },
  );

const postfetchInt = protectedApi.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
        const originalRequest = error.config;
            // prevent infinite loop
            // @ts-ignore
        if (originalRequest && error.response?.status === 401 && !originalRequest?._retry) {
            // @ts-ignore
            originalRequest._retry = true;
            try {
                await refreshToken();
                return protectedApi.request(originalRequest);
            } catch (error) {
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
}
