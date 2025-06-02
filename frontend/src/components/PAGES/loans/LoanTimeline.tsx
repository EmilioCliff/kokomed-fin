import { useState, useMemo } from 'react';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import FullCalendar from '@fullcalendar/react';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from '@/components/ui/sheet';
import getLoanEvents from '@/services/getLoanEvents';
import LoanEventDisplay from '@/components/UI/LoanEventDisplay';

export default function LoanTimeline() {
	const [sheetOpen, setSheetOpen] = useState(false);
	const [selectedDate, setSelectedDate] = useState<string | null>(null);
	const [selectedEvents, setSelectedEvents] = useState<any[]>([]);

	const { isLoading, error, data } = useQuery({
		queryKey: ['loans/events'],
		queryFn: () => getLoanEvents(),
		staleTime: 5 * 1000,
		placeholderData: keepPreviousData,
	});

	const eventsByDate = useMemo(() => {
		const map = new Map<string, any[]>();
		data?.data.forEach((event) => {
			if (!map.has(event.date)) {
				map.set(event.date, []);
			}
			map.get(event.date)?.push(event);
		});
		return map;
	}, [data?.data]);

	const handleDateClick = (arg: any) => {
		const eventsForDate = eventsByDate.get(arg.dateStr) || [];
		setSelectedEvents(eventsForDate);
		setSelectedDate(arg.dateStr);
		setSheetOpen(true);
	};

	const renderEventContent = (eventInfo: any) => {
		return (
			<div className="relative flex justify-center items-center">
				<span
					className={`h-2 w-2 rounded-full mx-2 ${
						eventInfo.event.extendedProps.type === 'disbursed'
							? 'bg-green-500'
							: 'bg-red-500'
					}`}
				></span>
				<p>{eventInfo.event.id}</p>
			</div>
		);
	};

	return (
		<div>
			<FullCalendar
				plugins={[dayGridPlugin, interactionPlugin]}
				initialView="dayGridMonth"
				dateClick={handleDateClick}
				selectable={true}
				events={data?.data}
				eventContent={renderEventContent}
				headerToolbar={{
					left: 'prev,next',
					center: 'title',
					right: 'dayGridYear,dayGridMonth,dayGridWeek',
				}}
			/>

			<Sheet
				open={sheetOpen}
				onOpenChange={() => setSheetOpen(!sheetOpen)}
			>
				<SheetContent className="overflow-auto no-scrollbar">
					<SheetHeader>
						<SheetTitle>{selectedDate}</SheetTitle>
						<SheetDescription>
							Here are the loans scheduled for this day.
						</SheetDescription>
					</SheetHeader>

					<div className="mt-4 space-y-3">
						{selectedEvents.length > 0 ? (
							selectedEvents.map((event) => (
								<LoanEventDisplay event={event} />
							))
						) : (
							<p className="text-gray-500 text-center">
								No loan events for this day.
							</p>
						)}
					</div>
				</SheetContent>
			</Sheet>
		</div>
	);
}
