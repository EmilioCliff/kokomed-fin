import { Card, CardContent, CardFooter, CardHeader } from '../ui/card';
import { Skeleton } from '../ui/skeleton';

function EditPaymentFormSkeleton() {
	return (
		<div className="container mx-auto p-6">
			<div className="flex flex-col gap-6">
				{/* Header Skeleton */}
				<div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
					<div>
						<Skeleton className="h-8 w-48 mb-2" />
						<Skeleton className="h-4 w-64" />
					</div>
				</div>

				<Card>
					<CardContent className="space-y-4">
						<div className="grid grid-cols-1 md:grid-cols-2 gap-8 py-4">
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />

								<Skeleton className="h-10 w-full" />
							</div>
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-10 w-full" />
							</div>
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-10 w-full" />
							</div>
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-10 w-full" />
							</div>
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-10 w-full" />
							</div>
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-10 w-full" />
							</div>
							<div className="col-span-2">
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-20 w-full" />
							</div>
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-10 w-full" />
							</div>
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-10 w-full" />
							</div>
							<div>
								<Skeleton className="h-4 w-1/4 mb-2" />
								<Skeleton className="h-10 w-full" />
							</div>
						</div>
					</CardContent>
					<CardFooter>
						<Skeleton className="h-8 w-1/6" />
						<Skeleton className="h-8 w-1/6 ml-auto" />
					</CardFooter>
				</Card>
			</div>
		</div>
	);
}

export default EditPaymentFormSkeleton;
