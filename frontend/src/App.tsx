import "./App.css";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import AppLayout from "./layouts/AppLayout";
import Dashboard from "./pages/Dashboard";
import LoansPage from "./pages/LoansPage";
import LoginPage from "./pages/LoginPage";
import CustomersPage from "./pages/CustomersPage";
import PaymentsPage from "./pages/PaymentsPage";
import UsersPage from "./pages/UsersPage";
import BranchesPage from "./pages/BranchesPage";
import ProductsPage from "./pages/ProductsPage";

function App() {
	return (
		<Router>
			<Routes>
				<Route path='/login' element={<LoginPage />} />
				<Route path='/' element={<AppLayout />}>
					<Route index element={<Dashboard />} />
					<Route path='loans/overview' element={<LoansPage />} />
					<Route path='customers/overview' element={<CustomersPage />} />
					<Route path='users/overview' element={<UsersPage />} />
					<Route path='branches/overview' element={<BranchesPage />} />
					<Route path='payments/overview' element={<PaymentsPage />} />
					<Route path='products/overview' element={<ProductsPage />} />
					<Route path='*' element={<h1 className='mt-10'>404</h1>} />
				</Route>
			</Routes>
		</Router>
	);
}

export default App;
