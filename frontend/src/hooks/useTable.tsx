import { useContext } from 'react';
import { TableContext, TableContextType } from '@/context/TableContext';

export const useTable = (): TableContextType => {
	const context = useContext(TableContext);

	if (!context) {
		throw new Error('useTable must be used within an TableProvider');
	}
	return context;
};
