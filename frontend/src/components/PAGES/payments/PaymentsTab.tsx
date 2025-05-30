import { DataTable } from '@/components/table/data-table';
import { Card, CardTitle } from '@/components/ui/card';
import { clientPaymentColumns } from './payment';
import { paymentSources } from '@/data/loan';
import { Payment } from './schema';

interface paymentsTabProps {
	payments: Payment[];
	searchType: string;
	total: number;
}

function PaymentsTab({ payments, searchType, total }: paymentsTabProps) {
	return (
		<Card className="p-4 mb-6">
			{payments.length === 0 ? (
				<div className="text-center py-6">
					<p className="text-gray-500">
						No payments found for this{' '}
						{searchType === 'id' ? 'client.' : 'number.'}
					</p>
				</div>
			) : (
				<div>
					<div className="flex justify-between items-center mb-4">
						<CardTitle>Payment History</CardTitle>

						<span className="text-lg font-semibold text-blue-600">
							Total: KES {`${total.toLocaleString()}`}
						</span>
					</div>
					<DataTable
						data={payments || []}
						columns={clientPaymentColumns}
						searchableColumns={[
							{
								id: 'payingName',
								title: 'Paying Name',
							},
							{
								id: 'accountNumber',
								title: 'Account Number',
							},
							{
								id: 'transactionNumber',
								title: 'Transaction Number',
							},
						]}
						facetedFilterColumns={[
							{
								id: 'transactionSource',
								title: 'Transaction Source',
								options: paymentSources,
							},
						]}
					/>
				</div>
			)}
		</Card>
	);
}

export default PaymentsTab;
