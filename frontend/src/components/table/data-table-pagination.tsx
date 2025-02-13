import { Table } from '@tanstack/react-table';
import {
	ChevronLeft,
	ChevronRight,
	ChevronsLeft,
	ChevronsRight,
} from 'lucide-react';

import { Button } from '../ui/button';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '../ui/select';

import { useContext } from 'react';
import { TableContext } from '@/context/TableContext';

interface DataTablePaginationProps<TData> {
	table: Table<TData>;
}

export function DataTablePagination<TData>({
	table,
}: DataTablePaginationProps<TData>) {
	const { rowsCount, pageIndex, setPageSize, setPageIndex } =
		useContext(TableContext);
	const totalPages = Math.ceil(
		rowsCount / table.getState().pagination.pageSize,
	);

	return (
		<div className="flex items-center justify-between px-2 mt-4 mb-2">
			<div className="flex-1 text-sm text-muted-foreground"></div>
			<div className="flex items-center space-x-6 lg:space-x-8">
				<div className="flex items-center space-x-2">
					<p className="text-sm font-medium">Rows per page</p>
					<Select
						value={`${table.getState().pagination.pageSize}`}
						onValueChange={(value) => {
							table.setPageSize(Number(value));
							setPageSize(Number(value));
							setPageIndex(0);
						}}
					>
						<SelectTrigger className="h-8 w-[70px]">
							<SelectValue
								placeholder={
									table.getState().pagination.pageSize
								}
							/>
						</SelectTrigger>
						<SelectContent side="top">
							{[10, 20, 30, 40, 50].map((pageSize) => (
								<SelectItem
									key={pageSize}
									value={`${pageSize}`}
								>
									{pageSize}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>
				<div className="flex w-[100px] items-center justify-center text-sm font-medium">
					Page {pageIndex + 1} of {totalPages}
				</div>
				<div className="flex items-center space-x-2">
					<Button
						variant="outline"
						className="hidden h-8 w-8 p-0 lg:flex"
						onClick={() => setPageIndex(0)}
						disabled={pageIndex === 0}
					>
						<span className="sr-only">Go to first page</span>
						<ChevronsLeft />
					</Button>
					<Button
						variant="outline"
						className="h-8 w-8 p-0"
						onClick={() => setPageIndex(pageIndex - 1)}
						disabled={pageIndex === 0}
					>
						<span className="sr-only">Go to previous page</span>
						<ChevronLeft />
					</Button>
					<Button
						variant="outline"
						className="h-8 w-8 p-0"
						onClick={() => setPageIndex(pageIndex + 1)}
						disabled={pageIndex + 1 >= totalPages}
					>
						<span className="sr-only">Go to next page</span>
						<ChevronRight />
					</Button>
					<Button
						variant="outline"
						className="hidden h-8 w-8 p-0 lg:flex"
						onClick={() => setPageIndex(totalPages - 1)}
						disabled={pageIndex + 1 >= totalPages}
					>
						<span className="sr-only">Go to last page</span>
						<ChevronsRight />
					</Button>
				</div>
			</div>
		</div>
	);
}
