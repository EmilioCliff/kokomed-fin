import { generateRandomBranch } from "@/lib/generator";
import { z } from "zod";
import { branchSchema, Branch } from "@/data/schema";
import { useEffect, useState } from "react";
import TableSkeleton from "@/components/TableSkeleton";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from "@/components/ui/sheet";
import { DataTable } from "../components/table/data-table";
import BranchForm from "@/components/forms/BranchForm";
import { branchColumns } from "@/components/table/columns/branch";

const branches = generateRandomBranch(30);
const validatedBranches = z.array(branchSchema).parse(branches);

function BranchesPage() {
	const [branches, setBranches] = useState<Branch[]>([]);
	const [loading, setLoading] = useState(true);
	const [selectedRow, setSelectedRow] = useState<Branch | null>(null);

	useEffect(() => {
		setBranches(validatedBranches);
		setLoading(false);
	}, []);

	if (loading) {
		return <TableSkeleton />;
	}

	return (
		<div className='px-4'>
			<div className='flex justify-between items-center mb-4'>
				<h1 className='text-3xl font-bold'>Branches</h1>
				<Dialog>
					<DialogTrigger asChild>
						<Button className='text-xs py-1 font-bold' size='sm'>
							Add New Branch
						</Button>
					</DialogTrigger>
					<DialogContent className='max-w-screen-lg '>
						<DialogHeader>
							<DialogTitle>Add New Customer</DialogTitle>
							<DialogDescription>
								Enter the details for the new user.
							</DialogDescription>
						</DialogHeader>
						<BranchForm />
					</DialogContent>
				</Dialog>
			</div>
			<DataTable
				data={branches}
				columns={branchColumns}
				setSelectedRow={setSelectedRow}
				searchableColumns={[
					{
						id: "branchName",
						title: "Branch Name",
					},
				]}
			/>
			<Sheet
				open={!!selectedRow}
				onOpenChange={(open: boolean) => {
					if (!open) {
						setSelectedRow(null);
					}
				}}
			>
				<SheetContent className='overflow-auto custom-sheet-class'>
					<SheetHeader>
						<SheetTitle>Branch Details</SheetTitle>
						<SheetDescription>Description goes here</SheetDescription>
					</SheetHeader>
					{selectedRow && (
						<div className='py-4'>
							{Object.entries(selectedRow).map(([key, value]) => {
								// if (key === "createdBy" || key === "updatedBy") {
								// 	return;
								// }
								// if (fieldRenderers[key]) {
								// 	return (
								// 		<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
								// 			<span className='font-medium capitalize'>{key}</span>
								// 			{fieldRenderers[key](value)}
								// 		</div>
								// 	);
								// }

								// return (
								// 	<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
								// 		<span className='font-medium capitalize'>{key}</span>
								// 		{typeof value === "string" ||
								// 		typeof value === "number" ||
								// 		typeof value === "boolean" ? (
								// 			<Input
								// 				readOnly
								// 				placeholder={value.toString()}
								// 				className='bg-gray-100 text-gray-500'
								// 			/>
								// 		) : (
								// 			JSON.stringify(value)
								// 		)}
								// 	</div>
								// );
								return (
									<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
										<span className='font-medium capitalize'>{key}</span>
										<p>
											{typeof value == "string" ? (
												<p>{value}</p>
											) : (
												JSON.stringify(value)
											)}
										</p>
									</div>
								);
							})}
							<Button
								size='lg'
								onClick={() => console.log("CLicked")}
								className='mt-8'
							>
								Save
							</Button>
						</div>
					)}
				</SheetContent>
			</Sheet>
		</div>
	);
}

export default BranchesPage;
