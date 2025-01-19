import { Skeleton } from "../ui/skeleton";
import TableSkeleton from "./TableSkeleton";

function DashboardSkeleton() {
	return (
		<div className='px-4'>
			{/* <h1 className='text-3xl font-bold mb-6 text-start'>Dashboard</h1> */}
			<Skeleton className='h-8 w-32 mb-6' />
			<div className='space-y-6'>
				<div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4'>
					<Skeleton className='aspect-video ' />
					<Skeleton className='aspect-video ' />
					<Skeleton className='aspect-video ' />
					<Skeleton className='aspect-video ' />
				</div>
				<div className='grid grid-cols-1 lg:grid-cols-2 gap-4'>
					<Skeleton className='aspect-square ' />
					<Skeleton className='aspect-square ' />
				</div>
				<TableSkeleton />
			</div>
		</div>
	);
}

export default DashboardSkeleton;
