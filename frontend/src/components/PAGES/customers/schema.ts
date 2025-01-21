import { z } from 'zod';
import { userSchema } from '../users/schema';

export const clientSchema = z.object({
  id: z.number(),
  fullName: z.string(),
  phoneNumber: z.string(),
  idNumber: z.string().optional(),
  dob: z.string().date().optional(),
  gender: z.enum(['MALE', 'FEMALE']),
  active: z.boolean(),
  branchName: z.string(),
  assignedStaff: userSchema,
  overpayment: z.number(),
  dueAmount: z.number(),
  createdBy: userSchema,
  createdAt: z.string().date('Invalid date string!'),
});

export type Client = z.infer<typeof clientSchema>;
