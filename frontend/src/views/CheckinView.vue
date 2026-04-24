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
        <button @click="showLogoutModal = true" class="btn-ghost text-xs px-3 py-2 text-rose-400">Salir</button>
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

      <!-- Location error / guide modal -->
      <Transition name="modal">
        <div v-if="showLocationErrorModal" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center px-4 pb-4 sm:pb-0"
          style="background: rgba(0,0,0,0.75); backdrop-filter: blur(6px);">
          <div class="w-full max-w-sm glass-card p-6 animate-in">
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
                style="background: rgba(239,68,68,0.15); border: 1px solid rgba(239,68,68,0.3);">
                <MapPinIcon class="w-5 h-5 text-rose-400" />
              </div>
              <div>
                <h3 class="font-bold text-base" style="color: var(--text);">
                  {{ locationErrorType === 'inapp' ? 'Abre la página en Safari' :
                     locationErrorType === 'unavailable' ? 'GPS no disponible' :
                     locationErrorType === 'timeout' ? 'Tiempo de espera agotado' :
                     'Ubicación bloqueada' }}
                </h3>
                <p class="text-xs" style="color: var(--text-muted);">
                  {{ locationErrorType === 'inapp' ? 'El browser de WhatsApp no tiene acceso al GPS' :
                     locationErrorType === 'unavailable' ? 'No se detectó señal GPS' :
                     locationErrorType === 'timeout' ? 'El GPS tardó demasiado en responder' :
                     'Sigue los pasos de abajo para activarla' }}
                  <span v-if="locationErrorCode" class="ml-1 opacity-50">(código {{ locationErrorCode }})</span>
                </p>
              </div>
            </div>

            <!-- In-app browser (WKWebView) -->
            <template v-if="locationErrorType === 'inapp'">
              <div class="rounded-xl p-4 mb-5" style="background: rgba(99,102,241,0.07); border: 1px solid rgba(99,102,241,0.2);">
                <p class="text-xs font-semibold mb-3 text-brand-400">Cómo abrir en Safari</p>
                <ol class="space-y-2">
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">1.</span>
                    Toca el ícono de <strong style="color: var(--text);">tres puntos (···)</strong> o el ícono de <strong style="color: var(--text);">compartir</strong> en la barra inferior de WhatsApp.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">2.</span>
                    Selecciona <strong style="color: var(--text);">"Abrir en Safari"</strong> o <strong style="color: var(--text);">"Abrir en navegador"</strong>.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">3.</span>
                    Inicia sesión nuevamente y el GPS funcionará correctamente.
                  </li>
                </ol>
              </div>
            </template>

            <!-- GPS unavailable or timeout -->
            <template v-else-if="locationErrorType === 'unavailable' || locationErrorType === 'timeout'">
              <div class="rounded-xl p-4 mb-5" style="background: rgba(245,158,11,0.07); border: 1px solid rgba(245,158,11,0.2);">
                <ol class="space-y-2">
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-amber-400 flex-shrink-0">1.</span>
                    Ve a un lugar con buena señal o al aire libre.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-amber-400 flex-shrink-0">2.</span>
                    Asegúrate de que el <strong style="color: var(--text);">GPS del teléfono</strong> esté activado.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-amber-400 flex-shrink-0">3.</span>
                    Vuelve a intentarlo.
                  </li>
                </ol>
              </div>
            </template>

            <!-- Permission denied -->
            <template v-else>
              <!-- Paso 0: Verificar Servicios de Localización globales -->
              <div class="rounded-xl p-3 mb-3 flex items-start gap-2"
                style="background: rgba(239,68,68,0.08); border: 1px solid rgba(239,68,68,0.25);">
                <ExclamationTriangleIcon class="w-4 h-4 text-rose-400 flex-shrink-0 mt-0.5" />
                <p class="text-xs leading-relaxed" style="color: var(--text-muted);">
                  iOS bloqueó la ubicación. Hay <strong style="color: var(--text);">dos lugares</strong> donde revisar:
                </p>
              </div>

              <!-- Paso 1: Configuración iOS global (causa más común) -->
              <div class="rounded-xl p-4 mb-3" style="background: rgba(99,102,241,0.10); border: 2px solid rgba(99,102,241,0.35);">
                <p class="text-xs font-bold mb-2" style="color: var(--text);">① Configuración del iPhone (revisa primero)</p>
                <ol class="space-y-1.5">
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">1.</span>
                    Abre <strong style="color: var(--text);">Configuración → Privacidad y Seguridad → Servicios de Localización</strong>.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">2.</span>
                    Asegúrate de que los Servicios de Localización estén <strong style="color: var(--text);">activados</strong> (toggle verde).
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">3.</span>
                    En la lista busca <strong style="color: var(--text);">Safari</strong> → ponlo en <strong style="color: var(--text);">"Al usar la app"</strong>.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">4.</span>
                    Regresa aquí y toca <strong style="color: var(--text);">"Recargar página"</strong>.
                  </li>
                </ol>
              </div>

              <!-- Paso 2: AA del sitio en Safari -->
              <div class="rounded-xl p-4 mb-3" style="background: rgba(99,102,241,0.05); border: 1px solid rgba(99,102,241,0.2);">
                <p class="text-xs font-bold mb-2" style="color: var(--text);">② Permiso del sitio en Safari</p>
                <ol class="space-y-1.5">
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">1.</span>
                    Toca las letras <strong style="color: var(--text);">"AA"</strong> a la izquierda de la barra de dirección de Safari.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">2.</span>
                    Toca <strong style="color: var(--text);">"Ajustes del sitio web"</strong>.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">3.</span>
                    En <strong style="color: var(--text);">Ubicación</strong> selecciona <strong style="color: var(--text);">"Permitir"</strong>.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-brand-400 flex-shrink-0">4.</span>
                    Regresa aquí y toca <strong style="color: var(--text);">"Recargar página"</strong>.
                  </li>
                </ol>
              </div>

              <!-- Android -->
              <div class="rounded-xl p-4 mb-5" style="background: rgba(16,185,129,0.07); border: 1px solid rgba(16,185,129,0.2);">
                <p class="text-xs font-semibold mb-2 text-emerald-400">En Android (Chrome)</p>
                <ol class="space-y-1.5">
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-emerald-400 flex-shrink-0">1.</span>
                    Toca el ícono de <strong style="color: var(--text);">candado</strong> junto a la URL → <strong style="color: var(--text);">Permisos → Ubicación → Permitir</strong>.
                  </li>
                  <li class="text-xs flex items-start gap-2" style="color: var(--text-muted);">
                    <span class="font-bold text-emerald-400 flex-shrink-0">2.</span>
                    Recarga la página y vuelve a intentarlo.
                  </li>
                </ol>
              </div>
            </template>

            <div class="space-y-2">
              <button v-if="locationErrorType === 'denied'" @click="reloadPage()" class="btn-primary btn-lg w-full">
                <ArrowPathRoundedSquareIcon class="w-5 h-5" /> Recargar página
              </button>
              <button v-else-if="locationErrorType !== 'inapp'" @click="retryLocation()" class="btn-primary btn-lg w-full">
                <ArrowPathRoundedSquareIcon class="w-5 h-5" /> Reintentar
              </button>
              <button @click="cancelProcess" class="btn-secondary btn-md w-full">Cerrar</button>
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
      <div v-if="processing && !showLocationErrorModal && !showCameraModal && !showSuccessModal"
        class="fixed inset-0 z-50 flex items-center justify-center"
        style="background: rgba(0,0,0,0.5); backdrop-filter: blur(4px);">
        <div class="glass-card p-6 text-center">
          <div class="w-10 h-10 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin mx-auto mb-3"></div>
          <p class="font-semibold" style="color: var(--text);">Procesando...</p>
          <p class="text-xs mt-1" style="color: var(--text-muted);">Por favor espera</p>
        </div>
      </div>

      <!-- Logout confirmation modal -->
      <Transition name="modal">
        <div v-if="showLogoutModal" class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0,0,0,0.85); backdrop-filter: blur(8px);">
          <div class="w-full max-w-sm glass-card p-6 animate-in text-center" style="background: var(--modal-bg);">
            <div class="w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center bg-rose-500/10 text-rose-500">
              <ArrowRightOnRectangleIcon class="w-8 h-8" />
            </div>
            <h3 class="font-bold text-lg mb-2" style="color: var(--text);">¿Cerrar sesión?</h3>
            <p class="text-sm mb-6" style="color: var(--text-muted);">Tendrás que volver a ingresar tus credenciales para registrar asistencia.</p>
            <div class="flex gap-3">
              <button @click="showLogoutModal = false" class="btn-secondary flex-1">Cancelar</button>
              <button @click="confirmLogout" class="btn-danger flex-1">Salir</button>
            </div>
          </div>
        </div>
      </Transition>

    </Teleport>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import {
  CheckCircleIcon, HandRaisedIcon, MapPinIcon, ArrowRightOnRectangleIcon,
  ExclamationTriangleIcon, XMarkIcon, LockClosedIcon, CloudArrowUpIcon,
  SignalIcon, MinusCircleIcon, CameraIcon, UserCircleIcon, BoltIcon, ArrowPathIcon,
  InformationCircleIcon, DevicePhoneMobileIcon, ArrowPathRoundedSquareIcon
} from '@heroicons/vue/24/outline'
import ThemeToggle from '@/components/ThemeToggle.vue'

