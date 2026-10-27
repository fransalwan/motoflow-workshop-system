<template>
  <header class="h-20 bg-white border-b border-slate-200/80 px-8 flex items-center justify-between shadow-xs">
    <div>
      <h2 class="text-xl font-bold text-slate-800 tracking-tight">{{ pageTitle }}</h2>
      <p class="text-xs text-slate-500 font-medium">Sistem Informasi Manajemen Bengkel Resmi Honda AHASS NSS</p>
    </div>

    <div class="flex items-center space-x-4">
      <!-- Live Workshop Clock -->
      <div class="hidden sm:flex items-center space-x-2 bg-slate-50 px-3.5 py-1.5 rounded-xl border border-slate-200 text-xs text-slate-600 font-mono">
        <Clock class="w-3.5 h-3.5 text-red-600" />
        <span>{{ currentTime }}</span>
      </div>

      <!-- Quick Action: Buat Work Order Baru -->
      <router-link
        to="/work-orders/create"
        class="inline-flex items-center space-x-2 bg-red-600 hover:bg-red-700 text-white font-semibold text-xs px-4 py-2.5 rounded-xl shadow-sm transition active:scale-95">
        <Plus class="w-4 h-4" />
        <span>Servis Baru</span>
      </router-link>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { Clock, Plus } from 'lucide-vue-next'

const route = useRoute()
const currentTime = ref('')

const pageTitle = computed(() => {
  switch (route.name) {
    case 'Dashboard':
      return 'Executive Workshop Dashboard'
    case 'WorkOrders':
      return 'Daftar Perintah Kerja & Kasir POS'
    case 'WorkOrderCreate':
      return 'Pendaftaran Servis & Estimasi Biaya'
    case 'WorkOrderDetail':
      return 'Kasir Billing & Detail Faktur Servis'
    case 'Inventory':
      return 'Inventaris Honda Genuine Parts (HGP)'
    case 'Customers':
      return 'Master Pelanggan & Unit Kendaraan'
    case 'Mechanics':
      return 'Kinerja Mekanik & Pit Bengkel'
    default:
      return 'MotoFlow AHASS'
  }
})

function updateTime() {
  const d = new Date()
  currentTime.value = new Intl.DateTimeFormat('id-ID', {
    weekday: 'short',
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(d) + ' WIB'
}

let timer: any = null
onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

