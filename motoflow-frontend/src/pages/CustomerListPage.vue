<template>
  <div class="space-y-6">
    <!-- Header Controls -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-slate-900 tracking-tight">Master Data Pelanggan & Unit Kendaraan Honda</h2>
        <p class="text-xs text-slate-500">Pencatatan riwayat motor, nomor plat, nomor rangka, dan kilometer odometer</p>
      </div>

      <div class="flex items-center space-x-3">
        <div class="relative">
          <input
            v-model="searchQuery"
            @input="onSearch"
            type="text"
            placeholder="Cari Plat, Nama, atau Model..."
            class="bg-white border border-slate-200 rounded-xl pl-9 pr-4 py-2 text-xs text-slate-900 placeholder-slate-400 focus:outline-none focus:border-red-500 w-72 shadow-xs" />
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
        </div>

        <button
          @click="showAddModal = true"
          class="inline-flex items-center space-x-1.5 bg-red-600 hover:bg-red-700 text-white font-semibold text-xs px-4 py-2 rounded-xl shadow-xs transition">
          <Plus class="w-4 h-4" />
          <span>Tambah Pelanggan</span>
        </button>
      </div>
    </div>

    <!-- Customer Cards / Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <!-- Desktop Table -->
      <div class="hidden md:block overflow-x-auto">
        <table class="w-full text-xs text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-500 font-bold uppercase text-[10px] tracking-wider">
              <th class="py-3 px-4">No. Polisi Plat</th>
              <th class="py-3 px-4">Nama Pelanggan</th>
              <th class="py-3 px-4">Model & Tahun</th>
              <th class="py-3 px-4">No. Rangka & Mesin</th>
              <th class="py-3 px-4 text-right">Odometer Terakhir</th>
              <th class="py-3 px-4">Kontak & Alamat</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="c in customerStore.customers" :key="c.id" class="hover:bg-slate-50/60 transition">
              <td class="py-3.5 px-4">
                <span class="font-mono font-black bg-slate-100 text-slate-900 px-2.5 py-1 rounded-lg border border-slate-300 text-xs">
                  {{ c.plate_number }}
                </span>
              </td>
              <td class="py-3.5 px-4 font-bold text-slate-800">{{ c.name }}</td>
              <td class="py-3.5 px-4">
                <span class="font-semibold text-slate-700">{{ c.motor_model }}</span>
                <span v-if="c.year" class="text-slate-400 ml-1">({{ c.year }})</span>
              </td>
              <td class="py-3.5 px-4 font-mono text-[11px] text-slate-500">
                <p>R: {{ c.frame_number || '-' }}</p>
                <p>M: {{ c.engine_number || '-' }}</p>
              </td>
              <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900">
                {{ c.current_km }} KM
              </td>
              <td class="py-3.5 px-4 text-slate-500 text-[11px]">
                <p class="font-medium text-slate-700">{{ c.phone }}</p>
                <p class="truncate max-w-xs">{{ c.address || '-' }}</p>
              </td>
            </tr>

            <tr v-if="customerStore.customers.length === 0">
              <td colspan="6" class="py-12 text-center text-slate-400">
                Belum ada data pelanggan terdaftar.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mobile Responsive Cards (Screens < 768px) -->
      <div class="md:hidden divide-y divide-slate-100">
        <div v-for="c in customerStore.customers" :key="c.id" class="p-4 space-y-2">
          <div class="flex items-center justify-between">
            <span class="font-mono font-black text-xs bg-slate-100 text-slate-900 px-2.5 py-0.5 rounded border border-slate-300">
              {{ c.plate_number }}
            </span>
            <span class="font-mono font-bold text-xs text-slate-700">{{ c.current_km }} KM</span>
          </div>

          <div>
            <h4 class="font-bold text-sm text-slate-900">{{ c.name }}</h4>
            <p class="text-xs text-slate-600 font-medium">{{ c.motor_model }} ({{ c.year || '2023' }})</p>
          </div>

          <div class="pt-2 border-t border-slate-100 flex justify-between items-center text-[11px] text-slate-500">
            <span>{{ c.phone }}</span>
            <span class="truncate max-w-[150px]">{{ c.address || '-' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Tambah Pelanggan -->
    <div v-if="showAddModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl shadow-xl max-w-lg w-full p-6 border border-slate-200">
        <h3 class="font-bold text-slate-900 text-base mb-4">Pendaftaran Pelanggan & Motor Honda Baru</h3>

        <form @submit.prevent="handleAddCustomer" class="space-y-4 text-xs">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block font-semibold text-slate-700 mb-1">Nama Pemilik *</label>
              <input v-model="form.name" type="text" required placeholder="Contoh: Budi Santoso" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2" />
            </div>
            <div>
              <label class="block font-semibold text-slate-700 mb-1">Nomor Telepon (WhatsApp) *</label>
              <input v-model="form.phone" type="text" required placeholder="Contoh: 08123456789" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2" />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block font-semibold text-slate-700 mb-1">No. Polisi Plat Motor *</label>
              <input v-model="form.plate_number" type="text" required placeholder="Contoh: B 1234 ABC" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2 font-mono uppercase" />
            </div>
            <div>
              <label class="block font-semibold text-slate-700 mb-1">Model Motor Honda *</label>
              <input v-model="form.motor_model" type="text" required placeholder="Contoh: Honda Vario 160 ABS" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2" />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block font-semibold text-slate-700 mb-1">Tahun Perakitan</label>
              <input v-model.number="form.year" type="number" placeholder="Contoh: 2023" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2" />
            </div>
            <div>
              <label class="block font-semibold text-slate-700 mb-1">Odometer Terkini (KM)</label>
              <input v-model.number="form.current_km" type="number" placeholder="Contoh: 12500" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2 font-mono" />
            </div>
          </div>

          <div>
            <label class="block font-semibold text-slate-700 mb-1">Alamat Tinggal</label>
            <textarea v-model="form.address" rows="2" placeholder="Alamat domisili pemilik" class="w-full bg-slate-50 border border-slate-200 rounded-xl px-3 py-2"></textarea>
          </div>

          <div class="mt-6 flex justify-end space-x-2 pt-2 border-t border-slate-100">
            <button type="button" @click="showAddModal = false" class="px-4 py-2 text-slate-600 hover:text-slate-800">
              Batal
            </button>
            <button type="submit" class="px-5 py-2 font-bold text-white bg-red-600 hover:bg-red-700 rounded-xl shadow-xs">
              Daftarkan Pelanggan
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search, Plus } from 'lucide-vue-next'
import { useCustomerStore } from '../stores/customer'
import { useToastStore } from '../stores/toast'
import { playSuccessChime } from '../utils/sound'

const customerStore = useCustomerStore()
const toast = useToastStore()
const searchQuery = ref('')
const showAddModal = ref(false)

const form = ref({
  name: '',
  phone: '',
  plate_number: '',
  motor_model: '',
  year: 2023,
  current_km: 10000,
  address: '',
})

onMounted(() => {
  customerStore.fetchCustomers()
})

let debounceTimer: any = null
function onSearch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    customerStore.fetchCustomers(searchQuery.value)
  }, 300)
}

async function handleAddCustomer() {
  try {
    const plate = form.value.plate_number.toUpperCase()
    await customerStore.createCustomer(form.value)
    playSuccessChime()
    toast.success(`Pelanggan & motor ${plate} berhasil didaftarkan!`)
    showAddModal.value = false
    form.value = {
      name: '',
      phone: '',
      plate_number: '',
      motor_model: '',
      year: 2023,
      current_km: 10000,
      address: '',
    }
  } catch (err: any) {
    toast.error(err.message || 'Gagal menambahkan pelanggan')
  }
}
</script>

