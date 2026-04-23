<template>
  <div class="min-h-screen min-h-dvh flex items-center justify-center px-4 py-8 bg-mesh" style="background-color: #0f1629;">
    <!-- Background blobs -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -left-40 w-80 h-80 rounded-full opacity-15 animate-pulse-slow"
        style="background: radial-gradient(circle, #8b5cf6, transparent 70%)"></div>
      <div class="absolute -bottom-20 -right-20 w-72 h-72 rounded-full opacity-15 animate-float"
        style="background: radial-gradient(circle, #06b6d4, transparent 70%)"></div>
    </div>

    <div class="relative w-full max-w-lg animate-in">
      <!-- Header -->
      <div class="text-center mb-6">
        <div class="inline-flex items-center justify-center w-14 h-14 rounded-2xl mb-3"
          style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">
          <span class="text-2xl">👤</span>
        </div>
        <h1 class="text-2xl font-bold text-white">Crear Perfil</h1>
        <p class="text-slate-400 text-sm mt-1">Regístrate para comenzar a checar</p>
      </div>

      <div class="glass-card p-7" style="background: rgba(20,29,53,0.9);">
        <form @submit.prevent="handleRegister" class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="input-label">Nombre(s)</label>
              <input v-model="form.first_name" type="text" class="input" placeholder="Juan" required />
            </div>
            <div>
              <label class="input-label">Apellido(s)</label>
              <input v-model="form.last_name" type="text" class="input" placeholder="Pérez García" required />
            </div>
          </div>

          <div>
            <label class="input-label">Nombre completo del proyecto</label>
            <input v-model="form.project_name" type="text" class="input" placeholder="Proyecto de Construcción Centro" required />
          </div>

          <div>
            <label class="input-label">Correo electrónico</label>
            <input v-model="form.email" type="email" class="input" placeholder="correo@ejemplo.com" required autocomplete="email" />
          </div>

          <div>
            <label class="input-label">Contraseña <span class="text-slate-500 text-xs font-normal">(mínimo 8 caracteres)</span></label>
            <div class="relative">
              <input v-model="form.password" :type="showPwd ? 'text' : 'password'" class="input pr-12" placeholder="••••••••" required minlength="8" />
              <button type="button" @click="showPwd = !showPwd" class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-300">
                {{ showPwd ? '🙈' : '👁️' }}
              </button>
            </div>
          </div>

          <div>
            <label class="input-label">Confirmar contraseña</label>
            <input v-model="form.confirm_password" type="password" class="input" placeholder="Repite tu contraseña" required />
            <p v-if="passwordMismatch" class="text-rose-400 text-xs mt-1.5 flex items-center gap-1">
              <span>⚠️</span> Las contraseñas no coinciden
            </p>
          </div>

          <!-- Math CAPTCHA -->
          <div class="rounded-xl p-4 border" style="background: rgba(99,102,241,0.08); border-color: rgba(99,102,241,0.25);">
            <label class="input-label mb-3">Verificación de seguridad</label>
            <div class="flex items-center gap-3 flex-wrap">
              <div class="flex items-center gap-2 px-4 py-2.5 rounded-xl font-mono text-lg font-bold text-white select-none"
                style="background: rgba(99,102,241,0.2); border: 1px solid rgba(99,102,241,0.3);">
                {{ captcha.a }} + {{ captcha.b }} = ?
              </div>
              <input v-model="captchaAnswer" type="number" class="input w-24 text-center" placeholder="?" required />
              <button type="button" @click="refreshCaptcha" class="btn-ghost text-xl" title="Nuevo captcha">🔄</button>
            </div>
            <p v-if="captchaError" class="text-rose-400 text-xs mt-2 flex items-center gap-1">
              <span>⚠️</span> Respuesta incorrecta, intenta de nuevo
            </p>
          </div>

          <div v-if="error" class="flex items-start gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-400 rounded-xl px-4 py-3 text-sm">
            <span>⚠️</span><span>{{ error }}</span>
          </div>
          <div v-if="success" class="flex items-start gap-2 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 rounded-xl px-4 py-3 text-sm">
            <span></span><span>¡Perfil creado! Redirigiendo al login...</span>
          </div>

          <button type="submit" class="btn-primary btn-lg w-full" :disabled="loading || passwordMismatch">
            <span v-if="loading" class="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span>{{ loading ? 'Creando perfil...' : 'Crear Perfil' }}</span>
          </button>
        </form>

        <div class="mt-5 pt-5 border-t border-white/8 text-center text-sm text-slate-400">
          ¿Ya tienes cuenta?
          <router-link to="/login" class="text-brand-400 font-semibold hover:text-brand-300 ml-1 transition-colors">
            Iniciar sesión →
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const loading = ref(false)
const error = ref('')
const success = ref(false)
const captchaAnswer = ref('')
const captchaError = ref(false)
const showPwd = ref(false)
const captcha = ref({ a: 0, b: 0 })

function refreshCaptcha() {
  captcha.value = {
    a: Math.floor(Math.random() * 10) + 1,
    b: Math.floor(Math.random() * 10) + 1
  }
  captchaAnswer.value = ''
  captchaError.value = false
}

onMounted(refreshCaptcha)

const form = ref({
  first_name: '', last_name: '', project_name: '',
  email: '', password: '', confirm_password: ''
})

const passwordMismatch = computed(
  () => form.value.confirm_password.length > 0 && form.value.password !== form.value.confirm_password
)

async function handleRegister() {
  captchaError.value = false
  error.value = ''
  if (parseInt(captchaAnswer.value) !== captcha.value.a + captcha.value.b) {
    captchaError.value = true
    return
  }
  if (form.value.password !== form.value.confirm_password) return
  loading.value = true
  try {
    await auth.register({
      first_name: form.value.first_name,
      last_name: form.value.last_name,
      project_name: form.value.project_name,
      email: form.value.email,
      password: form.value.password
    })
    success.value = true
    setTimeout(() => router.push('/login'), 1500)
  } catch (e) {
    error.value = e.response?.data?.error || 'Error al crear el perfil'
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}
</script>