const auth = useAuthStore()
const router = useRouter()

const showLogoutModal = ref(false)

async function confirmLogout() {
  showLogoutModal.value = false
  await auth.logout()
  router.push('/login')
}

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
const showLocationErrorModal = ref(false)
const locationErrorType = ref('denied')
const locationErrorCode = ref(0)
const showCameraModal = ref(false)
const showSuccessModal = ref(false)
const registeredAt = ref('')

// Location
let locationWatchId = null
const locationPoints = []

// UI extras
const recoveredPoints = ref(0)
const showHiddenWarning = ref(false)
const isTracking = ref(false)

// Wake lock
let wakeLock = null
const wakeLockActive = ref(false)
async function requestWakeLock() {
  if ('wakeLock' in navigator) {
    try {
      wakeLock = await navigator.wakeLock.request('screen')
      wakeLockActive.value = true
      wakeLock.addEventListener('release', () => { wakeLockActive.value = false })
    } catch { /* not supported */ }
  }
}
function releaseWakeLock() {
  if (wakeLock) { wakeLock.release(); wakeLock = null; wakeLockActive.value = false }
}

// Location buffer — persists points when app goes to background
function saveLocationBuffer(points) {
  try { localStorage.setItem('locationBuffer', JSON.stringify(points)) } catch {}
}
function loadLocationBuffer() {
  try { return JSON.parse(localStorage.getItem('locationBuffer') || '[]') } catch { return [] }
}
function clearLocationBuffer() {
  try { localStorage.removeItem('locationBuffer') } catch {}
}
function handleVisibilityChange() {
  if (document.visibilityState === 'hidden') {
    saveLocationBuffer(locationPoints)
  } else if (document.visibilityState === 'visible' && isTracking.value) {
    showHiddenWarning.value = true
  }
}

