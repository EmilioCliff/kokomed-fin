import { ReactNode } from 'react';
import { User, Loan } from '@/data/schema';

export interface authCtx {
  isLoading: boolean;
  isAuthenticated: boolean;
  error: any;
  // add login, logout and checkSession methods
}

export interface contextWrapperProps {
  children: ReactNode;
}

export interface tableMetaDataType {
  totalRows: number;
  rowsPerPage: number;
  currentPage: number;
}

export interface commonresponse {
  status: 'Success' | 'Failure';
  error?: string;
  data: any;
}

export interface getUserType extends Omit<commonresponse, 'data'> {
  data: User;
}

export interface refreshTokenRes extends Omit<commonresponse, 'data'> {
  accessToken: string;
}

export interface getLoansType extends Omit<commonresponse, 'data'> {
  data: Loan[];
}
