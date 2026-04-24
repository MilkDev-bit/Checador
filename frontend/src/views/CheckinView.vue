<template>
  <div class="min-h-screen min-h-dvh bg-mesh" style="background-color: var(--bg);">
    <!-- Top nav -->
    <header class="sticky top-0 z-20 px-4 py-3 flex items-center justify-between"
      style="background: var(--nav-bg); backdrop-filter: blur(12px); border-bottom: 1px solid var(--border-subtle);">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center"
          style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">
          <CheckCircleIcon class="w-4 h-4 text-white" />
        </div>
        <div>
          <p class="font-semibold text-sm leading-none" style="color: var(--text);">PaseLista</p>
          <p class="text-xs mt-0.5 truncate max-w-[160px]" style="color: var(--text-muted);">{{ auth.user?.project_name }}</p>
        </div>
      </div>
      <div class="flex items-center gap-1">
        <router-link to="/qr" class="btn-ghost text-xs px-3 py-2">QR</router-link>
        <router-link to="/history" class="btn-ghost text-xs px-3 py-2">Historial</router-link>
        <ThemeToggle />
        <button @click="auth.logout(); $router.push('/login')" class="btn-ghost text-xs px-3 py-2 text-rose-400">Salir</button>
      </div>
    </header>

    <!-- GPS background warning banner -->
    <Transition name="slide-banner">
      <div v-if="showHiddenWarning"
        class="px-4 py-3 flex items-center gap-3 text-sm"
        style="background: rgba(245,158,11,0.15); border-bottom: 1px solid rgba(245,158,11,0.3);">
        <ExclamationTriangleIcon class="w-5 h-5 text-amber-400 flex-shrink-0" />
        <div class="flex-1">
          <p class="text-amber-400 font-semibold text-xs">GPS posiblemente pausado</p>
          <p class="text-amber-300/70 text-xs">La app estuvo en segundo plano. Algunos puntos de recorrido podrían no haberse registrado.</p>
        </div>
        <button @click="showHiddenWarning = false" class="text-amber-400/60 hover:text-amber-300 flex-shrink-0">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
    </Transition>

    <!-- Active tracking indicator -->
    <div v-if="isTracking" class="px-4 py-2 flex items-center gap-2 text-xs"
      style="background: rgba(16,185,129,0.08); border-bottom: 1px solid rgba(16,185,129,0.15);">
      <div class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></div>
      <span class="text-emerald-400/80">Rastreo GPS activo · {{ locationPoints.length }} puntos</span>
      <span v-if="wakeLockActive" class="flex items-center gap-1 text-emerald-400/80">
        · <LockClosedIcon class="w-3 h-3" /> Pantalla activa
      </span>
    </div>

    <div class="max-w-lg mx-auto px-4 pt-6 pb-10 space-y-5">
      <!-- Greeting card -->
      <div class="glass-card p-5 animate-in">
        <div class="flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center flex-shrink-0"
            style="background: linear-gradient(135deg, rgba(99,102,241,0.3), rgba(139,92,246,0.3)); border: 1px solid rgba(99,102,241,0.3);">
            <HandRaisedIcon class="w-7 h-7 text-brand-400" />
          </div>
          <div class="min-w-0">
            <h2 class="text-lg font-bold truncate" style="color: var(--text);">Hola, {{ auth.user?.first_name }}</h2>
            <p class="text-sm" style="color: var(--text-muted);">{{ todayStr }}</p>
            <p class="text-brand-400 font-mono font-semibold text-base mt-0.5">{{ currentTime }}</p>
          </div>
        </div>
      </div>

      <!-- Recovered points notice -->
      <div v-if="recoveredPoints > 0" class="glass-card p-4 animate-in"
        style="background: rgba(99,102,241,0.1); border: 1px solid rgba(99,102,241,0.25);">
        <div class="flex items-start gap-3">
          <CloudArrowUpIcon class="w-5 h-5 text-brand-400 flex-shrink-0 mt-0.5" />
          <div>
            <p class="text-brand-300 font-semibold text-sm">Puntos GPS recuperados</p>
            <p class="text-sm mt-0.5" style="color: var(--text-muted);">Se recuperaron {{ recoveredPoints }} puntos de una sesión anterior.</p>
          </div>
        </div>
      </div>

      <!-- Status card -->
      <div class="glass-card p-4 animate-in" style="animation-delay: 0.05s;">
        <div class="flex items-center gap-3">
          <div :class="activeRecordId ? 'bg-emerald-500/20' : ''"
            class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
            :style="!activeRecordId ? 'background: var(--input-bg)' : ''">
            <SignalIcon v-if="activeRecordId" class="w-5 h-5 text-emerald-400" />
            <MinusCircleIcon v-else class="w-5 h-5" style="color: var(--text-muted);" />
          </div>
          <div>
            <p :class="activeRecordId ? 'text-emerald-400' : ''" class="font-semibold text-sm"
              :style="!activeRecordId ? 'color: var(--text-muted)' : ''">
              {{ activeRecordId ? 'Entrada registrada actualmente' : 'Sin entrada activa' }}
            </p>
            <p v-if="entryTime" class="text-xs" style="color: var(--text-muted);">Desde las {{ entryTime }}</p>
          </div>
        </div>
      </div>

      <!-- Main buttons -->
      <div class="space-y-3 animate-in" style="animation-delay: 0.1s;">
        <button
          v-if="!activeRecordId"
          @click="startCheckProcess('entry')"
          :disabled="processing"
          class="btn-success btn-xl w-full"
          style="background: linear-gradient(135deg, #059669, #10b981); box-shadow: 0 4px 20px rgba(16,185,129,0.3);"
        >
          <MapPinIcon class="w-6 h-6" />
          <span>Registrar Entrada</span>
        </button>

        <button
          v-if="activeRecordId"
          @click="startCheckProcess('exit')"
          :disabled="processing"
          class="btn-danger btn-xl w-full"
          style="background: linear-gradient(135deg, #dc2626, #f43f5e); box-shadow: 0 4px 20px rgba(244,63,94,0.3);"
        >
          <ArrowRightOnRectangleIcon class="w-6 h-6" />
          <span>Registrar Salida</span>
        </button>
      </div>
    </div>

    <!-- ===== Modals ===== -->
    <Teleport to="body">

      <!-- Camera modal -->
      <div v-if="showCameraModal" class="fixed inset-0 bg-black z-50 flex flex-col">
        <!-- Camera instruction bar -->
        <div class="relative z-10 text-center py-5 px-6 safe-top"
          style="background: linear-gradient(to bottom, rgba(0,0,0,0.85), transparent);">
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full text-white text-sm font-semibold"
            style="background: rgba(99,102,241,0.5); border: 1px solid rgba(99,102,241,0.5);">
            <CameraIcon v-if="cameraStep === 'site'" class="w-4 h-4" />
            <UserCircleIcon v-else class="w-4 h-4" />
            <span>{{ cameraStep === 'site' ? 'Toma fotografía del sitio' : 'Tómate una selfie en el sitio' }}</span>
          </div>
          <p class="text-white/60 text-xs mt-2">
            {{ cameraStep === 'site' ? 'Fotografía el lugar donde te encuentras' : 'Incluye claramente tu rostro' }}
          </p>
        </div>

        <!-- Video / Preview -->
        <div class="flex-1 relative overflow-hidden">
          <video ref="videoRef" autoplay playsinline muted
            class="absolute inset-0 w-full h-full object-cover"
            :class="{ 'scale-x-[-1]': facingFront }"></video>
          <canvas ref="canvasRef" class="hidden"></canvas>

          <!-- Photo preview overlay -->
          <div v-if="capturedPhoto" class="absolute inset-0 bg-black">
            <img :src="capturedPhoto" class="w-full h-full object-cover"
              :class="{ 'scale-x-[-1]': cameraStep === 'selfie' }" />
            <!-- Preview overlay badge -->
            <div class="absolute top-4 left-1/2 -translate-x-1/2 px-4 py-2 rounded-full text-sm font-semibold"
              style="background: rgba(0,0,0,0.7); border: 1px solid rgba(255,255,255,0.2); color: white;">
              Vista previa
            </div>
          </div>

          <!-- Guide frame -->
          <div v-if="!capturedPhoto && cameraStep === 'selfie'"
            class="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div class="w-48 h-60 rounded-3xl border-2 border-white/30 border-dashed"></div>
          </div>
        </div>

        <!-- Controls bar -->
        <div class="safe-bottom py-6 px-6"
          style="background: linear-gradient(to top, rgba(0,0,0,0.9), transparent 80%);">

          <!-- Before capture -->
          <div v-if="!capturedPhoto" class="flex items-center justify-between max-w-xs mx-auto">
            <button @click="toggleFlash"
              class="w-12 h-12 rounded-full flex items-center justify-center text-2xl transition-all"
              :style="flashOn ? 'background: rgba(251,191,36,0.3); border: 2px solid rgba(251,191,36,0.7);' : 'background: rgba(255,255,255,0.1); border: 2px solid rgba(255,255,255,0.2);'">
              <BoltIcon class="w-6 h-6" :class="flashOn ? 'text-amber-300' : 'text-white'" />
            </button>

            <button @click="capturePhoto"
              class="w-20 h-20 rounded-full flex items-center justify-center transition-all active:scale-90"
              style="background: white; border: 4px solid rgba(255,255,255,0.3); box-shadow: 0 0 0 6px rgba(255,255,255,0.1);">
              <div class="w-14 h-14 rounded-full border-2 border-gray-300" style="background: white;"></div>
            </button>

            <button @click="switchCamera"
              class="w-12 h-12 rounded-full flex items-center justify-center text-2xl transition-all"
              style="background: rgba(255,255,255,0.1); border: 2px solid rgba(255,255,255,0.2);">
              <ArrowPathIcon class="w-6 h-6 text-white" />
            </button>
          </div>

          <!-- After capture -->
          <div v-if="capturedPhoto" class="flex gap-3 max-w-xs mx-auto">
            <button @click="retakePhoto" class="btn-secondary btn-lg flex-1">
              <ArrowPathIcon class="w-4 h-4" /> Retomar
            </button>
            <button @click="confirmPhoto" class="btn-primary btn-lg flex-1">
              <CheckCircleIcon class="w-4 h-4" /> Usar foto
            </button>
          </div>
        </div>
      </div>

      <!-- Success modal -->
      <Transition name="modal">
        <div v-if="showSuccessModal" class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0,0,0,0.75); backdrop-filter: blur(8px);">
          <div class="w-full max-w-sm glass-card p-8 text-center animate-in">
            <div class="w-20 h-20 rounded-3xl flex items-center justify-center mx-auto mb-4"
              :style="checkType === 'entry'
                ? 'background: linear-gradient(135deg, rgba(5,150,105,0.3), rgba(16,185,129,0.3)); border: 1px solid rgba(16,185,129,0.4);'
                : 'background: linear-gradient(135deg, rgba(220,38,38,0.3), rgba(244,63,94,0.3)); border: 1px solid rgba(244,63,94,0.4);'">
              <CheckCircleIcon v-if="checkType === 'entry'" class="w-10 h-10 text-emerald-400" />
              <HandRaisedIcon v-else class="w-10 h-10 text-rose-400" />
            </div>
            <h3 class="text-2xl font-bold mb-1" style="color: var(--text);">
              {{ checkType === 'entry' ? '¡Entrada Registrada!' : '¡Salida Registrada!' }}
            </h3>
            <p style="color: var(--text-muted);" class="text-sm mb-1">{{ registeredAt }}</p>
            <p class="text-xs mb-6" style="color: var(--text-muted);">Ubicación y fotografías guardadas exitosamente</p>
            <button @click="closeSuccess" class="btn-primary btn-lg w-full">Aceptar</button>
          </div>
        </div>
      </Transition>

      <!-- Processing spinner -->
      <div v-if="processing && !showCameraModal && !showSuccessModal"
        class="fixed inset-0 z-50 flex items-center justify-center"
        style="background: rgba(0,0,0,0.5); backdrop-filter: blur(4px);">
        <div class="glass-card p-6 text-center">
          <div class="w-10 h-10 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin mx-auto mb-3"></div>
          <p class="font-semibold" style="color: var(--text);">Procesando...</p>
          <p class="text-xs mt-1" style="color: var(--text-muted);">Por favor espera</p>
        </div>
      </div>

    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import {
  CheckCircleIcon, HandRaisedIcon, MapPinIcon, ArrowRightOnRectangleIcon,
  ExclamationTriangleIcon, XMarkIcon, LockClosedIcon, CloudArrowUpIcon,
  SignalIcon, MinusCircleIcon, CameraIcon, UserCircleIcon, BoltIcon, ArrowPathIcon
} from '@heroicons/vue/24/outline'
import ThemeToggle from '@/components/ThemeToggle.vue'

