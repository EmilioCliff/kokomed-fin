import { z } from 'zod';

export const productSchema = z.object({
  id: z.number(),
  branchName: z.string(),
  loanAmount: z.number(),
  repayAMount: z.number(),
  interestAmount: z.number(),
});

export type Product = z.infer<typeof productSchema>;
