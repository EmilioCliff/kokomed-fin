import "./App.css";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import AppLayout from "./layouts/AppLayout";
import Dashboard from "./components/PAGES/dashboard/Dashboard";
import LoansPage from "./components/PAGES/loans/LoansPage";
import LoginPage from "./components/PAGES/login/LoginPage";
import CustomersPage from "./components/PAGES/customers/CustomersPage";
import PaymentsPage from "./components/PAGES/payments/PaymentsPage";
import UsersPage from "./components/PAGES/users/UsersPage";
import BranchesPage from "./components/PAGES/branches/BranchesPage";
import ProductsPage from "./components/PAGES/products/ProductsPage";
// import GetDataTest from "./pages/GetDataTest";

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
					{/* <Route path='getdata' element={<GetDataTest />} /> */}
					<Route path='*' element={<h1 className='mt-10'>404</h1>} />
				</Route>
			</Routes>
		</Router>
	);
}

export default App;
