import { z } from 'zod';

export const loginFormSchema = z.object({
	email: z.string().email({ message: 'Invalid email address' }),
	password: z
		.string()
		.min(3, { message: 'Must be 3 or more characters long' }),
});

export type LoginForm = z.infer<typeof loginFormSchema>;

export const resetPasswordFormSchema = z.object({
	token: z.string(),
	newPassword: z
		.string()
		.min(3, { message: 'Must be 3 or more characters long' }),
	confirmPassword: z
		.string()
		.min(3, { message: 'Must be 3 or more characters long' }),
});

export type ResetPassowordFormType = z.infer<typeof resetPasswordFormSchema>;
