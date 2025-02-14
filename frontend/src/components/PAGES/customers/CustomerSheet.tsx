import { useState } from 'react';
import { updateClientType } from '@/lib/types';
import updateClient from '@/services/updateClient';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import getFormData from '@/services/getFormData';
import { useTable } from '@/hooks/useTable';
import { useAuth } from '@/hooks/useAuth';
import { toast } from 'react-toastify';
import VirtualizeddSelect from '../../UI/VisualizedSelect';
import UserCardDisplay from '@/components/UI/UserCardDisplay';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from '@/components/ui/sheet';
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import { format } from 'date-fns';
import { Calendar } from '@/components/ui/calendar';
import { CalendarIcon } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { role } from '@/lib/types';

function CustomerSheet() {
	const [idNumber, setIdNumber] = useState<string | undefined>(undefined);
	const [dob, setDob] = useState<string | undefined>(undefined);
	const [branchId, setBranchId] = useState<number | undefined>(undefined);
	const [active, setActive] = useState<string | undefined>(undefined);
	const { selectedRow, setSelectedRow } = useTable();
	const { decoded } = useAuth();

	const { data } = useQuery({
		queryKey: ['loans/form'],
		queryFn: () => getFormData(false, false, false, true, false),
		staleTime: 5 * 1000,
	});

	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: updateClient,
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: ['clients'] });
			toast.success('User Updated');
		},
		onError: (error) => {
			toast.error(error.message);
		},
		onSettled: () => {
			setBranchId(undefined);
			setSelectedRow(null);
		},
	});

	const onSave = () => {
		const values: updateClientType = {
			id: Number(selectedRow.id),
			idNumber: idNumber ? idNumber : undefined,
			branchId: branchId ? branchId : undefined,
			dob: dob ? dob : undefined,
			active: active ? active : undefined,
		};

		mutation.mutate(values);
	};

	const fieldRenderers: Record<string, (value: any) => JSX.Element> = {
		id: (value) => (
			<Input
				readOnly
				placeholder={`CM${String(value).padStart(3, '0')}`}
				className="readonly-input"
			/>
		),
		idNumber: (value) => {
			return value ? (
				<Input
					readOnly
					placeholder={value}
					className="readonly-input"
				/>
			) : (
				<Input
					id="idNumber"
					type="text"
					value={idNumber!}
					placeholder="Enter ID Number"
					onChange={(e) => setIdNumber(e.target.value)}
					className="readonly-input"
				/>
			);
		},
		dob: (value) => {
			return value ? (
				<Button
					variant={'outline'}
					className={cn(
						'w-[240px] pl-3 text-left font-normal',
						!dob && 'text-muted-foreground',
					)}
				>
					{format(value, 'yyyy-MM-dd')}
					<CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
				</Button>
			) : (
				<Popover>
					<PopoverTrigger asChild>
						<Button
							variant={'outline'}
							className={cn(
								'w-[240px] pl-3 text-left font-normal',
								!dob && 'text-muted-foreground',
							)}
						>
							{dob ? (
								format(dob, 'PPP')
							) : (
								<span>Pick a date</span>
							)}
							<CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
						</Button>
					</PopoverTrigger>
					<PopoverContent className="w-auto p-0" align="start">
						<Calendar
							mode="single"
							selected={dob ? new Date(dob) : undefined}
							onSelect={(date: any) =>
								setDob(format(date, 'yyyy-MM-dd'))
							}
							disabled={(date) =>
								date > new Date() ||
								date < new Date('1900-01-01')
							}
							initialFocus
						/>
					</PopoverContent>
				</Popover>
			);
		},
		branchName: (value) => {
			return decoded?.role === role.ADMIN && data?.branch ? (
				<VirtualizeddSelect
					options={data.branch}
					placeholder={value}
					value={value}
					onChange={(id) => setBranchId(id)}
				/>
			) : (
				<Input
					readOnly
					placeholder={value.toString()}
					className="readonly-input"
				/>
			);
		},
		active: (value) => {
			return (
				<div>
					<Select
						defaultValue={value ? 'ACTIVE' : 'INACTIVE'}
						onValueChange={(value) => {
							const tmp = value === 'ACTIVE' ? 'true' : 'false';
							setActive(tmp);
						}}
					>
						<SelectTrigger className="w-[180px]">
							<SelectValue placeholder={value} />
						</SelectTrigger>
						<SelectContent>
							<SelectGroup>
								<SelectItem value="ACTIVE">ACTIVE</SelectItem>
								<SelectItem value="INACTIVE">
									INACTIVE
								</SelectItem>
							</SelectGroup>
						</SelectContent>
					</Select>
				</div>
			);
		},
		createdBy: (value) => <UserCardDisplay user={value} />,
		assignedStaff: (value) => <UserCardDisplay user={value} />,
		createdAt: (value) => {
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
					setIdNumber(undefined);
					setDob(undefined);
					setActive(undefined);
					setBranchId(undefined);
				}
			}}
		>
			<SheetContent className="overflow-auto custom-sheet-class">
				<SheetHeader>
					<SheetTitle>Customer Details</SheetTitle>
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

export default CustomerSheet;
