import { useTable } from '@/hooks/useTable';

import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from '@/components/ui/sheet';
import { Input } from '@/components/ui/input';

function ProductSheet() {
	const { selectedRow, setSelectedRow } = useTable();

	const fieldRenderers: Record<string, (value: any) => JSX.Element> = {
		id: (value) => (
			<Input
				readOnly
				placeholder={`P${String(value).padStart(3, '0')}`}
				className="readonly-input"
			/>
		),
	};

	return (
		<Sheet
			open={!!selectedRow}
			onOpenChange={(open: boolean) => {
				if (!open) {
					setSelectedRow(null);
				}
			}}
		>
			<SheetContent className="overflow-auto custom-sheet-class">
				<SheetHeader>
					<SheetTitle>Product Details</SheetTitle>
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
					</div>
				)}
			</SheetContent>
		</Sheet>
	);
}

export default ProductSheet;
