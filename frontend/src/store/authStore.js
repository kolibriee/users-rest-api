import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import jwtDecode from 'jwt-decode'

const useAuthStore = create(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      isAdmin: false,
      
      setToken: (token) => {
        if (token) {
          const decoded = jwtDecode(token)
          set({ 
            token,
            user: decoded,
            isAdmin: decoded.role === 'admin'
          })
        } else {
          set({ token: null, user: null, isAdmin: false })
        }
      },

      logout: () => {
        set({ token: null, user: null, isAdmin: false })
      },

      getToken: () => get().token,
      
      isAuthenticated: () => !!get().token,
    }),
    {
      name: 'auth-storage',
    }
  )
)

export default useAuthStore
