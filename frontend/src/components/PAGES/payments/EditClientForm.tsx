import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import getFormData from '@/services/getFormData';
import { editClientFormSchema, EditClientFormType } from './schema';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { toast } from 'react-toastify';
import updateClient from '@/services/updateClient';
import VirtualizeddSelect from '../../UI/VisualizedSelect';
import { CalendarIcon } from 'lucide-react';
import { format } from 'date-fns';
import { cn } from '@/lib/utils';
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form';
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import { Calendar } from '@/components/ui/calendar';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import Spinner from '@/components/UI/Spinner';
import { Client } from '../customers/schema';

interface editCustomerFormProps {
	onFormOpen: (isOpen: boolean) => void;
	clientData: Client;
}

const EditClientForm = ({ onFormOpen, clientData }: editCustomerFormProps) => {
	const { isLoading, data, error } = useQuery({
		queryKey: ['clients/form'],
		queryFn: () => getFormData(false, false, true, true, false),
		staleTime: 5 * 1000,
	});

	const form = useForm<EditClientFormType>({
		resolver: zodResolver(editClientFormSchema),
		defaultValues: {
			id: clientData.id,
			fullName: clientData.fullName,
			phoneNumber: clientData.phoneNumber,
			idNumber: clientData.idNumber === '' ? '' : clientData.idNumber,
			dob: clientData.dob === '0001-01-01' ? '' : clientData.dob,
			active: clientData.active === true ? 'true' : 'false',
			gender: clientData.gender,
			branchId: clientData.branchId,
			assignedStaffId: clientData.assignedStaff.id,
		},
	});
	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: updateClient,
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: ['clients'] });
			toast.success('Customer Edited Successful');
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
		onSettled: () => onFormOpen(false),
	});

	function onSubmit(values: EditClientFormType) {
		mutation.mutate(values);
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
							name="fullName"
							render={({ field }) => (
								<FormItem>
									<FormLabel>First Name</FormLabel>
									<FormControl>
										<Input
											placeholder="Full Name"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="phoneNumber"
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
							name="idNumber"
							render={({ field }) => (
								<FormItem>
									<FormLabel>ID Number</FormLabel>
									<FormControl>
										<Input
											placeholder="12345678"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="active"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Active</FormLabel>
									<FormControl>
										<Select
											value={field.value}
											onValueChange={(value) => {
												field.onChange(value);
											}}
										>
											<SelectTrigger className="w-[180px]">
												<SelectValue
													placeholder={field.value}
												/>
											</SelectTrigger>
											<SelectContent>
												<SelectGroup>
													<SelectItem value="true">
														ACTIVE
													</SelectItem>
													<SelectItem value="false">
														INACTIVE
													</SelectItem>
												</SelectGroup>
											</SelectContent>
										</Select>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="gender"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Gender</FormLabel>
									<FormControl>
										<Select
											value={field.value}
											onValueChange={(value) => {
												field.onChange(value);
											}}
										>
											<SelectTrigger className="w-[180px]">
												<SelectValue
													placeholder={field.value}
												/>
											</SelectTrigger>
											<SelectContent>
												<SelectGroup>
													<SelectItem value="MALE">
														MALE
													</SelectItem>
													<SelectItem value="FEMALE">
														FEMALE
													</SelectItem>
												</SelectGroup>
											</SelectContent>
										</Select>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="dob"
							render={({ field }) => (
								<FormItem className="flex flex-col">
									<FormLabel>Date of Birth</FormLabel>
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
									<FormMessage />
								</FormItem>
							)}
						/>
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
												placeholder="Select Customer Branch"
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
							name="assignedStaffId"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Assigned Agent</FormLabel>
									<FormControl>
										{data?.user && (
											<VirtualizeddSelect
												options={data.user}
												placeholder="Assign To Agent"
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
					</div>
					<Button className="ml-auto block" type="submit">
						Edit Customer
					</Button>
				</form>
			</Form>
		</>
	);
};

export default EditClientForm;
