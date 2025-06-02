import { Card, CardContent, CardHeader } from '../ui/card';
import { Skeleton } from '../ui/skeleton';

function LoanDetailsSkeleton() {
	return (
		<div className="container mx-auto p-6">
			<div className="flex flex-col gap-6">
				{/* Header Skeleton */}
				<div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
					<div>
						<Skeleton className="h-8 w-48 mb-2" />
						<Skeleton className="h-4 w-64" />
					</div>
					<Skeleton className="h-6 w-20" />
				</div>

				{/* Summary and Client Details Skeleton */}
				<div className="grid grid-cols-1 md:grid-cols-2 gap-6">
					<Card>
						<CardHeader>
							<Skeleton className="h-6 w-32 mb-2" />
						</CardHeader>
						<CardContent className="space-y-4">
							<Skeleton className="h-4 w-full" />
							<Skeleton className="h-4 w-full" />
							<Skeleton className="h-4 w-3/4" />
						</CardContent>
					</Card>
					<Card>
						<CardHeader>
							<Skeleton className="h-6 w-32 mb-2" />
						</CardHeader>
						<CardContent className="space-y-4">
							<Skeleton className="h-4 w-full" />
							<Skeleton className="h-4 w-full" />
							<Skeleton className="h-4 w-3/4" />
						</CardContent>
					</Card>
				</div>
				<div>
					<div className="border-b flex mb-4">
						<Skeleton className="h-10 w-1/2" />
						<Skeleton className="h-10 w-1/2" />
					</div>
					<Card>
						<CardHeader>
							<Skeleton className="h-6 w-48 mb-2" />
							<Skeleton className="h-4 w-64" />
						</CardHeader>
						<CardContent>
							<div className="space-y-2">
								<Skeleton className="h-10 w-full" />
								<Skeleton className="h-10 w-full" />
								<Skeleton className="h-10 w-full" />
							</div>
						</CardContent>
					</Card>
				</div>
			</div>
		</div>
	);
}

export default LoanDetailsSkeleton;
