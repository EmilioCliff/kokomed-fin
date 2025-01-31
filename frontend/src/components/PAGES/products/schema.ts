import { z } from 'zod';

export const productSchema = z.object({
	id: z.number(),
	branchName: z.string(),
	loanAmount: z.number(),
	repayAmount: z.number(),
	interestAmount: z.number(),
});

export type Product = z.infer<typeof productSchema>;

export const productFormSchema = z.object({
	branchId: z.number(),
	loanAmount: z.coerce.number().min(0),
	repayAmount: z.coerce.number().min(0),
});

export type ProductFormType = z.infer<typeof productFormSchema>;
