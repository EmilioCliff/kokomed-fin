import {
	Card,
	CardHeader,
	CardTitle,
	CardDescription,
	CardContent,
} from '@/components/ui/card';
import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Eye, EyeOff } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { toast } from 'react-toastify';
import { useParams } from 'react-router';
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormMessage,
} from '@/components/ui/form';
import resetPasswordService from '@/services/resetPassword';
import { ResetPassowordFormType, resetPasswordFormSchema } from './schema';

function ResetPassword() {
	const [showPassword, setShowPassword] = useState(false);
	const [showConfirmPassword, setShowConfirmPassword] = useState(false);

	const { token } = useParams();

	const mutation = useMutation({
		mutationFn: resetPasswordService,
	});

	function onSubmit(values: ResetPassowordFormType) {
		mutation.mutate(values, {
			onSuccess: () => {
				toast.success('Credentials Updated');
			},
			onError: (error: any) => {
				toast.error(error.message);
			},
			onSettled: () => mutation.reset(),
		});
		// navigate("/login")
	}

	function onError(errors: any) {
		console.log('Errors: ', errors);
	}

	const form = useForm<ResetPassowordFormType>({
		resolver: zodResolver(resetPasswordFormSchema),
		defaultValues: {
			token: token,
			newPassword: '',
			confirmPassword: '',
		},
	});

	return (
		<>
			{/* {mutation.isPending && <Spinner />} */}
			<div className="grid grid-cols-1 min-h-full sm:grid-cols-1 md:grid-cols-2">
				<div className="bg-slate-500 bg-custom-bg bg-cover bg-center somen"></div>
				<div className="w-full min-h-full flex justify-center items-center">
					<Card className="mx-auto max-w-4xl border-none shadow-none flex-col text-center">
						<CardHeader className="space-y-1 gap-4">
							<CardTitle className="text-2xl font-bold">
								Reset Password
							</CardTitle>
							<CardDescription className="">
								Welcome to Kokomed Finance System
							</CardDescription>
						</CardHeader>
						<CardContent>
							<Form {...form}>
								<div>
									<form
										onSubmit={form.handleSubmit(
											onSubmit,
											onError,
										)}
										className="space-y-8"
									>
										<FormField
											control={form.control}
											name="newPassword"
											render={({ field }) => (
												<FormItem>
													<p className="text-start">
														New Password
													</p>
													<FormControl>
														<div className="space-y-2 relative">
															<Input
																placeholder="Password"
																type={
																	showPassword
																		? 'text'
																		: 'password'
																}
																{...field}
															/>
															<button
																type="button"
																className="absolute right-3 top-1/3 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
																onClick={() =>
																	setShowPassword(
																		(
																			prevState,
																		) =>
																			!prevState,
																	)
																}
																aria-label={
																	showPassword
																		? 'Hide password'
																		: 'Show password'
																}
															>
																{showPassword ? (
																	<EyeOff className="h-5 w-5" />
																) : (
																	<Eye className="h-5 w-5" />
																)}
															</button>
														</div>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
										<FormField
											control={form.control}
											name="confirmPassword"
											render={({ field }) => (
												<FormItem>
													<p className="text-start">
														Confirm Password
													</p>
													<FormControl>
														<div className="space-y-2 relative">
															<Input
																placeholder="Password"
																type={
																	showConfirmPassword
																		? 'text'
																		: 'password'
																}
																{...field}
															/>
															<button
																type="button"
																className="absolute right-3 top-1/3 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
																onClick={() =>
																	setShowConfirmPassword(
																		(
																			prevState,
																		) =>
																			!prevState,
																	)
																}
																aria-label={
																	showConfirmPassword
																		? 'Hide password'
																		: 'Show password'
																}
															>
																{showConfirmPassword ? (
																	<EyeOff className="h-5 w-5" />
																) : (
																	<Eye className="h-5 w-5" />
																)}
															</button>
														</div>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<Button
											className=" ml-auto mr-auto"
											type="submit"
										>
											Submit
										</Button>
									</form>
								</div>
							</Form>
						</CardContent>
					</Card>
				</div>
			</div>
		</>
	);
}

export default ResetPassword;
