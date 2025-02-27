import { useEffect, useState } from 'react';
import LoanForm from '@/components/PAGES/loans/LoanForm';
import { loanColumns } from '@/components/PAGES/loans/loan';
import { DataTable } from '@/components/table/data-table';
import TableSkeleton from '@/components/UI/TableSkeleton';

import { statuses } from '@/data/loan';

import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import LoanSheet from './LoanSheet';

import { keepPreviousData, useQuery } from '@tanstack/react-query';
import { getLoans } from '@/services/getLoans';
import { useDebounce } from '@/hooks/useDebounce';
import { useTable } from '@/hooks/useTable';

export default function LoanPage() {
	const [formOpen, setFormOpen] = useState(false);
	const { pageIndex, pageSize, filter, search, updateTableContext } =
		useTable();

	const debouncedInput = useDebounce({ value: search, delay: 500 });

	const { isLoading, error, data } = useQuery({
		queryKey: ['loans', pageIndex, pageSize, filter, debouncedInput],
		queryFn: () => getLoans(pageIndex, pageSize, filter, debouncedInput),
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
				<h1 className="text-3xl font-bold">Loans</h1>
				<Dialog open={formOpen} onOpenChange={setFormOpen}>
					<DialogTrigger asChild>
						<Button className="text-xs py-1 font-bold" size="sm">
							Add New Loan
						</Button>
					</DialogTrigger>
					<DialogContent className="max-w-screen-lg max-h-screen overflow-y-auto">
						<DialogHeader>
							<DialogTitle>Add New Loan</DialogTitle>
							<DialogDescription>
								Enter the details for the new loan.
							</DialogDescription>
						</DialogHeader>
						<LoanForm onFormOpen={setFormOpen} />
					</DialogContent>
				</Dialog>
			</div>
			<DataTable
				data={data?.data || []}
				columns={loanColumns}
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
				facetedFilterColumns={[
					{
						id: 'status',
						title: 'Status',
						options: statuses,
					},
				]}
			/>
			<LoanSheet />
		</div>
	);
}
