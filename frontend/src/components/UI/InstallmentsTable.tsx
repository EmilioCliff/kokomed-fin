import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { CheckCircle, XCircle } from 'lucide-react';
import { formatCurrency, formatDate } from '@/lib/utils';
import { Card } from '../ui/card';

interface InstallmentsTableProps {
	installments: any[];
}

export function LoanDetailsInstallmentsTable({
	installments,
}: InstallmentsTableProps) {
	return (
		<div className="overflow-x-auto">
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead className="w-[80px]">No.</TableHead>
						<TableHead>Amount</TableHead>
						<TableHead>Remaining</TableHead>
						<TableHead>Due Date</TableHead>
						<TableHead>Status</TableHead>
						<TableHead>Paid Date</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{installments.length > 0 ? (
						installments.map((installment) => {
							const isPaid = installment.paid;
							const isPastDue =
								!isPaid &&
								new Date(installment.dueDate) < new Date();
							const isPartiallyPaid =
								!isPaid &&
								installment.remainingAmount <
									installment.amount;

							return (
								<TableRow key={installment.id}>
									<TableCell className="font-medium">
										{installment.installmentNo}
									</TableCell>
									<TableCell>
										{formatCurrency(installment.amount)}
									</TableCell>
									<TableCell>
										{formatCurrency(
											installment.remainingAmount,
										)}
									</TableCell>
									<TableCell>
										{formatDate(installment.dueDate)}
									</TableCell>
									<TableCell>
										{isPaid ? (
											<Badge
												// variant="success"
												className="flex items-center gap-1 w-fit"
											>
												<CheckCircle className="h-3.5 w-3.5" />
												Paid
											</Badge>
										) : isPartiallyPaid ? (
											<Badge
												// variant="warning"
												className="flex items-center gap-1 w-fit"
											>
												Partial
											</Badge>
										) : isPastDue ? (
											<Badge
												variant="destructive"
												className="flex items-center gap-1 w-fit"
											>
												<XCircle className="h-3.5 w-3.5" />
												Overdue
											</Badge>
										) : (
											<Badge
												variant="outline"
												className="flex items-center gap-1 w-fit"
											>
												Pending
											</Badge>
										)}
									</TableCell>
									<TableCell>
										{installment.paidAt &&
										installment.paidAt !== '0001-01-01'
											? formatDate(installment.paidAt)
											: '-'}
									</TableCell>
								</TableRow>
							);
						})
					) : (
						<TableRow>
							<TableCell colSpan={6} className="h-24 text-center">
								No Installments
							</TableCell>
						</TableRow>
					)}
				</TableBody>
			</Table>
		</div>
	);
}
