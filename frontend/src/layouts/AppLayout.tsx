import { Outlet } from "react-router-dom";
import { SidebarProvider } from "@/components/ui/sidebar";
import AppSidebar from "@/components/UI/AppSidebar";
import Navbar from "@/components/Navbar";

export default function AppLayout() {
	return (
		<SidebarProvider>
			<AppSidebar />
			<Navbar />
			<div className='p-4 overflow-x-auto my-28 px-2 w-full'>
				<Outlet />
			</div>
		</SidebarProvider>
	);
}
