import { Table, Column } from '@tanstack/react-table';
import { useState, useMemo, useContext } from 'react';
import { TableContext } from '@/context/TableContext';
import { X } from 'lucide-react';

import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { DataTableViewOptions } from './data-table-view-options';

import DataTableFacetedFilter from './data-table-faceted-filter';

interface DataTableToolbarProps<TData> {
  table: Table<TData>;
  searchableColumns?: {
    id: string;
    title: string;
  }[];
  facetedFilterColumns?: {
    id: string;
    title: string;
    options: {
      label: string;
      value: string;
    }[];
  }[];
}

export function DataTableToolbar<TData>({
  table,
  searchableColumns = [],
  facetedFilterColumns = [],
}: DataTableToolbarProps<TData>) {
  const { search, filter, setSearch, setFilter, setPageIndex } = useContext(TableContext);

  // const [searchValue, setSearchValue] = useState('');

  // const isFiltered = table.getState().columnFilters.length > 0;

  const isFiltered = filter.options.length > 0 || search;

  // const searchableColumnIds = useMemo(
  //   () => searchableColumns.map((column) => column.id),
  //   [searchableColumns]
  // );

  const handleSearch = (value: string) => {
    setPageIndex(0);
    setSearch(value);
    // searchableColumnIds.forEach((columnId) => {
    //   table.getColumn(columnId)?.setFilterValue(value);
    // });
  };

  return (
    <div className="flex items-center justify-between mt-2 mb-6">
      <div className="flex flex-1 items-center space-x-2">
        <Input
          placeholder={`Filter ${searchableColumns
            .map((column) => column.title)
            .join(' or ')}...`}
          value={search}
          onChange={(event) => handleSearch(event.target.value)}
          className="h-8 w-[150px] lg:w-[250px]"
        />
        {facetedFilterColumns.map((column) => {
          const tableColumn = table.getColumn(column.id) as Column<TData, unknown>;
          return (
            tableColumn && (
              <DataTableFacetedFilter
                key={column.id}
                column={tableColumn}
                title={column.title}
                options={column.options}
              />
            )
          );
        })}
        {isFiltered && (
          <Button
            variant="ghost"
            onClick={() => {
              table.resetColumnFilters();
              setSearch('');
              setFilter({ options: [] });
            }}
            className="h-8 px-2 lg:px-3"
          >
            Reset
            <X />
          </Button>
        )}
      </div>
      <DataTableViewOptions table={table} />
    </div>
  );
}
