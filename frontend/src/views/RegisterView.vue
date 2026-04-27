<template>
  <div class="min-h-screen min-h-dvh flex items-center justify-center p-4 relative overflow-hidden bg-mesh dark:bg-surface-950">
    
    <!-- Background accents -->
    <div class="absolute inset-0 z-0 overflow-hidden pointer-events-none">
      <div class="absolute w-[800px] h-[800px] rounded-full blur-[120px] opacity-10 bg-brand-400 animate-float-slow -top-40 -left-40 mix-blend-screen"></div>
      <div class="absolute w-[600px] h-[600px] rounded-full blur-[100px] opacity-[0.08] bg-sky-300 animate-float-delayed -bottom-20 -right-20 mix-blend-screen"></div>
      <div class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMSIgY3k9IjEiIHI9IjEiIGZpbGw9InJnYmEoMTI4LCAxMjgsIDEyOCwgMC4yKSIvPjwvc3ZnPg==')] opacity-50"></div>
    </div>

    <!-- Container -->
    <div class="relative z-10 w-full max-w-[600px] glass-panel border border-white/20 dark:border-white/5 p-8 lg:p-12 animate-slide-up">
      
      <div class="text-center mb-10">
        <div class="inline-flex w-16 h-16 rounded-2xl overflow-hidden bg-brand-50 text-brand-600 dark:bg-brand-900/40 dark:text-brand-400 mb-4 shadow-inner-light ring-1 ring-brand-500/10 animate-bounce-subtle">
          <img src="/checador-logo.png" alt="PaseLista Logo" class="w-full h-full object-cover" />
        </div>
        <h2 class="text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-2">Crear Cuenta</h2>
        <p class="text-slate-500 dark:text-slate-400 font-medium">Únete a nosotros completando tus datos</p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-5">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <div class="space-y-1 animate-slide-up" style="animation-delay: 0.1s; animation-fill-mode: both;">
            <label class="label-base">Nombre(s)</label>
            <input v-model="form.first_name" type="text" 
              class="input-base focus:ring-4 focus:ring-brand-500/20"
              placeholder="Ej. Carlos" required />
          </div>
          <div class="space-y-1 animate-slide-up" style="animation-delay: 0.2s; animation-fill-mode: both;">
            <label class="label-base">Apellido(s)</label>
            <input v-model="form.last_name" type="text" 
              class="input-base focus:ring-4 focus:ring-brand-500/20"
              placeholder="Ej. García" required />
          </div>
        </div>

        <div class="space-y-1 animate-slide-up" style="animation-delay: 0.3s; animation-fill-mode: both;">
          <label class="label-base">Proyecto</label>
          <input v-model="form.project_name" type="text" 
            class="input-base focus:ring-4 focus:ring-brand-500/20"
            placeholder="Nombre de tu empresa o sucursal" required />
        </div>
        
        <div class="space-y-1 animate-slide-up" style="animation-delay: 0.4s; animation-fill-mode: both;">
          <label class="label-base">Correo Electrónico</label>
          <input v-model="form.email" type="email" 
            class="input-base focus:ring-4 focus:ring-brand-500/20"
            placeholder="usuario@correo.com" required autocomplete="email" />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <div class="space-y-1 animate-slide-up" style="animation-delay: 0.5s; animation-fill-mode: both;">
            <label class="label-base">Contraseña</label>
            <div class="relative group">
              <input v-model="form.password" :type="showPwd ? 'text' : 'password'" 
                class="input-base tracking-widest focus:ring-4 focus:ring-brand-500/20"
                placeholder="••••••••" required minlength="8" />
              <button type="button" @click="showPwd = !showPwd"
                class="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg text-slate-400 hover:text-slate-700 dark:text-slate-500 dark:hover:text-slate-300 transition-colors focus:outline-none">
                <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
                <EyeIcon v-else class="w-5 h-5" />
              </button>
            </div>
          </div>

          <div class="space-y-1 animate-slide-up" style="animation-delay: 0.6s; animation-fill-mode: both;">
            <label class="label-base">Confirmar</label>
            <div class="relative group">
              <input v-model="form.confirm_password" type="password" 
                class="input-base tracking-widest focus:ring-4 focus:ring-brand-500/20"
                placeholder="••••••••" required />
            </div>
          </div>
        </div>
        
        <Transition name="fade">
          <p v-if="passwordMismatch" class="text-rose-500 dark:text-rose-400 text-xs font-bold flex items-center justify-center gap-1.5 animate-shake">
            <ExclamationTriangleIcon class="w-4 h-4" /> Las contraseñas no coinciden
          </p>
        </Transition>

        <!-- reCAPTCHA v2 Container -->
        <div class="space-y-1 animate-fade-in" style="animation-delay: 0.7s; animation-fill-mode: both;">
          <label class="label-base">Verificación de seguridad</label>
          <div class="flex items-center justify-center bg-slate-50 dark:bg-surface-800 rounded-xl border border-slate-200 dark:border-white/10 p-2 min-h-[90px] overflow-hidden transition-all shadow-sm">
            <div id="recaptcha-container" class="scale-[0.8] sm:scale-95 origin-center"></div>
          </div>
        </div>

        <Transition name="fade">
          <div v-if="error" class="flex items-start gap-3 bg-rose-50/80 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 text-rose-600 dark:text-rose-400 rounded-xl px-4 py-3 text-sm font-semibold animate-shake">
            <ExclamationTriangleIcon class="w-5 h-5 flex-shrink-0 mt-0.5" />
            <span class="leading-snug">{{ error }}</span>
          </div>
        </Transition>
        <Transition name="fade">
          <div v-if="success" class="flex items-start gap-3 bg-emerald-50/80 dark:bg-emerald-500/10 border border-emerald-200 dark:border-emerald-500/20 text-emerald-600 dark:text-emerald-400 rounded-xl px-4 py-3 text-sm font-semibold">
            <CheckCircleIcon class="w-5 h-5 flex-shrink-0 mt-0.5" />
            <span class="leading-snug">¡Cuenta creada con éxito! Redirigiendo...</span>
          </div>
        </Transition>

        <button type="submit" class="btn-primary w-full py-3.5 rounded-xl font-bold tracking-wide mt-2 animate-slide-up hover:ring-4 hover:ring-brand-500/20 transition-all" :disabled="loading || passwordMismatch" style="animation-delay: 0.8s; animation-fill-mode: both;">
          <div class="relative flex items-center justify-center gap-2">
            <span v-if="loading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span>{{ loading ? 'Procesando...' : 'Registrarme' }}</span>
          </div>
        </button>
      </form>

      <div class="mt-8 text-center animate-fade-in" style="animation-delay: 0.9s; animation-fill-mode: both;">
        <p class="text-slate-500 dark:text-slate-400 font-medium text-sm">
          ¿Ya tienes una cuenta? 
          <router-link to="/login" class="text-brand-600 dark:text-brand-400 font-bold hover:text-brand-700 dark:hover:text-brand-300 hover:underline underline-offset-4 transition-all">
            Inicia sesión aquí
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon, CheckCircleIcon } from '@heroicons/vue/24/outline'

