<template>
  <div class="max-w-6xl mx-auto space-y-6">
    <div v-if="workOrderStore.loading && !wo">
      <SkeletonLoader />
    </div>

    <div v-else-if="wo" class="space-y-6">
      <!-- Top Bar: WO Number, Status & Action Print Buttons -->
      <div class="bg-white p-4 sm:p-6 rounded-2xl border border-slate-200/80 shadow-xs flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div class="flex flex-wrap items-center gap-2">
            <h1 class="text-xl sm:text-2xl font-black text-slate-900 font-mono tracking-tight">{{ wo.wo_number }}</h1>
            <span 
              :class="[getStatusBadge(wo.status).bgClass, getStatusBadge(wo.status).textClass, getStatusBadge(wo.status).borderClass]"
              class="px-2.5 py-0.5 rounded-lg text-xs font-bold border uppercase tracking-wider">
              {{ getStatusBadge(wo.status).label }}
            </span>
          </div>
          <p class="text-xs text-slate-500 mt-1">Diterbitkan pada {{ formatDate(wo.created_at) }}</p>
        </div>

        <div class="flex items-center space-x-2 sm:space-x-3">
          <!-- Thermal Receipt Button -->
          <button
            @click="openThermalReceipt"
            title="Hotkey F9"
            class="flex-1 sm:flex-initial inline-flex items-center justify-center space-x-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs px-3.5 py-2 rounded-xl transition border border-slate-300">
            <Printer class="w-4 h-4 text-slate-600" />
            <span>Struk POS (F9)</span>
          </button>

          <!-- A4 Invoice Button -->
          <button
            @click="showA4Modal = true"
            class="flex-1 sm:flex-initial inline-flex items-center justify-center space-x-1.5 bg-red-600 hover:bg-red-700 text-white font-semibold text-xs px-4 py-2 rounded-xl shadow transition">
            <FileText class="w-4 h-4" />
            <span>Faktur A4</span>
          </button>
        </div>
      </div>

      <!-- Service Workflow Stepper -->
      <div class="bg-white p-4 sm:p-6 rounded-2xl border border-slate-200/80 shadow-xs">
        <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-3">Progres Servis Motor Honda</h4>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 text-center text-xs">
          <!-- Step 1: Draft -->
          <div :class="isStepActive('draft') ? 'border-red-600 bg-red-50 text-red-900 font-bold' : 'border-slate-200 text-slate-400'"
               class="p-2.5 sm:p-3 rounded-xl border transition">
            <span class="block text-[10px] uppercase tracking-wider">Tahap 1</span>
            <span>Pendaftaran SA</span>
          </div>

          <!-- Step 2: In Progress -->
          <div :class="isStepActive('in_progress') ? 'border-blue-600 bg-blue-50 text-blue-900 font-bold' : 'border-slate-200 text-slate-400'"
               class="p-2.5 sm:p-3 rounded-xl border transition">
            <span class="block text-[10px] uppercase tracking-wider">Tahap 2</span>
            <span>Pengerjaan Pit</span>
          </div>

          <!-- Step 3: Completed -->
          <div :class="isStepActive('completed') ? 'border-purple-600 bg-purple-50 text-purple-900 font-bold' : 'border-slate-200 text-slate-400'"
               class="p-2.5 sm:p-3 rounded-xl border transition">
            <span class="block text-[10px] uppercase tracking-wider">Tahap 3</span>
            <span>Inspeksi Akhir</span>
          </div>

          <!-- Step 4: Paid -->
          <div :class="isStepActive('paid') ? 'border-emerald-600 bg-emerald-50 text-emerald-900 font-bold' : 'border-slate-200 text-slate-400'"
               class="p-2.5 sm:p-3 rounded-xl border transition">
            <span class="block text-[10px] uppercase tracking-wider">Tahap 4</span>
            <span>Kasir Lunas</span>
          </div>
        </div>

        <!-- Workflow Action Buttons -->
        <div class="mt-4 pt-3 border-t border-slate-100 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <span class="text-xs text-slate-500">Tindakan Alur Status:</span>
          <div class="flex items-center space-x-2">
            <button
              v-if="wo.status === 'draft'"
              @click="handleStatusChange('in_progress')"
              class="w-full sm:w-auto bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold px-4 py-2 rounded-lg shadow-sm transition">
              Mulai Kerjakan di Pit Servis &rarr;
            </button>
            <button
              v-if="wo.status === 'in_progress'"
              @click="handleStatusChange('completed')"
              class="w-full sm:w-auto bg-purple-600 hover:bg-purple-700 text-white text-xs font-bold px-4 py-2 rounded-lg shadow-sm transition">
              Selesaikan Servis & Inspeksi QC &rarr;
            </button>
            <span v-if="wo.status === 'paid'" class="text-emerald-700 font-bold text-xs bg-emerald-50 px-3 py-1.5 rounded-lg border border-emerald-200">
              Transaksi Selesai & Lunas
            </span>
          </div>
        </div>
      </div>

      <!-- 2-Columns: Left Items, Right Customer & POS Cashier Terminal -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Left: Items & Parts Table (2 Cols) -->
        <div class="lg:col-span-2 space-y-6">
          <div class="bg-white p-4 sm:p-6 rounded-2xl border border-slate-200/80 shadow-xs space-y-4">
            <h3 class="font-bold text-slate-800 text-sm uppercase tracking-wider">Rincian Suku Cadang HGP & Jasa Mekanik</h3>

            <!-- Desktop Items Table -->
            <div class="hidden sm:block overflow-x-auto">
              <table class="w-full text-xs text-left">
                <thead>
                  <tr class="border-b border-slate-200 text-slate-500 font-bold uppercase text-[10px]">
                    <th class="py-2.5">Deskripsi</th>
                    <th class="py-2.5 text-center">Qty</th>
                    <th class="py-2.5 text-right">Harga</th>
                    <th class="py-2.5 text-center">Disc</th>
                    <th class="py-2.5 text-right">Subtotal</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr v-for="item in wo.items" :key="item.id" class="hover:bg-slate-50/50">
                    <td class="py-3">
                      <span v-if="item.item_type === 'service'" class="text-blue-600 font-bold mr-1">[JASA]</span>
                      <span class="font-semibold text-slate-800">{{ item.item_name_snapshot }}</span>
                      <p class="text-[10px] text-slate-400 font-mono">{{ item.part_number_snapshot }}</p>
                    </td>
                    <td class="py-3 text-center font-mono font-bold text-slate-700">{{ item.quantity }}</td>
                    <td class="py-3 text-right font-mono text-slate-600">{{ formatRupiah(item.price_snapshot) }}</td>
                    <td class="py-3 text-center font-mono text-slate-500">{{ Number(item.discount_percent) > 0 ? `${item.discount_percent}%` : '-' }}</td>
                    <td class="py-3 text-right font-mono font-bold text-slate-900">{{ formatRupiah(item.subtotal) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Mobile Item Cards (on screens < 640px) -->
            <div class="sm:hidden divide-y divide-slate-100">
              <div v-for="item in wo.items" :key="item.id" class="py-2.5 space-y-1 text-xs">
                <div class="flex justify-between items-start">
                  <div>
                    <span v-if="item.item_type === 'service'" class="text-blue-600 font-bold mr-1">[JASA]</span>
                    <span class="font-bold text-slate-800">{{ item.item_name_snapshot }}</span>
                    <p class="text-[10px] text-slate-400 font-mono">{{ item.part_number_snapshot }}</p>
                  </div>
                  <span class="font-mono font-bold text-slate-900">{{ formatRupiah(item.subtotal) }}</span>
                </div>
                <div class="flex justify-between text-[11px] text-slate-500">
                  <span>{{ item.quantity }} pcs &times; {{ formatRupiah(item.price_snapshot) }}</span>
                  <span v-if="Number(item.discount_percent) > 0" class="text-emerald-600 font-semibold">Diskon {{ item.discount_percent }}%</span>
                </div>
              </div>
            </div>

            <!-- Totals Breakdown -->
            <div class="border-t border-slate-200 pt-4 space-y-2 text-xs">
              <div class="flex justify-between text-slate-600">
                <span>Subtotal Suku Cadang:</span>
                <span class="font-mono font-medium">{{ formatRupiah(wo.parts_subtotal) }}</span>
              </div>
              <div class="flex justify-between text-slate-600">
                <span>Subtotal Jasa Servis:</span>
                <span class="font-mono font-medium">{{ formatRupiah(wo.services_subtotal) }}</span>
              </div>
              <div v-if="Number(wo.total_discount) > 0" class="flex justify-between text-emerald-600 font-medium">
                <span>Total Diskon Promosi:</span>
                <span class="font-mono">-{{ formatRupiah(wo.total_discount) }}</span>
              </div>
              <div class="flex justify-between text-slate-600">
                <span>PPN ({{ wo.tax_rate }}%):</span>
                <span class="font-mono font-medium">{{ formatRupiah(wo.tax_amount) }}</span>
              </div>
              <div class="flex justify-between text-base font-black text-slate-900 pt-2 border-t border-slate-300">
                <span class="text-red-600">TOTAL BIAYA:</span>
                <span class="font-mono">{{ formatRupiah(wo.grand_total) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Customer Info & POS Cashier Terminal (1 Col) -->
        <div class="space-y-6">
          <!-- Customer & Bike Card -->
          <div class="bg-white p-4 sm:p-6 rounded-2xl border border-slate-200/80 shadow-xs space-y-3 text-xs">
            <h4 class="font-bold text-slate-800 uppercase tracking-wider text-[11px] border-b border-slate-100 pb-2">
              Data Konsumen & Unit
            </h4>
            <div class="space-y-1.5">
              <p><span class="text-slate-500">Pemilik:</span> <span class="font-bold text-slate-800">{{ wo.customer?.name }}</span></p>
              <p><span class="text-slate-500">No. HP:</span> {{ wo.customer?.phone }}</p>
              <div class="flex items-center space-x-1.5 pt-1">
                <span class="font-mono font-black text-xs bg-slate-100 text-slate-800 px-2 py-0.5 rounded border border-slate-300">
                  {{ wo.customer?.plate_number }}
                </span>
                <span class="font-semibold text-slate-900">{{ wo.customer?.motor_model }}</span>
              </div>
              <p><span class="text-slate-500">Odometer:</span> <span class="font-mono font-semibold">{{ wo.entry_km }} KM</span></p>
              <p><span class="text-slate-500">Mekanik:</span> <span class="font-semibold">{{ wo.mechanic?.name }} (Pit {{ wo.mechanic?.pit_number }})</span></p>
            </div>
          </div>

          <!-- POS Cashier Terminal -->
          <div class="bg-white p-4 sm:p-6 rounded-2xl border border-slate-200/80 shadow-xs space-y-4">
            <div class="flex items-center justify-between border-b border-slate-100 pb-2">
              <h4 class="font-bold text-slate-800 uppercase tracking-wider text-[11px]">
                Terminal Kasir POS
              </h4>
              <span 
                :class="wo.payment_status === 'paid' ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'"
                class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider">
                {{ wo.payment_status === 'paid' ? 'LUNAS' : 'BELUM BAYAR' }}
              </span>
            </div>

            <!-- When Already Paid -->
            <div v-if="wo.payment_status === 'paid'" class="space-y-3 bg-emerald-50/50 p-4 rounded-xl border border-emerald-200/60 text-xs">
              <div class="flex justify-between text-slate-600">
                <span>Metode Bayar:</span>
                <span class="font-mono font-bold uppercase text-slate-800">{{ wo.payment_method }}</span>
              </div>
              <div class="flex justify-between text-slate-600">
                <span>Uang Diterima:</span>
                <span class="font-mono font-bold text-slate-800">{{ formatRupiah(wo.paid_amount) }}</span>
              </div>
              <div class="flex justify-between text-slate-600">
                <span>Kembalian:</span>
                <span class="font-mono font-bold text-emerald-700">{{ formatRupiah(wo.change_amount) }}</span>
              </div>
              <div class="text-[10px] text-slate-400 pt-2 border-t border-emerald-200">
                Lunas pada {{ formatDate(wo.paid_at) }}
              </div>
            </div>

            <!-- Cashier Payment Form (When Unpaid) -->
            <form v-else @submit.prevent="handleProcessPayment" class="space-y-4 text-xs">
              <div>
                <label class="block font-semibold text-slate-700 mb-1.5 flex justify-between">
                  <span>Metode Pembayaran</span>
                  <span class="text-[10px] text-slate-400 font-normal">Tekan F4 untuk ganti</span>
                </label>
                <div class="grid grid-cols-2 gap-2">
                  <button
                    type="button"
                    v-for="method in ['cash', 'qris', 'transfer', 'debit'] as const"
                    :key="method"
                    @click="setPaymentMethod(method)"
                    :class="paymentMethod === method ? 'bg-red-600 text-white font-bold shadow-xs' : 'bg-slate-50 text-slate-700 hover:bg-slate-100 border border-slate-200'"
                    class="py-2.5 px-2 rounded-xl text-center uppercase text-[11px] transition flex items-center justify-center space-x-1">
                    <QrCode v-if="method === 'qris'" class="w-3.5 h-3.5" />
                    <span>{{ method }}</span>
                  </button>
                </div>
              </div>

              <!-- Button to open QRIS Modal if QRIS selected -->
              <div v-if="paymentMethod === 'qris'">
                <button
                  type="button"
                  @click="showQRISModal = true"
                  class="w-full bg-slate-900 hover:bg-slate-800 text-white font-bold py-2.5 rounded-xl text-xs flex items-center justify-center space-x-2 transition">
                  <QrCode class="w-4 h-4 text-red-500" />
                  <span>Tampilkan QRIS di Layar</span>
                </button>
              </div>

              <div>
                <label class="block font-semibold text-slate-700 mb-1.5 flex justify-between">
                  <span>Nominal Uang Diterima</span>
                  <span class="text-[10px] text-slate-400 font-normal">F8: Uang Pas</span>
                </label>
                <input
                  v-model.number="paidAmount"
                  type="number"
                  min="0"
                  required
                  class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 font-mono text-sm font-bold text-slate-900 focus:outline-none focus:border-red-500" />
              </div>

              <!-- Quick Denomination Buttons -->
              <div class="grid grid-cols-4 gap-1.5">
                <button
                  type="button"
                  @click="setExactCash"
                  class="bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold py-1.5 rounded-lg text-[10px] transition">
                  Uang Pas
                </button>
                <button
                  type="button"
                  @click="setDenomination(100000)"
                  class="bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold py-1.5 rounded-lg text-[10px] transition">
                  100k
                </button>
                <button
                  type="button"
                  @click="setDenomination(200000)"
                  class="bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold py-1.5 rounded-lg text-[10px] transition">
                  200k
                </button>
                <button
                  type="button"
                  @click="setDenomination(500000)"
                  class="bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold py-1.5 rounded-lg text-[10px] transition">
                  500k
                </button>
              </div>

              <!-- Live Kembalian Display -->
              <div class="bg-slate-50 p-3 rounded-xl border border-slate-200 space-y-1">
                <div class="flex justify-between font-medium text-slate-500">
                  <span>Kembalian:</span>
                  <span class="font-mono font-bold" :class="liveChange >= 0 ? 'text-emerald-600' : 'text-red-500'">
                    {{ liveChange >= 0 ? formatRupiah(liveChange) : `Kurang ${formatRupiah(Math.abs(liveChange))}` }}
                  </span>
                </div>
              </div>

              <button
                type="submit"
                :disabled="workOrderStore.loading || liveChange < 0"
                class="w-full bg-emerald-600 hover:bg-emerald-700 text-white font-bold py-3 rounded-xl shadow-md transition disabled:opacity-50 text-xs">
                <span v-if="!workOrderStore.loading">Selesaikan Pembayaran Kasir</span>
                <span v-else>Memproses Transaksi...</span>
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Thermal Receipt Modal -->
    <ThermalReceiptModal 
      :is-open="showThermalModal" 
      :work-order="wo" 
      @close="showThermalModal = false" />

    <!-- A4 Formal Invoice Modal -->
    <A4InvoiceModal 
      :is-open="showA4Modal" 
      :work-order="wo" 
      @close="showA4Modal = false" />

    <!-- Dynamic QRIS Modal -->
    <QRISModal
      :is-open="showQRISModal"
      :amount="wo?.grand_total || 0"
      @close="showQRISModal = false"
      @verified="onQRISVerified" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { Printer, FileText, QrCode } from 'lucide-vue-next'
import { useWorkOrderStore } from '../stores/workOrder'
import { useToastStore } from '../stores/toast'
import { playScannerBeep, playSuccessChime, playClickSound } from '../utils/sound'
import { formatRupiah, formatDate, getStatusBadge } from '../utils/format'
import ThermalReceiptModal from '../components/ThermalReceiptModal.vue'
import A4InvoiceModal from '../components/A4InvoiceModal.vue'
import QRISModal from '../components/QRISModal.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'
import type { PaymentMethod, WorkOrderStatus } from '../types'

const route = useRoute()
const workOrderStore = useWorkOrderStore()
const toastStore = useToastStore()

const wo = computed(() => workOrderStore.currentWorkOrder)

const showThermalModal = ref(false)
const showA4Modal = ref(false)
const showQRISModal = ref(false)

const paymentMethod = ref<PaymentMethod>('cash')
const paidAmount = ref<number>(0)

onMounted(async () => {
  const id = Number(route.params.id)
  const data = await workOrderStore.fetchWorkOrder(id)
  if (data) {
    paidAmount.value = Number(data.grand_total)
  }
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})

function handleKeydown(e: KeyboardEvent) {
  // F9: Open thermal receipt modal
  if (e.key === 'F9') {
    e.preventDefault()
    openThermalReceipt()
  }
  // F8: Set exact cash
  if (e.key === 'F8') {
    e.preventDefault()
    setExactCash()
  }
  // F4: Cycle payment method
  if (e.key === 'F4') {
    e.preventDefault()
    const methods: PaymentMethod[] = ['cash', 'qris', 'transfer', 'debit']
    const idx = methods.indexOf(paymentMethod.value)
    paymentMethod.value = methods[(idx + 1) % methods.length]
    playClickSound()
  }
  // Escape: Close active modals
  if (e.key === 'Escape') {
    showThermalModal.value = false
    showA4Modal.value = false
    showQRISModal.value = false
  }
}

function openThermalReceipt() {
  playClickSound()
  showThermalModal.value = true
}

function setPaymentMethod(method: PaymentMethod) {
  playClickSound()
  paymentMethod.value = method
  if (method === 'qris') {
    showQRISModal.value = true
  }
}

function setExactCash() {
  if (!wo.value) return
  paidAmount.value = Number(wo.value.grand_total)
  playClickSound()
}

function setDenomination(val: number) {
  paidAmount.value = val
  playClickSound()
}

function isStepActive(step: string) {
  if (!wo.value) return false
  const order = ['draft', 'in_progress', 'completed', 'paid']
  const currentIdx = order.indexOf(wo.value.status)
  const targetIdx = order.indexOf(step)
  return currentIdx >= targetIdx
}

const liveChange = computed(() => {
  if (!wo.value) return 0
  return paidAmount.value - Number(wo.value.grand_total)
})

async function handleStatusChange(status: WorkOrderStatus) {
  if (!wo.value) return
  playClickSound()
  try {
    await workOrderStore.updateStatus(wo.value.id, status)
    playSuccessChime()
    toastStore.success(`Status Work Order berhasil diubah ke ${getStatusBadge(status).label}`)
  } catch (err: any) {
    toastStore.error(err.message || 'Gagal mengubah status')
  }
}

function onQRISVerified() {
  showQRISModal.value = false
  paidAmount.value = Number(wo.value?.grand_total || 0)
  handleProcessPayment()
}

async function handleProcessPayment() {
  if (!wo.value) return
  if (liveChange.value < 0) {
    toastStore.error('Jumlah uang pembayaran tidak mencukupi')
    return
  }

  playScannerBeep()
  try {
    await workOrderStore.processPayment(wo.value.id, paymentMethod.value, paidAmount.value)
    playSuccessChime()
    toastStore.success('Pembayaran kasir berhasil diselesaikan!')
    showThermalModal.value = true // Automatically show receipt on payment completion!
  } catch (err: any) {
    toastStore.error(err.message || 'Gagal memproses pembayaran')
  }
}
</script>
