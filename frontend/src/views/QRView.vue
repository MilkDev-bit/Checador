<template>
  <div class="min-h-screen min-h-dvh bg-mesh" style="background-color: var(--bg);">
    <header class="sticky top-0 z-10 px-4 py-3 safe-top flex items-center gap-3"
      style="background: var(--nav-bg); backdrop-filter: blur(12px); border-bottom: 1px solid var(--border-subtle); --header-base-padding: 0.75rem;">
      <router-link to="/checkin" class="w-8 h-8 rounded-xl flex items-center justify-center transition-all"
        style="color: var(--text-muted);">
        <ArrowLeftIcon class="w-5 h-5" />
      </router-link>
      <h1 class="font-bold text-base flex-1" style="color: var(--text);">Código QR de Acceso</h1>

      <router-link to="/profile"
        class="w-8 h-8 rounded-xl overflow-hidden flex-shrink-0 transition-all hover:ring-2 ring-brand-400/60"
        style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">
        <img v-if="auth?.user?.avatar_url" :src="auth.user.avatar_url" class="w-full h-full object-cover" alt="Perfil" />
        <span v-else class="w-full h-full flex items-center justify-center text-white text-xs font-black">
          {{ auth?.user?.first_name?.[0] ?? '' }}{{ auth?.user?.last_name?.[0] ?? '' }}
        </span>
      </router-link>
    </header>

    <div class="max-w-sm mx-auto px-4 pt-8 pb-10 text-center space-y-5">
      <!-- Card -->
      <div class="glass-card p-7 animate-in">
        <div class="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4"
          style="background: linear-gradient(135deg, rgba(99,102,241,0.3), rgba(6,182,212,0.3)); border: 1px solid rgba(99,102,241,0.3);">
          <QrCodeIcon class="w-8 h-8 text-brand-400" />
        </div>
        <h2 class="text-xl font-bold mb-1" style="color: var(--text);">Escanea para checar</h2>
        <p class="text-sm mb-6 leading-relaxed" style="color: var(--text-muted);">
          Comparte este código QR para acceder directamente al sistema de checado
        </p>

        <!-- QR image -->
        <div v-if="loading" class="w-64 h-64 mx-auto flex items-center justify-center rounded-2xl"
          style="background: var(--input-bg); border: 1px solid var(--input-border);">
          <div class="w-8 h-8 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin"></div>
        </div>
        <div v-else class="p-3 rounded-2xl inline-block" style="background: white;">
          <img :src="qrUrl" alt="QR Code" class="w-56 h-56 block" />
        </div>

        <p class="text-xs mt-4 break-all font-mono" style="color: var(--text-dim);">{{ checkinUrl }}</p>
      </div>

      <button @click="downloadQR" :disabled="loading" class="btn-primary btn-lg w-full">
        <ArrowDownTrayIcon class="w-5 h-5" /> Descargar QR
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { useAuthStore } from '@/stores/auth'
import { ArrowLeftIcon, QrCodeIcon, ArrowDownTrayIcon } from '@heroicons/vue/24/outline'

const auth = useAuthStore()

const qrUrl = ref('')
const checkinUrl = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    const urlRes = await api.get('/checkin-url')
    checkinUrl.value = urlRes.data.url
    const res = await fetch('/api/qr')
    const blob = await res.blob()
    qrUrl.value = URL.createObjectURL(blob)
  } finally {
    loading.value = false
  }
})

function downloadQR() {
  const a = document.createElement('a')
  a.href = qrUrl.value
  a.download = 'paselista-qr.png'
  a.click()
}
</script>

