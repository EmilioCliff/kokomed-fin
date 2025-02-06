import { ColumnDef } from '@tanstack/react-table';
import { User } from './schema';
import DataTableColumnHeader from '../../table/data-table-column-header';
import { format } from 'date-fns';
import { Badge } from '@/components/ui/badge';

export const userColumns: ColumnDef<User>[] = [
	{
		accessorKey: 'id',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="User ID" />
		),
		cell: ({ row }) => (
			<div className="text-center">{`AG${String(
				row.getValue('id'),
			).padStart(3, '0')}`}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		accessorKey: 'branchName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Branch Name" />
		),
		cell: ({ row }) => (
			<div className="w-[80] truncate">{row.getValue('branchName')}</div>
		),
		enableSorting: false,
	},
	{
		accessorKey: 'fullName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Username" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.getValue('fullName')}</div>
		),
		filterFn: (row, filterValue) => {
			const fullName = row.original.fullName.toLowerCase();
			const email = row.original.email.toLowerCase();
			return (
				fullName.includes(filterValue.toLowerCase()) ||
				email.includes(filterValue.toLowerCase())
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'email',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Email" />
		),
		cell: ({ row }) => <div className="">{row.getValue('email')}</div>,
		filterFn: (row, filterValue) => {
			const fullName = row.original.fullName.toLowerCase();
			const email = row.original.email.toLowerCase();
			return (
				fullName.includes(filterValue.toLowerCase()) ||
				email.includes(filterValue.toLowerCase())
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'phoneNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Phonenumber" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.getValue('phoneNumber')}</div>
		),
		enableHiding: true,
		enableSorting: false,
	},
	{
		accessorKey: 'role',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Role" />
		),
		cell: ({ row }) => {
			const role: string = row.getValue('role');
			return (
				<Badge variant={role === 'ADMIN' ? 'default' : 'secondary'}>
					{role}
				</Badge>
			);
		},
		filterFn: (row, filterValue) => {
			const role = row.original.role.toLowerCase();
			return role.includes(filterValue);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'createdAt',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Created At" />
		),
		cell: ({ row }) => (
			<div className="w-[50] truncate">
				{format(row.getValue('createdAt'), 'dd/MM/yyyy')}
			</div>
		),
		enableHiding: true,
		enableSorting: false,
	},
];
