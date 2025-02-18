import { Installment } from '../PAGES/loans/schema';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Badge } from '../ui/badge';

interface InstallmentDisplayProps {
	installment: Installment;
}

function InstallmentDisplay({ installment }: InstallmentDisplayProps) {
	return (
		<Card className="p-2">
			<CardHeader className="p-0 mb-2">
				<CardTitle className="text-lg font-semibold">
					<div className="flex">
						<p>Installment No: {installment.installmentNo}</p>
						<Badge
							className="ml-auto"
							variant={
								installment.paid ? 'default' : 'destructive'
							}
						>
							{installment.paid ? 'PAID' : 'UNPAID'}
						</Badge>
					</div>
				</CardTitle>
			</CardHeader>
			<div className="flex gap-x-2 justify-between">
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						Installment Amount
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">
							{installment.amount.toLocaleString()}
						</p>
					</CardContent>
				</div>
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						Remaining Amount
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">
							{installment.remainingAmount.toLocaleString()}
						</p>
					</CardContent>
				</div>
			</div>
			<div className="mt-6">
				{installment.paid ? (
					<p>
						Paid Date:{' '}
						<span className="text-green-500">
							{installment.paidAt}
						</span>
					</p>
				) : (
					<p>
						Due Date:{' '}
						<span className="text-red-500">
							{installment.dueDate}
						</span>
					</p>
				)}
			</div>
		</Card>
	);
}

export default InstallmentDisplay;
