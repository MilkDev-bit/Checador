<template>
  <router-view />

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
import { ref, onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { ArrowUpOnSquareIcon, XMarkIcon } from '@heroicons/vue/24/outline'

useThemeStore()

const showIosPrompt = ref(false)

const dismissIosPrompt = () => {
  showIosPrompt.value = false
  localStorage.setItem('iosPwaPromptDismissed', 'true')
}

onMounted(() => {
  // Check if device is iOS
  const isIos = () => {
    const userAgent = window.navigator.userAgent.toLowerCase()
    return /iphone|ipad|ipod/.test(userAgent)
  }
  
  // Check if app is open from Home Screen (standalone)
  const isStandalone = () => {
    return ('standalone' in window.navigator) && window.navigator.standalone
  }

  if (isIos() && !isStandalone()) {
    const hasDismissed = localStorage.getItem('iosPwaPromptDismissed')
    if (!hasDismissed) {
      setTimeout(() => {
        showIosPrompt.value = true
      }, 2500)
    }
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
