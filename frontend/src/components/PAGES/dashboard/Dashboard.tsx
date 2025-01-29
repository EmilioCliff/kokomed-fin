import Widgets from '@/components/UI/Widgets';
import LoanStatusChart from '@/components/UI/LoanStatusChart';
import RecentPayments from '@/components/UI/RecentPayments';
import { generateRandomInactiveLoans } from '@/lib/generator';
import { DataTable } from '@/components/table/data-table';
import { inactiveLoanColumns } from './inactive-loan';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import DashboardSkeleton from '@/components/UI/DashboardSkeleton';
import { DashboardData } from './schema';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import { getDashboardData } from '@/services/helpers';

// const widgets = [
//   {
//     title: 'Customers',
//     icon: Users,
//     mainAmount: 1000,
//     active: 600,
//     activeTitle: 'Active',
//     closed: 400,
//     closedTitle: 'Inactive',
//   },
//   {
//     title: 'Loans',
//     icon: Wallet,
//     mainAmount: 2000,
//     active: 20,
//     activeTitle: 'Active',
//     closed: 40,
//     closedTitle: 'Inactive',
//   },

// totalLoanAmount, _ := data.TotalLoanAmount.(float64)
// totalLoanDisbursed, _ := data.TotalLoanDisbursed.(float64)
// totalLoanPaid, _ := data.TotalLoanPaid.(float64)
// totalPaymentsReceived, _ := data.TotalPaymentsReceived.(float64)
// totalNonPosted, _ := data.TotalNonPosted.(float64)

//   {
//     title: 'Transactions',
//     icon: Flag,
//     mainAmount: 600,
//     active: 400,
//     activeTitle: 'Disbursed',
//     closed: 200,
//     closedTitle: 'Received',
//     currency: 'Ksh',
//   },
//   {
//     title: 'Payments',
//     icon: DollarSign,
//     mainAmount: 500000,
//     active: 400000,
//     activeTitle: 'Posted',
//     closed: 100000,
//     closedTitle: 'Non-Posted',
//     currency: 'Ksh',
//   },
// ];

// const recentPayments = [
//   { id: 1, borrower: 'John Doe', amount: 1000, date: '2023-04-15' },
//   { id: 2, borrower: 'Jane Smith', amount: 750, date: '2023-04-14' },
//   { id: 3, borrower: 'Bob Johnson', amount: 1200, date: '2023-04-13' },
//   { id: 4, borrower: 'Alice Brown', amount: 500, date: '2023-04-12' },
// ];

function Dashboard() {
	const { isLoading, error, data } = useQuery({
		queryKey: ['dashboard'],
		queryFn: getDashboardData,
		staleTime: 5 * 1000,
	});

	console.log(data);

	if (isLoading) {
		return <DashboardSkeleton />;
	}

	if (error) {
		return <div>Error: {error.message}</div>;
	}

	return (
		<div className="px-4">
			<h1 className="text-3xl font-bold mb-6 text-start">Dashboard</h1>
			<div className="space-y-6">
				<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
					{data?.widget_data.map((widget, index) => (
						<Widgets key={index} {...widget} />
					))}
				</div>
				<div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
					<LoanStatusChart />
					{data?.recent_payments && (
						<RecentPayments recentPayments={data.recent_payments} />
					)}
				</div>
				<Card className="col-span-1">
					<CardHeader>
						<CardTitle className="text-start">
							Recent Undisbursed Loans
						</CardTitle>
					</CardHeader>
					<CardContent>
						{data?.inactive_loans && (
							<DataTable
								data={data.inactive_loans}
								columns={inactiveLoanColumns}
							/>
						)}
					</CardContent>
				</Card>
			</div>
		</div>
	);
}

export default Dashboard;
