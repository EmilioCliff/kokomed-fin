import { generateRandomProduct } from '@/lib/generator';
import { z } from 'zod';
import { Product, productSchema } from './schema';
import { useEffect, useState, useContext } from 'react';
import { TableContext } from '@/context/TableContext';
import TableSkeleton from '@/components/UI/TableSkeleton';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { DataTable } from '../../table/data-table';
import ProductForm from '@/components/PAGES/products/ProductForm';
import { productColumns } from '@/components/PAGES/products/product';

// const products = generateRandomProduct(30);
// const validatedProducts = z.array(productSchema).parse(products);

function ProductsPage() {
  const { selectedRow, setSelectedRow } = useContext(TableContext);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  // const [selectedRow, setSelectedRow] = useState<Product | null>(null);

  useEffect(() => {
    // setProducts(validatedProducts);
    setLoading(false);
  }, []);

  if (loading) {
    return <TableSkeleton />;
  }

  return (
    <div className="px-4">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-3xl font-bold">Products</h1>
        <Dialog>
          <DialogTrigger asChild>
            <Button className="text-xs py-1 font-bold" size="sm">
              Add New Product
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-screen-lg ">
            <DialogHeader>
              <DialogTitle>Add New Product</DialogTitle>
              <DialogDescription>Enter the details for the new user.</DialogDescription>
            </DialogHeader>
            <ProductForm />
          </DialogContent>
        </Dialog>
      </div>
      <DataTable
        data={products}
        columns={productColumns}
        // setSelectedRow={setSelectedRow}
        searchableColumns={[
          {
            id: 'branchName',
            title: 'branch name',
          },
        ]}
      />
      <Sheet
        open={!!selectedRow}
        onOpenChange={(open: boolean) => {
          if (!open) {
            setSelectedRow(null);
          }
        }}
      >
        <SheetContent className="overflow-auto custom-sheet-class">
          <SheetHeader>
            <SheetTitle>Product Details</SheetTitle>
            <SheetDescription>Description goes here</SheetDescription>
          </SheetHeader>
          {selectedRow && (
            <div className="py-4">
              {Object.entries(selectedRow).map(([key, value]) => {
                // if (key === "createdBy" || key === "updatedBy") {
                // 	return;
                // }
                // if (fieldRenderers[key]) {
                // 	return (
                // 		<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
                // 			<span className='font-medium capitalize'>{key}</span>
                // 			{fieldRenderers[key](value)}
                // 		</div>
                // 	);
                // }

                // return (
                // 	<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
                // 		<span className='font-medium capitalize'>{key}</span>
                // 		{typeof value === "string" ||
                // 		typeof value === "number" ||
                // 		typeof value === "boolean" ? (
                // 			<Input
                // 				readOnly
                // 				placeholder={value.toString()}
                // 				className='bg-gray-100 text-gray-500'
                // 			/>
                // 		) : (
                // 			JSON.stringify(value)
                // 		)}
                // 	</div>
                // );
                return (
                  <div key={key} className="grid grid-cols-[0.5fr_1fr] mb-4">
                    <span className="font-medium capitalize">{key}</span>
                    <p>
                      {typeof value == 'string' ? <p>{value}</p> : JSON.stringify(value)}
                    </p>
                  </div>
                );
              })}
              <Button size="lg" onClick={() => console.log('CLicked')} className="mt-8">
                Save
              </Button>
            </div>
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}

export default ProductsPage;
