import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";

interface RecentPaymentsProps {
	id: number;
	borrower: string;
	amount: number;
	date: string;
}

function RecentPayments({
	recentPayments,
}: {
	recentPayments: RecentPaymentsProps[];
}) {
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
						{recentPayments.map((payment: RecentPaymentsProps) => (
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