const auth = useAuthStore()

// Clock
const currentTime = ref('')
const todayStr = ref('')
let timeInterval = null

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  todayStr.value = now.toLocaleDateString('es-MX', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
}


// State
const activeRecordId = ref(localStorage.getItem('activeRecordId') || null)
const entryTime = ref(localStorage.getItem('entryTime') || null)
const processing = ref(false)
const checkType = ref('entry')
const showCameraModal = ref(false)
const showSuccessModal = ref(false)
const registeredAt = ref('')

// Location
let locationWatchId = null
const locationPoints = []

function startCheckProcess(type) {
  checkType.value = type
  requestLocation()
}

function cancelProcess() {
  processing.value = false
}

async function requestLocation() {
  processing.value = true
  try {
    await new Promise((resolve, reject) =>
      navigator.geolocation.getCurrentPosition(resolve, reject, { enableHighAccuracy: true })
    )
    locationWatchId = navigator.geolocation.watchPosition(
      pos => locationPoints.push({
        latitude: pos.coords.latitude,
        longitude: pos.coords.longitude,
        accuracy: pos.coords.accuracy,
        recorded_at: new Date().toISOString()
      }),
      null, { enableHighAccuracy: true, maximumAge: 5000 }
    )
    cameraStep.value = 'site'
    await startCamera('environment')
    showCameraModal.value = true
    processing.value = false
  } catch {
    alert('No se pudo obtener la ubicación. Por favor permite el acceso en tu navegador.')
    processing.value = false
  }
}

