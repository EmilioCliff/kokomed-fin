import { ColumnDef } from '@tanstack/react-table';
import { Payment } from './schema';
import { format } from 'date-fns';
import DataTableColumnHeader from '@/components/table/data-table-column-header';
import { Badge } from '@/components/ui/badge';
import { Loan } from '../loans/schema';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import { Eye, MoreHorizontal } from 'lucide-react';
import { Link, useNavigate } from 'react-router';

export const paymentColumns: ColumnDef<Payment>[] = [
	{
		accessorKey: 'id',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Payment ID" />
		),
		cell: ({ row }) => (
			<div className="text-center">{`P${String(
				row.getValue('id'),
			).padStart(3, '0')}`}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		accessorKey: 'transactionNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Transaction Number" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.transactionNumber}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		accessorKey: 'accountNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Account Number" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.accountNumber}</div>
		),
		enableSorting: false,
		enableHiding: true,
	},
	{
		accessorKey: 'phoneNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Phone Number" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.phoneNumber}</div>
		),
		enableSorting: false,
		enableHiding: true,
	},
	{
		accessorKey: 'payingName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Paying Name" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.payingName}</div>
		),
		enableSorting: false,
	},
	{
		accessorKey: 'amount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Amount" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">
				{Number(row.original.amount).toLocaleString()}
			</div>
		),
		enableSorting: false,
	},
	{
		accessorKey: 'assigned',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Assigned" />
		),
		cell: ({ row }) => (
			<div className="">{row.original.assigned ? 'YES' : 'NO'}</div>
		),
		filterFn: (row, id, filterValues: string[]) => {
			const cellValue = row.getValue(id);
			return filterValues.includes(String(cellValue));
		},
		enableHiding: true,
	},
	{
		accessorKey: 'transactionSource',
		header: ({ column }) => (
			<DataTableColumnHeader
				column={column}
				className="w-50"
				title="Transaction Source"
			/>
		),
		cell: ({ row }) => (
			// <div className="w-[80]">{row.original.transactionSource}</div>
			<Badge
				variant={
					row.original.transactionSource === 'INTERNAL'
						? 'default'
						: 'secondary'
				}
			>
				{row.original.transactionSource}
			</Badge>
		),
		filterFn: (row, id, filterValues: string[]) => {
			const cellValue = row.getValue(id);
			return filterValues.includes(String(cellValue));
		},
		enableHiding: true,
	},
	{
		accessorKey: 'paidDate',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Paid Date" />
		),
		cell: ({ row }) => {
			return (
				<div className="w-50 truncate">
					{format(row.original.paidDate, 'dd MMM yyyy')}
				</div>
			);
		},
		enableSorting: true,
	},
];

export const clientLoanColumns: ColumnDef<Loan>[] = [
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
		filterFn: (row, filterValue) => {
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

			if (!date || date === '0001-01-01T00:00:00Z') {
				return <div className="w-50 truncate">-</div>;
			}

			return (
				<div className="w-50 truncate">
					{format(date, 'dd MMM yyyy')}
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

			if (!date || date === '0001-01-01T00:00:00Z') {
				return <div className="w-50 truncate">-</div>;
			}

			return (
				<div className="w-50 truncate">
					{format(date, 'dd MMM yyyy')}
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
	{
		accessorKey: 'view',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="View" />
		),
		cell: ({ row }) => (
			<Link to={`/loans/overview/${row.original.id}`} className="w-full">
				<Button variant="ghost" size="sm" className="h-8 w-8 p-0">
					<Eye className="h-4 w-4" />
					<span className="sr-only">View Loan</span>
				</Button>
			</Link>
		),
		enableSorting: false,
	},
];

export const clientPaymentColumns: ColumnDef<Payment>[] = [
	{
		accessorKey: 'id',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Payment ID" />
		),
		cell: ({ row }) => (
			<div className="text-center">{`P${String(
				row.getValue('id'),
			).padStart(3, '0')}`}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		accessorKey: 'transactionNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Transaction Number" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.transactionNumber}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		accessorKey: 'accountNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Account Number" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.accountNumber}</div>
		),
		enableSorting: false,
		enableHiding: true,
	},
	{
		accessorKey: 'phoneNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Phone Number" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.phoneNumber}</div>
		),
		enableSorting: false,
		enableHiding: true,
	},
	{
		accessorKey: 'payingName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Paying Name" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.payingName}</div>
		),
		enableSorting: false,
	},
	{
		accessorKey: 'amount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Amount" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">
				{Number(row.original.amount).toLocaleString()}
			</div>
		),
		enableSorting: false,
	},
	{
		accessorKey: 'transactionSource',
		header: ({ column }) => (
			<DataTableColumnHeader
				column={column}
				className="w-50"
				title="Transaction Source"
			/>
		),
		cell: ({ row }) => (
			// <div className="w-[80]">{row.original.transactionSource}</div>
			<Badge
				variant={
					row.original.transactionSource === 'INTERNAL'
						? 'default'
						: 'secondary'
				}
			>
				{row.original.transactionSource}
			</Badge>
		),
		filterFn: (row, id, filterValues: string[]) => {
			const cellValue = row.getValue(id);
			return filterValues.includes(String(cellValue));
		},
		enableHiding: true,
	},
	{
		accessorKey: 'paidDate',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Paid Date" />
		),
		cell: ({ row }) => {
			return (
				<div className="w-50 truncate">
					{format(row.original.paidDate, 'dd MMM yyyy')}
				</div>
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'assignedBy',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Assigned By" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">
				{row.original.assignedBy === 'APP'
					? 'System'
					: row.original.assignedBy}
			</div>
		),
		enableSorting: false,
	},
	{
		accessorKey: 'actions',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Actions" />
		),
		cell: ({ row }) => (
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost" className="h-8 w-8 p-0">
						<span className="sr-only">Open menu</span>
						<MoreHorizontal className="h-4 w-4" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuLabel>Actions</DropdownMenuLabel>
					<DropdownMenuItem
						disabled={
							row.original.transactionSource === 'MPESA'
								? true
								: false
						}
					>
						<Link
							to={`/payments/overview/${row.original.id}`}
							className="w-full"
						>
							Edit Payment
						</Link>
					</DropdownMenuItem>
					<DropdownMenuSeparator />
					<DropdownMenuItem className=" hover:bg-destructive">
						<Link
							to={`/sms/new?customer=${row.original.id}`}
							className="w-full"
						>
							Delete
						</Link>
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		),
		enableSorting: false,
	},
];
