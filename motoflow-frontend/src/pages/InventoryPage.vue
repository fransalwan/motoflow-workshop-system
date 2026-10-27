<template>
  <div class="space-y-6">
    <!-- Header Controls -->
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

      <!-- Search & Add Button -->
      <div class="flex items-center space-x-3">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari part number atau nama part..."
            class="bg-white border border-slate-200 rounded-xl pl-9 pr-4 py-2 text-xs text-slate-900 placeholder-slate-400 focus:outline-none focus:border-red-500 w-72 shadow-xs" />
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
        </div>

        <button
          @click="showAddModal = true"
          class="inline-flex items-center space-x-1.5 bg-red-600 hover:bg-red-700 text-white font-semibold text-xs px-4 py-2 rounded-xl shadow-xs transition">
          <Plus class="w-4 h-4" />
          <span>Tambah Part</span>
        </button>
      </div>
    </div>

    <!-- Inventory Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <!-- Desktop Table -->
      <div class="hidden md:block overflow-x-auto">
        <table class="w-full text-xs text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-500 font-bold uppercase text-[10px] tracking-wider">
              <th class="py-3 px-4">Nomor Part (HGP)</th>
              <th class="py-3 px-4">Nama Suku Cadang / Jasa</th>
              <th class="py-3 px-4">Kategori</th>
              <th class="py-3 px-4 text-right">Harga Beli</th>
              <th class="py-3 px-4 text-right">Harga Jual</th>
              <th class="py-3 px-4 text-center">Stok Fisik</th>
              <th class="py-3 px-4 text-center">Status Stok</th>
              <th class="py-3 px-4 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="item in filteredItems" :key="item.id" class="hover:bg-slate-50/60 transition">
              <td class="py-3 px-4 font-mono font-medium text-slate-600">{{ item.part_number }}</td>
              <td class="py-3 px-4">
                <span v-if="item.type === 'service'" class="text-blue-600 font-bold mr-1">[JASA]</span>
                <span class="font-semibold text-slate-800">{{ item.name }}</span>
                <p class="text-[10px] text-slate-400">{{ item.location || 'Area Bengkel' }}</p>
              </td>
              <td class="py-3 px-4 text-slate-600">{{ item.category?.name || '-' }}</td>
              <td class="py-3 px-4 text-right font-mono text-slate-500">{{ formatRupiah(item.cost_price) }}</td>
              <td class="py-3 px-4 text-right font-mono font-bold text-slate-900">{{ formatRupiah(item.selling_price) }}</td>
              <td class="py-3 px-4 text-center font-mono font-bold text-slate-800">
                {{ item.type === 'service' ? '∞' : `${item.stock} ${item.unit}` }}
              </td>
              <td class="py-3 px-4 text-center">
                <span v-if="item.type === 'service'" class="bg-blue-50 text-blue-700 px-2 py-0.5 rounded text-[10px] font-bold">
                  Layanan Jasa
                </span>
                <span 
                  v-else-if="item.stock <= item.min_stock"
                  class="bg-red-100 text-red-800 border border-red-200 px-2 py-0.5 rounded text-[10px] font-bold inline-flex items-center space-x-1">
                  <AlertCircle class="w-3 h-3 text-red-600" />
                  <span>Stok Menipis (&le;{{ item.min_stock }})</span>
                </span>
                <span v-else class="bg-emerald-50 text-emerald-700 px-2 py-0.5 rounded text-[10px] font-bold">
                  Stok Aman
                </span>
              </td>
              <td class="py-3 px-4 text-center">
                <button
                  v-if="item.type === 'part'"
                  @click="openAdjustStock(item)"
                  class="p-1.5 bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-lg text-[11px] font-semibold transition">
                  Sesuaikan Stok
                </button>
              </td>
            </tr>

            <tr v-if="filteredItems.length === 0">
              <td colspan="8" class="py-12 text-center text-slate-400">
                Tidak ada data suku cadang sesuai filter.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mobile Responsive Cards (Screens < 768px) -->
      <div class="md:hidden divide-y divide-slate-100">
        <div v-for="item in filteredItems" :key="item.id" class="p-4 space-y-2.5">
          <div class="flex items-start justify-between">
            <div>
              <span class="font-mono text-xs text-slate-500 block">{{ item.part_number }}</span>
              <span v-if="item.type === 'service'" class="text-blue-600 font-bold text-xs mr-1">[JASA]</span>
              <span class="font-bold text-slate-900 text-sm">{{ item.name }}</span>
            </div>
            <span 
              v-if="item.type === 'part' && item.stock <= item.min_stock"
              class="bg-red-100 text-red-800 text-[10px] font-bold px-2 py-0.5 rounded border border-red-200">
              Stok Menipis
            </span>
          </div>

          <div class="flex justify-between items-center text-xs">
            <span class="text-slate-500">Stok: <strong>{{ item.type === 'service' ? 'Tersedia' : `${item.stock} ${item.unit}` }}</strong></span>
            <span class="font-mono font-black text-sm text-slate-900">{{ formatRupiah(item.selling_price) }}</span>
          </div>

          <div v-if="item.type === 'part'" class="pt-2 border-t border-slate-100 flex justify-end">
            <button
              @click="openAdjustStock(item)"
              class="bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold text-xs px-3 py-1.5 rounded-lg transition">
              Sesuaikan Stok
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Adjust Stock -->
    <div v-if="showAdjustModal && selectedPart" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl shadow-xl max-w-sm w-full p-6 border border-slate-200">
        <h3 class="font-bold text-slate-900 text-sm mb-1">Penyesuaian Stok Gudang</h3>
        <p class="text-xs text-slate-500 font-mono mb-4">{{ selectedPart.part_number }} - {{ selectedPart.name }}</p>

        <div class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Jumlah Perubahan (+ Masuk / - Keluar)</label>
            <input
              v-model.number="adjustQty"
              type="number"
              required
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2 font-mono text-xs focus:outline-none focus:border-red-500" />
            <p class="text-[10px] text-slate-400 mt-1">Stok saat ini: {{ selectedPart.stock }} {{ selectedPart.unit }}.</p>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Alasan Penyesuaian</label>
            <input
              v-model="adjustReason"
              type="text"
              placeholder="Contoh: Penerimaan restock AHM / Rusak"
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2 text-xs focus:outline-none focus:border-red-500" />
          </div>
        </div>

        <div class="mt-6 flex justify-end space-x-2">
          <button
            @click="showAdjustModal = false"
            class="px-4 py-2 text-xs font-semibold text-slate-600 hover:text-slate-800">
            Batal
          </button>
          <button
            @click="handleSaveStock"
            class="px-4 py-2 text-xs font-bold text-white bg-red-600 hover:bg-red-700 rounded-lg">
            Simpan Stok
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Search, Plus, AlertCircle } from 'lucide-vue-next'
import { useSparepartStore } from '../stores/sparepart'
import { useToastStore } from '../stores/toast'
import { playClickSound, playSuccessChime } from '../utils/sound'
import { formatRupiah } from '../utils/format'
import type { Sparepart } from '../types'

