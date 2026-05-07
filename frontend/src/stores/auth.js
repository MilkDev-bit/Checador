import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const token = ref(localStorage.getItem('auth_token') || null)

  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function persistUser(u) {
    if (u !== undefined) user.value = u
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  function persistToken(t) {
    token.value = t
    if (t) {
      localStorage.setItem('auth_token', t)
    } else {
      localStorage.removeItem('auth_token')
    }
  }

  async function login(email, password, recaptcha_token) {
    const { data } = await api.post('/auth/login', { email, password, recaptcha_token })
    user.value = data.user
    persistUser()
    persistToken(data.token || null)
    return data
  }

  async function register(payload) {
    const { data } = await api.post('/auth/register', payload)
    if (data.user) {
      user.value = data.user
      persistUser()
    }
    persistToken(data.token || null)
    return data
  }

  async function logout() {
    try {
      await api.post('/auth/logout')
    } catch (e) {
      console.warn('Error during server logout', e)
    } finally {
      user.value = null
      token.value = null
      localStorage.removeItem('user')
      localStorage.removeItem('auth_token')
    }
  }

  return { user, token, isLoggedIn, isAdmin, login, register, logout, persistUser }
})
