<template>
  <div class="space-y-6">
    <!-- Top Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <div class="flex items-center space-x-2">
          <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold uppercase tracking-wider bg-red-100 text-red-700">
            Direct Over-The-Counter
          </span>
          <span class="text-xs text-slate-500 font-mono">PT Nusantara Sakti Group</span>
        </div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight mt-1 flex items-center">
          <ShoppingBag class="w-6 h-6 mr-2 text-red-600" />
          Kasir POS Penjualan Suku Cadang & Oli HGP
        </h1>
        <p class="text-xs text-slate-500">
          Transaksi langsung tanpa pendaftaran servis mekanik (Walk-in Counter POS).
        </p>
      </div>

      <div class="flex items-center space-x-2">
        <button
          @click="showHistoryModal = true"
          class="flex items-center px-4 py-2 text-xs font-bold text-slate-700 bg-white border border-slate-200 rounded-xl hover:bg-slate-50 shadow-xs"
        >
          <ReceiptText class="w-4 h-4 mr-1.5 text-slate-500" />
          Riwayat Struk Retail ({{ retailStore.sales.length }})
        </button>
      </div>
    </div>

    <!-- Mobile Segmented Tab Switcher (Visible on < lg) -->
    <div class="lg:hidden flex bg-slate-200/80 p-1 rounded-2xl text-xs font-bold">
      <button
        @click="mobileTab = 'catalog'"
        :class="mobileTab === 'catalog' ? 'bg-white text-slate-900 shadow-xs' : 'text-slate-600'"
        class="flex-1 py-2.5 rounded-xl transition flex items-center justify-center space-x-1.5"
      >
        <ShoppingBag class="w-4 h-4" />
        <span>Katalog Part</span>
      </button>
      <button
        @click="mobileTab = 'cart'"
        :class="mobileTab === 'cart' ? 'bg-white text-slate-900 shadow-xs' : 'text-slate-600'"
        class="flex-1 py-2.5 rounded-xl transition flex items-center justify-center space-x-1.5 relative"
      >
        <ShoppingCart class="w-4 h-4" />
        <span>Keranjang</span>
        <span
          v-if="retailStore.cartTotalItems > 0"
          class="bg-red-600 text-white text-[10px] px-2 py-0.5 rounded-full font-mono font-black"
        >
          {{ retailStore.cartTotalItems }}
        </span>
      </button>
    </div>

    <!-- Main Layout: 2 Columns (Catalog Left, Cart Right) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Left: Products Catalog (7 cols on lg) -->
      <div 
        :class="mobileTab === 'catalog' ? 'block' : 'hidden lg:block'"
        class="lg:col-span-7 space-y-4">
        <!-- Search & Filter Controls -->
        <div class="bg-white p-4 rounded-2xl border border-slate-200/80 shadow-xs space-y-3">
          <div class="relative">
            <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              v-model="searchQuery"
              @input="onSearch"
              type="text"
              placeholder="Scan barcode / ketik nama suku cadang, oli, nomor part..."
              class="w-full bg-slate-50 border border-slate-200 rounded-xl pl-10 pr-10 py-2.5 text-xs text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:bg-white transition-all font-mono"
            />
            <Barcode class="w-4 h-4 absolute right-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
          </div>

          <!-- Category Pills -->
          <div class="flex items-center space-x-1.5 overflow-x-auto pb-1 text-xs">
            <button
              v-for="cat in categories"
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

        <!-- Product Cards Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 max-h-[600px] overflow-y-auto pr-1">
          <div
            v-for="part in filteredParts"
            :key="part.id"
            class="bg-white p-3.5 rounded-xl border border-slate-200/80 hover:border-red-300 hover:shadow-md transition-all flex flex-col justify-between group"
          >
            <div>
              <div class="flex items-center justify-between mb-1.5">
                <span class="font-mono text-[10px] bg-slate-100 text-slate-600 px-2 py-0.5 rounded font-bold">
                  {{ part.part_number }}
                </span>
                <span
                  :class="[
                    'text-[10px] font-bold px-2 py-0.5 rounded-full',
                    part.stock > 10 ? 'bg-emerald-100 text-emerald-800' :
                    part.stock > 0 ? 'bg-amber-100 text-amber-800' : 'bg-red-100 text-red-800'
                  ]"
                >
                  Stok: {{ part.stock }}
                </span>
              </div>
              <h4 class="font-bold text-xs text-slate-900 group-hover:text-red-600 transition-colors line-clamp-2">
                {{ part.name }}
              </h4>
            </div>

            <div class="mt-3 pt-2.5 border-t border-slate-100 flex items-center justify-between">
              <div>
                <span class="text-[10px] text-slate-400 block font-medium">Harga HET</span>
                <span class="text-xs font-black text-slate-900 font-mono">
                  {{ formatRupiah(part.selling_price) }}
                </span>
              </div>

              <button
                @click="handleAddToCart(part)"
                :disabled="part.stock <= 0"
                class="px-3 py-1.5 rounded-lg text-xs font-bold text-white bg-red-600 hover:bg-red-700 disabled:bg-slate-300 disabled:cursor-not-allowed shadow-xs flex items-center transition-transform active:scale-95"
              >
                <Plus class="w-3.5 h-3.5 mr-1" />
                Beli
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Sticky Cart & Checkout Panel (5 cols on lg) -->
      <div 
        :class="mobileTab === 'cart' ? 'block' : 'hidden lg:block'"
        class="lg:col-span-5 space-y-4">
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-md p-4 sm:p-5 flex flex-col justify-between min-h-[580px]">
          <!-- Cart Header -->
          <div>
            <div class="flex items-center justify-between pb-3 border-b border-slate-100">
              <div class="flex items-center space-x-2">
                <ShoppingCart class="w-5 h-5 text-red-600" />
                <h3 class="font-bold text-sm text-slate-900">Keranjang Kasir</h3>
                <span class="bg-red-100 text-red-700 text-xs font-mono font-bold px-2 py-0.5 rounded-full">
                  {{ retailStore.cartTotalItems }} item
                </span>
              </div>
              <div class="flex items-center space-x-2">
                <button
                  v-if="mobileTab === 'cart'"
                  @click="mobileTab = 'catalog'"
                  class="lg:hidden text-xs font-bold text-slate-600 bg-slate-100 px-2.5 py-1 rounded-lg hover:bg-slate-200"
                >
                  + Tambah Part
                </button>
                <button
                  v-if="retailStore.cart.length > 0"
                  @click="retailStore.clearCart"
                  class="text-[11px] font-semibold text-slate-400 hover:text-red-600"
                >
                  Kosongkan
                </button>
              </div>
            </div>

            <!-- Customer Identity -->
            <div class="mt-3 grid grid-cols-2 gap-2 text-xs">
              <div>
                <label class="block text-[10px] font-semibold text-slate-500 mb-1">Nama Pembeli</label>
                <input
                  v-model="customerName"
                  type="text"
                  placeholder="Walk-in Customer"
                  class="w-full bg-slate-50 border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs text-slate-800"
                />
              </div>
              <div>
                <label class="block text-[10px] font-semibold text-slate-500 mb-1">No. WhatsApp</label>
                <input
                  v-model="customerPhone"
                  type="text"
                  placeholder="0812xxxx"
                  class="w-full bg-slate-50 border border-slate-200 rounded-lg px-2.5 py-1.5 text-xs text-slate-800"
                />
              </div>
            </div>

            <!-- Cart Items List -->
            <div class="mt-3 space-y-2 max-h-[220px] overflow-y-auto pr-1">
              <div
                v-if="retailStore.cart.length === 0"
                class="py-10 text-center text-slate-400 text-xs"
              >
                <ShoppingBag class="w-8 h-8 mx-auto mb-2 text-slate-300" />
                Keranjang belanja masih kosong.<br/>Pilih suku cadang atau scan barcode di sebelah kiri.
              </div>

              <div
                v-for="item in retailStore.cart"
                :key="item.sparepart.id"
                class="p-2.5 bg-slate-50 rounded-xl border border-slate-100 flex items-center justify-between text-xs"
              >
                <div class="flex-1 min-w-0 pr-2">
                  <h5 class="font-bold text-slate-900 truncate">{{ item.sparepart.name }}</h5>
                  <span class="text-[11px] text-slate-500 font-mono">
                    {{ formatRupiah(item.sparepart.selling_price) }}
                  </span>
                </div>

                <!-- Qty Stepper (Touch-Friendly Min 36px) -->
                <div class="flex items-center space-x-1.5">
                  <button
                    @click="retailStore.updateQty(item.sparepart.id, item.quantity - 1)"
                    class="w-7 h-7 sm:w-8 sm:h-8 rounded-lg bg-white border border-slate-200 text-slate-700 flex items-center justify-center font-bold text-sm hover:bg-slate-100 active:scale-95 shadow-xs transition"
                  >
                    -
                  </button>
                  <span class="font-mono font-bold text-xs w-6 text-center">{{ item.quantity }}</span>
                  <button
                    @click="retailStore.updateQty(item.sparepart.id, item.quantity + 1)"
                    class="w-7 h-7 sm:w-8 sm:h-8 rounded-lg bg-white border border-slate-200 text-slate-700 flex items-center justify-center font-bold text-sm hover:bg-slate-100 active:scale-95 shadow-xs transition"
                  >
                    +
                  </button>
                  <button
                    @click="retailStore.removeFromCart(item.sparepart.id)"
                    class="p-1.5 text-slate-400 hover:text-red-600 ml-1 active:scale-90 transition"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Promo Voucher Input -->
            <div class="mt-3 pt-3 border-t border-slate-100">
              <div class="flex space-x-2">
                <div class="relative flex-1">
                  <Tag class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input
                    v-model="promoInput"
                    type="text"
                    placeholder="Kode Promo (misal: NSSHEMAT10)"
                    class="w-full bg-slate-50 border border-slate-200 rounded-lg pl-8 pr-2 py-1.5 text-xs uppercase font-mono font-bold"
                  />
                </div>
                <button
                  @click="handleApplyPromo"
                  class="px-3 py-1.5 text-xs font-bold text-red-600 bg-red-50 hover:bg-red-100 border border-red-200 rounded-lg transition-colors"
                >
                  Pasang
                </button>
              </div>

              <!-- Promo Tag Chips -->
              <div class="flex items-center space-x-1.5 mt-2 overflow-x-auto text-[10px]">
                <span
                  v-for="p in retailStore.promotions"
                  :key="p.promo_code"
                  @click="promoInput = p.promo_code; handleApplyPromo()"
                  class="cursor-pointer bg-slate-100 hover:bg-red-50 text-slate-600 hover:text-red-700 px-2 py-0.5 rounded border border-slate-200 font-mono font-bold whitespace-nowrap"
                >
                  {{ p.promo_code }}
                </span>
              </div>

              <div v-if="retailStore.appliedPromo" class="mt-1.5 flex items-center justify-between text-[11px] text-emerald-700 font-bold bg-emerald-50 px-2 py-1 rounded">
                <span>Promo Aktif: {{ retailStore.appliedPromo.title }}</span>
                <span class="font-mono">-{{ formatRupiah(retailStore.promoDiscountAmount) }}</span>
              </div>
            </div>
          </div>

          <!-- Bottom: Totals & Payment Section -->
          <div class="pt-3 border-t border-slate-100 space-y-3">
            <div class="space-y-1.5 text-xs text-slate-600">
              <div class="flex justify-between">
                <span>Subtotal Part</span>
                <span class="font-mono font-bold">{{ formatRupiah(retailStore.cartSubtotal) }}</span>
              </div>
              <div v-if="retailStore.promoDiscountAmount > 0" class="flex justify-between text-emerald-600">
                <span>Diskon Promo</span>
                <span class="font-mono font-bold">-{{ formatRupiah(retailStore.promoDiscountAmount) }}</span>
              </div>
              <div class="flex justify-between">
                <span>PPN (11%)</span>
                <span class="font-mono font-bold">{{ formatRupiah(retailStore.cartTaxPPN) }}</span>
              </div>
              <div class="flex justify-between text-base font-black text-slate-900 pt-1.5 border-t border-slate-200">
                <span>TOTAL AKHIR</span>
                <span class="font-mono text-red-600">{{ formatRupiah(retailStore.cartGrandTotal) }}</span>
              </div>
            </div>

            <!-- Payment Methods -->
            <div>
              <label class="block text-[10px] font-bold text-slate-500 uppercase mb-1">Metode Pembayaran</label>
              <div class="grid grid-cols-4 gap-1.5 text-xs font-bold">
                <button
                  v-for="method in ['cash', 'qris', 'debit', 'transfer']"
                  :key="method"
                  @click="paymentMethod = method; handleMethodChange(method)"
                  :class="[
                    'py-1.5 rounded-lg border text-center uppercase tracking-wider text-[11px] transition-all',
                    paymentMethod === method
                      ? 'bg-red-600 text-white border-red-600 shadow-xs'
                      : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
                  ]"
                >
                  {{ method }}
                </button>
              </div>
            </div>

            <!-- Quick Cash Buttons -->
            <div v-if="paymentMethod === 'cash'" class="space-y-2">
              <div class="grid grid-cols-4 gap-1.5 text-[11px] font-mono font-bold">
                <button
                  @click="setExactCash"
                  class="py-1 bg-slate-100 hover:bg-slate-200 text-slate-800 rounded border border-slate-200"
                >
                  Uang Pas
                </button>
                <button
                  @click="paidAmount = 50000"
                  class="py-1 bg-slate-100 hover:bg-slate-200 text-slate-800 rounded border border-slate-200"
                >
                  50k
                </button>
                <button
                  @click="paidAmount = 100000"
                  class="py-1 bg-slate-100 hover:bg-slate-200 text-slate-800 rounded border border-slate-200"
                >
                  100k
                </button>
                <button
                  @click="paidAmount = 200000"
                  class="py-1 bg-slate-100 hover:bg-slate-200 text-slate-800 rounded border border-slate-200"
                >
                  200k
                </button>
              </div>

              <div class="flex items-center justify-between text-xs">
                <div class="w-1/2 pr-1">
                  <label class="block text-[10px] text-slate-400 font-semibold mb-0.5">Uang Diterima</label>
                  <input
                    v-model.number="paidAmount"
                    type="number"
                    class="w-full bg-slate-50 border border-slate-200 rounded-lg px-2.5 py-1.5 font-mono font-bold text-slate-900"
                  />
                </div>
                <div class="w-1/2 pl-1">
                  <label class="block text-[10px] text-slate-400 font-semibold mb-0.5">Kembalian</label>
                  <div class="bg-slate-100 border border-slate-200 rounded-lg px-2.5 py-1.5 font-mono font-bold text-emerald-700 truncate">
                    {{ formatRupiah(changeAmount) }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Submit Checkout Button -->
            <button
              @click="handleCheckout"
              :disabled="retailStore.cart.length === 0 || retailStore.loading || (paymentMethod === 'cash' && paidAmount < retailStore.cartGrandTotal)"
              class="w-full py-3 bg-red-600 hover:bg-red-700 disabled:bg-slate-300 disabled:cursor-not-allowed text-white font-black text-sm rounded-xl shadow-md transition-all flex items-center justify-center space-x-2 active:scale-98"
            >
              <CheckCircle class="w-4 h-4" />
              <span>Bayar & Terbitkan Struk POS</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Sticky Cart Bar (Mobile Only: Visible when on Catalog tab and Cart has items) -->
    <div
      v-if="mobileTab === 'catalog' && retailStore.cartTotalItems > 0"
      class="lg:hidden fixed bottom-[calc(4.5rem+env(safe-area-inset-bottom))] left-3 right-3 z-30 bg-slate-900/95 backdrop-blur-md text-white p-3 rounded-2xl shadow-2xl flex items-center justify-between border border-slate-800 animate-in slide-in-from-bottom duration-200"
    >
      <div class="flex items-center space-x-2.5">
        <div class="w-9 h-9 rounded-xl bg-red-600 text-white flex items-center justify-center font-bold text-xs shadow-md">
          {{ retailStore.cartTotalItems }}
        </div>
        <div>
          <span class="text-[10px] text-slate-400 block font-medium">Total Tagihan</span>
          <span class="text-sm font-black font-mono text-white">{{ formatRupiah(retailStore.cartGrandTotal) }}</span>
        </div>
      </div>
      <button
        @click="mobileTab = 'cart'"
        class="bg-red-600 hover:bg-red-500 text-white text-xs font-bold px-4 py-2 rounded-xl flex items-center space-x-1.5 shadow-md active:scale-95 transition"
      >
        <span>Lihat Keranjang</span>
        <ArrowRight class="w-3.5 h-3.5" />
      </button>
    </div>

    <!-- Modal Thermal Receipt POS 58mm/80mm -->
    <div v-if="showReceiptModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-end sm:items-center justify-center p-0 sm:p-4">
      <div class="bg-white rounded-t-3xl sm:rounded-2xl shadow-2xl max-w-md w-full p-6 border border-slate-200 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="w-12 h-1.5 bg-slate-300 rounded-full mx-auto -mt-2 mb-2 sm:hidden"></div>
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center space-x-2">
            <Printer class="w-5 h-5 text-red-600" />
            <h3 class="font-bold text-slate-900 text-sm">Struk Kasir Retail Honda AHASS</h3>
          </div>
          <button @click="showReceiptModal = false" class="text-slate-400 hover:text-slate-600">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- ESC/POS Thermal Paper Preview -->
        <div class="bg-slate-50 p-4 rounded-xl border border-slate-200 font-mono text-[11px] whitespace-pre-wrap leading-tight text-slate-800 max-h-[350px] overflow-y-auto">
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
            <span>Cetak Struk POS</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Riwayat Transaksi Retail -->
    <div v-if="showHistoryModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-end sm:items-center justify-center p-0 sm:p-4">
      <div class="bg-white rounded-t-3xl sm:rounded-2xl shadow-2xl max-w-3xl w-full p-6 border border-slate-200 space-y-4 max-h-[90vh] flex flex-col">
        <div class="w-12 h-1.5 bg-slate-300 rounded-full mx-auto -mt-2 mb-2 sm:hidden"></div>
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center space-x-2">
            <ReceiptText class="w-5 h-5 text-red-600" />
            <h3 class="font-bold text-slate-900 text-base">Riwayat Transaksi Kasir Retail</h3>
          </div>
          <button @click="showHistoryModal = false" class="text-slate-400 hover:text-slate-600">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="overflow-y-auto flex-1 divide-y divide-slate-100">
          <div
            v-for="sale in retailStore.sales"
            :key="sale.id"
            class="py-3 flex items-center justify-between text-xs"
          >
            <div>
              <div class="flex items-center space-x-2">
                <span class="font-mono font-bold text-red-600">{{ sale.invoice_number }}</span>
                <span class="text-slate-400">•</span>
                <span class="font-bold text-slate-900">{{ sale.customer_name }}</span>
                <span class="text-[10px] bg-slate-100 uppercase px-2 py-0.5 rounded font-bold text-slate-600">
                  {{ sale.payment_method }}
                </span>
              </div>
              <div class="text-[11px] text-slate-400 mt-0.5">
                {{ formatDate(sale.created_at) }} • {{ sale.items.length }} item suku cadang
              </div>
            </div>

            <div class="flex items-center space-x-3">
              <span class="font-mono font-black text-slate-900">{{ formatRupiah(sale.grand_total) }}</span>
              <button
                @click="openReceiptForSale(sale.id)"
                class="px-2.5 py-1 text-[11px] font-bold text-red-600 bg-red-50 hover:bg-red-100 rounded-lg border border-red-200"
              >
                Struk
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- QRIS Modal Dynamic Simulator -->
    <QRISModal
      :is-open="showQRISModal"
      :amount="retailStore.cartGrandTotal"
      @close="showQRISModal = false"
      @verified="onQRISVerified"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  ShoppingBag,
  ShoppingCart,
  Search,
  Barcode,
  Plus,
  Trash2,
  Tag,
  CheckCircle,
  Printer,
  ReceiptText,
  X,
  ArrowRight,
} from 'lucide-vue-next'
import { useRetailStore } from '../stores/retail'
import { useSparepartStore } from '../stores/sparepart'
import { useToastStore } from '../stores/toast'
import { formatRupiah, formatDate } from '../utils/format'
import { playScannerBeep, playSuccessChime, playClickSound } from '../utils/sound'
import QRISModal from '../components/QRISModal.vue'
import type { Sparepart } from '../types'

