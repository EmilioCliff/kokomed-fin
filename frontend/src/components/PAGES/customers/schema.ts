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
	dueAmount: z.number().optional(),
	createdBy: userSchema,
	createdAt: z.string().date('Invalid date string!'),
});

export type Client = z.infer<typeof clientSchema>;

export const clientFormSchema = z.object({
	firstName: z.string().min(3),
	lastName: z.string().min(3),
	phoneNumber: z
		.string()
		.length(10, 'Phone number must be exactly 10 digits')
		.regex(/^\d+$/, 'Phone number must contain only digits'),
	idNumber: z.string().optional(),
	// dob: z.string().date().optional(),
	dob: z.string().optional(),
	gender: z.enum(['MALE', 'FEMALE']),
	branchId: z.number().gte(0, { message: 'Select valid branch' }),
	assignedStaffId: z.number().gte(0, { message: 'Select valid staff' }),
});

export type ClientFormType = z.infer<typeof clientFormSchema>;
