import { useForm } from 'react-hook-form';
import { editPaymentFormSchema, EditPaymentFormType, Payment } from './schema';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import { editPayment, simulateEditPayment } from '@/services/editPayment';
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormDescription,
	FormMessage,
} from '@/components/ui/form';
import {
	AlertTriangle,
	ArrowLeft,
	ArrowLeftRight,
	CalendarIcon,
	CheckCircle,
	DollarSign,
	Lock,
	Navigation,
	RefreshCw,
	Trash2,
	TrendingDown,
	Zap,
} from 'lucide-react';
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
import Spinner from '@/components/UI/Spinner';
import { useAuth } from '@/hooks/useAuth';
import VirtualizeddSelect from '@/components/UI/VisualizedSelect';
import getFormData from '@/services/getFormData';
import { Textarea } from '@/components/ui/textarea';
import { useEffect, useState } from 'react';
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { useNavigate, useParams } from 'react-router';
import { role, simulationResult } from '@/lib/types';
import getNonPosted from '@/services/getNonPosted';
import EditPaymentFormSkeleton from '@/components/UI/EditPaymentFormSkeleton';

function EditPaymentForm() {
	const [step, setStep] = useState<'edit' | 'confirm'>('edit');
	const [originalPayment, setOriginalPayment] =
		useState<EditPaymentFormType | null>(null);
	const [editedPayment, setEditedPayment] =
		useState<EditPaymentFormType | null>(null);
	const [actions, setActions] = useState<simulationResult | null>(null);
	const [isBlocked, setIsBlocked] = useState(false);
	const [lacksData, setLacksData] = useState(false);

	const navigate = useNavigate();
	const { id } = useParams();
	const { decoded } = useAuth();

	const form = useForm<EditPaymentFormType>({
		resolver: zodResolver(editPaymentFormSchema),
		defaultValues: {
			id: 0,
			transactionSource: 'MPESA',
			transactionId: '',
			accountNumber: '',
			phoneNumber: '',
			payingName: '',
			amount: 0,
			assignedBy: '',
			assignedTo: 0,
			description: '',
			paidDate: format(new Date(), 'yyyy-MM-dd'),
		},
	});

	const nonPostedQuery = useQuery({
		queryKey: ['nonposted', id],
		queryFn: () => getNonPosted(Number(id)),
		staleTime: 5 * 1000,
	});

	const { data } = useQuery({
		queryKey: ['clients/form'],
		queryFn: () => getFormData(false, true, false, false, false),
		staleTime: 5 * 1000,
	});

	const mutation = useMutation({
		mutationFn: simulateEditPayment,
		onSuccess: (data) => {
			setActions(data.data);
			setIsBlocked(
				data.data.actions.some(
					(action) => action.actionType === 'loan_status_blocked',
				),
			);
			setLacksData(
				data.data.actions.some(
					(action) => action.actionType === 'not_enough_data',
				),
			);
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
	});

	const editMutation = useMutation({
		mutationFn: editPayment,
		onSuccess: () => {
			toast.success('Payment Updated Successful');
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
		onSettled: () => {},
	});

	useEffect(() => {
		if (!nonPostedQuery.data?.data || !decoded?.email) return;

		if (decoded.role !== role.ADMIN) {
			navigate('/payments/overview');
			toast.error('not authorized to edit payment');
		}

		const data = nonPostedQuery.data.data;

		const paymentData: EditPaymentFormType = {
			id: data.id,
			transactionSource: data.transactionSource,
			transactionId: data.transactionNumber,
			accountNumber: data.accountNumber,
			phoneNumber: data.phoneNumber,
			payingName: data.payingName,
			amount: data.amount,
			assignedBy: decoded.email,
			assignedTo: data.assignedTo.id ? data.assignedTo.id : 0,
			description: '',
			paidDate: format(data.paidDate, 'yyyy-MM-dd'),
		};

		setOriginalPayment(paymentData);
		setEditedPayment(paymentData);

		form.reset(paymentData);
	}, [nonPostedQuery.data, decoded]);

	if (!originalPayment || !editedPayment) {
		return <EditPaymentFormSkeleton />;
	}

	const hasChanges =
		JSON.stringify(originalPayment) !== JSON.stringify(editedPayment);
	const hasFinancialImpact = originalPayment.amount !== editedPayment.amount;
	const getChangeValue = (field: keyof EditPaymentFormType) => {
		const original = originalPayment[field];
		const edited = editedPayment[field];

		if (original === edited) return null;

		return { original, edited };
	};

	const handleFormSubmit = (data: EditPaymentFormType) => {
		mutation.mutate(data);
		setEditedPayment(data);
		setStep('confirm');
	};

	const handleConfirm = async () => {
		editMutation.mutate(editedPayment);
		navigate('/payments/overview');
	};

	const handleBack = () => {
		setStep('edit');
	};

	if (mutation.error) {
		return (
			<div className="container mx-auto p-6">
				<Card className="border-destructive">
					<CardHeader>
						<CardTitle className="text-destructive">
							Error Editing Payment
						</CardTitle>
					</CardHeader>
					<CardContent>
						<p>{mutation.error.message}</p>
					</CardContent>
					<CardFooter>
						<Button
							variant={'destructive'}
							onClick={() => {
								setStep('edit');
								mutation.reset();
							}}
						>
							Back
						</Button>
					</CardFooter>
				</Card>
			</div>
		);
	}

	return (
		<>
			{mutation.isPending && <Spinner />}
			<div className="p-4">
				<h1 className="text-3xl font-bold">Edit Payment</h1>
				<p className="text-muted-foreground mb-2">
					{step === 'edit'
						? 'Modify payment details below'
						: 'Review changes and confirm the actions that will be taken'}
				</p>
				{step === 'edit' ? (
					<Card>
						<Form {...form}>
							<form
								onSubmit={form.handleSubmit(handleFormSubmit)}
								className="space-y-8"
							>
								<CardContent>
									<div className="grid grid-cols-1 md:grid-cols-2 gap-4 py-4">
										<FormField
											control={form.control}
											name="amount"
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Transaction Amount
													</FormLabel>
													<FormControl>
														<Input
															type="number"
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
										<FormField
											control={form.control}
											name="transactionId"
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Transaction ID
													</FormLabel>
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
											name="accountNumber"
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Account Number
													</FormLabel>
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
											name="phoneNumber"
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Phone Number
													</FormLabel>
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
											name="payingName"
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Paying Name
													</FormLabel>
													<FormControl>
														<Input
															placeholder="John"
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name="assignedTo"
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Assigned To
													</FormLabel>
													<FormControl>
														{data?.client && (
															<VirtualizeddSelect
																options={
																	data.client
																}
																placeholder="Select Client"
																value={
																	field.value
																}
																onChange={(
																	id,
																) =>
																	field.onChange(
																		id,
																	)
																}
																disabled={true}
															/>
														)}
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
										<FormField
											control={form.control}
											name="description"
											render={({ field }) => (
												<FormItem className="col-span-2">
													<FormLabel>
														Description
													</FormLabel>
													<FormControl>
														<Textarea
															placeholder="Why are you updating this loan"
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
										<FormField
											control={form.control}
											name="assignedBy"
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Updated By
													</FormLabel>
													<FormControl>
														<Input
															readOnly
															placeholder={
																field.value
															}
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
										<FormField
											control={form.control}
											name="transactionSource"
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Transaction Source
													</FormLabel>
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
										<FormField
											control={form.control}
											name="paidDate"
											render={({ field }) => (
												<FormItem className="flex flex-col mt-3">
													<FormLabel>
														Paid On
													</FormLabel>
													<Popover>
														<PopoverTrigger asChild>
															<FormControl>
																<Button
																	variant={
																		'outline'
																	}
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
																		<span>
																			Pick
																			a
																			date
																		</span>
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
																		? new Date(
																				field.value,
																		  )
																		: undefined
																}
																onSelect={(
																	date,
																) =>
																	field.onChange(
																		format(
																			date!,
																			'yyyy-MM-dd',
																		),
																	)
																}
																disabled={(
																	date,
																) =>
																	date >
																		new Date() ||
																	date <
																		new Date(
																			'1900-01-01',
																		)
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
									</div>
								</CardContent>
								<CardFooter>
									<Button
										onClick={() =>
											navigate('/payments/overview')
										}
										type="button"
										variant="outline"
									>
										Cancel
									</Button>
									<Button
										className="ml-auto block"
										type="submit"
									>
										Review Changes
									</Button>
								</CardFooter>
							</form>
						</Form>
					</Card>
				) : (
					<div className="space-y-6">
						{/* Changes Summary */}
						<Card>
							<CardHeader>
								<CardTitle>Payment Changes Summary</CardTitle>
								<CardDescription>
									Review the changes you've made to this
									payment
								</CardDescription>
							</CardHeader>
							<CardContent className="space-y-4">
								{!hasChanges ? (
									<p className="text-muted-foreground">
										No changes detected
									</p>
								) : (
									<div className="space-y-3">
										{Object.keys(originalPayment).map(
											(key) => {
												const field =
													key as keyof EditPaymentFormType;
												const change =
													getChangeValue(field);

												if (!change || field === 'id')
													return null;

												return (
													<div
														key={field}
														className="flex justify-between items-center py-2 border-b"
													>
														<span className="font-medium capitalize">
															{field
																.replace(
																	/([A-Z])/g,
																	' $1',
																)
																.trim()}
														</span>
														<div className="text-right">
															<div className="text-sm text-muted-foreground line-through">
																{String(
																	change.original,
																)}
															</div>
															<div className="font-medium">
																{String(
																	change.edited,
																)}
															</div>
														</div>
													</div>
												);
											},
										)}
									</div>
								)}
							</CardContent>
						</Card>

						{/* Actions that will be taken */}
						{/* <Card>
							<CardHeader>
								<CardTitle className="flex items-center gap-2">
									<AlertTriangle className="h-5 w-5" />
									Actions to be Performed
								</CardTitle>
								<CardDescription>
									The following actions will be automatically
									executed when you confirm
								</CardDescription>
							</CardHeader>
							<CardContent>
								{hasFinancialImpact && (
									<Alert className="mb-4">
										<AlertTriangle className="h-4 w-4" />
										<AlertTitle>
											Financial Impact Detected
										</AlertTitle>
										<AlertDescription>
											This payment change will affect loan
											balances, installments, and
											overpayments. Please review the
											actions carefully.
										</AlertDescription>
									</Alert>
								)}

								<div className="space-y-3">
									{actions.map((action, index) => (
										<div
											key={index}
											className="flex items-start gap-3 p-3 rounded-lg border"
										>
											<div
												className={`mt-0.5 ${getSeverityColor(
													action.severity,
												)}`}
											>
												{action.icon}
											</div>
											<div className="flex-1">
												<p className="font-medium">
													{action.description}
												</p>
												{action.amount && (
													<p className="text-sm text-muted-foreground">
														Amount: KES{' '}
														{action.amount.toLocaleString()}
													</p>
												)}
											</div>
											<Badge
												variant={
													action.severity ===
													'warning'
														? 'destructive'
														: 'default'
												}
											>
												{action.severity}
											</Badge>
										</div>
									))}
								</div>
							</CardContent>
						</Card> */}
						<Card>
							<CardHeader>
								<CardTitle className="flex items-center gap-2">
									<AlertTriangle className="h-5 w-5" />
									Actions to be Performed
								</CardTitle>
								<CardDescription>
									The following actions will be automatically
									executed when you confirm
								</CardDescription>
							</CardHeader>
							<CardContent>
								{isBlocked && (
									<Alert
										variant="destructive"
										className="mb-4"
									>
										<AlertTriangle className="h-4 w-4" />
										<AlertTitle>Action Blocked</AlertTitle>
										<AlertDescription>
											Changes cannot be confirmed because
											another loan is currently active.
											You cannot have multiple active
											loans at the same time.
										</AlertDescription>
									</Alert>
								)}
								{lacksData && (
									<Alert
										variant="destructive"
										className="mb-4"
									>
										<AlertTriangle className="h-4 w-4" />
										<AlertTitle>Action Blocked</AlertTitle>
										<AlertDescription>
											Changes cannot be confirmed because
											the payment lacks enough data. Seek
											support for this action
										</AlertDescription>
									</Alert>
								)}
								{hasFinancialImpact && (
									<Alert className="mb-4">
										<AlertTriangle className="h-4 w-4" />
										<AlertTitle>
											Financial Impact Detected
										</AlertTitle>
										<AlertDescription>
											This payment change will affect loan
											balances, installments, and
											overpayments. Please review the
											actions carefully.
										</AlertDescription>
									</Alert>
								)}

								<div className="space-y-3">
									{actions &&
										actions.actions.map((action, index) => (
											<div
												key={index}
												className="flex items-start gap-3 p-3 rounded-lg border"
											>
												<div
													className={`mt-0.5 ${getSeverityColor(
														action.severity,
													)}`}
												>
													{getActionIcon(
														action.actionType,
													)}
												</div>
												<div className="flex-1">
													<p className="font-medium">
														{action.description}
													</p>
													{action.amount !== 0 && (
														<p className="text-sm text-muted-foreground">
															Amount: KES{' '}
															{action.amount.toLocaleString()}
														</p>
													)}
												</div>
												<Badge
													variant={
														action.severity ===
														'warning'
															? 'destructive'
															: 'default'
													}
												>
													{action.severity}
												</Badge>
											</div>
										))}
								</div>
							</CardContent>
						</Card>

						{/* Confirmation Actions */}
						<div className="flex justify-between">
							<Button
								variant="outline"
								onClick={handleBack}
								className="flex items-center gap-2"
							>
								<ArrowLeft className="h-4 w-4" />
								Back to Edit
							</Button>
							<Button
								onClick={handleConfirm}
								disabled={!hasChanges || isBlocked}
								className="flex items-center gap-2"
							>
								<CheckCircle className="h-4 w-4" />
								Confirm Changes
							</Button>
						</div>
					</div>
				)}
			</div>
		</>
	);
}

export default EditPaymentForm;

const getActionIcon = (type: string) => {
	switch (type) {
		case 'update_payment':
			return <CheckCircle className="h-4 w-4" />;
		case 'revert_installments':
			return <RefreshCw className="h-4 w-4" />;
		case 'reduce_overpayment':
			return <TrendingDown className="h-4 w-4" />;
		case 'allocate_additional':
			return <DollarSign className="h-4 w-4" />;
		case 'reactivate_loan':
			return <RefreshCw className="h-4 w-4" />;
		case 'loan_status_blocked':
			return <Lock className="h-4 w-4" />;
		case 'not_enough_data':
			return <Lock className="h-4 w-4" />;
		case 'loan_status_change':
			return <ArrowLeftRight className="h-4 w-4" />;
		case 'process_loan_payment':
			return <Zap className="h-4 w-4" />;
		case 'delete_non_posted':
			return <Trash2 className="h-4 w-4" />;
		default:
			return <CheckCircle className="h-4 w-4" />;
	}
};

const getSeverityColor = (severity: string) => {
	switch (severity) {
		case 'info':
			return 'text-blue-500';
		case 'warning':
			return 'text-yellow-600';
		case 'success':
			return 'text-green-600';
		case 'danger':
			return 'text-red-600';
		default:
			return 'text-gray-500';
	}
};
