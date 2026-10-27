<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col justify-between p-4 sm:p-6 lg:p-8 select-none font-sans relative overflow-hidden">
    <!-- Ambient Red Glow Accents -->
    <div class="absolute -top-32 -left-32 w-96 h-96 bg-red-600/10 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-32 -right-32 w-96 h-96 bg-red-800/10 rounded-full blur-3xl pointer-events-none"></div>

    <!-- TV Top Bar: AHASS Official Header & Live Clock -->
    <header class="flex items-center justify-between border-b border-slate-800/80 pb-4 mb-6">
      <div class="flex items-center space-x-3">
        <div class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-red-700 to-red-500 flex items-center justify-center font-black text-white text-2xl shadow-lg shadow-red-900/40">
          H
        </div>
        <div>
          <div class="flex items-center space-x-2">
            <h1 class="text-xl sm:text-2xl font-black text-white tracking-tight">AHASS 01234 - NUSANTARA SAKTI</h1>
            <span class="bg-red-600 text-white text-[10px] font-black px-2 py-0.5 rounded tracking-wider uppercase">
              LIVE MONITOR
            </span>
          </div>
          <p class="text-xs text-slate-400 font-medium">Layar Antrean & Progres Servis Motor Honda Ruang Tunggu Konsumen</p>
        </div>
      </div>

      <div class="flex items-center space-x-4">
        <!-- Live Clock -->
        <div class="text-right">
          <p class="text-xl sm:text-2xl font-black text-white font-mono tracking-tight">{{ currentTime }}</p>
          <p class="text-[11px] text-slate-400 font-medium">{{ currentDate }}</p>
        </div>

        <router-link
          to="/dashboard"
          title="Kembali ke Dashboard"
          class="p-2 text-slate-500 hover:text-white bg-slate-900 hover:bg-slate-800 rounded-xl border border-slate-800 transition">
          <Minimize class="w-4 h-4" />
        </router-link>
      </div>
    </header>

    <!-- Main Content Grid -->
    <div class="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      <!-- Left (8 Cols): Real-time Motorcycles in Service -->
      <div class="lg:col-span-8 space-y-4">
        <div class="flex items-center justify-between px-1">
          <div class="flex items-center space-x-2">
            <Wrench class="w-4 h-4 text-red-500" />
            <h2 class="text-xs font-bold text-slate-300 uppercase tracking-wider">Antrean & Progres Pengerjaan Motor</h2>
          </div>
          <span class="text-xs text-slate-400 font-mono">Total Unit: {{ activeOrders.length }}</span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
          <div
            v-for="wo in activeOrders"
            :key="wo.id"
            class="bg-slate-900/90 rounded-2xl p-4 border transition-all duration-300"
            :class="{
              'border-blue-500/50 shadow-lg shadow-blue-950/40': wo.status === 'in_progress',
              'border-purple-500/50 shadow-lg shadow-purple-950/40 animate-pulse': wo.status === 'completed',
              'border-slate-800': wo.status === 'draft',
              'border-emerald-500/50': wo.status === 'paid',
            }">
            <!-- Card Top: Plate Number & Status Badge -->
            <div class="flex items-center justify-between mb-2">
              <span class="font-mono font-black text-base bg-slate-950 text-white px-2.5 py-1 rounded-lg border border-slate-700 tracking-wider">
                {{ wo.customer?.plate_number }}
              </span>
              <span
                :class="{
                  'bg-amber-950 text-amber-300 border-amber-800': wo.status === 'draft',
                  'bg-blue-950 text-blue-300 border-blue-800': wo.status === 'in_progress',
                  'bg-purple-950 text-purple-300 border-purple-800': wo.status === 'completed',
                  'bg-emerald-950 text-emerald-300 border-emerald-800': wo.status === 'paid',
                }"
                class="px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-wider border">
                {{ getDisplayStatus(wo.status) }}
              </span>
            </div>

            <!-- Motorcycle Model & Owner -->
            <div class="space-y-0.5">
              <h3 class="font-bold text-white text-sm truncate">{{ wo.customer?.motor_model }}</h3>
              <p class="text-xs text-slate-400 truncate">Pemilik: {{ wo.customer?.name }}</p>
            </div>

            <!-- Pit & Mechanic Assignment -->
            <div class="mt-3 pt-2.5 border-t border-slate-800/80 flex items-center justify-between text-xs">
              <div class="flex items-center space-x-1.5 text-slate-300">
                <span class="w-5 h-5 rounded-full bg-red-600/80 text-white text-[10px] font-black flex items-center justify-center">
                  {{ wo.mechanic ? wo.mechanic.pit_number : '?' }}
                </span>
                <span class="truncate max-w-[120px]">{{ wo.mechanic ? wo.mechanic.name : 'Antrean Pit' }}</span>
              </div>
              <span class="text-[11px] font-mono font-bold" :class="getProgressColor(wo.status)">
                {{ getProgressPercent(wo.status) }}%
              </span>
            </div>

            <!-- Progress Bar -->
            <div class="mt-2 w-full bg-slate-800 rounded-full h-1.5 overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500"
                :class="getProgressBarClass(wo.status)"
                :style="{ width: `${getProgressPercent(wo.status)}%` }"></div>
            </div>
          </div>

          <div v-if="activeOrders.length === 0" class="col-span-2 bg-slate-900/40 rounded-2xl p-12 text-center border border-slate-800/60 text-slate-500">
            Saat ini tidak ada unit motor dalam antrean pengerjaan.
          </div>
        </div>
      </div>

      <!-- Right (4 Cols): Promo & NSS Service Info -->
      <div class="lg:col-span-4 space-y-4">
        <!-- Promo Banner Card -->
        <div class="bg-gradient-to-br from-red-950/80 via-slate-900 to-slate-950 p-6 rounded-2xl border border-red-900/40 shadow-xl space-y-3">
          <div class="inline-flex items-center space-x-1.5 bg-red-600 text-white text-[10px] font-black px-2.5 py-0.5 rounded-full uppercase">
            <Sparkles class="w-3 h-3" />
            <span>Promo Bulan Ini</span>
          </div>
          <h3 class="text-lg font-black text-white leading-tight">
            Ganti Oli AHM SPX2 Synthetic Hemat s/d 15%
          </h3>
          <p class="text-xs text-slate-300 leading-relaxed">
            Performa mesin lebih dingin dan tarikan enteng. Dapatkan gratis safety check 10 titik untuk pemilik motor Honda Vario 160 & PCX.
          </p>
          <div class="pt-2 text-[11px] text-red-300 font-medium">
            * Berlaku di seluruh jaringan AHASS Nusantara Sakti Group
          </div>
        </div>

        <!-- Service Guarantee Badge -->
        <div class="bg-slate-900/80 p-5 rounded-2xl border border-slate-800/80 space-y-2 text-xs">
          <div class="flex items-center space-x-2 text-emerald-400 font-bold">
            <ShieldCheck class="w-5 h-5" />
            <span>Garansi Servis Resmi AHASS</span>
          </div>
          <p class="text-slate-400 text-[11px] leading-relaxed">
            Setiap servis di AHASS NSS dilindungi garansi 500 KM atau 7 hari kerja. Pastikan menyimpan struk atau faktur resmi digital Anda.
          </p>
        </div>

        <!-- Workshop Daily Stats -->
        <div class="bg-slate-900/80 p-5 rounded-2xl border border-slate-800/80 text-xs flex justify-between items-center">
          <div>
            <span class="text-slate-500 block text-[10px] uppercase font-bold">Pit Aktif Saat Ini</span>
            <span class="text-xl font-black text-white font-mono">4 / 4 PIT</span>
          </div>
          <div class="text-right">
            <span class="text-slate-500 block text-[10px] uppercase font-bold">Kualitas Servis</span>
            <span class="text-emerald-400 font-bold text-xs">100% HGP Asli</span>
          </div>
        </div>
      </div>
    </div>

    <!-- TV Footer Running Text Announcement -->
    <footer class="mt-6 pt-3 border-t border-slate-800/80 flex items-center space-x-3 text-xs text-slate-400 bg-slate-950">
      <span class="bg-red-600 text-white font-black px-2.5 py-0.5 rounded text-[10px] uppercase shrink-0">
        INFO AHASS
      </span>
      <div class="overflow-hidden whitespace-nowrap">
        <p class="inline-block animate-marquee font-medium text-slate-300">
          Selamat datang di Bengkel Resmi Honda AHASS PT Nusantara Sakti Group • Gunakan Selalu Honda Genuine Parts & Oli AHM • Servis Rutin Menjaga Performa Mesin dan Keberlakuan Garansi Resmi Astra Honda Motor • Konsultasikan keluhan motor Anda kepada Service Advisor kami.
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  Wrench,
  Sparkles,
  ShieldCheck,
  Minimize,
} from 'lucide-vue-next'
import { useWorkOrderStore } from '../stores/workOrder'
import type { WorkOrder, WorkOrderStatus } from '../types'

