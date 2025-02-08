import Widgets from '@/components/UI/Widgets';
import LoanStatusChart from '@/components/UI/LoanStatusChart';
import RecentPayments from '@/components/PAGES/dashboard/RecentPayments';
import { DataTable } from '@/components/table/data-table';
import { inactiveLoanColumns } from './inactive-loan';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import DashboardSkeleton from '@/components/UI/DashboardSkeleton';
import { useQuery } from '@tanstack/react-query';
import { getDashboardData } from '@/services/helpers';

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
