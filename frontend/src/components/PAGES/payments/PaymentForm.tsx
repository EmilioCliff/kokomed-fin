import { useQueryClient, useMutation, useQuery } from '@tanstack/react-query';
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
	FormDescription,
	FormMessage,
} from '@/components/ui/form';
import { CalendarIcon } from 'lucide-react';
import { format } from 'date-fns';
import { cn } from '@/lib/utils';
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import { useAuth } from '@/hooks/useAuth';
import getFormData from '@/services/getFormData';
import VirtualizeddSelect from '@/components/UI/VisualizedSelect';

interface LoanFormProps {
	onFormOpen: (isOpen: boolean) => void;
}

function PaymentForm({ onFormOpen }: LoanFormProps) {
	const { decoded } = useAuth();

	const { isLoading, data, error } = useQuery({
		queryKey: ['loans/form'],
		queryFn: () => getFormData(false, true, false, false, false),
		staleTime: 5 * 1000,
	});

	const form = useForm<PaymentFormType>({
		resolver: zodResolver(paymentFormSchema),
		defaultValues: {
			TransAmount: 0,
			TransID: '',
			BillRefNumber: '',
			MSISDN: '***',
			FirstName: '',
			DatePaid: '',
			App: 'INTERNAL',
			Email: decoded?.email,
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

			{isLoading && <Spinner />}

			{error && <h1>error</h1>}

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
										{data?.client && (
											<VirtualizeddSelect
												options={data.client}
												placeholder="Select a loan officer"
												value={field.value}
												onPhoneChange={(
													phoneNumber,
												) => {
													field.onChange(phoneNumber);
												}}
											/>
										)}
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
							name="DatePaid"
							render={({ field }) => (
								<FormItem className="flex flex-col mt-3">
									<FormLabel>Paid On</FormLabel>
									<Popover>
										<PopoverTrigger asChild>
											<FormControl>
												<Button
													variant={'outline'}
													className={cn(
														'w-[240px] pl-3 text-left font-normal',
														!field.value &&
															'text-muted-foreground',
													)}
												>
													{field.value ? (
														format(
															field.value,
															'PPP',
														)
													) : (
														<span>Pick a date</span>
													)}
													<CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
												</Button>
											</FormControl>
										</PopoverTrigger>
										<PopoverContent
											className="w-auto p-0"
											align="start"
										>
											<Calendar
												mode="single"
												selected={
													field.value
														? new Date(field.value)
														: undefined
												}
												onSelect={(date) =>
													field.onChange(
														format(
															date!,
															'yyyy-MM-dd',
														),
													)
												}
												disabled={(date) =>
													date > new Date() ||
													date <
														new Date('1900-01-01')
												}
												initialFocus
											/>
										</PopoverContent>
									</Popover>
									<FormDescription>
										Defaults to today
									</FormDescription>
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
