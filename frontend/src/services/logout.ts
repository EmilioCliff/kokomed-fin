export const logout = () => {
  sessionStorage.removeItem('accessToken');
  document.cookie = 'refreshToken=; Max-Age=0; Path=/';
  console.log('Logged out');
};
