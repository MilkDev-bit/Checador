<template>
  <div class="min-h-screen min-h-dvh flex items-center justify-center px-4 py-8 bg-mesh"
    style="background-color: var(--bg);">
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
          <UserIcon class="w-7 h-7 text-white" />
        </div>
        <h1 class="text-2xl font-bold" style="color: var(--text);">Crear Perfil</h1>
        <p class="text-sm mt-1" style="color: var(--text-muted);">Regístrate para comenzar a checar</p>
      </div>

      <div class="glass-card p-7">
        <div class="flex justify-end mb-2">
          <ThemeToggle />
        </div>
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
            <label class="input-label">Contraseña <span class="text-xs font-normal" style="color: var(--text-muted);">(mínimo 8 caracteres)</span></label>
            <div class="relative">
              <input v-model="form.password" :type="showPwd ? 'text' : 'password'" class="input pr-12" placeholder="••••••••" required minlength="8" />
              <button type="button" @click="showPwd = !showPwd"
                class="absolute right-3 top-1/2 -translate-y-1/2 transition-colors"
                style="color: var(--text-muted);">
                <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
                <EyeIcon v-else class="w-5 h-5" />
              </button>
            </div>
          </div>

          <div>
            <label class="input-label">Confirmar contraseña</label>
            <input v-model="form.confirm_password" type="password" class="input" placeholder="Repite tu contraseña" required />
            <p v-if="passwordMismatch" class="text-rose-400 text-xs mt-1.5 flex items-center gap-1">
              <ExclamationTriangleIcon class="w-3.5 h-3.5" /> Las contraseñas no coinciden
            </p>
          </div>

          <!-- Math CAPTCHA -->
          <div class="rounded-xl p-4 border" style="background: rgba(99,102,241,0.08); border-color: rgba(99,102,241,0.25);">
            <label class="input-label mb-3">Verificación de seguridad</label>
            <div class="flex items-center gap-3 flex-wrap">
              <div class="flex items-center gap-2 px-4 py-2.5 rounded-xl font-mono text-lg font-bold select-none"
                style="background: rgba(99,102,241,0.2); border: 1px solid rgba(99,102,241,0.3); color: var(--text);">
                {{ captcha.a }} + {{ captcha.b }} = ?
              </div>
              <input v-model="captchaAnswer" type="number" class="input w-24 text-center" placeholder="?" required />
              <button type="button" @click="refreshCaptcha" class="btn-ghost" title="Nuevo captcha">
                <ArrowPathIcon class="w-5 h-5" />
              </button>
            </div>
            <p v-if="captchaError" class="text-rose-400 text-xs mt-2 flex items-center gap-1">
              <ExclamationTriangleIcon class="w-3.5 h-3.5" /> Respuesta incorrecta, intenta de nuevo
            </p>
          </div>

          <div v-if="error" class="flex items-start gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-400 rounded-xl px-4 py-3 text-sm">
            <ExclamationTriangleIcon class="w-4 h-4 flex-shrink-0 mt-0.5" /><span>{{ error }}</span>
          </div>
          <div v-if="success" class="flex items-start gap-2 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 rounded-xl px-4 py-3 text-sm">
            <CheckCircleIcon class="w-4 h-4 flex-shrink-0 mt-0.5" /><span>¡Perfil creado! Redirigiendo al login...</span>
          </div>

          <button type="submit" class="btn-primary btn-lg w-full" :disabled="loading || passwordMismatch">
            <span v-if="loading" class="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span>{{ loading ? 'Creando perfil...' : 'Crear Perfil' }}</span>
          </button>
        </form>

        <div class="mt-5 pt-5 text-center text-sm" style="border-top: 1px solid var(--border-subtle); color: var(--text-muted);">
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
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon, ArrowPathIcon, UserIcon, CheckCircleIcon } from '@heroicons/vue/24/outline'
import ThemeToggle from '@/components/ThemeToggle.vue'

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