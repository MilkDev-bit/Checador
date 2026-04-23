<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 pb-20">
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
