<template>
  <Teleport to="body">
    <!-- Backdrop -->
    <Transition name="izi-overlay">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex justify-center"
        :class="bottomSheet ? 'items-end' : 'items-center px-4'"
        style="background: rgba(0,0,0,0.8); backdrop-filter: blur(8px); -webkit-backdrop-filter: blur(8px);"
        @click.self="handleBackdrop"
      >
        <!-- Panel -->
        <Transition :name="bottomSheet ? 'izi-sheet' : 'izi-modal'">
          <div
            v-if="modelValue"
            class="glass-card w-full overflow-hidden"
            :class="[
              sizeClass,
              bottomSheet && 'rounded-t-3xl rounded-b-none sm:rounded-3xl mb-0 sm:mb-4',
            ]"
            style="background: var(--modal-bg);"
            role="dialog"
            aria-modal="true"
          >
            <!-- Header: shown only when there is title, icon, or closable button -->
            <div
              v-if="title || icon || closable"
              class="flex items-center justify-between px-5 py-4 flex-shrink-0"
              :style="headerBorderStyle"
            >
              <div class="flex items-center gap-3 min-w-0 flex-1">
                <!-- Icon badge -->
                <div
                  v-if="icon"
                  class="w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0"
                  :class="iconBgClass"
                >
                  <component :is="icon" class="w-5 h-5" :class="iconColorClass" />
                </div>
                <!-- Title / subtitle -->
                <div v-if="title || subtitle" class="min-w-0 flex-1">
                  <h3 class="font-bold text-base leading-tight" style="color: var(--text);">{{ title }}</h3>
                  <p v-if="subtitle" class="text-xs mt-0.5 truncate" style="color: var(--text-muted);">{{ subtitle }}</p>
                </div>
              </div>
              <!-- Close button -->
              <button
                v-if="closable"
                class="w-8 h-8 rounded-lg flex items-center justify-center transition-colors ml-3 flex-shrink-0 hover:bg-black/10 dark:hover:bg-white/10"
                style="color: var(--text-muted);"
                aria-label="Cerrar"
                @click="close"
              >
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>

            <!-- Body slot -->
            <div class="px-5 py-4">
              <slot />
            </div>

            <!-- Footer slot -->
            <div v-if="$slots.footer" class="px-5 pb-5 pt-0 flex gap-3">
              <slot name="footer" />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  /** v-model: controls visibility */
  modelValue:      { type: Boolean, default: false },
  /** Header title text (omit to hide header) */
  title:           { type: String,  default: '' },
  /** Small subtitle below title */
  subtitle:        { type: String,  default: '' },
  /** Heroicon component passed as prop, e.g. :icon="CheckCircleIcon" */
  icon:            { type: Object,  default: null },
  /**
   * Accent color for header border + icon badge.
   * 'default' | 'brand' | 'success' | 'danger' | 'warning' | 'info'
   */
  color:           { type: String,  default: 'default' },
  /** Panel width: 'sm' | 'md' | 'lg' | 'xl' */
  size:            { type: String,  default: 'sm' },
  /** Show the × close button in the header */
  closable:        { type: Boolean, default: true },
  /** Close when clicking the dark backdrop */
  closeOnBackdrop: { type: Boolean, default: true },
  /**
   * On mobile: slide up from bottom (sheet style).
   * On desktop it centres regardless.
   */
  bottomSheet:     { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'close'])

// ─── Computed ─────────────────────────────────────────────────────────────────
const sizeClass = computed(() => ({
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-2xl',
}[props.size] ?? 'max-w-sm'))

const headerBorderStyle = computed(() => {
  const borders = {
    default: 'var(--border-subtle)',
    brand:   'rgba(99,102,241,0.3)',
    success: 'rgba(16,185,129,0.3)',
    danger:  'rgba(239,68,68,0.3)',
    warning: 'rgba(245,158,11,0.3)',
    info:    'rgba(59,130,246,0.3)',
  }
  return `border-bottom: 1px solid ${borders[props.color] ?? borders.default}`
})

const iconBgClass = computed(() => ({
  default: 'bg-slate-100 dark:bg-white/10',
  brand:   'bg-indigo-100 dark:bg-indigo-500/15',
  success: 'bg-emerald-100 dark:bg-emerald-500/15',
  danger:  'bg-rose-100 dark:bg-rose-500/15',
  warning: 'bg-amber-100 dark:bg-amber-500/15',
  info:    'bg-blue-100 dark:bg-blue-500/15',
}[props.color] ?? 'bg-slate-100 dark:bg-white/10'))

const iconColorClass = computed(() => ({
  default: 'text-slate-500 dark:text-slate-400',
  brand:   'text-indigo-600 dark:text-indigo-400',
  success: 'text-emerald-600 dark:text-emerald-400',
  danger:  'text-rose-500 dark:text-rose-400',
  warning: 'text-amber-600 dark:text-amber-400',
  info:    'text-blue-600 dark:text-blue-400',
}[props.color] ?? 'text-slate-500 dark:text-slate-400'))

// ─── Actions ──────────────────────────────────────────────────────────────────
function close () {
  emit('update:modelValue', false)
  emit('close')
}

function handleBackdrop () {
  if (props.closeOnBackdrop) close()
}

function handleKeydown (e) {
  if (e.key === 'Escape' && props.modelValue && props.closable) close()
}

onMounted  (() => document.addEventListener('keydown', handleKeydown))
onUnmounted(() => document.removeEventListener('keydown', handleKeydown))
</script>
