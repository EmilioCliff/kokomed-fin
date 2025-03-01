import { ColumnDef } from '@tanstack/react-table';
import { Loan, ExpectedPayment, UnpaidInstallment } from './schema';
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
];

export const expectedPaymentColumns: ColumnDef<ExpectedPayment>[] = [
	{
		accessorKey: 'loanId',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Loan ID" />
		),
		cell: ({ row }) => (
			<div className="text-center">{`LN${String(
				row.getValue('loanId'),
			).padStart(3, '0')}`}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		id: 'branchName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Branch Name" />
		),
		cell: ({ row }) => <div>{row.original.branchName}</div>,
		enableHiding: true,
	},
	{
		id: 'clientName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Client Name" />
		),
		cell: ({ row }) => <div>{row.original.clientName}</div>,
		filterFn: (row, filterValue) => {
			const fullName = row.original.clientName.toLowerCase();
			const loanOfficerName = row.original.loanOfficerName.toLowerCase();
			return (
				fullName.includes(filterValue.toLowerCase()) ||
				loanOfficerName.includes(filterValue.toLowerCase())
			);
		},
		enableSorting: true,
	},
	{
		id: 'loanOfficerName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Loan Officer" />
		),
		cell: ({ row }) => <div>{row.original.loanOfficerName}</div>,
		filterFn: (row, filterValue) => {
			const fullName = row.original.clientName.toLowerCase();
			const loanOfficerName = row.original.loanOfficerName.toLowerCase();
			return (
				fullName.includes(filterValue.toLowerCase()) ||
				loanOfficerName.includes(filterValue.toLowerCase())
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'loanAmount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Loan Amount" />
		),
		cell: ({ row }) => (
			<div>{Number(row.original.loanAmount).toLocaleString()}</div>
		),
		enableSorting: true,
	},
	{
		accessorKey: 'repayAmount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Repay Amount" />
		),
		cell: ({ row }) => (
			<div className="">
				{Number(row.original.repayAmount).toLocaleString()}
			</div>
		),
		enableSorting: true,
	},
	{
		accessorKey: 'totalUnpaid',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Total Unpaid" />
		),
		cell: ({ row }) => (
			<div className="">
				{Number(row.original.totalUnpaid).toLocaleString()}
			</div>
		),
		enableSorting: true,
	},
	{
		accessorKey: 'dueDate',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Due Date" />
		),
		cell: ({ row }) => (
			<div className="truncate">{row.getValue('dueDate')}</div>
		),
		enableSorting: true,
	},
];

export const unpaidInstallments: ColumnDef<UnpaidInstallment>[] = [
	{
		accessorKey: 'LoanId',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Loan ID" />
		),
		cell: ({ row }) => (
			<div className="text-center">{`LN${String(
				row.original.loanId,
			).padStart(3, '0')}`}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		id: 'branchName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Branch Name" />
		),
		cell: ({ row }) => <div>{row.original.branchName}</div>,
		enableHiding: true,
	},
	{
		id: 'clientName',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Client Name" />
		),
		cell: ({ row }) => <div>{row.original.fullName}</div>,
		filterFn: (row, filterValue) => {
			const fullName = row.original.fullName.toLowerCase();
			const phoneNumber = row.original.phoneNumber.toLowerCase();
			return (
				fullName.includes(filterValue.toLowerCase()) ||
				phoneNumber.includes(filterValue.toLowerCase())
			);
		},
		enableSorting: true,
	},
	{
		accessorKey: 'phoneNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Phone Number" />
		),
		cell: ({ row }) => (
			<div className="w-[80]">{row.original.phoneNumber}</div>
		),
		filterFn: (row, filterValue) => {
			const fullName = row.original.fullName.toLowerCase();
			const phoneNumber = row.original.phoneNumber.toLowerCase();
			return (
				fullName.includes(filterValue.toLowerCase()) ||
				phoneNumber.includes(filterValue.toLowerCase())
			);
		},
		enableSorting: true,
		enableHiding: true,
	},
	{
		accessorKey: 'loanAmount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Loan Amount" />
		),
		cell: ({ row }) => (
			<div>{Number(row.original.loanAmount).toLocaleString()}</div>
		),
		enableHiding: true,
	},
	{
		accessorKey: 'repayAmount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Repay Amount" />
		),
		cell: ({ row }) => (
			<div className="">
				{Number(row.original.repayAmount).toLocaleString()}
			</div>
		),
		enableHiding: true,
	},
	{
		accessorKey: 'paidAmount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Loan Paid Amount" />
		),
		cell: ({ row }) => (
			<div className="">
				{Number(row.original.paidAmount).toLocaleString()}
			</div>
		),
		enableHiding: true,
	},
	{
		accessorKey: 'installmentNumber',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Installment Number" />
		),
		cell: ({ row }) => (
			<div className="">{row.original.installmentNumber}</div>
		),
		enableHiding: true,
	},
	{
		accessorKey: 'amountDue',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Installment Amount" />
		),
		cell: ({ row }) => <div className="">{row.original.amountDue}</div>,
		enableHiding: true,
	},
	{
		accessorKey: 'remainingAmount',
		header: ({ column }) => (
			<DataTableColumnHeader
				column={column}
				title="Installment Remaining Amount"
			/>
		),
		cell: ({ row }) => (
			<div className="">{row.original.remainingAmount}</div>
		),
		enableHiding: true,
	},
	{
		accessorKey: 'totalDueAmount',
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title="Total Due Loan" />
		),
		cell: ({ row }) => (
			<div className="">{row.original.totalDueAmount}</div>
		),
		enableHiding: true,
	},
	{
		accessorKey: 'dueDate',
		header: ({ column }) => (
			<DataTableColumnHeader
				column={column}
				title="Installment Due Date"
			/>
		),
		cell: ({ row }) => (
			<div className="truncate">{row.original.dueDate}</div>
		),
		enableSorting: true,
	},
];
