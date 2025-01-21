import { z } from 'zod';

export const userSchema = z.object({
  id: z.number(),
  fullName: z.string(),
  phoneNumber: z.string(),
  email: z.string().email(),
  role: z.enum(['ADMIN', 'AGENT']),
  branchName: z.string(),
  createdAt: z.string().date('Invalid date string!'),
});

export type User = z.infer<typeof userSchema>;
