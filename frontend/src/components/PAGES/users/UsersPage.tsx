import { useState } from 'react';
import TableSkeleton from '@/components/UI/TableSkeleton';
import { roles } from '@/data/loan';
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
import UserForm from '@/components/PAGES/users/UserForm';
import { userColumns } from '@/components/PAGES/users/user';
import { useDebounce } from '@/hooks/useDebounce';
import { useTable } from '@/hooks/useTable';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import getUsers from '@/services/getUsers';
import UserSheet from './UserSheet';
import { useAuth } from '@/hooks/useAuth';
import { role } from '@/lib/types';
import { toast } from 'react-toastify';

function UsersPage() {
	const { decoded } = useAuth();
	const [formOpen, setFormOpen] = useState(false);
	const { pageIndex, pageSize, filter, search, updateTableContext } =
		useTable();

	const debouncedInput = useDebounce({ value: search, delay: 500 });

	const { isLoading, error, data } = useQuery({
		queryKey: ['users', pageIndex, pageSize, filter, debouncedInput],
		queryFn: () => getUsers(pageIndex, pageSize, filter, debouncedInput),
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
				<h1 className="text-3xl font-bold">Users</h1>
				<Dialog
					open={formOpen}
					onOpenChange={() => {
						if (!formOpen && decoded?.role !== role.ADMIN) {
							toast.error('only admins can add users');
							return;
						}
						setFormOpen(!formOpen);
					}}
				>
					<DialogTrigger asChild>
						<Button className="text-xs py-1 font-bold" size="sm">
							Add New User
						</Button>
					</DialogTrigger>
					<DialogContent className="max-w-screen-lg max-h-[90vh] overflow-y-auto">
						<DialogHeader>
							<DialogTitle>Add New User</DialogTitle>
							<DialogDescription>
								Enter the details for the new user.
							</DialogDescription>
						</DialogHeader>
						<UserForm onFormOpen={setFormOpen} />
					</DialogContent>
				</Dialog>
			</div>
			<DataTable
				data={data?.data || []}
				columns={userColumns}
				searchableColumns={[
					{
						id: 'fullName',
						title: 'username',
					},
					{
						id: 'email',
						title: 'email',
					},
				]}
				facetedFilterColumns={[
					{
						id: 'role',
						title: 'Role',
						options: roles,
					},
				]}
			/>
			<UserSheet />
		</div>
	);
}

export default UsersPage;
