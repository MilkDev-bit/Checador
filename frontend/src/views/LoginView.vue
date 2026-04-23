<template>
  <div class="min-h-screen min-h-dvh bg-mesh flex items-center justify-center px-4 py-8" style="background-color: #0f1629;">
    <!-- Background blobs -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-80 h-80 rounded-full opacity-20 animate-pulse-slow"
        style="background: radial-gradient(circle, #6366f1, transparent 70%)"></div>
      <div class="absolute -bottom-40 -left-40 w-96 h-96 rounded-full opacity-15 animate-float"
        style="background: radial-gradient(circle, #06b6d4, transparent 70%)"></div>
    </div>

    <div class="relative w-full max-w-md animate-in">

      <!-- Card -->
      <div class="glass-card p-8" style="background: rgba(20,29,53,0.9);">
        <h2 class="text-xl font-semibold text-white mb-6">Iniciar Sesión</h2>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="input-label">Correo electrónico</label>
            <input v-model="form.email" type="email" class="input" placeholder="correo@ejemplo.com" required autocomplete="email" />
          </div>
          <div>
            <label class="input-label">Contraseña</label>
            <div class="relative">
              <input v-model="form.password" :type="showPwd ? 'text' : 'password'" class="input pr-12" placeholder="••••••••" required autocomplete="current-password" />
              <button type="button" @click="showPwd = !showPwd" class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-300 text-lg">
                {{ showPwd ? '🙈' : '👁️' }}
              </button>
            </div>
          </div>

          <div v-if="error" class="flex items-start gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-400 rounded-xl px-4 py-3 text-sm">
            <span>⚠️</span><span>{{ error }}</span>
          </div>

          <button type="submit" class="btn-primary btn-lg w-full mt-2" :disabled="loading">
            <span v-if="loading" class="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span>{{ loading ? 'Iniciando...' : 'Iniciar Sesión' }}</span>
          </button>
        </form>

        <div class="mt-6 pt-6 border-t border-white/8 text-center text-sm text-slate-400">
          ¿No tienes cuenta?
          <router-link to="/register" class="text-brand-400 font-semibold hover:text-brand-300 ml-1 transition-colors">
            Crear perfil →
          </router-link>
        </div>
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
const showPwd = ref(false)
const form = ref({ email: '', password: '' })

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    const data = await auth.login(form.value.email, form.value.password)
    if (data.user?.role === 'admin') {
      router.push('/admin')
    } else {
      router.push('/checkin')
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Error al iniciar sesión'
  } finally {
    loading.value = false
  }
}
</script>

