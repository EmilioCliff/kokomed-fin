import axios from 'axios';

const api = axios.create({
	baseURL: import.meta.env.VITE_API_URL,
	withCredentials: true,
	headers: {
		'Content-Type': 'application/json',
	},
});

export default api;

// Function to attach access token
export const setupAxiosInterceptors = (
	getAccessToken: () => string | null,
	refreshSession: () => Promise<void>,
) => {
	api.interceptors.request.use(
		(config) => {
			const token = getAccessToken();
			if (token) {
				config.headers.Authorization = `Bearer ${token}`;
			}
			return config;
		},
		(error) => Promise.reject(error),
	);

	api.interceptors.response.use(
		(response) => response,
		async (error) => {
			const originalRequest = error.config;

			if (originalRequest.url === '/login') {
				return Promise.reject(error);
			}

			if (error.response?.status === 401 && !originalRequest._retry) {
				originalRequest._retry = true;

				try {
					await refreshSession();
					return api(originalRequest);
				} catch (refreshError) {
					// window.location.href = "/login";
					return Promise.reject(refreshError);
				}
			}

			return Promise.reject(error);
		},
	);
};
