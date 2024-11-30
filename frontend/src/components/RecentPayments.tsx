import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";

const recentPayments = [
	{ id: 1, borrower: "John Doe", amount: 1000, date: "2023-04-15" },
	{ id: 2, borrower: "Jane Smith", amount: 750, date: "2023-04-14" },
	{ id: 3, borrower: "Bob Johnson", amount: 1200, date: "2023-04-13" },
	{ id: 4, borrower: "Alice Brown", amount: 500, date: "2023-04-12" },
];

function RecentPayments() {
	return (
		<Card className='col-span-1'>
			<CardHeader>
				<CardTitle className='text-start'>Recent Payments</CardTitle>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead className='text-center'>Borrower</TableHead>
							<TableHead className='text-center'>Amount</TableHead>
							<TableHead className='text-center'>Date</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						{recentPayments.map((payment) => (
							<TableRow key={payment.id}>
								<TableCell>{payment.borrower}</TableCell>
								<TableCell>${payment.amount}</TableCell>
								<TableCell>{payment.date}</TableCell>
							</TableRow>
						))}
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	);
}

export default RecentPayments;
