import { useMutation, useQueryClient } from '@tanstack/react-query';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { BranchFormType, branchFormSchema } from './schema';
import { toast } from 'react-toastify';
import addBranch from '@/services/addBranch';
import Spinner from '@/components/UI/Spinner';
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

interface LoanFormProps {
	onFormOpen: (isOpen: boolean) => void;
}

function BranchForm({ onFormOpen }: LoanFormProps) {
	const form = useForm<BranchFormType>({
		resolver: zodResolver(branchFormSchema),
		defaultValues: {
			name: '',
		},
	});

	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: addBranch,
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: ['branches'] });
			toast.success('Branch Added Successful');
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
		onSettled: () => onFormOpen(false),
	});

	function onSubmit(values: BranchFormType) {
		mutation.mutate(values);
	}

	function onError(errors: any) {
		console.log('Errors: ', errors);
	}

	return (
		<>
			{mutation.isPending && <Spinner />}

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
							name="name"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Branch Name</FormLabel>
									<FormControl>
										<Input {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					</div>
					<Button className="ml-auto block" type="submit">
						Add Branch
					</Button>
				</form>
			</Form>
		</>
	);
}

export default BranchForm;
