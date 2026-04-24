<template>
  <div class="min-h-screen min-h-dvh flex items-center justify-center px-4 py-8 relative overflow-hidden"
    style="background-color: #0b0e14;">
    <!-- Grid pattern background -->
    <div class="absolute inset-0 opacity-20"
      style="background-image: linear-gradient(rgba(255,255,255,0.05) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.05) 1px, transparent 1px); background-size: 40px 40px;">
    </div>
    
    <!-- Glow effects -->
    <div class="absolute -top-40 -left-40 w-96 h-96 rounded-full opacity-30 animate-pulse-slow"
      style="background: radial-gradient(circle, #f97316, transparent 70%)"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 rounded-full opacity-20 animate-pulse-slow"
      style="background: radial-gradient(circle, #8b5cf6, transparent 70%); animation-delay: 2s;"></div>

    <div class="relative w-full max-w-5xl grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-24 items-center">
      
      <!-- Left side: Branding -->
      <div class="flex flex-col items-center lg:items-start text-center lg:text-left animate-in">
        <div class="mb-4">
          <!-- Logo SVG roughly matching the MH Soluciones image -->
          <svg width="120" height="120" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 80V30L40 50L60 30V80" stroke="#475569" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M40 80V40L50 50L60 40V80" stroke="#475569" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M10 65C30 65 40 55 50 45C60 35 70 25 90 20" stroke="#f97316" stroke-width="4" stroke-linecap="round"/>
            <path d="M80 15H95V30" stroke="#f97316" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M20 80H80" stroke="#f97316" stroke-width="6" stroke-linecap="round"/>
            <path d="M30 80V90 M50 80V90 M70 80V90" stroke="#f97316" stroke-width="4" stroke-linecap="round"/>
          </svg>
        </div>
        <h1 class="text-3xl font-black tracking-wider text-[#f97316] mb-0 leading-none">SOLUCIONES</h1>
        <h2 class="text-xl tracking-widest text-[#64748b] mb-8 font-light">EMPRESARIALES</h2>
        <div class="text-4xl font-bold flex gap-3 items-center">
          <span class="text-white">Portal</span>
          <span class="text-[#f97316]">RRHH</span>
        </div>
      </div>

      <!-- Right side: Login Card -->
      <div class="w-full max-w-md mx-auto animate-in" style="animation-delay: 0.1s;">
        <div class="rounded-2xl p-8 lg:p-10"
          style="background-color: #121621; border: 1px solid rgba(255,255,255,0.05); border-top: 4px solid #f97316; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);">
          
          <div class="text-center mb-8">
            <h2 class="text-2xl font-bold text-white mb-2">Bienvenido</h2>
            <p class="text-sm text-[#94a3b8]">Ingresa tus credenciales para continuar</p>
          </div>

          <!-- Decorative dots -->
          <div class="flex justify-center gap-2 mb-8">
            <div v-for="i in 6" :key="i" class="w-2.5 h-2.5 rounded-full" style="background-color: #1e293b;"></div>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-5">
            <div>
              <input v-model="form.email" type="email" 
                class="w-full px-4 py-3 rounded-xl text-white text-center focus:outline-none transition-all" 
                style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.05); border-color: rgba(255,255,255,0.05);"
                placeholder="correo@ejemplo.com" required autocomplete="email" 
                onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.05)'" />
            </div>
            <div class="relative">
              <input v-model="form.password" :type="showPwd ? 'text' : 'password'" 
                class="w-full px-4 py-3 rounded-xl text-white text-center tracking-[0.2em] focus:outline-none transition-all" 
                style="background-color: #1e293b; border: 1px solid rgba(255,255,255,0.05); border-color: rgba(255,255,255,0.05);"
                placeholder="••••••••" required autocomplete="current-password"
                onfocus="this.style.borderColor='#f97316'" onblur="this.style.borderColor='rgba(255,255,255,0.05)'" />
              <button type="button" @click="showPwd = !showPwd"
                class="absolute right-4 top-1/2 -translate-y-1/2 transition-colors text-[#64748b] hover:text-white">
                <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
                <EyeIcon v-else class="w-5 h-5" />
              </button>
            </div>

            <div v-if="error" class="flex items-start gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-400 rounded-xl px-4 py-3 text-sm">
              <ExclamationTriangleIcon class="w-4 h-4 flex-shrink-0 mt-0.5" />
              <span>{{ error }}</span>
            </div>

            <button type="submit" class="w-full py-3 rounded-xl font-medium text-white transition-all active:scale-[0.98] mt-4 flex items-center justify-center gap-2" 
              style="background-color: #1a2235; border: 1px solid rgba(255,255,255,0.05);" :disabled="loading">
              <span v-if="loading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
              <span>{{ loading ? 'Iniciando...' : 'Ingresar' }}</span>
            </button>
          </form>

          <div class="mt-8 text-center text-xs" style="color: #475569;">
            ¿No tienes cuenta?
            <router-link to="/register" class="text-[#f97316] font-semibold hover:text-orange-400 ml-1 transition-colors">
              Crear perfil
            </router-link>
            <p class="mt-6">MH Soluciones Empresariales © {{ new Date().getFullYear() }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import ThemeToggle from '@/components/ThemeToggle.vue'

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

