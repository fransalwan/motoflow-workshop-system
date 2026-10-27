<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 bg-slate-950/70 backdrop-blur-sm flex items-center justify-center p-4">
    <div class="bg-white rounded-3xl shadow-2xl max-w-sm w-full overflow-hidden border border-slate-200 text-center animate-in fade-in zoom-in-95 duration-200">
      <!-- QRIS National Header Banner -->
      <div class="bg-gradient-to-r from-red-700 to-red-600 text-white p-4 relative">
        <button
          @click="$emit('close')"
          class="absolute top-4 right-4 text-white/80 hover:text-white transition">
          <X class="w-5 h-5" />
        </button>
        <div class="inline-flex items-center space-x-1.5 bg-white px-3 py-1 rounded-full text-slate-900 font-black text-xs shadow-sm mb-1">
          <span class="text-red-600">QRIS</span>
          <span class="text-[9px] text-slate-500 font-normal">| GPN</span>
        </div>
        <h3 class="font-bold text-sm">Pembayaran Digital AHASS NSS</h3>
        <p class="text-[11px] text-red-100">PT Nusantara Sakti Group - Cempaka Putih</p>
      </div>

      <!-- QR Code & Nominal Display -->
      <div class="p-6 space-y-4">
        <div>
          <span class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Total Tagihan Servis</span>
          <p class="text-2xl font-black text-slate-900 font-mono mt-0.5">{{ formatRupiah(amount) }}</p>
        </div>

        <!-- Simulated Scannable QR Frame -->
        <div class="p-4 bg-slate-50 rounded-2xl border-2 border-dashed border-slate-300 inline-block relative mx-auto group">
          <div class="w-48 h-48 bg-white p-2 rounded-xl shadow-inner flex items-center justify-center relative">
            <!-- Simulated QR Pattern Graphic -->
            <svg class="w-full h-full text-slate-900" viewBox="0 0 100 100" fill="currentColor">
              <!-- Position detection patterns (corners) -->
              <rect x="5" y="5" width="25" height="25" rx="3" fill="#0F172A" />
              <rect x="9" y="9" width="17" height="17" fill="#FFFFFF" />
              <rect x="13" y="13" width="9" height="9" fill="#CC0000" />

              <rect x="70" y="5" width="25" height="25" rx="3" fill="#0F172A" />
              <rect x="74" y="9" width="17" height="17" fill="#FFFFFF" />
              <rect x="78" y="13" width="9" height="9" fill="#CC0000" />

              <rect x="5" y="70" width="25" height="25" rx="3" fill="#0F172A" />
              <rect x="9" y="74" width="17" height="17" fill="#FFFFFF" />
              <rect x="13" y="78" width="9" height="9" fill="#CC0000" />

              <!-- Data matrix dots -->
              <rect x="36" y="8" width="6" height="6" />
              <rect x="48" y="14" width="6" height="6" />
              <rect x="58" y="8" width="6" height="6" />
              <rect x="36" y="24" width="6" height="6" />
              <rect x="50" y="24" width="6" height="6" />
              
              <!-- Center AHASS Emblem -->
              <rect x="40" y="40" width="20" height="20" rx="4" fill="#CC0000" />
              <text x="50" y="54" font-size="11" font-weight="900" fill="#FFFFFF" text-anchor="middle">H</text>

              <!-- Bottom Matrix Dots -->
              <rect x="36" y="72" width="6" height="6" />
              <rect x="50" y="72" width="6" height="6" />
              <rect x="62" y="72" width="6" height="6" />
              <rect x="40" y="84" width="6" height="6" />
              <rect x="54" y="84" width="6" height="6" />
              <rect x="70" y="84" width="6" height="6" />
              <rect x="84" y="84" width="6" height="6" />
              <rect x="72" y="42" width="6" height="6" />
              <rect x="84" y="52" width="6" height="6" />
              <rect x="80" y="62" width="6" height="6" />
            </svg>
          </div>
          <p class="text-[10px] text-slate-500 font-mono mt-2">NMID: ID1026098234812</p>
        </div>

        <!-- Countdown Timer -->
        <div class="flex items-center justify-center space-x-1.5 text-xs text-amber-600 bg-amber-50 py-1.5 px-3 rounded-xl border border-amber-200">
          <Clock class="w-3.5 h-3.5" />
          <span>Berlaku hingga: <strong>02:59</strong></span>
        </div>

        <p class="text-[11px] text-slate-400">
          Mendukung BCA, Mandiri, BRI, BNI, GoPay, OVO, ShopeePay & Dana
        </p>

        <!-- Instant Verification Button -->
        <button
          @click="$emit('verified')"
          class="w-full bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs py-2.5 rounded-xl shadow-md transition flex items-center justify-center space-x-1.5">
          <CheckCircle2 class="w-4 h-4" />
          <span>Simulasikan Pembayaran Berhasil</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { X, Clock, CheckCircle2 } from 'lucide-vue-next'
import { formatRupiah } from '../utils/format'

defineProps<{
  isOpen: boolean
  amount: number | string
}>()

defineEmits(['close', 'verified'])
</script>

