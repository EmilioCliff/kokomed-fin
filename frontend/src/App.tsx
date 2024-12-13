import "./App.css";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { SidebarProvider } from "@/components/ui/sidebar";
import AppSidebar from "./components/AppSidebar";
import Navbar from "./components/Navbar";
import Dashboard from "./pages/Dashboard";
import LoansPage from "./pages/LoansPage";
import LoginPage from "./pages/LoginPage";
import CustomersPage from "./pages/CustomersPage";
import PaymentsPage from "./pages/PaymentsPage";

function App() {
	return (
		<Router>
			<SidebarProvider>
				<AppSidebar />
				<Navbar />
				<div className='p-4 overflow-x-auto my-28 px-2'>
					<Routes>
						<Route path='/login' element={<LoginPage />} />
						<Route path='/' element={<Dashboard />} />
						<Route path='/loans/overview' element={<LoansPage />} />
						<Route path='/customers/overview' element={<CustomersPage />} />
						<Route path='/payments/overview' element={<PaymentsPage />} />
						<Route path='*' element={<h1 className='mt-10'>404</h1>} />
					</Routes>
				</div>
			</SidebarProvider>
		</Router>
	);
}

export default App;
