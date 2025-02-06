import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import VirtualizeddSelect from '../../UI/VisualizedSelect';
import Spinner from '@/components/UI/Spinner';
import addProduct from '@/services/addProduct';
import getFormData from '@/services/getFormData';
import { ProductFormType, productFormSchema } from './schema';
import { toast } from 'react-toastify';
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

interface ProductFormProps {
	onFormOpen: (isOpen: boolean) => void;
	onMutation: () => void;
}

function ProductForm({ onFormOpen, onMutation }: ProductFormProps) {
	const { isLoading, data, error } = useQuery({
		queryKey: ['loans/form'],
		queryFn: () => getFormData(false, false, false, true),
		staleTime: 5 * 1000,
	});

	const form = useForm<ProductFormType>({
		resolver: zodResolver(productFormSchema),
		defaultValues: {
			branchId: 0,
			loanAmount: 0,
			repayAmount: 0,
		},
	});
	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: addProduct,
	});

	function onSubmit(values: ProductFormType) {
		mutation.mutate(values, {
			onSuccess: () => {
				toast.success('Product Added Successful');
				onMutation();
				// queryClient.invalidateQueries({ queryKey: ['products'] });
			},
			onError: (error: any) => {
				toast.error(error.message);
			},
			onSettled: async () => {
				mutation.reset();
				return await queryClient.invalidateQueries({
					queryKey: ['products'],
				});
			},
		});
		onFormOpen(false);
	}

	function onError(errors: any) {
		console.log('Errors: ', errors);
	}

	return (
		<>
			{(mutation.isPending || isLoading) && <Spinner />}
			{error && <p>{error.message}</p>}

			{mutation.error && (
				<h5 onClick={() => mutation.reset()}>{`${mutation.error}`}</h5>
			)}

			<Form {...form}>
				<form
					onSubmit={form.handleSubmit(onSubmit, onError)}
					className="space-y-8"
				>
					<div className="grid grid-cols-1 md:grid-cols-2 gap-4 py-4">
						<FormField
							control={form.control}
							name="branchId"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Branch</FormLabel>
									<FormControl>
										{data?.branch && (
											<VirtualizeddSelect
												options={data.branch}
												placeholder="Select a product"
												value={field.value}
												onChange={(id) =>
													field.onChange(id)
												}
											/>
										)}
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="loanAmount"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Loan Amount(Ksh)</FormLabel>
									<FormControl>
										<Input type="number" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="repayAmount"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Repay Amount(Ksh)</FormLabel>
									<FormControl>
										<Input type="number" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					</div>
					<Button className="ml-auto block" type="submit">
						Add Product
					</Button>
				</form>
			</Form>
		</>
	);
}

export default ProductForm;
