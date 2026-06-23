<template>
  <div
    class="min-h-screen min-h-dvh flex items-center justify-center p-4 relative overflow-hidden bg-mesh dark:bg-surface-950">

    <!-- Animated background elements (refined) -->
    <div class="absolute inset-0 z-0 overflow-hidden pointer-events-none">
      <div
        class="absolute w-[800px] h-[800px] rounded-full blur-[120px] opacity-10 bg-brand-400 animate-float-slow -top-40 -left-40 mix-blend-screen">
      </div>
      <div
        class="absolute w-[600px] h-[600px] rounded-full blur-[100px] opacity-[0.08] bg-sky-300 animate-float-delayed -bottom-20 -right-20 mix-blend-screen">
      </div>
      <div
        class="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMSIgY3k9IjEiIHI9IjEiIGZpbGw9InJnYmEoMTI4LCAxMjgsIDEyOCwgMC4yKSIvPjwvc3ZnPg==')] opacity-50">
      </div>
    </div>

    <!-- Main Card -->
    <div
      class="relative z-10 w-full max-w-[420px] glass-panel border border-white/20 dark:border-white/5 p-8 sm:p-10 animate-slide-up">

      <div class="text-center mb-8">
        <div
          class="inline-flex w-16 h-16 rounded-2xl overflow-hidden bg-brand-50 text-brand-600 dark:bg-brand-900/40 dark:text-brand-400 mb-6 shadow-inner-light ring-1 ring-brand-500/10 animate-bounce-subtle">
          <img src="/LOGOC.png" alt="PaseLista Logo" class="w-full h-full object-cover" />
        </div>
        <h2 class="text-2xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight mb-2">Bienvenido de
          nuevo</h2>
        <p class="text-slate-500 dark:text-slate-400 font-medium text-sm">Ingresa a tu cuenta para continuar</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-5">
        <div class="space-y-1 animate-slide-up" style="animation-delay: 0.1s; animation-fill-mode: both;">
          <label class="label-base">Correo Electrónico</label>
          <div class="relative group">
            <input v-model="form.email" type="email" class="input-base focus:ring-4 focus:ring-brand-500/20"
              placeholder="ejemplo@correo.com" required autocomplete="email" />
          </div>
        </div>

        <div class="space-y-1 animate-slide-up" style="animation-delay: 0.2s; animation-fill-mode: both;">
          <label class="label-base">Contraseña</label>
          <div class="relative group">
            <input v-model="form.password" :type="showPwd ? 'text' : 'password'"
              class="input-base tracking-widest focus:ring-4 focus:ring-brand-500/20" placeholder="••••••••" required
              autocomplete="current-password" />
            <button type="button" @click="showPwd = !showPwd"
              class="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg text-slate-400 hover:text-slate-700 dark:text-slate-500 dark:hover:text-slate-300 transition-colors focus:outline-none">
              <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
              <EyeIcon v-else class="w-5 h-5" />
            </button>
          </div>
          <div class="flex justify-end mt-1.5">
            <button type="button" @click="showForgotModal = true"
              class="text-xs text-brand-600 dark:text-brand-400 hover:underline underline-offset-4 font-medium transition-colors">
              ¿Olvidaste tu contraseña?
            </button>
          </div>
        </div>

        <!-- reCAPTCHA v2 Container -->
        <div class="space-y-1 animate-fade-in" style="animation-delay: 0.3s; animation-fill-mode: both;">
          <label class="label-base">Verificación</label>
          <div
            class="flex items-center justify-center bg-slate-50 dark:bg-surface-800 rounded-xl border border-slate-200 dark:border-white/10 p-2 min-h-[90px] overflow-hidden transition-all shadow-sm">
            <!-- Cargando widget -->
            <div v-if="recaptchaLoading" class="flex flex-col items-center gap-2 py-2">
              <div class="w-5 h-5 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin"></div>
              <p class="text-xs text-slate-400 dark:text-slate-500">Cargando verificación...</p>
            </div>
            <!-- Error al cargar -->
            <div v-else-if="recaptchaLoadError" class="flex flex-col items-center gap-1.5 py-2 text-center px-4">
              <ExclamationTriangleIcon class="w-5 h-5 text-amber-500 flex-shrink-0" />
              <p class="text-xs text-slate-500 dark:text-slate-400 leading-tight">No se pudo cargar la verificación de
                seguridad.</p>
              <button type="button" @click="retryRecaptcha"
                class="text-xs text-brand-600 dark:text-brand-400 font-bold hover:underline mt-0.5">↺
                Reintentar</button>
            </div>
            <!-- Widget -->
            <div id="recaptcha-container" class="scale-[0.8] sm:scale-95 origin-center"
              :class="{ hidden: recaptchaLoading || recaptchaLoadError }"></div>
          </div>
        </div>

        <Transition name="fade">
          <div v-if="error"
            class="flex flex-col gap-2 bg-rose-50/80 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 text-rose-600 dark:text-rose-400 rounded-xl px-4 py-3 text-sm animate-shake">
            <div class="flex items-start gap-3 font-semibold">
              <ExclamationTriangleIcon class="w-5 h-5 flex-shrink-0 mt-0.5" />
              <span class="leading-snug">{{ error }}</span>
            </div>
            <!-- Ayuda contextual según tipo de error -->
            <ul v-if="errorHints.length"
              class="mt-1 ml-8 space-y-1 text-xs text-rose-500 dark:text-rose-400/80 list-disc list-outside">
              <li v-for="hint in errorHints" :key="hint">{{ hint }}</li>
            </ul>
          </div>
        </Transition>

        <button type="submit"
          class="btn-primary w-full py-3.5 rounded-xl font-bold tracking-wide mt-2 animate-slide-up hover:ring-4 hover:ring-brand-500/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="loading || (!recaptchaToken && !recaptchaLoadError)"
          style="animation-delay: 0.4s; animation-fill-mode: both;">
          <div class="relative flex items-center justify-center gap-2">
            <span v-if="loading"
              class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span>{{ loading ? 'Ingresando...' : 'Iniciar Sesión' }}</span>
          </div>
        </button>
      </form>

      <div class="mt-8 text-center animate-fade-in" style="animation-delay: 0.5s; animation-fill-mode: both;">
        <p class="text-slate-500 dark:text-slate-400 font-medium text-sm">
          ¿No tienes una cuenta?
          <router-link to="/register"
            class="text-brand-600 dark:text-brand-400 font-bold hover:text-brand-700 dark:hover:text-brand-300 hover:underline underline-offset-4 transition-all">
            Regístrate aquí
          </router-link>
        </p>
      </div>
    </div>
  </div>

  <!-- ===== Forgot Password Modal ===== -->
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="showForgotModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-slate-950/80 backdrop-blur-sm"
        @click.self="closeForgotModal">
        <div
          class="relative w-full max-w-[420px] glass-panel border border-white/20 dark:border-white/5 p-8 sm:p-10 animate-slide-up">

          <!-- Close -->
          <button @click="closeForgotModal"
            class="absolute top-4 right-4 p-1.5 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-white/5 transition-colors">
            <XMarkIcon class="w-5 h-5" />
          </button>

          <!-- Step 1: Enter email -->
          <template v-if="forgotStep === 1">
            <div class="text-center mb-7">
              <div
                class="inline-flex w-14 h-14 rounded-2xl items-center justify-center bg-brand-50 dark:bg-brand-900/40 text-brand-600 dark:text-brand-400 mb-4">
                <EnvelopeIcon class="w-7 h-7" />
              </div>
              <h3 class="text-xl font-extrabold text-slate-900 dark:text-white">Recuperar contraseña</h3>
              <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">Ingresa tu correo y te enviaremos un código de
                6 dígitos.</p>
            </div>
            <form @submit.prevent="handleForgotRequest" class="space-y-4">
              <div>
                <label class="label-base">Correo electrónico</label>
                <input v-model="forgotEmail" type="email" class="input-base" placeholder="ejemplo@correo.com" required
                  autocomplete="email" />
              </div>
              <Transition name="fade">
                <div v-if="forgotError"
                  class="flex items-start gap-3 bg-rose-50/80 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 text-rose-600 dark:text-rose-400 rounded-xl px-4 py-3 text-sm font-semibold">
                  <ExclamationTriangleIcon class="w-5 h-5 flex-shrink-0 mt-0.5" />
                  {{ forgotError }}
                </div>
              </Transition>
              <button type="submit" :disabled="forgotLoading"
                class="btn-primary w-full py-3.5 rounded-xl font-bold disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2">
                <span v-if="forgotLoading"
                  class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                {{ forgotLoading ? 'Enviando...' : 'Enviar código' }}
              </button>
            </form>
          </template>

          <!-- Step 2: Enter code + new password -->
          <template v-else-if="forgotStep === 2">
            <div class="text-center mb-7">
              <div
                class="inline-flex w-14 h-14 rounded-2xl items-center justify-center bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400 mb-4">
                <EnvelopeOpenIcon class="w-7 h-7" />
              </div>
              <h3 class="text-xl font-extrabold text-slate-900 dark:text-white">Revisa tu correo</h3>
              <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">
                Enviamos un código a
                <span class="font-semibold text-slate-700 dark:text-slate-300">{{ forgotEmail }}</span>
              </p>
            </div>
            <form @submit.prevent="handleResetPassword" class="space-y-4">
              <div>
                <label class="label-base">Código de 6 dígitos</label>
                <input v-model="forgotCode" type="text" inputmode="numeric"
                  class="input-base text-center tracking-[0.4em] font-bold text-xl" placeholder="000000" maxlength="6"
                  pattern="[0-9]{6}" required autocomplete="one-time-code" />
              </div>
              <div>
                <label class="label-base">Nueva contraseña</label>
                <div class="relative">
                  <input v-model="forgotNewPwd" :type="showForgotPwd ? 'text' : 'password'" class="input-base pr-12"
                    placeholder="Mínimo 8 caracteres" required minlength="8" />
                  <button type="button" @click="showForgotPwd = !showForgotPwd"
                    class="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg text-slate-400 hover:text-slate-700 dark:text-slate-500 dark:hover:text-slate-300 transition-colors">
                    <EyeSlashIcon v-if="showForgotPwd" class="w-5 h-5" />
                    <EyeIcon v-else class="w-5 h-5" />
                  </button>
                </div>
              </div>
              <div>
                <label class="label-base">Confirmar contraseña</label>
                <input v-model="forgotConfirmPwd" type="password" class="input-base"
                  :class="forgotConfirmPwd && forgotPwdMismatch ? 'ring-2 ring-rose-500/30 border-rose-500/40' : ''"
                  placeholder="••••••••" required />
                <p v-if="forgotConfirmPwd && forgotPwdMismatch" class="text-xs text-rose-400 mt-1 font-semibold">Las
                  contraseñas no coinciden</p>
              </div>
              <Transition name="fade">
                <div v-if="forgotError"
                  class="flex items-start gap-3 bg-rose-50/80 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 text-rose-600 dark:text-rose-400 rounded-xl px-4 py-3 text-sm font-semibold">
                  <ExclamationTriangleIcon class="w-5 h-5 flex-shrink-0 mt-0.5" />
                  {{ forgotError }}
                </div>
              </Transition>
              <button type="submit" :disabled="forgotLoading || forgotPwdMismatch"
                class="btn-primary w-full py-3.5 rounded-xl font-bold disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2">
                <span v-if="forgotLoading"
                  class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                {{ forgotLoading ? 'Cambiando...' : 'Cambiar contraseña' }}
              </button>
              <button type="button" @click="forgotStep = 1"
                class="w-full text-sm text-brand-600 dark:text-brand-400 hover:underline font-medium transition-colors">
                ← Usar otro correo
              </button>
            </form>
          </template>

          <!-- Step 3: Success -->
          <template v-else-if="forgotStep === 3">
            <div class="text-center py-4">
              <div
                class="inline-flex w-16 h-16 rounded-full items-center justify-center bg-emerald-50 dark:bg-emerald-900/30 text-emerald-500 mb-5">
                <CheckCircleIcon class="w-8 h-8" />
              </div>
              <h3 class="text-xl font-extrabold text-slate-900 dark:text-white mb-2">¡Contraseña actualizada!</h3>
              <p class="text-sm text-slate-500 dark:text-slate-400 mb-7">Ya puedes iniciar sesión con tu nueva
                contraseña.</p>
              <button @click="closeForgotModal" class="btn-primary w-full py-3.5 rounded-xl font-bold">
                Iniciar sesión
              </button>
            </div>
          </template>

        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon, EnvelopeIcon, EnvelopeOpenIcon, CheckCircleIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import api from '@/api'

