import { z } from 'zod';
import { clientSchema } from '../customers/schema';
import { userSchema } from '../users/schema';

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
      message: 'Invalid date string!',
    }),
});

export type InactiveLoan = z.infer<typeof inactiveLoanSchema>;
