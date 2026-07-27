import axios from 'axios'

function getBaseURL() {
  const envBase = import.meta.env.VITE_API_BASE
  if (!envBase || envBase === '' || envBase === '/') {
    return ''
  }
  
  if (typeof window !== 'undefined' && window.location && window.location.hostname) {
    const currentHost = window.location.hostname
    if (currentHost !== 'localhost' && currentHost !== '127.0.0.1') {
      return envBase.replace('localhost', currentHost).replace('127.0.0.1', currentHost)
    }
  }
  
  return envBase
}

const api = axios.create({
  baseURL: getBaseURL(),
  timeout: 10000,
})

// Request interceptor - add token from sessionStorage or localStorage
api.interceptors.request.use(
  (config) => {
    const token = sessionStorage.getItem('token') || localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor - handle 401 unauthorized status cleanly
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      sessionStorage.clear()
      localStorage.clear()
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
