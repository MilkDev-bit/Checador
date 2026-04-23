<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <header class="bg-white shadow-sm px-4 py-4 flex items-center gap-3">
      <router-link to="/checkin" class="text-blue-500">← Volver</router-link>
      <h1 class="font-bold text-gray-800 text-lg">Historial de Registros</h1>
    </header>

    <div class="max-w-lg mx-auto px-4 pt-6 space-y-4">
      <div v-if="loading" class="text-center text-gray-500 py-8">Cargando...</div>

      <div v-else-if="records.length === 0" class="card text-center text-gray-500 py-8">
        <div class="text-4xl mb-2">📋</div>
        <p>No tienes registros aún</p>
      </div>

      <div v-for="record in records" :key="record.id" class="card">
        <div class="flex items-center justify-between mb-3">
          <span :class="record.type === 'entry' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
            class="px-3 py-1 rounded-full text-sm font-semibold">
            {{ record.type === 'entry' ? '📍 Entrada' : '🚪 Salida' }}
          </span>
          <span class="text-sm text-gray-500">{{ formatDate(record.timestamp) }}</span>
        </div>

        <div class="flex gap-3 mt-3">
          <div v-if="record.photo_site_path" class="flex-1">
            <p class="text-xs text-gray-500 mb-1">Sitio</p>
            <img :src="record.photo_site_path" class="w-full h-32 object-cover rounded-lg" />
          </div>
          <div v-if="record.photo_selfie_path" class="flex-1">
            <p class="text-xs text-gray-500 mb-1">Selfie</p>
            <img :src="record.photo_selfie_path" class="w-full h-32 object-cover rounded-lg" />
          </div>
        </div>

        <button @click="viewRoute(record.id)" class="mt-3 text-sm text-blue-600 hover:underline">
          Ver recorrido →
        </button>
      </div>
    </div>

    <!-- Route modal -->
    <Teleport to="body">
      <div v-if="showRouteModal" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50 px-4">
        <div class="card max-w-lg w-full max-h-[80vh] overflow-y-auto">
          <div class="flex justify-between items-center mb-4">
            <h3 class="font-bold text-lg">Puntos de Recorrido</h3>
            <button @click="showRouteModal = false" class="text-gray-400 hover:text-gray-700 text-2xl">✕</button>
          </div>
          <div v-if="routePoints.length === 0" class="text-center text-gray-500 py-4">Sin puntos registrados</div>
          <div v-for="(point, i) in routePoints" :key="point.id" class="flex items-start gap-3 py-2 border-b last:border-0">
            <div class="w-6 h-6 rounded-full bg-blue-600 text-white text-xs flex items-center justify-center flex-shrink-0 mt-0.5">{{ i + 1 }}</div>
            <div>
              <p class="text-sm font-medium">{{ point.latitude.toFixed(6) }}, {{ point.longitude.toFixed(6) }}</p>
              <p class="text-xs text-gray-500">{{ formatDate(point.recorded_at) }} · Precisión: {{ Math.round(point.accuracy) }}m</p>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const records = ref([])
const loading = ref(true)
const showRouteModal = ref(false)
const routePoints = ref([])

onMounted(async () => {
  try {
    const { data } = await api.get('/checks')
    records.value = data
  } finally {
    loading.value = false
  }
})

async function viewRoute(id) {
  const { data } = await api.get(`/checks/${id}/route`)
  routePoints.value = data
  showRouteModal.value = true
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('es-MX', {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}
</script>
