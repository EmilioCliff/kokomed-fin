import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from '@/components/ui/sheet';
import { useQuery } from '@tanstack/react-query';
import InstallmentDisplay from '@/components/UI/InstallmentDisplay';
import getLoanInstallments from '@/services/getLoanInstallments';
import { useTable } from '@/hooks/useTable';
import TableSkeleton from '@/components/UI/TableSkeleton';

function LoanExpectedPaymentSheet() {
	const { selectedRow, setSelectedRow } = useTable();

	const { isLoading, error, data } = useQuery({
		queryKey: ['loan/installments', selectedRow],
		queryFn: () => getLoanInstallments(selectedRow.loanId),
		staleTime: 5 * 1000,
		enabled: !!selectedRow,
	});

	if (isLoading) {
		return <TableSkeleton />;
	}

	if (error) {
		return <div>Error: {error.message}</div>;
	}

	return (
		<Sheet
			open={!!selectedRow}
			onOpenChange={(open: boolean) => {
				if (!open) {
					setSelectedRow(null);
				}
			}}
		>
			<SheetContent className="overflow-auto no-scrollbar">
				<SheetHeader>
					<SheetTitle>Loan Details</SheetTitle>
					<SheetDescription>Loan Installments</SheetDescription>
				</SheetHeader>
				{selectedRow && (
					<div className="py-4">
						{data?.length ? (
							data.map((installment) => (
								<div
									key={installment.id}
									className="mt-4 space-y-3"
								>
									<InstallmentDisplay
										installment={installment}
									/>
								</div>
							))
						) : (
							<p className="text-gray-500">
								No installments available
							</p>
						)}
					</div>
				)}
			</SheetContent>
		</Sheet>
	);
}

export default LoanExpectedPaymentSheet;
