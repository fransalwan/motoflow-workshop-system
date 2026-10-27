<template>
  <div class="max-w-5xl mx-auto space-y-8">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900 tracking-tight">Formulir Pendaftaran Servis & Estimasi Biaya</h2>
        <p class="text-xs text-slate-500">Service Advisor - Penerimaan unit motor Honda & diagnosis kerusakan</p>
      </div>
      <router-link to="/work-orders" class="text-xs text-slate-500 hover:text-slate-800 font-medium">
        &larr; Kembali ke Daftar WO
      </router-link>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-8">
      <!-- Section 1: Data Pelanggan & Kendaraan -->
      <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs space-y-6">
        <div class="flex items-center space-x-2 border-b border-slate-100 pb-3">
          <UserIcon class="w-4 h-4 text-red-600" />
          <h3 class="font-bold text-slate-800 text-sm uppercase tracking-wider">1. Data Pemilik & Kendaraan Honda</h3>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Pilih Pelanggan Terdaftar -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-2">Pilih Pelanggan Terdaftar</label>
            <select
              v-model="selectedCustomerId"
              @change="onCustomerChange"
              required
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 text-xs text-slate-800 focus:outline-none focus:border-red-500 transition">
              <option :value="null" disabled>-- Pilih Pelanggan & Motor --</option>
              <option v-for="c in customerStore.customers" :key="c.id" :value="c.id">
                [{{ c.plate_number }}] {{ c.name }} - {{ c.motor_model }}
              </option>
            </select>
          </div>

          <!-- Odometer Masuk (KM) -->
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-2">Kilometer Odometer Masuk (KM)</label>
            <input
              v-model.number="entryKM"
              type="number"
              min="0"
              required
              placeholder="Contoh: 14500"
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 text-xs text-slate-800 focus:outline-none focus:border-red-500 transition font-mono" />
          </div>
        </div>

        <!-- Keluhan & Diagnosa SA -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 pt-2">
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-2">Keluhan Konsumen (Symptom)</label>
            <textarea
              v-model="complaints"
              rows="3"
              required
              placeholder="Contoh: Tarikan awal tersendat, rem depan berdecit saat macet..."
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 text-xs text-slate-800 focus:outline-none focus:border-red-500 transition"></textarea>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-2">Analisa & Rekomendasi Service Advisor</label>
            <textarea
              v-model="serviceDiagnosis"
              rows="3"
              placeholder="Contoh: Servis lengkap matic, ganti oli MPX2, cek kampas ganda CVT..."
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 text-xs text-slate-800 focus:outline-none focus:border-red-500 transition"></textarea>
          </div>
        </div>
      </div>

      <!-- Section 2: Penugasan Mekanik & Pit -->
      <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs space-y-6">
        <div class="flex items-center space-x-2 border-b border-slate-100 pb-3">
          <Wrench class="w-4 h-4 text-red-600" />
          <h3 class="font-bold text-slate-800 text-sm uppercase tracking-wider">2. Penugasan Mekanik & Pit AHASS</h3>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-2">Pilih Mekanik Penanggung Jawab</label>
            <select
              v-model="selectedMechanicId"
              class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 text-xs text-slate-800 focus:outline-none focus:border-red-500 transition">
              <option :value="null">-- Belum Ditentukan (Antrean Pit) --</option>
              <option v-for="m in mechanics" :key="m.id" :value="m.id">
                Pit {{ m.pit_number }} - {{ m.name }} ({{ m.specialization }})
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Section 3: Rincian Suku Cadang HGP & Jasa Servis -->
      <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs space-y-6">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div class="flex items-center space-x-2">
            <PackageCheck class="w-4 h-4 text-red-600" />
            <h3 class="font-bold text-slate-800 text-sm uppercase tracking-wider">3. Pemilihan Suku Cadang (HGP) & Jasa Servis</h3>
          </div>
          <button
            type="button"
            @click="addItemRow"
            class="inline-flex items-center space-x-1 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold text-xs px-3 py-1.5 rounded-lg transition">
            <Plus class="w-3.5 h-3.5" />
            <span>Tambah Item</span>
          </button>
        </div>

        <!-- Items Table -->
        <div class="overflow-x-auto">
          <table class="w-full text-xs text-left">
            <thead>
              <tr class="border-b border-slate-200 text-slate-500 font-bold uppercase text-[10px]">
                <th class="py-2.5 w-6/12">Suku Cadang HGP / Jasa Servis</th>
                <th class="py-2.5 w-1/12 text-center">Qty</th>
                <th class="py-2.5 w-2/12 text-right">Harga Satuan</th>
                <th class="py-2.5 w-2/12 text-right">Subtotal</th>
                <th class="py-2.5 w-1/12 text-center">Hapus</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="(row, idx) in selectedItems" :key="idx" class="align-top">
                <!-- Dropdown Suku Cadang -->
                <td class="py-3 pr-4">
                  <select
                    v-model="row.sparepart_id"
                    @change="onSparepartSelect(row)"
                    required
                    class="w-full bg-slate-50 border border-slate-200 rounded-lg px-3 py-2 text-xs focus:outline-none focus:border-red-500">
                    <option :value="null" disabled>-- Pilih Part / Jasa --</option>
                    <option v-for="part in sparepartStore.spareparts" :key="part.id" :value="part.id">
                      [{{ part.type === 'service' ? 'JASA' : 'PART' }}] {{ part.name }} ({{ formatRupiah(part.selling_price) }}) - Stok: {{ part.type === 'service' ? 'Tersedia' : part.stock }}
                    </option>
                  </select>
                </td>

                <!-- Quantity -->
                <td class="py-3 px-2 text-center">
                  <input
                    v-model.number="row.quantity"
                    type="number"
                    min="1"
                    required
                    class="w-16 bg-slate-50 border border-slate-200 rounded-lg py-1.5 px-2 text-center text-xs font-mono focus:outline-none focus:border-red-500" />
                </td>

                <!-- Price Snapshot Display -->
                <td class="py-3 px-2 text-right font-mono font-medium text-slate-700">
                  {{ formatRupiah(row.price) }}
                </td>

                <!-- Subtotal Display -->
                <td class="py-3 px-2 text-right font-mono font-bold text-slate-900">
                  {{ formatRupiah(Number(row.price) * row.quantity) }}
                </td>

                <!-- Delete Row -->
                <td class="py-3 text-center">
                  <button
                    type="button"
                    @click="removeItemRow(idx)"
                    :disabled="selectedItems.length <= 1"
                    class="text-slate-400 hover:text-red-600 disabled:opacity-30 transition">
                    <Trash2 class="w-4 h-4 mx-auto" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Financial Summary Calculator -->
        <div class="mt-6 pt-6 border-t border-slate-200 flex justify-end">
          <div class="w-80 space-y-2 text-xs">
            <div class="flex justify-between text-slate-600">
              <span>Subtotal Suku Cadang:</span>
              <span class="font-mono font-semibold">{{ formatRupiah(calculatedPartsSubtotal) }}</span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span>Subtotal Jasa Servis:</span>
              <span class="font-mono font-semibold">{{ formatRupiah(calculatedServicesSubtotal) }}</span>
            </div>
            <div class="flex justify-between text-slate-600">
              <span>PPN (11%):</span>
              <span class="font-mono font-semibold">{{ formatRupiah(calculatedTax) }}</span>
            </div>
            <div class="flex justify-between font-black text-slate-900 text-sm pt-2 border-t border-slate-300">
              <span class="text-red-600">ESTIMASI TOTAL:</span>
              <span class="font-mono">{{ formatRupiah(calculatedGrandTotal) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center justify-end space-x-4">
        <router-link
          to="/work-orders"
          class="px-5 py-2.5 text-xs font-semibold text-slate-600 hover:text-slate-900 bg-white border border-slate-300 rounded-xl shadow-xs transition">
          Batal
        </router-link>
        <button
          type="submit"
          :disabled="workOrderStore.loading"
          class="px-6 py-2.5 text-xs font-bold text-white bg-red-600 hover:bg-red-700 rounded-xl shadow-md transition disabled:opacity-50">
          <span v-if="!workOrderStore.loading">Terbitkan Work Order (WO)</span>
          <span v-else>Memproses WO...</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  User as UserIcon,
  Wrench,
  PackageCheck,
  Plus,
  Trash2,
} from 'lucide-vue-next'
import { useCustomerStore } from '../stores/customer'
import { useSparepartStore } from '../stores/sparepart'
import { useWorkOrderStore } from '../stores/workOrder'
import api from '../services/api'
import { formatRupiah } from '../utils/format'
import type { Mechanic } from '../types'

const router = useRouter()
const customerStore = useCustomerStore()
const sparepartStore = useSparepartStore()
const workOrderStore = useWorkOrderStore()

const mechanics = ref<Mechanic[]>([])
const selectedCustomerId = ref<number | null>(null)
const entryKM = ref<number>(0)
const complaints = ref('')
const serviceDiagnosis = ref('')
const selectedMechanicId = ref<number | null>(null)

interface RowItem {
  sparepart_id: number | null
  item_type: 'part' | 'service'
  part_number: string
  item_name: string
  quantity: number
  price: number
  discount_percent: number
  mechanic_notes: string
}

const selectedItems = ref<RowItem[]>([
  {
    sparepart_id: null,
    item_type: 'service',
    part_number: '',
    item_name: '',
    quantity: 1,
    price: 0,
    discount_percent: 0,
    mechanic_notes: '',
  },
])

onMounted(async () => {
  customerStore.fetchCustomers()
  sparepartStore.fetchSpareparts()
  try {
    const res = await api.get('/mechanics')
    mechanics.value = res.data.data
  } catch (err) {
    console.error('Failed to load mechanics', err)
  }
})

function onCustomerChange() {
  const cust = customerStore.customers.find((c) => c.id === selectedCustomerId.value)
  if (cust) {
    entryKM.value = cust.current_km || 0
  }
}

function addItemRow() {
  selectedItems.value.push({
    sparepart_id: null,
    item_type: 'part',
    part_number: '',
    item_name: '',
    quantity: 1,
    price: 0,
    discount_percent: 0,
    mechanic_notes: '',
  })
}

function removeItemRow(index: number) {
  if (selectedItems.value.length > 1) {
    selectedItems.value.splice(index, 1)
  }
}

function onSparepartSelect(row: RowItem) {
  const part = sparepartStore.spareparts.find((p) => p.id === row.sparepart_id)
  if (part) {
    row.part_number = part.part_number
    row.item_name = part.name
    row.item_type = part.type
    row.price = Number(part.selling_price)
  }
}

const calculatedPartsSubtotal = computed(() => {
  return selectedItems.value
    .filter((r) => r.item_type === 'part')
    .reduce((acc, r) => acc + Number(r.price) * r.quantity, 0)
})

const calculatedServicesSubtotal = computed(() => {
  return selectedItems.value
    .filter((r) => r.item_type === 'service')
    .reduce((acc, r) => acc + Number(r.price) * r.quantity, 0)
})

const calculatedTax = computed(() => {
  const taxable = calculatedPartsSubtotal.value + calculatedServicesSubtotal.value
  return Math.round(taxable * 0.11)
})

const calculatedGrandTotal = computed(() => {
  return calculatedPartsSubtotal.value + calculatedServicesSubtotal.value + calculatedTax.value
})

async function handleSubmit() {
  if (!selectedCustomerId.value) {
    alert('Silakan pilih pelanggan terlebih dahulu')
    return
  }

  // Validate items
  for (const item of selectedItems.value) {
    if (!item.sparepart_id) {
      alert('Terdapat baris suku cadang/jasa yang belum dipilih')
      return
    }
  }

  const payload = {
    customer_id: selectedCustomerId.value,
    mechanic_id: selectedMechanicId.value,
    entry_km: entryKM.value,
    complaints: complaints.value,
    service_diagnosis: serviceDiagnosis.value,
    items: selectedItems.value.map((i) => ({
      sparepart_id: i.sparepart_id,
      item_type: i.item_type,
      part_number: i.part_number,
      item_name: i.item_name,
      quantity: i.quantity,
      price: i.price,
      discount_percent: i.discount_percent,
      mechanic_notes: i.mechanic_notes,
    })),
  }

  try {
    const created = await workOrderStore.createWorkOrder(payload)
    router.push(`/work-orders/${created.id}`)
  } catch (err: any) {
    alert(err.message || 'Gagal membuat Work Order')
  }
}
</script>

