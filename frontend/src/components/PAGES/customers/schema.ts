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
  overpayment: z.number(),
  assignedStaff: userSchema,
  dueAmount: z.number(),
  createdBy: userSchema,
  createdAt: z.string().date('Invalid date string!'),
});

export type Client = z.infer<typeof clientSchema>;

export const clientFormSchema = z.object({
  firstName: z.string(),
  lastName: z.string(),
  phoneNumber: z.string(),
  idNumber: z.string().optional(),
  dob: z.string().date().optional(),
  gender: z.enum(['MALE', 'FEMALE']),
  active: z.boolean(),
  branchId: z.number().gte(0, { message: 'Select valid branch' }),
  assignedStaffId: z.number().gte(0, { message: 'Select valid staff' }),
  // updatedBy
});

export type ClientFormType = z.infer<typeof clientFormSchema>;
