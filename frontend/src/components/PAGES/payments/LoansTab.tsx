import { Card, CardTitle } from '@/components/ui/card';
import { Loan } from '../loans/schema';
import { clientLoanColumns } from './payment';
import { DataTable } from '@/components/table/data-table';
import { statuses } from '@/data/loan';

interface loansTabProps {
	loans: Loan[];
	searchType: string;
}

function LoansTab({ loans, searchType }: loansTabProps) {
	return (
		<Card className="p-4 mb-6">
			{loans.length === 0 ? (
				<div className="text-center py-6">
					<p className="text-gray-500">
						No loans found for this{' '}
						{searchType === 'id' ? 'client.' : 'number.'}
					</p>
				</div>
			) : (
				<div>
					<div className="flex justify-between items-center mb-4">
						<CardTitle>Loan History</CardTitle>
					</div>
					<DataTable
						data={loans || []}
						columns={clientLoanColumns}
						searchableColumns={[
							{
								id: 'loanOfficerName',
								title: 'Loan Officer',
							},
						]}
						facetedFilterColumns={[
							{
								id: 'status',
								title: 'Status',
								options: statuses,
							},
						]}
					/>
				</div>
			)}
		</Card>
	);
}

export default LoansTab;
