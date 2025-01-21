import { z } from 'zod';

export const branchSchema = z.object({
  id: z.number(),
  branchName: z.string(),
});

export type Branch = z.infer<typeof branchSchema>;
