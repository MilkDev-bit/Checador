import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true // Send cookies when available
})

api.interceptors.request.use(config => {
  // Send stored token as Authorization header — fallback for when cookies are blocked
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('user')
      localStorage.removeItem('auth_token')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