function stopLocationTracking() {
  if (locationWatchId !== null) {
    navigator.geolocation.clearWatch(locationWatchId)
    locationWatchId = null
  }
}

// Camera
const videoRef = ref(null)
const canvasRef = ref(null)
const capturedPhoto = ref(null)
const cameraStep = ref('site')
const facingFront = ref(false)
const flashOn = ref(false)
let stream = null
let photoSite = null
let photoSelfie = null

async function startCamera(facing = 'environment') {
  stopCamera()
  facingFront.value = facing === 'user'
  stream = await navigator.mediaDevices.getUserMedia({
    video: { facingMode: facing, width: { ideal: 1280 }, height: { ideal: 720 } },
    audio: false
  })
  if (videoRef.value) videoRef.value.srcObject = stream
}

function stopCamera() {
  if (stream) { stream.getTracks().forEach(t => t.stop()); stream = null }
}

async function switchCamera() {
  await startCamera(facingFront.value ? 'environment' : 'user')
}

function toggleFlash() {
  const track = stream?.getVideoTracks()[0]
  if (!track) return
  const caps = track.getCapabilities?.()
  if (caps?.torch) {
    flashOn.value = !flashOn.value
    track.applyConstraints({ advanced: [{ torch: flashOn.value }] })
  } else {
    alert('Tu dispositivo no soporta flash desde el navegador')
  }
}

