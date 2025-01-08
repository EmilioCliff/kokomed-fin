import { z } from "zod";
import { generateRandomLoans } from "@/lib/generator";
import { useState, useEffect } from "react";
import { cn } from "@/lib/utils";
import { format } from "date-fns";
import LoanForm from "@/components/forms/LoanForm";
import { CalendarIcon } from "lucide-react";
import { loanColumns } from "../components/table/columns/loan";
import { DataTable } from "../components/table/data-table";
import { loanSchema, Loan } from "../data/schema";
import { isUser, isClient } from "@/lib/validators";
import UserCardDisplay from "@/components/UserCardDisplay";
import ClientCardDisplay from "@/components/ClientCardDisplay";
import TableSkeleton from "@/components/TableSkeleton";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
} from "@/components/ui/sheet";
import { Calendar } from "@/components/ui/calendar";
import { statuses } from "@/data/loan";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@/components/ui/popover";
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectLabel,
	SelectTrigger,
	SelectValue,
	SelectSeparator,
} from "@/components/ui/select";
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

const isAdmin = true;

const generatedLoans = generateRandomLoans(30);
const validatedLoans = z.array(loanSchema).parse(generatedLoans);

export default function LoanPage() {
	const [loans, setLoans] = useState<Loan[]>([]);
	const [loading, setLoading] = useState(true);
	const [selectedRow, setSelectedRow] = useState<Loan | null>(null);
	const [loanStatus, setLoanStatus] = useState<string | null>(null);
	const [disbursedDate, setDisbursedDate] = useState<string | null>(null);
	const [feePaid, setFeePaid] = useState<boolean | null>(null);

	useEffect(() => {
		setLoans(validatedLoans);
		setLoading(false);
	}, [validatedLoans]);

	console.log("Here")

	const fieldRenderers: Record<string, (value: any) => JSX.Element> = {
		id: (value) => (
			<Input
				readOnly
				placeholder={`LN${String(value).padStart(3, "0")}`}
				className='bg-gray-100 text-gray-500'
			/>
		),
		client: (value) => <ClientCardDisplay client={isClient(value)} />,
		loanOfficer: (value) => <UserCardDisplay user={isUser(value)} />,
		loanPurpose: (value) => <>{value}</>,
		approvedBy: (value) => <UserCardDisplay user={isUser(value)} />,
		disbursedBy: (value) => <UserCardDisplay user={isUser(value)} />,
		status: (value: string) => {
			return value === "INACTIVE" && isAdmin ? (
				<div>
					<Select
						defaultValue='INACTIVE'
						onValueChange={(value) => setLoanStatus(value)}
					>
						<SelectTrigger className='w-[180px]'>
							<SelectValue placeholder='INACTIVE' />
						</SelectTrigger>
						<SelectContent>
							<SelectGroup>
								<SelectItem value='INACTIVE'>INACTIVE</SelectItem>
								<SelectItem value='ACTIVE'>ACTIVE</SelectItem>
								<SelectSeparator />
								<SelectLabel>Disbursed Date</SelectLabel>
								<Popover>
									<PopoverTrigger asChild>
										<Button
											variant={"outline"}
											className={cn(
												"w-[240px] pl-3 text-left font-normal",
												!disbursedDate && "text-muted-foreground"
											)}
										>
											{disbursedDate ? (
												format(disbursedDate, "PPP")
											) : (
												<span>Pick a date</span>
											)}
											<CalendarIcon className='ml-auto h-4 w-4 opacity-50' />
										</Button>
									</PopoverTrigger>
									<PopoverContent className='w-auto p-0' align='start'>
										<Calendar
											mode='single'
											selected={
												disbursedDate ? new Date(disbursedDate) : undefined
											}
											onSelect={(date) =>
												setDisbursedDate(format(date!, "yyyy-MM-dd"))
											}
											disabled={(date) =>
												date > new Date() || date < new Date("1900-01-01")
											}
											initialFocus
										/>
									</PopoverContent>
								</Popover>
							</SelectGroup>
						</SelectContent>
					</Select>
					{disbursedDate && (
						<p>Disbursed Date: {format(new Date(disbursedDate), "PPP")}</p>
					)}
				</div>
			) : (
				<Input
					readOnly
					placeholder={value.toString()}
					className='bg-gray-100 text-gray-500'
				/>
			);
		},
		feePaid: (value) => {
			return value === false ? (
				<Select
					defaultValue='NO'
					onValueChange={(value) => setFeePaid(value === "YES")}
				>
					<SelectTrigger className='w-[180px]'>
						<SelectValue placeholder='NO' />
					</SelectTrigger>
					<SelectContent>
						<SelectGroup>
							<SelectItem value='YES'>YES</SelectItem>
							<SelectItem value='NO'>NO</SelectItem>
						</SelectGroup>
					</SelectContent>
				</Select>
			) : (
				<Input
					readOnly
					placeholder='YES'
					className='bg-gray-100 text-gray-500'
				/>
			);
		},
		updatedBy: () => <></>,
		createdBy: () => <></>,
	};

	if (loading) {
		return <TableSkeleton />;
	}

	return (
		<div className='px-4'>
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
				searchableColumns={[
					{
						id: "clientName",
						title: "Client Name",
					},
					{
						id: "loanOfficerName",
						title: "Loan Officer",
					},
				]}
				facetedFilterColumns={[
					{
						id: "status",
						title: "Status",
						options: statuses,
					},
				]}
			/>
			<Sheet
				open={!!selectedRow}
				onOpenChange={(open: boolean) => {
					if (!open) {
						setSelectedRow(null);
						setDisbursedDate(null);
						setFeePaid(null);
					}
				}}
			>
				<SheetContent className='overflow-auto custom-sheet-class'>
					<SheetHeader>
						<SheetTitle>Loan Details</SheetTitle>
						<SheetDescription>Description goes here</SheetDescription>
					</SheetHeader>
					{selectedRow && (
						<div className='py-4'>
							{Object.entries(selectedRow).map(([key, value]) => {
								if (key === "createdBy" || key === "updatedBy") {
									return;
								}
								if (fieldRenderers[key]) {
									return (
										<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
											<span className='font-medium capitalize'>{key}</span>
											{fieldRenderers[key](value)}
										</div>
									);
								}

								return (
									<div key={key} className='grid grid-cols-[0.5fr_1fr] mb-4'>
										<span className='font-medium capitalize'>{key}</span>
										{typeof value === "string" ||
										typeof value === "number" ||
										typeof value === "boolean" ? (
											<Input
												readOnly
												placeholder={value.toString()}
												className='bg-gray-100 text-gray-500'
											/>
										) : (
											JSON.stringify(value)
										)}
									</div>
								);
							})}
							<Button
								size='lg'
								onClick={() =>
									console.log(loanStatus, "/n", disbursedDate, "/n", feePaid)
								}
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
