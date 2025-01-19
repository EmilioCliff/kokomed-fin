import { createContext, useContext, FC } from "react"
import { contextWrapperProps } from "@/lib/types"
import { authCtx } from "@/lib/types"
// import { useQuery } from "@tanstack/react-query"

export const AuthContext = createContext<authCtx>({ isLoading: false, isAuthenticated: false, error: null })

export const useAuthContext = () => {
    const ctx = useContext(AuthContext)
  return ctx
}  


export const AuthContextWrapper: FC<contextWrapperProps> = ({ children }) => {
    // check if user is authenticated from stored key in the local storage
    // if not authenticated, redirect to login page
  return (
    <AuthContext.Provider value={{ isLoading: false, isAuthenticated: false, error: null }}>{children}</AuthContext.Provider>
  )
}