function capturePhoto() {
  const video = videoRef.value
  const canvas = canvasRef.value
  if (!video || !canvas) return
  canvas.width = video.videoWidth
  canvas.height = video.videoHeight
  const ctx = canvas.getContext('2d')
  if (cameraStep.value === 'selfie') { ctx.translate(canvas.width, 0); ctx.scale(-1, 1) }
  ctx.drawImage(video, 0, 0)
  capturedPhoto.value = canvas.toDataURL('image/jpeg', 0.85)
}

function retakePhoto() { capturedPhoto.value = null }

async function confirmPhoto() {
  if (cameraStep.value === 'site') {
    photoSite = capturedPhoto.value
    capturedPhoto.value = null
    cameraStep.value = 'selfie'
    await startCamera('user')
  } else {
    photoSelfie = capturedPhoto.value
    capturedPhoto.value = null
    showCameraModal.value = false
    stopCamera()
    await submitCheck()
  }
}

function dataURLtoBlob(dataURL) {
  const [header, data] = dataURL.split(',')
  const mime = header.match(/:(.*?);/)[1]
  const bstr = atob(data)
  const u8arr = new Uint8Array(bstr.length)
  for (let i = 0; i < bstr.length; i++) u8arr[i] = bstr.charCodeAt(i)
  return new Blob([u8arr], { type: mime })
}

