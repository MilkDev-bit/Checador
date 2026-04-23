<template>
  <div class="min-h-screen min-h-dvh bg-mesh" style="background-color: #0f1629;">
    <header class="sticky top-0 z-10 px-4 py-3 flex items-center gap-3"
      style="background: rgba(15,22,41,0.9); backdrop-filter: blur(12px); border-bottom: 1px solid rgba(255,255,255,0.06);">
      <router-link to="/checkin" class="w-8 h-8 rounded-xl flex items-center justify-center text-slate-400 hover:text-white hover:bg-white/8 transition-all">
        ←
      </router-link>
      <h1 class="font-bold text-white text-base">Código QR de Acceso</h1>
    </header>

    <div class="max-w-sm mx-auto px-4 pt-8 pb-10 text-center space-y-5">
      <!-- Card -->
      <div class="glass-card p-7 animate-in" style="background: rgba(20,29,53,0.85);">
        <div class="w-16 h-16 rounded-2xl flex items-center justify-center text-3xl mx-auto mb-4"
          style="background: linear-gradient(135deg, rgba(99,102,241,0.3), rgba(6,182,212,0.3)); border: 1px solid rgba(99,102,241,0.3);">
          📲
        </div>
        <h2 class="text-xl font-bold text-white mb-1">Escanea para checar</h2>
        <p class="text-slate-400 text-sm mb-6 leading-relaxed">
          Comparte este código QR para acceder directamente al sistema de checado
        </p>

        <!-- QR image -->
        <div v-if="loading" class="w-64 h-64 mx-auto flex items-center justify-center rounded-2xl"
          style="background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.08);">
          <div class="w-8 h-8 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin"></div>
        </div>
        <div v-else class="p-3 rounded-2xl inline-block" style="background: white;">
          <img :src="qrUrl" alt="QR Code" class="w-56 h-56 block" />
        </div>

        <p class="text-slate-600 text-xs mt-4 break-all font-mono">{{ checkinUrl }}</p>
      </div>

      <button @click="downloadQR" :disabled="loading" class="btn-primary btn-lg w-full">
        <span>⬇️</span> Descargar QR
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

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

