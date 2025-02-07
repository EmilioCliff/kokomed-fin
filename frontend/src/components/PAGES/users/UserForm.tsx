import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import getFormData from '@/services/getFormData';
import { UserFormType, userFormSchema } from './schema';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import VirtualizeddSelect from '../../UI/VisualizedSelect';
import addUser from '@/services/addUser';
import { toast } from 'react-toastify';
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
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { role } from '@/lib/types';

interface UserFormProps {
	onFormOpen: (isOpen: boolean) => void;
}

function UserForm({ onFormOpen }: UserFormProps) {
	const { isLoading, data, error } = useQuery({
		queryKey: ['loans/form'],
		queryFn: () => getFormData(false, false, false, true),
		staleTime: 5 * 1000,
	});

	const form = useForm<UserFormType>({
		resolver: zodResolver(userFormSchema),
		defaultValues: {
			firstName: '',
			lastName: '',
			phoneNumber: '',
			email: '',
			branchId: 0,
			role: 'AGENT',
		},
	});
	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: addUser,
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: ['users'] });
			toast.success('User Added Successful');
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
		onSettled: () => onFormOpen(false),
	});

	function onSubmit(values: UserFormType) {
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
							name="firstName"
							render={({ field }) => (
								<FormItem>
									<FormLabel>First Name</FormLabel>
									<FormControl>
										<Input
											placeholder="First Name"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="lastName"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Last Name</FormLabel>
									<FormControl>
										<Input
											placeholder="Last Name"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="email"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Email</FormLabel>
									<FormControl>
										<Input
											placeholder="bob@example.com"
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
							name="branchId"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Branch</FormLabel>
									<FormControl>
										{data?.branch && (
											<VirtualizeddSelect
												options={data.branch}
												placeholder="Select User Branch"
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
							name="role"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Role</FormLabel>
									<FormControl>
										<Select
											value={field.value}
											onValueChange={(value) => {
												field.onChange(value);
											}}
										>
											<SelectTrigger className="w-[180px]">
												<SelectValue placeholder="Select Role" />
											</SelectTrigger>
											<SelectContent>
												<SelectGroup>
													<SelectItem
														value={role.ADMIN}
													>
														ADMIN
													</SelectItem>
													<SelectItem
														value={role.AGENT}
													>
														AGENT
													</SelectItem>
												</SelectGroup>
											</SelectContent>
										</Select>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					</div>
					<Button className="ml-auto block" type="submit">
						Add User
					</Button>
				</form>
			</Form>
		</>
	);
}

export default UserForm;