async function submitCheck() {
  processing.value = true
  try {
    const now = new Date().toISOString()
    const formData = new FormData()
    formData.append('type', checkType.value)
    formData.append('timestamp', now)
    if (photoSite) formData.append('photo_site', dataURLtoBlob(photoSite), 'site.jpg')
    if (photoSelfie) formData.append('photo_selfie', dataURLtoBlob(photoSelfie), 'selfie.jpg')

    const { data } = await api.post('/checks', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
    const recordId = data.record.id

    for (const point of locationPoints.value) {
      try {
        await api.post('/location-points', { ...point, check_record_id: recordId })
      } catch { /* Skip individual failed points */ }
    }
    locationPoints.value = []
    clearLocationBuffer()
    stopLocationTracking()
    recoveredPoints.value = 0

    if (checkType.value === 'entry') {
      activeRecordId.value = recordId
      entryTime.value = new Date(now).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' })
      localStorage.setItem('activeRecordId', recordId)
      localStorage.setItem('entryTime', entryTime.value)
    } else {
      activeRecordId.value = null
      entryTime.value = null
      localStorage.removeItem('activeRecordId')
      localStorage.removeItem('entryTime')
    }

    registeredAt.value = new Date(now).toLocaleString('es-MX')
    showSuccessModal.value = true
  } catch {
    alert('Error al registrar. Intenta de nuevo.')
  } finally {
    processing.value = false
    photoSite = null
    photoSelfie = null
  }
}

function closeSuccess() { showSuccessModal.value = false }

// ─── Lifecycle ────────────────────────────────────────────────────────────────
onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  document.addEventListener('visibilitychange', handleVisibilityChange)
  // Recover buffered points if app was closed during an active session
  const buffered = loadLocationBuffer()
  if (buffered.length > 0 && activeRecordId.value) {
    recoveredPoints.value = buffered.length
    locationPoints.value = buffered
  }
})

onUnmounted(() => {
  clearInterval(timeInterval)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  stopLocationTracking()
  stopCamera()
})
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.25s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.safe-top { padding-top: max(16px, env(safe-area-inset-top)); }
.safe-bottom { padding-bottom: max(24px, env(safe-area-inset-bottom)); }
.slide-banner-enter-active, .slide-banner-leave-active { transition: all 0.3s ease; }
.slide-banner-enter-from, .slide-banner-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
