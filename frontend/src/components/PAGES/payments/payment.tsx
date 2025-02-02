import { ColumnDef } from '@tanstack/react-table';
import { Payment } from './schema';
import { format } from 'date-fns';
import DataTableColumnHeader from '@/components/table/data-table-column-header';
import { Badge } from '@/components/ui/badge';

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
			console.log(filterValues);
			console.log(cellValue);
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
