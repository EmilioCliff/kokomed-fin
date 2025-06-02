import { Badge } from '@/components/ui/badge';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { ClientDetailsCard } from '@/components/UI/ClientDetailsCard';
import { LoanDetailsInstallmentsTable } from '@/components/UI/InstallmentsTable';
import { LoanDetailsCard } from '@/components/UI/LoanDetailsCard';
import { LoanDetailsPaymentsTable } from '@/components/UI/LoanDetailsPaymentTable';
import LoanDetailsSkeleton from '@/components/UI/LoanDetailsSkeleton';
import { LoanDetailsPaymentAllocationModal } from '@/components/UI/PaymentAllocationModal';
import { formatDate } from '@/lib/utils';
import getLoan from '@/services/getLoan';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { useQuery } from '@tanstack/react-query';
import { CalendarIcon, DollarSign } from 'lucide-react';
import { useState } from 'react';
import { useParams } from 'react-router';

function LoanDetailsPage() {
	const [selectedPayment, setSelectedPayment] = useState<any>(null);
	const [isAllocationModalOpen, setIsAllocationModalOpen] = useState(false);

	const { id } = useParams();
	const { isLoading, data, error } = useQuery({
		queryKey: ['clients/form'],
		queryFn: () => getLoan(Number(id)),
		staleTime: 5 * 1000,
	});

	const handlePaymentClick = (payment: any) => {
		setSelectedPayment(payment);
		setIsAllocationModalOpen(true);
	};

	if (isLoading || !data?.data) return <LoanDetailsSkeleton />;

	if (error) {
		return (
			<div className="container mx-auto p-6">
				<Card className="border-destructive">
					<CardHeader>
						<CardTitle className="text-destructive">
							Error Loading Loan
						</CardTitle>
					</CardHeader>
					<CardContent>
						<p>
							There was an error loading the loan details. Please
							try again later.
						</p>
					</CardContent>
				</Card>
			</div>
		);
	}

	// if (!data?.data) {
	// 	return (
	// 		<div className="container mx-auto p-6">
	// 			<Card>
	// 				<CardHeader>
	// 					<CardTitle>Loan Not Found</CardTitle>
	// 				</CardHeader>
	// 				<CardContent>
	// 					<p>The requested loan could not be found.</p>
	// 				</CardContent>
	// 			</Card>
	// 		</div>
	// 	);
	// }

	const progress = (data.data.paidAmount / data.data.repayAmount) * 100;
	const remainingAmount = data.data.repayAmount - data.data.paidAmount;

	return (
		<div className="container mx-auto p-6">
			<div className="flex flex-col gap-6">
				{/* Header */}
				<div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
					<div>
						<h1 className="text-3xl font-bold">
							Loan #{data.data.id}
						</h1>
						<p className="text-muted-foreground">
							Client: {data.data.clientDetails.fullName} |
							Disbursed: {formatDate(data.data.disbursedOn)}
						</p>
					</div>
					<div className="flex gap-2">
						<Badge>{data.data.status}</Badge>
					</div>
				</div>

				{/* Loan Summary and Client Details */}
				<div className="grid grid-cols-1 md:grid-cols-2 gap-6">
					<LoanDetailsCard loan={data.data} />
					<ClientDetailsCard client={data.data.clientDetails} />
				</div>

				{/* Tabs for Installments and Payments */}
				<Tabs defaultValue="installments" className="w-full">
					<TabsList className="grid w-full grid-cols-2">
						<TabsTrigger value="installments">
							Installments
						</TabsTrigger>
						<TabsTrigger value="payments">Payments</TabsTrigger>
					</TabsList>
					<TabsContent value="installments">
						<Card>
							<CardHeader>
								<CardTitle className="flex items-center gap-2">
									<CalendarIcon className="h-5 w-5" />
									Installment Schedule
								</CardTitle>
								<CardDescription>
									View all installments for this loan and
									their payment status
								</CardDescription>
							</CardHeader>
							<CardContent>
								<LoanDetailsInstallmentsTable
									installments={data.data.installments}
								/>
							</CardContent>
						</Card>
					</TabsContent>
					<TabsContent value="payments">
						<Card>
							<CardHeader>
								<CardTitle className="flex items-center gap-2">
									<DollarSign className="h-5 w-5" />
									Payment History
								</CardTitle>
								<CardDescription>
									Click on a payment to see how it was
									allocated to installments
								</CardDescription>
							</CardHeader>
							<CardContent>
								<LoanDetailsPaymentsTable
									payments={data.data.nonPosted}
									onPaymentClick={handlePaymentClick}
								/>
							</CardContent>
						</Card>
					</TabsContent>
				</Tabs>
			</div>

			{/* Payment Allocation Modal */}
			{selectedPayment && (
				<LoanDetailsPaymentAllocationModal
					payment={selectedPayment}
					allocations={data.data.paymentAllocations.filter(
						(a) => a.nonPostedId === selectedPayment.id,
					)}
					open={isAllocationModalOpen}
					onClose={() => setIsAllocationModalOpen(false)}
				/>
			)}
		</div>
	);
}

export default LoanDetailsPage;
