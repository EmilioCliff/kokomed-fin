import { useState } from "react";
import { Table } from "@/components/Table";
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
import {
	Sheet,
	SheetContent,
	SheetHeader,
	SheetTitle,
} from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";

const columns = [
	{ key: "id", label: "ID", sortable: true },
	{ key: "customerName", label: "Customer Name", sortable: true },
	{ key: "amount", label: "Amount", sortable: true },
	{ key: "date", label: "Date", sortable: true },
	{
		key: "status",
		label: "Status",
		sortable: true,
		render: (value: string) => (
			// variant={value === "Completed" ? "success" : "warning"}
			<Badge>{value}</Badge>
		),
	},
];

const initialData = [
	{
		id: 1,
		customerName: "John Doe",
		amount: 500,
		date: "2023-05-01",
		status: "Completed",
	},
	{
		id: 2,
		customerName: "Jane Smith",
		amount: 750,
		date: "2023-05-02",
		status: "Pending",
	},
	{
		id: 3,
		customerName: "Bob Johnson",
		amount: 1000,
		date: "2023-05-03",
		status: "Completed",
	},
];

export default function PaymentsPage() {
	const [data, setData] = useState(initialData);
	const [newPayment, setNewPayment] = useState({
		customerName: "",
		amount: "",
		date: "",
		status: "Pending",
	});
	const [selectedPayment, setSelectedPayment] = useState<
		(typeof initialData)[0] | null
	>(null);

	const handleAddPayment = () => {
		const newId = Math.max(...data.map((payment) => payment.id)) + 1;
		setData([
			...data,
			{ id: newId, ...newPayment, amount: Number(newPayment.amount) },
		]);
		setNewPayment({
			customerName: "",
			amount: "",
			date: "",
			status: "Pending",
		});
	};

	return (
		<div className='p-6'>
			<div className='flex justify-between items-center mb-4'>
				<h1 className='text-3xl font-bold'>Payments</h1>
				<Dialog>
					<DialogTrigger asChild>
						<Button>Add New Payment</Button>
					</DialogTrigger>
					<DialogContent>
						<DialogHeader>
							<DialogTitle>Add New Payment</DialogTitle>
							<DialogDescription>
								Enter the details for the new payment.
							</DialogDescription>
						</DialogHeader>
						<div className='grid gap-4 py-4'>
							<div className='grid grid-cols-4 items-center gap-4'>
								<Label htmlFor='customerName' className='text-right'>
									Customer Name
								</Label>
								<Input
									id='customerName'
									value={newPayment.customerName}
									onChange={(e) =>
										setNewPayment({
											...newPayment,
											customerName: e.target.value,
										})
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
									value={newPayment.amount}
									onChange={(e) =>
										setNewPayment({ ...newPayment, amount: e.target.value })
									}
									className='col-span-3'
								/>
							</div>
							<div className='grid grid-cols-4 items-center gap-4'>
								<Label htmlFor='date' className='text-right'>
									Date
								</Label>
								<Input
									id='date'
									type='date'
									value={newPayment.date}
									onChange={(e) =>
										setNewPayment({ ...newPayment, date: e.target.value })
									}
									className='col-span-3'
								/>
							</div>
						</div>
						<DialogFooter>
							<Button onClick={handleAddPayment}>Add Payment</Button>
						</DialogFooter>
					</DialogContent>
				</Dialog>
			</div>
			<Table
				columns={columns}
				data={data}
				onRowClick={(row) => setSelectedPayment(row)}
			/>
			<Sheet
				open={!!selectedPayment}
				onOpenChange={(open) => !open && setSelectedPayment(null)}
			>
				<SheetContent>
					<SheetHeader>
						<SheetTitle>Payment Details</SheetTitle>
					</SheetHeader>
					{selectedPayment && (
						<div className='py-4'>
							{Object.entries(selectedPayment).map(([key, value]) => (
								<div key={key} className='flex justify-between py-2'>
									<span className='font-medium'>{key}</span>
									<span>
										{typeof value === "string" ? value : JSON.stringify(value)}
									</span>
								</div>
							))}
						</div>
					)}
				</SheetContent>
			</Sheet>
		</div>
	);
}
