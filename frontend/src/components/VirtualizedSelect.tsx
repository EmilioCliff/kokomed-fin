import { useState, useMemo } from "react";
import {
	Select,
	SelectTrigger,
	SelectValue,
	SelectContent,
} from "@/components/ui/select";
import { FixedSizeList as List } from "react-window";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

const loanOfficers = Array.from({ length: 300 }, (_, i) => ({
	id: i,
	name: `Loan Officer ${i + 1}`,
}));

const VirtualizedSelect = ({ onValueChange }: any) => {
	const [searchQuery, setSearchQuery] = useState("");

	const filteredOfficers = useMemo(() => {
		return loanOfficers.filter((officer) =>
			officer.name.toLowerCase().includes(searchQuery.toLowerCase())
		);
	}, [searchQuery]);

	return (
		<Select onValueChange={onValueChange}>
			<SelectTrigger>
				<SelectValue placeholder='Select a loan officer' />
			</SelectTrigger>
			<SelectContent className='w-full'>
				{/* Search Bar */}
				<div className='p-2'>
					<Input
						placeholder='Search loan officers...'
						value={searchQuery}
						onChange={(e) => setSearchQuery(e.target.value)}
					/>
				</div>

				{/* Virtualized List */}
				<List
					height={200}
					itemCount={filteredOfficers.length}
					itemSize={40}
					width='100%'
				>
					{({ index, style }: any) => (
						<div
							style={style}
							className={cn(
								"px-4 py-2 cursor-pointer hover:bg-gray-100",
								index % 2 && "bg-gray-50"
							)}
							onClick={() => onValueChange(filteredOfficers[index].id)}
						>
							{filteredOfficers[index].name}
						</div>
					)}
				</List>
			</SelectContent>
		</Select>
	);
};

export default VirtualizedSelect;
