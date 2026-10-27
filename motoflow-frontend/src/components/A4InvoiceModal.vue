<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-sm flex items-end sm:items-center justify-center p-0 sm:p-4">
    <div class="bg-white rounded-t-3xl sm:rounded-2xl shadow-2xl max-w-4xl w-full overflow-hidden border border-slate-200 max-h-[92vh] flex flex-col">
      <div class="w-12 h-1.5 bg-slate-300 rounded-full mx-auto my-2 sm:hidden"></div>
      <!-- Modal Top Bar -->
      <div class="bg-slate-900 text-white px-4 sm:px-6 py-3.5 sm:py-4 flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <FileText class="w-5 h-5 text-red-500" />
          <h3 class="font-bold text-sm sm:text-lg">Faktur Resmi Servis & Suku Cadang (A4)</h3>
        </div>
        <div class="flex items-center space-x-2 sm:space-x-3">
          <button 
            @click="printDirect"
            title="Cetak langsung ke printer A4"
            class="inline-flex items-center space-x-1 sm:space-x-1.5 bg-slate-800 hover:bg-slate-700 text-white text-xs font-semibold px-2.5 sm:px-3 py-1.5 rounded-lg shadow-xs transition active:scale-95 border border-slate-700">
            <Printer class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-slate-300" />
            <span class="hidden xs:inline">Cetak A4</span>
          </button>
          <button 
            @click="downloadPDF"
            title="Unduh dokumen PDF resmi"
            class="inline-flex items-center space-x-1 sm:space-x-1.5 bg-red-600 hover:bg-red-700 text-white text-xs font-semibold px-2.5 sm:px-3.5 py-1.5 rounded-lg shadow-xs transition active:scale-95">
            <Download class="w-3.5 h-3.5 sm:w-4 sm:h-4" />
            <span class="hidden xs:inline">Unduh PDF</span>
          </button>
          <button @click="$emit('close')" class="text-slate-400 hover:text-white transition">
            <X class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- A4 Sheet Paper Preview (Horizontally scrollable on mobile) -->
      <div class="p-2 sm:p-8 bg-slate-100 max-h-[75vh] overflow-x-auto overflow-y-auto flex justify-start sm:justify-center">
        <div 
          ref="invoiceSheetRef"
          class="bg-white w-[210mm] min-w-[700px] min-h-[260mm] p-6 sm:p-8 shadow-xl border border-slate-200 rounded-sm text-slate-900 flex flex-col justify-between">
          <div>
            <!-- AHASS Official Red Header Bar -->
            <div class="bg-red-700 text-white text-center py-2 px-4 rounded font-bold text-xs tracking-wider uppercase mb-6">
              HONDA AHASS - PT NUSANTARA SAKTI GROUP (OFFICIAL AUTHORIZED WORKSHOP)
            </div>

            <!-- Header Info -->
            <div class="flex justify-between items-start border-b border-slate-200 pb-6 mb-6">
              <div>
                <h1 class="text-xl font-black text-slate-900 tracking-tight">FAKTUR BIAYA SERVIS & SUKU CADANG</h1>
                <p class="text-xs font-semibold text-red-600 mt-1">AHASS 01234 NSS Cempaka Putih</p>
                <p class="text-xs text-slate-500 mt-0.5">Jl. Cempaka Putih Raya No. 128, Jakarta Pusat</p>
                <p class="text-xs text-slate-500">Telp: (021) 420-5678 | Email: ahass.nss@nusantara-sakti.com</p>
              </div>
              <div class="bg-slate-50 border border-slate-200 p-4 rounded-xl text-right min-w-[220px]">
                <p class="text-xs font-mono font-bold text-slate-800 text-sm">{{ workOrder?.wo_number }}</p>
                <p class="text-xs text-slate-500 mt-1">Tanggal: {{ formatDate(workOrder?.created_at) }}</p>
                <span 
                  :class="workOrder?.payment_status === 'paid' ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'"
                  class="inline-block mt-2 px-2.5 py-0.5 rounded text-[11px] font-bold uppercase tracking-wider">
                  {{ workOrder?.payment_status === 'paid' ? 'LUNAS (PAID)' : 'MENUNGGU KASIR' }}
                </span>
              </div>
            </div>

            <!-- Customer & Motor Information Grid -->
            <div class="grid grid-cols-2 gap-4 bg-slate-50 p-4 rounded-xl border border-slate-200 text-xs mb-6">
              <div>
                <h4 class="font-bold text-slate-700 uppercase tracking-wider mb-2 text-[10px]">Data Konsumen & Kendaraan</h4>
                <div class="space-y-1">
                  <p><span class="text-slate-500 inline-block w-24">Pemilik:</span> <span class="font-semibold">{{ workOrder?.customer?.name }}</span></p>
                  <p><span class="text-slate-500 inline-block w-24">No. HP:</span> {{ workOrder?.customer?.phone }}</p>
                  <p><span class="text-slate-500 inline-block w-24">No. Polisi Plat:</span> <span class="font-mono font-bold bg-white px-1.5 py-0.5 rounded border border-slate-300 text-slate-900">{{ workOrder?.customer?.plate_number }}</span></p>
                  <p><span class="text-slate-500 inline-block w-24">Tipe Motor:</span> {{ workOrder?.customer?.motor_model }} ({{ workOrder?.customer?.year }})</p>
                  <p><span class="text-slate-500 inline-block w-24">Odometer:</span> <span class="font-mono font-semibold">{{ workOrder?.entry_km }} KM</span></p>
                </div>
              </div>

              <div>
                <h4 class="font-bold text-slate-700 uppercase tracking-wider mb-2 text-[10px]">Penanggung Jawab Bengkel</h4>
                <div class="space-y-1">
                  <p><span class="text-slate-500 inline-block w-28">Service Advisor:</span> <span class="font-semibold">{{ workOrder?.service_advisor?.full_name || 'Service Advisor AHASS' }}</span></p>
                  <p><span class="text-slate-500 inline-block w-28">Mekanik Handal:</span> <span class="font-semibold">{{ workOrder?.mechanic?.name }} (Pit {{ workOrder?.mechanic?.pit_number }})</span></p>
                  <p><span class="text-slate-500 inline-block w-28">Metode Bayar:</span> <span class="uppercase font-mono font-medium">{{ workOrder?.payment_method }}</span></p>
                  <p><span class="text-slate-500 inline-block w-28">Keluhan:</span> <span class="text-slate-700 italic">"{{ workOrder?.complaints || '-' }}"</span></p>
                </div>
              </div>
            </div>

            <!-- Items Table -->
            <table class="w-full text-xs text-left border-collapse mb-6">
              <thead>
                <tr class="bg-red-600 text-white font-bold uppercase text-[10px] tracking-wider">
                  <th class="py-2.5 px-3 rounded-l">No</th>
                  <th class="py-2.5 px-3">Kode Part/Jasa</th>
                  <th class="py-2.5 px-3">Deskripsi Pekerjaan & Suku Cadang</th>
                  <th class="py-2.5 px-3 text-center">Qty</th>
                  <th class="py-2.5 px-3 text-right">Harga Satuan</th>
                  <th class="py-2.5 px-3 text-center">Diskon</th>
                  <th class="py-2.5 px-3 text-right rounded-r">Subtotal</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-200">
                <tr v-for="(item, idx) in workOrder?.items" :key="idx" class="hover:bg-slate-50/50">
                  <td class="py-2 px-3 text-slate-400 font-mono">{{ idx + 1 }}</td>
                  <td class="py-2 px-3 font-mono font-medium text-slate-600">{{ item.part_number_snapshot }}</td>
                  <td class="py-2 px-3 font-medium text-slate-800">
                    <span v-if="item.item_type === 'service'" class="text-blue-600 font-bold mr-1">[JASA]</span>
                    {{ item.item_name_snapshot }}
                  </td>
                  <td class="py-2 px-3 text-center font-mono font-semibold">{{ item.quantity }}</td>
                  <td class="py-2 px-3 text-right font-mono">{{ formatRupiah(item.price_snapshot) }}</td>
                  <td class="py-2 px-3 text-center font-mono">{{ Number(item.discount_percent) > 0 ? `${item.discount_percent}%` : '-' }}</td>
                  <td class="py-2 px-3 text-right font-mono font-bold text-slate-900">{{ formatRupiah(item.subtotal) }}</td>
                </tr>
              </tbody>
            </table>

            <!-- Financial Calculations Breakdown -->
            <div class="flex justify-between items-start border-t border-slate-200 pt-4">
              <div class="w-1/2 text-[11px] text-slate-500 space-y-1 pr-6">
                <p class="font-bold text-slate-700">Catatan & Garansi Servis AHASS:</p>
                <p>1. Garansi servis berlaku selama 500 KM atau 7 hari kerja sejak tanggal pengerjaan.</p>
                <p>2. Seluruh suku cadang yang digunakan adalah 100% Asli Honda Genuine Parts (HGP) berstandar PT Astra Honda Motor.</p>
              </div>

              <div class="w-1/2 space-y-1.5 text-xs">
                <div class="flex justify-between text-slate-600">
                  <span>Subtotal Suku Cadang:</span>
                  <span class="font-mono">{{ formatRupiah(workOrder?.parts_subtotal) }}</span>
                </div>
                <div class="flex justify-between text-slate-600">
                  <span>Subtotal Jasa Mekanik:</span>
                  <span class="font-mono">{{ formatRupiah(workOrder?.services_subtotal) }}</span>
                </div>
                <div v-if="Number(workOrder?.total_discount) > 0" class="flex justify-between text-emerald-600 font-medium">
                  <span>Total Diskon Promosi:</span>
                  <span class="font-mono">-{{ formatRupiah(workOrder?.total_discount) }}</span>
                </div>
                <div class="flex justify-between text-slate-600">
                  <span>PPN (11%):</span>
                  <span class="font-mono">{{ formatRupiah(workOrder?.tax_amount) }}</span>
                </div>
                <div class="flex justify-between font-black text-sm text-slate-900 pt-2 border-t border-slate-300">
                  <span class="text-red-600">GRAND TOTAL:</span>
                  <span class="font-mono">{{ formatRupiah(workOrder?.grand_total) }}</span>
                </div>
                <div v-if="workOrder?.payment_status === 'paid'" class="flex justify-between text-slate-500 text-[11px] pt-1">
                  <span>Jumlah Bayar / Kembalian:</span>
                  <span class="font-mono">{{ formatRupiah(workOrder?.paid_amount) }} / {{ formatRupiah(workOrder?.change_amount) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Dual Signature Section & Digital Stamp Seal -->
          <div class="mt-12 pt-6 border-t border-slate-200 flex justify-between items-center text-xs text-center">
            <div class="w-48">
              <p class="text-slate-600">Tanda Tangan Pelanggan,</p>
              <div class="h-16"></div>
              <p class="font-bold border-t border-slate-300 pt-1 text-slate-800">( {{ workOrder?.customer?.name }} )</p>
            </div>

            <!-- Official Seal Stamp -->
            <div class="w-24 h-24 border-2 border-dashed border-red-500 rounded-full flex flex-col items-center justify-center p-2 text-red-600 opacity-80 transform -rotate-6">
              <span class="text-[9px] font-black uppercase">AHASS 01234</span>
              <span class="text-[8px] font-bold text-center leading-tight">NUSANTARA SAKTI</span>
              <span class="text-[7px] font-mono mt-0.5">VERIFIED</span>
            </div>

            <div class="w-48">
              <p class="text-slate-600">Service Advisor / Kasir AHASS,</p>
              <div class="h-16"></div>
              <p class="font-bold border-t border-slate-300 pt-1 text-slate-800">( {{ workOrder?.service_advisor?.full_name || 'Petugas Bengkel' }} )</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { FileText, Download, X, Printer } from 'lucide-vue-next'
import { formatRupiah, formatDate } from '../utils/format'
import { printA4Document } from '../utils/printer'
import type { WorkOrder } from '../types'

const props = defineProps<{
  isOpen: boolean
  workOrder: WorkOrder | null
}>()

defineEmits(['close'])

const invoiceSheetRef = ref<HTMLElement | null>(null)

function printDirect() {
  if (!invoiceSheetRef.value) return
  printA4Document(invoiceSheetRef.value, `Faktur Servis AHASS - ${props.workOrder?.wo_number || ''}`)
}

function downloadPDF() {
  if (!props.workOrder) return
  const token = localStorage.getItem('motoflow_token')
  const url = `/api/v1/work-orders/${props.workOrder.id}/invoice-pdf`
  
  // Open with fetch & blob or direct link
  fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`
    }
  })
    .then(res => res.blob())
    .then(blob => {
      const blobUrl = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = blobUrl
      a.download = `Faktur-AHASS-${props.workOrder?.wo_number}.pdf`
      document.body.appendChild(a)
      a.click()
      a.remove()
    })
    .catch(err => {
      console.error('Download PDF error:', err)
      alert('Gagal mengunduh file PDF faktur')
    })
}
</script>