const retailStore = useRetailStore()
const sparepartStore = useSparepartStore()
const toast = useToastStore()

const mobileTab = ref<'catalog' | 'cart'>('catalog')
const searchQuery = ref('')
const selectedCategory = ref('Semua')
const customerName = ref('Konsumen Walk-in')
const customerPhone = ref('')
const promoInput = ref('')
const paymentMethod = ref('cash')
const paidAmount = ref<number>(0)
const showReceiptModal = ref(false)
const showHistoryModal = ref(false)
const showQRISModal = ref(false)
const receiptText = ref('')

const categories = ['Semua', 'Oli & Pelumas', 'Suku Cadang', 'Busi & Kelistrikan']

onMounted(async () => {
  await Promise.all([
    sparepartStore.fetchSpareparts(),
    retailStore.fetchSales(),
    retailStore.fetchPromotions(),
  ])
  setExactCash()
})

const filteredParts = computed(() => {
  return sparepartStore.spareparts.filter(p => {
    // Exclude services from OTC retail parts catalog
    if (p.category?.name?.toLowerCase().includes('jasa')) return false

    const matchesSearch =
      p.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      p.part_number.toLowerCase().includes(searchQuery.value.toLowerCase())

    if (selectedCategory.value === 'Semua') return matchesSearch
    if (selectedCategory.value === 'Oli & Pelumas') {
      return matchesSearch && (p.name.toLowerCase().includes('oli') || p.name.toLowerCase().includes('oil'))
    }
    if (selectedCategory.value === 'Busi & Kelistrikan') {
      return matchesSearch && (p.name.toLowerCase().includes('busi') || p.name.toLowerCase().includes('battery'))
    }
    return matchesSearch
  })
})

