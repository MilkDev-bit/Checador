<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4">
    <div class="card w-full max-w-md">
      <div class="text-center mb-8">
        <div class="text-5xl mb-3">✅</div>
        <h1 class="text-3xl font-bold text-blue-700">PaseLista</h1>
        <p class="text-gray-500 mt-1">Checador en Línea</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Correo electrónico</label>
          <input v-model="form.email" type="email" class="input-field" placeholder="correo@ejemplo.com" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Contraseña</label>
          <input v-model="form.password" type="password" class="input-field" placeholder="••••••••" required />
        </div>

        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 rounded-xl px-4 py-3 text-sm">
          {{ error }}
        </div>

        <button type="submit" class="btn-primary w-full" :disabled="loading">
          <span v-if="loading">Iniciando sesión...</span>
          <span v-else>Iniciar Sesión</span>
        </button>
      </form>

      <div class="mt-6 text-center text-sm text-gray-600">
        ¿No tienes cuenta?
        <router-link to="/register" class="text-blue-600 font-semibold hover:underline ml-1">
          Crear perfil
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const loading = ref(false)
const error = ref('')
const form = ref({ email: '', password: '' })

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(form.value.email, form.value.password)
    router.push('/checkin')
  } catch (e) {
    error.value = e.response?.data?.error || 'Error al iniciar sesión'
  } finally {
    loading.value = false
  }
}
</script>
