import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Eye } from 'lucide-react';
import { formatCurrency, formatDate } from '@/lib/utils';

interface PaymentsTableProps {
	payments: any[];
	onPaymentClick: (payment: any) => void;
}

export function LoanDetailsPaymentsTable({
	payments,
	onPaymentClick,
}: PaymentsTableProps) {
	return (
		<div className="overflow-x-auto">
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead className="w-[80px]">ID</TableHead>
						<TableHead>Source</TableHead>
						<TableHead>Transaction #</TableHead>
						<TableHead>Paying Name</TableHead>
						<TableHead>Amount</TableHead>
						<TableHead>Date</TableHead>
						<TableHead className="text-right">Actions</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{payments.length > 0 ? (
						payments.map((payment) => (
							<TableRow
								key={payment.id}
								className="cursor-pointer hover:bg-muted/50"
							>
								<TableCell className="font-medium">
									{payment.id}
								</TableCell>
								<TableCell>
									<Badge
										variant={
											payment.transactionSource ===
											'MPESA'
												? 'default'
												: 'secondary'
										}
									>
										{payment.transactionSource}
									</Badge>
								</TableCell>
								<TableCell className="font-mono text-xs">
									{payment.transactionNumber}
								</TableCell>
								<TableCell>{payment.payingName}</TableCell>
								<TableCell className="font-medium">
									{formatCurrency(payment.amount)}
								</TableCell>
								<TableCell>
									{formatDate(payment.paidDate)}
								</TableCell>
								<TableCell className="text-right">
									<Button
										variant="ghost"
										size="sm"
										onClick={() => onPaymentClick(payment)}
										className="h-8 w-8 p-0"
									>
										<Eye className="h-4 w-4" />
										<span className="sr-only">
											View allocations
										</span>
									</Button>
								</TableCell>
							</TableRow>
						))
					) : (
						<TableRow>
							<TableCell colSpan={7} className="h-24 text-center">
								No Payments
							</TableCell>
						</TableRow>
					)}
				</TableBody>
			</Table>
		</div>
	);
}
