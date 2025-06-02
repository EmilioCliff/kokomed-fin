import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { formatCurrency, formatDate } from '@/lib/utils';

interface PaymentAllocationModalProps {
	payment: any;
	allocations: any[];
	open: boolean;
	onClose: () => void;
}

export function LoanDetailsPaymentAllocationModal({
	payment,
	allocations,
	open,
	onClose,
}: PaymentAllocationModalProps) {
	// Calculate total allocated amount
	const totalAllocated = allocations.reduce(
		(sum, allocation) => sum + allocation.amount,
		0,
	);

	return (
		<Dialog open={open} onOpenChange={onClose}>
			<DialogContent className="max-w-3xl max-h-screen overflow-y-auto no-scrollbar">
				<DialogHeader>
					<DialogTitle>Payment Allocation Details</DialogTitle>
					<DialogDescription>
						How payment #{payment.id} was allocated to installments
					</DialogDescription>
				</DialogHeader>

				<div className="space-y-6">
					{/* Payment Summary */}
					<Card>
						<CardHeader className="pb-2">
							<CardTitle className="text-base">
								Payment Summary
							</CardTitle>
						</CardHeader>
						<CardContent>
							<div className="grid grid-cols-2 md:grid-cols-3 gap-4">
								<div>
									<p className="text-sm text-muted-foreground">
										Transaction #
									</p>
									<p className="font-mono text-sm">
										{payment.transactionNumber}
									</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground">
										Source
									</p>
									<p>{payment.transactionSource}</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground">
										Paying Name
									</p>
									<p>{payment.payingName}</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground">
										Amount
									</p>
									<p className="font-semibold">
										{formatCurrency(payment.amount)}
									</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground">
										Date
									</p>
									<p>{formatDate(payment.paidDate)}</p>
								</div>
								<div>
									<p className="text-sm text-muted-foreground">
										Account
									</p>
									<p>{payment.accountNumber}</p>
								</div>
							</div>
						</CardContent>
					</Card>

					{/* Allocation Details */}
					<div>
						<div className="flex justify-between items-center mb-2">
							<h3 className="font-semibold">
								Allocation Breakdown
							</h3>
							<div className="text-sm">
								<span className="text-muted-foreground">
									Total Allocated:
								</span>{' '}
								<span className="font-medium">
									{formatCurrency(totalAllocated)}
								</span>
								{payment.amount !== totalAllocated && (
									<span className="text-amber-600 ml-2">
										(
										{formatCurrency(
											payment.amount - totalAllocated,
										)}{' '}
										unallocated)
									</span>
								)}
							</div>
						</div>

						{allocations.length > 0 ? (
							<div className="border rounded-md overflow-hidden no-scrollbar">
								<Table>
									<TableHeader>
										<TableRow>
											<TableHead>Installment</TableHead>
											<TableHead>Amount</TableHead>
											<TableHead>Date</TableHead>
											<TableHead>Description</TableHead>
										</TableRow>
									</TableHeader>
									<TableBody>
										{allocations.map((allocation) => (
											<TableRow key={allocation.id}>
												<TableCell>
													#{allocation.installmentId}
												</TableCell>
												<TableCell className="font-medium">
													{formatCurrency(
														allocation.amount,
													)}
												</TableCell>
												<TableCell>
													{formatDate(
														allocation.createdAt,
													)}
												</TableCell>
												<TableCell
													className="max-w-[250px] truncate"
													title={
														allocation.description
													}
												>
													{allocation.description}
												</TableCell>
											</TableRow>
										))}
									</TableBody>
								</Table>
							</div>
						) : (
							<div className="text-center py-6 border rounded-md bg-muted/50">
								<p className="text-muted-foreground">
									No allocations found for this payment
								</p>
							</div>
						)}
					</div>
				</div>
			</DialogContent>
		</Dialog>
	);
}
