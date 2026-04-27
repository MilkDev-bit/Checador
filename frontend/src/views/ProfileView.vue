<template>
  <div class="min-h-screen min-h-dvh pb-16" style="background-color: var(--bg);">

    <!-- Header -->
    <header class="sticky top-0 z-20 px-4 py-3 safe-top flex items-center gap-3"
      style="background: var(--nav-bg); backdrop-filter: blur(14px); border-bottom: 1px solid var(--border-subtle); --header-base-padding: 0.75rem;">
      <router-link to="/checkin"
        class="w-9 h-9 rounded-xl flex items-center justify-center transition-all"
        style="color: var(--text-muted);"
        @mouseenter="e => e.currentTarget.style.background='var(--input-bg)'"
        @mouseleave="e => e.currentTarget.style.background='transparent'">
        <ArrowLeftIcon class="w-5 h-5" />
      </router-link>
      <h1 class="font-bold text-base" style="color: var(--text);">Mi Perfil</h1>
    </header>

    <div class="max-w-lg mx-auto px-4 pt-6 space-y-5">

      <!-- ── Avatar + Cover card ───────────────────────── -->
      <div class="glass-card overflow-hidden animate-in">

        <!-- Cover / Banner -->
        <div class="relative h-28 group cursor-pointer" @click="triggerCoverPick">
          <!-- Cover image or gradient fallback -->
          <div class="absolute inset-0 overflow-hidden">
            <img v-if="coverPreview || auth.user?.cover_url"
              :src="coverPreview || auth.user?.cover_url"
              class="w-full h-full object-cover" alt="Portada" />
            <div v-else class="w-full h-full"
              style="background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #06b6d4 100%);">
              <div class="absolute inset-0 opacity-20"
                style="background-image: radial-gradient(#fff 1px, transparent 1px); background-size: 18px 18px;"></div>
            </div>
          </div>
          <!-- Hover overlay for cover -->
          <div class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-all duration-200 flex items-center justify-center">
            <div class="opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-2 bg-black/60 text-white text-xs font-semibold px-3 py-1.5 rounded-full backdrop-blur-sm">
              <CameraIcon class="w-3.5 h-3.5" /> Cambiar portada
            </div>
          </div>
          <input ref="coverInput" type="file" accept="image/*" class="hidden" @change="onCoverChange" />
        </div>

        <!-- Avatar overlapping banner -->
        <div class="px-5 pb-5">
          <div class="flex items-end gap-4 -mt-9 mb-4">
            <!-- Avatar wrapper — needs position:relative and group here -->
            <div class="relative flex-shrink-0 group w-[72px] h-[72px] cursor-pointer"
              @click="triggerAvatarPick"
              style="filter: drop-shadow(0 4px 12px rgba(0,0,0,0.3));">
              <div class="w-[72px] h-[72px] rounded-2xl overflow-hidden ring-4 ring-[var(--card-bg)]"
                style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">
                <img v-if="previewUrl || auth.user?.avatar_url"
                  :src="previewUrl || auth.user?.avatar_url"
                  class="w-full h-full object-cover" alt="Avatar" />
                <div v-else class="w-full h-full flex items-center justify-center text-white text-xl font-black select-none">
                  {{ initials }}
                </div>
              </div>
              <!-- Hover overlay over avatar -->
              <div class="absolute inset-0 rounded-2xl bg-black/0 group-hover:bg-black/50 transition-all duration-200 flex items-center justify-center">
                <CameraIcon class="w-5 h-5 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
              </div>
              <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onFileChange" />
            </div>

            <div class="mb-1 min-w-0 flex-1 pt-10">
              <p class="font-black text-lg truncate leading-tight" style="color: var(--text);">{{ fullName }}</p>
              <p class="text-sm truncate" style="color: var(--text-muted);">{{ auth.user?.email }}</p>
              <span class="badge mt-1.5" :class="auth.isAdmin ? 'badge-blue' : 'badge-gray'">
                {{ auth.isAdmin ? 'Administrador' : 'Usuario' }}
              </span>
            </div>
          </div>

          <!-- Cover save/cancel -->
          <Transition name="fade">
            <div v-if="coverPreview" class="flex gap-2 mb-3 p-3 rounded-xl" style="background: var(--input-bg);">
              <div class="flex items-center gap-2 flex-1 text-sm" style="color: var(--text-muted);">
                <PhotoIcon class="w-4 h-4 flex-shrink-0" /> Nueva portada lista
              </div>
              <button @click="saveCover" :disabled="savingCover" class="btn btn-primary btn-sm">
                <span v-if="savingCover" class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                <CheckIcon v-else class="w-3.5 h-3.5" />
                Guardar
              </button>
              <button @click="cancelCover" class="btn btn-secondary btn-sm px-3">
                <XMarkIcon class="w-3.5 h-3.5" />
              </button>
            </div>
          </Transition>

          <!-- Avatar save/cancel -->
          <Transition name="fade">
            <div v-if="previewUrl" class="flex gap-2 mb-3 p-3 rounded-xl" style="background: var(--input-bg);">
              <div class="flex items-center gap-2 flex-1 text-sm" style="color: var(--text-muted);">
                <UserCircleIcon class="w-4 h-4 flex-shrink-0" /> Nueva foto de perfil lista
              </div>
              <button @click="saveAvatar" :disabled="savingAvatar" class="btn btn-primary btn-sm">
                <span v-if="savingAvatar" class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                <CheckIcon v-else class="w-3.5 h-3.5" />
                Guardar
              </button>
              <button @click="cancelAvatar" class="btn btn-secondary btn-sm px-3">
                <XMarkIcon class="w-3.5 h-3.5" />
              </button>
            </div>
          </Transition>

          <Transition name="fade">
            <div v-if="avatarSuccess" class="flex items-center gap-2 bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 rounded-xl px-4 py-2.5 text-sm font-semibold">
              <CheckCircleIcon class="w-4 h-4 flex-shrink-0" /> Foto de perfil actualizada
            </div>
          </Transition>
          <Transition name="fade">
            <div v-if="avatarError" class="flex items-center gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-500 rounded-xl px-4 py-2.5 text-sm font-semibold">
              <ExclamationTriangleIcon class="w-4 h-4 flex-shrink-0" /> {{ avatarError }}
            </div>
          </Transition>
        </div>
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
      <!-- ── Appearance card ────────────────────────── -->
      <div class="glass-card p-6 animate-in" style="animation-delay: 0.15s">
        <div class="flex items-center gap-3 mb-5">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
            style="background: rgba(245,158,11,0.12); border: 1px solid rgba(245,158,11,0.2);">
            <SunIcon class="w-5 h-5 text-amber-400" />
          </div>
          <h2 class="font-bold text-base" style="color: var(--text);">Apariencia</h2>
        </div>
        <div class="flex items-center justify-between py-3 px-4 rounded-xl" style="background: var(--input-bg);">
          <div>
            <p class="text-sm font-semibold" style="color: var(--text);">Modo Oscuro</p>
            <p class="text-xs mt-0.5" style="color: var(--text-muted);">Cambiar tema visual</p>
          </div>
          <button
            @click="theme.toggle()"
            class="relative inline-flex h-7 w-12 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none"
            :class="theme.isDark ? 'bg-brand-500' : 'bg-slate-300 dark:bg-surface-600'"
            role="switch"
            :aria-checked="theme.isDark"
          >
            <span
              class="pointer-events-none inline-block h-6 w-6 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
              :class="theme.isDark ? 'translate-x-5' : 'translate-x-0'"
            ></span>
          </button>
        </div>
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
import { useThemeStore } from '@/stores/theme'
import {
  ArrowLeftIcon, CameraIcon, EnvelopeIcon, LockClosedIcon,
  CheckIcon, XMarkIcon, CheckCircleIcon, ExclamationTriangleIcon,
  EyeIcon, EyeSlashIcon, InformationCircleIcon, UserCircleIcon, PhotoIcon, SunIcon
} from '@heroicons/vue/24/outline'

const auth = useAuthStore()
const theme = useThemeStore()

// ── Cover / Banner ────────────────────────────────────────
const coverInput = ref(null)
const coverPreview = ref('')
const savingCover = ref(false)

function triggerCoverPick() { coverInput.value?.click() }

function onCoverChange(e) {
  const file = e.target.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) return
  if (file.size > 700_000) {
    avatarError.value = 'La portada supera el límite de 700 KB'
    return
  }
  const reader = new FileReader()
  reader.onload = ev => { coverPreview.value = ev.target.result }
  reader.readAsDataURL(file)
}

function cancelCover() {
  coverPreview.value = ''
  if (coverInput.value) coverInput.value.value = ''
}

async function saveCover() {
  if (!coverPreview.value) return
  savingCover.value = true
  try {
    await api.put('/profile/cover', { avatar_url: coverPreview.value })
    auth.user.cover_url = coverPreview.value
    auth.persistUser()
    coverPreview.value = ''
  } catch (e) {
    avatarError.value = e.response?.data?.error ?? 'Error al guardar la portada'
  } finally {
    savingCover.value = false
  }
}

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