function reloadPage() { window.location.reload() }

function isIOS() {
  return /iPhone|iPad|iPod/.test(navigator.userAgent)
}

function isIOSInAppBrowser() {
  const ua = navigator.userAgent
  if (!isIOS()) return false
  // Pure WKWebView (older WhatsApp, Telegram, etc.) — no Version/ token
  if (!/Version\//.test(ua)) return true
  // SFSafariViewController shares the Safari UA but lacks the actual Safari window object.
  // Detect by checking standalone is unavailable AND no window.safari (removed iOS 16 but still works on 15-)
  // Best heuristic: check known in-app UA tokens
  if (/FBAN|FBAV|Instagram|Twitter|LinkedInApp|Snapchat/.test(ua)) return true
  return false
}

function startCheckProcess(type) {
  checkType.value = type
  if (isIOSInAppBrowser()) {
    locationErrorType.value = 'inapp'
    showLocationErrorModal.value = true
    return
  }
  // Call requestLocation directly — iOS Safari requires getCurrentPosition
  // to be called in the same synchronous call stack as the user tap.
  requestLocation()
}

function cancelProcess() {
  showLocationErrorModal.value = false
  processing.value = false
}

// requestLocation MUST be a plain synchronous function so that
// navigator.geolocation.getCurrentPosition() is called as the VERY FIRST
// instruction — before any Vue reactive state changes — to stay within the
// iOS Safari user-gesture call stack. Even a single ref assignment before
// the call can cause WebKit to break the gesture chain.
function requestLocation() {
  if (!navigator.geolocation) {
    locationErrorCode.value = 0
    locationErrorType.value = 'unavailable'
    showLocationErrorModal.value = true
    return
  }
  // ← getCurrentPosition is the absolute first call. No state changes before this.
  navigator.geolocation.getCurrentPosition(
    async (position) => {
      processing.value = false
      locationPoints.push({
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
        accuracy: position.coords.accuracy,
        recorded_at: new Date().toISOString()
      })
      locationWatchId = navigator.geolocation.watchPosition(
        pos => locationPoints.push({
          latitude: pos.coords.latitude,
          longitude: pos.coords.longitude,
          accuracy: pos.coords.accuracy,
          recorded_at: new Date().toISOString()
        }),
        null,
        { enableHighAccuracy: true, maximumAge: 5000 }
      )
      isTracking.value = true
      requestWakeLock()
      cameraStep.value = 'site'
      showCameraModal.value = true
      await nextTick()
      await startCamera('environment')
    },
    (err) => {
      processing.value = false
      const code = err?.code ?? -1
      locationErrorCode.value = code
      if (code === 1) locationErrorType.value = 'denied'           // PERMISSION_DENIED
      else if (code === 2) locationErrorType.value = 'unavailable'  // POSITION_UNAVAILABLE
      else locationErrorType.value = 'timeout'                      // TIMEOUT or unknown
      showLocationErrorModal.value = true
    },
    { enableHighAccuracy: true, timeout: 30000, maximumAge: 0 }
  )
  // Show processing spinner after registering the geolocation listener
  processing.value = true
}

function stopLocationTracking() {
  if (locationWatchId !== null) {
    navigator.geolocation.clearWatch(locationWatchId)
    locationWatchId = null
  }
  isTracking.value = false
  releaseWakeLock()
}

function retryLocation() {
  showLocationErrorModal.value = false
  requestLocation()
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

  // Limit resolution to max 800px to keep payload under 1 MB
  const MAX_DIM = 800
  const scale = Math.min(1, MAX_DIM / Math.max(video.videoWidth, video.videoHeight))
  canvas.width = Math.round(video.videoWidth * scale)
  canvas.height = Math.round(video.videoHeight * scale)

  const ctx = canvas.getContext('2d')
  if (cameraStep.value === 'selfie') { ctx.translate(canvas.width, 0); ctx.scale(-1, 1) }
  ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
  capturedPhoto.value = canvas.toDataURL('image/jpeg', 0.72)
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
    // Send GPS coords so backend can cross-check against IP geolocation
    if (locationPoints.length > 0) {
      formData.append('latitude', String(locationPoints[0].latitude))
      formData.append('longitude', String(locationPoints[0].longitude))
    }
    if (photoSite) formData.append('photo_site', dataURLtoBlob(photoSite), 'site.jpg')
    if (photoSelfie) formData.append('photo_selfie', dataURLtoBlob(photoSelfie), 'selfie.jpg')

    const { data } = await api.post('/checks', formData)
    const recordId = data.record.id

    if (locationPoints.length > 0) {
      try {
        await api.post('/location-points/batch', {
          check_record_id: recordId,
          points: locationPoints
        })
      } catch { /* Points failed — non-critical */ }
    }
    locationPoints.length = 0
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
  } catch (err) {
    const status = err?.response?.status
    const msg = err?.response?.data?.error || ''
    alert(`Error al registrar. Intenta de nuevo.${status ? ` (${status}${msg ? ': ' + msg : ''})` : ''}`)
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
    locationPoints.push(...buffered)
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
