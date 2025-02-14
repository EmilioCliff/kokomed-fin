import { useState } from 'react';
import TableSkeleton from '@/components/UI/TableSkeleton';
import { clientStatus } from '@/data/loan';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { DataTable } from '../../table/data-table';
import CustomerForm from '@/components/PAGES/customers/CustomerForm';
import { clientColumns } from '@/components/PAGES/customers/customer';
import { useDebounce } from '@/hooks/useDebounce';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import getClients from '@/services/getClients';
import { useTable } from '@/hooks/useTable';
import CustomerSheet from './CustomerSheet';

function CustomersPage() {
	const [formOpen, setFormOpen] = useState(false);
	const { pageIndex, pageSize, filter, search, updateTableContext } =
		useTable();

	const debouncedInput = useDebounce({ value: search, delay: 500 });

	const { isLoading, error, data } = useQuery({
		queryKey: ['clients', pageIndex, pageSize, filter, debouncedInput],
		queryFn: () => getClients(pageIndex, pageSize, filter, debouncedInput),
		staleTime: 5 * 1000,
		placeholderData: keepPreviousData,
	});

	if (data?.metadata) {
		updateTableContext(data.metadata);
	}

	if (isLoading) {
		return <TableSkeleton />;
	}

	if (error) {
		return <div>Error: {error.message}</div>;
	}

	return (
		<div className="px-4">
			<div className="flex justify-between items-center mb-4">
				<h1 className="text-3xl font-bold">Customers</h1>
				<Dialog open={formOpen} onOpenChange={setFormOpen}>
					<DialogTrigger asChild>
						<Button className="text-xs py-1 font-bold" size="sm">
							Add New Customer
						</Button>
					</DialogTrigger>
					<DialogContent className="max-w-screen-lg max-h-screen overflow-y-auto">
						<DialogHeader>
							<DialogTitle>Add New Customer</DialogTitle>
							<DialogDescription>
								Enter the details for the customer user.
							</DialogDescription>
						</DialogHeader>
						<CustomerForm onFormOpen={setFormOpen} />
					</DialogContent>
				</Dialog>
			</div>
			<DataTable
				data={data?.data || []}
				columns={clientColumns}
				searchableColumns={[
					{
						id: 'fullName',
						title: 'Client Name',
					},
					{
						id: 'phoneNumber',
						title: 'Phone Number',
					},
				]}
				facetedFilterColumns={[
					{
						id: 'active',
						title: 'Active',
						options: clientStatus,
					},
				]}
			/>
			<CustomerSheet />
		</div>
	);
}

export default CustomersPage;