const auth = useAuthStore()
const themeStore = useThemeStore()
const router = useRouter()
const loading = ref(false)
const error = ref('')
const showPwd = ref(false)
const form = ref({ email: '', password: '' })
const recaptchaWidgetId = ref(null)
const recaptchaToken = ref('')

// ── Forgot Password ──────────────────────────────────────────────────────────
const showForgotModal = ref(false)
const forgotStep = ref(1)         // 1=email, 2=code+newpwd, 3=success
const forgotEmail = ref('')
const forgotCode = ref('')
const forgotNewPwd = ref('')
const forgotConfirmPwd = ref('')
const showForgotPwd = ref(false)
const forgotLoading = ref(false)
const forgotError = ref('')
const forgotPwdMismatch = computed(() =>
  forgotNewPwd.value && forgotConfirmPwd.value && forgotNewPwd.value !== forgotConfirmPwd.value
)

function closeForgotModal() {
  showForgotModal.value = false
  forgotStep.value = 1
  forgotEmail.value = ''
  forgotCode.value = ''
  forgotNewPwd.value = ''
  forgotConfirmPwd.value = ''
  showForgotPwd.value = false
  forgotError.value = ''
}

async function handleForgotRequest() {
  forgotError.value = ''
  forgotLoading.value = true
  try {
    await api.post('/auth/forgot-password', { email: forgotEmail.value.trim() })
    forgotStep.value = 2
  } catch (e) {
    forgotError.value = e.response?.data?.error || 'Error al enviar el código. Intenta de nuevo.'
  } finally {
    forgotLoading.value = false
  }
}

