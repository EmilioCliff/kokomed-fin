import { useState, useContext, useEffect } from 'react';
import TableSkeleton from '@/components/UI/TableSkeleton';
import { TableContext } from '@/context/TableContext';
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
import BranchForm from '@/components/PAGES/branches/BranchForm';
import { branchColumns } from '@/components/PAGES/branches/branch';
import { useDebounce } from '@/hooks/useDebounce';
import getBranches from '@/services/getBranches';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import BranchSheet from './BranchSheet';

function BranchesPage() {
	const [formOpen, setFormOpen] = useState(false);
	const { pageIndex, pageSize, filter, search, updateTableContext } =
		useContext(TableContext);

	const debouncedInput = useDebounce({ value: search, delay: 500 });

	const { isLoading, error, data } = useQuery({
		queryKey: ['branches', pageIndex, pageSize, filter, debouncedInput],
		queryFn: () => getBranches(pageIndex, pageSize, debouncedInput),
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
				<h1 className="text-3xl font-bold">Branches</h1>
				<Dialog open={formOpen} onOpenChange={setFormOpen}>
					<DialogTrigger asChild>
						<Button className="text-xs py-1 font-bold" size="sm">
							Add New Branch
						</Button>
					</DialogTrigger>
					<DialogContent className="max-w-screen-lg max-h-screen overflow-y-auto">
						<DialogHeader>
							<DialogTitle>Add New Branch</DialogTitle>
							<DialogDescription>
								Enter the details for the new branch.
							</DialogDescription>
						</DialogHeader>
						<BranchForm onFormOpen={setFormOpen} />
					</DialogContent>
				</Dialog>
			</div>
			<DataTable
				data={data?.data || []}
				columns={branchColumns}
				searchableColumns={[
					{
						id: 'name',
						title: 'Branch Name',
					},
				]}
			/>
			<BranchSheet />
		</div>
	);
}

export default BranchesPage;
