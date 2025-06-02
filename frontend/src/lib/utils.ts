import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// export function formatDate(input: string | number): string {
// 	const date = new Date(input);
// 	return date.toLocaleDateString("en-US", {
// 		month: "long",
// 		day: "numeric",
// 		year: "numeric",
// 	});
// }

export function formatCurrency(amount: number): string {
	return new Intl.NumberFormat('en-KE', {
		style: 'currency',
		currency: 'KES',
		minimumFractionDigits: 0,
		maximumFractionDigits: 0,
	}).format(amount);
}

export function formatDate(dateString: string): string {
	if (!dateString || dateString === '0001-01-01') return '-';

	const date = new Date(dateString);
	if (isNaN(date.getTime())) return '-';

	return new Intl.DateTimeFormat('en-KE', {
		month: 'short',
		day: 'numeric',
		year: 'numeric',
	}).format(date);
}
