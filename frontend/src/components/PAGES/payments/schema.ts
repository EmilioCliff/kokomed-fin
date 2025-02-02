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
	TransAmount: z.coerce.number(),
	TransID: z.string(),
	BillRefNumber: z.string(),
	MSISDN: z.string(),
	FirstName: z.string(),
	App: z.string(),
});

export type PaymentFormType = z.infer<typeof paymentFormSchema>;
