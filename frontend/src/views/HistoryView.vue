<template>
  <div class="min-h-screen min-h-dvh bg-mesh pb-10" style="background-color: #0f1629;">
    <!-- Header -->
    <header class="sticky top-0 z-10 px-4 py-3 flex items-center gap-3"
      style="background: rgba(15,22,41,0.9); backdrop-filter: blur(12px); border-bottom: 1px solid rgba(255,255,255,0.06);">
      <router-link to="/checkin" class="w-8 h-8 rounded-xl flex items-center justify-center text-slate-400 hover:text-white hover:bg-white/8 transition-all">
        ←
      </router-link>
      <h1 class="font-bold text-white text-base">Historial de Registros</h1>
    </header>

    <div class="max-w-lg mx-auto px-4 pt-5 space-y-3">
      <!-- Loading -->
      <div v-if="loading" class="flex flex-col items-center py-16 gap-3">
        <div class="w-8 h-8 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin"></div>
        <p class="text-slate-500 text-sm">Cargando registros...</p>
      </div>

      <!-- Empty -->
      <div v-else-if="records.length === 0" class="glass-card p-10 text-center animate-in"
        style="background: rgba(20,29,53,0.8);">
        <div class="text-5xl mb-3 opacity-50">📋</div>
        <p class="text-slate-400 font-medium">Sin registros aún</p>
        <p class="text-slate-600 text-sm mt-1">Tus checadas aparecerán aquí</p>
      </div>

      <!-- Records -->
      <div v-for="(record, i) in records" :key="record.id"
        class="glass-card overflow-hidden animate-in"
        :style="`background: rgba(20,29,53,0.8); animation-delay: ${i * 0.04}s`">

        <!-- Record header -->
        <div class="px-5 py-4 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center text-lg"
              :style="record.type === 'entry'
                ? 'background: rgba(16,185,129,0.15); border: 1px solid rgba(16,185,129,0.3);'
                : 'background: rgba(244,63,94,0.15); border: 1px solid rgba(244,63,94,0.3);'">
              {{ record.type === 'entry' ? '📍' : '🚪' }}
            </div>
            <div>
              <span :class="record.type === 'entry' ? 'badge-green' : 'badge-red'" class="badge mb-1">
                {{ record.type === 'entry' ? 'Entrada' : 'Salida' }}
              </span>
              <p class="text-slate-400 text-xs">{{ formatDate(record.timestamp) }}</p>
            </div>
          </div>
          <button @click="viewRoute(record.id)"
            class="btn-secondary text-xs px-3 py-1.5 rounded-lg">
            Ver ruta →
          </button>
        </div>

        <!-- Photos -->
        <div v-if="record.photo_site_path || record.photo_selfie_path"
          class="px-5 pb-4 grid grid-cols-2 gap-3">
          <div v-if="record.photo_site_path">
            <p class="text-slate-500 text-xs mb-1.5 flex items-center gap-1"><span>🏗️</span> Sitio</p>
            <img :src="record.photo_site_path" class="w-full h-28 object-cover rounded-xl" />
          </div>
          <div v-if="record.photo_selfie_path">
            <p class="text-slate-500 text-xs mb-1.5 flex items-center gap-1"><span>🤳</span> Selfie</p>
            <img :src="record.photo_selfie_path" class="w-full h-28 object-cover rounded-xl" />
          </div>
        </div>
      </div>
    </div>

    <!-- Route modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showRouteModal" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center px-4 pb-4 sm:pb-0"
          style="background: rgba(0,0,0,0.75); backdrop-filter: blur(6px);">
          <div class="w-full max-w-md glass-card max-h-[75vh] flex flex-col animate-in"
            style="background: rgba(20,29,53,0.98);">
            <div class="flex items-center justify-between px-5 py-4 border-b border-white/8 flex-shrink-0">
              <div>
                <h3 class="font-bold text-white">Puntos de Recorrido</h3>
                <p class="text-slate-500 text-xs mt-0.5">{{ routePoints.length }} puntos registrados</p>
              </div>
              <button @click="showRouteModal = false"
                class="w-8 h-8 rounded-lg flex items-center justify-center text-slate-400 hover:text-white hover:bg-white/8 transition-all">
                ✕
              </button>
            </div>

            <div class="overflow-y-auto flex-1 custom-scroll">
              <div v-if="routePoints.length === 0" class="text-center text-slate-500 py-8 text-sm">
                Sin puntos registrados
              </div>
              <div v-for="(point, i) in routePoints" :key="point.id"
                class="flex items-start gap-3 px-5 py-3 border-b border-white/5 last:border-0">
                <div class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0 mt-0.5"
                  style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">
                  {{ i + 1 }}
                </div>
                <div>
                  <p class="text-white text-sm font-mono">{{ point.latitude.toFixed(6) }}, {{ point.longitude.toFixed(6) }}</p>
                  <p class="text-slate-500 text-xs mt-0.5">{{ formatDate(point.recorded_at) }} · ±{{ Math.round(point.accuracy) }}m</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>
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

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.25s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>

