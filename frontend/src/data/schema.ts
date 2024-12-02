import { z } from "zod";

export const userSchema = z.object({
	id: z.number(),
	fullName: z.string(),
	phoneNumber: z.string(),
	email: z.string().email(),
	role: z.enum(["ADMIN", "AGENT"]),
});

export const clientSchema = z.object({
	id: z.number(),
	fullName: z.string(),
	phoneNumber: z.string(),
	active: z.boolean(),
	assignedStaff: userSchema,
	overpayment: z.number(),
});

export const loanSchema = z.object({
	id: z.number(),
	amount: z.number(),
	repayAmount: z.number(),
	client: clientSchema,
	loanOfficer: userSchema,
	loanPurpose: z.string().optional(),
	dueDate: z
		.string()
		.optional()
		.refine((dateString) => !dateString || !isNaN(Date.parse(dateString)), {
			message: "Invalid date string!",
		}),
	approvedBy: userSchema,
	disbursedOn: z
		.string()
		.optional()
		.refine((dateString) => !dateString || !isNaN(Date.parse(dateString)), {
			message: "Invalid date string!",
		}),
	disbursedBy: userSchema,
	noOfInstallments: z.number(),
	installmentsPeriod: z.number(),
	status: z.enum(["INACTIVE", "ACTIVE", "COMPLETED", "DEFAULTED"]),
	processingFee: z.number(),
	feePaid: z.boolean(),
	paidAmount: z.number(),
	updatedBy: userSchema,
	createdBy: userSchema,
	createdAt: z
		.string()
		.optional()
		.refine((dateString) => !dateString || !isNaN(Date.parse(dateString)), {
			message: "Invalid date string!",
		}),
});

export type User = z.infer<typeof userSchema>;
export type Client = z.infer<typeof clientSchema>;
export type Loan = z.infer<typeof loanSchema>;
