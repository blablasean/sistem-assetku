import axios from 'axios'

function getBaseURL() {
  const envBase = import.meta.env.VITE_API_BASE
  if (envBase && envBase !== '' && envBase !== '/') {
    if (typeof window !== 'undefined' && window.location && window.location.hostname) {
      const currentHost = window.location.hostname
      if (currentHost !== 'localhost' && currentHost !== '127.0.0.1') {
        return envBase.replace('localhost', currentHost).replace('127.0.0.1', currentHost)
      }
    }
    return envBase
  }

  // Fallback: If no envBase set, dynamically map frontend port to backend Go port 8080
  if (typeof window !== 'undefined' && window.location) {
    const currentHost = window.location.hostname || 'localhost'
    const port = window.location.port
    // In serve / build mode on separate port (e.g. 3000, 5000, 5173), direct API requests to port 8080
    if (port && port !== '8080') {
      return `http://${currentHost}:8080`
    }
  }

  return ''
}

const api = axios.create({
  baseURL: getBaseURL(),
  timeout: 10000,
})

// Request interceptor - add token from sessionStorage
api.interceptors.request.use(
  (config) => {
    const token = sessionStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

import { clearSessionAndForceLogin } from './utils/auth'

// Response interceptor - handle 401/403 unauthorized status cleanly (except for login requests)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const isLoginEndpoint = error.config && error.config.url && error.config.url.includes('/auth/login')
    if ((error.response?.status === 401 || error.response?.status === 403) && !isLoginEndpoint) {
      clearSessionAndForceLogin()
    }
    return Promise.reject(error)
  }
)

export default api