const workOrderStore = useWorkOrderStore()

const currentTime = ref('')
const currentDate = ref('')

function updateTime() {
  const d = new Date()
  currentTime.value = new Intl.DateTimeFormat('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(d) + ' WIB'

  currentDate.value = new Intl.DateTimeFormat('id-ID', {
    weekday: 'long',
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  }).format(d)
}

let timer: any = null
let pollTimer: any = null

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)

  // Poll work orders every 5 seconds to keep the TV screen updated automatically
  workOrderStore.fetchWorkOrders()
  pollTimer = setInterval(() => {
    workOrderStore.fetchWorkOrders()
  }, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
  if (pollTimer) clearInterval(pollTimer)
})

const activeOrders = computed<WorkOrder[]>(() => {
  // Show in-progress, completed, and draft orders on the live TV screen
  return workOrderStore.workOrders.filter(
    (wo) => wo.status === 'in_progress' || wo.status === 'completed' || wo.status === 'draft'
  )
})

function getDisplayStatus(status: WorkOrderStatus): string {
  switch (status) {
    case 'draft':
      return 'Antrean Pendaftaran'
    case 'in_progress':
      return 'Sedang Dikerjakan (Pit)'
    case 'completed':
      return 'Selesai - Siap di Kasir'
    case 'paid':
      return 'Lunas - Siap Pulang'
    default:
      return status
  }
}

function getProgressPercent(status: WorkOrderStatus): number {
  switch (status) {
    case 'draft':
      return 25
    case 'in_progress':
      return 65
    case 'completed':
      return 95
    case 'paid':
      return 100
    default:
      return 0
  }
}

function getProgressColor(status: WorkOrderStatus): string {
  switch (status) {
    case 'draft':
      return 'text-amber-400'
    case 'in_progress':
      return 'text-blue-400'
    case 'completed':
      return 'text-purple-400'
    case 'paid':
      return 'text-emerald-400'
    default:
      return 'text-slate-400'
  }
}

function getProgressBarClass(status: WorkOrderStatus): string {
  switch (status) {
    case 'draft':
      return 'bg-amber-500'
    case 'in_progress':
      return 'bg-blue-500'
    case 'completed':
      return 'bg-purple-500'
    case 'paid':
      return 'bg-emerald-500'
    default:
      return 'bg-slate-600'
  }
}
</script>

<style scoped>
@keyframes marquee {
  0% { transform: translateX(100%); }
  100% { transform: translateX(-100%); }
}
.animate-marquee {
  display: inline-block;
  white-space: nowrap;
  animation: marquee 30s linear infinite;
}
</style>

