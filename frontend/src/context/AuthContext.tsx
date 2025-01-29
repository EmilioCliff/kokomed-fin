import { createContext, FC, useState } from 'react';
import { contextWrapperProps } from '@/lib/types';
import { authCtx } from '@/lib/types';
import { role } from '@/lib/types';
import { tokenData } from '@/lib/types';
import { toast } from 'react-toastify';
import { z } from 'zod';
// import { useQuery } from "@tanstack/react-query"

export const AuthContext = createContext<authCtx>({
  isLoading: false,
  isAuthenticated: false,
  userRole: role.GUEST,
  error: null,
  updateAuthContext: () => {},
});

// export const useAuthContext = () => {
//   const ctx = useContext(AuthContext);
//   return ctx;
// };

export const AuthContextWrapper: FC<contextWrapperProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [userRole, setUserRole] = useState(role.GUEST);
  const [error, setError] = useState(null);

  const updateAuthContext = (tokenData: tokenData) => {
    setIsAuthenticated(true);
    setUserRole(tokenData.userData.role);

    toast.success(tokenData.accessToken);
  };

  return (
    <AuthContext.Provider
      value={{
        isLoading,
        isAuthenticated,
        userRole,
        error,
        updateAuthContext,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
