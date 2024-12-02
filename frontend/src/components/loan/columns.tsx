import { ColumnDef } from "@tanstack/react-table";
import { Loan } from "@/data/schema";
import DataTableColumnHeader from "./data-table-column-header";
import { Badge } from "@/components/ui/badge";

export const loanColumns: ColumnDef<Loan>[] = [
	{
		accessorKey: "id",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Loan ID' />
		),
		cell: ({ row }) => (
			<div className='w-[80px]'>{`LN${String(row.getValue("id")).padStart(
				3,
				"0"
			)}`}</div>
		),
		enableSorting: true,
		enableHiding: false,
	},
	// // Loan amount column
	// {
	// 	accessorKey: "amount",
	// 	header: ({ column }) => (
	// 		<DataTableColumnHeader column={column} title='Amount' />
	// 	),
	// 	cell: ({ row }) => (
	// 		<div className='text-right'>${row.getValue("amount").toFixed(2)}</div>
	// 	),
	// 	enableSorting: true,
	// },

	{
		accessorKey: "repayAmount",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Repay Amount' />
		),
		cell: ({ row }) => (
			<div className='text-right'>
				{/* .toFixed(2) */}${row.getValue("repayAmount")}
			</div>
		),
		enableSorting: true,
	},
	{
		id: "clientName",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Client Name' />
		),
		cell: ({ row }) => <div>{row.original.client.fullName}</div>,
		enableSorting: true,
	},
	{
		id: "loanOfficerName",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Loan Officer' />
		),
		cell: ({ row }) => <div>{row.original.loanOfficer.fullName}</div>,
		enableSorting: true,
	},
	// {
	// 	accessorKey: "loanPurpose",
	// 	header: ({ column }) => (
	// 		<DataTableColumnHeader column={column} title='Loan Purpose' />
	// 	),
	// 	cell: ({ row }) => <div>{row.getValue("loanPurpose") || "N/A"}</div>,
	// 	enableSorting: false,
	// },
	{
		accessorKey: "disbursedOn",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Start Date' />
		),
		cell: ({ row }) => (
			<div>
				{row.getValue("disbursedOn")
					? new Date(row.getValue("disbursedOn")).toLocaleDateString()
					: "N/A"}
			</div>
		),
		enableSorting: true,
	},
	{
		accessorKey: "dueDate",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Due Date' />
		),
		cell: ({ row }) => (
			<div>
				{row.getValue("dueDate")
					? new Date(row.getValue("dueDate")).toLocaleDateString()
					: "N/A"}
			</div>
		),
		enableSorting: true,
	},
	{
		accessorKey: "status",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Status' />
		),
		cell: ({ row }) => {
			const status = row.getValue("status");
			const statusColors: Record<string, string> = {
				INACTIVE: "text-gray-500",
				ACTIVE: "text-green-500",
				COMPLETED: "text-blue-500",
				DEFAULTED: "text-red-500",
			};
			return (
				<Badge className={statusColors[String(status)] || ""}>
					{String(status)}
				</Badge>
			);
		},
		filterFn: (row, id, value) => value.includes(row.getValue(id)),
	},
	{
		accessorKey: "feePaid",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Fee Paid' />
		),
		cell: ({ row }) => <div>{row.getValue("feePaid") ? "Yes" : "No"}</div>,
		enableSorting: false,
	},
	{
		accessorKey: "paidAmount",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Paid Amount' />
		),
		cell: ({ row }) => (
			<div className='text-right'>${row.getValue("paidAmount")}</div>
		),
		enableSorting: true,
	},
];
