<template>
  <div class="min-h-screen min-h-dvh flex items-center justify-center p-4 relative overflow-hidden bg-[#f8fafc]">
    
    <!-- Background accents -->
    <div class="absolute inset-0 z-0 overflow-hidden pointer-events-none">
      <div class="absolute w-[800px] h-[800px] rounded-full blur-[150px] opacity-15 bg-gradient-to-tr from-cyan-200 to-blue-200 animate-float-slow -bottom-40 -left-40"></div>
      <div class="absolute w-[600px] h-[600px] rounded-full blur-[120px] opacity-10 bg-gradient-to-bl from-indigo-300 to-purple-200 animate-float-delayed -top-40 -right-40"></div>
    </div>

    <!-- Container -->
    <div class="relative z-10 w-full max-w-[600px] bg-white shadow-[0_30px_60px_rgba(0,0,0,0.05)] rounded-[3rem] border border-white p-8 lg:p-12 animate-in-up">
      
      <div class="text-center mb-10">
        <h2 class="text-3xl font-black text-slate-800 tracking-tight mb-2">Crear Cuenta</h2>
        <p class="text-slate-500 font-medium">Únete a nosotros completando tus datos</p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-5">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <div class="space-y-1 animate-slide-in" style="animation-delay: 0.1s">
            <label class="block text-sm font-bold text-slate-700 ml-1">Nombre(s)</label>
            <input v-model="form.first_name" type="text" 
              class="w-full px-5 py-4 rounded-2xl text-slate-800 bg-slate-50 border border-transparent focus:border-blue-400 focus:bg-white focus:ring-4 focus:ring-blue-100 transition-all duration-300 outline-none placeholder-slate-400 font-medium"
              placeholder="Ej. Carlos" required />
          </div>
          <div class="space-y-1 animate-slide-in" style="animation-delay: 0.2s">
            <label class="block text-sm font-bold text-slate-700 ml-1">Apellido(s)</label>
            <input v-model="form.last_name" type="text" 
              class="w-full px-5 py-4 rounded-2xl text-slate-800 bg-slate-50 border border-transparent focus:border-blue-400 focus:bg-white focus:ring-4 focus:ring-blue-100 transition-all duration-300 outline-none placeholder-slate-400 font-medium"
              placeholder="Ej. García" required />
          </div>
        </div>

        <div class="space-y-1 animate-slide-in" style="animation-delay: 0.3s">
          <label class="block text-sm font-bold text-slate-700 ml-1">Proyecto</label>
          <input v-model="form.project_name" type="text" 
            class="w-full px-5 py-4 rounded-2xl text-slate-800 bg-slate-50 border border-transparent focus:border-blue-400 focus:bg-white focus:ring-4 focus:ring-blue-100 transition-all duration-300 outline-none placeholder-slate-400 font-medium"
            placeholder="Nombre de tu empresa o sucursal" required />
        </div>
        
        <div class="space-y-1 animate-slide-in" style="animation-delay: 0.4s">
          <label class="block text-sm font-bold text-slate-700 ml-1">Correo Electrónico</label>
          <input v-model="form.email" type="email" 
            class="w-full px-5 py-4 rounded-2xl text-slate-800 bg-slate-50 border border-transparent focus:border-blue-400 focus:bg-white focus:ring-4 focus:ring-blue-100 transition-all duration-300 outline-none placeholder-slate-400 font-medium"
            placeholder="usuario@correo.com" required autocomplete="email" />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <div class="space-y-1 animate-slide-in" style="animation-delay: 0.5s">
            <label class="block text-sm font-bold text-slate-700 ml-1">Contraseña</label>
            <div class="relative">
              <input v-model="form.password" :type="showPwd ? 'text' : 'password'" 
                class="w-full px-5 py-4 rounded-2xl text-slate-800 bg-slate-50 border border-transparent focus:border-blue-400 focus:bg-white focus:ring-4 focus:ring-blue-100 transition-all duration-300 outline-none placeholder-slate-400 font-medium tracking-tight"
                placeholder="••••••••" required minlength="8" />
              <button type="button" @click="showPwd = !showPwd"
                class="absolute right-4 top-1/2 -translate-y-1/2 p-2 rounded-xl text-slate-400 hover:text-blue-500 hover:bg-blue-50 transition-colors">
                <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
                <EyeIcon v-else class="w-5 h-5" />
              </button>
            </div>
          </div>

          <div class="space-y-1 animate-slide-in" style="animation-delay: 0.6s">
            <label class="block text-sm font-bold text-slate-700 ml-1">Confirmar</label>
            <input v-model="form.confirm_password" type="password" 
              class="w-full px-5 py-4 rounded-2xl text-slate-800 bg-slate-50 border border-transparent focus:border-blue-400 focus:bg-white focus:ring-4 focus:ring-blue-100 transition-all duration-300 outline-none placeholder-slate-400 font-medium tracking-tight"
              placeholder="••••••••" required />
          </div>
        </div>
        
        <Transition name="fade">
          <p v-if="passwordMismatch" class="text-red-500 text-xs font-bold flex items-center justify-center gap-1.5 animate-shake">
            <ExclamationTriangleIcon class="w-4 h-4" /> Las contraseñas no coinciden
          </p>
        </Transition>

        <div class="flex justify-center py-2 animate-fade-in" style="animation-delay: 0.7s">
          <div id="recaptcha-container" class="shadow-sm rounded-lg overflow-hidden border border-slate-100"></div>
        </div>

        <Transition name="fade">
          <div v-if="error" class="bg-red-50 border border-red-100 text-red-500 rounded-2xl px-5 py-4 text-sm font-bold text-center animate-shake">
            {{ error }}
          </div>
        </Transition>
        <Transition name="fade">
          <div v-if="success" class="bg-emerald-50 border border-emerald-100 text-emerald-600 rounded-2xl px-5 py-4 text-sm font-bold text-center">
            ¡Cuenta creada con éxito! Redirigiendo...
          </div>
        </Transition>

        <button type="submit" class="relative w-full py-4 rounded-2xl font-black text-white overflow-hidden group transition-all duration-300 hover:shadow-2xl hover:shadow-blue-200 active:scale-95 animate-slide-up" :disabled="loading || passwordMismatch" style="animation-delay: 0.8s">
          <div class="absolute inset-0 bg-gradient-to-r from-indigo-500 to-blue-600 transition-transform duration-500 group-hover:scale-110"></div>
          <div class="relative flex items-center justify-center gap-3">
            <span v-if="loading" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span class="tracking-widest uppercase text-sm">{{ loading ? 'Procesando...' : 'Registrarme' }}</span>
          </div>
        </button>
      </form>

      <div class="mt-8 text-center animate-fade-in" style="animation-delay: 0.9s">
        <p class="text-slate-400 font-medium text-sm">
          ¿Ya tienes una cuenta? 
          <router-link to="/login" class="text-blue-600 font-bold hover:text-blue-700 hover:underline transition-all">
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
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon, CheckCircleIcon } from '@heroicons/vue/24/outline'

const auth = useAuthStore()
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
          theme: 'light'
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
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, -30px); }
}
@keyframes float-delayed {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-40px, 40px); }
}
.animate-float-slow { animation: float-slow 15s ease-in-out infinite; }
.animate-float-delayed { animation: float-delayed 18s ease-in-out infinite; }

@keyframes in-up {
  from { opacity: 0; transform: translateY(40px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-in-up { animation: in-up 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards; }

@keyframes slide-in {
  from { opacity: 0; transform: translateX(-30px); }
  to { opacity: 1; transform: translateX(0); }
}
.animate-slide-in { opacity: 0; animation: slide-in 0.6s ease-out forwards; }

@keyframes slide-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-slide-up { opacity: 0; animation: slide-up 0.5s ease-out forwards; }

@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}
.animate-fade-in { opacity: 0; animation: fade-in 1s ease-out forwards; }

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}
.animate-shake { animation: shake 0.4s ease-in-out; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>