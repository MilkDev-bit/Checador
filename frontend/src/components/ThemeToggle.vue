<template>
  <button
    @click="theme.toggle()"
    class="w-9 h-9 rounded-xl flex items-center justify-center transition-all"
    :class="theme.isDark
      ? 'text-slate-400 hover:text-amber-300 hover:bg-white/8'
      : 'text-slate-500 hover:text-amber-500 hover:bg-black/8'"
    :title="theme.isDark ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'"
  >
    <SunIcon v-if="theme.isDark" class="w-5 h-5" />
    <MoonIcon v-else class="w-5 h-5" />
  </button>
</template>

<script setup>
import { SunIcon, MoonIcon } from '@heroicons/vue/24/outline'
import { useThemeStore } from '@/stores/theme'

const theme = useThemeStore()
</script>
