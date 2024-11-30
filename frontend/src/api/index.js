import axios from 'axios'
import useAuthStore from '../store/authStore'
import jwt_decode from 'jwt-decode'

export const API_URL = import.meta.env.VITE_API_HOST || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_URL,
  withCredentials: true,
})

// Add a request interceptor
api.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Auth API
export const signIn = (data) => api.post('/auth/sign-in', data)
export const signUp = (data) => api.post('/auth/sign-up', data)
export const refreshToken = () => api.get('/auth/refresh')

// User API
export const getCurrentUser = () => {
  const token = useAuthStore.getState().token
  if (!token) return Promise.reject('No token found')
  const decoded = jwt_decode(token)
  return api.get(`/api/users/${decoded.sub}`)
}

export const updateCurrentUser = (data) => {
  const token = useAuthStore.getState().token
  if (!token) return Promise.reject('No token found')
  const decoded = jwt_decode(token)
  return api.put(`/api/users/${decoded.sub}`, data)
}

export const deleteCurrentUser = () => {
  const token = useAuthStore.getState().token
  if (!token) return Promise.reject('No token found')
  const decoded = jwt_decode(token)
  return api.delete(`/api/users/${decoded.sub}`)
}

// Admin API
export const getAllUsers = () => api.get('/admin/users')
export const getUserById = (id) => api.get(`/admin/users/${id}`)
export const createUser = (data) => api.post('/admin/users', data)
export const updateUser = (id, data) => api.put(`/admin/users/${id}`, data)
export const deleteUser = (id) => api.delete(`/admin/users/${id}`)

export default api
