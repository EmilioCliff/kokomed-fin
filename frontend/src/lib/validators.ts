import { User, Client, userSchema, clientSchema } from "@/data/schema";

export function isUser(value: unknown): User {
	return userSchema.parse(value);
}

export function isClient(value: unknown): Client {
	return clientSchema.parse(value);
}
