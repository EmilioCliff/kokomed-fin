import { z } from 'zod';
import { clientSchema } from '../customers/schema';

export const paymentSchema = z.object({
  id: z.number(),
  transactionSource: z.string(),
  transactionNumber: z.string(),
  accountNumber: z.string(),
  phoneNumber: z.string(),
  payingName: z.string(),
  amount: z.number(),
  paidDate: z.string().date(),
  assignedTo: clientSchema,
});

export type Payment = z.infer<typeof paymentSchema>;
