import { useState, useEffect } from 'react';
import { useDebounce } from '@/hooks/useDebounce';
import { useTable } from '@/hooks/useTable';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import TableSkeleton from '@/components/UI/TableSkeleton';
import getPayments from '@/services/getPayments';
import PaymentSheet from './PaymentSheet';
import PaymentForm from './PaymentForm';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { DataTable } from '@/components/table/data-table';
import { paymentColumns } from './payment';
import { paymentSources } from '@/data/loan';
import { useAuth } from '@/hooks/useAuth';
import { role } from '@/lib/types';
import { DateRangePicker } from '@/components/UI/DateRangePicker';

function PaymentsPage() {
	const [formOpen, setFormOpen] = useState(false);
	const {
		pageIndex,
		pageSize,
		filter,
		search,
		fromDate,
		toDate,
		updateTableContext,
	} = useTable();
	const { decoded } = useAuth();

	const debouncedInput = useDebounce({ value: search, delay: 500 });

	const { isLoading, error, data } = useQuery({
		queryKey: [
			'payments',
			pageIndex,
			pageSize,
			filter,
			fromDate,
			toDate,
			debouncedInput,
		],
		queryFn: () =>
			getPayments(
				pageIndex,
				pageSize,
				filter,
				fromDate,
				toDate,
				debouncedInput,
			),
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
			<div className="flex justify-end mb-4">
				<DateRangePicker />
			</div>
			<div className="flex justify-between items-center mb-4">
				<h1 className="text-3xl font-bold">Payments</h1>
				{decoded?.role === role.ADMIN && (
					<Dialog open={formOpen} onOpenChange={setFormOpen}>
						<DialogTrigger asChild>
							<Button
								className="text-xs py-1 font-bold"
								size="sm"
							>
								Add New Payment
							</Button>
						</DialogTrigger>
						<DialogContent className="max-w-screen-lg max-h-screen overflow-y-auto">
							<DialogHeader>
								<DialogTitle>Add New Payment</DialogTitle>
								<DialogDescription>
									Submiting this form creates a Clients
									Payment
								</DialogDescription>
							</DialogHeader>
							<PaymentForm onFormOpen={setFormOpen} />
						</DialogContent>
					</Dialog>
				)}
			</div>
			<DataTable
				data={data?.data || []}
				columns={paymentColumns}
				searchableColumns={[
					{
						id: 'payingName',
						title: 'Paying Name',
					},
					{
						id: 'accountNumber',
						title: 'Account Number',
					},
					{
						id: 'transactionNumber',
						title: 'Transaction Number',
					},
				]}
				facetedFilterColumns={[
					{
						id: 'transactionSource',
						title: 'Transaction Source',
						options: paymentSources,
					},
					// {
					// 	id: 'assigned',
					// 	title: 'Assigned Payments',
					// 	options: assignedStatus,
					// },
				]}
			/>
			<PaymentSheet />
		</div>
	);
}

export default PaymentsPage;
