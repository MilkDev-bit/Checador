<template>
  <div class="min-h-screen min-h-dvh bg-mesh" style="background-color: #0f1629;">
    <!-- Top nav -->
    <header class="sticky top-0 z-20 px-4 py-3 flex items-center justify-between"
      style="background: rgba(15,22,41,0.85); backdrop-filter: blur(12px); border-bottom: 1px solid rgba(255,255,255,0.06);">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center text-sm"
          style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">✅</div>
        <div>
          <p class="text-white font-semibold text-sm leading-none">PaseLista</p>
          <p class="text-slate-500 text-xs mt-0.5 truncate max-w-[160px]">{{ auth.user?.project_name }}</p>
        </div>
      </div>
      <div class="flex items-center gap-1">
        <router-link to="/qr" class="btn-ghost text-xs px-3 py-2">QR</router-link>
        <router-link to="/history" class="btn-ghost text-xs px-3 py-2">Historial</router-link>
        <button @click="auth.logout(); $router.push('/login')" class="btn-ghost text-xs px-3 py-2 text-rose-400">Salir</button>
      </div>
    </header>

    <div class="max-w-lg mx-auto px-4 pt-6 pb-10 space-y-5">
      <!-- Greeting card -->
      <div class="glass-card p-5 animate-in" style="background: rgba(20,29,53,0.8);">
        <div class="flex items-center gap-4">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-2xl flex-shrink-0"
            style="background: linear-gradient(135deg, rgba(99,102,241,0.3), rgba(139,92,246,0.3)); border: 1px solid rgba(99,102,241,0.3);">
            👋
          </div>
          <div class="min-w-0">
            <h2 class="text-lg font-bold text-white truncate">Hola, {{ auth.user?.first_name }}</h2>
            <p class="text-slate-400 text-sm">{{ todayStr }}</p>
            <p class="text-brand-400 font-mono font-semibold text-base mt-0.5">{{ currentTime }}</p>
          </div>
        </div>
      </div>

      <!-- Status card -->
      <div class="glass-card p-4 animate-in" style="animation-delay: 0.05s; background: rgba(20,29,53,0.8);">
        <div class="flex items-center gap-3">
          <div :class="activeRecordId ? 'bg-emerald-500/20' : 'bg-slate-700/50'"
            class="w-10 h-10 rounded-xl flex items-center justify-center text-xl flex-shrink-0">
            {{ activeRecordId ? '🟢' : '⚪' }}
          </div>
          <div>
            <p :class="activeRecordId ? 'text-emerald-400' : 'text-slate-400'" class="font-semibold text-sm">
              {{ activeRecordId ? 'Entrada registrada actualmente' : 'Sin entrada activa' }}
            </p>
            <p v-if="entryTime" class="text-slate-500 text-xs">Desde las {{ entryTime }}</p>
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
          <span class="text-2xl">📍</span>
          <span>Registrar Entrada</span>
        </button>

        <button
          v-if="activeRecordId"
          @click="startCheckProcess('exit')"
          :disabled="processing"
          class="btn-danger btn-xl w-full"
          style="background: linear-gradient(135deg, #dc2626, #f43f5e); box-shadow: 0 4px 20px rgba(244,63,94,0.3);"
        >
          <span class="text-2xl">🚪</span>
          <span>Registrar Salida</span>
        </button>
      </div>
    </div>

    <!-- ===== Modals ===== -->
    <Teleport to="body">

      <!-- Location permission modal -->
      <Transition name="modal">
        <div v-if="showLocationModal" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center px-4 pb-4 sm:pb-0"
          style="background: rgba(0,0,0,0.7); backdrop-filter: blur(6px);">
          <div class="w-full max-w-sm glass-card p-7 text-center animate-in"
            style="background: rgba(20,29,53,0.97);">
            <div class="w-16 h-16 rounded-2xl flex items-center justify-center text-3xl mx-auto mb-4"
              style="background: linear-gradient(135deg, rgba(99,102,241,0.3), rgba(6,182,212,0.3)); border: 1px solid rgba(99,102,241,0.3);">
              📍
            </div>
            <h3 class="text-xl font-bold text-white mb-2">Compartir Ubicación</h3>
            <p class="text-slate-400 text-sm mb-6 leading-relaxed">
              Para registrar tu {{ checkType === 'entry' ? 'entrada' : 'salida' }}, necesitamos tu ubicación en tiempo real para registrar el recorrido completo.
            </p>
            <div class="space-y-3">
              <button @click="requestLocation" class="btn-primary btn-lg w-full">
                <span>📍</span> Permitir Ubicación
              </button>
              <button @click="cancelProcess" class="btn-secondary btn-md w-full">Cancelar</button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Camera modal -->
      <div v-if="showCameraModal" class="fixed inset-0 bg-black z-50 flex flex-col">
        <!-- Camera instruction bar -->
        <div class="relative z-10 text-center py-5 px-6 safe-top"
          style="background: linear-gradient(to bottom, rgba(0,0,0,0.85), transparent);">
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full text-white text-sm font-semibold"
            style="background: rgba(99,102,241,0.5); border: 1px solid rgba(99,102,241,0.5);">
            <span>{{ cameraStep === 'site' ? '📸' : '🤳' }}</span>
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
              {{ flashOn ? '⚡' : '🔦' }}
            </button>

            <button @click="capturePhoto"
              class="w-20 h-20 rounded-full flex items-center justify-center transition-all active:scale-90"
              style="background: white; border: 4px solid rgba(255,255,255,0.3); box-shadow: 0 0 0 6px rgba(255,255,255,0.1);">
              <div class="w-14 h-14 rounded-full border-2 border-gray-300" style="background: white;"></div>
            </button>

            <button @click="switchCamera"
              class="w-12 h-12 rounded-full flex items-center justify-center text-2xl transition-all"
              style="background: rgba(255,255,255,0.1); border: 2px solid rgba(255,255,255,0.2);">
              🔄
            </button>
          </div>

          <!-- After capture -->
          <div v-if="capturedPhoto" class="flex gap-3 max-w-xs mx-auto">
            <button @click="retakePhoto" class="btn-secondary btn-lg flex-1">🔄 Retomar</button>
            <button @click="confirmPhoto" class="btn-primary btn-lg flex-1">✅ Usar foto</button>
          </div>
        </div>
      </div>

      <!-- Success modal -->
      <Transition name="modal">
        <div v-if="showSuccessModal" class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0,0,0,0.75); backdrop-filter: blur(8px);">
          <div class="w-full max-w-sm glass-card p-8 text-center animate-in"
            style="background: rgba(20,29,53,0.97);">
            <div class="w-20 h-20 rounded-3xl flex items-center justify-center text-4xl mx-auto mb-4"
              :style="checkType === 'entry'
                ? 'background: linear-gradient(135deg, rgba(5,150,105,0.3), rgba(16,185,129,0.3)); border: 1px solid rgba(16,185,129,0.4);'
                : 'background: linear-gradient(135deg, rgba(220,38,38,0.3), rgba(244,63,94,0.3)); border: 1px solid rgba(244,63,94,0.4);'">
              {{ checkType === 'entry' ? '✅' : '👋' }}
            </div>
            <h3 class="text-2xl font-bold text-white mb-1">
              {{ checkType === 'entry' ? '¡Entrada Registrada!' : '¡Salida Registrada!' }}
            </h3>
            <p class="text-slate-400 text-sm mb-1">{{ registeredAt }}</p>
            <p class="text-slate-500 text-xs mb-6">Ubicación y fotografías guardadas exitosamente</p>
            <button @click="closeSuccess" class="btn-primary btn-lg w-full">Aceptar</button>
          </div>
        </div>
      </Transition>

      <!-- Processing spinner -->
      <div v-if="processing && !showLocationModal && !showCameraModal && !showSuccessModal"
        class="fixed inset-0 z-50 flex items-center justify-center"
        style="background: rgba(0,0,0,0.5); backdrop-filter: blur(4px);">
        <div class="glass-card p-6 text-center" style="background: rgba(20,29,53,0.95);">
          <div class="w-10 h-10 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin mx-auto mb-3"></div>
          <p class="text-white font-semibold">Procesando...</p>
          <p class="text-slate-400 text-xs mt-1">Por favor espera</p>
        </div>
      </div>

    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

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

