<template>
  <div class="min-h-screen min-h-dvh pb-16" style="background-color: var(--bg);">

    <!-- Header -->
    <header class="sticky top-0 z-20 px-4 py-3 flex items-center gap-3"
      style="background: var(--nav-bg); backdrop-filter: blur(14px); border-bottom: 1px solid var(--border-subtle);">
      <router-link to="/checkin"
        class="w-9 h-9 rounded-xl flex items-center justify-center transition-all"
        style="color: var(--text-muted);"
        @mouseenter="e => e.currentTarget.style.background='var(--input-bg)'"
        @mouseleave="e => e.currentTarget.style.background='transparent'">
        <ArrowLeftIcon class="w-5 h-5" />
      </router-link>
      <h1 class="font-bold text-base flex-1" style="color: var(--text);">Mi Perfil</h1>
      <ThemeToggle />
    </header>

    <div class="max-w-lg mx-auto px-4 pt-6 space-y-5">

      <!-- ── Avatar card ─────────────────────────────────── -->
      <div class="glass-card p-6 animate-in">
        <!-- Gradient banner -->
        <div class="h-24 rounded-2xl mb-0 -mx-0 relative overflow-hidden"
          style="background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #06b6d4 100%);">
          <div class="absolute inset-0 opacity-20"
            style="background-image: radial-gradient(#fff 1px, transparent 1px); background-size: 18px 18px;"></div>
        </div>

        <!-- Avatar overlapping banner -->
        <div class="flex items-end gap-4 -mt-10 mb-4">
          <div class="relative flex-shrink-0">
            <div class="w-20 h-20 rounded-2xl ring-4 overflow-hidden cursor-pointer group transition-transform hover:scale-105"
              style="ring-color: var(--card-bg); background: linear-gradient(135deg, #6366f1, #8b5cf6);"
              @click="triggerAvatarPick">
              <img v-if="previewUrl || auth.user?.avatar_url"
                :src="previewUrl || auth.user?.avatar_url"
                class="w-full h-full object-cover" alt="Avatar" />
              <div v-else class="w-full h-full flex items-center justify-center text-white text-2xl font-black">
                {{ initials }}
              </div>
              <!-- Overlay on hover -->
              <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center rounded-2xl">
                <CameraIcon class="w-6 h-6 text-white" />
              </div>
            </div>
            <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onFileChange" />
          </div>
          <div class="mb-1 min-w-0">
            <p class="font-black text-lg truncate" style="color: var(--text);">{{ fullName }}</p>
            <p class="text-sm truncate" style="color: var(--text-muted);">{{ auth.user?.email }}</p>
            <span class="badge mt-1" :class="auth.isAdmin ? 'badge-blue' : 'badge-gray'">
              {{ auth.isAdmin ? 'Administrador' : 'Usuario' }}
            </span>
          </div>
        </div>

        <!-- Avatar action buttons -->
        <Transition name="fade">
          <div v-if="previewUrl" class="flex gap-3 mt-1">
            <button @click="saveAvatar" :disabled="savingAvatar" class="btn btn-primary btn-sm flex-1">
              <span v-if="savingAvatar" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
              <CheckIcon v-else class="w-4 h-4" />
              {{ savingAvatar ? 'Guardando...' : 'Guardar foto' }}
            </button>
            <button @click="cancelAvatar" class="btn btn-secondary btn-sm px-4">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </Transition>

        <Transition name="fade">
          <div v-if="avatarSuccess" class="mt-3 flex items-center gap-2 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 rounded-xl px-4 py-2.5 text-sm font-semibold">
            <CheckCircleIcon class="w-4 h-4 flex-shrink-0" /> Foto de perfil actualizada
          </div>
        </Transition>
        <Transition name="fade">
          <div v-if="avatarError" class="mt-3 flex items-center gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-400 rounded-xl px-4 py-2.5 text-sm font-semibold">
            <ExclamationTriangleIcon class="w-4 h-4 flex-shrink-0" /> {{ avatarError }}
          </div>
        </Transition>
      </div>

      <!-- ── Change email card ───────────────────────────── -->
      <div class="glass-card p-6 animate-in" style="animation-delay: 0.06s">
        <div class="flex items-center gap-3 mb-5">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
            style="background: rgba(99,102,241,0.12); border: 1px solid rgba(99,102,241,0.2);">
            <EnvelopeIcon class="w-5 h-5 text-brand-400" />
          </div>
          <div>
            <h2 class="font-bold text-base" style="color: var(--text);">Cambiar correo</h2>
            <p class="text-xs" style="color: var(--text-muted);">Actual: {{ auth.user?.email }}</p>
          </div>
        </div>

        <form @submit.prevent="handleEmailUpdate" class="space-y-4">
          <div>
            <label class="input-label">Nuevo correo electrónico</label>
            <input v-model="emailForm.newEmail" type="email"
              class="input" placeholder="nuevo@correo.com" required autocomplete="email" />
          </div>
          <Transition name="fade">
            <div v-if="emailMsg.text" class="flex items-center gap-2 rounded-xl px-4 py-2.5 text-sm font-semibold"
              :class="emailMsg.ok
                ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400'
                : 'bg-rose-500/10 border border-rose-500/20 text-rose-400'">
              <CheckCircleIcon v-if="emailMsg.ok" class="w-4 h-4 flex-shrink-0" />
              <ExclamationTriangleIcon v-else class="w-4 h-4 flex-shrink-0" />
              {{ emailMsg.text }}
            </div>
          </Transition>
          <button type="submit" :disabled="savingEmail" class="btn btn-primary btn-md w-full">
            <span v-if="savingEmail" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <CheckIcon v-else class="w-4 h-4" />
            {{ savingEmail ? 'Guardando...' : 'Actualizar correo' }}
          </button>
        </form>
      </div>

      <!-- ── Change password card ───────────────────────── -->
      <div class="glass-card p-6 animate-in" style="animation-delay: 0.12s">
        <div class="flex items-center gap-3 mb-5">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
            style="background: rgba(139,92,246,0.12); border: 1px solid rgba(139,92,246,0.2);">
            <LockClosedIcon class="w-5 h-5 text-violet-400" />
          </div>
          <div>
            <h2 class="font-bold text-base" style="color: var(--text);">Cambiar contraseña</h2>
            <p class="text-xs" style="color: var(--text-muted);">Elige una contraseña segura</p>
          </div>
        </div>

        <form @submit.prevent="handlePasswordUpdate" class="space-y-4">
          <div>
            <label class="input-label">Contraseña actual</label>
            <div class="relative">
              <input v-model="pwdForm.current" :type="showCurrent ? 'text' : 'password'"
                class="input pr-12" placeholder="••••••••" required />
              <button type="button" @click="showCurrent = !showCurrent"
                class="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg transition-colors"
                style="color: var(--text-muted);">
                <EyeSlashIcon v-if="showCurrent" class="w-4 h-4" />
                <EyeIcon v-else class="w-4 h-4" />
              </button>
            </div>
          </div>
          <div>
            <label class="input-label">Nueva contraseña</label>
            <div class="relative">
              <input v-model="pwdForm.newPwd" :type="showNew ? 'text' : 'password'"
                class="input pr-12" placeholder="Mínimo 8 caracteres" required minlength="8" />
              <button type="button" @click="showNew = !showNew"
                class="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg transition-colors"
                style="color: var(--text-muted);">
                <EyeSlashIcon v-if="showNew" class="w-4 h-4" />
                <EyeIcon v-else class="w-4 h-4" />
              </button>
            </div>
            <!-- Strength bar -->
            <div v-if="pwdForm.newPwd" class="mt-2 flex items-center gap-2">
              <div class="flex gap-1 flex-1">
                <div v-for="n in 4" :key="n" class="h-1 flex-1 rounded-full transition-all duration-300"
                  :style="n <= pwdStrength.level
                    ? `background: ${pwdStrength.color}`
                    : 'background: var(--input-border)'"></div>
              </div>
              <span class="text-xs font-semibold" :style="`color: ${pwdStrength.color}`">{{ pwdStrength.label }}</span>
            </div>
          </div>
          <div>
            <label class="input-label">Confirmar nueva contraseña</label>
            <input v-model="pwdForm.confirm" type="password"
              class="input" :class="pwdForm.confirm && pwdMismatch ? 'ring-2 ring-rose-500/30 border-rose-500/40' : ''"
              placeholder="••••••••" required />
            <Transition name="fade">
              <p v-if="pwdForm.confirm && pwdMismatch" class="text-xs text-rose-400 mt-1 font-semibold flex items-center gap-1">
                <ExclamationTriangleIcon class="w-3.5 h-3.5" /> Las contraseñas no coinciden
              </p>
            </Transition>
          </div>
          <Transition name="fade">
            <div v-if="pwdMsg.text" class="flex items-center gap-2 rounded-xl px-4 py-2.5 text-sm font-semibold"
              :class="pwdMsg.ok
                ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400'
                : 'bg-rose-500/10 border border-rose-500/20 text-rose-400'">
              <CheckCircleIcon v-if="pwdMsg.ok" class="w-4 h-4 flex-shrink-0" />
              <ExclamationTriangleIcon v-else class="w-4 h-4 flex-shrink-0" />
              {{ pwdMsg.text }}
            </div>
          </Transition>
          <button type="submit" :disabled="savingPwd || pwdMismatch" class="btn w-full btn-md text-white font-bold"
            style="background: linear-gradient(135deg, #7c3aed, #6366f1); box-shadow: 0 4px 15px rgba(124,58,237,0.3);">
            <span v-if="savingPwd" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <LockClosedIcon v-else class="w-4 h-4" />
            {{ savingPwd ? 'Guardando...' : 'Cambiar contraseña' }}
          </button>
        </form>
      </div>

      <!-- ── Account info card ─────────────────────────── -->
      <div class="glass-card p-6 animate-in" style="animation-delay: 0.18s">
        <div class="flex items-center gap-3 mb-5">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
            style="background: rgba(6,182,212,0.12); border: 1px solid rgba(6,182,212,0.2);">
            <InformationCircleIcon class="w-5 h-5 text-cyan-400" />
          </div>
          <h2 class="font-bold text-base" style="color: var(--text);">Información de cuenta</h2>
        </div>
        <div class="space-y-3">
          <div class="flex items-center justify-between py-2.5 px-3 rounded-xl" style="background: var(--input-bg);">
            <span class="text-sm" style="color: var(--text-muted);">Nombre completo</span>
            <span class="text-sm font-semibold" style="color: var(--text);">{{ fullName }}</span>
          </div>
          <div class="flex items-center justify-between py-2.5 px-3 rounded-xl" style="background: var(--input-bg);">
            <span class="text-sm" style="color: var(--text-muted);">Proyecto</span>
            <span class="text-sm font-semibold truncate max-w-[180px]" style="color: var(--text);">{{ auth.user?.project_name }}</span>
          </div>
          <div class="flex items-center justify-between py-2.5 px-3 rounded-xl" style="background: var(--input-bg);">
            <span class="text-sm" style="color: var(--text-muted);">Miembro desde</span>
            <span class="text-sm font-semibold" style="color: var(--text);">{{ memberSince }}</span>
          </div>
          <div class="flex items-center justify-between py-2.5 px-3 rounded-xl" style="background: var(--input-bg);">
            <span class="text-sm" style="color: var(--text-muted);">Rol</span>
            <span class="badge" :class="auth.isAdmin ? 'badge-blue' : 'badge-gray'">
              {{ auth.isAdmin ? 'Administrador' : 'Usuario' }}
            </span>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import ThemeToggle from '@/components/ThemeToggle.vue'
