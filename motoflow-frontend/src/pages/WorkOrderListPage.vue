<template>
  <div class="space-y-6">
    <!-- Controls Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <!-- Filter Tabs -->
      <div class="flex items-center space-x-2 bg-slate-200/70 p-1 rounded-xl text-xs font-semibold overflow-x-auto">
        <button
          v-for="tab in filterTabs"
          :key="tab.value"
          @click="selectedTab = tab.value"
          :class="selectedTab === tab.value ? 'bg-white text-slate-900 shadow-xs' : 'text-slate-600 hover:text-slate-900'"
          class="px-3.5 py-1.5 rounded-lg transition whitespace-nowrap">
          {{ tab.label }}
        </button>
      </div>

      <!-- Search & New Button -->
      <div class="flex items-center space-x-3">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari WO, Plat, atau Nama..."
            class="bg-white border border-slate-200 rounded-xl pl-9 pr-4 py-2 text-xs text-slate-900 placeholder-slate-400 focus:outline-none focus:border-red-500 w-64 shadow-xs" />
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
        </div>

        <router-link
          to="/work-orders/create"
          class="inline-flex items-center space-x-1.5 bg-red-600 hover:bg-red-700 text-white font-semibold text-xs px-4 py-2 rounded-xl shadow-xs transition">
          <Plus class="w-4 h-4" />
          <span>Buat WO</span>
        </router-link>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="workOrderStore.loading">
      <SkeletonLoader />
    </div>

    <!-- Table of Work Orders -->
    <div v-else class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
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
              <!-- WO Number & Time -->
              <td class="py-3.5 px-4">
                <router-link :to="`/work-orders/${wo.id}`" class="font-mono font-bold text-red-600 hover:text-red-700 block">
                  {{ wo.wo_number }}
                </router-link>
                <span class="text-[11px] text-slate-400">{{ formatDate(wo.created_at) }}</span>
              </td>

              <!-- Customer & Bike -->
              <td class="py-3.5 px-4">
                <div class="flex items-center space-x-2">
                  <span class="font-mono font-bold bg-slate-100 text-slate-800 px-2 py-0.5 rounded border border-slate-200 text-[11px]">
                    {{ wo.customer?.plate_number }}
                  </span>
                  <span class="font-semibold text-slate-900">{{ wo.customer?.motor_model }}</span>
                </div>
                <p class="text-[11px] text-slate-500 mt-0.5">{{ wo.customer?.name }} ({{ wo.customer?.phone }})</p>
              </td>

              <!-- Mechanic -->
              <td class="py-3.5 px-4">
                <div v-if="wo.mechanic" class="flex items-center space-x-1.5">
                  <span class="w-5 h-5 rounded-full bg-red-100 text-red-700 font-bold text-[10px] flex items-center justify-center">
                    {{ wo.mechanic.pit_number }}
                  </span>
                  <span class="font-medium text-slate-700">{{ wo.mechanic.name }}</span>
                </div>
                <span v-else class="text-slate-400 italic">Belum ditentukan</span>
              </td>

              <!-- Status Badge -->
              <td class="py-3.5 px-4">
                <span 
                  :class="[getStatusBadge(wo.status).bgClass, getStatusBadge(wo.status).textClass, getStatusBadge(wo.status).borderClass]"
                  class="inline-block px-2.5 py-1 rounded-lg text-[10px] font-bold border uppercase tracking-wider">
                  {{ getStatusBadge(wo.status).label }}
                </span>
              </td>

              <!-- Grand Total -->
              <td class="py-3.5 px-4 text-right">
                <p class="font-mono font-bold text-slate-900 text-sm">{{ formatRupiah(wo.grand_total) }}</p>
                <span 
                  :class="wo.payment_status === 'paid' ? 'text-emerald-600 font-bold' : 'text-amber-600 font-medium'"
                  class="text-[10px] block uppercase">
                  {{ wo.payment_status === 'paid' ? 'Lunas' : 'Belum Bayar' }}
                </span>
              </td>

              <!-- Actions -->
              <td class="py-3.5 px-4 text-center">
                <div class="flex items-center justify-center space-x-1.5">
                  <!-- POS Billing / Detail -->
                  <router-link
                    :to="`/work-orders/${wo.id}`"
                    title="Buka Kasir POS"
                    class="p-1.5 bg-slate-100 hover:bg-red-50 hover:text-red-600 text-slate-600 rounded-lg transition">
                    <CreditCard class="w-4 h-4" />
                  </router-link>

                  <!-- Thermal Receipt Modal -->
                  <button
                    @click="openThermalModal(wo)"
                    title="Struk Thermal POS"
                    class="p-1.5 bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-lg transition">
                    <Printer class="w-4 h-4" />
                  </button>

                  <!-- A4 Invoice Modal -->
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
} from 'lucide-vue-next'
import { useWorkOrderStore } from '../stores/workOrder'
import { formatRupiah, formatDate, getStatusBadge } from '../utils/format'
import ThermalReceiptModal from '../components/ThermalReceiptModal.vue'
import A4InvoiceModal from '../components/A4InvoiceModal.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'
import type { WorkOrder } from '../types'

const workOrderStore = useWorkOrderStore()
const searchQuery = ref('')
const selectedTab = ref('all')

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
  selectedWO.value = wo
  showThermalModal.value = true
}

function openA4Modal(wo: WorkOrder) {
  selectedWO.value = wo
  showA4Modal.value = true
}

onMounted(() => {
  workOrderStore.fetchWorkOrders()
})

const filteredOrders = computed(() => {
  let list = workOrderStore.workOrders

  if (selectedTab.value !== 'all') {
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

