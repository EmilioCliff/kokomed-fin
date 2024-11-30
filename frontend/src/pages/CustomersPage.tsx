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

const columns = [
	{ key: "id", label: "ID", sortable: true },
	{ key: "name", label: "Name", sortable: true },
	{ key: "email", label: "Email", sortable: true },
	{ key: "phone", label: "Phone", sortable: true },
	{ key: "totalLoans", label: "Total Loans", sortable: true },
];

const initialData = [
	{
		id: 1,
		name: "John Doe",
		email: "john@example.com",
		phone: "123-456-7890",
		totalLoans: 2,
	},
	{
		id: 2,
		name: "Jane Smith",
		email: "jane@example.com",
		phone: "098-765-4321",
		totalLoans: 1,
	},
	{
		id: 3,
		name: "Bob Johnson",
		email: "bob@example.com",
		phone: "555-555-5555",
		totalLoans: 3,
	},
];

export default function CustomersPage() {
	const [data, setData] = useState(initialData);
	const [newCustomer, setNewCustomer] = useState({
		name: "",
		email: "",
		phone: "",
	});
	const [selectedCustomer, setSelectedCustomer] = useState<
		(typeof initialData)[0] | null
	>(null);

	const handleAddCustomer = () => {
		const newId = Math.max(...data.map((customer) => customer.id)) + 1;
		setData([...data, { id: newId, ...newCustomer, totalLoans: 0 }]);
		setNewCustomer({ name: "", email: "", phone: "" });
	};

	return (
		<div className='p-6'>
			<div className='flex justify-between items-center mb-4'>
				<h1 className='text-3xl font-bold'>Customers</h1>
				<Dialog>
					<DialogTrigger asChild>
						<Button>Add New Customer</Button>
					</DialogTrigger>
					<DialogContent>
						<DialogHeader>
							<DialogTitle>Add New Customer</DialogTitle>
							<DialogDescription>
								Enter the details for the new customer.
							</DialogDescription>
						</DialogHeader>
						<div className='grid gap-4 py-4'>
							<div className='grid grid-cols-4 items-center gap-4'>
								<Label htmlFor='name' className='text-right'>
									Name
								</Label>
								<Input
									id='name'
									value={newCustomer.name}
									onChange={(e) =>
										setNewCustomer({ ...newCustomer, name: e.target.value })
									}
									className='col-span-3'
								/>
							</div>
							<div className='grid grid-cols-4 items-center gap-4'>
								<Label htmlFor='email' className='text-right'>
									Email
								</Label>
								<Input
									id='email'
									type='email'
									value={newCustomer.email}
									onChange={(e) =>
										setNewCustomer({ ...newCustomer, email: e.target.value })
									}
									className='col-span-3'
								/>
							</div>
							<div className='grid grid-cols-4 items-center gap-4'>
								<Label htmlFor='phone' className='text-right'>
									Phone
								</Label>
								<Input
									id='phone'
									type='tel'
									value={newCustomer.phone}
									onChange={(e) =>
										setNewCustomer({ ...newCustomer, phone: e.target.value })
									}
									className='col-span-3'
								/>
							</div>
						</div>
						<DialogFooter>
							<Button onClick={handleAddCustomer}>Add Customer</Button>
						</DialogFooter>
					</DialogContent>
				</Dialog>
			</div>
			<Table
				columns={columns}
				data={data}
				onRowClick={(row) => setSelectedCustomer(row)}
			/>
			<Sheet
				open={!!selectedCustomer}
				onOpenChange={(open) => !open && setSelectedCustomer(null)}
			>
				<SheetContent>
					<SheetHeader>
						<SheetTitle>Customer Details</SheetTitle>
					</SheetHeader>
					{selectedCustomer && (
						<div className='py-4'>
							{Object.entries(selectedCustomer).map(([key, value]) => (
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