import {
  ArrowLeftIcon, CameraIcon, EnvelopeIcon, LockClosedIcon,
  CheckIcon, XMarkIcon, CheckCircleIcon, ExclamationTriangleIcon,
  EyeIcon, EyeSlashIcon, InformationCircleIcon
} from '@heroicons/vue/24/outline'

const auth = useAuthStore()

// ── Avatar ────────────────────────────────────────────────
const fileInput = ref(null)
const previewUrl = ref('')
const pendingFile = ref(null)
const savingAvatar = ref(false)
const avatarSuccess = ref(false)
const avatarError = ref('')

const initials = computed(() => {
  const u = auth.user
  if (!u) return '?'
  return `${u.first_name?.[0] ?? ''}${u.last_name?.[0] ?? ''}`.toUpperCase()
})

function triggerAvatarPick() { fileInput.value?.click() }

function onFileChange(e) {
  const file = e.target.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) {
    avatarError.value = 'Solo se aceptan imágenes (jpg, png, webp…)'
    return
  }
  if (file.size > 500_000) {
    avatarError.value = 'La imagen supera el límite de 500 KB'
    return
  }
  avatarError.value = ''
  pendingFile.value = file
  const reader = new FileReader()
  reader.onload = ev => { previewUrl.value = ev.target.result }
  reader.readAsDataURL(file)
}