async function handleResetPassword() {
  forgotError.value = ''
  if (forgotPwdMismatch.value) return
  forgotLoading.value = true
  try {
    await api.post('/auth/reset-password', {
      email: forgotEmail.value.trim(),
      code: forgotCode.value.trim(),
      password: forgotNewPwd.value,
    })
    forgotStep.value = 3
  } catch (e) {
    forgotError.value = e.response?.data?.error || 'Código incorrecto o expirado. Intenta de nuevo.'
  } finally {
    forgotLoading.value = false
  }
}

// Sugerencias contextuales basadas en el mensaje de error
const errorHints = computed(() => {
  const msg = error.value.toLowerCase()
  if (msg.includes('captcha') || msg.includes('robot')) {
    return [
      'Asegúrate de marcar la casilla "No soy un robot" antes de presionar Iniciar Sesión.',
      'Si aparece un puzzle de imágenes, resuélvelo completamente.',
      'Desactiva bloqueadores de anuncios o VPN si el captcha no carga.',
      'Si el problema persiste, recarga la página (F5) e intenta nuevamente.',
    ]
  }
  if (msg.includes('contraseña incorrecta') || msg.includes('password')) {
    return [
      'Las contraseñas distinguen MAYÚSCULAS y minúsculas.',
      'Verifica que no haya espacios al inicio o al final.',
      'Si olvidaste tu contraseña, usa el botón "¿Olvidaste tu contraseña?" que aparece debajo del campo de contraseña.',
    ]
  }
  if (msg.includes('correo no registrado') || msg.includes('no registrado')) {
    return [
      'Verifica que el correo sea exactamente el que usaste al registrarte.',
      'Si nunca te has registrado, usa el enlace "Regístrate aquí" de abajo.',
      'Si crees que ya tienes cuenta, contacta al administrador.',
    ]
  }
  return []
})

