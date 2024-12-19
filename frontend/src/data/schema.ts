import { z } from "zod";

export const userSchema = z.object({
	id: z.number(),
	fullName: z.string(),
	phoneNumber: z.string(),
	email: z.string().email(),
	role: z.enum(["ADMIN", "AGENT"]),
	branchName: z.string(),
	createdAt: z.string().date("Invalid date string!"),
});

export const clientSchema = z.object({
	id: z.number(),
	fullName: z.string(),
	phoneNumber: z.string(),
	idNumber: z.string().optional(),
	dob: z.string().date().optional(),
	gender: z.enum(["MALE", "FEMALE"]),
	active: z.boolean(),
	branchName: z.string(),
	assignedStaff: userSchema,
	overpayment: z.number(),
	dueAmount: z.number(),
	createdBy: userSchema,
	createdAt: z.string().date("Invalid date string!"),
});

export const branchSchema = z.object({
	id: z.number(),
	branchName: z.string(),
});

export const productSchema = z.object({
	id: z.number(),
	branchName: z.string(),
	loanAmount: z.number(),
	repayAMount: z.number(),
	interestAmount: z.number(),
});

export const loanSchema = z.object({
	id: z.number(),
	amount: z.number(),
	repayAmount: z.number(),
	client: clientSchema,
	loanOfficer: userSchema,
	loanPurpose: z.string().optional(),
	dueDate: z.string().date().optional(),
	approvedBy: userSchema,
	disbursedOn: z.string().date().optional(),
	disbursedBy: userSchema.optional(),
	noOfInstallments: z.number(),
	installmentsPeriod: z.number(),
	status: z.enum(["INACTIVE", "ACTIVE", "COMPLETED", "DEFAULTED"]),
	processingFee: z.number(),
	feePaid: z.boolean(),
	paidAmount: z.number(),
	updatedBy: userSchema,
	createdBy: userSchema,
	createdAt: z.string().date("Invalid date string!"),
});

export const inactiveLoanSchema = z.object({
	id: z.number(),
	amount: z.number(),
	repayAmount: z.number(),
	client: clientSchema,
	loanOfficer: userSchema,
	approvedBy: userSchema,
	approvedOn: z
		.string()
		.optional()
		.refine((dateString) => !dateString || !isNaN(Date.parse(dateString)), {
			message: "Invalid date string!",
		}),
});

export type User = z.infer<typeof userSchema>;
export type Client = z.infer<typeof clientSchema>;
export type Branch = z.infer<typeof branchSchema>;
export type Product = z.infer<typeof productSchema>;
export type Loan = z.infer<typeof loanSchema>;
export type InactiveLoan = z.infer<typeof inactiveLoanSchema>;
