import { ColumnDef } from "@tanstack/react-table";
import { Client } from "@/data/schema";
import { Badge } from "@/components/ui/badge";
import { format } from "date-fns";
import DataTableColumnHeader from "@/components/table/data-table-column-header";

export const clientColumns: ColumnDef<Client>[] = [
	{
		accessorKey: "id",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Client ID' />
		),
		cell: ({ row }) => (
			<div className='text-center'>{`${String(row.getValue("id")).padStart(
				3,
				"0"
			)}`}</div>
		),
		enableSorting: true,
		enableHiding: true,
	},
	{
		accessorKey: "branchName",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Branch Name' />
		),
		cell: ({ row }) => (
			<div className='w-[80] truncate'>{row.getValue("branchName")}</div>
		),
	},
	{
		accessorKey: "fullName",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Client Name' />
		),
		cell: ({ row }) => <div className='w-[80]'>{row.getValue("fullName")}</div>,
		filterFn: (row, columnId, filterValue) => {
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
		accessorKey: "phoneNumber",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Phone Number' />
		),
		cell: ({ row }) => (
			<div className='w-[80]'>{row.getValue("phoneNumber")}</div>
		),
		filterFn: (row, columnId, filterValue) => {
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
		accessorKey: "idNumber",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='ID Number' />
		),
		cell: ({ row }) => (
			<div className='w-[80]'>
				{row.getValue("idNumber") ? row.getValue("idNumber") : "-"}
			</div>
		),
		enableHiding: true,
	},
	{
		accessorKey: "gender",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Gender' />
		),
		cell: ({ row }) => <div>{row.getValue("gender")}</div>,
		enableHiding: true,
	},
	{
		accessorKey: "active",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Active' />
		),
		cell: ({ row }) => {
			const status: boolean = row.getValue("active");
			return (
				<Badge variant={status ? "default" : "destructive"}>
					{status ? "Active" : "Inactive"}
				</Badge>
			);
		},
		filterFn: (row, columnId, filterValue) => {
			return filterValue.includes(Boolean(row.getValue(columnId)).toString());
		},
	},
	{
		id: "assignedStaff",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Assigned Staff' />
		),
		cell: ({ row }) => (
			<div className='truncate'>{row.original.assignedStaff.fullName}</div>
		),
		enableHiding: true,
	},
	{
		id: "dueAmount",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Loan Balance' />
		),
		cell: ({ row }) => <div>{row.original.dueAmount}</div>,
		enableHiding: true,
	},
	{
		accessorKey: "createdAt",
		header: ({ column }) => (
			<DataTableColumnHeader column={column} title='Created At' />
		),
		cell: ({ row }) => (
			<div className='w-[50] truncate'>
				{format(row.getValue("createdAt"), "dd MMM yyyy")}
			</div>
		),
		enableHiding: true,
	},
];