// reCAPTCHA tokens expire after 2 minutes. Auto-reset at 110s so the user
// never unknowingly submits with an expired token.
let recaptchaResetTimer = null
let captchaInitInterval = null
let captchaInitTimeout = null

const recaptchaLoading = ref(true)  // true while widget is initializing
const recaptchaLoadError = ref(false) // true when widget failed to load

function clearCaptchaTimers() {
  clearInterval(captchaInitInterval)
  clearTimeout(captchaInitTimeout)
}

function scheduleRecaptchaReset() {
  clearTimeout(recaptchaResetTimer)
  recaptchaResetTimer = setTimeout(() => {
    if (recaptchaWidgetId.value !== null && window.grecaptcha) {
      window.grecaptcha.reset(recaptchaWidgetId.value)
      recaptchaToken.value = ''
    }
  }, 110_000)
}

function onRecaptchaSolved(token) {
  recaptchaToken.value = token
  scheduleRecaptchaReset()
}

function onRecaptchaExpired() {
  recaptchaToken.value = ''
  clearTimeout(recaptchaResetTimer)
  // Auto-reset to a clean widget so the user doesn't see the confusing
  // native red "expired" state — they just re-check the box.
  setTimeout(() => {
    if (recaptchaWidgetId.value !== null && window.grecaptcha) {
      window.grecaptcha.reset(recaptchaWidgetId.value)
    }
  }, 800)
}

