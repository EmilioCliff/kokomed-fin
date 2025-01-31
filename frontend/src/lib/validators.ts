import { User, userSchema } from '@/components/PAGES/users/schema';
import { Client, clientSchema } from '@/components/PAGES/customers/schema';

export function isUser(value: unknown): User {
	return userSchema.parse(value);
}

export function isClient(value: unknown): Client {
	return clientSchema.parse(value);
}
