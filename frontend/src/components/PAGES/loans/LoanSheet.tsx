import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from '@/components/ui/sheet';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { useState } from 'react';
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectLabel,
	SelectTrigger,
	SelectValue,
	SelectSeparator,
} from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Calendar } from '@/components/ui/calendar';
import UserCardDisplay from '@/components/UI/UserCardDisplay';
import ClientCardDisplay from '@/components/UI/ClientCardDisplay';
import ProductCardDisplay from '@/components/UI/ProductCardDisplay';
import { CalendarIcon } from 'lucide-react';
import { cn } from '@/lib/utils';
import { format } from 'date-fns';
import updateLoan from '@/services/updateLoan';
import { updateLoanType, loanStatus } from '@/lib/types';
import { toast } from 'react-toastify';
import { useAuth } from '@/hooks/useAuth';
import { useTable } from '@/hooks/useTable';
import { role } from '@/lib/types';
import { useNavigate } from 'react-router';

function LoanSheet() {
	const [status, setStatus] = useState<loanStatus | null>(null);
	const [disbursedDate, setDisbursedDate] = useState<string | null>(null);
	const [feePaid, setFeePaid] = useState<boolean | null>(null);
	const { selectedRow, setSelectedRow } = useTable();
	const { decoded } = useAuth();
	const navigation = useNavigate();

	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: updateLoan,
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: ['loans'] });
			toast.success('Load Updated');
		},
		onError: (error) => {
			toast.error(error.message);
		},
		onSettled: () => {
			setStatus(null);
			setDisbursedDate(null);
			setFeePaid(null);
			setSelectedRow(null);
		},
	});

	const onSave = () => {
		const values: updateLoanType = {
			id: Number(selectedRow.id),
			status: status ? status : undefined,
			disburseDate: disbursedDate
				? format(disbursedDate, 'yyyy-MM-dd')
				: undefined,
			feePaid: feePaid ? feePaid : undefined,
		};

		mutation.mutate(values);
	};

	const fieldRenderers: Record<string, (value: any) => JSX.Element> = {
		id: (value) => (
			<Input
				readOnly
				placeholder={`LN${String(value).padStart(3, '0')}`}
				className="readonly-input"
			/>
		),
		product: (value) => <ProductCardDisplay product={value} />,
		client: (value) => <ClientCardDisplay client={value} />,
		loanOfficer: (value) => <UserCardDisplay user={value} />,
		loanPurpose: (value) => <>{value}</>,
		approvedBy: (value) => <UserCardDisplay user={value} />,
		disbursedBy: (value) => {
			if (value.id === 0) {
				return (
					<Input
						readOnly
						placeholder="-"
						className="readonly-input"
					/>
				);
			}
			return <UserCardDisplay user={value} />;
		},
		status: (value: string) => {
			return value === 'INACTIVE' && decoded?.role === role.ADMIN ? (
				<div>
					<Select
						defaultValue="INACTIVE"
						onValueChange={(value) => {
							const tmp =
								value === 'INACTIVE'
									? loanStatus.INACTIVE
									: loanStatus.ACTIVE;
							setStatus(tmp);
						}}
					>
						<SelectTrigger className="w-[180px]">
							<SelectValue placeholder="INACTIVE" />
						</SelectTrigger>
						<SelectContent>
							<SelectGroup>
								<SelectItem value="INACTIVE">
									INACTIVE
								</SelectItem>
								<SelectItem value="ACTIVE">ACTIVE</SelectItem>
								<SelectSeparator />
								<SelectLabel>Disbursed Date</SelectLabel>
								<Popover>
									<PopoverTrigger asChild>
										<Button
											variant={'outline'}
											className={cn(
												'w-[240px] pl-3 text-left font-normal',
												!disbursedDate &&
													'text-muted-foreground',
											)}
										>
											{disbursedDate ? (
												format(
													disbursedDate,
													'yyyy-MM-dd',
												)
											) : (
												<span>Pick a date</span>
											)}
											<CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
										</Button>
									</PopoverTrigger>
									<PopoverContent
										className="w-auto p-0"
										align="start"
									>
										<Calendar
											mode="single"
											selected={
												disbursedDate
													? new Date(disbursedDate)
													: undefined
											}
											onSelect={(date: any) =>
												setDisbursedDate(
													format(date, 'yyyy-MM-dd'),
												)
											}
											disabled={(date) =>
												date > new Date() ||
												date < new Date('1900-01-01')
											}
											initialFocus
										/>
									</PopoverContent>
								</Popover>
							</SelectGroup>
						</SelectContent>
					</Select>
					{disbursedDate && (
						<p>
							Disbursed Date:{' '}
							{format(new Date(disbursedDate), 'PPP')}
						</p>
					)}
				</div>
			) : (
				<Input
					readOnly
					placeholder={value.toString()}
					className="readonly-input"
				/>
			);
		},
		feePaid: (value) => {
			return value === false ? (
				<Select
					defaultValue="NO"
					onValueChange={(value) => setFeePaid(value === 'YES')}
				>
					<SelectTrigger className="w-[180px]">
						<SelectValue placeholder="NO" />
					</SelectTrigger>
					<SelectContent>
						<SelectGroup>
							<SelectItem value="YES">YES</SelectItem>
							<SelectItem value="NO">NO</SelectItem>
						</SelectGroup>
					</SelectContent>
				</Select>
			) : (
				<Input readOnly placeholder="YES" className="readonly-input" />
			);
		},
		dueDate: (value) => {
			if (!value || value === '0001-01-01T00:00:00Z') {
				value = '-';
			}

			return (
				<Input
					readOnly
					placeholder={value === '-' ? value : format(value, 'PPP')}
					className="readonly-input"
				/>
			);
		},
		disbursedOn: (value) => {
			if (!value || value === '0001-01-01T00:00:00Z') {
				value = '-';
			}

			return (
				<Input
					readOnly
					placeholder={value === '-' ? value : format(value, 'PPP')}
					className="readonly-input"
				/>
			);
		},
		createdAt: (value) => {
			return (
				<Input
					readOnly
					placeholder={format(value, 'PPP')}
					className="readonly-input"
				/>
			);
		},
		updatedBy: () => <></>,
		createdBy: () => <></>,
	};

	return (
		<Sheet
			open={!!selectedRow}
			onOpenChange={(open: boolean) => {
				if (!open) {
					setSelectedRow(null);
					setDisbursedDate(null);
					setFeePaid(null);
				}
			}}
		>
			<SheetContent className="overflow-auto no-scrollbar custom-sheet-class">
				<SheetHeader>
					<SheetTitle>Loan Details</SheetTitle>
					<SheetDescription>Description goes here</SheetDescription>
				</SheetHeader>
				{selectedRow && (
					<div className="py-4">
						{Object.entries(selectedRow).map(([key, value]) => {
							if (key === 'createdBy' || key === 'updatedBy') {
								return;
							}
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
						<div className="flex justify-between">
							<Button
								variant={'outline'}
								size="lg"
								onClick={() =>
									navigation(
										`/loans/overview/${selectedRow.id}`,
									)
								}
								className="mt-8"
							>
								View
							</Button>
							<Button size="lg" onClick={onSave} className="mt-8">
								Save
							</Button>
						</div>
					</div>
				)}
			</SheetContent>
		</Sheet>
	);
}

export default LoanSheet;
