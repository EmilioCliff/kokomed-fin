import * as React from 'react';
import { CalendarIcon } from 'lucide-react';
import { format } from 'date-fns';
import type { DateRange } from 'react-day-picker';

import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import { useTable } from '@/hooks/useTable';

// get date range
export function DateRangePicker({
	className,
}: React.HTMLAttributes<HTMLDivElement>) {
	const { setTimeRange } = useTable();

	const [date, setDate] = React.useState<DateRange | undefined>({
		from: new Date(2024, 12, 1),
		to: new Date(),
	});

	React.useEffect(() => {
		setTimeRange(
			format(date?.from ?? '01/01/2025', 'P'),
			format(date?.to ?? new Date(), 'P'),
		);
	}, [date]);

	return (
		<div className={cn('grid gap-2', className)}>
			<Popover>
				<PopoverTrigger asChild>
					<Button
						id="date"
						variant={'outline'}
						className={cn(
							'w-[300px] justify-start text-left font-normal',
							!date && 'text-muted-foreground',
						)}
					>
						<CalendarIcon className="mr-2 h-4 w-4" />
						{date?.from ? (
							date.to ? (
								<>
									{format(date.from, 'LLL dd, y')} -{' '}
									{format(date.to, 'LLL dd, y')}
								</>
							) : (
								format(date.from, 'LLL dd, y')
							)
						) : (
							<span>Pick a date</span>
						)}
					</Button>
				</PopoverTrigger>
				<PopoverContent className="w-auto p-0" align="start">
					<Calendar
						initialFocus
						mode="range"
						defaultMonth={date?.from}
						selected={date}
						onSelect={setDate}
						numberOfMonths={2}
					/>
				</PopoverContent>
			</Popover>
		</div>
	);
}
