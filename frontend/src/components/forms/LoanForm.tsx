import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import VirtualizeddSelect from "../VisualizeddSelect";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form";

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

	const isAdmin = true;

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className='space-y-8'>
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
					{isAdmin && (
						<FormField
							control={form.control}
							name='processingFeePaid'
							render={({ field }) => (
								<FormItem className='space-y-3'>
									<FormLabel>Processing Fee Paid</FormLabel>
									<FormControl>
										<RadioGroup
											onValueChange={(value) => field.onChange(value === "yes")}
											value={field.value ? "yes" : "no"}
											className='flex space-x-4'
										>
											<FormItem className='flex items-center space-x-3 space-y-0'>
												<FormControl>
													<RadioGroupItem
														className='peer'
														value='yes'
														id='processing-fee-yes'
													/>
												</FormControl>
												<FormLabel
													htmlFor='processing-fee-yes'
													className='font-normal peer-checked:text-blue-600'
												>
													Yes
												</FormLabel>
											</FormItem>
											<FormItem className='flex items-center space-x-3 space-y-0'>
												<FormControl>
													<RadioGroupItem
														className='peer'
														value='no'
														id='processing-fee-no'
													/>
												</FormControl>
												<FormLabel
													htmlFor='processing-fee-no'
													className='font-normal peer-checked:text-blue-600'
												>
													No
												</FormLabel>
											</FormItem>
										</RadioGroup>
									</FormControl>
								</FormItem>
							)}
						/>
					)}
				</div>
				<Button className='ml-auto block' type='submit'>
					Add Loan
				</Button>
			</form>
		</Form>
	);
}
