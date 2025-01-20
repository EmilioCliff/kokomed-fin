import { ReactNode } from 'react';
import { User, Loan } from '@/data/schema';

enum role {}

export interface authCtx {
  isLoading: boolean;
  isAuthenticated: boolean;
  role: string;
  error: any;
  // add login, logout and checkSession methods
}

export interface tableCtx {
  search: string;
  filter: tableFilterType;
  pageIndex: number;
  pageSize: number;
  selectedRow: any;
  rowsCount: number;
  setSearch: (value: string) => void;
  setFilter: React.Dispatch<React.SetStateAction<tableFilterType>>;
  setPageIndex: (value: number) => void;
  setPageSize: (value: number) => void;
  setRowsCount: (value: number) => void;
  setSelectedRow: (value: any) => void;
  resetTableState: () => void;
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
