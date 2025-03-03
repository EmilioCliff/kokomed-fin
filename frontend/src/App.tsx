import './App.css';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import AppLayout from './layouts/AppLayout';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import Dashboard from './components/PAGES/dashboard/Dashboard';
import LoansPage from './components/PAGES/loans/LoansPage';
import LoginPage from './components/PAGES/login/LoginPage';
import ResetPassword from './components/PAGES/login/ResetPassword';
import CustomersPage from './components/PAGES/customers/CustomersPage';
import PaymentsPage from './components/PAGES/payments/PaymentsPage';
import UsersPage from './components/PAGES/users/UsersPage';
import BranchesPage from './components/PAGES/branches/BranchesPage';
import ProductsPage from './components/PAGES/products/ProductsPage';
// import LoanExpectedPayments from './components/PAGES/loans/LoanExpectedPayments';
import { AuthContextWrapper } from './context/AuthContext';
import { TableContextWrapper } from './context/TableContext';
import { ReportsPage } from './components/PAGES/reports/ReportsPage';
import LoanTimeline from './components/PAGES/loans/LoanTimeline';
import { ToastContainer } from 'react-toastify';
import PaymentClient from './components/PAGES/payments/PaymentClient';
// import GetDataTest from './pages/GetDataTest';
import LoanUnpaidInstallments from './components/PAGES/loans/LoanUnpaidInstallments';

const queryClient = new QueryClient();

function App() {
	return (
		<AuthContextWrapper>
			<TableContextWrapper>
				<QueryClientProvider client={queryClient}>
					<Router>
						<Routes>
							<Route path="/login" element={<LoginPage />} />
							<Route
								path="/reset-password/:token"
								element={<ResetPassword />}
							/>
							<Route path="/" element={<AppLayout />}>
								<Route index element={<Dashboard />} />
								<Route
									path="loans/overview"
									element={<LoansPage />}
								/>
								<Route
									path="loans/unpaid-installments"
									element={<LoanUnpaidInstallments />}
								/>
								<Route
									path="loans/timeline"
									element={<LoanTimeline />}
								/>
								{/* <Route
									path="loans/expected-payments"
									element={<LoanExpectedPayments />}
								/> */}
								<Route
									path="customers/overview"
									element={<CustomersPage />}
								/>
								<Route
									path="users/overview"
									element={<UsersPage />}
								/>
								<Route
									path="branches/overview"
									element={<BranchesPage />}
								/>
								<Route
									path="payments/overview"
									element={<PaymentsPage />}
								/>
								<Route
									path="payments/client"
									element={<PaymentClient />}
								/>
								<Route
									path="products/overview"
									element={<ProductsPage />}
								/>
								<Route
									path="reports"
									element={<ReportsPage />}
								/>
								{/* <Route
									path="getdata"
									element={<GetDataTest />}
								/> */}
								<Route
									path="*"
									element={<h1 className="mt-10">404</h1>}
								/>
							</Route>
						</Routes>
					</Router>
					<ToastContainer />
					<ReactQueryDevtools />
				</QueryClientProvider>
			</TableContextWrapper>
		</AuthContextWrapper>
	);
}

export default App;
