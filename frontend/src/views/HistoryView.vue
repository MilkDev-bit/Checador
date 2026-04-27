<template>
  <div class="min-h-screen min-h-dvh bg-mesh pb-10" style="background-color: var(--bg);">
    <!-- Header -->
    <header class="sticky top-0 z-10 px-4 py-3 flex items-center gap-3"
      style="background: var(--nav-bg); backdrop-filter: blur(12px); border-bottom: 1px solid var(--border-subtle);">
      <router-link to="/checkin" class="w-8 h-8 rounded-xl flex items-center justify-center transition-all"
        style="color: var(--text-muted);"
        @mouseenter="e => e.currentTarget.style.background='var(--input-bg)'"
        @mouseleave="e => e.currentTarget.style.background='transparent'">
        <ArrowLeftIcon class="w-5 h-5" />
      </router-link>
      <h1 class="font-bold text-base flex-1" style="color: var(--text);">Historial de Registros</h1>
      <ThemeToggle />
      <router-link to="/profile"
        class="w-8 h-8 rounded-xl overflow-hidden flex-shrink-0 transition-all hover:ring-2 ring-brand-400/60"
        style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">
        <img v-if="auth?.user?.avatar_url" :src="auth.user.avatar_url" class="w-full h-full object-cover" alt="Perfil" />
        <span v-else class="w-full h-full flex items-center justify-center text-white text-xs font-black">
          {{ auth?.user?.first_name?.[0] ?? '' }}{{ auth?.user?.last_name?.[0] ?? '' }}
        </span>
      </router-link>
    </header>

    <div class="max-w-lg mx-auto px-4 pt-5 space-y-3">
      <!-- Loading -->
      <div v-if="loading" class="flex flex-col items-center py-16 gap-3">
        <div class="w-8 h-8 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin"></div>
        <p class="text-sm" style="color: var(--text-muted);">Cargando registros...</p>
      </div>

      <!-- Empty -->
      <div v-else-if="records.length === 0" class="glass-card p-10 text-center animate-in">
        <div class="flex justify-center mb-3 opacity-50">
          <ClipboardDocumentListIcon class="w-12 h-12" style="color: var(--text-muted);" />
        </div>
        <p class="font-medium" style="color: var(--text-muted);">Sin registros aún</p>
        <p class="text-sm mt-1" style="color: var(--text-dim);">Tus checadas aparecerán aquí</p>
      </div>

      <!-- Records -->
      <div v-for="(record, i) in records" :key="record.id"
        class="glass-card overflow-hidden animate-in"
        :style="`animation-delay: ${i * 0.04}s`">

        <!-- Record header -->
        <div class="px-5 py-4 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center"
              :style="record.type === 'entry'
                ? 'background: rgba(16,185,129,0.15); border: 1px solid rgba(16,185,129,0.3);'
                : 'background: rgba(244,63,94,0.15); border: 1px solid rgba(244,63,94,0.3);'">
              <MapPinIcon v-if="record.type === 'entry'" class="w-5 h-5 text-emerald-400" />
              <ArrowRightOnRectangleIcon v-else class="w-5 h-5 text-rose-400" />
            </div>
            <div>
              <span :class="record.type === 'entry' ? 'badge-green' : 'badge-red'" class="badge mb-1">
                {{ record.type === 'entry' ? 'Entrada' : 'Salida' }}
              </span>
              <p class="text-xs" style="color: var(--text-muted);">{{ formatDate(record.timestamp) }}</p>
            </div>
          </div>
          <button @click="viewRoute(record.id)" class="btn-secondary text-xs px-3 py-1.5 rounded-lg">
            Ver ruta →
          </button>
        </div>

        <!-- Photos -->
        <div v-if="record.photo_site_path || record.photo_selfie_path"
          class="px-5 pb-4 grid grid-cols-2 gap-3">
          <div v-if="record.photo_site_path">
            <p class="text-xs mb-1.5 flex items-center gap-1" style="color: var(--text-muted);">
              <BuildingOffice2Icon class="w-3 h-3" /> Sitio
            </p>
            <img :src="record.photo_site_path" class="w-full h-28 object-cover rounded-xl" />
          </div>
          <div v-if="record.photo_selfie_path">
            <p class="text-xs mb-1.5 flex items-center gap-1" style="color: var(--text-muted);">
              <UserCircleIcon class="w-3 h-3" /> Selfie
            </p>
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
            style="background: var(--modal-bg);">
            <div class="flex items-center justify-between px-5 py-4 flex-shrink-0"
              style="border-bottom: 1px solid var(--border-subtle);">
              <div>
                <h3 class="font-bold" style="color: var(--text);">Puntos de Recorrido</h3>
                <p class="text-xs mt-0.5" style="color: var(--text-muted);">{{ routePoints.length }} puntos registrados</p>
              </div>
              <button @click="showRouteModal = false"
                class="w-8 h-8 rounded-lg flex items-center justify-center transition-all"
                style="color: var(--text-muted);">
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>

            <div class="overflow-y-auto flex-1 custom-scroll">
              <div v-if="routePoints.length === 0" class="text-center py-8 text-sm" style="color: var(--text-muted);">
                Sin puntos registrados
              </div>
              <div v-for="(point, i) in routePoints" :key="point.id"
                class="flex items-start gap-3 px-5 py-3 last:border-0"
                style="border-bottom: 1px solid var(--border-subtle);">
                <div class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0 mt-0.5 text-white"
                  style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">
                  {{ i + 1 }}
                </div>
                <div>
                  <p class="text-sm font-mono" style="color: var(--text);">{{ point.latitude.toFixed(6) }}, {{ point.longitude.toFixed(6) }}</p>
                  <p class="text-xs mt-0.5" style="color: var(--text-muted);">{{ formatDate(point.recorded_at) }} · ±{{ Math.round(point.accuracy) }}m</p>
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
import { useAuthStore } from '@/stores/auth'
import {
  ArrowLeftIcon, MapPinIcon, ArrowRightOnRectangleIcon,
  ClipboardDocumentListIcon, BuildingOffice2Icon, UserCircleIcon, XMarkIcon
} from '@heroicons/vue/24/outline'
import ThemeToggle from '@/components/ThemeToggle.vue'

const auth = useAuthStore()
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
