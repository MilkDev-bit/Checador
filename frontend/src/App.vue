<template>
  <router-view />

  <!-- Fix 4: SW update banner — shown when a new version takes control -->
  <Transition name="slide-up">
    <div v-if="showUpdateBanner"
      class="fixed bottom-4 left-4 right-4 z-50 glass-card p-4 border border-brand-500/30 shadow-2xl flex items-center gap-4">
      <div class="bg-brand-50 dark:bg-brand-500/10 text-brand-600 dark:text-brand-400 p-2 rounded-xl flex-shrink-0">
        <ArrowPathRoundedSquareIcon class="w-6 h-6" />
      </div>
      <div class="flex-1">
        <h4 class="font-bold text-slate-900 dark:text-white text-sm">Nueva versión disponible</h4>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Recarga para aplicar la actualización.</p>
      </div>
      <button @click="reloadForUpdate"
        class="text-xs font-bold px-3 py-1.5 rounded-lg bg-brand-500 text-white hover:bg-brand-600 transition-colors flex-shrink-0">
        Actualizar
      </button>
    </div>
  </Transition>
  <Transition name="slide-up">
    <div v-if="showAndroidPrompt" class="fixed bottom-4 left-4 right-4 z-50 glass-card p-4 border border-brand-500/20 shadow-2xl flex items-start gap-4">
      <div class="bg-brand-50 dark:bg-brand-500/10 text-brand-600 dark:text-brand-400 p-2 rounded-xl flex-shrink-0">
        <ArrowDownOnSquareIcon class="w-6 h-6" />
      </div>
      <div class="flex-1 mt-0.5">
        <h4 class="font-bold text-slate-900 dark:text-white text-sm">Instalar PaseLista</h4>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 leading-relaxed">
          Instala la app para acceso rápido desde tu pantalla de inicio.
        </p>
      </div>
      <div class="flex items-center gap-2 flex-shrink-0">
        <button @click="installPwa"
          class="text-xs font-bold px-3 py-1.5 rounded-lg bg-brand-500 text-white hover:bg-brand-600 transition-colors">
          Instalar
        </button>
        <button @click="dismissAndroidPrompt" class="text-slate-400 hover:text-slate-600 dark:hover:text-white p-1">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
    </div>
  </Transition>

  <!-- iOS PWA Install Prompt -->
  <Transition name="slide-up">
    <div v-if="showIosPrompt" class="fixed bottom-4 left-4 right-4 z-50 glass-card p-4 border border-brand-500/20 shadow-2xl flex items-start gap-4">
      <div class="bg-brand-50 dark:bg-brand-500/10 text-brand-600 dark:text-brand-400 p-2 rounded-xl flex-shrink-0">
        <ArrowUpOnSquareIcon class="w-6 h-6" />
      </div>
      <div class="flex-1 mt-0.5">
        <h4 class="font-bold text-slate-900 dark:text-white text-sm">Instalar PaseLista</h4>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 leading-relaxed">
          Toca <strong>Compartir</strong> en la barra inferior y selecciona <br/><strong>"Añadir a la pantalla de inicio"</strong>.
        </p>
      </div>
      <button @click="dismissIosPrompt" class="text-slate-400 hover:text-slate-600 dark:hover:text-white p-1">
        <XMarkIcon class="w-5 h-5" />
      </button>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { ArrowUpOnSquareIcon, ArrowDownOnSquareIcon, XMarkIcon, ArrowPathRoundedSquareIcon } from '@heroicons/vue/24/outline'

useThemeStore()

// Fix 4: SW update — when a new service worker takes control, prompt reload
const showUpdateBanner = ref(false)
function reloadForUpdate() { window.location.reload() }
function handleControllerChange() { showUpdateBanner.value = true }

// ── Android / Chrome install prompt ──────────────────────────────────────────
const showAndroidPrompt = ref(false)
let deferredPrompt = null

function handleBeforeInstallPrompt(e) {
  e.preventDefault()
  deferredPrompt = e
  const dismissed = localStorage.getItem('androidPwaPromptDismissed')
  // Don't show if already installed or dismissed in last 30 days
  if (!dismissed || Date.now() - Number(dismissed) > 30 * 24 * 60 * 60 * 1000) {
    setTimeout(() => { showAndroidPrompt.value = true }, 2500)
  }
}

async function installPwa() {
  showAndroidPrompt.value = false
  if (!deferredPrompt) return
  deferredPrompt.prompt()
  await deferredPrompt.userChoice
  deferredPrompt = null
}

function dismissAndroidPrompt() {
  showAndroidPrompt.value = false
  localStorage.setItem('androidPwaPromptDismissed', String(Date.now()))
}

// ── iOS install prompt ────────────────────────────────────────────────────────
const showIosPrompt = ref(false)

const dismissIosPrompt = () => {
  showIosPrompt.value = false
  localStorage.setItem('iosPwaPromptDismissed', 'true')
}

onMounted(() => {
  const ua = window.navigator.userAgent.toLowerCase()
  const isIos = /iphone|ipad|ipod/.test(ua)
  const isAndroidChrome = /android/.test(ua) && /chrome/.test(ua)
  const isStandalone = window.matchMedia('(display-mode: standalone)').matches
    || ('standalone' in window.navigator && window.navigator.standalone)

  if (isStandalone) return // already installed

  // Fix 4: listen for SW controller change (new version took over)
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.addEventListener('controllerchange', handleControllerChange)
  }

  if (isAndroidChrome) {
    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  } else if (isIos) {
    const hasDismissed = localStorage.getItem('iosPwaPromptDismissed')
    if (!hasDismissed) {
      setTimeout(() => { showIosPrompt.value = true }, 2500)
    }
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.removeEventListener('controllerchange', handleControllerChange)
  }
})
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(120%);
}
</style>
