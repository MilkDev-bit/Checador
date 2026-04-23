<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow-sm px-4 py-4 flex items-center gap-3">
      <router-link to="/checkin" class="text-blue-500">← Volver</router-link>
      <h1 class="font-bold text-gray-800 text-lg">Código QR de Acceso</h1>
    </header>

    <div class="max-w-sm mx-auto px-4 pt-10 text-center space-y-6">
      <div class="card">
        <div class="text-4xl mb-3">📲</div>
        <h2 class="text-xl font-bold mb-2">Escanea para checar</h2>
        <p class="text-gray-500 text-sm mb-6">Comparte este QR para que otros accedan al sistema de checado</p>

        <div v-if="loading" class="py-12 text-gray-400">Generando QR...</div>
        <img v-else :src="qrUrl" alt="QR Code" class="w-64 h-64 mx-auto rounded-xl shadow" />

        <p class="text-xs text-gray-400 mt-4 break-all">{{ checkinUrl }}</p>
      </div>

      <button @click="downloadQR" class="btn-primary w-full" :disabled="loading">
        ⬇️ Descargar QR
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
