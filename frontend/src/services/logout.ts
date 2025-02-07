import api from '@/API/api';

const logoutService = async () => {
	await api.get('/logout');
};

export default logoutService;
