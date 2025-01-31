import { ColumnDef } from '@tanstack/react-table';
import { Loan } from './schema';
import { format } from 'date-fns';
import DataTableColumnHeader from '@/components/table/data-table-column-header';
import { Badge } from '@/components/ui/badge';

export const loanColumns: ColumnDef<Loan>[] = [
	{
		accessorKey: 'id',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Loan ID" />
		),
		cell: ({ row }) => (
			<div className="text-center">{`LN${String(
				row.getValue('id'),
			).padStart(3, '0')}`}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		id: 'clientName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Client Name" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.client.fullName}</div>
		),
		filterFn: (row, columnId, filterValue) => {
			const fullName = row.original.client.fullName.toLowerCase();
			const loanOfficerName =
				row.original.loanOfficer.fullName.toLowerCase();
			return (
				fullName.includes(filterValue.toLowerCase()) ||
				loanOfficerName.includes(filterValue.toLowerCase())
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'repayAmount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Amount(Ksh)" />
		),
		cell: ({ row }) => (
			<div className="">
				{Number(row.original.product.repayAmount).toLocaleString()}
			</div>
		),
		enableSorting: true,
	},
	{
		id: 'loanOfficerName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Loan Officer" />
		),
		cell: ({ row }) => <div>{row.original.loanOfficer.fullName}</div>,
		filterFn: (row, columnId, filterValue) => {
			const fullName = row.original.client.fullName.toLowerCase();
			const loanOfficerName =
				row.original.loanOfficer.fullName.toLowerCase();
			return (
				fullName.includes(filterValue.toLowerCase()) ||
				loanOfficerName.includes(filterValue.toLowerCase())
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'disbursedOn',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Start Date" />
		),
		cell: ({ row }) => {
			const date: string = row.getValue('disbursedOn');
			return (
				<div className="w-50 truncate">
					{date !== '' ? format(date, 'dd MMM yyyy') : 'N/A'}
				</div>
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'dueDate',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Due Date" />
		),
		cell: ({ row }) => {
			const date: string = row.getValue('dueDate');
			return (
				<div className="w-50 truncate">
					{date !== '' ? format(date, 'dd MMM yyyy') : 'N/A'}
				</div>
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'status',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Status" />
		),
		cell: ({ row }) => {
			const status: string = row.getValue('status');
			const statusColors: Record<string, string> = {
				INACTIVE: 'text-gray-500',
				ACTIVE: 'text-green-500',
				COMPLETED: 'text-blue-500',
				DEFAULTED: 'text-red-500',
			};
			return (
				<Badge className={statusColors[status] || ''}>{status}</Badge>
			);
		},
		filterFn: (row, id, filterValues: string[]) => {
			const cellValue = row.getValue(id);
			return filterValues.includes(String(cellValue).toLowerCase());
		},
	},
	{
		accessorKey: 'feePaid',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Fee Paid" />
		),
		cell: ({ row }) => <div>{row.getValue('feePaid') ? 'Yes' : 'No'}</div>,
		enableSorting: false,
	},
	{
		accessorKey: 'paidAmount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Paid Amount" />
		),
		cell: ({ row }) => (
			<div className="">
				{Number(row.getValue('paidAmount')).toLocaleString()}
			</div>
		),
		enableSorting: true,
	},
];
