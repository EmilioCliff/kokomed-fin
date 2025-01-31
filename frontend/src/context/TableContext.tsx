import { createContext, FC, useState } from 'react';
import { contextWrapperProps } from '@/lib/types';
import { tableCtx, tableFilterType, pagination } from '@/lib/types';

const defaultContext: tableCtx = {
	search: '',
	filter: { options: [] },
	pageIndex: 0,
	pageSize: 10,
	rowsCount: 0,
	selectedRow: null,
	setSearch: () => {},
	setFilter: () => {},
	setPageIndex: () => {},
	setPageSize: () => {},
	setRowsCount: () => {},
	setSelectedRow: () => {},
	resetTableState: () => {},
	updateTableContext: () => {},
};

export const TableContext = createContext<tableCtx>(defaultContext);

export const TableContextWrapper: FC<contextWrapperProps> = ({ children }) => {
	// store table state being used for calls
	const [search, setSearch] = useState(defaultContext.search);
	const [filter, setFilter] = useState<tableFilterType>(
		defaultContext.filter,
	);
	const [pageIndex, setPageIndex] = useState(defaultContext.pageIndex);
	const [pageSize, setPageSize] = useState(defaultContext.pageSize);
	const [rowsCount, setRowsCount] = useState(defaultContext.rowsCount);
	const [selectedRow, setSelectedRow] = useState(defaultContext.selectedRow);

	const resetTableState = () => {
		setSearch(defaultContext.search);
		setFilter(defaultContext.filter);
		setPageIndex(defaultContext.pageIndex);
		setPageSize(defaultContext.pageSize);
		setSelectedRow(defaultContext.selectedRow);
	};

	const updateTableContext = (pagination: pagination | undefined) => {
		if (pagination) {
			setRowsCount(pagination.totalData);
		}
	};

	return (
		<TableContext.Provider
			value={{
				search,
				filter,
				pageIndex,
				pageSize,
				rowsCount,
				selectedRow,
				setSearch,
				setFilter,
				setPageIndex,
				setPageSize,
				setRowsCount,
				setSelectedRow,
				resetTableState,
				updateTableContext,
			}}
		>
			{children}
		</TableContext.Provider>
	);
};
