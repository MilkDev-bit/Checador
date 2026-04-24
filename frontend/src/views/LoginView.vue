<template>
  <div class="min-h-screen min-h-dvh flex items-center justify-center p-4 relative overflow-hidden"
    style="background-color: #030712;">
    
    <!-- Dynamic Animated Background -->
    <div class="absolute inset-0 z-0 overflow-hidden pointer-events-none">
      <!-- Grid overlay -->
      <div class="absolute inset-0 opacity-10"
        style="background-image: linear-gradient(rgba(255,255,255,0.08) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.08) 1px, transparent 1px); background-size: 60px 60px;">
      </div>
      <!-- Floating glowing orbs -->
      <div class="absolute w-[500px] h-[500px] rounded-full blur-[100px] opacity-40 bg-gradient-to-br from-cyan-400 to-blue-600 animate-orb-1"></div>
      <div class="absolute w-[600px] h-[600px] rounded-full blur-[120px] opacity-30 bg-gradient-to-br from-purple-500 to-indigo-600 animate-orb-2"></div>
      <div class="absolute w-[400px] h-[400px] rounded-full blur-[80px] opacity-20 bg-gradient-to-br from-fuchsia-500 to-pink-500 animate-orb-3"></div>
    </div>

    <!-- Main Container -->
    <div class="relative z-10 w-full max-w-6xl grid grid-cols-1 lg:grid-cols-5 gap-8 lg:gap-0 items-stretch bg-white/5 backdrop-blur-xl rounded-[2.5rem] border border-white/10 shadow-2xl overflow-hidden animate-in-up">
      
      <!-- Left side: Interactive Branding area -->
      <div class="lg:col-span-2 relative p-8 lg:p-12 flex flex-col justify-center items-center text-center overflow-hidden border-b lg:border-b-0 lg:border-r border-white/10"
        style="background: radial-gradient(circle at top left, rgba(255,255,255,0.05), transparent);">
        <div class="relative z-10 animate-float">
          <!-- The new generated 3D logo -->
          <img src="/brand_logo.png" alt="Logo" class="w-40 h-40 lg:w-56 lg:h-56 object-contain drop-shadow-[0_0_30px_rgba(6,182,212,0.6)] mb-6 transition-transform duration-500 hover:scale-110" />
        </div>
        <div class="relative z-10 mt-4">
          <h1 class="text-3xl lg:text-4xl font-black text-transparent bg-clip-text bg-gradient-to-r from-cyan-300 to-purple-400 tracking-tight mb-3">
            PASELISTA
          </h1>
          <p class="text-slate-400 text-sm lg:text-base font-light max-w-[250px] mx-auto">
            El sistema de control de asistencia de nueva generación.
          </p>
        </div>
      </div>

      <!-- Right side: Login Form -->
      <div class="lg:col-span-3 p-8 lg:p-16 flex flex-col justify-center relative bg-[#0f172a]/50">
        <div class="max-w-sm w-full mx-auto">
          <div class="mb-10 text-center lg:text-left">
            <h2 class="text-3xl font-bold text-white mb-2">Bienvenido de vuelta</h2>
            <p class="text-slate-400">Ingresa a tu panel de control</p>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-6">
            <div class="group">
              <label class="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2 group-focus-within:text-cyan-400 transition-colors">Correo Electrónico</label>
              <input v-model="form.email" type="email" 
                class="w-full px-5 py-4 rounded-2xl text-white bg-slate-900/50 border border-slate-700/50 focus:border-cyan-400 focus:bg-slate-900/80 focus:ring-4 focus:ring-cyan-500/10 transition-all duration-300 outline-none placeholder-slate-600 shadow-inner"
                placeholder="usuario@correo.com" required autocomplete="email" />
            </div>
            
            <div class="group">
              <label class="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2 group-focus-within:text-purple-400 transition-colors">Contraseña</label>
              <div class="relative">
                <input v-model="form.password" :type="showPwd ? 'text' : 'password'" 
                  class="w-full px-5 py-4 rounded-2xl text-white bg-slate-900/50 border border-slate-700/50 focus:border-purple-400 focus:bg-slate-900/80 focus:ring-4 focus:ring-purple-500/10 transition-all duration-300 outline-none placeholder-slate-600 shadow-inner tracking-widest"
                  placeholder="••••••••" required autocomplete="current-password" />
                <button type="button" @click="showPwd = !showPwd"
                  class="absolute right-4 top-1/2 -translate-y-1/2 p-2 rounded-lg text-slate-500 hover:text-white hover:bg-white/10 transition-colors">
                  <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
                  <EyeIcon v-else class="w-5 h-5" />
                </button>
              </div>
            </div>

            <!-- reCAPTCHA v2 Container -->
            <div class="flex justify-center mt-2 overflow-hidden rounded-xl">
              <div id="recaptcha-container"></div>
            </div>

            <Transition name="fade">
              <div v-if="error" class="flex items-start gap-3 bg-rose-500/10 border border-rose-500/20 text-rose-400 rounded-2xl px-5 py-4 text-sm animate-shake">
                <ExclamationTriangleIcon class="w-5 h-5 flex-shrink-0" />
                <span>{{ error }}</span>
              </div>
            </Transition>

            <button type="submit" class="relative w-full py-4 rounded-2xl font-bold text-white overflow-hidden group transition-all duration-300 hover:scale-[1.02] active:scale-[0.98] shadow-lg shadow-cyan-500/25" :disabled="loading">
              <div class="absolute inset-0 bg-gradient-to-r from-cyan-500 to-purple-600 transition-transform duration-500 group-hover:scale-110"></div>
              <div class="absolute inset-0 opacity-0 group-hover:opacity-20 bg-white transition-opacity duration-300"></div>
              <div class="relative flex items-center justify-center gap-3">
                <span v-if="loading" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                <span class="tracking-wide">{{ loading ? 'AUTENTICANDO...' : 'INICIAR SESIÓN' }}</span>
              </div>
            </button>
          </form>

          <div class="mt-10 text-center text-sm text-slate-500">
            ¿Nuevo en la plataforma?
            <router-link to="/register" class="text-cyan-400 font-semibold hover:text-cyan-300 hover:underline transition-all ml-1">
              Crear una cuenta
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
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
          sitekey: '6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI', // Google Test Site Key
          theme: 'dark'
        })
      } catch (e) {
        // Handle double render gracefully
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
    error.value = 'Por favor, completa la verificación de seguridad reCAPTCHA.'
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
    error.value = e.response?.data?.error || 'Error al iniciar sesión'
    if (recaptchaWidgetId.value !== null) {
      window.grecaptcha.reset(recaptchaWidgetId.value)
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@keyframes float {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-15px) scale(1.02); }
}
.animate-float {
  animation: float 6s ease-in-out infinite;
}

@keyframes orb-1 {
  0%, 100% { transform: translate(-10%, -10%) scale(1); }
  50% { transform: translate(10%, 20%) scale(1.1); }
}
@keyframes orb-2 {
  0%, 100% { transform: translate(80%, -20%) scale(1); }
  50% { transform: translate(60%, 10%) scale(0.9); }
}
@keyframes orb-3 {
  0%, 100% { transform: translate(40%, 80%) scale(1); }
  50% { transform: translate(20%, 60%) scale(1.2); }
}
.animate-orb-1 { animation: orb-1 15s ease-in-out infinite alternate; }
.animate-orb-2 { animation: orb-2 20s ease-in-out infinite alternate-reverse; }
.animate-orb-3 { animation: orb-3 18s ease-in-out infinite alternate; }

@keyframes in-up {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-in-up {
  animation: in-up 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}
.animate-shake {
  animation: shake 0.4s ease-in-out;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
