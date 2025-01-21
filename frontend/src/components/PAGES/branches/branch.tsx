import { ColumnDef } from '@tanstack/react-table';
import { Branch } from './schema';
import DataTableColumnHeader from '../../table/data-table-column-header';

export const branchColumns: ColumnDef<Branch>[] = [
  {
    accessorKey: 'id',
    header: ({ column }) => <DataTableColumnHeader column={column} title="ID" />,
    cell: ({ row }) => (
      <div className="text-center">{`LN${String(row.getValue('id')).padStart(
        3,
        '0'
      )}`}</div>
    ),
    enableSorting: true,
  },
  {
    id: 'branchName',
    header: ({ column }) => <DataTableColumnHeader column={column} title="Branch Name" />,
    cell: ({ row }) => <div className="w-[80]">{row.original.branchName}</div>,
    filterFn: (row, columnId, filterValue) => {
      const name = row.original.branchName.toLowerCase();
      return name.includes(filterValue.toLowerCase());
    },
    enableSorting: true,
  },
];
