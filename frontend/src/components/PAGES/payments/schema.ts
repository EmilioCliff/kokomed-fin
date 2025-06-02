import { z } from 'zod';
import { clientSchema } from '../customers/schema';

export const paymentSchema = z.object({
	id: z.number(),
	transactionSource: z.enum(['MPESA', 'INTERNAL']),
	transactionNumber: z.string(),
	accountNumber: z.string(),
	phoneNumber: z.string(),
	payingName: z.string(),
	amount: z.number(),
	paidDate: z.string().date(),
	assigned: z.boolean(),
	assignedTo: clientSchema,
	assignedBy: z.string(),
});

export type Payment = z.infer<typeof paymentSchema>;

export const paymentAllocationSchema = z.object({
	id: z.number(),
	nonPostedId: z.number(),
	loanId: z.number().optional(),
	installmentId: z.number().optional(),
	amount: z.number(),
	description: z.string(),
	deletedDescription: z.string().optional(),
	deletedAt: z.string().optional(),
	createdAt: z.string(),
});

export const paymentFormSchema = z.object({
	TransAmount: z.coerce.number().gt(0, { message: 'Select valid amount' }),
	TransID: z
		.string()
		.min(3, { message: 'Must be 3 or more characters long' }),
	BillRefNumber: z
		.string()
		.min(10, { message: 'Must be 10 characters long' }),
	MSISDN: z.string().min(10, { message: 'Must be 10 characters long' }),
	FirstName: z
		.string()
		.min(3, { message: 'Must be 3 or more characters long' }),
	DatePaid: z.string(),
	App: z.string(),
	Email: z.string().email(),
});

export type PaymentFormType = z.infer<typeof paymentFormSchema>;

export const editClientFormSchema = z.object({
	id: z.number().gt(0, { message: 'Select valid client' }),
	fullName: z.string().min(3),
	phoneNumber: z
		.string()
		.length(10, 'Phone number must be exactly 10 digits')
		.regex(/^\d+$/, 'Phone number must contain only digits'),
	idNumber: z.string().optional(),
	dob: z.string().optional(),
	gender: z.enum(['MALE', 'FEMALE']),
	active: z.enum(['true', 'false']),
	branchId: z.number().gt(0, { message: 'Select valid branch' }),
	assignedStaffId: z.number().gt(0, { message: 'Select valid staff' }),
});

export type EditClientFormType = z.infer<typeof editClientFormSchema>;

export const editPaymentFormSchema = z.object({
	id: z.number(),
	transactionSource: z.enum(['MPESA', 'INTERNAL']),
	transactionId: z
		.string()
		.min(3, { message: 'Must be 3 or more characters long' }),
	accountNumber: z
		.string()
		.min(10, { message: 'Must be 10 characters long' }),
	phoneNumber: z.string().min(3, { message: 'Must be 3 characters long' }),
	payingName: z
		.string()
		.min(3, { message: 'Must be 3 or more characters long' }),
	amount: z.coerce.number().gt(0, { message: 'Select valid amount' }),
	assignedBy: z.string(),
	assignedTo: z.number(),
	description: z
		.string()
		.min(3, { message: 'Must be 3 or more characters long' }),
	paidDate: z.string(),
});

export type EditPaymentFormType = z.infer<typeof editPaymentFormSchema>;
