import { useQueryClient, useMutation } from '@tanstack/react-query';
import Spinner from '@/components/UI/Spinner';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import { PaymentFormType, paymentFormSchema } from './schema';
import addPayment from '@/services/addPayment';
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

function PaymentForm({ onFormOpen }: LoanFormProps) {
	const form = useForm<PaymentFormType>({
		resolver: zodResolver(paymentFormSchema),
		defaultValues: {
			TransAmount: 0,
			TransID: '',
			BillRefNumber: '',
			MSISDN: '',
			FirstName: '',
			App: 'INTERNAL',
		},
	});
	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: addPayment,
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: ['payments'] });
			toast.success('Payment Added Successful');
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
		onSettled: () => onFormOpen(false),
	});

	function onSubmit(values: PaymentFormType) {
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
							name="TransAmount"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Transaction Amount</FormLabel>
									<FormControl>
										<Input type="number" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="TransID"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Transaction ID</FormLabel>
									<FormControl>
										<Input
											placeholder="TAV39ELQCV"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="BillRefNumber"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Account Number</FormLabel>
									<FormControl>
										<Input
											placeholder="0712345678"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="MSISDN"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Phone Number</FormLabel>
									<FormControl>
										<Input
											placeholder="0712345678"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="FirstName"
							render={({ field }) => (
								<FormItem>
									<FormLabel>First Name</FormLabel>
									<FormControl>
										<Input placeholder="John" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="App"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Transaction Source</FormLabel>
									<FormControl>
										<Input
											readOnly
											placeholder="INTERNAL"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					</div>
					<Button className="ml-auto block" type="submit">
						Add Payment
					</Button>
				</form>
			</Form>
		</>
	);
}

export default PaymentForm;
