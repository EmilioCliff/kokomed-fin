import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { CalendarIcon, CreditCard } from 'lucide-react';
import { formatCurrency, formatDate } from '@/lib/utils';

interface LoanSummaryProps {
	loan: any;
}

export function LoanDetailsCard({ loan }: LoanSummaryProps) {
	// Calculate loan progress
	const progress = Math.min(100, (loan.paidAmount / loan.repayAmount) * 100);
	const remainingAmount = loan.repayAmount - loan.paidAmount;

	return (
		<Card>
			<CardHeader>
				<CardTitle className="flex items-center gap-2">
					<CreditCard className="h-5 w-5" />
					Loan Summary
				</CardTitle>
				<CardDescription>
					Overview of loan amount, repayment, and progress
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<div className="grid grid-cols-2 gap-4">
					<div>
						<p className="text-sm text-muted-foreground">
							Loan Amount
						</p>
						<p className="text-lg font-semibold">
							{formatCurrency(loan.loanAmount)}
						</p>
					</div>
					<div>
						<p className="text-sm text-muted-foreground">
							Repay Amount
						</p>
						<p className="text-lg font-semibold">
							{formatCurrency(loan.repayAmount)}
						</p>
					</div>
					<div>
						<p className="text-sm text-muted-foreground">
							Disbursed On
						</p>
						<p className="text-lg font-semibold flex items-center gap-1">
							<CalendarIcon className="h-4 w-4" />
							{formatDate(loan.disbursedOn)}
						</p>
					</div>
					<div>
						<p className="text-sm text-muted-foreground">
							Due Date
						</p>
						<p className="text-lg font-semibold flex items-center gap-1">
							<CalendarIcon className="h-4 w-4" />
							{formatDate(loan.dueDate)}
						</p>
					</div>
				</div>

				<div className="space-y-2">
					<div className="flex justify-between">
						<span className="text-sm font-medium">
							Repayment Progress
						</span>
						<span className="text-sm font-medium">
							{progress.toFixed(0)}%
						</span>
					</div>
					<Progress value={progress} className="h-2" />
					<div className="flex justify-between text-sm text-muted-foreground">
						<span>Paid: {formatCurrency(loan.paidAmount)}</span>
						<span>
							Remaining: {formatCurrency(remainingAmount)}
						</span>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
