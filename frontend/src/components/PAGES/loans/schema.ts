import { z } from 'zod';
import { clientSchema } from '../customers/schema';
import { userSchema } from '../users/schema';
import { productSchema } from '../products/schema';

export const loanSchema = z.object({
	id: z.number(),
	product: productSchema,
	client: clientSchema,
	loanOfficer: userSchema,
	loanPurpose: z.string().optional(),
	dueDate: z.string().date().optional(),
	approvedBy: userSchema,
	disbursedOn: z.string().date().optional(),
	disbursedBy: userSchema.optional(),
	noOfInstallments: z.number(),
	installmentsPeriod: z.number(),
	status: z.enum(['INACTIVE', 'ACTIVE', 'COMPLETED', 'DEFAULTED']),
	processingFee: z.number(),
	feePaid: z.boolean(),
	paidAmount: z.number(),
	updatedBy: userSchema,
	createdBy: userSchema,
	createdAt: z.string().date('Invalid date string!'),
});

export type Loan = z.infer<typeof loanSchema>;

export const loanFormSchema = z.object({
	productId: z.number().gt(0, { message: 'Select valid product' }),
	clientId: z.number().gt(0, { message: 'Select valid client' }),
	loanOfficerId: z.number().gt(0, { message: 'Select valid loan officer' }),
	loanPurpose: z.string().optional(),
	disburse: z.boolean(),
	disburseOn: z
		.string()
		.optional()
		.refine((dateString) => !dateString || !isNaN(Date.parse(dateString)), {
			message: 'Invalid date string!',
		}),
	installments: z.coerce.number().gt(0),
	installmentsPeriod: z.coerce.number().gt(0),
	processingFee: z.coerce.number().gt(0),
	processingFeePaid: z.boolean(),
});

export type LoanFormType = z.infer<typeof loanFormSchema>;

export const expectedPaymentsSchema = z.object({
	loanId: z.number(),
	branchName: z.string(),
	clientName: z.string(),
	loanOfficerName: z.string(),
	loanAmount: z.number(),
	repayAmount: z.number(),
	totalUnpaid: z.number(),
	dueDate: z.string(),
});

export type ExpectedPayment = z.infer<typeof expectedPaymentsSchema>;

export const installmentSchema = z.object({
	id: z.number(),
	loanId: z.number(),
	installmentNo: z.number(),
	amount: z.number(),
	remainingAmount: z.number(),
	paid: z.boolean(),
	paidAt: z.string(),
	dueDate: z.string(),
});

export type Installment = z.infer<typeof installmentSchema>;

export const loanShortSchema = z.object({
	id: z.number(),
	loanAmount: z.number(),
	repayAmount: z.number(),
	disbursedOn: z.string(),
	dueDate: z.string(),
	paidAmount: z.number(),
	installments: z.array(installmentSchema),
});

export type LoanShort = z.infer<typeof loanShortSchema>;

export const unpaidInstallmentSchema = z.object({
	installmentNumber: z.number(),
	amountDue: z.number(),
	remainingAmount: z.number(),
	dueDate: z.string(),
	loanId: z.number(),
	loanAmount: z.number(),
	repayAmount: z.number(),
	paidAmount: z.number(),
	clientId: z.number(),
	fullName: z.string(),
	phoneNumber: z.string(),
	branchName: z.string(),
	totalDueAmount: z.number(),
});

export type UnpaidInstallment = z.infer<typeof unpaidInstallmentSchema>;
