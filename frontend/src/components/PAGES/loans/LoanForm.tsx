import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, Controller } from "react-hook-form";
import VirtualizeddSelect from "../../UI/VisualizedSelect";
import { CalendarIcon } from "lucide-react";
import { format } from "date-fns";
import { cn } from "@/lib/utils";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const loanFormSchema = z.object({
	productId: z.number().gte(0, { message: "Select valid product" }),
	clientId: z.number().gte(0, { message: "Select valid client" }),
	loanOfficerId: z.number().gte(0, { message: "Select valid loan officer" }),
	loanPurpose: z.string(),
	// approvedBy: z.number(), add approved by when form is submitted
	// disburseBy: z.number().optional(), add disbursed by when form is submitted
	disburse: z.boolean(),
	disburseOn: z
		.string()
		.optional()
		.refine((dateString) => !dateString || !isNaN(Date.parse(dateString)), {
			message: "Invalid date string!",
		}),
	noOfInstallments: z.coerce.number().gt(0),
	installmentsPeriod: z.coerce.number().gt(0),
	processingFee: z.coerce.number().gt(0),
	processingFeePaid: z.boolean(),
	dob: z.string().date().optional(),
});

const products = Array.from({ length: 200 }, (_, i) => ({
	id: i,
	name: `Product ${i + 1}`,
}));

const clients = Array.from({ length: 200 }, (_, i) => ({
	id: i,
	name: `Client ${i + 1}`,
}));

const loanOfficers = Array.from({ length: 200 }, (_, i) => ({
	id: i,
	name: `Loan Officer ${i + 1}`,
}));

export default function LoanForm() {
	const form = useForm<z.infer<typeof loanFormSchema>>({
		resolver: zodResolver(loanFormSchema),
		defaultValues: {
			productId: 0,
			clientId: 0,
			loanOfficerId: 0,
			loanPurpose: "",
			// approvedBy: 0,
			disburse: false,
			// disburseBy: 0,
			disburseOn: "",
			noOfInstallments: 4,
			installmentsPeriod: 30,
			processingFee: 0,
			processingFeePaid: false,
		},
	});

	function onSubmit(values: z.infer<typeof loanFormSchema>) {
		console.log(values);
	}

	function onError(errors: any) {
		console.log(errors);
	}

	const isAdmin = true;

	return (
		<Form {...form}>
			<form
				onSubmit={form.handleSubmit(onSubmit, onError)}
				className='space-y-8'
			>
				<div className='grid grid-cols-1 md:grid-cols-2 gap-4 py-4'>
					<FormField
						control={form.control}
						name='productId'
						render={({ field }) => (
							<FormItem>
								<FormLabel>Product</FormLabel>
								<FormControl>
									<VirtualizeddSelect
										options={products}
										placeholder='Select a product'
										value={field.value}
										onChange={(id) => field.onChange(id)}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name='clientId'
						render={({ field }) => (
							<FormItem>
								<FormLabel>Client</FormLabel>
								<FormControl>
									<VirtualizeddSelect
										options={clients}
										placeholder='Select a client'
										value={field.value}
										onChange={(id) => field.onChange(id)}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name='loanOfficerId'
						render={({ field }) => (
							<FormItem>
								<FormLabel>Loan Officer</FormLabel>
								<FormControl>
									<VirtualizeddSelect
										options={loanOfficers}
										placeholder='Select a loan officer'
										value={field.value}
										onChange={(id) => field.onChange(id)}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name='loanPurpose'
						render={({ field }) => (
							<FormItem>
								<FormLabel>Loan Purpose</FormLabel>
								<FormControl>
									<Input placeholder='Loan Purpose' {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name='noOfInstallments'
						render={({ field }) => (
							<FormItem>
								<FormLabel>No of Installments</FormLabel>
								<FormControl>
									<Input type='number' {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name='installmentsPeriod'
						render={({ field }) => (
							<FormItem>
								<FormLabel>Installment Period(days)</FormLabel>
								<FormControl>
									<Input type='number' {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name='processingFee'
						render={({ field }) => (
							<FormItem>
								<FormLabel>Processing Fee</FormLabel>
								<FormControl>
									<Input type='number' {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<div>
						<p className='mb-4'>Processing Fee Paid</p>
						<Controller
							control={form.control}
							name='processingFeePaid'
							render={({ field }) => (
								<div className='flex gap-x-8'>
									<label>
										<input
											type='radio'
											onBlur={field.onBlur}
											onChange={() => field.onChange(true)}
											checked={field.value === true}
											className='mr-2'
										/>
										Yes
									</label>
									<label>
										<input
											type='radio'
											onBlur={field.onBlur}
											onChange={() => field.onChange(false)}
											checked={field.value === false}
											className='mr-2'
										/>
										No
									</label>
								</div>
							)}
						/>
					</div>
				</div>
				<FormField
					control={form.control}
					name='dob'
					render={({ field }) => (
						<FormItem className='flex flex-col'>
							<FormLabel>Date of birth</FormLabel>
							<Popover>
								<PopoverTrigger asChild>
									<FormControl>
										<Button
											variant={"outline"}
											className={cn(
												"w-[240px] pl-3 text-left font-normal",
												!field.value && "text-muted-foreground"
											)}
										>
											{field.value ? (
												format(field.value, "PPP")
											) : (
												<span>Pick a date</span>
											)}
											<CalendarIcon className='ml-auto h-4 w-4 opacity-50' />
										</Button>
									</FormControl>
								</PopoverTrigger>
								<PopoverContent className='w-auto p-0' align='start'>
									<Calendar
										mode='single'
										selected={field.value ? new Date(field.value) : undefined}
										onSelect={(date) =>
											field.onChange(format(date!, "yyyy-MM-dd"))
										}
										disabled={(date) =>
											date > new Date() || date < new Date("1900-01-01")
										}
										initialFocus
									/>
								</PopoverContent>
							</Popover>
							<FormDescription>
								Your date of birth is used to calculate your age.
							</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>
				<Button className='ml-auto block' type='submit'>
					Add Loan
				</Button>
			</form>
		</Form>
	);
}
