<template>
  <div class="min-h-screen min-h-dvh flex items-center justify-center p-4 relative overflow-hidden bg-[#f8fafc]">
    
    <!-- Animated background elements -->
    <div class="absolute inset-0 z-0 overflow-hidden pointer-events-none">
      <div class="absolute w-[600px] h-[600px] rounded-full blur-[120px] opacity-20 bg-gradient-to-br from-blue-400 to-indigo-300 animate-float-slow -top-20 -left-20"></div>
      <div class="absolute w-[500px] h-[500px] rounded-full blur-[100px] opacity-15 bg-gradient-to-br from-purple-300 to-pink-200 animate-float-delayed -bottom-20 -right-20"></div>
      <div class="absolute inset-0 opacity-[0.03]" style="background-image: radial-gradient(#000 1px, transparent 1px); background-size: 30px 30px;"></div>
    </div>

    <!-- Main Card -->
    <div class="relative z-10 w-full max-w-[450px] bg-white/80 backdrop-blur-2xl rounded-[3rem] shadow-[0_20px_50px_rgba(0,0,0,0.05)] border border-white p-8 lg:p-12 animate-in-scale">
      
      <div class="text-center mb-10">
        <div class="inline-flex p-4 rounded-3xl bg-blue-50 text-blue-500 mb-6 animate-bounce-subtle">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
          </svg>
        </div>
        <h2 class="text-3xl font-extrabold text-slate-800 tracking-tight mb-2">Bienvenido</h2>
        <p class="text-slate-500 font-medium">Ingresa tus credenciales para continuar</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <div class="space-y-1 animate-slide-right" style="animation-delay: 0.1s">
          <label class="block text-sm font-bold text-slate-700 ml-1">Correo Electrónico</label>
          <div class="relative group">
            <input v-model="form.email" type="email" 
              class="w-full px-6 py-4 rounded-2xl text-slate-800 bg-slate-100/50 border border-transparent focus:border-blue-400 focus:bg-white focus:ring-4 focus:ring-blue-100 transition-all duration-300 outline-none placeholder-slate-400 font-medium"
              placeholder="ejemplo@correo.com" required autocomplete="email" />
          </div>
        </div>
        
        <div class="space-y-1 animate-slide-right" style="animation-delay: 0.2s">
          <label class="block text-sm font-bold text-slate-700 ml-1">Contraseña</label>
          <div class="relative group">
            <input v-model="form.password" :type="showPwd ? 'text' : 'password'" 
              class="w-full px-6 py-4 rounded-2xl text-slate-800 bg-slate-100/50 border border-transparent focus:border-blue-400 focus:bg-white focus:ring-4 focus:ring-blue-100 transition-all duration-300 outline-none placeholder-slate-400 font-medium tracking-tight"
              placeholder="••••••••" required autocomplete="current-password" />
            <button type="button" @click="showPwd = !showPwd"
              class="absolute right-4 top-1/2 -translate-y-1/2 p-2 rounded-xl text-slate-400 hover:text-blue-500 hover:bg-blue-50 transition-colors">
              <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
              <EyeIcon v-else class="w-5 h-5" />
            </button>
          </div>
        </div>

        <!-- reCAPTCHA v2 Container -->
        <div class="flex justify-center py-2 animate-fade-in" style="animation-delay: 0.3s">
          <div id="recaptcha-container" class="shadow-sm rounded-lg overflow-hidden border border-slate-100"></div>
        </div>

        <Transition name="fade">
          <div v-if="error" class="flex items-center gap-3 bg-red-50 border border-red-100 text-red-500 rounded-2xl px-5 py-4 text-sm font-bold animate-shake">
            <ExclamationTriangleIcon class="w-5 h-5 flex-shrink-0" />
            <span>{{ error }}</span>
          </div>
        </Transition>

        <button type="submit" class="relative w-full py-4 rounded-2xl font-black text-white overflow-hidden group transition-all duration-300 hover:shadow-2xl hover:shadow-blue-200 active:scale-95 animate-slide-up" :disabled="loading" style="animation-delay: 0.4s">
          <div class="absolute inset-0 bg-gradient-to-r from-blue-500 to-indigo-600 transition-transform duration-500 group-hover:scale-110"></div>
          <div class="relative flex items-center justify-center gap-3">
            <span v-if="loading" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span class="tracking-widest uppercase text-sm">{{ loading ? 'Cargando...' : 'Entrar ahora' }}</span>
          </div>
        </button>
      </form>

      <div class="mt-10 text-center animate-fade-in" style="animation-delay: 0.5s">
        <p class="text-slate-400 font-medium text-sm">
          ¿Aún no tienes cuenta? 
          <router-link to="/register" class="text-blue-600 font-bold hover:text-blue-700 hover:underline transition-all">
            Crea una aquí
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'

const auth = useAuthStore()
const router = useRouter()
const loading = ref(false)
const error = ref('')
const showPwd = ref(false)
const form = ref({ email: '', password: '' })
const recaptchaWidgetId = ref(null)

onMounted(() => {
  const checkRecaptcha = setInterval(() => {
    if (window.grecaptcha && window.grecaptcha.render) {
      clearInterval(checkRecaptcha)
      try {
        recaptchaWidgetId.value = window.grecaptcha.render('recaptcha-container', {
          sitekey: '6LdUjcgsAAAAAI0pgOSk7QMEmq-zjD4saihqzaa-', // Tu Clave pública
          theme: 'light'
        })
      } catch (e) {
        console.warn('Recaptcha error:', e)
      }
    }
  }, 100)
})

async function handleLogin() {
  error.value = ''
  
  let recaptchaToken = ''
  if (recaptchaWidgetId.value !== null) {
    recaptchaToken = window.grecaptcha.getResponse(recaptchaWidgetId.value)
  }

  if (!recaptchaToken) {
    error.value = 'Por favor, verifica que no eres un robot.'
    return
  }

  loading.value = true
  try {
    const data = await auth.login(form.value.email, form.value.password, recaptchaToken)
    if (data.user?.role === 'admin') {
      router.push('/admin')
    } else {
      router.push('/checkin')
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Credenciales inválidas'
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

@keyframes in-scale {
  from { opacity: 0; transform: scale(0.9); }
  to { opacity: 1; transform: scale(1); }
}
.animate-in-scale { animation: in-scale 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards; }

@keyframes slide-right {
  from { opacity: 0; transform: translateX(-20px); }
  to { opacity: 1; transform: translateX(0); }
}
.animate-slide-right { opacity: 0; animation: slide-right 0.5s ease-out forwards; }

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
