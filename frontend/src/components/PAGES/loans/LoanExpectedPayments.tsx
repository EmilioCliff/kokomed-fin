import { useTable } from '@/hooks/useTable';
import { useDebounce } from '@/hooks/useDebounce';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import TableSkeleton from '@/components/UI/TableSkeleton';
import { DataTable } from '@/components/table/data-table';
import LoanExpectedPaymentSheet from './LoanExpectedPaymentSheet';
import getExpectedPayments from '@/services/getExpectedPayments';
import { expectedPaymentColumns } from './loan';
import { useEffect } from 'react';

function LoanExpectedPayments() {
	const { pageIndex, pageSize, search, updateTableContext } = useTable();

	const debouncedInput = useDebounce({ value: search, delay: 500 });

	const { isLoading, error, data } = useQuery({
		queryKey: [
			'loans/expected-payments',
			pageIndex,
			pageSize,
			debouncedInput,
		],
		queryFn: () => getExpectedPayments(pageIndex, pageSize, debouncedInput),
		staleTime: 5 * 1000,
		placeholderData: keepPreviousData,
	});

	useEffect(() => {
		if (data?.metadata) {
			updateTableContext(data.metadata);
		}
	}, [data]);

	if (isLoading) {
		return <TableSkeleton />;
	}

	if (error) {
		return <div>Error: {error.message}</div>;
	}

	return (
		<div className="px-4">
			<div className="flex justify-between items-center mb-4">
				<h1 className="text-3xl font-bold">Expected Payments</h1>
			</div>
			<DataTable
				data={data?.data || []}
				columns={expectedPaymentColumns}
				searchableColumns={[
					{
						id: 'clientName',
						title: 'Client Name',
					},
					{
						id: 'loanOfficerName',
						title: 'Loan Officer',
					},
				]}
			/>
			<LoanExpectedPaymentSheet />
		</div>
	);
}

export default LoanExpectedPayments;
