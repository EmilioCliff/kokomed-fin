import { Outlet, Navigate } from 'react-router-dom';
import { SidebarProvider } from '@/components/ui/sidebar';
import AppSidebar from '@/components/UI/AppSidebar';
import Navbar from '@/components/UI/Navbar';
import { useAuth } from '@/hooks/useAuth';
import Spinner from '@/components/UI/Spinner';

export default function AppLayout() {
	const { isAuthenticated, isChecking } = useAuth();

	if (isChecking) {
		return <Spinner />;
	}

	return (
		<>
			{isAuthenticated ? (
				<SidebarProvider>
					<AppSidebar />
					<Navbar />
					<div className="p-4 overflow-x-auto my-28 px-2 w-full no-scrollbar">
						<Outlet />
					</div>
				</SidebarProvider>
			) : (
				<Navigate to="/login" />
			)}
		</>
	);
}
