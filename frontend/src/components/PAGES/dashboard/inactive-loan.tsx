import { ColumnDef } from '@tanstack/react-table';
import { InactiveLoan } from './schema';
import DataTableColumnHeader from '@/components/table/data-table-column-header';
import { formatDate } from '@/lib/utils';

export const inactiveLoanColumns: ColumnDef<InactiveLoan>[] = [
  {
    accessorKey: 'id',
    header: ({ column }) => <DataTableColumnHeader column={column} title="Loan ID" />,
    cell: ({ row }) => (
      <div className="text-center">{`LN${String(row.getValue('id')).padStart(
        3,
        '0'
      )}`}</div>
    ),
    enableSorting: true,
    enableHiding: true,
  },
  {
    id: 'clientName',
    header: ({ column }) => <DataTableColumnHeader column={column} title="Client Name" />,
    cell: ({ row }) => <div className="w-[80]">{row.original.client.fullName}</div>,
    enableSorting: true,
    filterFn: (row, columnId, filterValue) => {
      const fullName = row.original.client.fullName.toLowerCase();
      return fullName.includes(filterValue.toLowerCase());
    },
  },
  {
    accessorKey: 'amount',
    header: ({ column }) => <DataTableColumnHeader column={column} title="Amount(Ksh)" />,
    cell: ({ row }) => <div>{Number(row.getValue('amount')).toLocaleString()}</div>,
    enableSorting: true,
  },
  {
    accessorKey: 'repayAmount',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Repay Amount(Ksh)" />
    ),
    cell: ({ row }) => <div>{Number(row.getValue('repayAmount')).toLocaleString()}</div>,
    enableSorting: true,
  },
  {
    id: 'loanOfficerName',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Loan Officer" />
    ),
    cell: ({ row }) => <div>{row.original.loanOfficer.fullName}</div>,
    enableSorting: true,
  },
  {
    id: 'approvedByName',
    header: ({ column }) => <DataTableColumnHeader column={column} title="Approved By" />,
    cell: ({ row }) => <div>{row.original.approvedBy.fullName}</div>,
    enableSorting: true,
  },
  {
    accessorKey: 'approvedOn',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Approved Date" />
    ),
    cell: ({ row }) => (
      <div>
        {row.getValue('approvedOn')
          ? new Date(row.getValue('approvedOn')).toLocaleDateString()
          : 'N/A'}
      </div>
    ),
    enableSorting: true,
  },
];
