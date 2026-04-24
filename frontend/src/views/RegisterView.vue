<template>
  <div class="min-h-screen min-h-dvh flex items-center justify-center px-4 py-12 relative overflow-hidden"
    style="background-color: #0b0e14;">
    <!-- Grid pattern background -->
    <div class="absolute inset-0 opacity-20 pointer-events-none"
      style="background-image: linear-gradient(rgba(255,255,255,0.05) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.05) 1px, transparent 1px); background-size: 40px 40px;">
    </div>
    
    <!-- Glow effects -->
    <div class="absolute -top-40 -left-40 w-96 h-96 rounded-full opacity-30 animate-pulse-slow pointer-events-none"
      style="background: radial-gradient(circle, #f97316, transparent 70%)"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 rounded-full opacity-20 animate-pulse-slow pointer-events-none"
      style="background: radial-gradient(circle, #8b5cf6, transparent 70%); animation-delay: 2s;"></div>

    <div class="relative w-full max-w-6xl grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-center">
      
      <!-- Left side: Branding -->
      <div class="hidden lg:flex flex-col items-start animate-in">
        <div class="mb-4">
          <!-- Logo SVG -->
          <svg width="120" height="120" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 80V30L40 50L60 30V80" stroke="#475569" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M40 80V40L50 50L60 40V80" stroke="#475569" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M10 65C30 65 40 55 50 45C60 35 70 25 90 20" stroke="#f97316" stroke-width="4" stroke-linecap="round"/>
            <path d="M80 15H95V30" stroke="#f97316" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M20 80H80" stroke="#f97316" stroke-width="6" stroke-linecap="round"/>
            <path d="M30 80V90 M50 80V90 M70 80V90" stroke="#f97316" stroke-width="4" stroke-linecap="round"/>
          </svg>
        </div>
        <h1 class="text-4xl font-black tracking-wider text-[#f97316] mb-0 leading-none">SOLUCIONES</h1>
        <h2 class="text-2xl tracking-widest text-[#64748b] mb-8 font-light">EMPRESARIALES</h2>
        <div class="text-5xl font-bold flex gap-3 items-center">
          <span class="text-white">Portal</span>
          <span class="text-[#f97316]">RRHH</span>
        </div>
      </div>

      <!-- Right side: Register Card -->
      <div class="w-full max-w-lg mx-auto animate-in" style="animation-delay: 0.1s;">
        <!-- Mobile Logo (hidden on desktop) -->
        <div class="lg:hidden flex flex-col items-center text-center mb-8">
          <svg width="80" height="80" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg" class="mb-2">
            <path d="M20 80V30L40 50L60 30V80" stroke="#475569" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M40 80V40L50 50L60 40V80" stroke="#475569" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M10 65C30 65 40 55 50 45C60 35 70 25 90 20" stroke="#f97316" stroke-width="4" stroke-linecap="round"/>
            <path d="M80 15H95V30" stroke="#f97316" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M20 80H80" stroke="#f97316" stroke-width="6" stroke-linecap="round"/>
          </svg>
          <div class="text-2xl font-bold flex gap-2 justify-center">
            <span class="text-white">Portal</span>
            <span class="text-[#f97316]">RRHH</span>
          </div>
        </div>

        <div class="rounded-2xl p-6 sm:p-10"
          style="background-color: #121621; border: 1px solid rgba(255,255,255,0.05); border-top: 4px solid #f97316; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);">
          
          <div class="text-center mb-6">
            <h2 class="text-2xl font-bold text-white mb-2">Crear Perfil</h2>
            <p class="text-sm text-[#94a3b8]">Ingresa tus datos para registrarte</p>
          </div>

          <form @submit.prevent="handleRegister" class="space-y-4">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <input v-model="form.first_name" type="text" class="w-full px-4 py-3 rounded-xl text-white text-center sm:text-left focus:outline-none transition-all placeholder-[#475569]" style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.05);" placeholder="Nombre(s)" required onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.05)'" />
              <input v-model="form.last_name" type="text" class="w-full px-4 py-3 rounded-xl text-white text-center sm:text-left focus:outline-none transition-all placeholder-[#475569]" style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.05);" placeholder="Apellido(s)" required onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.05)'" />
            </div>

            <input v-model="form.project_name" type="text" class="w-full px-4 py-3 rounded-xl text-white text-center focus:outline-none transition-all placeholder-[#475569]" style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.05);" placeholder="Nombre del Proyecto" required onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.05)'" />
            
            <input v-model="form.email" type="email" class="w-full px-4 py-3 rounded-xl text-white text-center focus:outline-none transition-all placeholder-[#475569]" style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.05);" placeholder="correo@ejemplo.com" required autocomplete="email" onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.05)'" />

            <div class="relative">
              <input v-model="form.password" :type="showPwd ? 'text' : 'password'" class="w-full px-4 py-3 rounded-xl text-white text-center tracking-[0.2em] focus:outline-none transition-all placeholder-[#475569]" style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.05);" placeholder="••••••••" required minlength="8" onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.05)'" />
              <button type="button" @click="showPwd = !showPwd" class="absolute right-4 top-1/2 -translate-y-1/2 transition-colors text-[#64748b] hover:text-white">
                <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
                <EyeIcon v-else class="w-5 h-5" />
              </button>
            </div>

            <div>
              <input v-model="form.confirm_password" type="password" class="w-full px-4 py-3 rounded-xl text-white text-center tracking-[0.2em] focus:outline-none transition-all placeholder-[#475569]" style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.05);" placeholder="Confirmar ••••••••" required onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.05)'" />
              <p v-if="passwordMismatch" class="text-rose-400 text-xs mt-1.5 flex items-center gap-1 justify-center">
                <ExclamationTriangleIcon class="w-3.5 h-3.5" /> Las contraseñas no coinciden
              </p>
            </div>

            <!-- Math CAPTCHA -->
            <div class="rounded-xl p-4 mt-2" style="background-color: rgba(249, 115, 22, 0.05); border: 1px solid rgba(249, 115, 22, 0.2);">
              <p class="text-xs text-[#94a3b8] mb-3 text-center">Verificación de seguridad</p>
              <div class="flex items-center justify-center gap-3">
                <div class="px-4 py-2 rounded-lg font-mono text-lg font-bold text-white bg-[#1e293b] border border-gray-700 select-none">
                  {{ captcha.a }} + {{ captcha.b }} =
                </div>
                <input v-model="captchaAnswer" type="number" class="w-20 px-3 py-2 rounded-lg text-white text-center focus:outline-none transition-all placeholder-[#475569]" style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.1);" placeholder="?" required onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.1)'" />
                <button type="button" @click="refreshCaptcha" class="text-[#64748b] hover:text-[#f97316] transition-colors p-2" title="Nuevo captcha">
                  <ArrowPathIcon class="w-5 h-5" />
                </button>
              </div>
              <p v-if="captchaError" class="text-rose-400 text-xs mt-2 flex items-center justify-center gap-1">
                <ExclamationTriangleIcon class="w-3.5 h-3.5" /> Respuesta incorrecta
              </p>
            </div>

            <div v-if="error" class="flex items-start gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-400 rounded-xl px-4 py-3 text-sm">
              <ExclamationTriangleIcon class="w-4 h-4 flex-shrink-0 mt-0.5" /><span>{{ error }}</span>
            </div>
            <div v-if="success" class="flex items-start gap-2 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 rounded-xl px-4 py-3 text-sm">
              <CheckCircleIcon class="w-4 h-4 flex-shrink-0 mt-0.5" /><span>¡Perfil creado! Redirigiendo...</span>
            </div>

            <button type="submit" class="w-full py-3 rounded-xl font-medium text-white transition-all active:scale-[0.98] mt-4 flex items-center justify-center gap-2" 
              style="background-color: #1a2235; border: 1px solid rgba(255,255,255,0.05);" :disabled="loading || passwordMismatch">
              <span v-if="loading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
              <span>{{ loading ? 'Creando perfil...' : 'Crear Perfil' }}</span>
            </button>
          </form>

          <div class="mt-6 text-center text-xs" style="color: #475569;">
            ¿Ya tienes cuenta?
            <router-link to="/login" class="text-[#f97316] font-semibold hover:text-orange-400 ml-1 transition-colors">
              Iniciar sesión
            </router-link>
            <p class="mt-5">MH Soluciones Empresariales © {{ new Date().getFullYear() }}</p>
          </div>
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