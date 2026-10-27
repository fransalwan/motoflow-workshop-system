<template>
  <div class="space-y-6">
    <!-- Controls Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <!-- Filter Tabs / Mode Switcher -->
      <div class="flex flex-wrap items-center gap-2">
        <!-- View Mode: Table vs Kanban -->
        <div class="bg-slate-200/80 p-1 rounded-xl flex items-center space-x-1 text-xs font-semibold">
          <button
            @click="viewMode = 'table'"
            :class="viewMode === 'table' ? 'bg-white text-slate-900 shadow-xs' : 'text-slate-600 hover:text-slate-900'"
            class="px-3 py-1.5 rounded-lg transition flex items-center space-x-1.5">
            <Table class="w-3.5 h-3.5" />
            <span>Tabel</span>
          </button>
          <button
            @click="viewMode = 'kanban'"
            :class="viewMode === 'kanban' ? 'bg-white text-slate-900 shadow-xs' : 'text-slate-600 hover:text-slate-900'"
            class="px-3 py-1.5 rounded-lg transition flex items-center space-x-1.5">
            <Kanban class="w-3.5 h-3.5" />
            <span>Papan Pit (Kanban)</span>
          </button>
        </div>

        <!-- Filter Status Tabs (When in Table View) -->
        <div v-if="viewMode === 'table'" class="flex items-center space-x-1 bg-slate-200/80 p-1 rounded-xl text-xs font-semibold overflow-x-auto max-w-full">
          <button
            v-for="tab in filterTabs"
            :key="tab.value"
            @click="selectedTab = tab.value"
            :class="selectedTab === tab.value ? 'bg-white text-slate-900 shadow-xs' : 'text-slate-600 hover:text-slate-900'"
            class="px-3 py-1.5 rounded-lg transition whitespace-nowrap">
            {{ tab.label }}
          </button>
        </div>
      </div>

      <!-- Search & New Button -->
      <div class="flex items-center space-x-3 w-full md:w-auto">
        <div class="relative flex-1 md:w-64">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari WO, Plat, atau Nama..."
            class="w-full bg-white border border-slate-200 rounded-xl pl-9 pr-4 py-2 text-xs text-slate-900 placeholder-slate-400 focus:outline-none focus:border-red-500 shadow-xs" />
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
        </div>

        <router-link
          to="/work-orders/create"
          class="inline-flex items-center space-x-1.5 bg-red-600 hover:bg-red-700 text-white font-semibold text-xs px-4 py-2 rounded-xl shadow-xs transition shrink-0">
          <Plus class="w-4 h-4" />
          <span>Buat WO</span>
        </router-link>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="workOrderStore.loading && workOrderStore.workOrders.length === 0">
      <SkeletonLoader />
    </div>

    <!-- VIEW 1: KANBAN WORKSHOP PIT BOARD -->
    <div v-else-if="viewMode === 'kanban'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-start">
      <!-- Column 1: Draft / Pendaftaran SA -->
      <div class="bg-slate-100/80 p-3.5 rounded-2xl border border-slate-200/80 space-y-3">
        <div class="flex items-center justify-between px-1">
          <div class="flex items-center space-x-2">
            <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span>
            <h4 class="font-bold text-xs text-slate-800 uppercase tracking-wider">1. Pendaftaran SA</h4>
          </div>
          <span class="bg-amber-100 text-amber-800 text-[10px] font-black px-2 py-0.5 rounded-full font-mono">
            {{ getOrdersByStatus('draft').length }}
          </span>
        </div>

        <div class="space-y-3 min-h-[300px]">
          <div 
            v-for="wo in getOrdersByStatus('draft')" 
            :key="wo.id" 
            class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-xs hover:border-red-400 transition space-y-2.5 group">
            <div class="flex items-center justify-between">
              <span class="font-mono font-bold text-xs bg-slate-100 px-2 py-0.5 rounded border border-slate-200 text-slate-900">
                {{ wo.customer?.plate_number }}
              </span>
              <span class="text-[10px] font-mono text-slate-400">{{ formatDate(wo.created_at) }}</span>
            </div>

            <div>
              <h5 class="font-bold text-slate-800 text-xs truncate">{{ wo.customer?.motor_model }}</h5>
              <p class="text-[11px] text-slate-500 truncate">{{ wo.customer?.name }}</p>
            </div>

            <div class="pt-2 border-t border-slate-100 flex items-center justify-between">
              <span class="font-mono font-bold text-xs text-slate-900">{{ formatRupiah(wo.grand_total) }}</span>
              <button 
                @click="advanceStatus(wo.id, 'in_progress')"
                class="bg-blue-600 hover:bg-blue-700 text-white text-[10px] font-bold px-2.5 py-1 rounded-lg transition flex items-center space-x-1">
                <span>Ke Pit</span>
                <ArrowRight class="w-3 h-3" />
              </button>
            </div>
          </div>

          <div v-if="getOrdersByStatus('draft').length === 0" class="text-center py-10 text-slate-400 text-xs">
            Antrean kosong
          </div>
        </div>
      </div>

      <!-- Column 2: In Progress (Pit Servis) -->
      <div class="bg-blue-50/50 p-3.5 rounded-2xl border border-blue-200/60 space-y-3">
        <div class="flex items-center justify-between px-1">
          <div class="flex items-center space-x-2">
            <span class="w-2.5 h-2.5 rounded-full bg-blue-600 animate-ping"></span>
            <h4 class="font-bold text-xs text-blue-900 uppercase tracking-wider">2. Dikerjakan (Pit)</h4>
          </div>
          <span class="bg-blue-100 text-blue-800 text-[10px] font-black px-2 py-0.5 rounded-full font-mono">
            {{ getOrdersByStatus('in_progress').length }}
          </span>
        </div>

        <div class="space-y-3 min-h-[300px]">
          <div 
            v-for="wo in getOrdersByStatus('in_progress')" 
            :key="wo.id" 
            class="bg-white p-4 rounded-xl border border-blue-200 shadow-xs hover:border-blue-400 transition space-y-2.5">
            <div class="flex items-center justify-between">
              <span class="font-mono font-black text-xs bg-blue-100 text-blue-900 px-2 py-0.5 rounded border border-blue-300">
                {{ wo.customer?.plate_number }}
              </span>
              <span v-if="wo.mechanic" class="text-[10px] font-bold bg-red-100 text-red-700 px-1.5 py-0.5 rounded">
                PIT {{ wo.mechanic.pit_number }}
              </span>
            </div>

            <div>
              <h5 class="font-bold text-slate-800 text-xs">{{ wo.customer?.motor_model }}</h5>
              <p class="text-[11px] text-slate-500">{{ wo.mechanic?.name || 'Mekanik Servis' }}</p>
            </div>

            <div class="pt-2 border-t border-slate-100 flex items-center justify-between">
              <span class="font-mono font-bold text-xs text-slate-900">{{ formatRupiah(wo.grand_total) }}</span>
              <button 
                @click="advanceStatus(wo.id, 'completed')"
                class="bg-purple-600 hover:bg-purple-700 text-white text-[10px] font-bold px-2.5 py-1 rounded-lg transition flex items-center space-x-1">
                <span>Selesai QC</span>
                <Check class="w-3 h-3" />
              </button>
            </div>
          </div>

          <div v-if="getOrdersByStatus('in_progress').length === 0" class="text-center py-10 text-slate-400 text-xs">
            Pit sedang kosong
          </div>
        </div>
      </div>

      <!-- Column 3: Selesai / Menunggu Kasir -->
      <div class="bg-purple-50/50 p-3.5 rounded-2xl border border-purple-200/60 space-y-3">
        <div class="flex items-center justify-between px-1">
          <div class="flex items-center space-x-2">
            <span class="w-2.5 h-2.5 rounded-full bg-purple-600"></span>
            <h4 class="font-bold text-xs text-purple-900 uppercase tracking-wider">3. Selesai (Kasir)</h4>
          </div>
          <span class="bg-purple-100 text-purple-800 text-[10px] font-black px-2 py-0.5 rounded-full font-mono">
            {{ getOrdersByStatus('completed').length }}
          </span>
        </div>

        <div class="space-y-3 min-h-[300px]">
          <div 
            v-for="wo in getOrdersByStatus('completed')" 
            :key="wo.id" 
            class="bg-white p-4 rounded-xl border border-purple-200 shadow-xs hover:border-purple-400 transition space-y-2.5">
            <div class="flex items-center justify-between">
              <span class="font-mono font-bold text-xs bg-slate-100 text-slate-900 px-2 py-0.5 rounded border border-slate-200">
                {{ wo.customer?.plate_number }}
              </span>
              <span class="text-[10px] font-bold text-purple-700 bg-purple-50 px-1.5 py-0.5 rounded">
                Lolos QC
              </span>
            </div>

            <div>
              <h5 class="font-bold text-slate-800 text-xs">{{ wo.customer?.motor_model }}</h5>
              <p class="text-[11px] text-slate-500">{{ wo.customer?.name }}</p>
            </div>

            <div class="pt-2 border-t border-slate-100 flex items-center justify-between">
              <span class="font-mono font-bold text-xs text-purple-900">{{ formatRupiah(wo.grand_total) }}</span>
              <router-link 
                :to="`/work-orders/${wo.id}`"
                class="bg-emerald-600 hover:bg-emerald-700 text-white text-[10px] font-bold px-2.5 py-1 rounded-lg transition flex items-center space-x-1">
                <CreditCard class="w-3 h-3" />
                <span>Buka Kasir</span>
              </router-link>
            </div>
          </div>

          <div v-if="getOrdersByStatus('completed').length === 0" class="text-center py-10 text-slate-400 text-xs">
            Tidak ada motor di kasir
          </div>
        </div>
      </div>

      <!-- Column 4: Lunas / Paid -->
      <div class="bg-emerald-50/50 p-3.5 rounded-2xl border border-emerald-200/60 space-y-3">
        <div class="flex items-center justify-between px-1">
          <div class="flex items-center space-x-2">
            <span class="w-2.5 h-2.5 rounded-full bg-emerald-600"></span>
            <h4 class="font-bold text-xs text-emerald-900 uppercase tracking-wider">4. Lunas (Paid)</h4>
          </div>
          <span class="bg-emerald-100 text-emerald-800 text-[10px] font-black px-2 py-0.5 rounded-full font-mono">
            {{ getOrdersByStatus('paid').length }}
          </span>
        </div>

        <div class="space-y-3 min-h-[300px]">
          <div 
            v-for="wo in getOrdersByStatus('paid').slice(0, 6)" 
            :key="wo.id" 
            class="bg-white p-4 rounded-xl border border-emerald-200/80 shadow-xs space-y-2.5">
            <div class="flex items-center justify-between">
              <span class="font-mono font-bold text-xs bg-emerald-50 text-emerald-800 px-2 py-0.5 rounded border border-emerald-200">
                {{ wo.customer?.plate_number }}
              </span>
              <span class="text-[10px] uppercase font-mono font-bold text-emerald-600">
                {{ wo.payment_method }}
              </span>
            </div>

            <div>
              <h5 class="font-bold text-slate-800 text-xs">{{ wo.customer?.motor_model }}</h5>
              <p class="text-[11px] text-slate-500">{{ wo.customer?.name }}</p>
            </div>

            <div class="pt-2 border-t border-slate-100 flex items-center justify-between">
              <span class="font-mono font-bold text-xs text-emerald-700">{{ formatRupiah(wo.grand_total) }}</span>
              <div class="flex items-center space-x-1">
                <button 
                  @click="openThermalModal(wo)"
                  title="Struk Thermal"
                  class="p-1 bg-slate-100 hover:bg-slate-200 rounded text-slate-600">
                  <Printer class="w-3.5 h-3.5" />
                </button>
                <button 
                  @click="openA4Modal(wo)"
                  title="Faktur A4"
                  class="p-1 bg-slate-100 hover:bg-slate-200 rounded text-slate-600">
                  <FileText class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- VIEW 2: TABLE & MOBILE CARDS VIEW -->
    <div v-else class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <!-- Desktop Table -->
      <div class="hidden md:block overflow-x-auto">
        <table class="w-full text-xs text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-500 font-bold uppercase text-[10px] tracking-wider">
              <th class="py-3 px-4">No. Work Order</th>
              <th class="py-3 px-4">Kendaraan & Pelanggan</th>
              <th class="py-3 px-4">Mekanik & Pit</th>
              <th class="py-3 px-4">Status Alur</th>
              <th class="py-3 px-4 text-right">Grand Total</th>
              <th class="py-3 px-4 text-center">Aksi Cepat</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr 
              v-for="wo in filteredOrders" 
              :key="wo.id" 
              class="hover:bg-slate-50/70 transition group">
              <td class="py-3.5 px-4">
                <router-link :to="`/work-orders/${wo.id}`" class="font-mono font-bold text-red-600 hover:text-red-700 block">
                  {{ wo.wo_number }}
                </router-link>
                <span class="text-[11px] text-slate-400">{{ formatDate(wo.created_at) }}</span>
              </td>

              <td class="py-3.5 px-4">
                <div class="flex items-center space-x-2">
                  <span class="font-mono font-bold bg-slate-100 text-slate-800 px-2 py-0.5 rounded border border-slate-200 text-[11px]">
                    {{ wo.customer?.plate_number }}
                  </span>
                  <span class="font-semibold text-slate-900">{{ wo.customer?.motor_model }}</span>
                </div>
                <p class="text-[11px] text-slate-500 mt-0.5">{{ wo.customer?.name }} ({{ wo.customer?.phone }})</p>
              </td>

              <td class="py-3.5 px-4">
                <div v-if="wo.mechanic" class="flex items-center space-x-1.5">
                  <span class="w-5 h-5 rounded-full bg-red-100 text-red-700 font-bold text-[10px] flex items-center justify-center">
                    {{ wo.mechanic.pit_number }}
                  </span>
                  <span class="font-medium text-slate-700">{{ wo.mechanic.name }}</span>
                </div>
                <span v-else class="text-slate-400 italic">Belum ditentukan</span>
              </td>

              <td class="py-3.5 px-4">
                <span 
                  :class="[getStatusBadge(wo.status).bgClass, getStatusBadge(wo.status).textClass, getStatusBadge(wo.status).borderClass]"
                  class="inline-block px-2.5 py-1 rounded-lg text-[10px] font-bold border uppercase tracking-wider">
                  {{ getStatusBadge(wo.status).label }}
                </span>
              </td>

              <td class="py-3.5 px-4 text-right">
                <p class="font-mono font-bold text-slate-900 text-sm">{{ formatRupiah(wo.grand_total) }}</p>
                <span 
                  :class="wo.payment_status === 'paid' ? 'text-emerald-600 font-bold' : 'text-amber-600 font-medium'"
                  class="text-[10px] block uppercase">
                  {{ wo.payment_status === 'paid' ? 'Lunas' : 'Belum Bayar' }}
                </span>
              </td>

              <td class="py-3.5 px-4 text-center">
                <div class="flex items-center justify-center space-x-1.5">
                  <router-link
                    :to="`/work-orders/${wo.id}`"
                    title="Buka Kasir POS"
                    class="p-1.5 bg-slate-100 hover:bg-red-50 hover:text-red-600 text-slate-600 rounded-lg transition">
                    <CreditCard class="w-4 h-4" />
                  </router-link>
                  <button
                    @click="openThermalModal(wo)"
                    title="Struk Thermal POS"
                    class="p-1.5 bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-lg transition">
                    <Printer class="w-4 h-4" />
                  </button>
                  <button
                    @click="openA4Modal(wo)"
                    title="Faktur A4 Resmi"
                    class="p-1.5 bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-lg transition">
                    <FileText class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredOrders.length === 0">
              <td colspan="6" class="py-12 text-center text-slate-400">
                Tidak ada data Work Order ditemukan.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mobile Responsive Cards (Visible on screens < 768px) -->
      <div class="md:hidden divide-y divide-slate-100">
        <div 
          v-for="wo in filteredOrders" 
          :key="wo.id" 
          class="p-4 space-y-3">
          <div class="flex items-center justify-between">
            <span class="font-mono font-bold text-xs text-red-600">{{ wo.wo_number }}</span>
            <span 
              :class="[getStatusBadge(wo.status).bgClass, getStatusBadge(wo.status).textClass, getStatusBadge(wo.status).borderClass]"
              class="px-2 py-0.5 rounded text-[10px] font-bold border uppercase">
              {{ getStatusBadge(wo.status).label }}
            </span>
          </div>

          <div class="flex items-center space-x-2">
            <span class="font-mono font-black text-xs bg-slate-100 px-2 py-0.5 rounded border border-slate-300">
              {{ wo.customer?.plate_number }}
            </span>
            <span class="font-bold text-xs text-slate-900">{{ wo.customer?.motor_model }}</span>
          </div>

          <div class="text-[11px] text-slate-500 flex justify-between">
            <span>Pemilik: {{ wo.customer?.name }}</span>
            <span class="font-semibold text-slate-700">{{ wo.mechanic ? `Pit ${wo.mechanic.pit_number}` : 'Antrean' }}</span>
          </div>

          <div class="pt-2 border-t border-slate-100 flex items-center justify-between">
            <div>
              <span class="text-[10px] text-slate-400 block">Total Biaya:</span>
              <span class="font-mono font-black text-sm text-slate-900">{{ formatRupiah(wo.grand_total) }}</span>
            </div>

            <div class="flex items-center space-x-1.5">
              <router-link
                :to="`/work-orders/${wo.id}`"
                class="bg-red-600 text-white font-bold text-xs px-3 py-1.5 rounded-lg">
                Kasir &rarr;
              </router-link>
              <button
                @click="openThermalModal(wo)"
                class="p-1.5 bg-slate-100 text-slate-700 rounded-lg">
                <Printer class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <ThermalReceiptModal 
      :is-open="showThermalModal" 
      :work-order="selectedWO" 
      @close="showThermalModal = false" />

    <A4InvoiceModal 
      :is-open="showA4Modal" 
      :work-order="selectedWO" 
      @close="showA4Modal = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Search,
  Plus,
  CreditCard,
  Printer,
  FileText,
  Table,
  Kanban,
  ArrowRight,
  Check,
} from 'lucide-vue-next'
import { useWorkOrderStore } from '../stores/workOrder'
import { useToastStore } from '../stores/toast'
import { playClickSound, playSuccessChime } from '../utils/sound'
import { formatRupiah, formatDate, getStatusBadge } from '../utils/format'
import ThermalReceiptModal from '../components/ThermalReceiptModal.vue'
import A4InvoiceModal from '../components/A4InvoiceModal.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'
import type { WorkOrder, WorkOrderStatus } from '../types'

