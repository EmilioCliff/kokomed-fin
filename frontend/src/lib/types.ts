import { ReactNode } from 'react';
import {
	Loan,
	ExpectedPayment,
	Installment,
} from '@/components/PAGES/loans/schema';
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

export enum ReportType {
	LOANS = 'loans',
	USERS = 'users',
	CLIENTS = 'clients',
	PRODUCTS = 'products',
	BRANCHES = 'branches',
	PAYMENTS = 'payments',
}

export enum ReportTag {
	LOANS = 'Loan Reports',
	USERS = 'User Reports',
	CLIENTS = 'Client Reports',
	PRODUCTS = 'Product Reports',
	BRANCHES = 'Branch Reports',
	PAYMENTS = 'Payments Reports',
}

export type ReportFilter = {
	reportName: string;
	startDate?: string;
	endDate?: string;
	clientId?: string;
	userId?: string;
	loanId?: string;
	format: string;
};

export type Report = {
	id: string;
	title: string;
	description: string;
	type: ReportType;
	lastGenerated: string;
	tag: ReportTag;
	filters: ReportFilter;
};

export type LoanEvent = {
	id: string;
	clientName: string;
	loanAmount: number;
	date: string;
	paymentDue?: number;
	type: string;
	allDay: boolean;
	title: string;
};

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

export interface getExpectedPaymentsType extends Omit<commonresponse, 'data'> {
	data: ExpectedPayment[];
}

export interface getLoanInstallmentsType extends Omit<commonresponse, 'data'> {
	data: Installment[];
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

export interface getLoanEventsType extends Omit<commonresponse, 'data'> {
	data: LoanEvent[];
}

export interface getFormDataType extends Omit<commonresponse, 'data'> {
	product?: commonDataResponse[];
	client?: commonDataResponse[];
	user?: commonDataResponse[];
	branch?: commonDataResponse[];
	loan?: commonDataResponse[];
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
