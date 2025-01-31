import { z } from 'zod';

export const branchSchema = z.object({
	id: z.number(),
	name: z.string(),
});

export type Branch = z.infer<typeof branchSchema>;

export const branchFormSchema = z.object({
	name: z.string(),
});

export type BranchFormType = z.infer<typeof branchFormSchema>;