function cancelAvatar() {
  previewUrl.value = ''
  pendingFile.value = null
  if (fileInput.value) fileInput.value.value = ''
}

async function saveAvatar() {
  if (!previewUrl.value) return
  savingAvatar.value = true
  avatarError.value = ''
  try {
    await api.put('/profile/avatar', { avatar_url: previewUrl.value })
    auth.user.avatar_url = previewUrl.value
    auth.persistUser()
    avatarSuccess.value = true
    previewUrl.value = ''
    pendingFile.value = null
    setTimeout(() => { avatarSuccess.value = false }, 3000)
  } catch (e) {
    avatarError.value = e.response?.data?.error ?? 'Error al guardar la imagen'
  } finally {
    savingAvatar.value = false
  }
}

// ── Email ─────────────────────────────────────────────────
const emailForm = ref({ newEmail: '' })
const savingEmail = ref(false)
const emailMsg = ref({ text: '', ok: false })

async function handleEmailUpdate() {
  savingEmail.value = true
  emailMsg.value = { text: '', ok: false }
  try {
    await api.put('/profile/email', { email: emailForm.value.newEmail })
    auth.user.email = emailForm.value.newEmail
    auth.persistUser()
    emailForm.value.newEmail = ''
    emailMsg.value = { text: 'Correo actualizado correctamente', ok: true }
    setTimeout(() => { emailMsg.value.text = '' }, 4000)
  } catch (e) {
    emailMsg.value = { text: e.response?.data?.error ?? 'Error al actualizar el correo', ok: false }
  } finally {
    savingEmail.value = false
  }
}

