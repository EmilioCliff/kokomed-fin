import Widgets from "@/components/Widgets";
import LoanStatusChart from "@/components/LoanStatusChart";
import RecentPayments from "@/components/RecentPayments";
import { Wallet, Flag, DollarSign, Users } from "lucide-react";

const widgets = [
	{
		title: "Customers",
		icon: Users,
		mainAmount: 1000,
		active: 600,
		activeTitle: "Active",
		closed: 400,
		closedTitle: "Inactive",
	},
	{
		title: "Loans",
		icon: Wallet,
		mainAmount: 2000,
		active: 20,
		activeTitle: "Active",
		closed: 40,
		closedTitle: "Inactive",
	},
	{
		title: "Transactions",
		icon: Flag,
		mainAmount: 600,
		active: 400,
		activeTitle: "Disbursed",
		closed: 200,
		closedTitle: "Received",
		currency: "Ksh",
	},
	{
		title: "Payments",
		icon: DollarSign,
		mainAmount: 500000,
		active: 400000,
		activeTitle: "Posted",
		closed: 100000,
		closedTitle: "Non-Posted",
		currency: "Ksh",
	},
];

function Dashboard() {
	return (
		<div className='p-6 mt-10'>
			<h1 className='text-3xl font-bold mb-6 text-start'>Dashboard</h1>
			<div className='space-y-6'>
				<div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4'>
					{widgets.map((widget, index) => (
						<Widgets key={index} {...widget} />
					))}
				</div>
				<div className='grid grid-cols-1 lg:grid-cols-2 gap-4'>
					<LoanStatusChart />
					<RecentPayments />
				</div>
			</div>
		</div>
	);
}

export default Dashboard;
