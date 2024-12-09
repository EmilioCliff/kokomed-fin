import { z } from "zod";
import { generateRandomLoans } from "@/components/lib/generator";
import { useState, useEffect } from "react";

import LoanForm from "@/components/forms/LoanForm";

import { loanColumns } from "../components/table/columns";
import { DataTable } from "../components/table/data-table";
import { loanSchema, Loan } from "../data/schema";
import { isUser, isClient } from "@/components/lib/validators";
import UserCardDisplay from "@/components/UserCardDisplay";
import ClientCardDisplay from "@/components/ClientCardDisplay";
import {
	Sheet,
	SheetContent,
	SheetHeader,
	SheetTitle,
} from "@/components/ui/sheet";
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

const fieldRenderers: Record<string, (value: any) => JSX.Element> = {
	client: (value) => <ClientCardDisplay client={isClient(value)} />,
	loanOfficer: (value) => <UserCardDisplay user={isUser(value)} />,
	approvedBy: (value) => <UserCardDisplay user={isUser(value)} />,
	disbursedBy: (value) => <UserCardDisplay user={isUser(value)} />,
	updatedBy: (value) => <UserCardDisplay user={isUser(value)} />,
	createdBy: (value) => <UserCardDisplay user={isUser(value)} />,
};

const generatedLoans = generateRandomLoans(30);
const validatedLoans = z.array(loanSchema).parse(generatedLoans);

export default function LoanPage() {
	const [loans, setLoans] = useState<Loan[]>([]);
	const [loading, setLoading] = useState(true);
	const [selectedRow, setSelectedRow] = useState<Loan | null>(null);

	useEffect(() => {
		setLoans(validatedLoans);
		setLoading(false);
	}, [validatedLoans]);

	if (loading) {
		return <div>Loading...</div>;
	}

	return (
		<>
			<div className='flex justify-between items-center mb-4'>
				<h1 className='text-3xl font-bold'>Loans</h1>
				<Dialog>
					<DialogTrigger asChild>
						<Button className='text-xs py-1 font-bold' size='sm'>
							Add New Loan
						</Button>
					</DialogTrigger>
					<DialogContent className='max-w-screen-lg '>
						<DialogHeader>
							<DialogTitle>Add New Loan</DialogTitle>
							<DialogDescription>
								Enter the details for the new loan.
							</DialogDescription>
						</DialogHeader>
						<LoanForm />
					</DialogContent>
				</Dialog>
			</div>
			<DataTable
				data={loans}
				columns={loanColumns}
				setSelectedRow={setSelectedRow}
			/>
			<Sheet
				open={!!selectedRow}
				onOpenChange={(open) => !open && setSelectedRow(null)}
			>
				<SheetContent className='overflow-auto custom-sheet-class'>
					<SheetHeader>
						<SheetTitle>Payment Details</SheetTitle>
					</SheetHeader>
					{selectedRow && (
						<div className='py-4'>
							{Object.entries(selectedRow).map(([key, value]) => {
								if (fieldRenderers[key]) {
									return (
										<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
											<span className='font-medium'>{key}</span>
											{fieldRenderers[key](value)}
										</div>
									);
								}

								// Default rendering for keys without a custom renderer
								return (
									<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
										<span className='font-medium'>{key}</span>
										{key === "id"
											? `LN${String(value).padStart(3, "0")}`
											: typeof value === "string"
											? value
											: JSON.stringify(value)}
									</div>
								);
							})}
						</div>
					)}
				</SheetContent>
			</Sheet>
		</>
	);
}