// ── Password ──────────────────────────────────────────────
const pwdForm = ref({ current: '', newPwd: '', confirm: '' })
const savingPwd = ref(false)
const pwdMsg = ref({ text: '', ok: false })
const showCurrent = ref(false)
const showNew = ref(false)

const pwdMismatch = computed(() =>
  pwdForm.value.confirm.length > 0 && pwdForm.value.newPwd !== pwdForm.value.confirm
)

const pwdStrength = computed(() => {
  const p = pwdForm.value.newPwd
  let score = 0
  if (p.length >= 8) score++
  if (/[A-Z]/.test(p)) score++
  if (/[0-9]/.test(p)) score++
  if (/[^A-Za-z0-9]/.test(p)) score++
  const levels = [
    { level: 1, color: '#ef4444', label: 'Débil' },
    { level: 2, color: '#f97316', label: 'Regular' },
    { level: 3, color: '#eab308', label: 'Buena' },
    { level: 4, color: '#22c55e', label: 'Fuerte' },
  ]
  return levels[Math.max(score - 1, 0)]
})

async function handlePasswordUpdate() {
  if (pwdMismatch.value) return
  savingPwd.value = true
  pwdMsg.value = { text: '', ok: false }
  try {
    await api.put('/profile/password', {
      current_password: pwdForm.value.current,
      new_password: pwdForm.value.newPwd,
    })
    pwdForm.value = { current: '', newPwd: '', confirm: '' }
    pwdMsg.value = { text: 'Contraseña cambiada correctamente', ok: true }
    setTimeout(() => { pwdMsg.value.text = '' }, 4000)
  } catch (e) {
    pwdMsg.value = { text: e.response?.data?.error ?? 'Error al cambiar la contraseña', ok: false }
  } finally {
    savingPwd.value = false
  }
}

// ── Helpers ───────────────────────────────────────────────
const fullName = computed(() => {
  const u = auth.user
  return u ? `${u.first_name} ${u.last_name}` : ''
})

const memberSince = computed(() => {
  const d = auth.user?.created_at
  if (!d) return '—'
  return new Intl.DateTimeFormat('es-MX', { year: 'numeric', month: 'long', day: 'numeric' }).format(new Date(d))
})
</script>
