<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <div class="flex items-center space-x-2">
          <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold uppercase tracking-wider bg-red-100 text-red-700">
            Dealer Resmi Honda & Pembiayaan NSC Finance
          </span>
          <span class="text-xs text-slate-500 font-mono">PT Nusantara Sakti Group</span>
        </div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight mt-1 flex items-center">
          <Bike class="w-6 h-6 mr-2 text-red-600" />
          Showroom Penjualan Motor Honda & Simulasi Kredit NSC
        </h1>
        <p class="text-xs text-slate-500">
          Katalog fisik unit motor Honda baru, alokasi Nomor Rangka (VIN), dan simulasi kredit cicilan NSC Finance.
        </p>
      </div>

      <!-- Navigation Tabs -->
      <div class="flex items-center space-x-2 bg-slate-100 p-1 rounded-xl">
        <button
          @click="activeTab = 'catalog'"
          :class="[
            'px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all',
            activeTab === 'catalog'
              ? 'bg-white text-slate-900 shadow-xs'
              : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          Katalog Showroom ({{ showroomStore.units.length }})
        </button>
        <button
          @click="activeTab = 'spk'"
          :class="[
            'px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all',
            activeTab === 'spk'
              ? 'bg-white text-slate-900 shadow-xs'
              : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          Daftar SPK Terbit ({{ showroomStore.spkList.length }})
        </button>
      </div>
    </div>

    <!-- TAB 1: SHOWROOM CATALOG -->
    <div v-if="activeTab === 'catalog'" class="space-y-4">
      <!-- Search & Filters -->
      <div class="bg-white p-4 rounded-2xl border border-slate-200/80 shadow-xs flex flex-col md:flex-row gap-3 items-center justify-between">
        <div class="relative w-full md:w-96">
          <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari model motor, nomor rangka, warna..."
            class="w-full bg-slate-50 border border-slate-200 rounded-xl pl-10 pr-4 py-2 text-xs text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-red-500"
          />
        </div>

        <div class="flex items-center space-x-1.5 w-full md:w-auto overflow-x-auto text-xs">
          <button
            v-for="cat in ['Semua', 'Matic', 'Sport', 'EV Listrik']"
            :key="cat"
            @click="selectedCategory = cat"
            :class="[
              'px-3 py-1.5 rounded-lg font-semibold whitespace-nowrap transition-all',
              selectedCategory === cat
                ? 'bg-red-600 text-white shadow-xs'
                : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
            ]"
          >
            {{ cat }}
          </button>
        </div>
      </div>

      <!-- Motor Unit Cards Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        <div
          v-for="unit in filteredUnits"
          :key="unit.id"
          class="bg-white rounded-2xl border border-slate-200/80 hover:border-red-300 hover:shadow-lg transition-all overflow-hidden flex flex-col justify-between group"
        >
          <div>
            <!-- Image & Badges -->
            <div class="relative h-44 bg-slate-100 overflow-hidden">
              <img
                :src="unit.image_url"
                :alt="unit.model_name"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              />
              <div class="absolute top-3 left-3 flex items-center space-x-1.5">
                <span class="px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider bg-slate-900/80 backdrop-blur-md text-white">
                  {{ unit.category }}
                </span>
                <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-white/90 backdrop-blur-md text-slate-800">
                  {{ unit.year }}
                </span>
              </div>
              <div class="absolute top-3 right-3">
                <span
                  :class="[
                    'px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider',
                    unit.status === 'ready'
                      ? 'bg-emerald-500 text-white'
                      : 'bg-amber-500 text-white'
                  ]"
                >
                  {{ unit.status === 'ready' ? 'Ready Showroom' : 'Booked SPK' }}
                </span>
              </div>
            </div>

            <!-- Details -->
            <div class="p-4 space-y-3">
              <div>
                <h3 class="font-black text-base text-slate-900 group-hover:text-red-600 transition-colors">
                  {{ unit.model_name }}
                </h3>
                <p class="text-xs text-slate-500 font-medium mt-0.5 flex items-center">
                  <Palette class="w-3.5 h-3.5 mr-1 text-slate-400" />
                  Warna: <span class="font-bold text-slate-700 ml-1">{{ unit.color }}</span>
                </p>
              </div>

              <div class="bg-slate-50 p-2.5 rounded-xl border border-slate-100 space-y-1 text-[11px] font-mono">
                <div class="flex justify-between">
                  <span class="text-slate-400">No. Rangka:</span>
                  <span class="font-bold text-slate-800">{{ unit.frame_number }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-slate-400">No. Mesin:</span>
                  <span class="font-bold text-slate-800">{{ unit.engine_number }}</span>
                </div>
              </div>

              <p class="text-xs text-slate-500 line-clamp-2 leading-relaxed">
                {{ unit.specs_summary }}
              </p>

              <!-- Price Box -->
              <div class="pt-2 border-t border-slate-100 flex items-baseline justify-between">
                <div>
                  <span class="text-[10px] text-slate-400 font-bold block uppercase">Harga OTR Jakarta</span>
                  <span class="text-lg font-black text-slate-900 font-mono">
                    {{ formatRupiah(unit.otr_price) }}
                  </span>
                </div>
                <div class="text-right">
                  <span class="text-[10px] text-slate-400 font-bold block uppercase">DP Mulai</span>
                  <span class="text-xs font-bold text-red-600 font-mono">
                    {{ formatRupiah(unit.min_dp) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="p-4 pt-0 grid grid-cols-2 gap-2">
            <button
              @click="openCreditSimulator(unit)"
              class="w-full py-2.5 px-3 text-xs font-bold text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-xl transition-all flex items-center justify-center space-x-1"
            >
              <Calculator class="w-3.5 h-3.5 text-slate-600" />
              <span>Simulasi NSC</span>
            </button>
            <button
              @click="openSPKModal(unit)"
              :disabled="unit.status !== 'ready'"
              class="w-full py-2.5 px-3 text-xs font-bold text-white bg-red-600 hover:bg-red-700 disabled:bg-slate-300 disabled:cursor-not-allowed rounded-xl transition-all shadow-xs flex items-center justify-center space-x-1"
            >
              <FileText class="w-3.5 h-3.5" />
              <span>Buat SPK</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: SPK LIST -->
    <div v-if="activeTab === 'spk'" class="space-y-4">
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
        <div class="p-4 border-b border-slate-100 flex justify-between items-center">
          <h3 class="font-bold text-sm text-slate-900">Daftar Surat Pemesanan Kendaraan (SPK)</h3>
          <span class="text-xs text-slate-500 font-mono">{{ showroomStore.spkList.length }} SPK terdaftar</span>
        </div>

        <div v-if="showroomStore.spkList.length === 0" class="p-10 text-center text-slate-400 text-xs">
          Belum ada SPK yang diterbitkan. Silakan buat SPK dari tab Katalog Showroom.
        </div>

        <div v-else class="divide-y divide-slate-100">
          <div
            v-for="spk in showroomStore.spkList"
            :key="spk.id"
            class="p-4 flex flex-col md:flex-row md:items-center justify-between gap-3 text-xs"
          >
            <div class="space-y-1">
              <div class="flex items-center space-x-2">
                <span class="font-mono font-bold text-red-600">{{ spk.spk_number }}</span>
                <span class="text-slate-400">•</span>
                <span class="font-bold text-slate-900">{{ spk.customer?.name || 'Pelanggan' }}</span>
                <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider bg-slate-100 text-slate-700">
                  {{ spk.purchase_type === 'credit_nsc' ? 'Kredit NSC Finance' : 'Cash Keras' }}
                </span>
              </div>
              <div class="text-slate-600 font-medium">
                {{ spk.motor_unit?.model_name }} ({{ spk.motor_unit?.color }})
                <span class="text-slate-400 font-mono ml-2">VIN: {{ spk.motor_unit?.frame_number }}</span>
              </div>
              <div class="text-[11px] text-slate-400">
                Diterbitkan: {{ formatDate(spk.created_at) }} • Sales: {{ spk.sales?.full_name || 'Counter NSS' }}
              </div>
            </div>

            <div class="flex items-center space-x-4">
              <div v-if="spk.purchase_type === 'credit_nsc'" class="text-right">
                <span class="text-[10px] text-slate-400 block font-bold">Angsuran / Bulan</span>
                <span class="font-mono font-black text-red-600 text-sm">
                  {{ formatRupiah(spk.monthly_installment) }} <span class="text-[10px] text-slate-500 font-normal">({{ spk.tenor_months }} bln)</span>
                </span>
              </div>
              <div v-else class="text-right">
                <span class="text-[10px] text-slate-400 block font-bold">Tanda Jadi</span>
                <span class="font-mono font-black text-slate-900 text-sm">
                  {{ formatRupiah(spk.booking_fee) }}
                </span>
              </div>

              <button
                @click="openSPKReceipt(spk.id)"
                class="px-3 py-1.5 text-xs font-bold text-red-600 bg-red-50 hover:bg-red-100 rounded-xl border border-red-200 flex items-center space-x-1"
              >
                <Printer class="w-3.5 h-3.5" />
                <span>Dokumen</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL 1: INTERACTIVE CREDIT SIMULATOR NSC FINANCE -->
    <div v-if="showSimulatorModal && activeUnit" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full p-6 border border-slate-200 space-y-5">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center space-x-2">
            <Calculator class="w-5 h-5 text-red-600" />
            <h3 class="font-bold text-slate-900 text-sm">Kalkulator Kredit NSC Finance</h3>
          </div>
          <button @click="showSimulatorModal = false" class="text-slate-400 hover:text-slate-600">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Selected Unit Header -->
        <div class="p-3 bg-slate-50 rounded-xl border border-slate-100 flex items-center justify-between">
          <div>
            <h4 class="font-bold text-xs text-slate-900">{{ activeUnit.model_name }}</h4>
            <span class="text-[11px] text-slate-500">Harga OTR Resmi</span>
          </div>
          <span class="font-mono font-black text-slate-900 text-base">
            {{ formatRupiah(activeUnit.otr_price) }}
          </span>
        </div>

        <!-- Slider Uang Muka (DP) -->
        <div class="space-y-2">
          <div class="flex justify-between items-center text-xs">
            <label class="font-bold text-slate-700">Uang Muka (DP): {{ simDPPercent }}%</label>
            <span class="font-mono font-black text-red-600">{{ formatRupiah(simDPAmount) }}</span>
          </div>
          <input
            v-model.number="simDPPercent"
            @input="runSimulation"
            type="range"
            min="10"
            max="50"
            step="5"
            class="w-full accent-red-600 cursor-pointer"
          />
          <div class="flex justify-between text-[10px] text-slate-400 font-mono">
            <span>Min 10% (OJK)</span>
            <span>20%</span>
            <span>30%</span>
            <span>40%</span>
            <span>Max 50%</span>
          </div>
        </div>

        <!-- Tenor Pills -->
        <div class="space-y-2">
          <label class="block text-xs font-bold text-slate-700">Pilihan Jangka Waktu (Tenor)</label>
          <div class="grid grid-cols-5 gap-1.5 text-xs font-bold">
            <button
              v-for="t in [11, 17, 23, 29, 35]"
              :key="t"
              @click="simTenor = t; runSimulation()"
              :class="[
                'py-2 rounded-xl border text-center transition-all',
                simTenor === t
                  ? 'bg-red-600 text-white border-red-600 shadow-xs'
                  : 'bg-slate-50 text-slate-700 border-slate-200 hover:bg-slate-100'
              ]"
            >
              {{ t }} Bln
            </button>
          </div>
        </div>

        <!-- Calculation Results Card -->
        <div v-if="showroomStore.simulationResult" class="bg-red-50/50 p-4 rounded-xl border border-red-100 space-y-2 text-xs">
          <div class="flex justify-between text-slate-600">
            <span>Pokok Pinjaman</span>
            <span class="font-mono font-bold">{{ formatRupiah(showroomStore.simulationResult.principal_loan) }}</span>
          </div>
          <div class="flex justify-between text-slate-600">
            <span>Bunga Pembiayaan NSC (18% Flat p.a.)</span>
            <span class="font-mono font-bold">{{ formatRupiah(showroomStore.simulationResult.total_interest) }}</span>
          </div>
          <div class="flex justify-between text-slate-600">
            <span>Biaya Fidusia & Asuransi Kehilangan (TLO)</span>
            <span class="font-mono font-bold">
              {{ formatRupiah(Number(showroomStore.simulationResult.admin_fee) + Number(showroomStore.simulationResult.insurance_tlo)) }}
            </span>
          </div>

          <div class="pt-3 border-t border-red-200 flex items-center justify-between">
            <div>
              <span class="text-[10px] text-slate-500 font-bold uppercase block">Cicilan Per Bulan</span>
              <span class="text-xl font-black text-red-600 font-mono">
                {{ formatRupiah(showroomStore.simulationResult.monthly_installment) }}
              </span>
              <span class="text-[10px] text-slate-400 font-medium block">selama {{ simTenor }} bulan</span>
            </div>
            <div class="text-right">
              <span class="text-[10px] text-slate-500 font-bold uppercase block">Total Bayar Pertama (TDP)</span>
              <span class="text-xs font-black text-slate-900 font-mono">
                {{ formatRupiah(showroomStore.simulationResult.total_first_payment) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Action -->
        <div class="flex space-x-2 pt-2">
          <button
            @click="showSimulatorModal = false"
            class="flex-1 py-2.5 text-xs font-bold text-slate-600 bg-slate-100 hover:bg-slate-200 rounded-xl"
          >
            Tutup
          </button>
          <button
            @click="proceedToSPKFromSimulation"
            class="flex-1 py-2.5 text-xs font-bold text-white bg-red-600 hover:bg-red-700 rounded-xl shadow-xs flex items-center justify-center space-x-1"
          >
            <span>Ajukan SPK Cicilan Ini</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL 2: PENERBITAN SPK RESMI -->
    <div v-if="showSPKModal && activeUnit" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full p-6 border border-slate-200 space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center space-x-2">
            <FileText class="w-5 h-5 text-red-600" />
            <h3 class="font-bold text-slate-900 text-sm">Penerbitan Surat Pemesanan Kendaraan (SPK)</h3>
          </div>
          <button @click="showSPKModal = false" class="text-slate-400 hover:text-slate-600">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Unit Summary -->
        <div class="p-3 bg-slate-50 rounded-xl border border-slate-100 text-xs flex justify-between items-center">
          <div>
            <h4 class="font-bold text-slate-900">{{ activeUnit.model_name }}</h4>
            <span class="text-slate-500 font-mono text-[11px]">VIN: {{ activeUnit.frame_number }}</span>
          </div>
          <span class="font-mono font-bold text-slate-900">{{ formatRupiah(activeUnit.otr_price) }}</span>
        </div>

        <!-- Form SPK -->
        <form @submit.prevent="handleCreateSPK" class="space-y-3 text-xs">
          <!-- Customer Selection -->
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Pilih Calon Pembeli *</label>
            <select
              v-model.number="spkForm.customer_id"
              required
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2 text-xs"
            >
              <option disabled value="0">-- Pilih Pelanggan Terdaftar --</option>
              <option v-for="c in customerStore.customers" :key="c.id" :value="c.id">
                {{ c.name }} - {{ c.phone }} ({{ c.plate_number }})
              </option>
            </select>
          </div>

          <!-- Purchase Type -->
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Skema Pembelian *</label>
            <div class="grid grid-cols-2 gap-2">
              <button
                type="button"
                @click="spkForm.purchase_type = 'cash'"
                :class="[
                  'py-2 rounded-xl border text-center font-bold transition-all',
                  spkForm.purchase_type === 'cash'
                    ? 'bg-red-600 text-white border-red-600'
                    : 'bg-slate-50 text-slate-700 border-slate-200'
                ]"
              >
                Cash Keras (Tunai)
              </button>
              <button
                type="button"
                @click="spkForm.purchase_type = 'credit_nsc'"
                :class="[
                  'py-2 rounded-xl border text-center font-bold transition-all',
                  spkForm.purchase_type === 'credit_nsc'
                    ? 'bg-red-600 text-white border-red-600'
                    : 'bg-slate-50 text-slate-700 border-slate-200'
                ]"
              >
                Kredit NSC Finance
              </button>
            </div>
          </div>

          <!-- If Cash: Booking Fee -->
          <div v-if="spkForm.purchase_type === 'cash'">
            <label class="block font-semibold text-slate-700 mb-1">Tanda Jadi / Booking Fee (Rp)</label>
            <input
              v-model.number="spkForm.booking_fee"
              type="number"
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2 font-mono font-bold"
            />
          </div>

          <!-- If Credit: DP, Tenor & Cicilan -->
          <div v-else class="space-y-3 bg-red-50/50 p-3 rounded-xl border border-red-100">
            <div class="grid grid-cols-2 gap-2">
              <div>
                <label class="block font-semibold text-slate-700 mb-1">Uang Muka DP (Rp)</label>
                <input
                  v-model.number="spkForm.down_payment"
                  type="number"
                  class="w-full bg-white border border-slate-200 rounded-lg px-2.5 py-1.5 font-mono font-bold"
                />
              </div>
              <div>
                <label class="block font-semibold text-slate-700 mb-1">Tenor Cicilan</label>
                <select
                  v-model.number="spkForm.tenor_months"
                  class="w-full bg-white border border-slate-200 rounded-lg px-2.5 py-1.5 font-bold"
                >
                  <option :value="11">11 Bulan</option>
                  <option :value="17">17 Bulan</option>
                  <option :value="23">23 Bulan</option>
                  <option :value="29">29 Bulan</option>
                  <option :value="35">35 Bulan</option>
                </select>
              </div>
            </div>

            <div>
              <label class="block font-semibold text-slate-700 mb-1">Angsuran Per Bulan (Estimasi NSC)</label>
              <input
                v-model.number="spkForm.monthly_installment"
                type="number"
                class="w-full bg-white border border-slate-200 rounded-lg px-2.5 py-1.5 font-mono font-black text-red-600"
              />
            </div>
          </div>

          <div>
            <label class="block font-semibold text-slate-700 mb-1">Catatan Khusus (Alamat Pengiriman / STNK)</label>
            <textarea
              v-model="spkForm.notes"
              rows="2"
              placeholder="Contoh: Pembuatan STNK atas nama istri, plat nomor B..."
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2 text-xs"
            ></textarea>
          </div>

          <div class="flex space-x-2 pt-2 border-t border-slate-100">
            <button
              type="button"
              @click="showSPKModal = false"
              class="flex-1 py-2.5 text-xs font-bold text-slate-600 bg-slate-100 hover:bg-slate-200 rounded-xl"
            >
              Batal
            </button>
            <button
              type="submit"
              class="flex-1 py-2.5 text-xs font-bold text-white bg-red-600 hover:bg-red-700 rounded-xl shadow-xs"
            >
              Terbitkan SPK Resmi
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL 3: SPK RECEIPT / DOKUMEN -->
    <div v-if="showReceiptModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full p-6 border border-slate-200 space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center space-x-2">
            <Printer class="w-5 h-5 text-red-600" />
            <h3 class="font-bold text-slate-900 text-sm">Dokumen SPK PT Nusantara Sakti Group</h3>
          </div>
          <button @click="showReceiptModal = false" class="text-slate-400 hover:text-slate-600">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="bg-slate-50 p-4 rounded-xl border border-slate-200 font-mono text-[11px] whitespace-pre-wrap leading-tight text-slate-800 max-h-[380px] overflow-y-auto">
          {{ receiptText }}
        </div>

        <div class="flex space-x-2 pt-2">
          <button
            @click="showReceiptModal = false"
            class="flex-1 py-2.5 text-xs font-bold text-slate-600 bg-slate-100 hover:bg-slate-200 rounded-xl"
          >
            Tutup
          </button>
          <button
            @click="printReceipt"
            class="flex-1 py-2.5 text-xs font-bold text-white bg-red-600 hover:bg-red-700 rounded-xl shadow-xs flex items-center justify-center space-x-1.5"
          >
            <Printer class="w-4 h-4" />
            <span>Cetak SPK</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Bike,
  Search,
  Calculator,
  FileText,
  Palette,
  Printer,
  X,
} from 'lucide-vue-next'
import { useShowroomStore } from '../stores/showroom'
import { useCustomerStore } from '../stores/customer'
import { useToastStore } from '../stores/toast'
import { formatRupiah, formatDate } from '../utils/format'
import { playSuccessChime, playClickSound } from '../utils/sound'
import type { MotorUnit } from '../types/showroom'

const showroomStore = useShowroomStore()
const customerStore = useCustomerStore()
const toast = useToastStore()

const activeTab = ref<'catalog' | 'spk'>('catalog')
const searchQuery = ref('')
const selectedCategory = ref('Semua')

const showSimulatorModal = ref(false)
const showSPKModal = ref(false)
const showReceiptModal = ref(false)
const receiptText = ref('')

const activeUnit = ref<MotorUnit | null>(null)
const simDPPercent = ref(20)
const simTenor = ref(35)

const spkForm = ref({
  customer_id: 0,
  purchase_type: 'credit_nsc' as 'cash' | 'credit_nsc',
  booking_fee: 500000,
  down_payment: 3000000,
  tenor_months: 35,
  monthly_installment: 1045000,
  notes: '',
})

onMounted(async () => {
  await Promise.all([
    showroomStore.fetchUnits(),
    showroomStore.fetchSPKList(),
    customerStore.fetchCustomers(),
  ])
})

const filteredUnits = computed(() => {
  return showroomStore.units.filter(u => {
    const matchesSearch =
      u.model_name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      u.frame_number.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      u.color.toLowerCase().includes(searchQuery.value.toLowerCase())

    if (selectedCategory.value === 'Semua') return matchesSearch
    return matchesSearch && u.category.toLowerCase().includes(selectedCategory.value.toLowerCase())
  })
})

const simDPAmount = computed(() => {
  if (!activeUnit.value) return 0
  return (Number(activeUnit.value.otr_price) * simDPPercent.value) / 100
})

async function openCreditSimulator(unit: MotorUnit) {
  playClickSound()
  activeUnit.value = unit
  simDPPercent.value = 20
  simTenor.value = 35
  await runSimulation()
  showSimulatorModal.value = true
}

async function runSimulation() {
  if (!activeUnit.value) return
  await showroomStore.simulateCredit(
    Number(activeUnit.value.otr_price),
    simDPPercent.value,
    simTenor.value
  )
}

function proceedToSPKFromSimulation() {
  showSimulatorModal.value = false
  if (!activeUnit.value) return
  openSPKModal(activeUnit.value)
  if (showroomStore.simulationResult) {
    spkForm.value.purchase_type = 'credit_nsc'
    spkForm.value.down_payment = Number(showroomStore.simulationResult.dp_amount)
    spkForm.value.tenor_months = showroomStore.simulationResult.tenor_months
    spkForm.value.monthly_installment = Number(showroomStore.simulationResult.monthly_installment)
  }
}

function openSPKModal(unit: MotorUnit) {
  playClickSound()
  activeUnit.value = unit
  spkForm.value.booking_fee = 500000
  spkForm.value.down_payment = Number(unit.min_dp)
  if (customerStore.customers.length > 0 && spkForm.value.customer_id === 0) {
    spkForm.value.customer_id = customerStore.customers[0].id
  }
  showSPKModal.value = true
}

async function handleCreateSPK() {
  if (!activeUnit.value) return
  if (spkForm.value.customer_id === 0) {
    toast.error('Silakan pilih calon pembeli')
    return
  }

  try {
    const spk = await showroomStore.createSPK({
      customer_id: spkForm.value.customer_id,
      motor_unit_id: activeUnit.value.id,
      purchase_type: spkForm.value.purchase_type,
      booking_fee: spkForm.value.booking_fee,
      down_payment: spkForm.value.down_payment,
      tenor_months: spkForm.value.tenor_months,
      monthly_installment: spkForm.value.monthly_installment,
      notes: spkForm.value.notes,
    })

    playSuccessChime()
    toast.success(`SPK ${spk.spk_number} berhasil diterbitkan!`)
    showSPKModal.value = false

    // Load receipt
    const rData = await showroomStore.fetchSPKReceipt(spk.id)
    receiptText.value = rData.receipt_text
    showReceiptModal.value = true
  } catch (err: any) {
    toast.error(err.message || 'Gagal menerbitkan SPK')
  }
}

async function openSPKReceipt(spkId: number) {
  try {
    const rData = await showroomStore.fetchSPKReceipt(spkId)
    receiptText.value = rData.receipt_text
    showReceiptModal.value = true
  } catch (err: any) {
    toast.error('Gagal mengambil dokumen SPK')
  }
}

function printReceipt() {
  window.print()
}
</script>

