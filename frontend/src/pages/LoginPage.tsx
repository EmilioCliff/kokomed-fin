import { Table } from "@/components/Table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import {
	Sheet,
	SheetContent,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@/components/ui/sheet";

const columns = [
	{ key: "id", label: "ID", sortable: true },
	{ key: "customerName", label: "Customer Name", sortable: true },
	{ key: "amount", label: "Amount", sortable: true },
	{
		key: "status",
		label: "Status",
		sortable: true,
		render: (value: string) => (
			// variant={value === "Active" ? "success" : "secondary"}
			<Badge>{value}</Badge>
		),
	},
	{ key: "dueDate", label: "Due Date", sortable: true },
];

const initialData = [
	{
		id: 1,
		customerName: "John Doe",
		amount: 5000,
		status: "Active",
		dueDate: "2023-12-31",
	},
	{
		id: 2,
		customerName: "Jane Smith",
		amount: 7500,
		status: "Inactive",
		dueDate: "2023-11-30",
	},
	{
		id: 3,
		customerName: "Bob Johnson",
		amount: 10000,
		status: "Active",
		dueDate: "2024-01-15",
	},
];

export default function LoansPage() {
	const [data, setData] = useState(initialData);
	const [newLoan, setNewLoan] = useState({
		customerName: "",
		amount: "",
		status: "Active",
		dueDate: "",
	});
	const [selectedRow, setSelectedRow] = useState<
		(typeof initialData)[0] | null
	>(null);

	const handleAddLoan = () => {
		const newId = Math.max(...data.map((loan) => loan.id)) + 1;
		setData([
			...data,
			{ id: newId, ...newLoan, amount: Number(newLoan.amount) },
		]);
		setNewLoan({ customerName: "", amount: "", status: "Active", dueDate: "" });
	};

	return (
		<div className='p-6'>
			<div className='flex justify-between items-center mb-4'>
				<h1 className='text-3xl font-bold'>Loans</h1>
				<Dialog>
					<DialogTrigger asChild>
						<Button>Add New Loan</Button>
					</DialogTrigger>
					<DialogContent>
						<DialogHeader>
							<DialogTitle>Add New Loan</DialogTitle>
							<DialogDescription>
								Enter the details for the new loan.
							</DialogDescription>
						</DialogHeader>
						<div className='grid gap-4 py-4'>
							<div className='grid grid-cols-4 items-center gap-4'>
								<Label htmlFor='customerName' className='text-right'>
									Customer Name
								</Label>
								<Input
									id='customerName'
									value={newLoan.customerName}
									onChange={(e) =>
										setNewLoan({ ...newLoan, customerName: e.target.value })
									}
									className='col-span-3'
								/>
							</div>
							<div className='grid grid-cols-4 items-center gap-4'>
								<Label htmlFor='amount' className='text-right'>
									Amount
								</Label>
								<Input
									id='amount'
									type='number'
									value={newLoan.amount}
									onChange={(e) =>
										setNewLoan({ ...newLoan, amount: e.target.value })
									}
									className='col-span-3'
								/>
							</div>
							<div className='grid grid-cols-4 items-center gap-4'>
								<Label htmlFor='dueDate' className='text-right'>
									Due Date
								</Label>
								<Input
									id='dueDate'
									type='date'
									value={newLoan.dueDate}
									onChange={(e) =>
										setNewLoan({ ...newLoan, dueDate: e.target.value })
									}
									className='col-span-3'
								/>
							</div>
						</div>
						<DialogFooter>
							<Button onClick={handleAddLoan}>Add Loan</Button>
						</DialogFooter>
					</DialogContent>
				</Dialog>
			</div>
			<Table
				columns={columns}
				data={data}
				onRowClick={(row) => setSelectedRow(row)}
			/>
			<Sheet
				open={!!selectedRow}
				onOpenChange={(open) => !open && setSelectedRow(null)}
			>
				<SheetContent>
					<SheetHeader>
						<SheetTitle>Loan Details</SheetTitle>
					</SheetHeader>
					{selectedRow && (
						<div className='py-4'>
							{Object.entries(selectedRow).map(([key, value]) => (
								<div key={key} className='flex justify-between py-2'>
									<span className='font-medium'>{key}</span>
									<span>{value}</span>
								</div>
							))}
						</div>
					)}
				</SheetContent>
			</Sheet>
		</div>
	);
}
