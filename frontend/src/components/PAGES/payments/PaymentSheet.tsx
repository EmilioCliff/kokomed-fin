import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTable } from '@/hooks/useTable';
import { useState } from 'react';
import { useAuth } from '@/hooks/useAuth';
import updatePayment from '@/services/updatePayment';
import { updatePaymentType } from '@/lib/types';
import { toast } from 'react-toastify';
import { role } from '@/lib/types';
import getFormData from '@/services/getFormData';
import { Input } from '@/components/ui/input';
import ClientCardDisplay from '@/components/UI/ClientCardDisplay';
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from '@/components/ui/sheet';
import { Button } from '@/components/ui/button';
import VirtualizeddSelect from '../../UI/VisualizedSelect';
import { format } from 'date-fns';

function PaymentSheet() {
	const [clientId, setClientId] = useState(0);
	const { selectedRow, setSelectedRow } = useTable();
	const { decoded } = useAuth();

	const { data } = useQuery({
		queryKey: ['payments/form'],
		queryFn: () => getFormData(false, true, false, false, false),
		staleTime: 5 * 1000,
	});

	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: updatePayment,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['payments'] });
			toast.success('Payment Assigned Successfully');
		},
		onError: (error) => {
			toast.error(error.message);
		},
		onSettled: () => {
			setClientId(0);
			setSelectedRow(null);
		},
	});

	const onSave = () => {
		const values: updatePaymentType = {
			id: Number(selectedRow.id),
			clientId: clientId ? clientId : 0,
		};

		mutation.mutate(values);
	};

	const fieldRenderers: Record<string, (value: any) => JSX.Element> = {
		id: (value) => (
			<Input
				readOnly
				placeholder={`P${String(value).padStart(3, '0')}`}
				className="readonly-input"
			/>
		),
		assignedTo: (value) => {
			return value.fullName ? (
				<ClientCardDisplay client={value} />
			) : decoded?.role === role.ADMIN && data?.client ? (
				<VirtualizeddSelect
					options={data.client}
					placeholder="Assign To"
					value={value}
					onChange={(id) => setClientId(id)}
				/>
			) : (
				<Input readOnly placeholder="N/A" className="readonly-input" />
			);
		},
		amount: (value) => {
			return (
				<Input
					readOnly
					placeholder={value.toLocaleString()}
					className="readonly-input"
				/>
			);
		},
		assigned: (value) => {
			return (
				<Input
					readOnly
					placeholder={value ? 'YES' : 'NO'}
					className="readonly-input"
				/>
			);
		},
		assignedBy: (value) => {
			return (
				<Input
					readOnly
					placeholder={value}
					className="readonly-input"
				/>
			);
		},
		paidDate: (value) => {
			return (
				<Input
					readOnly
					placeholder={format(value, 'PPP')}
					className="readonly-input"
				/>
			);
		},
	};

	return (
		<Sheet
			open={!!selectedRow}
			onOpenChange={(open: boolean) => {
				if (!open) {
					setSelectedRow(null);
					setClientId(0);
				}
			}}
		>
			<SheetContent className="overflow-auto custom-sheet-class">
				<SheetHeader>
					<SheetTitle>Payment Details</SheetTitle>
					<SheetDescription>Description goes here</SheetDescription>
				</SheetHeader>
				{selectedRow && (
					<div className="py-4">
						{Object.entries(selectedRow).map(([key, value]) => {
							if (fieldRenderers[key]) {
								return (
									<div
										key={key}
										className="grid grid-cols-[0.5fr_1fr] mb-4"
									>
										<span className="font-medium capitalize">
											{key}
										</span>
										{fieldRenderers[key](value)}
									</div>
								);
							}

							return (
								<div
									key={key}
									className="grid grid-cols-[0.5fr_1fr] mb-4"
								>
									<span className="font-medium capitalize">
										{key}
									</span>
									{typeof value === 'string' ||
									typeof value === 'number' ||
									typeof value === 'boolean' ? (
										<Input
											readOnly
											placeholder={value.toString()}
											className="readonly-input"
										/>
									) : (
										JSON.stringify(value)
									)}
								</div>
							);
						})}
						<Button size="lg" onClick={onSave} className="mt-8">
							Save
						</Button>
					</div>
				)}
			</SheetContent>
		</Sheet>
	);
}

export default PaymentSheet;
