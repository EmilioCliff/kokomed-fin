import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { RecentPayment } from './schema';
import { format } from 'date-fns';

function RecentPayments({
	recentPayments,
}: {
	recentPayments: RecentPayment[];
}) {
	return (
		<Card className="col-span-1">
			<CardHeader>
				<CardTitle className="text-start">Recent Payments</CardTitle>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead className="text-center">
								Borrower
							</TableHead>
							<TableHead className="text-center">
								Amount
							</TableHead>
							<TableHead className="text-center">Date</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						{recentPayments.map((payment: RecentPayment) => (
							<TableRow className="text-center" key={payment.id}>
								<TableCell>{payment.payingName}</TableCell>
								<TableCell>{payment.amount}</TableCell>
								<TableCell>
									{format(payment.paidDate, 'PPP')}
								</TableCell>
							</TableRow>
						))}
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	);
}

export default RecentPayments;
