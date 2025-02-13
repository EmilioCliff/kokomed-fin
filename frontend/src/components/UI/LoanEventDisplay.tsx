import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { LoanEvent } from '@/lib/types';

interface LoanEventDisplayProps {
	event: LoanEvent;
}

function LoanEventDisplay({ event }: LoanEventDisplayProps) {
	return (
		<Card className="p-2">
			<CardHeader className="p-0 mb-2">
				<CardTitle className="text-lg font-semibold">
					{event.id}
				</CardTitle>
			</CardHeader>
			<div className="flex gap-x-2 justify-between">
				<div>
					<CardDescription className="text-sm text-muted-foreground">
						Client Name
					</CardDescription>
					<CardContent className="p-0 text-start">
						<p className="text-sm">{event.clientName}</p>
					</CardContent>
				</div>
				<div>
					{event.type === 'due' ? (
						<div>
							<CardDescription className="text-sm text-muted-foreground">
								Payment Due
							</CardDescription>
							<CardContent className="p-0 text-start">
								<p className="text-sm">
									{event.paymentDue?.toLocaleString()}
								</p>
							</CardContent>
						</div>
					) : (
						<div>
							<CardDescription className="text-sm text-muted-foreground">
								Disburse Amount
							</CardDescription>
							<CardContent className="p-0 text-start">
								<p className="text-sm">
									{event.loanAmount.toLocaleString()}
								</p>
							</CardContent>
						</div>
					)}
				</div>
			</div>
			<CardDescription className="text-sm text-muted-foreground my-4">
				<p
					className={`text-sm font-medium ${
						event.type === 'disbursed'
							? 'text-green-600'
							: 'text-red-600'
					}`}
				>
					{event.type === 'disbursed'
						? 'Loan Disbursed'
						: 'Payment Due'}
				</p>
			</CardDescription>
		</Card>
	);
}

export default LoanEventDisplay;