const sparepartStore = useSparepartStore()
const toastStore = useToastStore()
const searchQuery = ref('')
const selectedTab = ref('all')

const filterTabs = [
  { label: 'Semua Suku Cadang & Jasa', value: 'all' },
  { label: 'Hanya Honda Genuine Parts', value: 'part' },
  { label: 'Hanya Jasa Servis', value: 'service' },
  { label: 'Peringatan Stok Menipis', value: 'low_stock' },
]

const showAdjustModal = ref(false)
const showAddModal = ref(false)
const selectedPart = ref<Sparepart | null>(null)
const adjustQty = ref<number>(5)
const adjustReason = ref('Restock rutin AHASS NSS')

onMounted(() => {
  sparepartStore.fetchSpareparts()
  sparepartStore.fetchCategories()
})

const filteredItems = computed(() => {
  let list = sparepartStore.spareparts

  if (selectedTab.value === 'part') {
    list = list.filter((i) => i.type === 'part')
  } else if (selectedTab.value === 'service') {
    list = list.filter((i) => i.type === 'service')
  } else if (selectedTab.value === 'low_stock') {
    list = list.filter((i) => i.type === 'part' && i.stock <= i.min_stock)
  }

  if (searchQuery.value.trim() !== '') {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(
      (i) => i.name.toLowerCase().includes(q) || i.part_number.toLowerCase().includes(q)
    )
  }

  return list
})

function openAdjustStock(item: Sparepart) {
  playClickSound()
  selectedPart.value = item
  adjustQty.value = 5
  showAdjustModal.value = true
}

async function handleSaveStock() {
  if (!selectedPart.value) return
  playClickSound()
  try {
    await sparepartStore.updateStock(selectedPart.value.id, adjustQty.value, adjustReason.value)
    playSuccessChime()
    toastStore.success(`Stok suku cadang ${selectedPart.value.name} berhasil diperbarui!`)
    showAdjustModal.value = false
    sparepartStore.fetchSpareparts()
  } catch (err: any) {
    toastStore.error(err.message || 'Gagal mengubah stok')
  }
}
</script>

