import { useState, useContext } from 'react';
import { TableContext } from '@/context/TableContext';
import TableSkeleton from '@/components/UI/TableSkeleton';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { DataTable } from '../../table/data-table';
import ProductForm from '@/components/PAGES/products/ProductForm';
import { productColumns } from '@/components/PAGES/products/product';
import { useDebounce } from '@/hooks/useDebounce';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import getProducts from '@/services/getProducts';
import ProductSheet from './ProductSheet';

function ProductsPage() {
	const [formOpen, setFormOpen] = useState(false);
	const { pageIndex, pageSize, search, updateTableContext } =
		useContext(TableContext);

	const debouncedInput = useDebounce({ value: search, delay: 500 });

	const { isLoading, error, data, refetch } = useQuery({
		queryKey: ['products', pageIndex, pageSize, debouncedInput],
		queryFn: () => getProducts(pageIndex, pageSize, debouncedInput),
		staleTime: 0,
		placeholderData: keepPreviousData,
	});

	if (data?.metadata) {
		updateTableContext(data.metadata);
	}

	if (isLoading) {
		return <TableSkeleton />;
	}

	if (error) {
		return <div>Error: {error.message}</div>;
	}

	return (
		<div className="px-4">
			<div className="flex justify-between items-center mb-4">
				<h1 className="text-3xl font-bold">Products</h1>
				<Dialog open={formOpen} onOpenChange={setFormOpen}>
					<DialogTrigger asChild>
						<Button className="text-xs py-1 font-bold" size="sm">
							Add New Product
						</Button>
					</DialogTrigger>
					<DialogContent className="max-w-screen-lg ">
						<DialogHeader>
							<DialogTitle>Add New Product</DialogTitle>
							<DialogDescription>
								Enter the details for the new product.
							</DialogDescription>
						</DialogHeader>
						<ProductForm onFormOpen={setFormOpen} />
					</DialogContent>
				</Dialog>
			</div>
			<DataTable
				data={data?.data || []}
				columns={productColumns}
				searchableColumns={[
					{
						id: 'branchName',
						title: 'branch name',
					},
				]}
			/>
			<ProductSheet />
		</div>
	);
}

export default ProductsPage;
