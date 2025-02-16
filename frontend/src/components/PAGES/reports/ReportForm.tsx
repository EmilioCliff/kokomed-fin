import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { useState } from 'react';
import { CalendarIcon } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Report, ReportFilter } from '@/lib/types';
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { format } from 'date-fns';
import { Calendar } from '@/components/ui/calendar';
import getFormData from '@/services/getFormData';
import { useQuery, useMutation } from '@tanstack/react-query';
import VirtualizeddSelect from '@/components/UI/VisualizedSelect';
import { toast } from 'react-toastify';
import generateReport from '@/services/generateReport';

type ReportFormProps = {
	report: Report;
	onClose: () => void;
};

export function ReportForm({ report, onClose }: ReportFormProps) {
	const [filters, setFilters] = useState<ReportFilter>(report.filters);

	const { isLoading, data, error } = useQuery({
		queryKey: ['loans/form'],
		queryFn: () =>
			getFormData(
				false,
				filters.clientId !== undefined,
				filters.userId !== undefined,
				false,
				filters.loanId !== undefined,
			),
		staleTime: 5 * 1000,
	});

	const mutation = useMutation({
		mutationFn: generateReport,
		onSuccess: async () => {
			toast.success('Report Generated');
		},
		onError: (error) => {
			toast.error(error.message);
		},
		onSettled: () => {
			onClose();
		},
	});

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		setFilters({ ...filters, [e.target.name]: e.target.value });
	};

	const handleDateChange = (
		name: keyof ReportFilter,
		date: Date | undefined,
	) => {
		if (date) {
			setFilters({ ...filters, [name]: format(date, 'yyyy-MM-dd') });
		}
	};

	const handleVirtualizedChange = (name: keyof ReportFilter, id: number) => {
		setFilters({ ...filters, [name]: id });
	};

	const handleSelectChange = (name: keyof ReportFilter, format: string) => {
		setFilters({ ...filters, [name]: format });
	};

	const handleSubmit = () => {
		if (
			(filters.startDate !== undefined && !filters.startDate) ||
			(filters.endDate !== undefined && !filters.endDate)
		) {
			toast.error('Start Date and End Date are required!');
			return;
		}
		if (filters.loanId !== undefined && !filters.loanId) {
			toast.error('Loan is required for this report!');
			return;
		}
		if (filters.clientId !== undefined && !filters.clientId) {
			toast.error('Client is required for this report!');
			return;
		}
		if (filters.userId !== undefined && !filters.userId) {
			toast.error('User is required for this report!');
			return;
		}
		mutation.mutate(filters);
	};

	return (
		<div>
			<div className="grid grid-cols-1 md:grid-cols-2 gap-4 py-4">
				{filters.startDate !== undefined && (
					<Popover>
						<PopoverTrigger asChild>
							<Button
								variant={'outline'}
								className={cn(
									'pl-3 text-left font-normal',
									!filters.startDate &&
										'text-muted-foreground',
								)}
							>
								{filters.startDate ? (
									format(filters.startDate, 'PPP')
								) : (
									<span>Pick a date</span>
								)}
								<CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
							</Button>
						</PopoverTrigger>
						<PopoverContent className="w-auto p-0" align="start">
							<Calendar
								mode="single"
								selected={
									filters.startDate
										? new Date(filters.startDate)
										: undefined
								}
								onSelect={(date) =>
									handleDateChange('startDate', date)
								}
								disabled={(date) =>
									date < new Date('1900-01-01')
								}
								initialFocus
							/>
						</PopoverContent>
					</Popover>
				)}
				{filters.endDate !== undefined && (
					<Popover>
						<PopoverTrigger asChild>
							<Button
								variant={'outline'}
								className={cn(
									'pl-3 text-left font-normal',
									!filters.endDate && 'text-muted-foreground',
								)}
							>
								{filters.endDate ? (
									format(new Date(filters.endDate), 'PPP')
								) : (
									<span>Pick an end date</span>
								)}
								<CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
							</Button>
						</PopoverTrigger>
						<PopoverContent className="w-auto p-0" align="start">
							<Calendar
								mode="single"
								selected={
									filters.endDate
										? new Date(filters.endDate)
										: undefined
								}
								onSelect={(date) =>
									handleDateChange('endDate', date)
								}
								disabled={(date) =>
									date < new Date('1900-01-01')
								}
								initialFocus
							/>
						</PopoverContent>
					</Popover>
				)}

				{filters.loanId !== undefined && data?.loan && (
					<VirtualizeddSelect
						options={data.loan}
						placeholder="Select Loan"
						value={0}
						onChange={(id) => handleVirtualizedChange('loanId', id)}
					/>
				)}

				{filters.clientId !== undefined && data?.client && (
					<VirtualizeddSelect
						options={data.client}
						placeholder="Select Customer"
						value={0}
						onChange={(id) =>
							handleVirtualizedChange('clientId', id)
						}
					/>
				)}

				{filters.userId !== undefined && data?.user && (
					<VirtualizeddSelect
						options={data.user}
						placeholder="Select Staff"
						value={0}
						onChange={(id) => handleVirtualizedChange('userId', id)}
					/>
				)}
			</div>
			<Label>Report Format</Label>
			<Select
				onValueChange={(format) => handleSelectChange('format', format)}
			>
				<SelectTrigger className="w-full sm:w-[180px]">
					<SelectValue placeholder="Excel" />
				</SelectTrigger>
				<SelectContent>
					<SelectItem value="excel">Excel</SelectItem>
					<SelectItem value="pdf">PDF</SelectItem>
				</SelectContent>
			</Select>
			<Button className="ml-auto block" onClick={handleSubmit}>
				Generate Report
			</Button>
		</div>
	);
}
