import { z } from 'zod';
import { clientSchema } from '../customers/schema';
import { userSchema } from '../users/schema';

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
  status: z.enum(['INACTIVE', 'ACTIVE', 'COMPLETED', 'DEFAULTED']),
  processingFee: z.number(),
  feePaid: z.boolean(),
  paidAmount: z.number(),
  updatedBy: userSchema,
  createdBy: userSchema,
  createdAt: z.string().date('Invalid date string!'),
});

export type Loan = z.infer<typeof loanSchema>;

// approvedBy: z.number(), add approved by when form is submitted
// disburseBy: z.number().optional(), add disbursed by when form is submitted
export const loanFormSchema = z.object({
  productId: z.number().gte(0, { message: 'Select valid product' }),
  clientId: z.number().gte(0, { message: 'Select valid client' }),
  loanOfficerId: z.number().gte(0, { message: 'Select valid loan officer' }),
  loanPurpose: z.string(),
  disburse: z.boolean(),
  disburseOn: z
    .string()
    .optional()
    .refine((dateString) => !dateString || !isNaN(Date.parse(dateString)), {
      message: 'Invalid date string!',
    }),
  noOfInstallments: z.coerce.number().gt(0),
  installmentsPeriod: z.coerce.number().gt(0),
  processingFee: z.coerce.number().gt(0),
  processingFeePaid: z.boolean(),
  dob: z.string().date().optional(),
});

export type LoanFormType = z.infer<typeof loanFormSchema>;
