<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-slate-900 tracking-tight">Kinerja Mekanik & Pit Bengkel AHASS</h2>
        <p class="text-xs text-slate-500">Monitoring status Pit servis motor, spesialisasi teknisi Honda, dan ketersediaan</p>
      </div>
    </div>

    <!-- Pit Grid Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div 
        v-for="m in mechanics" 
        :key="m.id"
        class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs relative overflow-hidden flex flex-col justify-between">
        
        <!-- Pit Number Badge -->
        <div>
          <div class="flex items-center justify-between mb-4">
            <span class="inline-flex items-center space-x-1.5 bg-red-100 text-red-700 font-extrabold text-xs px-2.5 py-1 rounded-lg">
              <Wrench class="w-3.5 h-3.5" />
              <span>PIT {{ m.pit_number }}</span>
            </span>
            <span 
              :class="m.is_available ? 'bg-emerald-100 text-emerald-800' : 'bg-blue-100 text-blue-800'"
              class="text-[10px] font-bold px-2 py-0.5 rounded uppercase">
              {{ m.is_available ? 'Siap Pengerjaan' : 'Sedang Servis' }}
            </span>
          </div>

          <h3 class="font-bold text-slate-900 text-base">{{ m.name }}</h3>
          <p class="text-xs text-slate-500 font-mono mt-0.5">{{ m.code }}</p>

          <div class="mt-4 pt-4 border-t border-slate-100 space-y-1.5 text-xs text-slate-600">
            <p><span class="text-slate-400">Spesialisasi:</span></p>
            <p class="font-semibold text-slate-800 bg-slate-50 p-2 rounded-lg border border-slate-200/60">
              {{ m.specialization }}
            </p>
            <p class="pt-2"><span class="text-slate-400">Telepon:</span> <span class="font-mono">{{ m.phone || '-' }}</span></p>
          </div>
        </div>

        <div class="mt-6 pt-3 border-t border-slate-100 flex items-center justify-between text-xs">
          <span class="text-slate-400 text-[11px]">Sertifikasi Astra Honda</span>
          <span class="text-emerald-600 font-bold text-[11px]">Level 2 / 3</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Wrench } from 'lucide-vue-next'
import api from '../services/api'
import type { Mechanic } from '../types'

const mechanics = ref<Mechanic[]>([])

onMounted(async () => {
  try {
    const res = await api.get('/mechanics')
    mechanics.value = res.data.data
  } catch (err) {
    console.error('Failed to load mechanics', err)
  }
})
</script>