const auth = useAuthStore()
const themeStore = useThemeStore()
const router = useRouter()
const loading = ref(false)
const error = ref('')
const success = ref(false)
const showPwd = ref(false)
const recaptchaWidgetId = ref(null)

onMounted(() => {
  const checkRecaptcha = setInterval(() => {
    if (window.grecaptcha && window.grecaptcha.render) {
      clearInterval(checkRecaptcha)
      try {
        recaptchaWidgetId.value = window.grecaptcha.render('recaptcha-container', {
          sitekey: '6LdUjcgsAAAAAI0pgOSk7QMEmq-zjD4saihqzaa-', 
          theme: themeStore.isDark ? 'dark' : 'light'
        })
      } catch (e) {
        console.warn('Recaptcha error:', e)
      }
    }
  }, 100)
})

const form = ref({
  first_name: '', last_name: '', project_name: '',
  email: '', password: '', confirm_password: ''
})

const passwordMismatch = computed(
  () => form.value.confirm_password.length > 0 && form.value.password !== form.value.confirm_password
)

async function handleRegister() {
  error.value = ''
  
  let recaptchaToken = ''
  if (recaptchaWidgetId.value !== null) {
    recaptchaToken = window.grecaptcha.getResponse(recaptchaWidgetId.value)
  }

  if (!recaptchaToken) {
    error.value = 'Por favor, completa la verificación de seguridad.'
    return
  }

  if (form.value.password !== form.value.confirm_password) return
  
  loading.value = true
  try {
    await auth.register({
      ...form.value,
      recaptcha_token: recaptchaToken
    })
    success.value = true
    setTimeout(() => router.push('/login'), 1500)
  } catch (e) {
    error.value = e.response?.data?.error || 'Error al crear el perfil'
    if (recaptchaWidgetId.value !== null) {
      window.grecaptcha.reset(recaptchaWidgetId.value)
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@keyframes float-slow {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(20px, -20px) scale(1.05); }
}
@keyframes float-delayed {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-30px, 30px) scale(1.1); }
}
.animate-float-slow { animation: float-slow 10s ease-in-out infinite; }
.animate-float-delayed { animation: float-delayed 12s ease-in-out infinite; }

@keyframes bounce-subtle {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}
.animate-bounce-subtle { animation: bounce-subtle 4s ease-in-out infinite; }

@keyframes slide-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-slide-up { opacity: 0; animation: slide-up 0.5s ease-out forwards; }

@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}
.animate-fade-in { opacity: 0; animation: fade-in 0.8s ease-out forwards; }

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}
.animate-shake { animation: shake 0.4s ease-in-out; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>