const changeAmount = computed(() => {
  const diff = paidAmount.value - retailStore.cartGrandTotal
  return diff > 0 ? diff : 0
})

function onSearch() {
  playClickSound()
}

function handleAddToCart(part: Sparepart) {
  playScannerBeep()
  retailStore.addToCart(part)
  setExactCash()
}

function setExactCash() {
  paidAmount.value = Math.ceil(retailStore.cartGrandTotal)
}

function handleMethodChange(method: string) {
  playClickSound()
  if (method === 'cash') {
    setExactCash()
  } else if (method === 'qris') {
    paidAmount.value = Math.ceil(retailStore.cartGrandTotal)
  } else {
    paidAmount.value = Math.ceil(retailStore.cartGrandTotal)
  }
}

function handleApplyPromo() {
  if (!promoInput.value) return
  const res = retailStore.applyPromoCode(promoInput.value)
  if (res.success) {
    toast.success(res.message)
    playScannerBeep()
    setExactCash()
  } else {
    toast.error(res.message)
  }
}

async function handleCheckout() {
  if (paymentMethod.value === 'qris') {
    showQRISModal.value = true
    return
  }

  await executeCheckout()
}

async function onQRISVerified() {
  showQRISModal.value = false
  await executeCheckout()
}

async function executeCheckout() {
  try {
    const sale = await retailStore.checkout({
      customer_name: customerName.value,
      customer_phone: customerPhone.value,
      payment_method: paymentMethod.value,
      paid_amount: paidAmount.value,
    })

    playSuccessChime()
    toast.success(`Transaksi ${sale.invoice_number} berhasil lunas!`)

    // Load receipt
    const rData = await retailStore.fetchReceipt(sale.id)
    receiptText.value = rData.receipt_text
    showReceiptModal.value = true
  } catch (err: any) {
    toast.error(err.message || 'Gagal memproses transaksi retail')
  }
}

async function openReceiptForSale(saleId: number) {
  try {
    const rData = await retailStore.fetchReceipt(saleId)
    receiptText.value = rData.receipt_text
    showReceiptModal.value = true
  } catch (err: any) {
    toast.error('Gagal mengambil struk')
  }
}

function printReceipt() {
  window.print()
}
</script>
