import { ColumnDef } from '@tanstack/react-table';
import { Product } from './schema';
import DataTableColumnHeader from '../../table/data-table-column-header';

export const productColumns: ColumnDef<Product>[] = [
  {
    accessorKey: 'id',
    header: ({ column }) => <DataTableColumnHeader column={column} title="ID" />,
    cell: ({ row }) => (
      <div className="text-center">{`${String(row.getValue('id')).padStart(
        3,
        '0'
      )}`}</div>
    ),
    enableSorting: true,
    enableHiding: true,
  },
  {
    accessorKey: 'branchName',
    header: ({ column }) => <DataTableColumnHeader column={column} title="Branch Name" />,
    cell: ({ row }) => <div className="w-[80]">{row.getValue('branchName')}</div>,
    filterFn: (row, columnId, filterValue) => {
      const name = row.original.branchName.toLowerCase();
      return name.includes(filterValue.toLowerCase());
    },
    enableSorting: true,
  },
  {
    accessorKey: 'loanAmount',
    header: ({ column }) => <DataTableColumnHeader column={column} title="Loan Amount" />,
    cell: ({ row }) => (
      <div className="">{Number(row.getValue('loanAmount')).toLocaleString()}</div>
    ),
    enableSorting: true,
  },
  {
    accessorKey: 'repayAMount',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Repay Amount" />
    ),
    cell: ({ row }) => (
      <div className="">{Number(row.getValue('repayAMount')).toLocaleString()}</div>
    ),
    // enableSorting: true,
  },
  {
    accessorKey: 'interestAmount',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Interest Amount" />
    ),
    cell: ({ row }) => (
      <div className="">{Number(row.getValue('interestAmount')).toLocaleString()}</div>
    ),
    // enableSorting: true,
  },
];