// Shared init logic used by onMounted and retryRecaptcha.
// Polls for grecaptcha readiness with a 15-second timeout.
function initRecaptcha() {
  recaptchaLoading.value = true
  recaptchaLoadError.value = false
  // Clear container in case this is a retry
  const container = document.getElementById('recaptcha-container')
  if (container) container.innerHTML = ''
  recaptchaWidgetId.value = null

  captchaInitInterval = setInterval(() => {
    if (window.grecaptcha && window.grecaptcha.render) {
      clearCaptchaTimers()
      recaptchaLoading.value = false
      try {
        recaptchaWidgetId.value = window.grecaptcha.render('recaptcha-container', {
          sitekey: '6LdUjcgsAAAAAI0pgOSk7QMEmq-zjD4saihqzaa-',
          theme: themeStore.isDark ? 'dark' : 'light',
          callback: onRecaptchaSolved,
          'expired-callback': onRecaptchaExpired,
        })
      } catch (e) {
        recaptchaLoadError.value = true
        console.warn('reCAPTCHA render error:', e)
      }
    }
  }, 100)

  // After 15 s give up and let the user know
  captchaInitTimeout = setTimeout(() => {
    clearCaptchaTimers()
    if (recaptchaLoading.value) {
      recaptchaLoading.value = false
      recaptchaLoadError.value = true
    }
  }, 15000)
}

// Injects a fresh script tag (in case it never loaded) and re-initializes.
function retryRecaptcha() {
  if (!window.grecaptcha || !window.grecaptcha.render) {
    const old = document.querySelector('script[src*="recaptcha/api.js"]')
    if (old) old.remove()
    const script = document.createElement('script')
    script.src = 'https://www.google.com/recaptcha/api.js'
    script.async = true
    document.head.appendChild(script)
  }
  initRecaptcha()
}

onMounted(() => { initRecaptcha() })

onUnmounted(() => {
  clearCaptchaTimers()
  clearTimeout(recaptchaResetTimer)
})

async function handleLogin() {
  error.value = ''

  // Usar el token guardado en el callback; como fallback intentar getResponse()
  // por si el callback no se disparó en algún navegador móvil.
  let token = recaptchaToken.value
  if (!token && recaptchaWidgetId.value !== null && window.grecaptcha) {
    token = window.grecaptcha.getResponse(recaptchaWidgetId.value)
  }

  if (!token) {
    error.value = 'Completa la verificación reCAPTCHA. Si aparece un puzzle de imágenes, resuélvelo antes de continuar.'
    return
  }

  loading.value = true
  try {
    // Trim both fields to prevent invisible-whitespace credential mismatches
    const data = await auth.login(form.value.email.trim(), form.value.password.trim(), token)
    clearTimeout(recaptchaResetTimer)
    recaptchaToken.value = ''
    if (data.user?.role === 'admin') {
      router.push('/admin')
    } else {
      router.push('/checkin')
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Credenciales inválidas'
    recaptchaToken.value = ''
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

  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }

  50% {
    transform: translate(20px, -20px) scale(1.05);
  }
}

@keyframes float-delayed {

  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }

  50% {
    transform: translate(-30px, 30px) scale(1.1);
  }
}

.animate-float-slow {
  animation: float-slow 10s ease-in-out infinite;
}

.animate-float-delayed {
  animation: float-delayed 12s ease-in-out infinite;
}

@keyframes bounce-subtle {

  0%,
  100% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-10px);
  }
}

.animate-bounce-subtle {
  animation: bounce-subtle 4s ease-in-out infinite;
}

@keyframes in-scale {
  from {
    opacity: 0;
    transform: scale(0.9);
  }

  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-in-scale {
  animation: in-scale 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes slide-right {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }

  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.animate-slide-right {
  opacity: 0;
  animation: slide-right 0.5s ease-out forwards;
}

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-slide-up {
  opacity: 0;
  animation: slide-up 0.5s ease-out forwards;
}

@keyframes fade-in {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

.animate-fade-in {
  opacity: 0;
  animation: fade-in 0.8s ease-out forwards;
}

@keyframes shake {

  0%,
  100% {
    transform: translateX(0);
  }

  25% {
    transform: translateX(-5px);
  }

  75% {
    transform: translateX(5px);
  }
}

.animate-shake {
  animation: shake 0.4s ease-in-out;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
