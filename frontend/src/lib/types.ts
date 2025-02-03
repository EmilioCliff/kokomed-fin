import { ReactNode } from 'react';
import { Loan } from '@/components/PAGES/loans/schema';
import { User } from '@/components/PAGES/users/schema';
import { Product } from '@/components/PAGES/products/schema';
import { Branch } from '@/components/PAGES/branches/schema';
import { Client } from '@/components/PAGES/customers/schema';
import { Payment } from '@/components/PAGES/payments/schema';

export enum role {
	AGENT = 'AGENT',
	ADMIN = 'ADMIN',
	GUEST = 'GUEST',
}

export enum loanStatus {
	ACTIVE = 'ACTIVE',
	INACTIVE = 'INACTIVE',
	COMPLETED = 'COMPLETED',
	DEFAULTED = 'DEFAULTED',
}

export interface tableFilterType {
	options: {
		label: string;
		value: string;
	}[];
}

export interface contextWrapperProps {
	children: ReactNode;
}

export interface pagination {
	pageSize: number;
	currentPage: number;
	totalData: number;
	totalPages: number;
}

export interface commonresponse {
	statusCode?: string;
	message?: string;
	metadata?: pagination;
	data: any;
}

export interface userResponse {
	id: number;
	fullname: string;
	email: string;
	phoneNumber: string;
	role: role;
	branchName: string;
	createdAt: string;
}

export interface tokenData extends Omit<commonresponse, 'data'> {
	accessToken: string;
}

export interface refreshTokenRes extends Omit<commonresponse, 'data'> {
	accessToken: string;
}

export interface getUsersType extends Omit<commonresponse, 'data'> {
	data: User[];
}

export interface getClientsType extends Omit<commonresponse, 'data'> {
	data: Client[];
}

export interface getLoansType extends Omit<commonresponse, 'data'> {
	data: Loan[];
}

export interface getProductsType extends Omit<commonresponse, 'data'> {
	data: Product[];
}

export interface getBranchesType extends Omit<commonresponse, 'data'> {
	data: Branch[];
}

export interface getPaymentsType extends Omit<commonresponse, 'data'> {
	data: Payment[];
}

export interface getFormDataType extends Omit<commonresponse, 'data'> {
	product?: commonDataResponse[];
	client?: commonDataResponse[];
	user?: commonDataResponse[];
	branch?: commonDataResponse[];
}

export interface commonDataResponse {
	id: number;
	name: string;
}

export interface updateLoanType extends Omit<commonresponse, 'data'> {
	id: number;
	status?: loanStatus;
	disburseDate?: string;
	feePaid?: boolean;
}

export interface updateUserType extends Omit<commonresponse, 'data'> {
	id: number;
	role?: role;
	branchId?: number;
}

export interface updateClientType extends Omit<commonresponse, 'data'> {
	id: number;
	idNumber?: string;
	dob?: string;
	branchId?: number;
	active?: string;
}

export interface updatePaymentType extends Omit<commonresponse, 'data'> {
	id: number;
	clientId: number;
}
