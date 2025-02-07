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
});

export type Payment = z.infer<typeof paymentSchema>;

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
	App: z.string(),
});

export type PaymentFormType = z.infer<typeof paymentFormSchema>;
