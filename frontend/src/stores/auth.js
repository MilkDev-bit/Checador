import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(email, password, recaptcha_token) {
    const { data } = await api.post('/auth/login', { email, password, recaptcha_token })
    user.value = data.user
    localStorage.setItem('user', JSON.stringify(data.user))
    return data
  }

  async function register(payload) {
    const { data } = await api.post('/auth/register', payload)
    return data
  }

  async function logout() {
    try {
      await api.post('/auth/logout')
    } catch (e) {
      console.warn('Error during server logout', e)
    } finally {
      user.value = null
      localStorage.removeItem('user')
      // Redirect to login handled by components or router
    }
  }

  return { user, isLoggedIn, isAdmin, login, register, logout }
})
