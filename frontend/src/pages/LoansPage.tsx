// import { Table } from "@/components/Table";
// import { Badge } from "@/components/ui/badge";
// import { Button } from "@/components/ui/button";
// import {
// 	Dialog,
// 	DialogContent,
// 	DialogDescription,
// 	DialogFooter,
// 	DialogHeader,
// 	DialogTitle,
// 	DialogTrigger,
// } from "@/components/ui/dialog";
// import { Input } from "@/components/ui/input";
// import { Label } from "@/components/ui/label";
// import { useState } from "react";

// const columns = [
// 	{ key: "id", label: "ID", sortable: true },
// 	{ key: "customerName", label: "Customer Name", sortable: true },
// 	{ key: "amount", label: "Amount", sortable: true },
// 	{
// 		key: "status",
// 		label: "Status",
// 		sortable: true,
// 		render: (value: string) => (
// 			<Badge variant={value === "Active" ? "default" : "secondary"}>
// 				{value}
// 			</Badge>
// 		),
// 	},
// 	{ key: "dueDate", label: "Due Date", sortable: true },
// ];

// const initialData = [
// 	{
// 		id: 1,
// 		customerName: "John Doe",
// 		amount: 5000,
// 		status: "Active",
// 		dueDate: "2023-12-31",
// 	},
// 	{
// 		id: 2,
// 		customerName: "Jane Smith",
// 		amount: 7500,
// 		status: "Inactive",
// 		dueDate: "2023-11-30",
// 	},
// 	{
// 		id: 3,
// 		customerName: "Bob Johnson",
// 		amount: 10000,
// 		status: "Active",
// 		dueDate: "2024-01-15",
// 	},
// ];

// export default function LoansPage() {
// 	const [data, setData] = useState(initialData);
// 	const [newLoan, setNewLoan] = useState({
// 		customerName: "",
// 		amount: "",
// 		status: "Active",
// 		dueDate: "",
// 	});

// 	const handleAddLoan = () => {
// 		const newId = Math.max(...data.map((loan) => loan.id)) + 1;
// 		setData([
// 			...data,
// 			{ id: newId, ...newLoan, amount: Number(newLoan.amount) },
// 		]);
// 		setNewLoan({ customerName: "", amount: "", status: "Active", dueDate: "" });
// 	};

// 	return (
// 		<div className='p-2 mt-14'>
// 			<div className='flex justify-between items-center mb-4'>
// 				<h1 className='text-3xl font-bold'>Loans</h1>
// 				<Dialog>
// 					<DialogTrigger asChild>
// 						<Button>Add New Loan</Button>
// 					</DialogTrigger>
// 					<DialogContent>
// 						<DialogHeader>
// 							<DialogTitle>Add New Loan</DialogTitle>
// 							<DialogDescription>
// 								Enter the details for the new loan.
// 							</DialogDescription>
// 						</DialogHeader>
// 						<div className='grid gap-4 py-4'>
// 							<div className='grid grid-cols-4 items-center gap-4'>
// 								<Label htmlFor='customerName' className='text-right'>
// 									Customer Name
// 								</Label>
// 								<Input
// 									id='customerName'
// 									value={newLoan.customerName}
// 									onChange={(e) =>
// 										setNewLoan({ ...newLoan, customerName: e.target.value })
// 									}
// 									className='col-span-3'
// 								/>
// 							</div>
// 							<div className='grid grid-cols-4 items-center gap-4'>
// 								<Label htmlFor='amount' className='text-right'>
// 									Amount
// 								</Label>
// 								<Input
// 									id='amount'
// 									type='number'
// 									value={newLoan.amount}
// 									onChange={(e) =>
// 										setNewLoan({ ...newLoan, amount: e.target.value })
// 									}
// 									className='col-span-3'
// 								/>
// 							</div>
// 							<div className='grid grid-cols-4 items-center gap-4'>
// 								<Label htmlFor='dueDate' className='text-right'>
// 									Due Date
// 								</Label>
// 								<Input
// 									id='dueDate'
// 									type='date'
// 									value={newLoan.dueDate}
// 									onChange={(e) =>
// 										setNewLoan({ ...newLoan, dueDate: e.target.value })
// 									}
// 									className='col-span-3'
// 								/>
// 							</div>
// 						</div>
// 						<DialogFooter>
// 							<Button onClick={handleAddLoan}>Add Loan</Button>
// 						</DialogFooter>
// 					</DialogContent>
// 				</Dialog>
// 			</div>
// 			<Table columns={columns} data={data} />
// 		</div>
// 	);
// }

import { z } from "zod";
import { generateRandomLoans } from "@/components/lib/generator";
import { useState, useEffect } from "react";

import { loanColumns } from "../components/loan/columns";
import { DataTable } from "../components/loan/data-table";
import { loanSchema, Loan } from "../data/schema";

export default function TaskPage() {
	const [loans, setLoans] = useState<Loan[]>([]); // State to hold the loan data
	const [loading, setLoading] = useState(true); // Loading state
	const [error, setError] = useState<string | null>(null); // Error state

	useEffect(() => {
		async function fetchLoans() {
			try {
				const generatedLoans = generateRandomLoans(30);
				const validatedLoans = z.array(loanSchema).parse(generatedLoans); // Validate loans
				setLoans(validatedLoans);
			} catch (err: unknown) {
				setError("Failed to fetch loans");
				console.error(err);
			} finally {
				setLoading(false);
			}
		}

		fetchLoans();
	}, []);

	return (
		<>
			{/* <div className='overflow-x-auto overflow-y-hidden flex-1 flex-col space-y-8 p-8 md:flex mt-16'> */}
			<DataTable data={loans} columns={loanColumns} />
			{/* </div> */}
		</>
	);
}
