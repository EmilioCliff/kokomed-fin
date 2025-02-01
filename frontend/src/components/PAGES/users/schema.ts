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

export const userFormSchema = z.object({
	firstName: z.string(),
	lastName: z.string(),
	phoneNumber: z
		.string()
		.length(10, 'Phone number must be exactly 10 digits')
		.regex(/^\d+$/, 'Phone number must contain only digits'),
	email: z.string().email({ message: 'Invalid email address' }),
	branchId: z.number().gt(0),
	role: z.enum(['ADMIN', 'AGENT']),
});

export type UserFormType = z.infer<typeof userFormSchema>;
