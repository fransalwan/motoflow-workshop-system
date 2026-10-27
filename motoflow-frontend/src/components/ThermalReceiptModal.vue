<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-sm flex items-end sm:items-center justify-center p-0 sm:p-4">
    <div class="bg-white rounded-t-3xl sm:rounded-2xl shadow-2xl max-w-lg w-full overflow-hidden border border-slate-200 max-h-[92vh] flex flex-col">
      <div class="w-12 h-1.5 bg-slate-300 rounded-full mx-auto my-2 sm:hidden"></div>
      <!-- Modal Header -->
      <div class="bg-slate-900 text-white px-5 sm:px-6 py-4 flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <Printer class="w-5 h-5 text-red-500" />
          <h3 class="font-bold text-base sm:text-lg">Cetak Struk Kasir POS Thermal</h3>
        </div>
        <button @click="$emit('close')" class="text-slate-400 hover:text-white transition">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Controls & Roll Width Selection -->
      <div class="bg-slate-100 px-4 sm:px-6 py-2.5 border-b border-slate-200 flex flex-wrap items-center justify-between gap-2 text-xs font-medium text-slate-600">
        <div class="flex items-center space-x-2">
          <span>Lebar Kertas:</span>
          <button 
            @click="changeWidth(58)" 
            :class="selectedWidth === 58 ? 'bg-red-600 text-white font-bold' : 'bg-white text-slate-700 hover:bg-slate-200'"
            class="px-2.5 py-1 rounded shadow-xs transition border border-slate-300">
            58 mm
          </button>
          <button 
            @click="changeWidth(80)" 
            :class="selectedWidth === 80 ? 'bg-red-600 text-white font-bold' : 'bg-white text-slate-700 hover:bg-slate-200'"
            class="px-2.5 py-1 rounded shadow-xs transition border border-slate-300">
            80 mm
          </button>
        </div>
        <span class="text-emerald-700 bg-emerald-100 px-2 py-0.5 rounded font-mono text-[10px]">ESC/POS Ready</span>
      </div>

      <!-- Thermal Paper Preview Area -->
      <div class="p-3 sm:p-6 bg-slate-200 flex justify-center max-h-[55vh] overflow-y-auto">
        <div 
          :style="{ width: selectedWidth === 58 ? 'min(300px, 100%)' : 'min(400px, 100%)' }"
          class="bg-white p-4 sm:p-5 shadow-md border border-slate-300 font-mono text-xs leading-relaxed text-slate-900 print-area transition-all duration-200">
          
          <pre class="whitespace-pre-wrap font-mono">{{ receiptText }}</pre>

          <!-- Barcode simulation at bottom -->
          <div class="mt-4 pt-3 border-t border-dashed border-slate-400 text-center">
            <div class="h-8 bg-repeat-x bg-[length:6px_100%] inline-block w-48 opacity-80" 
                 style="background-image: linear-gradient(90deg, #000 50%, transparent 50%);"></div>
            <p class="text-[10px] tracking-widest mt-1 text-slate-600 font-mono">{{ workOrder?.wo_number }}</p>
          </div>
        </div>
      </div>

      <!-- Action Footer -->
      <div class="px-6 py-4 bg-slate-50 border-t border-slate-200 flex items-center justify-between">
        <button 
          @click="copyText" 
          class="inline-flex items-center space-x-1.5 text-xs text-slate-600 hover:text-slate-900 bg-white border border-slate-300 px-3 py-2 rounded-lg font-medium shadow-sm hover:bg-slate-100 transition">
          <Copy class="w-3.5 h-3.5" />
          <span>{{ copied ? 'Tersalin!' : 'Salin Teks Mentah' }}</span>
        </button>

        <div class="flex items-center space-x-3">
          <button 
            @click="$emit('close')" 
            class="px-4 py-2 text-xs font-medium text-slate-600 hover:text-slate-800 transition">
            Tutup
          </button>
          <button 
            @click="printReceipt" 
            class="inline-flex items-center space-x-2 bg-red-600 hover:bg-red-700 text-white font-semibold text-xs px-5 py-2.5 rounded-lg shadow transition">
            <Printer class="w-4 h-4" />
            <span>Kirim ke Printer Kasir</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Printer, X, Copy } from 'lucide-vue-next'
import { useWorkOrderStore } from '../stores/workOrder'
import { printThermalReceipt } from '../utils/printer'
import type { WorkOrder } from '../types'

const props = defineProps<{
  isOpen: boolean
  workOrder: WorkOrder | null
}>()

const emit = defineEmits(['close'])

const workOrderStore = useWorkOrderStore()
const selectedWidth = ref<number>(58)
const receiptText = ref<string>('Memuat struk...')
const copied = ref<boolean>(false)

async function loadReceipt() {
  if (!props.workOrder) return
  try {
    const data = await workOrderStore.getThermalReceipt(props.workOrder.id, selectedWidth.value)
    receiptText.value = data.raw_text
  } catch (err) {
    receiptText.value = 'Gagal memuat teks struk.'
  }
}

watch(
  () => [props.isOpen, props.workOrder],
  ([open]) => {
    if (open) {
      loadReceipt()
    }
  }
)

function changeWidth(w: number) {
  selectedWidth.value = w
  loadReceipt()
}

function copyText() {
  navigator.clipboard.writeText(receiptText.value)
  copied.value = true
  setTimeout(() => (copied.value = false), 2000)
}

function printReceipt() {
  printThermalReceipt(receiptText.value, {
    widthMm: selectedWidth.value,
    title: `Struk POS AHASS - ${props.workOrder?.wo_number || ''}`,
    barcodeText: props.workOrder?.wo_number,
  })
}
</script>

