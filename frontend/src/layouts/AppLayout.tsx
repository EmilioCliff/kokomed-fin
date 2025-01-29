import { Outlet, Navigate } from 'react-router-dom';
import { SidebarProvider } from '@/components/ui/sidebar';
import AppSidebar from '@/components/UI/AppSidebar';
import Navbar from '@/components/UI/Navbar';
import { AuthContext } from '@/context/AuthContext';
import { useContext } from 'react';
import { useQuery } from '@tanstack/react-query';
import { refreshToken } from '@/services/refreshToken';

export default function AppLayout() {
  const { isAuthenticated, updateAuthContext } = useContext(AuthContext);

  const { isLoading, data, error } = useQuery({
    queryKey: ['refreshToken'],
    queryFn: refreshToken,
    staleTime: 0,
    enabled: !isAuthenticated,
    retry: 2,
  });

  // check if there is a refresh token stored if yes try to refresh the token then udate it
  if (data) {
    updateAuthContext(data);
  }

  return (
    <>
      {isAuthenticated ? (
        <SidebarProvider>
          <AppSidebar />
          <Navbar />
          <div className="p-4 overflow-x-auto my-28 px-2 w-full">
            <Outlet />
          </div>
        </SidebarProvider>
      ) : (
        <Navigate to="/login" />
      )}
    </>
  );
}
