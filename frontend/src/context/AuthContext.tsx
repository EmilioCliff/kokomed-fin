import { createContext, FC, useState, useEffect } from "react";
import { contextWrapperProps } from "@/lib/types";
import { setupAxiosInterceptors } from "@/API/api";
import { tokenData } from "@/lib/types";
import loginService from "@/services/login";
import refreshTokenService from "@/services/refreshToken";
import logoutService from "@/services/logout";

export interface AuthContextType {
  accessToken: string | null;
  isAuthenticated: boolean;
  isChecking: boolean;
  login: (email: string, password: string) => Promise<tokenData>;
  logout: () => Promise<void>;
  refreshSession: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthContextWrapper: FC<contextWrapperProps> = ({ children }) => {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [isChecking, setIsChecking] = useState(true);
  const isAuthenticated = !!accessToken;

  // Fetch session on app load
  useEffect(() => {
    const checkSession = async () => {
      await refreshSession();
      setIsChecking(false);
    };

    checkSession();
  }, []);

  const getAccessToken = () => accessToken;

  const refreshSession = async () => {
    try {
      const data = await refreshTokenService();
      setAccessToken(data.accessToken);
    } catch (error) {
      console.error("Session refresh failed", error);
      setAccessToken(null);
      setIsChecking(false);
      throw error;
    }
  };

  const login = async (email: string, password: string) => {
    try {
      const data = await loginService(email, password);
      setAccessToken(data.accessToken);
      return data;
    } catch (error: any) {
      throw error;
    }
  };

  const logout = async () => {
    await logoutService();
    setAccessToken(null);
  };

  useEffect(() => {
    setupAxiosInterceptors(getAccessToken, refreshSession);
  }, [accessToken]);

  return (
    <AuthContext.Provider
      value={{ accessToken, isAuthenticated, isChecking, login, logout, refreshSession }}
    >
      {children}
    </AuthContext.Provider>
  );
};
