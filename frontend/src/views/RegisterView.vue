<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4 py-8">
    <div class="card w-full max-w-lg">
      <div class="text-center mb-8">
        <div class="text-5xl mb-3">👤</div>
        <h1 class="text-2xl font-bold text-blue-700">Crear Perfil</h1>
        <p class="text-gray-500 mt-1">Regístrate para comenzar a checar</p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nombre(s)</label>
            <input v-model="form.first_name" type="text" class="input-field" placeholder="Juan" required />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Apellido(s)</label>
            <input v-model="form.last_name" type="text" class="input-field" placeholder="Pérez García" required />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Nombre completo del proyecto</label>
          <input v-model="form.project_name" type="text" class="input-field" placeholder="Proyecto de Construcción Centro" required />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Correo electrónico</label>
          <input v-model="form.email" type="email" class="input-field" placeholder="correo@ejemplo.com" required />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Contraseña</label>
          <input v-model="form.password" type="password" class="input-field" placeholder="Mínimo 8 caracteres" required minlength="8" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Confirmar contraseña</label>
          <input v-model="form.confirm_password" type="password" class="input-field" placeholder="Repite tu contraseña" required />
          <p v-if="passwordMismatch" class="text-red-500 text-xs mt-1">Las contraseñas no coinciden</p>
        </div>

        <!-- Simple Math CAPTCHA -->
        <div class="bg-gray-50 rounded-xl p-4 border border-gray-200">
          <label class="block text-sm font-medium text-gray-700 mb-2">Verificación de seguridad</label>
          <div class="flex items-center gap-3">
            <span class="text-lg font-bold text-gray-800 bg-white px-3 py-2 rounded-lg border select-none">
              {{ captcha.a }} + {{ captcha.b }} = ?
            </span>
            <input
              v-model="captchaAnswer"
              type="number"
              class="input-field w-24"
              placeholder="?"
              required
            />
            <button type="button" @click="refreshCaptcha" class="text-blue-500 hover:text-blue-700 text-xl" title="Nuevo captcha">🔄</button>
          </div>
          <p v-if="captchaError" class="text-red-500 text-xs mt-1">Respuesta incorrecta</p>
        </div>

        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 rounded-xl px-4 py-3 text-sm">
          {{ error }}
        </div>
        <div v-if="success" class="bg-green-50 border border-green-200 text-green-700 rounded-xl px-4 py-3 text-sm">
          ¡Perfil creado! Redirigiendo al login...
        </div>

        <button type="submit" class="btn-primary w-full" :disabled="loading || passwordMismatch">
          <span v-if="loading">Creando perfil...</span>
          <span v-else>Crear Perfil</span>
        </button>
      </form>

      <div class="mt-6 text-center text-sm text-gray-600">
        ¿Ya tienes cuenta?
        <router-link to="/login" class="text-blue-600 font-semibold hover:underline ml-1">
          Iniciar sesión
        </router-link>
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
  first_name: '',
  last_name: '',
  project_name: '',
  email: '',
  password: '',
  confirm_password: ''
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