const workOrderStore = useWorkOrderStore()
const toastStore = useToastStore()

const searchQuery = ref('')
const selectedTab = ref('all')
const viewMode = ref<'table' | 'kanban'>('kanban') // Default to Kanban for visual impact!

const filterTabs = [
  { label: 'Semua Status', value: 'all' },
  { label: 'Draft / Masuk', value: 'draft' },
  { label: 'In Progress (Pit)', value: 'in_progress' },
  { label: 'Selesai Servis', value: 'completed' },
  { label: 'Lunas Kasir', value: 'paid' },
]

const showThermalModal = ref(false)
const showA4Modal = ref(false)
const selectedWO = ref<WorkOrder | null>(null)

function openThermalModal(wo: WorkOrder) {
  playClickSound()
  selectedWO.value = wo
  showThermalModal.value = true
}

function openA4Modal(wo: WorkOrder) {
  playClickSound()
  selectedWO.value = wo
  showA4Modal.value = true
}

function getOrdersByStatus(status: WorkOrderStatus) {
  return filteredOrders.value.filter((wo) => wo.status === status)
}

async function advanceStatus(id: number, nextStatus: WorkOrderStatus) {
  playClickSound()
  try {
    await workOrderStore.updateStatus(id, nextStatus)
    playSuccessChime()
    toastStore.success(`Status Work Order berhasil diubah ke ${getStatusBadge(nextStatus).label}`)
    workOrderStore.fetchWorkOrders()
  } catch (err: any) {
    toastStore.error(err.message || 'Gagal mengubah status')
  }
}

onMounted(() => {
  workOrderStore.fetchWorkOrders()
})

const filteredOrders = computed(() => {
  let list = workOrderStore.workOrders

  if (selectedTab.value !== 'all' && viewMode.value === 'table') {
    list = list.filter((wo) => wo.status === selectedTab.value)
  }

  if (searchQuery.value.trim() !== '') {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(
      (wo) =>
        wo.wo_number.toLowerCase().includes(q) ||
        wo.customer?.name.toLowerCase().includes(q) ||
        wo.customer?.plate_number.toLowerCase().includes(q)
    )
  }

  return list
})
</script>
