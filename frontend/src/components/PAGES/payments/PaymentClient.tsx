import { useState, useEffect } from 'react';
import { Card } from '@/components/ui/card';
import VirtualizeddSelect from '@/components/UI/VisualizedSelect';
import { useQuery, useMutation, keepPreviousData } from '@tanstack/react-query';
import getFormData from '@/services/getFormData';
import { Input } from '@/components/ui/input';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
	getClientNonPosted,
	getClientPayment,
} from '@/services/getClientNonPosted';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';
import { Edit } from 'lucide-react';
import { toast } from 'react-toastify';
import { useTable } from '@/hooks/useTable';
import { useDebounce } from '@/hooks/useDebounce';
import EditClientForm from './EditClientForm';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import PaymentsTab from './PaymentsTab';
import { getClientLoans } from '@/services/getClientLoans';
import LoansTab from './LoansTab';

function PaymentClient() {
	const [searchType, setSearchType] = useState('id');
	const [clientId, setClientId] = useState(0);
	const [phoneNumber, setPhoneNumber] = useState('');
	const { pageIndex, pageSize, filter, updateTableContext, resetTableState } =
		useTable();
	const [searched, setSearched] = useState(false);
	const [editForm, setEditForm] = useState(false);
	const [currentTab, setCurrentTab] = useState('loans');

	const debouced = useDebounce({ value: phoneNumber, delay: 2000 });

	const paymentsQuery = useQuery({
		queryKey: ['payments/client', pageIndex, pageSize, clientId],
		queryFn: () =>
			getClientPayment(clientId, debouced, pageIndex, pageSize),
		staleTime: 5 * 1000,
		placeholderData: keepPreviousData,
		enabled:
			searched &&
			(!!phoneNumber || !!clientId) &&
			currentTab === 'payments',
	});

	const loansQuery = useQuery({
		queryKey: ['loans/client', pageIndex, pageSize, clientId, filter],
		queryFn: () => getClientLoans(clientId, pageIndex, pageSize, filter),
		staleTime: 5 * 1000,
		placeholderData: keepPreviousData,
		enabled:
			searched && (!!phoneNumber || !!clientId) && currentTab === 'loans',
	});

	useEffect(() => {
		if (currentTab === 'payments') {
			if (paymentsQuery.data) {
				if (paymentsQuery.data.metadata) {
					updateTableContext(paymentsQuery.data.metadata);
				}
			}
		} else {
			if (loansQuery.data) {
				if (loansQuery.data.metadata) {
					updateTableContext(loansQuery.data.metadata);
				}
			}
		}
	}, [paymentsQuery.data?.data, loansQuery.data?.data, currentTab]);

	const { data } = useQuery({
		queryKey: ['payments/form'],
		queryFn: () => getFormData(false, true, false, false, false),
		staleTime: 5 * 1000,
	});

	const mutation = useMutation({
		mutationFn: getClientNonPosted,
		onError: (error: any) => {
			setSearched(false);
			toast.error(error.message);
		},
	});

	const handleSubmit = () => {
		setSearched(true);
		resetTableState();
		mutation.mutate({ id: clientId, phoneNumber: phoneNumber });
	};

	return (
		<div className="px-4">
			<Card className="p-4 mb-6">
				<h1 className="text-2xl font-bold mb-6">
					Client Payment Management
				</h1>
				<form
					className="flex flex-col sm:flex-row gap-4"
					onSubmit={(e) => {
						e.preventDefault();
						handleSubmit();
					}}
				>
					<div className="flex-1">
						{searchType === 'id' ? (
							<>
								{data?.client && (
									<VirtualizeddSelect
										options={data.client}
										placeholder="Search Client"
										value={clientId}
										onChange={(id) => {
											setSearched(false);
											setClientId(id);
										}}
									/>
								)}
							</>
						) : (
							<Input
								placeholder="Enter Phone Number"
								type="text"
								value={phoneNumber}
								onChange={(e) => {
									setSearched(false);
									setPhoneNumber(e.target.value);
								}}
								className="w-full p-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
							/>
						)}
					</div>
					<div className="flex-none w-full sm:w-auto">
						<Select
							value={searchType}
							onValueChange={(value: string) => {
								setClientId(0);
								setPhoneNumber('');
								setSearched(false);
								setSearchType(value);
							}}
						>
							<SelectTrigger className="w-[180px]">
								<SelectValue placeholder={searchType} />
							</SelectTrigger>
							<SelectContent>
								<SelectGroup>
									<SelectItem value="id">
										By Client
									</SelectItem>
									<SelectItem value="phoneNumber">
										By PhoneNumber
									</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
					</div>
					<div className="flex-none">
						<Button type="submit">Search</Button>
					</div>
				</form>
			</Card>
			{mutation.data && searched && (
				<>
					<Card className="p-4 mb-6">
						{mutation.data.data.clientDetails.id === 0 ? (
							<div className="text-center py-6">
								<p className="text-gray-500">
									No client data found for this{' '}
									{searchType === 'id'
										? 'client.'
										: 'number.'}
								</p>
							</div>
						) : (
							<>
								<div className="flex justify-between">
									<h2 className="text-xl font-semibold mb-4">
										Client Details
									</h2>
									<Dialog
										open={editForm}
										onOpenChange={setEditForm}
									>
										<DialogTrigger asChild>
											<Button variant="outline" size="sm">
												<Edit className="mr-2 h-4 w-4" />
												Edit
											</Button>
										</DialogTrigger>
										<DialogContent>
											<DialogHeader>
												<DialogTitle>
													Edit Customer
												</DialogTitle>
												<DialogDescription>
													Update the customer's
													details below.
												</DialogDescription>
											</DialogHeader>

											{mutation.data?.data ? (
												<EditClientForm
													onFormOpen={setEditForm}
													clientData={
														mutation.data.data
															.clientDetails
													}
												/>
											) : (
												<p className="text-sm text-muted-foreground">
													Loading client data...
												</p>
											)}
										</DialogContent>
									</Dialog>
								</div>
								<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
									<div>
										<p className="text-sm text-muted-foreground">
											Client ID
										</p>
										<p className="font-medium">
											{`CM${String(
												mutation.data.data.clientDetails
													.id,
											).padStart(3, '0')}` || 'N/A'}
										</p>
									</div>
									<div>
										<p className="text-sm text-muted-foreground">
											Name
										</p>
										<p className="font-medium">
											{mutation.data.data.clientDetails
												.fullName || 'N/A'}
										</p>
									</div>
									<div>
										<p className="text-sm text-muted-foreground">
											Phone Number
										</p>
										<p className="font-medium">
											{mutation.data.data.clientDetails
												.phoneNumber || 'N/A'}
										</p>
									</div>
									<div>
										<p className="text-sm text-muted-foreground">
											Branch
										</p>
										<p className="font-medium">
											{mutation.data.data.clientDetails
												.branchName || 'N/A'}
										</p>
									</div>
									<div>
										<p className="text-sm text-muted-foreground">
											Overpayment
										</p>
										<p className="font-medium text-green-600">
											{mutation.data.data.clientDetails
												.overpayment
												? `KES ${mutation.data.data.clientDetails.overpayment.toLocaleString()}`
												: 'KES 0.00'}
										</p>
									</div>
								</div>
							</>
						)}
					</Card>
					<Card className="p-4 mb-6">
						{mutation.data.data.loanShort.id === 0 ? (
							<div className="text-center py-6">
								<p className="text-gray-500">
									No active loans found for this{' '}
									{searchType === 'id'
										? 'client.'
										: 'number.'}
								</p>
							</div>
						) : (
							<>
								<div className="flex flex-wrap justify-between items-center mb-4">
									<h2 className="text-xl text-center font-semibold">
										Active Loan
									</h2>
									<div className="mt-2 lg:mt-0">
										<Badge>
											{`LN${String(
												mutation.data.data.loanShort.id,
											).padStart(3, '0')}` || 'N/A'}
										</Badge>
									</div>
								</div>
								<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
									<div>
										<p className="text-sm text-gray-500">
											Loan Amount
										</p>
										<p className="font-medium">
											KES{' '}
											{mutation.data.data.loanShort.loanAmount.toLocaleString()}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Repayment Amount
										</p>
										<p className="font-medium">
											KES{' '}
											{mutation.data.data.loanShort.repayAmount.toLocaleString()}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Amount Paid
										</p>
										<p className="font-medium">
											KES{' '}
											{mutation.data.data.loanShort.paidAmount.toLocaleString()}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Unpaid Amount
										</p>
										<p className="font-medium">
											KES{' '}
											{(
												mutation.data.data.loanShort
													.repayAmount -
												mutation.data.data.loanShort
													.paidAmount
											).toLocaleString()}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Disbursed On
										</p>
										<p className="font-medium">
											{
												mutation.data.data.loanShort
													.disbursedOn
											}
										</p>
									</div>
									<div>
										<p className="text-sm text-gray-500">
											Due Date
										</p>
										<p className="font-medium">
											{
												mutation.data.data.loanShort
													.dueDate
											}
										</p>
									</div>
								</div>
								<h3 className="text-lg font-medium mb-3">
									Installments
								</h3>
								<div className="overflow-x-auto">
									<Table>
										<TableHeader>
											<TableRow>
												<TableHead className="text-center">
													No.
												</TableHead>
												<TableHead className="text-center">
													Amount Due
												</TableHead>
												<TableHead className="text-center">
													Remaining
												</TableHead>
												<TableHead className="text-center">
													Due Date
												</TableHead>
												<TableHead className="text-center">
													Status
												</TableHead>
												<TableHead className="text-center">
													Paid Date
												</TableHead>
											</TableRow>
										</TableHeader>
										<TableBody className="text-center">
											{mutation.data.data.loanShort.installments.map(
												(installment) => (
													<TableRow
														key={installment.id}
													>
														<TableCell className="px-4 py-3 whitespace-nowrap text-sm">
															{installment.id}
														</TableCell>
														<TableCell className="px-4 py-3 whitespace-nowrap text-sm">
															KES{' '}
															{installment.amount.toLocaleString()}
														</TableCell>
														<TableCell className="px-4 py-3 whitespace-nowrap text-sm">
															KES{' '}
															{installment.remainingAmount.toLocaleString()}
														</TableCell>
														<TableCell className="px-4 py-3 whitespace-nowrap text-sm">
															{
																installment.dueDate
															}
														</TableCell>
														<TableCell className="px-4 py-3 whitespace-nowrap text-sm">
															<Badge
																variant={
																	installment.paid
																		? 'default'
																		: 'destructive'
																}
															>
																{installment.paid
																	? 'Paid'
																	: 'Pending'}
															</Badge>
														</TableCell>
														<TableCell className="px-4 py-3 whitespace-nowrap text-sm">
															{installment.paid
																? installment.paidAt
																: '-'}
														</TableCell>
													</TableRow>
												),
											)}
										</TableBody>
									</Table>
								</div>
							</>
						)}
					</Card>
					<Tabs
						onValueChange={(value: string) => {
							resetTableState();
							setCurrentTab(value);
						}}
						className="w-full"
						defaultValue="loans"
					>
						<TabsList className="grid w-full grid-cols-2">
							<TabsTrigger value="loans">Loans</TabsTrigger>
							<TabsTrigger value="payments">Payments</TabsTrigger>
						</TabsList>
						<TabsContent value="loans">
							<LoansTab
								loans={loansQuery.data?.data || []}
								searchType={searchType}
							/>
						</TabsContent>
						<TabsContent value="payments">
							<PaymentsTab
								payments={paymentsQuery.data?.data || []}
								total={mutation.data.data.totalPaid}
								searchType={searchType}
							/>
						</TabsContent>
					</Tabs>
				</>
			)}
			{searched === false && (
				<Card className="mt-10">
					<div className="p-6 rounded-lg shadow-md text-center">
						<p className="text-muted-foreground">
							Enter a client ID or phone number to view payment
							details.
						</p>
					</div>
				</Card>
			)}
		</div>
	);
}

export default PaymentClient;