onMounted(() => { updateTime(); timeInterval = setInterval(updateTime, 1000) })
onUnmounted(() => { clearInterval(timeInterval); stopLocationTracking(); stopCamera() })

// State
const activeRecordId = ref(localStorage.getItem('activeRecordId') || null)
const entryTime = ref(localStorage.getItem('entryTime') || null)
const processing = ref(false)
const checkType = ref('entry')
const showLocationModal = ref(false)
const showCameraModal = ref(false)
const showSuccessModal = ref(false)
const registeredAt = ref('')

// Location
let locationWatchId = null
const locationPoints = []

function startCheckProcess(type) {
  checkType.value = type
  showLocationModal.value = true
}

function cancelProcess() {
  showLocationModal.value = false
  processing.value = false
}

async function requestLocation() {
  showLocationModal.value = false
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

    for (const point of locationPoints) {
      await api.post('/location-points', { ...point, check_record_id: recordId })
    }
    locationPoints.length = 0
    stopLocationTracking()

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
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.25s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.safe-top { padding-top: max(16px, env(safe-area-inset-top)); }
.safe-bottom { padding-bottom: max(24px, env(safe-area-inset-bottom)); }
</style>

    <!-- Header -->
    <header class="bg-white shadow-sm px-4 py-4 flex items-center justify-between">
      <div>
        <h1 class="font-bold text-blue-700 text-lg">PaseLista</h1>
        <p class="text-xs text-gray-500">{{ auth.user?.project_name }}</p>
      </div>
      <div class="flex items-center gap-3">
        <router-link to="/qr" class="text-blue-500 hover:text-blue-700 text-sm font-medium">QR</router-link>
        <router-link to="/history" class="text-blue-500 hover:text-blue-700 text-sm font-medium">Historial</router-link>
        <button @click="auth.logout(); $router.push('/login')" class="text-red-500 hover:text-red-700 text-sm font-medium">Salir</button>
      </div>
    </header>

    <div class="max-w-lg mx-auto px-4 pt-8 space-y-6">
      <!-- User greeting -->
      <div class="card text-center">
        <div class="text-4xl mb-2">👋</div>
        <h2 class="text-xl font-bold text-gray-800">Hola, {{ auth.user?.first_name }}</h2>
        <p class="text-gray-500 text-sm">{{ new Date().toLocaleDateString('es-MX', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}</p>
        <p class="text-gray-700 font-medium mt-1">{{ currentTime }}</p>
      </div>

      <!-- Status -->
      <div class="card text-center">
        <div :class="activeRecordId ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'" class="rounded-full px-4 py-2 inline-block text-sm font-semibold mb-3">
          {{ activeRecordId ? '🟢 Entrada registrada' : '⚪ Sin entrada activa' }}
        </div>
        <p v-if="entryTime" class="text-sm text-gray-500">Entrada: {{ entryTime }}</p>
      </div>

      <!-- Main action buttons -->
      <div class="space-y-4">
        <button
          v-if="!activeRecordId"
          @click="startCheckProcess('entry')"
          :disabled="processing"
          class="btn-success w-full flex items-center justify-center gap-3"
        >
          <span class="text-2xl">📍</span>
          <span>Registrar Entrada</span>
        </button>

        <button
          v-if="activeRecordId"
          @click="startCheckProcess('exit')"
          :disabled="processing"
          class="btn-danger w-full flex items-center justify-center gap-3"
        >
          <span class="text-2xl">🚪</span>
          <span>Registrar Salida</span>
        </button>
      </div>
    </div>

    <!-- Modal: Location permission -->
    <Teleport to="body">
      <div v-if="showLocationModal" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50 px-4">
        <div class="card max-w-sm w-full text-center">
          <div class="text-5xl mb-4">📍</div>
          <h3 class="text-xl font-bold mb-2">Compartir Ubicación</h3>
          <p class="text-gray-600 mb-6">
            Para registrar tu {{ checkType === 'entry' ? 'entrada' : 'salida' }}, necesitamos acceso a tu ubicación en tiempo real para registrar tu recorrido.
          </p>
          <div class="space-y-3">
            <button @click="requestLocation" class="btn-primary w-full">Permitir Ubicación</button>
            <button @click="cancelProcess" class="btn-secondary w-full">Cancelar</button>
          </div>
        </div>
      </div>

      <!-- Camera: Site photo -->
      <div v-if="showCameraModal" class="fixed inset-0 bg-black z-50 flex flex-col">
        <div class="bg-black/80 text-white text-center py-4 px-6">
          <p class="text-lg font-bold">
            {{ cameraStep === 'site' ? '📸 Toma fotografía del sitio' : '🤳 Tómate una selfie en el sitio' }}
          </p>
          <p class="text-sm text-gray-300 mt-1">
            {{ cameraStep === 'site' ? 'Fotografía el lugar donde te encuentras' : 'Incluye tu rostro en la foto' }}
          </p>
        </div>

        <div class="flex-1 relative">
          <video ref="videoRef" autoplay playsinline class="w-full h-full object-cover" :class="facingFront ? 'scale-x-[-1]' : ''"></video>
          <canvas ref="canvasRef" class="hidden"></canvas>

          <!-- Photo preview -->
          <div v-if="capturedPhoto" class="absolute inset-0 bg-black">
            <img :src="capturedPhoto" class="w-full h-full object-cover" :class="cameraStep === 'selfie' ? 'scale-x-[-1]' : ''" />
          </div>
        </div>

        <div class="bg-black/90 py-6 px-6 space-y-3">
          <!-- Camera controls -->
          <div v-if="!capturedPhoto" class="flex items-center justify-between">
            <button @click="toggleFlash" class="text-white text-3xl p-2" title="Flash">
              {{ flashOn ? '⚡' : '🔦' }}
            </button>
            <button @click="capturePhoto" class="w-20 h-20 rounded-full bg-white border-4 border-gray-300 flex items-center justify-center shadow-lg active:scale-95 transition-transform">
              <div class="w-14 h-14 rounded-full bg-white border-2 border-gray-400"></div>
            </button>
            <button @click="switchCamera" class="text-white text-3xl p-2" title="Cambiar cámara">🔄</button>
          </div>

          <!-- After capture -->
          <div v-if="capturedPhoto" class="flex gap-3">
            <button @click="retakePhoto" class="btn-secondary flex-1">Retomar</button>
            <button @click="confirmPhoto" class="btn-primary flex-1">Usar foto</button>
          </div>
        </div>
      </div>

      <!-- Success modal -->
      <div v-if="showSuccessModal" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50 px-4">
        <div class="card max-w-sm w-full text-center">
          <div class="text-6xl mb-4">✅</div>
          <h3 class="text-2xl font-bold text-green-700 mb-2">
            {{ checkType === 'entry' ? '¡Entrada Registrada!' : '¡Salida Registrada!' }}
          </h3>
          <p class="text-gray-600 mb-2">{{ registeredAt }}</p>
          <p class="text-sm text-gray-500 mb-6">Tu {{ checkType === 'entry' ? 'entrada' : 'salida' }} ha sido registrada exitosamente con tu ubicación y fotografías.</p>
          <button @click="closeSuccess" class="btn-primary w-full">Aceptar</button>
        </div>
      </div>

      <!-- Processing overlay -->
      <div v-if="processing && !showLocationModal && !showCameraModal && !showSuccessModal" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
        <div class="card text-center">
          <div class="text-4xl mb-3 animate-spin">⏳</div>
          <p class="font-semibold text-gray-700">Procesando...</p>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const auth = useAuthStore()

// Time
const currentTime = ref('')
let timeInterval = null

function updateTime() {
  currentTime.value = new Date().toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
})
onUnmounted(() => {
  clearInterval(timeInterval)
  stopLocationTracking()
  stopCamera()
})

// State
const activeRecordId = ref(localStorage.getItem('activeRecordId') || null)
const entryTime = ref(localStorage.getItem('entryTime') || null)
const processing = ref(false)
const checkType = ref('entry')

// Modals
const showLocationModal = ref(false)
const showCameraModal = ref(false)
const showSuccessModal = ref(false)
const registeredAt = ref('')

// Location tracking
let locationWatchId = null
const locationPoints = []

function startCheckProcess(type) {
  checkType.value = type
  showLocationModal.value = true
}

function cancelProcess() {
  showLocationModal.value = false
  processing.value = false
}

async function requestLocation() {
  showLocationModal.value = false
  processing.value = true

  try {
    await new Promise((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resolve, reject, { enableHighAccuracy: true })
    })

    // Start tracking route
    locationWatchId = navigator.geolocation.watchPosition(
      pos => {
        locationPoints.push({
          latitude: pos.coords.latitude,
          longitude: pos.coords.longitude,
          accuracy: pos.coords.accuracy,
          recorded_at: new Date().toISOString()
        })
      },
      null,
      { enableHighAccuracy: true, maximumAge: 5000 }
    )

    // Proceed to camera
    cameraStep.value = 'site'
    await startCamera()
    showCameraModal.value = true
    processing.value = false
  } catch (e) {
    alert('No se pudo obtener la ubicación. Por favor permite el acceso a la ubicación.')
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
const cameraStep = ref('site') // 'site' | 'selfie'
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
  if (videoRef.value) {
    videoRef.value.srcObject = stream
  }
}

function stopCamera() {
  if (stream) {
    stream.getTracks().forEach(t => t.stop())
    stream = null
  }
}

async function switchCamera() {
  const newFacing = facingFront.value ? 'environment' : 'user'
  await startCamera(newFacing)
}

function toggleFlash() {
  const track = stream?.getVideoTracks()[0]
  if (!track) return
  const caps = track.getCapabilities()
  if (caps.torch) {
    flashOn.value = !flashOn.value
    track.applyConstraints({ advanced: [{ torch: flashOn.value }] })
  } else {
    alert('Tu dispositivo no soporta flash')
  }
}

function capturePhoto() {
  const video = videoRef.value
  const canvas = canvasRef.value
  if (!video || !canvas) return

  canvas.width = video.videoWidth
  canvas.height = video.videoHeight
  const ctx = canvas.getContext('2d')

  if (cameraStep.value === 'selfie') {
    ctx.translate(canvas.width, 0)
    ctx.scale(-1, 1)
  }

  ctx.drawImage(video, 0, 0)
  capturedPhoto.value = canvas.toDataURL('image/jpeg', 0.85)
}

function retakePhoto() {
  capturedPhoto.value = null
}

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
  const arr = dataURL.split(',')
  const mime = arr[0].match(/:(.*?);/)[1]
  const bstr = atob(arr[1])
  let n = bstr.length
  const u8arr = new Uint8Array(n)
  while (n--) u8arr[n] = bstr.charCodeAt(n)
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

    const { data } = await api.post('/checks', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    const recordId = data.record.id

    // Send all accumulated location points
    for (const point of locationPoints) {
      await api.post('/location-points', { ...point, check_record_id: recordId })
    }
    locationPoints.length = 0

    stopLocationTracking()

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
  } catch (e) {
    alert('Error al registrar. Intenta de nuevo.')
  } finally {
    processing.value = false
    photoSite = null
    photoSelfie = null
  }
}

function closeSuccess() {
  showSuccessModal.value = false
}
</script>
