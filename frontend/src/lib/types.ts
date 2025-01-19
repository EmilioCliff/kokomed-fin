import { ReactNode } from "react";
import { User } from "@/data/schema";

export interface authCtx {
    isLoading: boolean;
    isAuthenticated: boolean;
    error: any;
    // add login, logout and checkSession methods
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

  