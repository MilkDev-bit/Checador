<template>
  <Teleport to="body">
    <!-- Overlay backdrop — fade in/out -->
    <Transition name="izi-overlay">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-[999]"
        style="background: rgba(0,0,0,0.8);"
        @click="handleBackdrop"
      />
    </Transition>

    <!-- Centering wrapper + panel — iziModal comingIn/comingOut keyframe animation -->
    <Transition @enter="onEnter" @leave="onLeave" :css="false">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-[1000] flex items-center justify-center px-4 pointer-events-none"
      >
        <div
          ref="panelEl"
          class="flex flex-col overflow-hidden shadow-2xl pointer-events-auto w-full"
          :style="panelStyle"
          role="dialog"
          aria-modal="true"
        >
          <!-- Header — iziModal colored bar style -->
          <div
            v-if="title || subtitle || closable"
            class="flex items-center justify-between flex-shrink-0"
            style="padding: 14px 18px; box-shadow: inset 0 -10px 15px -12px rgba(0,0,0,0.3);"
            :style="{ background: headerColor }"
          >
            <div class="flex-1 min-w-0">
              <span
                v-if="title"
                class="block font-semibold truncate"
                style="color: #fff; font-size: 18px; line-height: 1.3; font-family: 'Lato', Arial, sans-serif;"
              >{{ title }}</span>
              <p
                v-if="subtitle"
                class="block truncate mt-0.5"
                style="color: rgba(255,255,255,0.6); font-size: 12px; line-height: 1.45; font-family: 'Lato', Arial, sans-serif;"
              >{{ subtitle }}</p>
            </div>
            <button
              v-if="closable"
              class="flex-shrink-0 ml-3 w-8 h-8 flex items-center justify-center rounded transition-opacity"
              style="opacity: 0.7;"
              aria-label="Cerrar"
              @click="close"
              @mouseenter="$event.target.style.opacity = '1'"
              @mouseleave="$event.target.style.opacity = '0.7'"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="white" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Body (scrollable) -->
          <div class="flex-1 overflow-y-auto custom-scroll" style="overscroll-behavior: contain;">
            <div class="px-5 py-4">
              <slot />
            </div>
          </div>

          <!-- Footer -->
          <div
            v-if="$slots.footer"
            class="px-5 pb-5 pt-3 flex gap-3 flex-shrink-0 border-t"
            style="border-color: var(--border-subtle);"
          >
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue:      { type: Boolean, default: false },
  title:           { type: String,  default: '' },
  subtitle:        { type: String,  default: '' },
  color:           { type: String,  default: 'default' },
  size:            { type: String,  default: 'sm' },
  closable:        { type: Boolean, default: true },
  closeOnBackdrop: { type: Boolean, default: true },
})

const emit = defineEmits(['update:modelValue', 'close'])
const panelEl = ref(null)

const headerColors = {
  default: '#64748b',
  brand:   '#6366f1',
  success: '#10b981',
  danger:  '#ef4444',
  warning: '#f59e0b',
  info:    '#3b82f6',
}

const widthMap = { sm: 400, md: 520, lg: 640, xl: 780 }

const headerColor = computed(() => headerColors[props.color] ?? headerColors.default)

const panelStyle = computed(() => ({
  maxWidth: (widthMap[props.size] ?? 400) + 'px',
  maxHeight: 'calc(100vh - 4rem)',
  borderRadius: '16px',
  background: 'var(--modal-bg)',
}))

function close() {
  emit('update:modelValue', false)
  emit('close')
}

function handleBackdrop() {
  if (props.closeOnBackdrop) close()
}

function handleKeydown(e) {
  if (e.key === 'Escape' && props.modelValue && props.closable) close()
}

// iziModal CSS keyframe animations — loaded via izimodal/css/iziModal.min.css in main.js
function onEnter(el, done) {
  const panel = panelEl.value
  if (!panel) return done()
  panel.style.animation = 'iziM-comingIn 0.5s ease forwards'
  panel.addEventListener('animationend', done, { once: true })
}

function onLeave(el, done) {
  const panel = panelEl.value
  if (!panel) return done()
  panel.style.animation = 'iziM-comingOut 0.5s cubic-bezier(.16,.81,.32,1) forwards'
  panel.addEventListener('animationend', done, { once: true })
}

onMounted  (() => document.addEventListener('keydown', handleKeydown))
onUnmounted(() => document.removeEventListener('keydown', handleKeydown))
</script>
