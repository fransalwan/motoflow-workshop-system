<template>
  <header class="h-16 lg:h-20 bg-white border-b border-slate-200/80 px-4 sm:px-8 flex items-center justify-between shadow-xs sticky top-0 z-30">
    <div class="flex items-center space-x-3">
      <!-- Mobile Hamburger Menu Button -->
      <button 
        @click="$emit('toggle-mobile-menu')"
        class="lg:hidden p-2 text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded-xl transition">
        <Menu class="w-5 h-5" />
      </button>

      <div>
        <h2 class="text-base sm:text-xl font-bold text-slate-800 tracking-tight truncate max-w-[200px] sm:max-w-none">
          {{ pageTitle }}
        </h2>
        <p class="hidden sm:block text-xs text-slate-500 font-medium">Sistem Informasi Manajemen Bengkel Resmi Honda AHASS NSS</p>
      </div>
    </div>

    <div class="flex items-center space-x-2 sm:space-x-4">
      <!-- Live Workshop Clock (Desktop only) -->
      <div class="hidden md:flex items-center space-x-2 bg-slate-50 px-3.5 py-1.5 rounded-xl border border-slate-200 text-xs text-slate-600 font-mono">
        <Clock class="w-3.5 h-3.5 text-red-600" />
        <span>{{ currentTime }}</span>
      </div>

      <!-- Quick Action: TV Monitor Ruang Tunggu -->
      <router-link
        to="/waiting-display"
        title="Layar TV Ruang Tunggu Konsumen"
        class="hidden sm:inline-flex items-center space-x-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs px-3 py-2 rounded-xl border border-slate-200 transition">
        <Tv class="w-4 h-4 text-red-600" />
        <span class="hidden md:inline">TV Ruang Tunggu</span>
      </router-link>

      <!-- Quick Action: Buat Work Order Baru -->
      <router-link
        to="/work-orders/create"
        class="inline-flex items-center space-x-1.5 bg-red-600 hover:bg-red-700 text-white font-semibold text-xs px-3 sm:px-4 py-2 sm:py-2.5 rounded-xl shadow-xs transition active:scale-95">
        <Plus class="w-4 h-4" />
        <span>Servis Baru</span>
      </router-link>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { Clock, Plus, Menu, Tv } from 'lucide-vue-next'

defineEmits(['toggle-mobile-menu'])

const route = useRoute()
const currentTime = ref('')

const pageTitle = computed(() => {
  switch (route.name) {
    case 'Dashboard':
      return 'Executive Workshop Dashboard'
    case 'WorkOrders':
      return 'Perintah Kerja & Kasir POS'
    case 'WorkOrderCreate':
      return 'Pendaftaran Servis Baru'
    case 'WorkOrderDetail':
      return 'Kasir Billing & Detail Servis'
    case 'Inventory':
      return 'Inventaris Suku Cadang HGP'
    case 'Customers':
      return 'Master Pelanggan & Motor'
    case 'Mechanics':
      return 'Kinerja Mekanik & Pit Bengkel'
    case 'WaitingDisplay':
      return 'Monitor Ruang Tunggu Konsumen'
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
