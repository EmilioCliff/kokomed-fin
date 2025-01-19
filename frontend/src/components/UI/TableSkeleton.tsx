import { Skeleton } from '../ui/skeleton';

function TableSkeleton() {
  return (
    <div className="w-full">
      <div className="flex justify-between items-center mb-4">
        <Skeleton className="h-8 w-32" />
        <Skeleton className="h-8 w-32" />
      </div>
      <div className="flex items-center justify-between mt-2 mb-6">
        <div className="flex flex-1 items-center space-x-2">
          <Skeleton className="h-8 w-[250px]" />
          <Skeleton className="h-8 w-32" />
        </div>
        <Skeleton className="h-8 w-32" />
      </div>
      <div className="border border-gray-200 rounded-lg overflow-hidden">
        <div className="flex items-center justify-between px-4 py-3 bg-gray-100">
          <Skeleton className="h-6 w-40" />
          <Skeleton className="h-6 w-28" />
          <Skeleton className="h-6 w-28" />
          <Skeleton className="h-6 w-28" />
          <Skeleton className="h-6 w-20" />
        </div>
        {Array.from({ length: 10 }).map((_, index) => (
          <div
            key={index}
            className="flex items-center justify-between px-4 py-3 border-t border-gray-200"
          >
            <Skeleton className="h-6 w-40" />
            <Skeleton className="h-6 w-28" />
            <Skeleton className="h-6 w-28" />
            <Skeleton className="h-6 w-28" />
            <Skeleton className="h-6 w-20" />
          </div>
        ))}
      </div>
    </div>
  );
}

export default TableSkeleton;
