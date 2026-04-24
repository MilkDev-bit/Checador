import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true // Ensures cookies are sent with every request
})

api.interceptors.request.use(config => {
  // Token is now handled automatically via HttpOnly cookies
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('user') // Clear user info on session expiration
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
