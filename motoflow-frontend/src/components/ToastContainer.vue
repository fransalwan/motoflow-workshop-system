<template>
  <div class="fixed bottom-5 right-5 z-50 flex flex-col space-y-2 pointer-events-none max-w-sm w-full px-4">
    <transition-group
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0">
      <div
        v-for="toast in toastStore.toasts"
        :key="toast.id"
        class="pointer-events-auto flex items-center justify-between p-4 rounded-xl shadow-xl border text-xs font-semibold text-white transition-all backdrop-blur-md"
        :class="{
          'bg-slate-900/95 border-emerald-500/50 text-emerald-300': toast.type === 'success',
          'bg-slate-900/95 border-red-500/50 text-red-300': toast.type === 'error',
          'bg-slate-900/95 border-amber-500/50 text-amber-300': toast.type === 'warning',
          'bg-slate-900/95 border-blue-500/50 text-blue-300': toast.type === 'info',
        }">
        <div class="flex items-center space-x-2.5">
          <CheckCircle2 v-if="toast.type === 'success'" class="w-4 h-4 text-emerald-400 shrink-0" />
          <AlertCircle v-else-if="toast.type === 'error'" class="w-4 h-4 text-red-400 shrink-0" />
          <AlertTriangle v-else-if="toast.type === 'warning'" class="w-4 h-4 text-amber-400 shrink-0" />
          <Info v-else class="w-4 h-4 text-blue-400 shrink-0" />
          <span class="text-slate-100 leading-snug">{{ toast.message }}</span>
        </div>
        <button
          @click="toastStore.remove(toast.id)"
          class="ml-3 text-slate-400 hover:text-white transition">
          <X class="w-3.5 h-3.5" />
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { CheckCircle2, AlertCircle, AlertTriangle, Info, X } from 'lucide-vue-next'
import { useToastStore } from '../stores/toast'

const toastStore = useToastStore()
</script>

