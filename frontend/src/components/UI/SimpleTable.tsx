import React, { useState, useMemo } from 'react';
import {
	Table as ShadcnTable,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { ChevronDown, ChevronUp, Download } from 'lucide-react';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

interface Column {
	key: string;
	label: string;
	sortable?: boolean;
	render?: (value: any) => React.ReactNode;
}

interface TableProps {
	columns: Column[];
	data: any[];
	onRowClick?: (row: any) => void;
}

export function SimpleTable({ columns, data, onRowClick }: TableProps) {
	const [sortColumn, setSortColumn] = useState<string | null>(null);
	const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');
	// const [selectedRow, setSelectedRow] = useState<any | null>(null);

	const handleSort = (column: string) => {
		if (sortColumn === column) {
			setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
		} else {
			setSortColumn(column);
			setSortDirection('asc');
		}
	};

	const sortedData = useMemo(() => {
		if (!sortColumn) return data;

		return [...data].sort((a, b) => {
			if (a[sortColumn] < b[sortColumn])
				return sortDirection === 'asc' ? -1 : 1;
			if (a[sortColumn] > b[sortColumn])
				return sortDirection === 'asc' ? 1 : -1;
			return 0;
		});
	}, [data, sortColumn, sortDirection]);

	return (
		<div>
			<div className="mb-4 flex justify-between items-center">
				<h2 className="text-2xl font-bold">Data Table</h2>
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="outline">
							<Download className="mr-2 h-4 w-4" />
							Download
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent>
						<DropdownMenuItem>PDF</DropdownMenuItem>
						<DropdownMenuItem>Excel</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</div>
			<ShadcnTable>
				<TableHeader>
					<TableRow>
						{columns.map((column) => (
							<TableHead key={column.key}>
								<div className="flex items-center">
									{column.label}
									{column.sortable && (
										<Button
											variant="ghost"
											size="sm"
											className="ml-2"
											onClick={() =>
												handleSort(column.key)
											}
										>
											{sortColumn === column.key ? (
												sortDirection === 'asc' ? (
													<ChevronUp className="h-4 w-4" />
												) : (
													<ChevronDown className="h-4 w-4" />
												)
											) : (
												<ChevronDown className="h-4 w-4" />
											)}
										</Button>
									)}
								</div>
							</TableHead>
						))}
					</TableRow>
				</TableHeader>
				<TableBody>
					{sortedData.map((row, index) => (
						<TableRow
							key={index}
							onClick={() => onRowClick && onRowClick(row)}
							className="cursor-pointer"
						>
							{columns.map((column) => (
								<TableCell key={column.key}>
									{column.render
										? column.render(row[column.key])
										: row[column.key]}
								</TableCell>
							))}
						</TableRow>
					))}
				</TableBody>
			</ShadcnTable>
		</div>
	);
}
